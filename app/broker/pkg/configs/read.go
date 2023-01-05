package configs

import (
	"bytes"
	"encoding/json"
	"log"

	conf "github.com/beihai0xff/pudding/configs"
	"github.com/beihai0xff/pudding/types"
)

func (c *Config) JSON() []byte {
	b, err := json.Marshal(c)
	if err != nil {
		log.Panicf("marshal config failed: %v", err)
		return nil
	}

	return b
}

func Init(filePath string, opts ...OptionFunc) {
	conf.Parse(filePath, "yaml", c, conf.ReadFromFile)

	if c.ServerConfig.MessageTopic == "" {
		c.ServerConfig.MessageTopic = types.DefaultTopic
	}
	if c.ServerConfig.TokenTopic == "" {
		c.ServerConfig.TokenTopic = types.TokenTopic
	}

	for _, opt := range opts {
		opt(c)
	}

	producers := make(map[string]struct{})
	for _, v := range c.Pulsar.ProducersConfig {
		producers[v.Topic] = struct{}{}
	}

	if _, ok := producers[c.ServerConfig.MessageTopic]; !ok {
		c.Pulsar.ProducersConfig = append(c.Pulsar.ProducersConfig, conf.ProducerConfig{
			Topic:                   types.DefaultTopic,
			BatchingMaxPublishDelay: 20,
			BatchingMaxMessages:     100,
			BatchingMaxSize:         1024,
		})
	}

	if _, ok := producers[c.ServerConfig.TokenTopic]; !ok {
		c.Pulsar.ProducersConfig = append(c.Pulsar.ProducersConfig, conf.ProducerConfig{
			Topic:                   types.TokenTopic,
			BatchingMaxPublishDelay: 20,
			BatchingMaxMessages:     100,
			BatchingMaxSize:         1024,
		})
	}

	var str bytes.Buffer
	_ = json.Indent(&str, c.JSON(), "", "    ")
	log.Printf("pudding scheduler config:\n %s \n", str.String())
}

type OptionFunc func(config *Config)

func WithRedisURL(url string) OptionFunc {
	return func(config *Config) {
		if url != "" {
			config.Redis.URL = url
		}
	}
}

func WithPulsarURL(url string) OptionFunc {
	return func(config *Config) {
		if url != "" {
			config.Pulsar.URL = url
		}
	}
}
