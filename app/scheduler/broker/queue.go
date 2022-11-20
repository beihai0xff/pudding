package broker

import (
	"context"
	"time"

	"github.com/beihai0xff/pudding/app/scheduler/broker/pulsar_broker"
	"github.com/beihai0xff/pudding/app/scheduler/broker/redis_broker"
	"github.com/beihai0xff/pudding/app/scheduler/pkg/configs"
	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/pkg/mq/pulsar"
	rdb "github.com/beihai0xff/pudding/pkg/redis"
	"github.com/beihai0xff/pudding/types"
)

// nolint:lll
//go:generate mockgen -destination=../../../test/mock/queue_mock.go --package=mock github.com/beihai0xff/pudding/app/scheduler/broker DelayQueue,RealTimeQueue

// DelayQueue is a queue to store messages with delay time
type DelayQueue interface {
	// Produce produce a Message to DelayQueue
	Produce(ctx context.Context, msg *types.Message) error
	// Consume consume Messages from the queue
	Consume(ctx context.Context, now, batchSize int64, fn types.HandleMessage) error
	// Close the queue
	Close() error
}

// RealTimeQueue is a queue to store messages in realtime
type RealTimeQueue interface {
	// Produce produce a Message to the queue in real time
	Produce(ctx context.Context, msg *types.Message) error
	// NewConsumer new a consumer to consume Messages from the realtime queue in background
	NewConsumer(topic, group string, batchSize int, fn types.HandleMessage) error
	// Close the queue
	Close() error
}

// NewDelayQueue create a new DelayQueue
func NewDelayQueue(broker string) DelayQueue {
	switch broker {
	case "redis":
		// parse Polling delay queue interval
		interval := configs.GetSchedulerConfig().TimeSliceInterval
		t, err := time.ParseDuration(interval)
		if err != nil {
			log.Fatalf("failed to parse '%s' to time.Duration: %w", interval, err)
		}
		log.Infof("timeSlice interval is: %d seconds", t.Seconds())
		return redis_broker.NewDelayQueue(rdb.New(configs.GetRedisConfig()), int64(t.Seconds()))
	default:
		log.Fatalf("unknown broker type: [%s]", broker)
	}
	return nil
}

// NewRealTimeQueue create a new RealTimeQueue
func NewRealTimeQueue(connector string) RealTimeQueue {
	switch connector {
	case "pulsar":
		return pulsar_broker.NewRealTimeQueue(pulsar.New(configs.GetPulsarConfig()))
	default:
		log.Fatalf("unknown connector type: [%s]", connector)
	}
	return nil

}
