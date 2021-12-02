package logging

import (
	"fmt"
	"go.uber.org/zap/zapcore"
	"gopkg.in/ini.v1"
	"strconv"
	"strings"
	"sync"
)

type __logKey struct {
	_logger    Log
	_flushFunc func() error
}

const (
	ClassRotateFile = "logging.NewFileRotatingLogger"
	ClassConsole    = "logging.NewConsoleStreamingLogger"
)

var (
	__containers           = make(map[string]*__logKey)
	__loggerContainersOnce sync.Once
	__l                    = new(zapcore.Level)
)

func __getCfgKey(__cfg *ini.File, section, key string) string {
	return __cfg.Section(section).Key(key).String()
}

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

func initLoggersContainers(conf string) error {
	cfg, err := ini.Load(conf)
	if err != nil {
		panic(err)
	}

	loggerKeys := __getCfgKey(cfg, "loggers", "keys")
	for _, logName := range strings.Split(loggerKeys, ",") {
		if _, ok := __containers[logName]; ok {
			return fmt.Errorf("replicated logger `%s`", logName)
		}

		__containers[logName] = new(__logKey)
		loggerSection := fmt.Sprintf("logger_%s", logName)
		loggerLevel := __convertStr2Level(__getCfgKey(cfg, loggerSection, "level"))
		loggerStackLevel := __convertStr2Level(__getCfgKey(cfg, loggerSection, "stack_level"))
		handlersKeys := __getCfgKey(cfg, loggerSection, "handler")

		var handlers []interface{}
		handlersKeysArr := strings.Split(handlersKeys, ",")
		for _, handlerName := range handlersKeysArr {
			if !__in(handlerName, handlersKeysArr) {
				return fmt.Errorf("why handler `%s` not in `handlers` section", handlerName)
			}

			handlerSection := fmt.Sprintf("handler_%s", handlerName)
			handlerClass := __getCfgKey(cfg, handlerSection, "class")
			handlerLevel := __convertStr2Level(__getCfgKey(cfg, handlerSection, "level"))

			switch handlerClass {
			case ClassRotateFile:
				fp := __getCfgKey(cfg, handlerSection, "log_file")
				ms, _ := strconv.Atoi(__getCfgKey(cfg, handlerSection, "max_size"))
				mb, _ := strconv.Atoi(__getCfgKey(cfg, handlerSection, "max_backups"))
				ma, _ := strconv.Atoi(__getCfgKey(cfg, handlerSection, "max_age"))
				handlers = append(handlers, NewRotateWriter(handlerLevel, fp, ms, mb, ma))
				continue
			case ClassConsole:
				handlers = append(handlers, NewConsoleWriter(handlerLevel))
				continue
			default:
				return fmt.Errorf("unsupported handler class `%s`,"+
					"only `logging.NewFileRotatingLogger` and `logging.NewConsoleStreamingLogger` are valid",
					handlerClass)
			}
		}

		if len(handlers) == 0 {
			return fmt.Errorf("why logger(`%s`) has no handlers", logName)
		}

		if logger, flushFunc, err := NewLogger(loggerLevel, loggerStackLevel,
			logName, "", false, EncodeConsole, handlers...); err != nil {
			return err
		} else {
			__containers[logName]._logger = logger
			__containers[logName]._flushFunc = flushFunc
		}
	}

	return nil
}
