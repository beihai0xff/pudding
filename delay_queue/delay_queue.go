package delay_queue

import (
	"context"

	"github.com/beihai0xff/pudding/types"
)

type delayQueue interface {
	// Produce produce a Message to the queue
	Produce(ctx context.Context, msg *types.Message) error
	// NewConsumer New a consumer to consume Messages from the queue
	NewConsumer(topic string, batchSize int, fn func(msg *types.Message) error)
}
