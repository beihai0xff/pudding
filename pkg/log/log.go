// Package log provides the log
package log

import (
	"encoding/json"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	defaultLogger *logger

	loggers = map[string]Logger{}
)

// DefaultLoggerName is the default logger name
const DefaultLoggerName = "default"

func init() {
	defaultLogger = newLogger(&DefaultConfig)
	loggers[DefaultLoggerName] = defaultLogger
}

type logger struct {
	*zap.SugaredLogger
}

// WithFields add customs fields to logger
func (l *logger) WithFields(fields ...any) Logger {
	return &logger{l.WithOptions(zap.AddStacktrace(zapcore.WarnLevel)).With(fields...)}
}

// RegisterLogger register a logger with name
func RegisterLogger(loggerName string, c *Config, opts ...OptionFunc) {
	cjson, _ := json.Marshal(c)
	Infof("Register Logger [%s] with config: %s", loggerName, string(cjson))

	for _, opt := range opts {
		opt(c)
	}

	log := newLogger(c)
	if loggerName == DefaultLoggerName {
		defaultLogger = log
	}

	loggers[loggerName] = log
}

// GetLoggerByName get logger by name
func GetLoggerByName(loggerName string) Logger {
	if logger, ok := loggers[loggerName]; ok {
		return logger
	}

	Warnf("logger %s not found, use default logger", loggerName)

	return defaultLogger
}

// OptionFunc is the option function for LogConfig
type OptionFunc func(config *Config)

// WithCallerSkip set caller skip
func WithCallerSkip(callerSkip int) OptionFunc {
	return func(c *Config) {
		c.CallerSkip = callerSkip
	}
}
