package log

import (
	"errors"
	"syscall"
)

// Logger is the interface that wraps the basic log methods.
type Logger interface {
	// Debug logs to DEBUG log
	Debug(args ...interface{})
	// Info logs to INFO log
	Info(args ...interface{})
	// Warn logs to WARN log
	Warn(args ...interface{})
	// Error logs to ERROR log
	Error(args ...interface{})
	// Fatal logs to FATAL log
	Fatal(args ...interface{})
	// Debugln logs to DEBUG log
	Debugln(args ...interface{})
	// Infoln logs to INFO log
	Infoln(args ...interface{})
	// Warnln logs to WARN log
	Warnln(args ...interface{})
	// Errorln logs to ERROR log
	Errorln(args ...interface{})
	// Fatalln logs to FATAL log
	Fatalln(args ...interface{})
	// Debugf logs to DEBUG log, Arguments will be formatted with Sprint, Sprintf, or neither.
	Debugf(format string, args ...interface{})
	// Infof logs to INFO log, Arguments will be formatted with Sprint, Sprintf, or neither.
	Infof(format string, args ...interface{})
	// Warnf logs to WARN log, Arguments will be formatted with Sprint, Sprintf, or neither.
	Warnf(format string, args ...interface{})
	// Errorf logs to ERROR log, Arguments will be formatted with Sprint, Sprintf, or neither.
	Errorf(format string, args ...interface{})
	// Fatalf logs to FATAL log, Arguments will be formatted with Sprint, Sprintf, or neither.
	Fatalf(format string, args ...interface{})
	// Sync calls the zap defaultLogger's Sync method, flushing any buffered log entries.
	// Applications should take care to call Sync before exiting.
	Sync() error
	// WithFields set some business custom data to each log
	WithFields(fields ...interface{}) Logger
}

// Debug logs to DEBUG log
func Debug(args ...interface{}) {
	defaultLogger.Debug(args...)
}

// Info logs to INFO log
func Info(args ...interface{}) {
	defaultLogger.Info(args...)
}

// Warn logs to WARN log
func Warn(args ...interface{}) {
	defaultLogger.Warn(args...)
}

// Error logs to ERROR log
func Error(args ...interface{}) {
	defaultLogger.Error(args...)
}

// Fatal logs to FATAL log
func Fatal(args ...interface{}) {
	defaultLogger.Fatal(args...)
}

// Debugf logs to DEBUG log, Arguments will be formatted with Sprint, Sprintf, or neither.
func Debugf(format string, args ...interface{}) {
	defaultLogger.Debugf(format, args...)
}

// Infof logs to INFO log, Arguments will be formatted with Sprint, Sprintf, or neither.
func Infof(format string, args ...interface{}) {
	defaultLogger.Infof(format, args...)
}

// Warnf logs to WARN log, Arguments will be formatted with Sprint, Sprintf, or neither.
func Warnf(format string, args ...interface{}) {
	defaultLogger.Warnf(format, args...)
}

// Errorf logs to ERROR log, Arguments will be formatted with Sprint, Sprintf, or neither.
func Errorf(format string, args ...interface{}) {
	defaultLogger.Errorf(format, args...)
}

// Panicf logs to PANIC log, Arguments will be formatted with Sprint, Sprintf, or neither.
func Panicf(format string, args ...interface{}) {
	defaultLogger.Panicf(format, args...)
}

// Fatalf logs to FATAL log, Arguments will be formatted with Sprint, Sprintf, or neither.
func Fatalf(format string, args ...interface{}) {
	defaultLogger.Fatalf(format, args...)
}

// Sync calls the zap defaultLogger's Sync method, flushing any buffered log entries.
// Applications should take care to call Sync before exiting.
func Sync() {
	for loggerName, logger := range loggers {
		if err := logger.Sync(); err != nil {
			// https://github.com/uber-go/zap/issues/1026
			// Sync is not allowed on os.Stdout if it's being fed to a terminal.
			if !errors.Is(err, syscall.ENOTTY) {
				Errorf("sync logger [%s] error: %v", loggerName, err)
			}
		}
	}
}

// WithFields 设置一些业务自定义数据到每条 log 中
func WithFields(fields ...interface{}) Logger {
	return defaultLogger.WithFields(fields...)
}
