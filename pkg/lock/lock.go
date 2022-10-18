package lock

import (
	"context"
	"errors"
	"time"
)

var (
	// ErrNotObtained is returned when a lock cannot be obtained.
	ErrNotObtained = errors.New("lock: not obtained")

	// ErrLockNotHeld is returned when trying to release an inactive lock.
	ErrLockNotHeld = errors.New("lock: lock not held")
)

type Lock interface {
	// Release manually releases the lock.
	// May return ErrLockNotHeld.
	Release(ctx context.Context) error
	// Refresh extends the lock with a new TTL.
	// May return ErrNotObtained if refresh is unsuccessful.
	Refresh(ctx context.Context, ttl time.Duration) error
}
