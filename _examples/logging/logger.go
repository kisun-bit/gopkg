package main

import (
	"github.com/kisunSea/gopkg/logging"
	"time"
)

func ExampleFileRotateLogger() {

	// logging.NewLogger(logging.DebugLevel, logging.ErrorLevel,"custom-logger", "", true, logging.EncodeConsole)

	// logging.NewConsoleStreamingLogger("console-logger", logging.DebugLevel, logging.WarnLevel, logging.DebugLevel)

	logger, err := logging.NewFileRotatingLogger(
		"file-rotate-logger", `D:\workspace\gopkg\_examples\logging\logger_test.log`,
		logging.DebugLevel, logging.WarnLevel, logging.DebugLevel, 30, 6, 30)
	if err != nil {
		logging.GLogger().FatalF("failed to call `GetFileRotatingLogger`: %s", err)
	}
	defer logger.Sync()

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

// configured mode to get logger ...

func ExampleConfLoggers() {
	var (
		lp         *logging.LoggerPool
		loggerRoot *logging.Logger
		loggerLog1 *logging.Logger
		err        error
	)

	if lp, err = logging.SetConf(`D:\workspace\gopkg\_examples\logging\log.ini`); err != nil {
		logging.GLogger().FatalF("call `logging.SetConf` failed: %v", err)
	}

	if loggerRoot, err = lp.GetLogger("root"); err != nil {
		logging.GLogger().FatalF("call `logging.GetLogger` failed: %v", err)
	}

	if loggerLog1, err = lp.GetLogger("log1"); err != nil {
		logging.GLogger().FatalF("call `logging.GetLogger` failed: %v", err)
	}

	defer loggerRoot.Sync()
	defer loggerLog1.Sync()

	loggerRoot.Debug("loggerRoot start //////////////////////////////////////////////////////")
	loggerRoot.Debug("this is a debug message...")
	loggerRoot.Info("this is a info message...")
	loggerRoot.Warn("this is a warn message...")
	loggerRoot.Error("this is a error message...")
	loggerRoot.DebugF("this is a debug format message: level-`%s`", logging.DebugLevel)
	loggerRoot.InfoF("this is a info format message: level-`%s`", logging.InfoLevel)
	loggerRoot.WarnF("this is a warn format message: level-`%s`", logging.WarnLevel)
	loggerRoot.ErrorF("this is a error format message: level-`%s`", logging.ErrorLevel)
	loggerRoot.DebugW("this is a DebugW format message: ",
		"number", struct{ T time.Time }{T: time.Now()}, "total", 4)
	loggerRoot.InfoW("this is a InfoW format message: ",
		"number", struct{ T time.Time }{T: time.Now()}, "total", 4)
	loggerRoot.ErrorW("this is a ErrorF format message: ",
		"number", struct{ T time.Time }{T: time.Now()}, "total", 4)
	loggerRoot.WarnW("this is a WarnF format message: ",
		"number", struct{ T time.Time }{T: time.Now()}, "total", 4)
	loggerRoot.Debug("loggerRoot end   //////////////////////////////////////////////////////")

	loggerLog1.DebugW("this is a DebugW format message: ",
		"number", struct{ T time.Time }{T: time.Now()}, "total", 4)
	loggerLog1.InfoW("this is a InfoW format message: ",
		"number", struct{ T time.Time }{T: time.Now()}, "total", 4)
	loggerLog1.ErrorW("this is a ErrorF format message: ",
		"number", struct{ T time.Time }{T: time.Now()}, "total", 4)
	loggerLog1.WarnW("this is a WarnF format message: ",
		"number", struct{ T time.Time }{T: time.Now()}, "total", 4)
	loggerLog1.Debug("loggerRoot end   //////////////////////////////////////////////////////")
}

func main() {
	//ExampleFileRotateLogger()
	ExampleConfLoggers()
}
