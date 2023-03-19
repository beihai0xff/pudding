package logger

import (
	"sync"

	"github.com/beihai0xff/pudding/pkg/log"
)

func getMessageLog() log.Logger {
	return log.GetLoggerByName(KafkaLoggerName).WithFields("module", "kafka")
}

var (
	messageLoggerOnce sync.Once
	messageLogger     *MessageLogger
)

// MessageLogger 日志结构体
type MessageLogger struct {
	l log.Logger
}

// NewMessageLogger 构造一个自定义 Message Logger
func NewMessageLogger() *MessageLogger {
	messageLoggerOnce.Do(func() {
		messageLogger = &MessageLogger{
			l: getMessageLog(),
		}
	})
	return messageLogger
}

// RecordMessageInfoLog print Info messages
func (l *MessageLogger) RecordMessageInfoLog(format string, args ...interface{}) {
	l.l.Infof(format, args...)
}

// RecordMessageErrorLog print Error messages
func (l *MessageLogger) RecordMessageErrorLog(format string, args ...interface{}) {
	l.l.Errorf(format, args...)
}
