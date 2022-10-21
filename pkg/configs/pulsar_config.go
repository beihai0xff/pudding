package configs

import (
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
)

// PulsarConfig pulsar config
type PulsarConfig struct {
	PulsarURL       string                   `json:"pulsarURL" yaml:"pulsarURL"`
	DialTimeout     int                      `json:"dialTimeout" yaml:"dialTimeout"`
	ProducersConfig []pulsar.ProducerOptions `json:"producersConfig" yaml:"producersConfig"`
}

type ProducerConfig struct {
	// Topic specifies the topic this producer will be publishing on.
	// This argument is required when constructing the producer.
	Topic string `json:"topic" yaml:"topic"`
	// BatchingMaxPublishDelay specifies the time period within which the messages sent will be batched (default: 10ms)
	// if batch messages are enabled. If set to a non-zero value, messages will be queued until this time
	// interval or until
	BatchingMaxPublishDelay time.Duration

	// BatchingMaxMessages specifies the maximum number of messages permitted in a batch. (default: 1000)
	// If set to a value greater than 1, messages will be queued until this threshold is reached or
	// BatchingMaxSize (see below) has been reached or the batch interval has elapsed.
	BatchingMaxMessages uint

	// BatchingMaxSize specifies the maximum number of bytes permitted in a batch. (default 128 KB)
	// If set to a value greater than 1, messages will be queued until this threshold is reached or
	// BatchingMaxMessages (see above) has been reached or the batch interval has elapsed.
	BatchingMaxSize uint
}

func GetPulsarConfig() *PulsarConfig {
	return c.pulsar
}
