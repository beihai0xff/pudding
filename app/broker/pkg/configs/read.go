package configs

import (
	"bytes"
	"encoding/json"
	"log"

	"github.com/samber/lo"

	conf "github.com/beihai0xff/pudding/configs"
	"github.com/beihai0xff/pudding/types"
)

// ReadFrom read the config from the given configPath.
func ReadFrom(configPath string, opts ...OptionFunc) {
	conf.Parse(configPath, "yaml", c, conf.ReadFromFile)
	c.ServerConfig.SetFlags()

	for _, opt := range opts {
		opt(c)
	}
	setConnector()

	var str bytes.Buffer
	_ = json.Indent(&str, c.JSON(), "", "    ")
	log.Printf("pudding scheduler config:\n %s \n", str.String())
}

func setConnector() {
	if c.ServerConfig.MessageTopic == "" {
		c.ServerConfig.MessageTopic = types.DefaultTopic
	}
	if c.ServerConfig.TokenTopic == "" {
		c.ServerConfig.TokenTopic = types.TokenTopic
	}

	if !lo.ContainsBy[conf.ProducerConfig](c.Pulsar.ProducersConfig, func(v conf.ProducerConfig) bool {
		return v.Topic == c.ServerConfig.MessageTopic
	}) {
		c.Pulsar.ProducersConfig = append(c.Pulsar.ProducersConfig, conf.ProducerConfig{
			Topic:                   c.ServerConfig.MessageTopic,
			BatchingMaxPublishDelay: 20,
			BatchingMaxMessages:     100,
			BatchingMaxSize:         1024,
		})
	}

	if !lo.ContainsBy[conf.ProducerConfig](c.Pulsar.ProducersConfig, func(v conf.ProducerConfig) bool {
		return v.Topic == c.ServerConfig.TokenTopic
	}) {
		c.Pulsar.ProducersConfig = append(c.Pulsar.ProducersConfig, conf.ProducerConfig{
			Topic:                   c.ServerConfig.TokenTopic,
			BatchingMaxPublishDelay: 20,
			BatchingMaxMessages:     100,
			BatchingMaxSize:         1024,
		})
	}
}

// OptionFunc is the option function for config.
type OptionFunc func(config *Config)

// WithRedisURL set the redis url.
func WithRedisURL(url string) OptionFunc {
	return func(config *Config) {
		if url != "" {
			config.Redis.URL = url
		}
	}
}

// WithPulsarURL set the pulsar url.
func WithPulsarURL(url string) OptionFunc {
	return func(config *Config) {
		if url != "" {
			config.Pulsar.URL = url
		}
	}
}
