package log

var logger Logger

func init() {
	logger = NewLogger(defaultConfig)
}
