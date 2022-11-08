package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/beihai0xff/pudding/configs"
)

var defaultLogger *logger

var defaultConfig = &configs.LogConfig{
	Writers:    []string{configs.OutputConsole},
	Format:     configs.EncoderTypeConsole,
	Level:      "debug",
	CallerSkip: 1,
}

func init() {
	defaultLogger = newLog(defaultConfig)
	loggers["default"] = defaultLogger
}

type logger struct {
	*zap.SugaredLogger
}

// WithFields add customs fields to logger
func (l *logger) WithFields(fields ...interface{}) Logger {
	return &logger{l.WithOptions(zap.AddStacktrace(zapcore.WarnLevel)).With(fields...)}
}

var loggers = map[string]Logger{}

func RegisterLogger(logName string, opts ...OptionFunc) {
	c := configs.GetLogConfig(logName)

	for _, opt := range opts {
		opt(c)
	}
	loggers[logName] = newLog(c)
}

func GerLoggerByName(logName string) Logger {
	if logger, ok := loggers[logName]; ok {
		return logger
	}
	Warnf("logger %s not found, use default logger", logName)
	return defaultLogger
}

/*
	Functional Options Pattern
*/

type OptionFunc func(config *configs.LogConfig)

func WithCallerSkip(callerSkip int) OptionFunc {
	return func(c *configs.LogConfig) {
		c.CallerSkip = callerSkip
	}
}

type Option func(*configs.LogConfig)
