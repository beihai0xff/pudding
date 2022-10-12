package delay_queue

import (
	"context"
	"time"

	"github.com/beihai0xff/pudding/delay_queue/broker/redis_broker"
	"github.com/beihai0xff/pudding/types"
	"github.com/google/uuid"
)

type DelayQueue interface {
	// Produce produce a Message to DelayQueue
	Produce(ctx context.Context, bucketID int64, msg *types.Message) error
	// Consume New a consumer to consume Messages from the queue
	Consume(ctx context.Context, bucketID, batchSize int64,
		fn func(msg *types.Message) error) error
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

type Queue struct {
	delay    DelayQueue
	realtime RealTimeQueue
}

func NewQueue() *Queue {
	return &Queue{
		delay:    NewDelayQueue(),
		realtime: NewRealTimeQueue(),
	}
}

func (q *Queue) Produce(ctx context.Context, msg *types.Message) error {
	// 如果设置了 ReadyTime，则使用 ReadyTime
	// 否则使用当前时间
	if !msg.ReadyTime.IsZero() {
		msg.ReadyTime = time.Now().Add(time.Duration(msg.Delay) * time.Second)
	}
	if msg.Key == "" {
		msg.Key = uuid.NewString()
	}

	return q.delay.Produce(ctx, msg.ReadyTime.Unix()/60, msg)
}

// NewDelayQueue create a new DelayQueue
func NewDelayQueue() DelayQueue {
	return redis_broker.NewDelayQueue()
}

// NewRealTimeQueue create a new RealTimeQueue
func NewRealTimeQueue() RealTimeQueue {
	return nil
}
