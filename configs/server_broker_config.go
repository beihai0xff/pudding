// Package configs provide config management
// server_broker_config.go contains the config of broker server
package configs

import (
	"fmt"

	"github.com/beihai0xff/pudding/pkg/log"
)

// BrokerConfig BrokerConfig Config
type BrokerConfig struct {
	// BaseConfig server base Config
	BaseConfig `json:"server_config" yaml:"server_config" mapstructure:"server_config"`
	// ServerConfig server config
	// use same struct tag merge BaseConfig to ServerConfig
	//nolint:govet,revive
	ServerConfig struct {
		// TimeSliceInterval broker loop time interval
		TimeSliceInterval string `json:"time_slice_interval" yaml:"time_slice_interval" mapstructure:"time_slice_interval"`
		// MessageTopic default message topic, if no topic set in message, use this topic
		MessageTopic string `json:"message_topic" yaml:"message_topic" mapstructure:"message_topic"`
		// TokenTopic TimeSlice token topic
		TokenTopic string `json:"token_topic" yaml:"token_topic" mapstructure:"token_topic"`
		// Broker type
		Broker string `json:"broker" yaml:"broker" mapstructure:"broker"`
		// Connector type
		Connector string `json:"connector" yaml:"connector" mapstructure:"connector"`
		// EtcdURLs etcd connection urls
		EtcdURLs []string `json:"etcd_urls" yaml:"etcd_urls" mapstructure:"etcd_urls"`
	} `json:"server_config" yaml:"server_config" mapstructure:"server_config"`

	// Logger log config for output config message
	Logger []log.Config `json:"log_config" yaml:"log_config" mapstructure:"log_config"`

	// RedisConfig redis config
	RedisConfig RedisConfig `json:"redis_config" yaml:"redis_config" mapstructure:"redis_config"`
	// Kafka kafka config
	KafkaConfig KafkaConfig `json:"kafka_config" yaml:"kafka_config" mapstructure:"kafka_config"`
}

// ParseBrokerConfig read the config from the given configPath.
func ParseBrokerConfig(configPath string, opts ...OptionFunc) *BrokerConfig {
	if err := Parse(configPath, ConfigFormatYAML, ReadFromFile, opts...); err != nil {
		panic(err)
	}

	// unmarshal all config to BrokerConfig{}
	c := BrokerConfig{}
	if err := UnmarshalToStruct("", &c); err != nil {
		panic(err)
	}

	fmt.Printf("c: %+v\n", c)

	return &c
}
