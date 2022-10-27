package scheduler

import (
	"context"

	"github.com/beihai0xff/pudding/pkg/mq/pulsar"
	rdb "github.com/beihai0xff/pudding/pkg/redis"
	"github.com/beihai0xff/pudding/scheduler/broker/pulsar_broker"
	"github.com/beihai0xff/pudding/scheduler/broker/redis_broker"
	"github.com/beihai0xff/pudding/types"
)

//go:generate mockgen -destination=../test/mock/queue_mock.go -package=scheduler_test github.com/beihai0xff/pudding/scheduler DelayQueue,RealTimeQueue

type DelayQueue interface {
	// Produce produce a Message to DelayQueue
	Produce(ctx context.Context, timeSlice string, msg *types.Message) error
	// Consume consume Messages from the queue
	Consume(ctx context.Context, timeSlice string, now, batchSize int64, fn types.HandleMessage) error
	// Close the queue
	Close() error
}

type RealTimeQueue interface {
	// Produce produce a Message to the queue in real time
	Produce(ctx context.Context, msg *types.Message) error
	// NewConsumer new a consumer to consume Messages from the realtime queue in background
	NewConsumer(topic, group string, batchSize int, fn types.HandleMessage) error
	// Close the queue
	Close() error
}

// NewDelayQueue create a new DelayQueue
func NewDelayQueue(c *rdb.Client) DelayQueue {
	return redis_broker.NewDelayQueue(c)
}

// NewRealTimeQueue create a new RealTimeQueue
func NewRealTimeQueue(c *pulsar.Client) RealTimeQueue {
	return pulsar_broker.NewRealTimeQueue(c)
}
