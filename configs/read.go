package configs

import (
	"bytes"
	"encoding/json"
	"log"

	"github.com/beihai0xff/pudding/pkg/yaml"
	"github.com/beihai0xff/pudding/types"
)

var c = &Config{
	Broker:       "",
	MessageQueue: "",
	Scheduler:    nil,
	Redis:        nil,
	Pulsar: &PulsarConfig{
		PulsarURL:         "",
		ConnectionTimeout: 0,
		ProducersConfig:   nil,
		Log: &LogConfig{
			Writers:    []string{"console"},
			Format:     OutputConsole,
			Level:      "info",
			CallerSkip: 1,
		},
	},
	MySQL: &MySQLConfig{
		DSN: "",
		Log: &LogConfig{
			Writers:    []string{"console"},
			Format:     OutputConsole,
			Level:      "info",
			CallerSkip: 3,
		},
	},
}

type Config struct {
	Broker       string           `json:"broker" yaml:"broker"`
	MessageQueue string           `json:"messageQueue" yaml:"messageQueue"`
	Scheduler    *SchedulerConfig `json:"schedulerConfig" yaml:"schedulerConfig"`

	Redis  *RedisConfig  `json:"redisConfig" yaml:"redisConfig"`
	Pulsar *PulsarConfig `json:"pulsarConfig" yaml:"pulsarConfig"`
	MySQL  *MySQLConfig  `json:"mysqlConfig" yaml:"mysqlConfig"`
}

func (c *Config) JSON() []byte {
	b, err := json.Marshal(c)
	if err != nil {
		log.Panicf("marshal config failed: %v", err)
		return nil
	}

	return b
}

func Init(filePath string) {
	if err := yaml.Parse(filePath, c); err != nil {
		panic(err)
	}
	c.Pulsar.ProducersConfig = append(c.Pulsar.ProducersConfig, ProducerConfig{
		Topic:                   types.DefaultTopic,
		BatchingMaxPublishDelay: 20,
		BatchingMaxMessages:     100,
		BatchingMaxSize:         1024,
	})

	var str bytes.Buffer
	_ = json.Indent(&str, c.JSON(), "", "    ")
	log.Printf("config: %s \n", str.String())
}
