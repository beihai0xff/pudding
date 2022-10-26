package scheduler

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/beihai0xff/pudding/pkg/mq/pulsar"
	rdb "github.com/beihai0xff/pudding/pkg/redis"
)

var s *Scheduler

func TestMain(m *testing.M) {

	s = &Scheduler{
		delay:    NewDelayQueue(rdb.NewMockRdb()),
		realtime: NewRealTimeQueue(pulsar.NewMockPulsar()),
		interval: 60,
	}

	exitCode := m.Run()
	// Exit
	os.Exit(exitCode)
}

func TestScheduler_getTimeSlice(t *testing.T) {
	type args struct {
		readyTime int64
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"test0", args{readyTime: 0}, "0~60"},
		{"test1", args{readyTime: 1}, "0~60"},
		{"test2", args{readyTime: 2}, "0~60"},
		{"test59", args{readyTime: 59}, "0~60"},
		{"test60", args{readyTime: 60}, "60~120"},
		{"test61", args{readyTime: 61}, "60~120"},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, s.getTimeSlice(tt.args.readyTime))
	}
}
