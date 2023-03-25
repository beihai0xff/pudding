package configs

// KafkaConfig kafka client config
type KafkaConfig struct {
	// Address is the address list of kafka server
	Address []string `json:"address" yaml:"address" mapstructure:"address"`
	// Network is the network of kafka server
	Network string `json:"network" yaml:"network" mapstructure:"network"`
	// WriteTimeout is to write timeout of kafka server
	WriteTimeout int `json:"writeTimeout" yaml:"write_timeout" mapstructure:"write_timeout"`
	// ReadTimeout is the read timeout of kafka server
	ReadTimeout int `json:"readTimeout" yaml:"read_timeout" mapstructure:"read_timeout"`
	// NumPartitions is the number of partitions of kafka server
	NumPartitions int `json:"numPartitions" yaml:"num_partitions" mapstructure:"num_partitions"`
	// ReplicationFactor is the number of replicas of kafka server
	ReplicationFactor int `json:"replicationFactor" yaml:"replication_factor" mapstructure:"replication_factor"`
	// ConsumerMaxWaitTime is the max wait time of consumer, milliseconds
	ConsumerMaxWaitTime int `json:"consumerMaxWaitTime" yaml:"consumer_max_wait_time" mapstructure:"consumer_max_wait_time"`
	// ProducerBatchTimeout is the min wait time of producer, milliseconds
	ProducerBatchTimeout int `json:"producerBatchTimeout" yaml:"producer_batch_timeout" mapstructure:"producer_batch_timeout"`
	// BatchSize is the max number of messages to send in a single batch
	BatchSize int `json:"batchSize" yaml:"batch_size" mapstructure:"batch_size"`
}
