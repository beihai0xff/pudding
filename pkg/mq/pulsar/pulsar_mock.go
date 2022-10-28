package pulsar

import (
	"github.com/beihai0xff/pudding/pkg/configs"
	"github.com/beihai0xff/pudding/types"
)

// NewMockPulsar create a mock pulsar client
func NewMockPulsar() *Client {
	config := &configs.PulsarConfig{
		PulsarURL:         "pulsar://localhost:6650",
		ConnectionTimeout: 10,
		ProducersConfig: []configs.ProducerConfig{
			{Topic: "test_topic_1"},
			{Topic: "test_topic_2"},
			{Topic: types.DefaultTopic},
		},
	}
	return New(config)
}
