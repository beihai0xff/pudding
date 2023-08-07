package configs

import (
	"github.com/jinzhu/copier"
	"github.com/samber/lo"

	"github.com/beihai0xff/pudding/pkg/log"
)

var defaultLogFileConfig = log.FileConfig{
	MaxAge:     7, // days
	MaxBackups: 10,
	Compress:   false,
	MaxSize:    256, // megabytes
}

// GetLogConfig get specify log config by log name
func GetLogConfig(logName string) *log.Config {
	var logConfig []log.Config
	if err := UnmarshalToStruct("server_config.base_config.log_config", &logConfig); err != nil {
		panic(err)
	}

	c := log.DefaultConfig

	if len(logConfig) != 0 {
		v, ok := lo.Find(logConfig, func(conf log.Config) bool {
			return conf.LogName == logName
		})

		if ok {
			// if log writers contain file, then set file config
			if lo.Contains[string](v.Writers, log.OutputFile) {
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
