// Package connector provides the connector of the broker.
package connector

import (
	"context"

	"github.com/beihai0xff/pudding/api/gen/pudding/types/v1"
	type2 "github.com/beihai0xff/pudding/app/broker/pkg/types"
)

// RealTimeConnector is a connector which can send messages to the realtime queue
// the realtime queue can store or consume messages in realtime
type RealTimeConnector interface {
	// Produce produce a Message to the queue in real time
	Produce(ctx context.Context, msg *types.Message) error
	// NewConsumer new a consumer to consume Messages from the realtime queue in background
	NewConsumer(topic, group string, batchSize int, fn type2.HandleMessage) error
	// Close the queue
	Close() error
}
