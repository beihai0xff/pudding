package lock

import (
	"context"
	"errors"
	"time"

	"github.com/bsm/redislock"

	"github.com/beihai0xff/pudding/pkg/configs"
	rdb "github.com/beihai0xff/pudding/pkg/redis"
)

var r *rdb.Client

func Init() {
	r = rdb.New(configs.GetRedisConfig())
}

type RedLock struct {
	locker *redislock.Lock
}

func NewRedLock(ctx context.Context, name string, expireTime time.Duration) (Lock, error) {
	locker, err := r.GetDistributeLock(ctx, name, expireTime)
	if err != nil {
		if errors.Is(err, redislock.ErrNotObtained) {
			return nil, ErrNotObtained
		}
		return nil, err
	}
	return &RedLock{locker: locker}, nil
}
func (r *RedLock) Release(ctx context.Context) error {
	if err := r.locker.Release(ctx); err != nil {
		if errors.Is(err, redislock.ErrLockNotHeld) {
			return ErrLockNotHeld
		}
		return err
	}

	return nil
}
func (r *RedLock) Refresh(ctx context.Context, ttl time.Duration) error {
	return r.locker.Refresh(ctx, ttl, nil)
}
