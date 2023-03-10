package broker

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/beihai0xff/pudding/api/gen/pudding/types/v1"
	type2 "github.com/beihai0xff/pudding/app/broker/pkg/types"
	"github.com/beihai0xff/pudding/pkg/lock"
	"github.com/beihai0xff/pudding/pkg/log"
)

// token is the time of the token bucket
// used to consume delayed message every second

const (
	prefixToken       = "pudding_token:"
	prefixTokenLocker = "pudding_locker_token:"
)

/*
	Produce or Consume token
*/

// try to produce token to bucket
func (s *scheduler) tryProduceToken() {
	log.Infof("start produce token goroutine")

	now := s.wallClock.Now()
	timer := time.NewTimer(time.Until(now) + time.Second)

	// wait for the next second
	<-timer.C

	tick := time.NewTicker(1 * time.Second)
	for {
		select {
		case t := <-tick.C:
			// get token name
			tokenName := s.formatTokenName(uint64(t.Unix()))

			// try to lock the token
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			locker, err := lock.NewRedLock(context.Background(), s.formatTokenLockerName(t.Unix()), time.Millisecond*500)
			if err != nil {
				// if the token is not locked, but the error is not ErrNotObtained, log it
				if !errors.Is(err, lock.ErrNotObtained) {
					log.Errorf("failed to get token lock: %s, caused by %v", tokenName, err)
				}

				// if the token is locked, skip it
				continue
			}

			// if got the token lock, send it to token topic
			if err := s.produceRealTime(ctx, &types.Message{Topic: s.tokenTopic, Payload: []byte(tokenName)}); err != nil {
				log.Errorf("failed to produce token: %s, caused by %v", tokenName, err)
				continue
			}

			log.Infof("success produce token: %s", tokenName)

			// extends the lock with a new TTL
			if err := locker.Refresh(ctx, 3*time.Second); err != nil {
				log.Errorf("failed to refresh locker: %s, caused by %v", tokenName, err)
			}
			cancel()
		case <-s.quit:
			break
		}
	}
}

// try to consume token and send to token channel
func (s *scheduler) getToken() {
	log.Infof("start consume token")

	if err := s.connector.NewConsumer(type2.TokenTopic, type2.TokenGroup, 1,
		func(ctx context.Context, msg *types.Message) error {
			log.Infof("get token: %s", string(msg.Payload))

			t := s.parseNowFromToken(string(msg.Payload))
			if t <= 0 {
				return fmt.Errorf("failed to parse token: %s", string(msg.Payload))
			}
			s.token <- t
			return nil
		}); err != nil {
		log.Errorf("failed to get token, caused by %v", err)
		panic(err)
	}
}

func (s *scheduler) formatTokenName(timestamp uint64) string {
	return fmt.Sprintf(prefixToken+"%d", timestamp)
}

func (s *scheduler) formatTokenLockerName(timestamp int64) string {
	return fmt.Sprintf(prefixTokenLocker+"%d", timestamp)
}

// parseNowFromToken parse token from token name
// if return value is -1, means parse failed
func (s *scheduler) parseNowFromToken(token string) uint64 {
	if strings.HasPrefix(token, prefixToken) {
		t, err := strconv.ParseUint(token[len(prefixToken):], 10, 64)
		if err != nil {
			log.Errorf("failed to parse token token: %s, caused by %v", token, err)
			return 0
		}
		return t
	}

	log.Errorf("failed to parse token, token token: %s is invalid", token)
	return 0
}
