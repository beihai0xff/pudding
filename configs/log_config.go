// Package configs provides config management
package configs

import (
	"github.com/jinzhu/copier"
	"github.com/samber/lo"
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

// LogConfig log output: console file remote
type LogConfig struct {
	LogName string `json:"log_name" yaml:"log_name" mapstructure:"log_name"`
	// Writers log output(console, file)
	Writers []string `yaml:"writers" mapstructure:"writers" json:"writers"`
	// FileConfig log file config, if writers has file, must set file config
	FileConfig LogFileConfig `yaml:"file_config" mapstructure:"file_config" json:"file_config"`

	// Format log format type (console, json)
	Format string `yaml:"format" mapstructure:"format" json:"format"`

	// Level log level debug info error...
	Level string `yaml:"level" mapstructure:"level" json:"level"`

	// CallerSkip log caller skip
	CallerSkip int `yaml:"caller_skip" mapstructure:"caller_skip" json:"caller_skip"`
}

// LogFileConfig log file config
type LogFileConfig struct {
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
	defaultConfig = LogConfig{
		Writers:    []string{OutputConsole},
		Format:     EncoderTypeConsole,
		Level:      "info",
		CallerSkip: 1,
	}
	defaultLogFileConfig = LogFileConfig{
		MaxAge:     7, // days
		MaxBackups: 10,
		Compress:   false,
		MaxSize:    256, // megabytes
	}
)

// GetLogConfig get specify log config by log name
func GetLogConfig(logName string) *LogConfig {
	// get_logger_configs
	var logConfig []LogConfig
	if err := UnmarshalToStruct("server_config.base_config.log_config", &logConfig); err != nil {
		panic(err)
	}

	c := defaultConfig
	if len(logConfig) != 0 {
		v, ok := lo.Find(logConfig, func(conf LogConfig) bool {
			return conf.LogName == logName
		})

		if ok {
			// if log writers contains file, then set file config
			if lo.Contains[string](v.Writers, OutputFile) {
				fileConfig := defaultLogFileConfig
				c.FileConfig = fileConfig
			}
			if err := copier.CopyWithOption(&c, v, copier.Option{IgnoreEmpty: true, DeepCopy: true}); err != nil {
				panic(err)
			}
		}
	}

	c.LogName = logName

	return &c
}
