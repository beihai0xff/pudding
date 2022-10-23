package pulsar_broker

import (
	"context"

	"github.com/beihai0xff/pudding/pkg/mq/pulsar"
	"github.com/beihai0xff/pudding/types"
)

type RealTimeQueue struct {
	pulsar *pulsar.Client
}

func New(client *pulsar.Client) *RealTimeQueue {
	return &RealTimeQueue{
		pulsar: client,
	}
}

// Produce produce a Message to the queue in real time
func (q *RealTimeQueue) Produce(ctx context.Context, msg *types.Message) error {
	return q.pulsar.Produce(ctx, msg.Topic, convertToPulsarProducerMessage(msg))
}

// NewConsumer consume Messages from the queue in real time
func (q *RealTimeQueue) NewConsumer(topic, group string, batchSize int, fn types.HandleMessage) {

}

// Close the queue
func (q *RealTimeQueue) Close() error {
	return nil
}
