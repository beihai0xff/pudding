package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/beihai0xff/pudding/pkg/configs"
)

var defaultLogger *logger

var defaultConfig = &configs.OutputConfig{
	Writers:    []string{configs.OutputConsole},
	Formatter:  configs.EncoderTypeConsole,
	Level:      "debug",
	CallerSkip: 1,
}

func init() {
	defaultLogger = NewLog(defaultConfig)
}

type logger struct {
	*zap.SugaredLogger
}

// WithFields add customs fields to logger
func (l *logger) WithFields(fields ...interface{}) Logger {
	return &logger{l.WithOptions(zap.AddStacktrace(zapcore.WarnLevel)).With(fields...)}
}
