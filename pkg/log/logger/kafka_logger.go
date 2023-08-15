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

// MessageLogger kafka logger
type MessageLogger struct {
	l log.Logger
}

// NewMessageLogger create a kafka Message Logger
func NewMessageLogger() *MessageLogger {
	messageLoggerOnce.Do(func() {
		messageLogger = &MessageLogger{
			l: getMessageLog(),
		}
	})

	return messageLogger
}

// RecordMessageInfoLog print Info messages
func (l *MessageLogger) RecordMessageInfoLog(format string, args ...any) {
	l.l.Infof(format, args...)
}

// RecordMessageErrorLog print Error messages
func (l *MessageLogger) RecordMessageErrorLog(format string, args ...any) {
	l.l.Errorf(format, args...)
}
