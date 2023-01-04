package configs

import conf "github.com/beihai0xff/pudding/configs"

var c = &Config{
	BrokerConfig: &conf.BrokerConfig{
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
}

// Config is the config for scheduler module.
type Config struct {
	// BrokerConfig config
	BrokerConfig *conf.BrokerConfig `json:"broker_config" yaml:"broker_config" mapstructure:"broker_config"`

	Redis  *conf.RedisConfig  `json:"redis_config" yaml:"redis_config" mapstructure:"redis_config"`
	Pulsar *conf.PulsarConfig `json:"pulsar_config" yaml:"pulsar_config" mapstructure:"pulsar_config"`
}

// GetRedisConfig returns the redis config.
func GetRedisConfig() *conf.RedisConfig {
	return c.Redis
}

// GetPulsarConfig returns the pulsar config.
func GetPulsarConfig() *conf.PulsarConfig {
	return c.Pulsar
}

// GetBrokerConfig returns the scheduler config.
func GetBrokerConfig() *conf.BrokerConfig {
	return c.BrokerConfig
}

// GetNameServerURL returns the name server url.
func GetNameServerURL() string {
	return c.BrokerConfig.NameServerURL
}
