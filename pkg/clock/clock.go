// Package clock provides a clock interface and a real clock implementation.
package clock

import (
	"sync"
	"time"
)

// Clock is an interface for getting the current time.
type Clock interface {
	Now() time.Time
}

var (
	realClockOnce sync.Once
	realClock     *RealClock
)

// RealClock is a clock that returns the current time.
type RealClock struct{}

// New returns a new clock that returns the current time.
func New() Clock {
	realClockOnce.Do(func() {
		realClock = &RealClock{}
	})

	return realClock
}

// Now returns the current time.
func (c *RealClock) Now() time.Time {
	return time.Now().UTC()
}

// FakeClock is a clock that returns a fake time.
type FakeClock struct {
	t time.Time
}

// NewFakeClock returns a new fake clock.
func NewFakeClock(t time.Time) Clock {
	return &FakeClock{t: t.UTC()}
}

// Now returns the fake time.
func (c *FakeClock) Now() time.Time {
	return c.t
}
