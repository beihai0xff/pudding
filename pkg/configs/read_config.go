package configs

import "github.com/beihai0xff/pudding/pkg/yaml"

var c *Config

type Config struct {
	redis *RedisConfig
}

func Init(filePath string) {
	yaml.Parse(filePath, c)
}
