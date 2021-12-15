package logging

import "go.uber.org/zap/zapcore"

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
