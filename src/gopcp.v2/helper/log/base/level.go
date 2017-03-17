package base

// LoggerLevel 代表日志输出级别。
type LogLevel uint8

const (
	// LEVEL_DEBUG 代表调试级别，是最低的日志等级。
	LEVEL_DEBUG LogLevel = iota + 1
	// LEVEL_INFO 代表信息级别，是最常用的日志等级。
	LEVEL_INFO
	// LEVEL_WARN 代表警告级别，是适合输出到错误输出的日志等级。
	LEVEL_WARN
	// LEVEL_ERROR 代表普通错误级别，是建议输出到错误输出的日志等级。
	LEVEL_ERROR
	// LEVEL_FATAL 代表致命错误级别，是建议输出到错误输出的日志等级。
	// 此级别的日志一旦输出就意味着`os.Exit(1)`立即会被调用。
	LEVEL_FATAL
	// LEVEL_PANIC 代表恐慌级别，是最高的日志等级。
	// 此级别的日志一旦输出就意味着运行时恐慌立即会被引发。
	LEVEL_PANIC
)
