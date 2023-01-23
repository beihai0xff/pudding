// Package lock implements the distributed lock
package lock

import (
	"context"
	"errors"
	"time"

	"github.com/bsm/redislock"

	rdb "github.com/beihai0xff/pudding/pkg/redis"
)

var r *redislock.Client

// Init init the lock module
func Init(client *rdb.Client) {
	r = redislock.New(client.GetClient())
}

// RedLock is the redis distributed lock implement
type RedLock struct {
	locker *redislock.Lock
	// Create a new locker client.
}

// NewRedLock create a new redlock
func NewRedLock(ctx context.Context, name string, expireTime time.Duration) (Lock, error) {
	if r == nil {
		return nil, errors.New("redis client is nil, please init the lock module first")
	}

	// Retry every 100ms, for up-to 3x
	backoff := redislock.LimitRetry(redislock.LinearBackoff(100*time.Millisecond), 3)
	locker, err := r.Obtain(ctx, name, expireTime, &redislock.Options{
		RetryStrategy: backoff,
	})
	if err != nil {
		if errors.Is(err, redislock.ErrNotObtained) {
			return nil, ErrNotObtained
		}
		return nil, err
	}
	return &RedLock{locker: locker}, nil
}

// Release release the lock
func (r *RedLock) Release(ctx context.Context) error {
	if err := r.locker.Release(ctx); err != nil {
		if errors.Is(err, redislock.ErrLockNotHeld) {
			return ErrLockNotHeld
		}
		return err
	}

	return nil
}

// Refresh extend the lock TTL
func (r *RedLock) Refresh(ctx context.Context, ttl time.Duration) error {
	return r.locker.Refresh(ctx, ttl, nil)
}
