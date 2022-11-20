package pulsar_connector

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

func convertPulsarMessageToDelayMessage(msg pulsar.Message) *types.Message {
	return &types.Message{
		Topic:   msg.Topic(),
		Key:     msg.Key(),
		Payload: msg.Payload(),
	}
}
