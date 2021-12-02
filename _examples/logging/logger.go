package main

import (
	"github.com/kisunSea/gopkg/logging"
	"time"
)

func main() {

	loggerA, syncA, errA := logging.NewLogger(logging.DebugLevel, logging.WarnLevel,
		"logging-demoA", "",
		true, logging.EncodeConsole)
	if errA != nil {
		panic(errA)
	}

	defer func() {
		loggerA.DPanic("catch th panic stack...")
		_ = syncA()
	}()

	loggerA.Debug("loggerA start //////////////////////////////////////////////////////")
	loggerA.Debug("this is a debug message...")
	loggerA.Info("this is a info message...")
	loggerA.Warn("this is a warn message...")
	loggerA.Error("this is a error message...")

	loggerA.DebugF("this is a debug format message: level-`%s`", logging.DebugLevel)
	loggerA.InfoF("this is a info format message: level-`%s`", logging.InfoLevel)
	loggerA.WarnF("this is a warn format message: level-`%s`", logging.WarnLevel)
	loggerA.ErrorF("this is a error format message: level-`%s`", logging.ErrorLevel)

	loggerA.DebugW("this is a DebugW format message: ",
		"number", struct{ T time.Time }{T: time.Now()}, "total", 4)
	loggerA.InfoW("this is a InfoW format message: ",
		"number", struct{ T time.Time }{T: time.Now()}, "total", 4)
	loggerA.ErrorW("this is a ErrorF format message: ",
		"number", struct{ T time.Time }{T: time.Now()}, "total", 4)
	loggerA.WarnW("this is a WarnF format message: ",
		"number", struct{ T time.Time }{T: time.Now()}, "total", 4)
	loggerA.Debug("loggerA end   //////////////////////////////////////////////////////")

	//////////////////////////////////////////////////

	loggerB, syncB, errB := logging.NewConsoleStreamingLogger("logging-demoB",
		logging.DebugLevel, logging.WarnLevel, logging.DebugLevel)
	if errB != nil {
		panic(errB)
	}
	defer syncB()

	loggerB.Debug("loggerB start //////////////////////////////////////////////////////")
	loggerB.Debug("this is a debug message...")
	loggerB.Info("this is a info message...")
	loggerB.Warn("this is a warn message...")
	loggerB.Error("this is a error message...")

	loggerB.DebugF("this is a debug format message: level-`%s`", logging.DebugLevel)
	loggerB.InfoF("this is a info format message: level-`%s`", logging.InfoLevel)
	loggerB.WarnF("this is a warn format message: level-`%s`", logging.WarnLevel)
	loggerB.ErrorF("this is a error format message: level-`%s`", logging.ErrorLevel)

	loggerB.DebugW("this is a DebugW format message: ",
		"number", struct{ T time.Time }{T: time.Now()}, "total", 4)
	loggerB.InfoW("this is a InfoW format message: ",
		"number", struct{ T time.Time }{T: time.Now()}, "total", 4)
	loggerB.ErrorW("this is a ErrorF format message: ",
		"number", struct{ T time.Time }{T: time.Now()}, "total", 4)
	loggerB.WarnW("this is a WarnF format message: ",
		"number", struct{ T time.Time }{T: time.Now()}, "total", 4)
	loggerB.Debug("loggerB end   //////////////////////////////////////////////////////")

	loggerC, syncC, errC := logging.NewFileRotatingLogger(
		"logging-demoC", `D:\workspace\gopkg\_examples\logging\logger_test.log`,
		logging.DebugLevel, logging.WarnLevel, logging.DebugLevel, 30, 6, 30)
	if errC != nil {
		panic(errC)
	}
	defer syncC()

	loggerC.Debug("loggerC start //////////////////////////////////////////////////////")
	loggerC.Debug("this is a debug message...")
	loggerC.Info("this is a info message...")
	loggerC.Warn("this is a warn message...")
	loggerC.Error("this is a error message...")

	loggerC.DebugF("this is a debug format message: level-`%s`", logging.DebugLevel)
	loggerC.InfoF("this is a info format message: level-`%s`", logging.InfoLevel)
	loggerC.WarnF("this is a warn format message: level-`%s`", logging.WarnLevel)
	loggerC.ErrorF("this is a error format message: level-`%s`", logging.ErrorLevel)

	loggerC.DebugW("this is a DebugW format message: ",
		"number", struct{ T time.Time }{T: time.Now()}, "total", 4)
	loggerC.InfoW("this is a InfoW format message: ",
		"number", struct{ T time.Time }{T: time.Now()}, "total", 4)
	loggerC.ErrorW("this is a ErrorF format message: ",
		"number", struct{ T time.Time }{T: time.Now()}, "total", 4)
	loggerC.WarnW("this is a WarnF format message: ",
		"number", struct{ T time.Time }{T: time.Now()}, "total", 4)
	loggerC.Debug("loggerC end   //////////////////////////////////////////////////////")

	logger_, ok := loggerC.(*logging.Logger)
	if ok {
		logger_.SetStacktrace(logging.PanicLevel)
		logger_.SetStacktrace(logging.DPanicLevel)
		logger_.SetStacktrace(logging.FatalLevel)
	}

	defaultLogger, flush := logging.GLogger()
	defer flush()

	loggerD, syncD, errD := logging.GetLoggerFromConf(
		"root",
		`D:\workspace\gopkg\_examples\logging\log.ini`)
	if errD != nil || syncD == nil || loggerD == nil {
		defaultLogger.ErrorF("err in %s", defaultLogger.Stack())
		return
	}
	defer syncD()
	//if loggerD == nil {
	//	panic("1111")
	//}

	loggerD.Debug("loggerD start //////////////////////////////////////////////////////")
	loggerD.Debug("this is a debug message...")
	loggerD.Info("this is a info message...")
	loggerD.Warn("this is a warn message...")
	loggerD.Error("this is a error message...")

	loggerD.DebugF("this is a debug format message: level-`%s`", logging.DebugLevel)
	loggerD.InfoF("this is a info format message: level-`%s`", logging.InfoLevel)
	loggerD.WarnF("this is a warn format message: level-`%s`", logging.WarnLevel)
	loggerD.ErrorF("this is a error format message: level-`%s`", logging.ErrorLevel)

	loggerD.DebugW("this is a DebugW format message: ",
		"number", struct{ T time.Time }{T: time.Now()}, "total", 4)
	loggerD.InfoW("this is a InfoW format message: ",
		"number", struct{ T time.Time }{T: time.Now()}, "total", 4)
	loggerD.ErrorW("this is a ErrorF format message: ",
		"number", struct{ T time.Time }{T: time.Now()}, "total", 4)
	loggerD.WarnW("this is a WarnF format message: ",
		"number", struct{ T time.Time }{T: time.Now()}, "total", 4)
	loggerD.Debug("loggerD end   //////////////////////////////////////////////////////")
}
