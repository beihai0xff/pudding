// Package storage provides the storage of the broker.
package storage

import (
	"context"

	"github.com/beihai0xff/pudding/api/gen/pudding/types/v1"
	type2 "github.com/beihai0xff/pudding/app/broker/pkg/types"
)

// DelayStorage is a queue to store messages with delay time
// the message will be delivered to the realtime queue after the delay time
type DelayStorage interface {
	// Produce produce a Message to DelayStorage
	Produce(ctx context.Context, msg *types.Message) error
	// Consume consume Messages from the queue
	Consume(ctx context.Context, now, batchSize int64, fn type2.HandleMessage) error
	// Close the queue
	Close() error
}
