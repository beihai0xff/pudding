// Package configs provides config management
// server_broker_config.go contains the config of broker server
package configs

// BrokerConfig BrokerConfig Config
type BrokerConfig struct {
	ServerConfig struct {
		// BaseConfig server base Config
		BaseConfig `json:"base_config" yaml:"base_config" mapstructure:"base_config"`

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

	// RedisConfig redis config
	RedisConfig *RedisConfig `json:"redis_config" yaml:"redis_config" mapstructure:"redis_config"`
	// Pulsar pulsar config
	PulsarConfig *PulsarConfig `json:"pulsar_config" yaml:"pulsar_config" mapstructure:"pulsar_config"`
	// Kafka kafka config
	KafkaConfig *KafkaConfig `json:"kafka_config" yaml:"kafka_config" mapstructure:"kafka_config"`
}

// ParseBrokerConfig read the config from the given configPath.
func ParseBrokerConfig(configPath string, opts ...OptionFunc) *BrokerConfig {
	if err := Parse(configPath, ConfigFormatYAML, ReadFromFile, opts...); err != nil {
		panic(err)
	}

	// unmarshal all config to BrokerConfig{}
	var c = BrokerConfig{}
	if err := UnmarshalToStruct("", &c); err != nil {
		panic(err)
	}

	// set flags
	c.ServerConfig.BaseConfig.SetFlags()

	return &c
}
