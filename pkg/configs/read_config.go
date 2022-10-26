package configs

import (
	"github.com/apache/pulsar-client-go/pulsar"

	"github.com/beihai0xff/pudding/pkg/yaml"
	"github.com/beihai0xff/pudding/types"
)

var c = &Config{
	redis:      &RedisConfig{},
	delayQueue: &DelayQueueConfig{},
	pulsar:     &PulsarConfig{},
}

type Config struct {
	redis      *RedisConfig
	delayQueue *DelayQueueConfig
	pulsar     *PulsarConfig
}

func Init(filePath string) {
	if err := yaml.Parse(filePath, c); err != nil {
		panic(err)
	}
	c.pulsar.ProducersConfig = append(c.pulsar.ProducersConfig, pulsar.ProducerOptions{
		Topic:           types.DefaultTopic,
		SendTimeout:     10,
		CompressionType: pulsar.ZSTD,
	})
}
