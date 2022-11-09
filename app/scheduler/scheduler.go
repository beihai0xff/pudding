package scheduler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/beihai0xff/pudding/app/scheduler/broker"
	"github.com/beihai0xff/pudding/configs"
	"github.com/beihai0xff/pudding/pkg/clock"
	"github.com/beihai0xff/pudding/pkg/lock"
	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/types"
)

var (
	// error when the message delay is invalid
	errInvalidMessageDelay = errors.New("message delay must be greater than 0")
	// error when the message ready time is invalid
	errInvalidMessageReady = errors.New("DeliverAt must be greater than the current time")
)

const (
	prefixTimeSliceLocker = "pudding_locker_time:%d"
)

// nolint:lll
//go:generate mockgen -destination=../../test/mock/scheduler_mock.go --package=mock github.com/beihai0xff/pudding/app/scheduler Scheduler

// Scheduler interface
type Scheduler interface {
	// Run start the scheduler
	Run()
	// Produce produce a Message to DelayQueue
	Produce(ctx context.Context, msg *types.Message) error
	// NewConsumer consume Messages from the realtime queue
	NewConsumer(topic, group string, batchSize int, fn types.HandleMessage) error
}

type Schedule struct {
	delay    broker.DelayQueue
	realtime broker.RealTimeQueue
	// wallClock wall wallClock time
	wallClock clock.Clock
	// interval timeSlice interval (Seconds)
	interval int64

	// limiter rate limiter
	// limiter *redis_rate.Limiter

	// token timeSlice token channel
	token chan int64
	// quit signal channel
	quit chan int64
}

func New(config *configs.SchedulerConfig, delay broker.DelayQueue, realtime broker.RealTimeQueue) *Schedule {
	q := &Schedule{
		delay:     delay,
		realtime:  realtime,
		wallClock: clock.New(),
		token:     make(chan int64),
		quit:      make(chan int64),
	}

	// parse Polling delay queue interval
	t, err := time.ParseDuration(config.TimeSliceInterval)
	if err != nil {
		panic(fmt.Errorf("failed to parse '%s' to time.Duration: %w", config.TimeSliceInterval, err))
	}
	q.interval = int64(t.Seconds())
	log.Debugf("timeSlice interval is: %d seconds", q.interval)

	// init rate limiter
	// q.limiter = redisClient.GetLimiter()
	// if err != nil {
	// 	panic(err)
	// }

	return q
}

func (s *Schedule) Run() {
	go s.tryProduceToken()
	s.getToken(s.token)
	if err := s.startSchedule(s.quit, s.token); err != nil {
		log.Errorf("start Schedule failed: %v", err)
		panic(err)
	}
}

/*
	Produce or Consume DeliverAfter Schedule
*/

// Produce produce a Message to the delay queue
func (s *Schedule) Produce(ctx context.Context, msg *types.Message) error {
	var err error

	if err = s.checkParams(msg); err != nil {
		log.Errorf("check message params failed: %v", err)
		return fmt.Errorf("check message params failed: %w", err)
	}

	timeSlice := s.getTimeSlice(msg.DeliverAt)
	for i := 0; i < 3; i++ {
		err = s.delay.Produce(ctx, timeSlice, msg)
		if err == nil {
			break
		}
		// if produce failed, retry in three times
		log.Warnf("DelayQueue: failed to produce message to timeSlice %s, err: %w, retry in %d times",
			timeSlice, err, i)
	}
	return err
}

// checkParams check the Produced message params
func (s *Schedule) checkParams(msg *types.Message) error {
	// if Message.DeliverAt is set, use DeliverAt
	// otherwise use current time + DeliverAfter Seconds
	if msg.DeliverAt <= 0 {
		if msg.DeliverAfter <= 0 {
			return errInvalidMessageDelay
		}
		msg.DeliverAt = s.wallClock.Now().Unix() + msg.DeliverAfter
	} else {
		if time.Unix(msg.DeliverAt, 0).Before(s.wallClock.Now()) {
			return errInvalidMessageReady
		}
	}

	// if Message.Key is not set, generate a random ID
	if msg.Key == "" {
		msg.Key = uuid.NewString()
	}

	if msg.Topic == "" {
		msg.Topic = types.DefaultTopic
	}

	return nil
}

// getTimeSlice get the time slice of the given time
// Left closed right open interval
// e.g. the given interval is 60, the range is [0, 60)、[60, 120)、[120, 180)...
// 59 => 0~60
// 60 => 60~120
// 61 => 60~120
func (s *Schedule) getTimeSlice(readyTime int64) string {
	startAt := (readyTime / s.interval) * s.interval
	endAt := startAt + s.interval
	return fmt.Sprintf("%d~%d", startAt, endAt)
}

// startSchedule start a scheduler to consume DelayQueue
// and move delayed messages to RealTimeQueue
func (s *Schedule) startSchedule(quit, token chan int64) error {
	log.Infof("start Scheduler")

	for {
		select {
		case t := <-token:
			ctx := context.Background()
			timeSlice := s.getTimeSlice(t)

			// lock the timeSlice
			name := s.getLockerName(t)
			locker, err := lock.NewRedLock(ctx, name, time.Second*3)
			if err != nil {
				if !errors.Is(err, lock.ErrNotObtained) {
					log.Errorf("failed to get timeSlice locker: %s, caused by %v", name, err)
				}
				continue
			}

			if err := s.delay.Consume(ctx, timeSlice, t, 100, s.produceRealTime); err != nil {
				log.Errorf("failed to consume timeSlice: %s, time is: %d, caused by %w", timeSlice, t, err)
			}

			// Release the lock
			if err := locker.Release(ctx); err != nil {
				log.Errorf("failed to release timeSlice locker: %s, caused by %w", name, err)
			}

		case <-quit:
			break
		}
	}

}

func (s *Schedule) getLockerName(t int64) string {
	return fmt.Sprintf(prefixTimeSliceLocker, t)
}

/*
	Produce or Consume RealTime Schedule
*/

// produceRealTime produce a Message to the queue in realtime
func (s *Schedule) produceRealTime(ctx context.Context, msg *types.Message) error {
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

func (s *Schedule) NewConsumer(topic, group string, batchSize int, fn types.HandleMessage) error {
	return s.realtime.NewConsumer(topic, group, batchSize, fn)
}

/*
	rate limit
*/
//
// func (s *Schedule) startLimiter(token chan int) {
//
// 	for {
// 		res, err := s.limiter.Allow(context.Background(), "pudding:rate_every_second", redis_rate.PerSecond(1))
// 		if err != nil {
// 			log.Errorf("failed to allow limiter: %v", err)
// 		}
// 		if res.Allowed == 1 {
// 			token <- 1
// 		}
// 		time.Sleep(time.Duration(random.GetRand(500, 1000)) * time.Millisecond)
// 	}
// }
