package pulsar

import (
	"github.com/apache/pulsar-client-go/pulsar"

	"github.com/beihai0xff/pudding/pkg/configs"
)

// NewMockPulsar create a mock pulsar client
func NewMockPulsar() *Client {
	config := &configs.PulsarConfig{
		PulsarURL:   "pulsar://localhost:6650",
		DialTimeout: 10,
		ProducersConfig: []pulsar.ProducerOptions{
			{Topic: "test_topic_1"},
			{Topic: "test_topic_2"},
		},
	}
	return New(config)
}
