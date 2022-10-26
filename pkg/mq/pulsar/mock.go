package pulsar

import (
	"github.com/apache/pulsar-client-go/pulsar"

	"github.com/beihai0xff/pudding/pkg/configs"
	"github.com/beihai0xff/pudding/types"
)

// NewMockPulsar create a mock pulsar client
func NewMockPulsar() *Client {
	config := &configs.PulsarConfig{
		PulsarURL:   "pulsar://localhost:6650",
		DialTimeout: 10,
		ProducersConfig: []pulsar.ProducerOptions{
			{Topic: "test_topic_1"},
			{Topic: "test_topic_2"},
			{Topic: types.DefaultTopic, SendTimeout: 10, CompressionType: pulsar.ZSTD},
		},
	}
	return New(config)
}
