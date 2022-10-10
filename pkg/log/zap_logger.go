package log

import (
	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// zapLog 基于 zap 的 Logger 实现
type zapLog struct {
	level  zap.AtomicLevel
	logger *zap.Logger
}

func getLogMsg(args ...interface{}) string {
	return fmt.Sprint(args...)
}

func getLogMsgf(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

// Trace logs to TRACE log, Arguments are handled in the manner of fmt.Print
func (l *zapLog) Trace(args ...interface{}) {
	if l.logger.Core().Enabled(zapcore.DebugLevel) {
		l.logger.Debug(getLogMsg(args...))
	}
}

// Tracef logs to TRACE log, Arguments are handled in the manner of fmt.Printf
func (l *zapLog) Tracef(format string, args ...interface{}) {
	if l.logger.Core().Enabled(zapcore.DebugLevel) {
		l.logger.Debug(getLogMsgf(format, args...))
	}
}

// Debug logs to DEBUG log, Arguments are handled in the manner of fmt.Print
func (l *zapLog) Debug(args ...interface{}) {
	if l.logger.Core().Enabled(zapcore.DebugLevel) {
		l.logger.Debug(getLogMsg(args...))
	}
}

// Debugf logs to DEBUG log, Arguments are handled in the manner of fmt.Printf
func (l *zapLog) Debugf(format string, args ...interface{}) {
	if l.logger.Core().Enabled(zapcore.DebugLevel) {
		l.logger.Debug(getLogMsgf(format, args...))
	}
}

// Info logs to INFO log, Arguments are handled in the manner of fmt.Print
func (l *zapLog) Info(args ...interface{}) {
	if l.logger.Core().Enabled(zapcore.InfoLevel) {
		l.logger.Info(getLogMsg(args...))
	}
}

// Infof logs to INFO log, Arguments are handled in the manner of fmt.Printf
func (l *zapLog) Infof(format string, args ...interface{}) {
	if l.logger.Core().Enabled(zapcore.InfoLevel) {
		l.logger.Info(getLogMsgf(format, args...))
	}
}

// Warn logs to WARNING log, Arguments are handled in the manner of fmt.Print
func (l *zapLog) Warn(args ...interface{}) {
	if l.logger.Core().Enabled(zapcore.WarnLevel) {
		l.logger.Warn(getLogMsg(args...))
	}
}

// Warnf logs to WARNING log, Arguments are handled in the manner of fmt.Printf
func (l *zapLog) Warnf(format string, args ...interface{}) {
	if l.logger.Core().Enabled(zapcore.WarnLevel) {
		l.logger.Warn(getLogMsgf(format, args...))
	}
}

// Error logs to ERROR log, Arguments are handled in the manner of fmt.Print
func (l *zapLog) Error(args ...interface{}) {
	if l.logger.Core().Enabled(zapcore.ErrorLevel) {
		l.logger.Error(getLogMsg(args...))
	}
}

// Errorf logs to ERROR log, Arguments are handled in the manner of fmt.Printf
func (l *zapLog) Errorf(format string, args ...interface{}) {
	if l.logger.Core().Enabled(zapcore.ErrorLevel) {
		l.logger.Error(getLogMsgf(format, args...))
	}
}

// Fatal logs to FATAL log, Arguments are handled in the manner of fmt.Print
func (l *zapLog) Fatal(args ...interface{}) {
	if l.logger.Core().Enabled(zapcore.FatalLevel) {
		l.logger.Fatal(getLogMsg(args...))
	}
}

// Fatalf logs to FATAL log, Arguments are handled in the manner of fmt.Printf
func (l *zapLog) Fatalf(format string, args ...interface{}) {
	if l.logger.Core().Enabled(zapcore.FatalLevel) {
		l.logger.Fatal(getLogMsgf(format, args...))
	}
}

// Sync calls the zap logger's Sync method, flushing any buffered log entries.
// Applications should take care to call Sync before exiting.
func (l *zapLog) Sync() error {
	return l.logger.Sync()
}

// WithFields 设置一些业务自定义数据到每条 log 中
func (l *zapLog) WithFields(fields ...string) *zap.Logger {

	zapFields := make([]zap.Field, len(fields)/2)
	for index := range zapFields {
		zapFields[index] = zap.String(fields[2*index], fields[2*index+1])
	}

	return l.logger.WithOptions(zap.AddCallerSkip(-1)).With(zapFields...)
}
