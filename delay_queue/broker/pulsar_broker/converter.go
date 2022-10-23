package pulsar_broker

import (
	"github.com/apache/pulsar-client-go/pulsar"

	"github.com/beihai0xff/pudding/types"
)

func convertToPulsarProducerMessage(msg *types.Message) *pulsar.ProducerMessage {
	return &pulsar.ProducerMessage{
		Payload: msg.Payload,
		Key:     msg.Key,
	}
}
