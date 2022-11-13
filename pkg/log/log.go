package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/beihai0xff/pudding/configs"
)

var (
	defaultLogger *logger
	defaultConfig = &configs.LogConfig{
		Writers:    []string{configs.OutputConsole},
		Format:     configs.EncoderTypeConsole,
		Level:      "debug",
		CallerSkip: 1,
	}

	loggers = map[string]Logger{}
)

const DefaultLoggerName = "default"

func init() {
	defaultLogger = newLogger(defaultConfig)
	loggers[DefaultLoggerName] = defaultLogger
}

type logger struct {
	*zap.SugaredLogger
}

// WithFields add customs fields to logger
func (l *logger) WithFields(fields ...interface{}) Logger {
	return &logger{l.WithOptions(zap.AddStacktrace(zapcore.WarnLevel)).With(fields...)}
}

func RegisterLogger(loggerName string, opts ...OptionFunc) {
	c := configs.GetLogConfig(loggerName)

	for _, opt := range opts {
		opt(c)
	}

	log := newLogger(c)
	if loggerName == DefaultLoggerName {
		defaultLogger = log
	}
	loggers[loggerName] = log
}

func GerLoggerByName(loggerName string) Logger {
	if logger, ok := loggers[loggerName]; ok {
		return logger
	}
	Warnf("logger %s not found, use default logger", loggerName)
	return defaultLogger
}

// OptionFunc is the option function for LogConfig
type OptionFunc func(config *configs.LogConfig)

// WithCallerSkip set caller skip
func WithCallerSkip(callerSkip int) OptionFunc {
	return func(c *configs.LogConfig) {
		c.CallerSkip = callerSkip
	}
}
