package broker

import (
	"context"

	"github.com/beihai0xff/pudding/types"
)

// nolint:lll
//go:generate mockgen -destination=../../../test/mock/broker_mock.go -package=mock github.com/beihai0xff/pudding/app/scheduler/broker DelayBroker

// DelayBroker is a queue to store messages with delay time
// the message will be delivered to the realtime queue after the delay time
type DelayBroker interface {
	// Produce produce a Message to DelayBroker
	Produce(ctx context.Context, msg *types.Message) error
	// Consume consume Messages from the queue
	Consume(ctx context.Context, now, batchSize int64, fn types.HandleMessage) error
	// Close the queue
	Close() error
}
