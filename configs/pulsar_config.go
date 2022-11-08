//nolint:lll
package configs

// PulsarConfig pulsar config
type PulsarConfig struct {
	URL               string           `json:"url" yaml:"url" mapstructure:"url"`
	ConnectionTimeout int              `json:"connection_timeout" yaml:"connection_timeout" mapstructure:"connection_timeout"`
	ProducersConfig   []ProducerConfig `json:"producers_config" yaml:"producers_config" mapstructure:"producers_config"`
}

type ProducerConfig struct {
	// Topic specifies the topic this producer will be publishing on.
	// This argument is required when constructing the producer.
	Topic string `json:"topic" yaml:"topic" mapstructure:"topic"`
	// BatchingMaxPublishDelay specifies the time period within which the messages sent will be batched (default: 10ms)
	// if batch messages are enabled. If set to a non-zero value, messages will be queued until this time
	// interval or until
	BatchingMaxPublishDelay uint `json:"batching_max_publish_delay" yaml:"batching_max_publish_delay" mapstructure:"batching_max_publish_delay"`

	// BatchingMaxMessages specifies the maximum number of messages permitted in a batch. (default: 1000)
	// If set to a value greater than 1, messages will be queued until this threshold is reached or
	// BatchingMaxSize (see below) has been reached or the batch interval has elapsed.
	BatchingMaxMessages uint `json:"batching_max_messages" yaml:"batching_max_messages" mapstructure:"batching_max_messages"`

	// BatchingMaxSize specifies the maximum number of bytes permitted in a batch. (default 128 KB)
	// If set to a value greater than 1, messages will be queued until this threshold is reached or
	// BatchingMaxMessages (see above) has been reached or the batch interval has elapsed.
	BatchingMaxSize uint `json:"batching_max_size" yaml:"batching_max_size" mapstructure:"batching_max_size"`
}
