package pulsar

import (
	configs2 "github.com/beihai0xff/pudding/configs"
)

// newMockPulsar create a mock pulsar client
func newMockPulsar() *Client {
	configs2.Init("../../../test/config.test.yaml")

	return New(configs2.GetPulsarConfig())
}
