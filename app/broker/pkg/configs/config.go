package configs

import conf "github.com/beihai0xff/pudding/configs"

var c = &Config{
	Scheduler: &conf.SchedulerConfig{
		TimeSliceInterval: "",
	},
	Redis: &conf.RedisConfig{
		URL:         "",
		DialTimeout: 10,
	},
	Pulsar: &conf.PulsarConfig{
		URL:               "",
		ConnectionTimeout: 10,
		ProducersConfig:   nil,
	},
	ConsulURL: "",
}

// Config is the config for scheduler module.
type Config struct {
	// Scheduler config
	Scheduler *conf.SchedulerConfig `json:"scheduler_config" yaml:"scheduler_config" mapstructure:"scheduler_config"`

	Redis  *conf.RedisConfig  `json:"redis_config" yaml:"redis_config" mapstructure:"redis_config"`
	Pulsar *conf.PulsarConfig `json:"pulsar_config" yaml:"pulsar_config" mapstructure:"pulsar_config"`

	// Logger log config for output config message, do not use it
	Logger    map[string]*conf.LogConfig `json:"log_config" yaml:"log_config" mapstructure:"log_config"`
	ConsulURL string                     `json:"consul_url" yaml:"consul_url" mapstructure:"consul_url"`
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

// GetConsulURL returns the consul url.
func GetConsulURL() string {
	return c.ConsulURL
}
