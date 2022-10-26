package scheduler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis_rate/v10"
	"github.com/google/uuid"

	"github.com/beihai0xff/pudding/pkg/configs"
	"github.com/beihai0xff/pudding/pkg/lock"
	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/pkg/mq/pulsar"
	"github.com/beihai0xff/pudding/pkg/random"
	rdb "github.com/beihai0xff/pudding/pkg/redis"
	"github.com/beihai0xff/pudding/scheduler/broker/pulsar_broker"
	"github.com/beihai0xff/pudding/scheduler/broker/redis_broker"
	"github.com/beihai0xff/pudding/types"
)

const prefixTimeSliceLocker = "key_locker_time_%d"

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

type Scheduler struct {
	delay    DelayQueue
	realtime RealTimeQueue
	config   *configs.DelayQueueConfig
	// timeSlice interval (Seconds)
	interval int64

	// rate limiter
	limiter *redis_rate.Limiter

	// timeSlice token channel
	token chan int64
	// quit signal channel
	quit chan int64
}

func New() *Scheduler {
	redisClient := rdb.New(configs.GetRedisConfig())
	pulsarClient := pulsar.New(configs.GetPulsarConfig())
	q := &Scheduler{
		delay:    NewDelayQueue(redisClient),
		realtime: NewRealTimeQueue(pulsarClient),
		config:   configs.GetDelayQueueConfig(),
		token:    make(chan int64),
		quit:     make(chan int64),
	}

	// parse Polling delay queue interval
	t, err := time.ParseDuration(q.config.PartitionInterval)
	if err != nil {
		panic(fmt.Errorf("failed to parse '%s' to time.Duration: %v", q.config.PartitionInterval, err))
	}
	q.interval = int64(t.Seconds())

	// init rate limiter
	q.limiter = redisClient.GetLimiter()
	if err != nil {
		panic(err)
	}

	return q
}

func (s *Scheduler) Run() {
	go s.tryProduceToken()
	s.getToken(s.token)
	if err := s.startScheduler(s.quit, s.token); err != nil {
		log.Errorf("start Scheduler failed: %v", err)
		panic(err)
	}
}

/*
	Produce or Consume Delay Scheduler
*/

func (s *Scheduler) Produce(ctx context.Context, msg *types.Message) error {
	var err error

	if err = s.checkParams(msg); err != nil {
		log.Errorf("check message params failed: %v", err)
		return fmt.Errorf("check message params failed: %w", err)
	}

	timeSlice := s.getTimeSlice(msg.ReadyTime)
	for i := 0; i < 3; i++ {
		err = s.delay.Produce(ctx, timeSlice, msg)
		if err == nil {
			break
		}
		// if produce failed, retry in three times
		log.Errorf("DelayQueue: failed to produce message to timeSlice %s, err: %w, retry in %d times",
			timeSlice, err, i)
	}
	return err
}

func (s *Scheduler) checkParams(msg *types.Message) error {
	// if Message.ReadyTime is set, use ReadyTime
	// otherwise use current time + Delay Seconds
	if msg.ReadyTime <= 0 {
		if msg.Delay <= 0 {
			return errors.New("message delay must be greater than 0")
		}
		msg.ReadyTime = time.Now().Unix() + msg.Delay
	} else {
		if time.Unix(msg.ReadyTime, 0).Before(time.Now()) {
			return errors.New("ReadyTime must be greater than the current time")
		}
	}

	// if Message.Key is not set, generate a random ID
	if msg.Key == "" {
		msg.Key = uuid.NewString()
	}

	return nil
}

// getTimeSlice get the time slice of the given time
// Left closed right open interval
// e.g. the given interval is 60, the range is [0, 60)、[60, 120)、[120, 180)...
// 59 => 0~60
// 60 => 60~120
// 61 => 60~120
func (s *Scheduler) getTimeSlice(readyTime int64) string {
	startAt := (readyTime / s.interval) * s.interval
	endAt := startAt + s.interval
	return fmt.Sprintf("%d~%d", startAt, endAt)
}

// startScheduler start a scheduler to consume DelayQueue
// and move delayed messages to RealTimeQueue
func (s *Scheduler) startScheduler(quit, token chan int64) error {
	for {
		select {
		case t := <-token:
			ctx := context.Background()
			timeSlice := s.getTimeSlice(t)

			// lock the timeSlice
			name := s.getLockerName(t)
			locker, err := lock.NewRedLock(context.Background(), name, time.Second*3)
			if err != nil {
				if err != lock.ErrNotObtained {
					log.Errorf("failed to get timeSlice locker: %s, caused by %v", name, err)
				}
				continue
			}

			if err := s.delay.Consume(ctx, timeSlice, t, 100, s.produceRealTime); err != nil {
				log.Errorf("failed to consume timeSlice: %s, time is: %d, caused by %w", timeSlice, t, err)
			}

			// Release the lock
			locker.Release(ctx)
		case <-quit:
			break
		}
	}

}

func (s *Scheduler) getLockerName(t int64) string {
	return fmt.Sprintf(prefixTimeSliceLocker, t)
}

/*
	Produce or Consume RealTime Scheduler
*/

// produceRealTime produce a Message to the queue in realtime
func (s *Scheduler) produceRealTime(ctx context.Context, msg *types.Message) error {
	var err error
	for i := 0; i < 3; i++ {
		if err = s.realtime.Produce(ctx, msg); err != nil {
			log.Errorf("RealTimeQueue: failed to produce message to topic %s, err: %v, retry in %d times",
				msg.Topic, err, i)
		} else {
			break
		}
	}
	return err
}

func (s *Scheduler) NewConsumer(topic, group string, batchSize int, fn types.HandleMessage) error {
	return s.realtime.NewConsumer(topic, group, batchSize, fn)
}

/*
	rate limit
*/

func (s *Scheduler) startLimiter(token chan int) {

	for {
		res, err := s.limiter.Allow(context.Background(), "pudding:rate_every_second", redis_rate.PerSecond(1))
		if err != nil {
			log.Errorf("failed to allow limiter: %v", err)
		}
		if res.Allowed == 1 {
			token <- 1
		}
		time.Sleep(time.Duration(random.GetRand(500, 1000)) * time.Millisecond)
	}
}
