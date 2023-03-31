package configs

// KafkaConfig kafka client config
type KafkaConfig struct {
	// Address is the address list of kafka server
	Address []string `json:"address" yaml:"address" mapstructure:"address"`
	// Network is the network of kafka server
	Network string `json:"network" yaml:"network" mapstructure:"network"`
	// WriteTimeout is to write timeout of kafka server
	WriteTimeout int `json:"write_timeout" yaml:"write_timeout" mapstructure:"write_timeout"`
	// ReadTimeout is the read timeout of kafka server
	ReadTimeout int `json:"read_timeout" yaml:"read_timeout" mapstructure:"read_timeout"`
	// NumPartitions is the number of partitions of kafka server
	NumPartitions int `json:"num_partitions" yaml:"num_partitions" mapstructure:"num_partitions"`
	// ReplicationFactor is the number of replicas of kafka server
	ReplicationFactor int `json:"replication_factor" yaml:"replication_factor" mapstructure:"replication_factor"`
	// ConsumerMaxWaitTime is the max wait time of consumer, milliseconds
	ConsumerMaxWaitTime int `json:"consumer_max_wait_time" yaml:"consumer_max_wait_time" mapstructure:"consumer_max_wait_time"`
	// ProducerBatchTimeout is the min wait time of producer, milliseconds
	ProducerBatchTimeout int `json:"producer_batch_timeout" yaml:"producer_batch_timeout" mapstructure:"producer_batch_timeout"`
	// BatchSize is the max number of messages to send in a single batch
	BatchSize int `json:"batch_size" yaml:"batch_size" mapstructure:"batch_size"`
}
