package scheduler

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/beihai0xff/pudding/pkg/lock"
	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/types"
)

// token is the time of the token bucket
// used to consume delayed message every second

const prefixToken = "pudding_token:"

/*
	Produce or Consume token
*/

// try to produce token to bucket
func (s *Schedule) tryProduceToken() {
	now := time.Now()
	timer := time.NewTimer(time.Unix(now.Unix()+1, 0).Sub(time.Now()) - time.Millisecond)

	// wait for the next second
	_ = <-timer.C

	tick := time.NewTicker(time.Duration(s.interval) * time.Second)
	for {
		select {
		case t := <-tick.C:
			// get token name
			tokenName := s.formatTokenName(t.Unix())

			// try to lock the token
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			locker, err := lock.NewRedLock(context.Background(), tokenName, time.Millisecond*500)
			if err != nil {
				// if the token is not locked, but the error is not ErrNotObtained, log it
				if err != lock.ErrNotObtained {
					log.Errorf("failed to get token lock: %s, caused by %v", tokenName, err)
				}

				// if the token is locked, skip it
				continue
			}

			// if get token lock, send it to token topic
			if err := s.produceRealTime(ctx, &types.Message{Topic: types.TokenTopic, Payload: []byte(tokenName)}); err != nil {
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

// try to consume token and send to channel
func (s *Schedule) getToken(token chan int64) {
	if err := s.realtime.NewConsumer(types.TokenTopic, types.TokenGroup, 1,
		func(ctx context.Context, msg *types.Message) error {
			t := s.parseNowFromToken(string(msg.Payload))
			if t <= 0 {
				return fmt.Errorf("failed to parse token: %s", string(msg.Payload))
			}
			token <- t
			return nil
		}); err != nil {
		log.Errorf("failed to get token, caused by %v", err)
		panic(err)
	}
	return
}

func (s *Schedule) formatTokenName(time int64) string {
	return fmt.Sprintf(prefixToken+"%d", time)
}

// parseNowFromToken parse token from token name
// if return value is -1, means parse failed
func (s *Schedule) parseNowFromToken(token string) int64 {
	if strings.HasPrefix(token, prefixToken) {
		t, err := strconv.ParseInt(token[len(prefixToken):], 10, 64)
		if err != nil {
			log.Errorf("failed to parse token token: %s, caused by %v", token, err)
			return -1
		}
		return t
	}

	log.Errorf("failed to parse token, token token: %s is invalid", token)
	return -1
}
