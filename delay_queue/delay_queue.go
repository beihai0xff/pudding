package delay_queue

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis_rate/v10"
	"github.com/google/uuid"

	"github.com/beihai0xff/pudding/delay_queue/broker/redis_broker"
	"github.com/beihai0xff/pudding/pkg/configs"
	"github.com/beihai0xff/pudding/pkg/lock"
	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/pkg/random"
	rdb "github.com/beihai0xff/pudding/pkg/redis"
	"github.com/beihai0xff/pudding/types"
)

type DelayQueue interface {
	// Produce produce a Message to DelayQueue
	Produce(ctx context.Context, quantum string, msg *types.Message) error
	// Consume New a consumer to consume Messages from the queue
	Consume(ctx context.Context, quantum string, now, batchSize int64, fn types.HandleMessage) error
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

// NewDelayQueue create a new DelayQueue
func NewDelayQueue(rdb *rdb.Client) DelayQueue {
	return redis_broker.NewDelayQueue(rdb)
}

// NewRealTimeQueue create a new RealTimeQueue
func NewRealTimeQueue() RealTimeQueue {
	return nil
}

type Queue struct {
	delay    DelayQueue
	realtime RealTimeQueue
	config   *configs.DelayQueueConfig
	// partition interval (Seconds)
	interval int64

	// rate limiter
	limiter *redis_rate.Limiter
}

func NewQueue() *Queue {
	client := rdb.New(configs.GetRedisConfig())
	q := &Queue{
		delay:    NewDelayQueue(client),
		realtime: NewRealTimeQueue(),
		config:   configs.GetDelayQueueConfig(),
	}

	// parse Polling delay queue interval
	t, err := time.ParseDuration(q.config.PartitionInterval)
	if err != nil {
		panic(fmt.Errorf("failed to parse '%s' to time.Duration: %v", q.config.PartitionInterval, err))
	}
	q.interval = int64(t.Seconds())

	// init rate limiter
	q.limiter = client.GetLimiter()
	if err != nil {
		panic(err)
	}

	return q
}

/*
	Produce or Consume Delay Queue
*/

func (q *Queue) Produce(ctx context.Context, msg *types.Message) error {
	var err error

	if err = q.checkParams(msg); err != nil {
		log.Errorf("check message params failed: %v", err)
		return err
	}

	quantum := q.getQuantum(msg.ReadyTime)
	for i := 0; i < 3; i++ {
		if err = q.delay.Produce(ctx, quantum, msg); err != nil {
			log.Errorf("DelayQueue: failed to produce message to Partition %s, err: %v, retry in %d times",
				quantum, err, i)
		} else {
			break
		}
	}
	return err
}

func (q *Queue) checkParams(msg *types.Message) error {
	// if ReadyTime is set, use it
	// otherwise use current time
	if msg.ReadyTime <= 0 {
		if msg.Delay == 0 {
			return errors.New("message delay must be greater than 0")
		}
		msg.ReadyTime = time.Now().Unix() + msg.Delay
	} else {
		if time.Unix(msg.ReadyTime, 0).Before(time.Now()) {
			return errors.New("ReadyTime must be greater than the current time")
		}
	}

	if msg.Key == "" {
		msg.Key = uuid.NewString()
	}

	return nil
}

func (q *Queue) getQuantum(readyTime int64) string {
	startAt := (readyTime / q.interval) * q.interval
	endAt := startAt + q.interval
	return fmt.Sprintf("%d-%d", startAt, endAt)
}

func (q *Queue) startConsumer(quit, token chan int64) error {
	for {
		select {
		case t := <-token:
			ctx := context.Background()
			quantum := q.getQuantum(t)

			// lock the token
			consumer := "consumer" + q.getConsumer(time.Now().Unix())
			locker, err := lock.NewRedLock(context.Background(), consumer, time.Second*3)
			if err != nil {
				if err != lock.ErrNotObtained {
					log.Errorf("failed to get consumer lock: %s, caused by %v", consumer, err)
				}
				continue
			}

			if err := q.delay.Consume(ctx, quantum, t, 100,
				q.ProduceRealTime); err != nil {
				log.Errorf("failed to consume quantum: %s, time token is: %d,caused by %v", quantum, t, err)
			}

			// Release the lock
			locker.Release(ctx)
		case <-quit:
			break
		}
	}

}

func (q *Queue) getConsumer(time int64) string {
	return fmt.Sprintf("key_token_%d", time)
}

/*
	Produce or Consume RealTime Queue
*/

// ProduceRealTime produce a Message to the queue in realtime
func (q *Queue) ProduceRealTime(ctx context.Context, msg *types.Message) error {
	var err error
	for i := 0; i < 3; i++ {
		if err = q.realtime.Produce(ctx, msg); err != nil {
			log.Errorf("RealTimeQueue: failed to produce message to topic %s, err: %v, retry in %d times",
				msg.Topic, err, i)
		} else {
			break
		}
	}
	return err
}

/*
	Produce or Consume token
*/

// try to produce token to bucket
func (q *Queue) tryProduceToken() {
	now := time.Now()
	timer := time.NewTimer(time.Unix(now.Unix()+1, 0).Sub(time.Now()) - time.Millisecond)

	_ = <-timer.C // 从定时器拿数据

	tick := time.NewTicker(1 * time.Second)
	for {
		select {
		case t := <-tick.C:
			// get token name
			tokenName := q.getConsumer(t.Unix())

			// try to lock the token
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			locker, err := lock.NewRedLock(context.Background(), tokenName, time.Millisecond*500)
			if err != nil {
				if err != lock.ErrNotObtained {
					log.Errorf("failed to get token lock: %s, caused by %v", tokenName, err)
				}

				continue
			}

			// TODO: set token to realtime queue

			if err := q.ProduceRealTime(ctx, &types.Message{Topic: "token", Key: tokenName}); err != nil {
				log.Errorf("failed to produce token: %s, caused by %v", tokenName, err)
			}

			// extends the lock with a new TTL
			if err := locker.Refresh(ctx, 3*time.Second); err != nil {
				log.Errorf("failed to refresh locker: %s, caused by %v", tokenName, err)
			}
			cancel()
		}
	}
}

func (q *Queue) getTokenName(time int64) string {
	return fmt.Sprintf("key_token_%d", time)
}

/*
	rate limit
*/

func (q *Queue) startLimiter(token chan int) {

	for {
		res, err := q.limiter.Allow(context.Background(), "pudding:rate_every_second", redis_rate.PerSecond(1))
		if err != nil {
			log.Errorf("failed to allow limiter: %v", err)
		}
		if res.Allowed == 1 {
			token <- 1
		}
		time.Sleep(time.Duration(random.GetRand(500, 1000)) * time.Millisecond)
	}
}
