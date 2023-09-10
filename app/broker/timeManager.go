package broker

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/beihai0xff/pudding/pkg/cluster"
	"github.com/beihai0xff/pudding/pkg/log"
)

// token is the time of the token bucket
// used to consume delayed message every second

const (
	prefixToken       = "pudding_token:"
	prefixTokenLocker = "pudding_locker_token:"
)

type timeManager struct {
	// tokenTopic default token topic
	tokenTopic string
	// tokenQueue token queue
	tokenQueue cluster.Queue

	cluster cluster.Cluster

	// msgChan token message channel
	msgChan chan *cluster.Message
	// quit signal quit channel
	quit chan struct{}
}

func newTimeManager(tokenTopic string, clus cluster.Cluster, quit chan struct{}) (*timeManager, error) {
	queue, err := clus.Queue(tokenTopic)
	if err != nil {
		return nil, fmt.Errorf("failed to new token queue: %w", err)
	}

	return &timeManager{
		tokenTopic: tokenTopic,
		tokenQueue: queue,
		msgChan:    make(chan *cluster.Message, 1),
		cluster:    clus,
		quit:       quit,
	}, nil
}

/*
	Produce or Consume token
*/

// produceTokenWorker try to produce token to bucket
func (s *timeManager) produceTokenWorker() {
	log.Infof("start produce token goroutine")

	now := s.cluster.WallClock()
	timer := time.NewTimer(time.Until(now) + time.Second)

	// wait for the next second
	<-timer.C

	tick := time.NewTicker(time.Second)

	for {
		select {
		case t := <-tick.C:
			if err := s.sendToken(&t); err != nil {
				log.Infof("produce token failed: %v", err)
			}
		case <-s.quit:
			break
		}
	}
}

func (s *timeManager) sendToken(t *time.Time) error {
	// get token name
	tokenName := s.formatTokenName(uint64(t.Unix()))

	// try to lock the token
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	locker, err := s.cluster.Mutex(s.formatTokenLockerName(t.Unix()), time.Second,
		cluster.WithDisableKeepalive())
	if err != nil {
		return fmt.Errorf("failed to get token locker [%s]: %w", tokenName, err)
	}

	if err = locker.Lock(ctx); err != nil {
		return fmt.Errorf("failed to get token locker [%s]: %w", tokenName, err)
	}

	// if got the token locker, send it to the token topic
	if _, err := s.tokenQueue.Produce(&cluster.Message{Key: tokenName, Unique: true}); err != nil {
		return fmt.Errorf("failed to produce token [%s]: %w", tokenName, err)
	}

	log.Infof("success produce token: %s", tokenName)

	// extends the locker with a new TTL
	if err := locker.Refresh(ctx); err != nil {
		log.Errorf("failed to refresh locker [%s]: %v", tokenName, err)
	}

	return nil
}

// try to consume token
func (s *timeManager) consumeToken(cb func(uint64) error) {
	log.Infof("start consume token")

	go s.reader()

	for {
		select {
		case <-s.quit:
			break
		case msg := <-s.msgChan:
			t := s.parseTimeFromToken(msg.Key)
			if t <= 0 {
				log.Errorf("failed to parse token from token key [%s]", msg.Key)
				continue
			}

			for err := cb(t); err != nil; err = cb(t) {
				log.Errorf("failed to consume token: %s, caused by %v", msg.Key, err)
			}
			log.Infof("success consume token: %s", msg.Key)

			if err := s.tokenQueue.Commit(msg); err != nil {
				log.Errorf("failed to commit token: %s, caused by %v", msg.Key, err)
			}
		}
	}
}

func (s *timeManager) reader() {
	for {
		msg, err := s.tokenQueue.Consume(context.Background())

		if err != nil {
			log.Errorf("failed to consume token, caused by %v", err)
			time.Sleep(time.Second)

			continue
		}

		log.Infof("get token: %s", msg.Key)
		s.msgChan <- msg
	}
}

func (s *timeManager) formatTokenName(timestamp uint64) string {
	return fmt.Sprintf(prefixToken+"%d", timestamp)
}

func (s *timeManager) formatTokenLockerName(timestamp int64) string {
	return fmt.Sprintf(prefixTokenLocker+"%d", timestamp)
}

// parseTimeFromToken parse token from token name
// if return value is -1, means parse failed
func (s *timeManager) parseTimeFromToken(token string) uint64 {
	if strings.HasPrefix(token, prefixToken) {
		t, err := strconv.ParseUint(token[len(prefixToken):], 10, 64)
		if err != nil {
			log.Errorf("failed to parse token [%s]: %v", token, err)

			return 0
		}

		return t
	}

	log.Errorf("failed to parse token, token [%s] is invalid", token)

	return 0
}
