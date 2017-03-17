package log

import (
	"fmt"
	"io"
	"os"
	"sync"

	"gopcp.v2/helper/log/base"
	"gopcp.v2/helper/log/logrus"
)

// LoggerCreator 代表日志记录器的创建器。
type LoggerCreator func(
	level base.LogLevel,
	format base.LogFormat,
	writer io.Writer,
	options []base.Option) base.MyLogger

// loggerCreatorMap 代表日志记录器创建器的映射。
var loggerCreatorMap = map[base.LoggerType]LoggerCreator{}

// rwm 代表日志记录器创建器映射的专用锁。
var rwm sync.RWMutex

// RegisterLogger 用于注册日志记录器。
func RegisterLogger(
	loggerType base.LoggerType,
	creator LoggerCreator,
	cover bool) error {
	if loggerType == "" {
		return fmt.Errorf("logger register error: invalid logger type")
	}
	if creator == nil {
		return fmt.Errorf("logger register error: invalid logger creator (logger type: %s)", loggerType)
	}
	rwm.Lock()
	defer rwm.Unlock()
	if _, ok := loggerCreatorMap[loggerType]; ok || !cover {
		return fmt.Errorf("logger register error: already existing logger for type %q", loggerType)
	}
	loggerCreatorMap[loggerType] = creator
	return nil
}

// DLogger 会返回一个新的默认日志记录器。
func DLogger() base.MyLogger {
	return Logger(
		base.TYPE_LOGRUS,
		base.LEVEL_INFO,
		base.FORMAT_TEXT,
		os.Stdout,
		nil)
}

// Logger 会新建一个日志记录器。
func Logger(
	loggerType base.LoggerType,
	level base.LogLevel,
	format base.LogFormat,
	writer io.Writer,
	options []base.Option) base.MyLogger {
	var logger base.MyLogger
	rwm.RLock()
	creator, ok := loggerCreatorMap[loggerType]
	rwm.RUnlock()
	if ok {
		logger = creator(level, format, writer, options)
	} else {
		logger = logrus.NewLoggerBy(level, format, writer, options)
	}
	return logger
}
