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
}

type logger struct {
	*zap.SugaredLogger
}

// WithFields add customs fields to logger
func (l *logger) WithFields(fields ...interface{}) Logger {
	return &logger{l.WithOptions(zap.AddStacktrace(zapcore.WarnLevel)).With(fields...)}
}
