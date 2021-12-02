# `gopkg`-`logging`

> A high-performance logger broker library that integrates 
`zap` and `lumberjack` (inspired by `logging` module style in `python`), 
supports dynamic creation of loggers,
support created based on configuration files    

## Install

```shell script
go get github.com/kisunSea/gopkg
```

## Example

### create a console logger

```go
package main

import (
	"github.com/kisunSea/gopkg/logging"
    "time"
)

func main() {

    logger, sync, err := logging.NewConsoleStreamingLogger("demo",
        logging.DebugLevel, logging.WarnLevel, logging.DebugLevel)

    if err != nil || logger == nil || sync == nil {
        panic("failed to create console logger")
    }   

    defer sync()
    
    logger.Debug("logger start //////////////////////////////////////////////////////")
    logger.Debug("this is a debug message...")
    logger.Info("this is a info message...")
    logger.Warn("this is a warn message...")
    logger.Error("this is a error message...")

    logger.DebugF("this is a debug format message: level-`%s`", logging.DebugLevel)
    logger.InfoF("this is a info format message: level-`%s`", logging.InfoLevel)
    logger.WarnF("this is a warn format message: level-`%s`", logging.WarnLevel)
    logger.ErrorF("this is a error format message: level-`%s`", logging.ErrorLevel)

    logger.DebugW("this is a DebugW format message: ",
        "number", struct{ T time.Time }{T: time.Now()}, "total", 4)
    logger.InfoW("this is a InfoW format message: ",
        "number", struct{ T time.Time }{T: time.Now()}, "total", 4)
    logger.ErrorW("this is a ErrorF format message: ",
        "number", struct{ T time.Time }{T: time.Now()}, "total", 4)
    logger.WarnW("this is a WarnF format message: ",
        "number", struct{ T time.Time }{T: time.Now()}, "total", 4)
    logger.Debug("logger end   //////////////////////////////////////////////////////")

}
```

### create a rotate file logger

```go
package main

import (
	"github.com/kisunSea/gopkg/logging"
    "time"
)

func main() {

    logger, sync, err := logging.NewFileRotatingLogger(
        "logging-demoC", `D:\workspace\gopkg\_examples\logging\logger_test.log`,
        logging.DebugLevel, logging.WarnLevel, logging.DebugLevel, 30, 6, 30)

    if err != nil || logger == nil || sync == nil {
        panic("failed to create rotate file logger")
    }   
    
    defer sync()
    
    logger.Debug("logger start //////////////////////////////////////////////////////")
    logger.Debug("this is a debug message...")
    logger.Info("this is a info message...")
    logger.Warn("this is a warn message...")
    logger.Error("this is a error message...")

    logger.DebugF("this is a debug format message: level-`%s`", logging.DebugLevel)
    logger.InfoF("this is a info format message: level-`%s`", logging.InfoLevel)
    logger.WarnF("this is a warn format message: level-`%s`", logging.WarnLevel)
    logger.ErrorF("this is a error format message: level-`%s`", logging.ErrorLevel)

    logger.DebugW("this is a DebugW format message: ",
        "number", struct{ T time.Time }{T: time.Now()}, "total", 4)
    logger.InfoW("this is a InfoW format message: ",
        "number", struct{ T time.Time }{T: time.Now()}, "total", 4)
    logger.ErrorW("this is a ErrorF format message: ",
        "number", struct{ T time.Time }{T: time.Now()}, "total", 4)
    logger.WarnW("this is a WarnF format message: ",
        "number", struct{ T time.Time }{T: time.Now()}, "total", 4)
    logger.Debug("logger end   //////////////////////////////////////////////////////")

}
```

### create loggers with configuration file

D:\workspace\gopkg\_examples\logging\log.ini
```ini
[loggers]
keys = root,log1,console

[handlers]
keys = root_handler,log1_handler,console_handler,log2_handler

[logger_root]
level = debug
stack_level = error
handler = root_handler,console_handler

[logger_log1]
level = debug
stack_level = error
handler = log1_handler

[logger_console]
level = debug
stack_level = error
handler = console_handler

[handler_root_handler]
class = logging.NewFileRotatingLogger
log_file = D:\workspace\gopkg\_examples\logging\logger_root.log
max_age = 30
max_size = 30
max_backups = 6
level = debug

[handler_log1_handler]
class = logging.NewFileRotatingLogger
log_file = D:\workspace\gopkg\_examples\logging\logger_log1.log
max_age = 30
max_size = 30
max_backups = 6
level = info

[handler_console_handler]
class = logging.NewConsoleStreamingLogger
level = warn
```
main.go

```go
package main

import (
	"github.com/kisunSea/gopkg/logging"
    "time"
)

func main() {

    logger, sync, err := logging.GetLoggerFromConf(
        "root", // `log1`, `console`
        `D:\workspace\gopkg\_examples\logging\log.ini`)

    if err != nil || logger == nil || sync == nil {
        panic("failed to create config logger")
    }   
    
    defer sync()
    
    logger.Debug("logger start //////////////////////////////////////////////////////")
    logger.Debug("this is a debug message...")
    logger.Info("this is a info message...")
    logger.Warn("this is a warn message...")
    logger.Error("this is a error message...")

    logger.DebugF("this is a debug format message: level-`%s`", logging.DebugLevel)
    logger.InfoF("this is a info format message: level-`%s`", logging.InfoLevel)
    logger.WarnF("this is a warn format message: level-`%s`", logging.WarnLevel)
    logger.ErrorF("this is a error format message: level-`%s`", logging.ErrorLevel)

    logger.DebugW("this is a DebugW format message: ",
        "number", struct{ T time.Time }{T: time.Now()}, "total", 4)
    logger.InfoW("this is a InfoW format message: ",
        "number", struct{ T time.Time }{T: time.Now()}, "total", 4)
    logger.ErrorW("this is a ErrorF format message: ",
        "number", struct{ T time.Time }{T: time.Now()}, "total", 4)
    logger.WarnW("this is a WarnF format message: ",
        "number", struct{ T time.Time }{T: time.Now()}, "total", 4)
    logger.Debug("logger end   //////////////////////////////////////////////////////")

}
```