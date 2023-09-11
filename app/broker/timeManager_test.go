package broker

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSchedule_formatTokenName(t *testing.T) {
	assert.Equal(t, "pudding_token:1", s.timeManager.formatTokenName(1))
	assert.Equal(t, "pudding_token:50", s.timeManager.formatTokenName(50))
	assert.Equal(t, "pudding_token:100", s.timeManager.formatTokenName(100))
	assert.Equal(t, "pudding_token:10000000", s.timeManager.formatTokenName(10000000))
}

func TestSchedule_formatTokenLockerName(t *testing.T) {
	assert.Equal(t, "pudding_locker_token:1", s.timeManager.formatTokenLockerName(1))
	assert.Equal(t, "pudding_locker_token:50", s.timeManager.formatTokenLockerName(50))
	assert.Equal(t, "pudding_locker_token:100", s.timeManager.formatTokenLockerName(100))
	assert.Equal(t, "pudding_locker_token:10000000", s.timeManager.formatTokenLockerName(10000000))
}

func TestSchedule_parseNowFromToken(t *testing.T) {
	assert.Equal(t, uint64(10000000), s.timeManager.parseTimeFromToken("pudding_token:10000000"))
	assert.Equal(t, uint64(100), s.timeManager.parseTimeFromToken("pudding_token:100"))
	assert.Equal(t, uint64(50), s.timeManager.parseTimeFromToken("pudding_token:50"))
	assert.Equal(t, uint64(0), s.timeManager.parseTimeFromToken("pudding_token:-2"))
	assert.Equal(t, uint64(0), s.timeManager.parseTimeFromToken("pudding_token:wewq"))
	assert.Equal(t, uint64(0), s.timeManager.parseTimeFromToken("broken_token:100"))
}

func Test_timeManager_sendToken(t *testing.T) {
	quit := make(chan struct{})
	timeManager, err := newTimeManager("pudding_token", s.cluster, quit)
	assert.NoError(t, err)

	start := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	assert.NoError(t, timeManager.sendToken(&start))

	timeManager.consumeToken(func(time uint64) error {
		assert.Equal(t, uint64(start.Unix()), time)
		close(quit)
		return nil
	})

}
