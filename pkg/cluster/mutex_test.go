package cluster

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_mutex_Lock(t *testing.T) {
	mutex, err := testCluster.Mutex("test", 2*time.Second)
	assert.NoError(t, err)

	ctx := context.Background()
	assert.NoError(t, mutex.Lock(ctx))
	assert.NoError(t, mutex.Unlock(ctx))

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	assert.NoError(t, mutex.Lock(ctx))
	assert.ErrorIs(t, ErrLockedBySelf, mutex.Lock(ctx))
	assert.NoError(t, mutex.Unlock(ctx))
}

func Test_mutexUnlock(t *testing.T) {
	mutex, err := testCluster.Mutex("test/unlock", 2*time.Second)
	assert.NoError(t, err)

	ctx := context.Background()
	assert.NoError(t, mutex.Lock(ctx))
	assert.NoError(t, mutex.Unlock(ctx))
	// unlock again
	assert.ErrorIs(t, mutex.Unlock(ctx), ErrLockNotHeld)
	assert.ErrorIs(t, mutex.Unlock(ctx), ErrLockNotHeld)
}

func Test_mutex_Refresh(t *testing.T) {
	mutex, err := testCluster.Mutex("test", 3*time.Second, WithDisableKeepalive())
	assert.NoError(t, err)
	ctx := context.Background()
	assert.NoError(t, mutex.Lock(ctx))
	refreshMutexAfter(t, mutex, time.Second)
	refreshMutexAfter(t, mutex, time.Second)
	refreshMutexAfter(t, mutex, 2*time.Second)

	// mutex2 try to lock the same key
	ctx2, cancel2 := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel2()
	mutex2, err := testCluster.Mutex("test", 5*time.Second, WithDisableKeepalive())
	assert.Error(t, mutex2.Lock(ctx2))

	ctx3, cancel3 := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel3()
	assert.NoError(t, mutex2.Lock(ctx3))
}

func refreshMutexAfter(t *testing.T, m Mutex, d time.Duration) {
	ctx := context.Background()
	time.Sleep(d)
	assert.NoError(t, m.Refresh(ctx))
	mm := m.(*mutex)
	resp, err := mm.session.Client().TimeToLive(ctx, mm.session.Lease())
	if err != nil {
		t.Errorf("refresh lock failed: %v", err)
	} else {
		t.Logf("refresh lock ttl: %v", resp.TTL)
	}
}

func Test_mutex_IsLocked(t *testing.T) {
	mutex1, err := testCluster.Mutex("test", 2*time.Second)
	assert.NoError(t, err)
	mutex2, err := testCluster.Mutex("test", 2*time.Second)
	assert.NoError(t, err)

	ctx := context.Background()
	assert.NoError(t, mutex1.Lock(ctx))
	assert.Equal(t, true, mutex1.IsLocked())
	assert.Equal(t, false, mutex2.IsLocked())

	assert.NoError(t, mutex1.Unlock(ctx))
	assert.Equal(t, false, mutex1.IsLocked())

	assert.NoError(t, mutex2.Lock(ctx))
	assert.Equal(t, true, mutex2.IsLocked())

	mutex3 := mutex2.(*mutex)
	assert.NoError(t, mutex3.session.Close())
	assert.Equal(t, false, mutex2.IsLocked())
}
