package redis_broker

import (
	"context"
	"os"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/beihai0xff/pudding/pkg/configs"
	rdb "github.com/beihai0xff/pudding/pkg/redis"
	"github.com/beihai0xff/pudding/types"
)

var q *DelayQueue

func TestMain(m *testing.M) {
	// initial Redis DB
	s, _ := miniredis.Run()

	q = &DelayQueue{
		rdb: rdb.New(&configs.RedisConfig{
			RedisURL: "redis://" + s.Addr(),
		}),
		bucket: map[string]int8{},
	}

	exitCode := m.Run()
	// 退出
	os.Exit(exitCode)
}

func TestRealTimeQueue_Produce(t *testing.T) {

	msg := &types.Message{
		Topic:     "test_Topic",
		Partition: 0,
		Payload:   []byte("12345678900987654321"),
		Delay:     0,
		ReadyTime: 10,
	}

	for i := 1; i <= 100; i++ {
		msg.Key = uuid.New().String()
		msg.ReadyTime = int64(i)
		assert.Equal(t, nil, q.Produce(context.Background(), "test_bucket_Produce", msg))
	}

	q.Consume(context.Background(), "test_bucket", 11, 10, func(ctx context.Context, msg *types.Message) error {
		assert.Equal(t, []byte("12345678900987654321"), msg.Payload)
		return nil
	})

}

func TestDelayQueue_getFromZSetByScore(t *testing.T) {
	bucket := "test_bucket_getFromZSetByScore"
	msg := &types.Message{
		Topic:     "test_Topic",
		Partition: 0,
		Payload:   []byte("12345678900987654321"),
		Delay:     0,
		ReadyTime: 10,
	}

	for i := 1; i <= 100; i++ {
		msg.Key = uuid.New().String()
		msg.ReadyTime = int64(i)
		assert.Equal(t, nil, q.Produce(context.Background(), bucket, msg))
	}

	for i := 1; i <= 5; i++ {
		msg.Key = uuid.New().String()
		msg.ReadyTime = 200
		assert.Equal(t, nil, q.Produce(context.Background(), bucket, msg))
	}

	if msgs, err := q.getFromZSetByScore(bucket, 10, 10); err != nil {
		t.Errorf("getFromZSetByScore error: %v", err)
	} else if len(msgs) != 1 {
		t.Errorf("getFromZSetByScore length is: %d", len(msgs))
	}

	msgs, _ := q.getFromZSetByScore(bucket, 10, 11)
	assert.Equal(t, 1, len(msgs))
	msgs, _ = q.getFromZSetByScore(bucket, 10, 12)
	assert.Equal(t, 1, len(msgs))
	msgs, _ = q.getFromZSetByScore(bucket, 200, 10)
	assert.Equal(t, 5, len(msgs))
	msgs, _ = q.getFromZSetByScore(bucket, 200, 3)
	assert.Equal(t, 3, len(msgs))
	msgs, _ = q.getFromZSetByScore(bucket, 200, 200)
	assert.Equal(t, 5, len(msgs))
}

func TestRealTimeQueue_getZSetName(t *testing.T) {

	assert.Equal(t, "zset_timeSlice_10_bucket_1", q.getZSetName("10"))
	assert.Equal(t, "zset_timeSlice_11_bucket_1", q.getZSetName("11"))
	assert.Equal(t, "zset_timeSlice_12_bucket_1", q.getZSetName("12"))
	assert.Equal(t, "zset_timeSlice_13_bucket_1", q.getZSetName("13"))
}

func TestRealTimeQueue_getHashtableName(t *testing.T) {

	assert.Equal(t, "hashTable_timeSlice_10_bucket_1", q.getHashtableName("10"))
	assert.Equal(t, "hashTable_timeSlice_11_bucket_1", q.getHashtableName("11"))
	assert.Equal(t, "hashTable_timeSlice_12_bucket_1", q.getHashtableName("12"))
	assert.Equal(t, "hashTable_timeSlice_13_bucket_1", q.getHashtableName("13"))
}

func TestRealTimeQueue_getBucket(t *testing.T) {

	assert.Equal(t, int8(1), q.getBucket("10"))
	assert.Equal(t, int8(1), q.getBucket("10"))
	assert.Equal(t, int8(1), q.getBucket("10"))
	assert.Equal(t, int8(1), q.getBucket("10"))
	assert.Equal(t, int8(1), q.getBucket("10"))
	assert.Equal(t, int8(1), q.getBucket("10"))
	assert.Equal(t, int8(1), q.getBucket("10"))
	assert.Equal(t, int8(1), q.getBucket("10"))
}
