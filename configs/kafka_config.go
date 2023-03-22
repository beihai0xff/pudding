package configs

// KafkaConfig kafka client config
type KafkaConfig struct {
	// Address is the address list of kafka server
	Address []string `json:"address"`
	// Network is the network of kafka server
	Network string `json:"network"`
	// WriteTimeout is to write timeout of kafka server
	WriteTimeout int `json:"writeTimeout"`
	// ReadTimeout is the read timeout of kafka server
	ReadTimeout int `json:"readTimeout"`
	// NumPartitions is the number of partitions of kafka server
	NumPartitions int `json:"numPartitions"`
	// ReplicationFactor is the number of replicas of kafka server
	ReplicationFactor int `json:"replicationFactor"`
	// ConsumerMaxWaitTime is the max wait time of consumer, milliseconds
	ConsumerMaxWaitTime int `json:"consumerMaxWaitTime"`
	// ProducerBatchTimeout is the min wait time of producer, milliseconds
	ProducerBatchTimeout int `json:"producerBatchTimeout"`
	// BatchSize is the max number of messages to send in a single batch
	BatchSize int `json:"batchSize"`
}
