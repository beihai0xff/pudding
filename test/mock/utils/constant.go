// Package test_utils provides some constants for testing.
package test_utils

import "github.com/beihai0xff/pudding/configs"

const (
	// TestConfigFilePath is the path of test config file.
	TestConfigFilePath = "./test/data/config.test.yaml"
)

// TestMySQLConfig is the test config for MySQL.
var TestMySQLConfig = &configs.MySQLConfig{
	DSN: "root:my-secret-pw@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=UTC",
}

// TestKafkaConfig is the test config for Kafka.
var TestKafkaConfig = &configs.KafkaConfig{
	Address:              []string{"localhost:9092"},
	Network:              "tcp",
	ConsumerMaxWaitTime:  1000,
	ProducerBatchTimeout: 1000,
	BatchSize:            1,
}
