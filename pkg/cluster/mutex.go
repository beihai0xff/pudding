// Package cluster provides a cluster manager.
// mutex.go contains the implementation of Mutex.
package cluster

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"go.etcd.io/etcd/client/v3/concurrency"

	"github.com/beihai0xff/pudding/pkg/log"
)

const minLockTTL = time.Second

var (
	// ErrInvalidTTL is returned when a TTL is less than minLockTTL.
	ErrInvalidTTL = fmt.Errorf("invalid TTL, must be greater than %s", minLockTTL.String())

	// ErrLocked is returned when a lock locked by another session.
	ErrLocked = concurrency.ErrLocked
	// ErrLockNotHeld is returned when a lock is not held
	ErrLockNotHeld = errors.New("lock not held")
)

// Mutex is a cluster level mutex.
type Mutex interface {
	Lock(ctx context.Context) error
	// Unlock releases the lock.
	// May return ErrLockNotHeld.
	Unlock(ctx context.Context) error

	// Refresh extends the lock with a new TTL.
	// recommended use it when keepAlive is false
	// will return ErrLockNotHeld if refresh is unsuccessful.
	Refresh(ctx context.Context) error
}

type mutex struct {
	// concurrency.Mutex is a session level mutex,
	// on client side without concurrency protection,
	// so sync.Mutex is required to make it goroutine safe
	lock    sync.Mutex
	m       *concurrency.Mutex
	session *concurrency.Session
	// keepAlive indicates whether the lock is refreshed periodically.
	// If keepAlive is true, the lock will be refreshed periodically.
	// If keepAlive is false, the lock will be released when the session is closed.
	keepAlive bool
}

func (m *mutex) Lock(ctx context.Context) (err error) {
	isTimeout := true

	m.lock.Lock()
	defer func() {
		if isTimeout || err != nil {
			if errors.Is(err, concurrency.ErrLocked) {
				log.Errorf("lock failed: %v", err)
				err = ErrLocked
			}
			m.lock.Unlock()
		}
	}()

	err = m.m.Lock(ctx)
	if !m.keepAlive {
		m.session.Orphan()
	}

	isTimeout = false
	return
}

func (m *mutex) Unlock(ctx context.Context) error {
	defer m.lock.Unlock()

	return m.m.Unlock(ctx)
}

func (m *mutex) Refresh(ctx context.Context) error {
	_, err := m.session.Client().KeepAliveOnce(ctx, m.session.Lease())
	if errors.Is(err, concurrency.ErrSessionExpired) {
		return ErrLockNotHeld
	}
	return err
}

func (c *cluster) Mutex(name string, ttl time.Duration, opts ...MutexOption) (Mutex, error) {
	if ttl <= minLockTTL {
		return nil, ErrInvalidTTL
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.opts.requestTimeout)
	defer cancel()
	session, err := c.getSession(ctx, int64(ttl.Seconds()))
	if err != nil {
		return nil, err
	}

	m := mutex{
		m:         concurrency.NewMutex(session, name),
		lock:      sync.Mutex{},
		session:   session,
		keepAlive: true,
	}

	for _, opt := range opts {
		opt(&m)
	}

	return &m, nil
}

// MutexOption is a function that applies an option to Mutex.
type MutexOption func(*mutex)

// WithDisableKeepalive disables the keepalive feature of Mutex.
func WithDisableKeepalive() MutexOption {
	return func(m *mutex) {
		m.keepAlive = false
	}
}
