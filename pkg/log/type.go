package log

import "go.uber.org/zap/zapcore"

// Levels zapcore level
var Levels = map[string]zapcore.Level{
	"":      zapcore.DebugLevel,
	"debug": zapcore.DebugLevel,
	"info":  zapcore.InfoLevel,
	"warn":  zapcore.WarnLevel,
	"error": zapcore.ErrorLevel,
	"fatal": zapcore.FatalLevel,
}

const (
	// OutputConsole 控制台输出日志
	OutputConsole = "console"
	// OutputFile 文件输出日志
	OutputFile = "file"

	// EncoderTypeConsole 日志输出格式：控制台
	EncoderTypeConsole = "console"
	// EncoderTypeJson 日志输出格式：json
	EncoderTypeJson = "json"
)
