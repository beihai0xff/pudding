package broker

import (
	"errors"
	"os"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	"github.com/beihai0xff/pudding/api/gen/pudding/types/v1"
	"github.com/beihai0xff/pudding/app/broker/storage/redis_storage"
	"github.com/beihai0xff/pudding/pkg/clock"
	rdb "github.com/beihai0xff/pudding/pkg/redis"
	mock "github.com/beihai0xff/pudding/test/mock/app/broker/connector"
	. "github.com/beihai0xff/pudding/types"
)

var s *scheduler

func TestMain(m *testing.M) {

	s = &scheduler{
		delay:        redis_storage.NewDelayQueue(rdb.NewMockRdb(), 60),
		messageTopic: DefaultTopic,
		tokenTopic:   TokenTopic,
		wallClock:    clock.NewFakeClock(time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)),
	}

	exitCode := m.Run()
	// Exit
	os.Exit(exitCode)
}

func beforeEach(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	realtime := mock.NewMockRealTimeConnector(mockCtrl)
	realtime.EXPECT().Produce(gomock.Any(), gomock.Any()).Return(nil).Times(1)
	realtime.EXPECT().Produce(gomock.Any(), &types.Message{
		Topic:        "test_topic",
		Key:          "produce failed three times",
		Payload:      []byte("test_payload"),
		DeliverAfter: 10,
		DeliverAt:    0,
	}).Return(errors.New("broken connection")).Times(3)

	s.connector = realtime
}

func TestSchedule_checkParams(t *testing.T) {
	var err error
	// no delay no ready time
	err = s.checkParams(&types.Message{Key: "12345"})
	assert.Error(t, err)
	if err != nil {
		assert.EqualError(t, err, errInvalidMessageDelay.Error())
	}

	// test DeliverAt less than now
	err = nil
	err = s.checkParams(&types.Message{Key: "12345", DeliverAfter: 50, DeliverAt: 60})
	assert.Error(t, err)
	if err != nil {
		assert.EqualError(t, err, errInvalidMessageReady.Error())
	}

	// test DeliverAt greater than now
	err = nil
	msg := &types.Message{DeliverAfter: 50, DeliverAt: 60000000000}
	err = s.checkParams(msg)
	assert.NoError(t, err)
	// test DeliverAt
	assert.Equal(t, int64(60000000000), msg.DeliverAt)
	// test no topic set
	assert.Equalf(t, DefaultTopic, msg.Topic, "msg.topic: %s", msg.Topic)
	// test no uuid set
	assert.NotEqualf(t, "", msg.Key, "msg.key: %s", msg.Key)

	msg = &types.Message{Topic: "test_topic", Key: "test_key", DeliverAfter: 50, DeliverAt: 60000000000}
	err = s.checkParams(msg)
	assert.NoError(t, err)
	// test DeliverAt
	assert.Equal(t, int64(60000000000), msg.DeliverAt)
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

	// test DeliverAt less than now
	err = nil
	err = s.Produce(ctx, &types.Message{Key: "12345", DeliverAfter: 50, DeliverAt: 60})
	assert.Error(t, err)
	if err != nil {
		assert.EqualError(t, err, "check message params failed: "+errInvalidMessageReady.Error())
	}

	// test DeliverAt greater than now
	err = nil
	msg := &types.Message{DeliverAfter: 50, DeliverAt: 60000000000}
	err = s.Produce(ctx, msg)
	assert.NoError(t, err)
	// test DeliverAt
	assert.Equal(t, int64(60000000000), msg.DeliverAt)
	// test no topic set
	assert.Equalf(t, DefaultTopic, msg.Topic, "msg.topic: %s", msg.Topic)
	// test no uuid set
	assert.NotEqualf(t, "", msg.Key, "msg.key: %s", msg.Key)

	msg = &types.Message{Topic: "test_topic", Key: "test_key", DeliverAfter: 50, DeliverAt: 60000000000}
	err = s.Produce(ctx, msg)
	assert.NoError(t, err)
	// test DeliverAt
	assert.Equal(t, int64(60000000000), msg.DeliverAt)
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
