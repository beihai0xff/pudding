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
	ErrLocked = errors.New("lock already held by another session")
	// ErrLockedBySelf is returned when a lock locked by self.
	ErrLockedBySelf = errors.New("lock already held by self")
	// ErrLockNotHeld is returned when a lock is not held
	ErrLockNotHeld = errors.New("lock not held")
)

// Mutex is a cluster level mutex.
type Mutex interface {
	Lock(ctx context.Context) error
	// Unlock releases the lock.
	// May return ErrLockNotHeld.
	Unlock(ctx context.Context) error

	// IsLocked returns whether the lock is held.
	IsLocked() bool

	// Refresh extends the lock with TTL.
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

	if !m.lock.TryLock() {
		return ErrLockedBySelf
	}

	defer func() {
		if isTimeout || err != nil {
			if errors.Is(err, concurrency.ErrLocked) {
				err = ErrLocked
			}

			log.Errorf("lock [%s] failed: %v", m.m.Key(), err)
			m.lock.Unlock()

			return
		}

		log.Infof("mutex [%s] locked", m.m.Key())
	}()

	err = m.m.Lock(ctx)
	if !m.keepAlive {
		m.session.Orphan()
	}

	isTimeout = false

	return
}

// IsLocked returns whether the lock is held.
func (m *mutex) IsLocked() bool {
	// if the mutex is locked by self, return true
	if m.lock.TryLock() {
		m.lock.Unlock()
		return false
	}

	if m.session == nil {
		return false
	}

	// check whether the session is expired
	select {
	case <-m.session.Done():
		return false
	default:
		return true
	}
}

func (m *mutex) Unlock(ctx context.Context) error {
	var err error
	lockerKey := m.m.Key()
	defer func() {
		if err == nil {
			m.lock.Unlock()
			log.Infof("mutex [%s] unlocked", lockerKey)
		}
	}()

	err = m.m.Unlock(ctx)

	return err
}

func (m *mutex) Refresh(ctx context.Context) error {
	_, err := m.session.Client().KeepAliveOnce(ctx, m.session.Lease())
	if errors.Is(err, concurrency.ErrSessionExpired) {
		return ErrLockNotHeld
	}

	return err
}

func (c *cluster) Mutex(name string, ttl time.Duration, opts ...MutexOption) (Mutex, error) {
	if ttl < minLockTTL {
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
