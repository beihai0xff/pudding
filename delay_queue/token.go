package delay_queue

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
			tokenName := q.formatTokenName(t.Unix())

			// try to lock the token
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			locker, err := lock.NewRedLock(context.Background(), tokenName, time.Millisecond*500)
			if err != nil {
				// if the token is locked, skip it

				// if the token is not locked, but the error is not ErrNotObtained, log it
				if err != lock.ErrNotObtained {
					log.Errorf("failed to get token lock: %s, caused by %v", tokenName, err)
				}

				continue
			}

			// if get token produce lock, send it to token topic
			if err := q.ProduceRealTime(ctx, &types.Message{Topic: types.TokenTopic, Payload: []byte(tokenName)}); err != nil {
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
func (q *Queue) getToken(token chan int64) {
	if err := q.realtime.NewConsumer(types.TokenTopic, types.TokenGroup, 1,
		func(ctx context.Context, msg *types.Message) error {
			token <- q.parseNowFromToken(string(msg.Payload))
			return nil
		}); err != nil {
		log.Errorf("failed to get token, caused by %v", err)
		panic(err)
	}
	return
}

func (q *Queue) formatTokenName(time int64) string {
	return fmt.Sprintf("key_token_%d", time)
}

// parseNowFromToken parse token from token name
// if return value is -1, means parse failed
func (q *Queue) parseNowFromToken(token string) int64 {
	if strings.HasPrefix(token, "key_token_") {
		t, err := strconv.ParseInt(token[10:], 10, 64)
		if err != nil {
			log.Errorf("failed to parse token token: %s, caused by %v", token, err)
			return -1
		}
		return t
	}

	log.Errorf("failed to parse token, token token: %s is invalid", token)
	return -1
}
