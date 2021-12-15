package logging

import (
	"gopkg.in/ini.v1"
	"strconv"
	"strings"
)

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

func __getCfgKey(__cfg *ini.File, section, key string) string {
	return __cfg.Section(section).Key(key).String()
}
