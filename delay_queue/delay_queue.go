package delay_queue

import (
	"context"

	"github.com/beihai0xff/pudding/types"
)

type DelayQueue interface {
	// Produce produce a Message to DelayQueue
	Produce(ctx context.Context, msg *types.Message) error
	// NewConsumer New a consumer to consume Messages from the queue
	NewConsumer(topic string, partition, batchSize int, fn func(msg *types.Message) error)
	// Close the queue
	Close() error
}

type RealTimeQueue interface {
	// Produce produce a Message to the queue in real time
	Produce(ctx context.Context, msg *types.Message) error
	// NewConsumer consume Messages from the queue in real time
	NewConsumer(topic, group, consumerName string, batchSize int, fn func(msg *types.Message) error)
	// Close the queue
	Close() error
}

type Queue interface {
	*DelayQueue
	*RealTimeQueue
}

// NewDelayQueue create a new DelayQueue
func NewDelayQueue() DelayQueue {
	return nil
}

// NewRealTimeQueue create a new RealTimeQueue
func NewRealTimeQueue() RealTimeQueue {
	return nil
}
