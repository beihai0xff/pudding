package configs

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/spf13/viper"

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

func Init(filePath string) {
	viper.SetConfigFile(filePath)

	// Find and read the config file
	if err := viper.ReadInConfig(); err != nil { // Handle errors reading the config file
		if errors.As(err, &viper.ConfigFileNotFoundError{}) {
			panic(fmt.Errorf("config file not found: %w", err))
		} else {
			// Config file was found but another error was produced
			panic(fmt.Errorf("failed to read config file: %w", err))
		}
	}
	if err := viper.Unmarshal(c); err != nil {
		panic(fmt.Errorf("failed to unmarshal config: %w", err))
	}

	c.Pulsar.ProducersConfig = append(c.Pulsar.ProducersConfig, conf.ProducerConfig{
		Topic:                   types.DefaultTopic,
		BatchingMaxPublishDelay: 20,
		BatchingMaxMessages:     100,
		BatchingMaxSize:         1024,
	})

	var str bytes.Buffer
	_ = json.Indent(&str, c.JSON(), "", "    ")
	log.Printf("config: %s \n", str.String())
}

func GetRedisConfig() *conf.RedisConfig {
	return c.Redis
}

func GetPulsarConfig() *conf.PulsarConfig {
	return c.Pulsar
}

func GetSchedulerConfig() *conf.SchedulerConfig {
	return c.Scheduler
}
