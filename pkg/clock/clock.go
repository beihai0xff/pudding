package clock

import "time"

type Clock interface {
	Now() time.Time
}

type RealClock struct{}

func New() Clock {
	return &RealClock{}
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
