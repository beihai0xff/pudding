package delay_queue

import (
	"context"

	"github.com/beihai0xff/pudding/types"
)

type DelayQueue interface {
	// Produce produce a Message to the queue
	Produce(ctx context.Context, msg *types.Message) error
	// NewConsumer New a consumer to consume Messages from the queue
	NewConsumer(topic string, partition, batchSize int, fn func(msg *types.Message) error)

	// NewRealTimeConsumer consume Messages from the queue in real time
	NewRealTimeConsumer(topic, group, consumerName string, batchSize int, fn func(msg *types.Message) error)
	// RealTimeProduce produce a Message to the queue in real time
	RealTimeProduce(ctx context.Context, topic string, msg *types.Message) error
	// Close the queue
	Close() error
}
