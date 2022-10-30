package pulsar

import (
	"github.com/beihai0xff/pudding/configs"
)

// newMockPulsar create a mock pulsar client
func newMockPulsar() *Client {
	configs.Init("../../../test/config.test.yaml")

	return New(configs.GetPulsarConfig())
}
