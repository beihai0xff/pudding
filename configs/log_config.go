package configs

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
	// OutputConsole console output
	OutputConsole = "console"
	// OutputFile file output
	OutputFile = "file"

	// EncoderTypeConsole log format console encoder
	EncoderTypeConsole = "console"
	// EncoderTypeJSON log format json encoder
	EncoderTypeJSON = "json"
)

// LogConfig log output: console file remote
type LogConfig struct {
	// Writers log output(console, file)
	Writers []string `yaml:"writers"`
	// FileConfig 日志文件配置，如果 Writers 为 file 则该配置不能为空
	FileConfig *FileConfig `yaml:"file_config"`

	// Format log format type (console, json)
	Format string `yaml:"format"`

	// Level log level debug info error...
	Level string `yaml:"level"`

	// CallerSkip 控制 log 函数嵌套深度
	CallerSkip int `yaml:"caller_skip"`
}

// FileConfig 日志文件的配置
type FileConfig struct {
	// Filepath log file path
	Filepath string `yaml:"filepath"`
	// MaxAge log file max age, days
	MaxAge int `yaml:"max_age"`
	// MaxBackups max backup files
	MaxBackups int `yaml:"max_backups"`
	// Compress log file is compress
	Compress bool `yaml:"compress"`
	// MaxSize max file size, MB
	MaxSize int `yaml:"max_size"`
}
