package configs

// KafkaConfig kafka client config
type KafkaConfig struct {
	// Host is the host of kafka server
	Host string `json:"host"`
	// Port is the port of kafka server
	Port int `json:"port"`
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
	// ConsumerMaxWaitTime is the max wait time of consumer
	ConsumerMaxWaitTime int `json:"consumerMaxWaitTime"`
	// ProducerBatchTimeout is the min wait time of producer
	ProducerBatchTimeout int `json:"producerBatchTimeout"`
	// BatchSize is the max number of messages to send in a single batch
	BatchSize int `json:"batchSize"`
}
