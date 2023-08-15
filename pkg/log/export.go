package log

import (
	"errors"
	"syscall"
)

// Logger is the interface that wraps the basic log methods.
type Logger interface {
	// Debug logs to DEBUG log
	Debug(args ...any)
	// Info logs to INFO log
	Info(args ...any)
	// Warn logs to WARN log
	Warn(args ...any)
	// Error logs to ERROR log
	Error(args ...any)
	// Fatal logs to FATAL log
	Fatal(args ...any)
	// Debugln logs to DEBUG log
	Debugln(args ...any)
	// Infoln logs to INFO log
	Infoln(args ...any)
	// Warnln logs to WARN log
	Warnln(args ...any)
	// Errorln logs to ERROR log
	Errorln(args ...any)
	// Fatalln logs to FATAL log
	Fatalln(args ...any)
	// Debugf logs to DEBUG log, Arguments will be formatted with Sprint, Sprintf, or neither.
	Debugf(format string, args ...any)
	// Infof logs to INFO log, Arguments will be formatted with Sprint, Sprintf, or neither.
	Infof(format string, args ...any)
	// Warnf logs to WARN log, Arguments will be formatted with Sprint, Sprintf, or neither.
	Warnf(format string, args ...any)
	// Errorf logs to ERROR log, Arguments will be formatted with Sprint, Sprintf, or neither.
	Errorf(format string, args ...any)
	// Fatalf logs to FATAL log, Arguments will be formatted with Sprint, Sprintf, or neither.
	Fatalf(format string, args ...any)
	// Sync calls the zap defaultLogger's Sync method, flushing any buffered log entries.
	// Applications should take care to call Sync before exiting.
	Sync() error
	// WithFields set some business custom data to each log
	WithFields(fields ...any) Logger
}

// Debug logs to DEBUG log
func Debug(args ...any) {
	defaultLogger.Debug(args...)
}

// Info logs to INFO log
func Info(args ...any) {
	defaultLogger.Info(args...)
}

// Warn logs to WARN log
func Warn(args ...any) {
	defaultLogger.Warn(args...)
}

// Error logs to ERROR log
func Error(args ...any) {
	defaultLogger.Error(args...)
}

// Fatal logs to FATAL log
func Fatal(args ...any) {
	defaultLogger.Fatal(args...)
}

// Debugf logs to DEBUG log, Arguments will be formatted with Sprint, Sprintf, or neither.
func Debugf(format string, args ...any) {
	defaultLogger.Debugf(format, args...)
}

// Infof logs to INFO log, Arguments will be formatted with Sprint, Sprintf, or neither.
func Infof(format string, args ...any) {
	defaultLogger.Infof(format, args...)
}

// Warnf logs to WARN log, Arguments will be formatted with Sprint, Sprintf, or neither.
func Warnf(format string, args ...any) {
	defaultLogger.Warnf(format, args...)
}

// Errorf logs to ERROR log, Arguments will be formatted with Sprint, Sprintf, or neither.
func Errorf(format string, args ...any) {
	defaultLogger.Errorf(format, args...)
}

// Panicf logs to PANIC log, Arguments will be formatted with Sprint, Sprintf, or neither.
func Panicf(format string, args ...any) {
	defaultLogger.Panicf(format, args...)
}

// Fatalf logs to FATAL log, Arguments will be formatted with Sprint, Sprintf, or neither.
func Fatalf(format string, args ...any) {
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

// WithFields set some business custom data to each log
func WithFields(fields ...any) Logger {
	return defaultLogger.WithFields(fields...)
}
