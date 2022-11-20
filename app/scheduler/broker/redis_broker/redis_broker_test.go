package redis_broker

import (
	"context"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	rdb "github.com/beihai0xff/pudding/pkg/redis"
	"github.com/beihai0xff/pudding/types"
)

var q *DelayQueue

func TestMain(m *testing.M) {

	q = &DelayQueue{
		rdb:      rdb.NewMockRdb(),
		interval: 60,
		bucket:   map[string]int8{},
	}

	exitCode := m.Run()
	// 退出
	os.Exit(exitCode)
}

func TestRealTimeQueue_Produce(t *testing.T) {

	msg := &types.Message{
		Topic:        "test_Topic_Produce",
		Payload:      []byte("12345678900987654321"),
		DeliverAfter: 0,
		DeliverAt:    1,
	}

	for i := 1; i <= 100; i++ {
		msg.Key = uuid.New().String()
		msg.DeliverAt = int64(i)
		assert.Equal(t, nil, q.Produce(context.Background(), msg))
	}

	q.Consume(context.Background(), 2, 10, func(ctx context.Context, msg *types.Message) error {
		assert.Equal(t, []byte("12345678900987654321"), msg.Payload)
		return nil
	})

}

func TestDelayQueue_getTimeSlice(t *testing.T) {
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
		assert.Equal(t, tt.want, q.getTimeSlice(tt.args.readyTime))
	}
}

func TestDelayQueue_getFromZSetByScore(t *testing.T) {
	msg := &types.Message{
		Topic:   "test_Topic_getFromZSetByScore",
		Payload: []byte("12345678900987654321"),
	}

	for i := 120; i <= 180; i++ {
		testm := *msg
		testm.Key = uuid.New().String()
		testm.DeliverAt = int64(i)
		assert.Equal(t, nil, q.Produce(context.Background(), &testm))
	}

	for i := 1; i <= 5; i++ {
		testm := *msg
		testm.Key = uuid.New().String()
		testm.DeliverAt = 300
		assert.Equal(t, nil, q.Produce(context.Background(), &testm))
	}

	if msgs, err := q.getFromZSetByScore(q.getTimeSlice(10), 9, 10); err != nil {
		t.Errorf("getFromZSetByScore error: %v", err)
	} else if len(msgs) != 1 {
		t.Errorf("getFromZSetByScore length is: %d", len(msgs))
	}

	msgs, _ := q.getFromZSetByScore(q.getTimeSlice(60), 60, 11)
	assert.Equal(t, 1, len(msgs))
	msgs, _ = q.getFromZSetByScore(q.getTimeSlice(60), 60, 12)
	assert.Equal(t, 1, len(msgs))
	msgs, _ = q.getFromZSetByScore(q.getTimeSlice(300), 300, 3)
	assert.Equal(t, 3, len(msgs))
	msgs, _ = q.getFromZSetByScore(q.getTimeSlice(300), 300, 10)
	assert.Equal(t, 5, len(msgs))
	msgs, _ = q.getFromZSetByScore(q.getTimeSlice(300), 300, 200)
	assert.Equal(t, 5, len(msgs))
}

func TestDelayQueue_getZSetName(t *testing.T) {

	assert.Equal(t, "zset_timeSlice_60~120_bucket_1", q.getZSetName("60~120"))
	assert.Equal(t, "zset_timeSlice_120~180_bucket_1", q.getZSetName("120~180"))
	assert.Equal(t, "zset_timeSlice_12~13_bucket_1", q.getZSetName("12~13"))
	assert.Equal(t, "zset_timeSlice_13~15_bucket_1", q.getZSetName("13~15"))
}

func TestDelayQueue_getHashtableName(t *testing.T) {

	assert.Equal(t, "hashTable_timeSlice_10_bucket_1", q.getHashtableName("10"))
	assert.Equal(t, "hashTable_timeSlice_11_bucket_1", q.getHashtableName("11"))
	assert.Equal(t, "hashTable_timeSlice_12_bucket_1", q.getHashtableName("12"))
	assert.Equal(t, "hashTable_timeSlice_13_bucket_1", q.getHashtableName("13"))
}

func TestDelayQueue_getBucket(t *testing.T) {

	assert.Equal(t, int8(1), q.getBucket("10"))
	assert.Equal(t, int8(1), q.getBucket("10"))
	assert.Equal(t, int8(1), q.getBucket("10"))
	assert.Equal(t, int8(1), q.getBucket("10"))
	assert.Equal(t, int8(1), q.getBucket("10"))
	assert.Equal(t, int8(1), q.getBucket("10"))
	assert.Equal(t, int8(1), q.getBucket("10"))
	assert.Equal(t, int8(1), q.getBucket("10"))
}
