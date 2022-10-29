package configs

import (
	"bytes"
	"encoding/json"

	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/pkg/yaml"
	"github.com/beihai0xff/pudding/types"
)

var c = &Config{}

type Config struct {
	Broker       string           `json:"broker" yaml:"broker"`
	MessageQueue string           `json:"messageQueue" yaml:"messageQueue"`
	Redis        *RedisConfig     `json:"redisConfig" yaml:"redisConfig"`
	Scheduler    *SchedulerConfig `json:"schedulerConfig" yaml:"schedulerConfig"`
	Pulsar       *PulsarConfig    `json:"pulsarConfig" yaml:"pulsarConfig"`
}

func (c *Config) JSON() []byte {
	b, err := json.Marshal(c)
	if err != nil {
		log.Errorf("marshal config failed: %v", err)
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
	log.Infof("config: %s \n", str.String())
}
