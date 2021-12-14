package logging

import (
	"errors"
	"fmt"
	"go.uber.org/zap/zapcore"
	"gopkg.in/ini.v1"
	"strconv"
	"strings"
	"sync"
)

var (
	__l = new(zapcore.Level)
	_lc = NewLoggerContainer()
)

func __convertStr2Level(levelStr string) Level {
	if err := __l.Set(levelStr); err != nil {
		panic(err)
	}
	return __l.Get().(zapcore.Level)
}

func __in(target string, origin []string) bool {
	for _, o := range origin {
		if o == target {
			return true
		}
	}
	return false
}

func __getCfgKey(__cfg *ini.File, section, key string) string {
	return __cfg.Section(section).Key(key).String()
}

// parser ...

const (
	SectionLoggers              = "loggers"
	SectionLoggersValKeys       = "keys"
	SectionHandlers             = "handlers"
	SectionHandlersValKeys      = "keys"
	SectionLoggerPrefix         = "logger_"
	SectionLoggerValLevel       = "level"
	SectionLoggerValStackLevel  = "stack_level"
	SectionLoggerValHandler     = "handler"
	SectionHandlerPrefix        = "handler_"
	SectionHandlerValClass      = "class"
	SectionHandlerValLevel      = "level"
	SectionHandlerValMaxAge     = "max_age"
	SectionHandlerValMaxSize    = "max_size"
	SectionHandlerValMaxBackups = "max_backups"
	SectionHandlerValLogFile    = "log_file"

	ClassRotateFile = "logging.NewFileRotatingLogger"
	ClassConsole    = "logging.NewConsoleStreamingLogger"
)

type confParser struct {
	conf  string
	iniFp *ini.File
}

func NewConfParser(conf string) (c *confParser, err error) {
	c = new(confParser)
	c.conf = conf

	if c.iniFp, err = ini.Load(conf); err != nil {
		return nil, err
	}
	return c, nil
}

func (c *confParser) LoggerKeys() []string {
	return strings.Split(
		__getCfgKey(c.iniFp, SectionLoggers, SectionLoggersValKeys), ",")
}

func (c *confParser) HandlerKeys() []string {
	return strings.Split(
		__getCfgKey(c.iniFp, SectionHandlers, SectionHandlersValKeys), ",")
}

func (c *confParser) ValLoggerLevel(loggerKey string) Level {
	return __convertStr2Level(
		__getCfgKey(c.iniFp, SectionLoggerPrefix+loggerKey, SectionLoggerValLevel))
}

func (c *confParser) ValLoggerStackLevel(loggerKey string) Level {
	return __convertStr2Level(
		__getCfgKey(c.iniFp, SectionLoggerPrefix+loggerKey, SectionLoggerValStackLevel))
}

func (c *confParser) ValLoggerHandler(loggerKey string) []string {
	return strings.Split(
		__getCfgKey(c.iniFp, SectionLoggerPrefix+loggerKey, SectionLoggerValHandler), ",")
}

func (c *confParser) ValHandlerLogFile(handlerKey string) string {
	return __getCfgKey(c.iniFp, SectionHandlerPrefix+handlerKey, SectionHandlerValLogFile)
}

func (c *confParser) ValHandlerClass(handlerKey string) string {
	return __getCfgKey(c.iniFp, SectionHandlerPrefix+handlerKey, SectionHandlerValClass)
}

func (c *confParser) ValHandlerLevel(handlerKey string) Level {
	return __convertStr2Level(
		__getCfgKey(c.iniFp, SectionHandlerPrefix+handlerKey, SectionHandlerValLevel))
}

func (c *confParser) ValHandlerMaxAge(handlerKey string) int {
	_r := __getCfgKey(c.iniFp, SectionHandlerPrefix+handlerKey, SectionHandlerValMaxAge)
	if r, err := strconv.Atoi(_r); err == nil {
		return r
	}
	return 0
}

func (c *confParser) ValHandlerMaxSize(handlerKey string) int {
	_r := __getCfgKey(c.iniFp, SectionHandlerPrefix+handlerKey, SectionHandlerValMaxSize)
	if r, err := strconv.Atoi(_r); err == nil {
		return r
	}
	return 0
}

func (c *confParser) ValHandlerMaxBackups(handlerKey string) int {
	_r := __getCfgKey(c.iniFp, SectionHandlerPrefix+handlerKey, SectionHandlerValMaxBackups)
	if r, err := strconv.Atoi(_r); err == nil {
		return r
	}
	return 0
}

type LoggerPool struct {
	// __containers only initialize once (write), it's safe to load `logger` (read) multiple times later
	__containers           map[string]*Logger
	__loggerContainersOnce sync.Once
}

func NewLoggerContainer() (_ *LoggerPool) {
	l := new(LoggerPool)
	l.__containers = make(map[string]*Logger)
	return l
}

func (lc *LoggerPool) GetLogger(name string) (logger *Logger, err error) {
	if __k, ok := lc.__containers[name]; ok {
		if __k == nil {
			return nil, errors.New("nil logger")
		}
		return __k, nil
	}
	return nil, fmt.Errorf("failed to get logger named `%s`", name)
}

func SetConf(conf string) (_ *LoggerPool, err error) {
	if err := initLoggersContainers(conf); err != nil {
		return nil, fmt.Errorf("set logging conf `%s` failed", conf)
	}
	return _lc, nil
}

// initLoggersContainers initializes a pool of log objects
func initLoggersContainers(conf string) (err error) {

	var c *confParser
	if c, err = NewConfParser(conf); err != nil {
		return err
	}

	for _, loggerName := range c.LoggerKeys() {
		if _, ok := _lc.__containers[loggerName]; ok {
			return fmt.Errorf("replicated logger `%s`", loggerName)
		}
		_lc.__containers[loggerName] = new(Logger)

		var (
			handlersKeysArr = c.ValLoggerHandler(loggerName)
			handlers        = make([]interface{}, 0, len(handlersKeysArr))
		)

		for _, handlerName := range handlersKeysArr {
			if !__in(handlerName, c.HandlerKeys()) {
				return fmt.Errorf(
					"why handler `%s` not in `handlers` section", handlerName)
			}
			switch c.ValHandlerClass(handlerName) {
			case ClassRotateFile:
				handlers = append(handlers, NewRotateWriter(
					c.ValHandlerLevel(handlerName),
					c.ValHandlerLogFile(handlerName),
					c.ValHandlerMaxSize(handlerName),
					c.ValHandlerMaxBackups(handlerName),
					c.ValHandlerMaxAge(handlerName)))
			case ClassConsole:
				handlers = append(handlers, NewConsoleWriter(
					c.ValHandlerLevel(handlerName)))
			default:
				return errors.New("unsupported handler class `%s`," +
					"only `logging.NewFileRotatingLogger` and `logging.NewConsoleStreamingLogger` are valid")
			}
		}

		if len(handlers) == 0 {
			return fmt.Errorf("why logger(`%s`) has no handlers", loggerName)
		}

		if logger, err := NewLogger(
			c.ValLoggerLevel(loggerName),
			c.ValLoggerStackLevel(loggerName),
			loggerName, "", false, EncodeConsole, handlers...); err != nil {
			return err
		} else {
			_lc.__containers[loggerName] = logger
		}
	}

	return nil
}
