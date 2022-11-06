package scheduler

import (
	"errors"
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	"github.com/beihai0xff/pudding/internal/scheduler/broker"
	rdb "github.com/beihai0xff/pudding/pkg/redis"
	"github.com/beihai0xff/pudding/test/mock"
	"github.com/beihai0xff/pudding/types"
)

var s *Schedule

func TestMain(m *testing.M) {

	s = &Schedule{
		delay:    broker.NewDelayQueue(rdb.NewMockRdb()),
		interval: 60,
	}

	exitCode := m.Run()
	// Exit
	os.Exit(exitCode)
}

func beforeEach(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	realtime := mock.NewMockRealTimeQueue(mockCtrl)
	realtime.EXPECT().Produce(gomock.Any(), gomock.Any()).Return(nil).Times(1)
	realtime.EXPECT().Produce(gomock.Any(), &types.Message{
		Topic:     "test_topic",
		Key:       "produce failed three times",
		Payload:   []byte("test_payload"),
		Delay:     10,
		ReadyTime: 0,
	}).Return(errors.New("broken connection")).Times(3)

	s.realtime = realtime
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

func TestSchedule_checkParams(t *testing.T) {
	var err error
	// no delay no ready time
	err = s.checkParams(&types.Message{Key: "12345"})
	assert.Error(t, err)
	if err != nil {
		assert.EqualError(t, err, errInvalidMessageDelay.Error())
	}

	// test ReadyTime less than now
	err = nil
	err = s.checkParams(&types.Message{Key: "12345", Delay: 50, ReadyTime: 60})
	assert.Error(t, err)
	if err != nil {
		assert.EqualError(t, err, errInvalidMessageReady.Error())
	}

	// test ReadyTime greater than now
	err = nil
	msg := &types.Message{Delay: 50, ReadyTime: 60000000000}
	err = s.checkParams(msg)
	assert.NoError(t, err)
	// test ReadyTime
	assert.Equal(t, int64(60000000000), msg.ReadyTime)
	// test no topic set
	assert.Equalf(t, types.DefaultTopic, msg.Topic, "msg.topic: %s", msg.Topic)
	// test no uuid set
	assert.NotEqualf(t, "", msg.Key, "msg.key: %s", msg.Key)

	msg = &types.Message{Topic: "test_topic", Key: "test_key", Delay: 50, ReadyTime: 60000000000}
	err = s.checkParams(msg)
	assert.NoError(t, err)
	// test ReadyTime
	assert.Equal(t, int64(60000000000), msg.ReadyTime)
	// test no topic set
	assert.Equalf(t, "test_topic", msg.Topic, "msg.topic: %s", msg.Topic)
	// test no uuid set
	assert.Equalf(t, "test_key", msg.Key, "msg.Key: %s", msg.Key)
}

func TestSchedule_Produce(t *testing.T) {
	// beforeEach(t)

	ctx := context.Background()
	var err error
	// no delay no ready time
	err = s.Produce(ctx, &types.Message{Key: "12345"})
	assert.Error(t, err)
	if err != nil {
		assert.EqualError(t, err, "check message params failed: "+errInvalidMessageDelay.Error())
	}

	// test ReadyTime less than now
	err = nil
	err = s.Produce(ctx, &types.Message{Key: "12345", Delay: 50, ReadyTime: 60})
	assert.Error(t, err)
	if err != nil {
		assert.EqualError(t, err, "check message params failed: "+errInvalidMessageReady.Error())
	}

	// test ReadyTime greater than now
	err = nil
	msg := &types.Message{Delay: 50, ReadyTime: 60000000000}
	err = s.Produce(ctx, msg)
	assert.NoError(t, err)
	// test ReadyTime
	assert.Equal(t, int64(60000000000), msg.ReadyTime)
	// test no topic set
	assert.Equalf(t, types.DefaultTopic, msg.Topic, "msg.topic: %s", msg.Topic)
	// test no uuid set
	assert.NotEqualf(t, "", msg.Key, "msg.key: %s", msg.Key)

	msg = &types.Message{Topic: "test_topic", Key: "test_key", Delay: 50, ReadyTime: 60000000000}
	err = s.Produce(ctx, msg)
	assert.NoError(t, err)
	// test ReadyTime
	assert.Equal(t, int64(60000000000), msg.ReadyTime)
	// test no topic set
	assert.Equalf(t, "test_topic", msg.Topic, "msg.topic: %s", msg.Topic)
	// test no uuid set
	assert.Equalf(t, "test_key", msg.Key, "msg.Key: %s", msg.Key)
}

func TestSchedule_getLockerName(t *testing.T) {
	assert.Equal(t, "pudding_locker_time:5", s.getLockerName(5))
	assert.Equal(t, "pudding_locker_time:60", s.getLockerName(60))
	assert.Equal(t, "pudding_locker_time:60000000000", s.getLockerName(60000000000))
	assert.Equal(t, "pudding_locker_time:1321131", s.getLockerName(1321131))
}
