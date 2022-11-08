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
	},
}

// Config is the config for scheduler module.
type Config struct {
	// Broker is the broker type, e.g. redis
	Broker string `json:"broker" yaml:"broker" mapstructure:"broker"`
	// MessageQueue is the message queue type, e.g. pulsar, kakfa, etc.
	MessageQueue string `json:"messageQueue" yaml:"messageQueue" mapstructure:"messageQueue"`
	// Scheduler config
	Scheduler *conf.SchedulerConfig `json:"schedulerConfig" yaml:"schedulerConfig" mapstructure:"schedulerConfig"`

	Redis  *conf.RedisConfig  `json:"redisConfig" yaml:"redisConfig" mapstructure:"redisConfig"`
	Pulsar *conf.PulsarConfig `json:"pulsarConfig" yaml:"pulsarConfig" mapstructure:"pulsarConfig"`

	// Logger log config for output config message, do not use it
	Logger map[string]*conf.LogConfig `json:"log_config" yaml:"log_config" mapstructure:"log_config"`
}

// GetRedisConfig returns the redis config.
func GetRedisConfig() *conf.RedisConfig {
	return c.Redis
}

// GetPulsarConfig returns the pulsar config.
func GetPulsarConfig() *conf.PulsarConfig {
	return c.Pulsar
}

// GetSchedulerConfig returns the scheduler config.
func GetSchedulerConfig() *conf.SchedulerConfig {
	return c.Scheduler
}
