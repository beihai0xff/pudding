package configs

import (
	"github.com/beihai0xff/pudding/pkg/yaml"
)

var c = &Config{}

type Config struct {
	redis      *RedisConfig
	delayQueue *DelayQueueConfig
	pulsar     *PulsarConfig
}

func Init(filePath string) {
	if err := yaml.Parse(filePath, c); err != nil {
		panic(err)
	}
}
