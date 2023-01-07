package configs

import conf "github.com/beihai0xff/pudding/configs"

var c = &Config{
	ServerConfig: &conf.BrokerConfig{
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
	// ServerConfig server config
	ServerConfig *conf.BrokerConfig `json:"server_config" yaml:"server_config" mapstructure:"server_config"`

	// Redis redis config
	Redis *conf.RedisConfig `json:"redis_config" yaml:"redis_config" mapstructure:"redis_config"`
	// Pulsar pulsar config
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

// GetServerConfig returns the scheduler config.
func GetServerConfig() *conf.BrokerConfig {
	return c.ServerConfig
}

// GetNameServerURL returns the name server url.
func GetNameServerURL() string {
	return c.ServerConfig.NameServerURL
}
