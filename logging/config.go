package logging

import (
	"errors"
	"fmt"
	"go.uber.org/zap/zapcore"
	"sync"
)

var (
	__l = new(zapcore.Level)
	_lc = NewLoggerContainer() // global loggers pool
)

func NewLoggerContainer() (_ *LoggerPool) {
	l := new(LoggerPool)
	l.__containers = make(map[string]*Logger)
	return l
}

type LoggerPool struct {
	// __containers only initialize once (write), it's safe to load `logger` (read) multiple times later
	__containers           map[string]*Logger
	__loggerContainersOnce sync.Once
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

func SetConf(conf string) (lp *LoggerPool, err error) {
	_lc.__loggerContainersOnce.Do(func() {
		if err := initLoggersContainers(conf); err != nil {
			lp, err = nil, fmt.Errorf("set logging conf `%s` failed", conf)
			return
		}
		lp, err = _lc, nil
	})
	return
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
				return fmt.Errorf("why handler `%s` not in `handlers` section", handlerName)
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
