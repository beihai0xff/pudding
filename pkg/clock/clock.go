package clock

import (
	"sync"
	"time"
)

type Clock interface {
	Now() time.Time
}

var (
	realClockOnce sync.Once
	realClock     *RealClock
)

type RealClock struct{}

func New() Clock {
	realClockOnce.Do(func() {
		realClock = &RealClock{}
	})
	return realClock
}

func (c *RealClock) Now() time.Time {
	return time.Now().UTC()
}

type FakeClock struct {
	t time.Time
}

func NewFakeClock(t time.Time) Clock {
	return &FakeClock{t: t.UTC()}
}

func (c *FakeClock) Now() time.Time {
	return c.t
}
