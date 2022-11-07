package configs

import conf "github.com/beihai0xff/pudding/configs"

var c = &Config{
	Broker:       "",
	MessageQueue: "",
	Scheduler: &conf.SchedulerConfig{
		TimeSliceInterval: "",
	},
	Redis: &conf.RedisConfig{
		RedisURL:    "",
		DialTimeout: 10,
	},
	Pulsar: &conf.PulsarConfig{
		PulsarURL:         "",
		ConnectionTimeout: 10,
		ProducersConfig:   nil,
		Log: &conf.LogConfig{
			Writers:    []string{"console"},
			Format:     conf.OutputConsole,
			Level:      "info",
			CallerSkip: 1,
		},
	},
}

type Config struct {
	Broker       string `json:"broker" yaml:"broker" mapstructure:"broker"`
	MessageQueue string `json:"messageQueue" yaml:"messageQueue" mapstructure:"messageQueue"`
	// Scheduler config
	Scheduler *conf.SchedulerConfig `json:"schedulerConfig" yaml:"schedulerConfig" mapstructure:"schedulerConfig"`

	Redis  *conf.RedisConfig  `json:"redisConfig" yaml:"redisConfig" mapstructure:"redisConfig"`
	Pulsar *conf.PulsarConfig `json:"pulsarConfig" yaml:"pulsarConfig" mapstructure:"pulsarConfig"`
}
