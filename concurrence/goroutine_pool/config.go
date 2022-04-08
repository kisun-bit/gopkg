package goroutine_pool

const defaultScalaThreshold = 1

type Config struct {
	// 当等待的任务数大于ScaleThreshold时，就启动新的goroutine
	ScaleThreshold int32
}

func NewDefaultConfig() *Config {
	return &Config{ScaleThreshold: defaultScalaThreshold}
}
