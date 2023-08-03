package broker

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/beihai0xff/pudding/api/gen/pudding/types/v1"
	"github.com/beihai0xff/pudding/app/broker/connector"
	type2 "github.com/beihai0xff/pudding/app/broker/pkg/types"
	"github.com/beihai0xff/pudding/app/broker/storage"
	"github.com/beihai0xff/pudding/configs"
	"github.com/beihai0xff/pudding/pkg/clock"
	"github.com/beihai0xff/pudding/pkg/cluster"
	"github.com/beihai0xff/pudding/pkg/errno"
	"github.com/beihai0xff/pudding/pkg/log"
)

var (
	// error when the message delay is invalid
	errInvalidMessageDelay = errors.New("message delay must be greater than 0")
	// error when the message ready time is invalid
	errInvalidMessageReady = errors.New("DeliverAt must be greater than the current time")
	// check scheduler is implemented Scheduler interface
	_ Scheduler = (*scheduler)(nil)
)

const (
	prefixTimeLocker = "pudding_locker_time:%d"
)

// Scheduler interface
type Scheduler interface {
	// Run start the scheduler
	Run()
	// Produce produce a Message to DelayStorage
	Produce(ctx context.Context, msg *types.Message) error
	// NewConsumer consume Messages from the RealTime Connector queue
	NewConsumer(topic, group string, batchSize int, fn type2.HandleMessage) error
	// Close close the scheduler
	Close() error
}

type scheduler struct {
	delay     storage.DelayStorage
	connector connector.RealTimeConnector
	// wallClock wall wallClock time
	wallClock clock.Clock

	cluster cluster.Cluster

	// messageTopic default message topic
	messageTopic string
	// tokenTopic default token topic
	tokenTopic string

	// token timeSlice token channel
	token chan uint64
	// quit signal quit channel
	quit chan int64
}

// New create a new scheduler
func New(config *configs.BrokerConfig, delay storage.DelayStorage,
	realtime connector.RealTimeConnector) (Scheduler, error) {
	clusterManager, err := cluster.New(config.ServerConfig.EtcdURLs)
	if err != nil {
		return nil, err
	}

	q := &scheduler{
		delay:        delay,
		connector:    realtime,
		wallClock:    clock.New(),
		cluster:      clusterManager,
		messageTopic: config.ServerConfig.MessageTopic,
		tokenTopic:   config.ServerConfig.TokenTopic,
		token:        make(chan uint64),
		quit:         make(chan int64),
	}

	return q, nil
}

// Run start the scheduler
func (s *scheduler) Run() {
	go s.tryProduceToken()

	s.getToken()
	go s.startSchedule()
}

/*
	Produce or Consume DeliverAfter scheduler
*/

// Produce a Message to the delay queue
func (s *scheduler) Produce(ctx context.Context, msg *types.Message) error {
	var err error

	if err = s.checkParams(msg); err != nil {
		log.Errorf("check message params failed: %v", err)
		return fmt.Errorf("check message params failed: %w", err)
	}

	for i := 0; i < 3; i++ {
		err = s.delay.Produce(ctx, msg)
		if err == nil {
			log.Infof("success produce message: %s", msg.String())
			break
		}

		log.Errorf("DelayStorage: failed to produce message: %v, retry in [%d] times", err, i)
		// if produce failed, will retry in three times
		// but if the error is ErrDuplicateMessage, will not retry
		if errors.Is(err, errno.ErrDuplicateMessage) {
			break
		}
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

		msg.DeliverAt = uint64(s.wallClock.Now().Unix()) + msg.DeliverAfter
	} else if time.Unix(int64(msg.DeliverAt), 0).Before(s.wallClock.Now()) {
		return errInvalidMessageReady
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

// startSchedule start a scheduler to consume DelayStorage
// and move delayed messages to RealTimeConnector
func (s *scheduler) startSchedule() {
	log.Infof("start Scheduling")

	for {
		select {
		case t := <-s.token:
			go s.consumeDelayMessage(t)
		case <-s.quit:
			break
		}
	}
}

func (s *scheduler) consumeDelayMessage(t uint64) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	// lock the timeSlice
	name := s.getLockerName(t)

	locker, err := s.cluster.Mutex(name, time.Second*3)
	if err != nil {
		return
	}

	if err = locker.Lock(ctx); err != nil {
		if !errors.Is(err, cluster.ErrLocked) {
			log.Errorf("failed to get timeSlice locker [%s]: %v", name, err)
		}

		return
	}

	if err := s.delay.Consume(ctx, t, 100, s.produceRealTime); err != nil {
		log.Errorf("failed to consume timeSlice [%d] message: %v", t, err)
	}

	// Release the cluster
	if err := locker.Unlock(ctx); err != nil {
		log.Errorf("failed to release timeSlice locker [%s]: %v", name, err)
	}
}

func (s *scheduler) getLockerName(t uint64) string {
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

// NewConsumer create a consumer to consume the queue in connector
func (s *scheduler) NewConsumer(topic, group string, batchSize int, fn type2.HandleMessage) error {
	return s.connector.NewConsumer(topic, group, batchSize, fn)
}

// Close close the scheduler
func (s *scheduler) Close() error {
	close(s.quit)
	return s.connector.Close()
}
