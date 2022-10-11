package log

import "go.uber.org/zap"

var defaultLogger *zap.SugaredLogger

func init() {
	defaultLogger = NewLogger(defaultConfig).Sugar()
}
