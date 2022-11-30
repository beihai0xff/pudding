package broker

import (
	"context"
	"time"

	"github.com/beihai0xff/pudding/app/scheduler/broker/redis_broker"
	"github.com/beihai0xff/pudding/app/scheduler/pkg/configs"
	"github.com/beihai0xff/pudding/pkg/log"
	rdb "github.com/beihai0xff/pudding/pkg/redis"
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

// NewDelayBroker create a new DelayBroker
func NewDelayBroker(broker string) DelayBroker {
	switch broker {
	case "redis":
		// parse Polling delay queue interval
		interval := configs.GetSchedulerConfig().TimeSliceInterval
		t, err := time.ParseDuration(interval)
		if err != nil {
			log.Fatalf("failed to parse '%s' to time.Duration: %v", interval, err)
		}
		log.Infof("timeSlice interval is: %d seconds", t.Seconds())
		return redis_broker.NewDelayQueue(rdb.New(configs.GetRedisConfig()), int64(t.Seconds()))
	default:
		log.Fatalf("unknown broker type: [%s]", broker)
	}
	return nil
}
