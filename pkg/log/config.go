// Package log provides the log
// config.go contains the log config
package log

import (
	"go.uber.org/zap/zapcore"
)

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

// Config log config
type Config struct {
	LogName string `json:"log_name" yaml:"log_name" mapstructure:"log_name"`
	// Writers log output(console, file)
	Writers []string `yaml:"writers" mapstructure:"writers" json:"writers"`
	// FileConfig log file config, if writers has file, must set file config
	FileConfig FileConfig `yaml:"file_config" mapstructure:"file_config" json:"file_config"`

	// Format log format type (console, json)
	Format string `yaml:"format" mapstructure:"format" json:"format"`

	// Level log level debug info error...
	Level string `yaml:"level" mapstructure:"level" json:"level"`

	// CallerSkip log caller skip
	CallerSkip int `yaml:"caller_skip" mapstructure:"caller_skip" json:"caller_skip"`
}

// FileConfig log file config
type FileConfig struct {
	// Filepath log file path
	Filepath string `yaml:"filepath" mapstructure:"filepath" json:"filepath"`
	// MaxAge log file max age, days
	MaxAge int `yaml:"max_age" mapstructure:"max_age" json:"max_age"`
	// MaxBackups max backup files
	MaxBackups int `yaml:"max_backups" mapstructure:"max_backups" json:"max_backups"`
	// Compress log file is compress
	Compress bool `yaml:"compress" mapstructure:"compress" json:"compress"`
	// MaxSize max file size, MB
	MaxSize int `yaml:"max_size" mapstructure:"max_size" json:"max_size"`
}

var (
	// DefaultConfig default log config
	DefaultConfig = Config{
		Writers:    []string{OutputConsole},
		Format:     EncoderTypeConsole,
		Level:      "info",
		CallerSkip: 1,
	}
)
