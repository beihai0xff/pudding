package pulsar_connector

import (
	"context"
	"fmt"

	p "github.com/apache/pulsar-client-go/pulsar"

	"github.com/beihai0xff/pudding/pkg/mq/pulsar"
	"github.com/beihai0xff/pudding/types"
)

type RealTimeQueue struct {
	pulsar *pulsar.Client
}

func NewRealTimeQueue(client *pulsar.Client) *RealTimeQueue {
	return &RealTimeQueue{
		pulsar: client,
	}
}

// Produce produce a Message to the queue in real time
func (q *RealTimeQueue) Produce(ctx context.Context, msg *types.Message) error {
	if msg.Payload == nil || len(msg.Payload) == 0 {
		return fmt.Errorf("message payload can not be empty")
	}
	return q.pulsar.Produce(ctx, msg.Topic, convertToPulsarProducerMessage(msg))
}

// NewConsumer consume Messages from the queue in real time
func (q *RealTimeQueue) NewConsumer(topic, group string, batchSize int, fn types.HandleMessage) error {
	var f pulsar.HandleMessage = func(ctx context.Context, msg p.Message) error {
		return fn(ctx, convertPulsarMessageToDelayMessage(msg))
	}

	return q.pulsar.NewConsumer(topic, group, f)
}

// Close the queue
func (q *RealTimeQueue) Close() error {
	return nil
}
