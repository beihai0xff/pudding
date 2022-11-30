package scheduler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/beihai0xff/pudding/app/scheduler/broker"
	"github.com/beihai0xff/pudding/app/scheduler/connector"
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
	prefixTimeLocker = "pudding_locker_time:%d"
)

// nolint:lll
//go:generate mockgen -destination=../../test/mock/scheduler_mock.go -package=mock github.com/beihai0xff/pudding/app/scheduler Scheduler

// Scheduler interface
type Scheduler interface {
	// Run start the scheduler
	Run()
	// Produce produce a Message to DelayBroker
	Produce(ctx context.Context, msg *types.Message) error
	// NewConsumer consume Messages from the RealTime Connector queue
	NewConsumer(topic, group string, batchSize int, fn types.HandleMessage) error
}

type scheduler struct {
	delay     broker.DelayBroker
	connector connector.RealTimeConnector
	// wallClock wall wallClock time
	wallClock clock.Clock

	// messageTopic default message topic
	messageTopic string
	// tokenTopic default token topic
	tokenTopic string

	// token timeSlice token channel
	token chan int64
	// quit signal quit channel
	quit chan int64
}

func NewQueue(config *configs.SchedulerConfig) (broker.DelayBroker, connector.RealTimeConnector) {
	return broker.NewDelayBroker(config.Broker), connector.NewConnector(config.Connector)
}

func New(config *configs.SchedulerConfig, delay broker.DelayBroker, realtime connector.RealTimeConnector) Scheduler {
	q := &scheduler{
		delay:        delay,
		connector:    realtime,
		wallClock:    clock.New(),
		messageTopic: config.MessageTopic,
		tokenTopic:   config.TokenTopic,
		token:        make(chan int64),
		quit:         make(chan int64),
	}

	return q
}

func (s *scheduler) Run() {
	go s.tryProduceToken()
	s.getToken()
	if err := s.startSchedule(); err != nil {
		log.Errorf("start scheduler failed: %v", err)
		panic(err)
	}
}

/*
	Produce or Consume DeliverAfter scheduler
*/

// Produce produce a Message to the delay queue
func (s *scheduler) Produce(ctx context.Context, msg *types.Message) error {
	var err error

	if err = s.checkParams(msg); err != nil {
		log.Errorf("check message params failed: %v", err)
		return fmt.Errorf("check message params failed: %w", err)
	}

	for i := 0; i < 3; i++ {
		err = s.delay.Produce(ctx, msg)
		if err == nil {
			break
		}
		// if produce failed, retry in three times
		log.Errorf("DelayBroker: failed to produce message: %v, retry in [%d] times", err, i)
	}
	return err
}

// checkParams check the Produced message params
func (s *scheduler) checkParams(msg *types.Message) error {
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
		msg.Topic = s.messageTopic
	}

	return nil
}

// startSchedule start a scheduler to consume DelayBroker
// and move delayed messages to RealTimeConnector
func (s *scheduler) startSchedule() error {
	log.Infof("start Scheduler")

	for {
		select {
		case t := <-s.token:
			ctx := context.Background()

			// lock the timeSlice
			name := s.getLockerName(t)
			locker, err := lock.NewRedLock(ctx, name, time.Second*3)
			if err != nil {
				if !errors.Is(err, lock.ErrNotObtained) {
					log.Errorf("failed to get timeSlice locker: %s, caused by %v", name, err)
				}
				continue
			}

			if err := s.delay.Consume(ctx, t, 100, s.produceRealTime); err != nil {
				log.Errorf("failed to consume time: %d, caused by %v", t, err)
			}

			// Release the lock
			if err := locker.Release(ctx); err != nil {
				log.Errorf("failed to release time locker: %s, caused by %v", name, err)
			}

		case <-s.quit:
			break
		}
	}

}

func (s *scheduler) getLockerName(t int64) string {
	return fmt.Sprintf(prefixTimeLocker, t)
}

/*
	Produce or Consume RealTime scheduler
*/

// produceRealTime produce a Message to the queue in connector
func (s *scheduler) produceRealTime(ctx context.Context, msg *types.Message) error {
	var err error
	for i := 0; i < 3; i++ {
		if err = s.connector.Produce(ctx, msg); err != nil {
			log.Errorf("RealTimeConnector: failed to produce message to topic %s, err: %v, retry in %d times",
				msg.Topic, err, i)
		} else {
			break
		}
	}
	return err
}

func (s *scheduler) NewConsumer(topic, group string, batchSize int, fn types.HandleMessage) error {
	return s.connector.NewConsumer(topic, group, batchSize, fn)
}

func (s *scheduler) Close() error {
	close(s.quit)
	return s.connector.Close()
}
