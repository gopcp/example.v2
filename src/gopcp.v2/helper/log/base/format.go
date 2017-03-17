package base

// LogFormat 表示日志格式的类型。
type LogFormat string

const (
	// FORMAT_TEXT 代表普通文本的日志格式。
	FORMAT_TEXT LogFormat = "text"
	// FORMAT_JSON 代表JSON的日志格式。
	FORMAT_JSON LogFormat = "json"
)

const (
	// TIMESTAMP_FORMAT 代表时间戳格式化字符串。
	TIMESTAMP_FORMAT = "2006-01-02T15:04:05.999"
)
