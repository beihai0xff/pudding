package lock

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	etcd "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"

	"github.com/beihai0xff/pudding/pkg/log"
)

// ETCDLockClient is a distributed lock client based on etcd.
type ETCDLockClient struct {
	*etcd.Client
}

// NewETCDLockClient creates a new ETCDLockClient.
func NewETCDLockClient(client *etcd.Client) *ETCDLockClient {
	return &ETCDLockClient{
		Client: client,
	}
}

// NewETCDLock creates a new etcdLock.
func (c ETCDLockClient) NewETCDLock(ctx context.Context, name string, expireTime time.Duration) (Lock, error) {
	var err error
	var session *concurrency.Session
	defer func() {
		if err != nil {
			if session != nil {
				_ = session.Close()
			}
			log.Error(err)
		}
	}()

	session, err = concurrency.NewSession(c.Client, concurrency.WithTTL(int(expireTime.Seconds())))
	if err != nil {
		return nil, fmt.Errorf("create etcd session failed: %w", err)
	}

	session.Lease()
	lock := concurrency.NewMutex(session, name)
	if err = lock.Lock(ctx); err != nil {
		return nil, fmt.Errorf("get mutex failed %w", err)
	}

	if err := lock.TryLock(ctx); err != nil {
		if errors.As(err, &concurrency.ErrLocked) {
			return nil, ErrNotObtained
		}
		return nil, fmt.Errorf("try lock failed: %w", err)
	}

	log.Infof("get locker [%s] success", name)

	return &etcdLock{lock}, nil
}

// etcdLock is a distributed lock based on etcd.
type etcdLock struct {
	locker *concurrency.Mutex
}

// Release releases the lock.
func (l *etcdLock) Release(ctx context.Context) error {
	if err := l.locker.Unlock(ctx); err != nil {
		if errors.As(err, &concurrency.ErrSessionExpired) {
			return ErrLockNotHeld
		}
		return fmt.Errorf("unlock [%s] failed: %w", l.locker.Key(), err)
	}

	return nil
}

// Refresh refreshes the lock.
func (l *etcdLock) Refresh(ctx context.Context, ttl time.Duration) error {
	go func() {
		timeoutCtx, cancel := context.WithTimeout(context.Background(), ttl)
		defer cancel()
		<-timeoutCtx.Done()
		if err := l.Release(ctx); err != nil {
			log.Errorf("refresh lock [%s] failed: %v", l.locker.Key(), err)
		}
	}()

	return nil
}
