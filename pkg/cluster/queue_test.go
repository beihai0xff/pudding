package cluster

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	v3 "go.etcd.io/etcd/client/v3"
)

func Test_queue_Produce_Unique_key(t *testing.T) {
	q, _ := newQueue(testCluster, "/test/Produce")

	uMsg := &Message{
		Key:    "1123",
		Value:  "value",
		Unique: true,
	}

	defer q.client.Delete(context.Background(), "/test/Produce", v3.WithPrefix())

	got, err := q.Produce(uMsg)
	assert.NoError(t, err)
	assert.Equal(t, got, uMsg.Key, "the messageID should Equal the key")
	got, err = q.Produce(uMsg)
	assert.ErrorIs(t, ErrKeyExists, err)
	assert.Equal(t, "", got)
}

func Test_queue_Produce_non_unique_key(t *testing.T) {
	q, _ := newQueue(testCluster, "/test/Produce")
	defer q.client.Delete(context.Background(), "/test/Produce", v3.WithPrefix())

	uMsg := &Message{
		Key:    "1123",
		Value:  "value",
		Unique: false,
	}

	for i := 0; i < 3; i++ {
		got, err := q.Produce(uMsg)
		assert.NoError(t, err)
		assert.Contains(t, got, uMsg.Key, "the messageID should Contains the key")
		fmt.Println(got)
	}
}

func Test_queue_Consume(t *testing.T) {
	q, _ := newQueue(testCluster, "/test/Produce")
	defer q.client.Delete(context.Background(), "/test/Produce", v3.WithPrefix())

	go func() {
		for i := 0; i < 4; i++ {
			uMsg := &Message{
				Key:    fmt.Sprintf("%d", i),
				Value:  "value",
				Unique: true,
			}
			got, err := q.Produce(uMsg)
			assert.NoError(t, err)
			assert.Equal(t, got, uMsg.Key)
			time.Sleep(time.Second)
		}

	}()

	for i := 0; i < 4; i++ {
		got, err := q.Consume(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("%d", i), got.Key)
		assert.NoError(t, q.Commit(got))
	}
}

func Test_queue_Consume_wait_key(t *testing.T) {
	q, _ := newQueue(testCluster, "/test/Produce")
	defer q.client.Delete(context.Background(), "/test/Produce", v3.WithPrefix())

	for i := 0; i < 3; i++ {
		uMsg := &Message{
			Key:    fmt.Sprintf("%d", i),
			Value:  "value",
			Unique: true,
		}
		messageID, err := q.Produce(uMsg)
		assert.NoError(t, err)
		assert.Equal(t, messageID, uMsg.Key)
	}

	got, err := q.Consume(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "0", got.Key)
	go func() {
		start := time.Now()
		got, err := q.Consume(context.Background())
		assert.Equal(t, true, time.Since(start) > time.Second)
		assert.NoError(t, err)
		assert.Equal(t, "1", got.Key)
		assert.NoError(t, q.Commit(got))
		time.Sleep(time.Second)
	}()
	time.Sleep(time.Second)
	assert.NoError(t, q.Commit(got))
	time.Sleep(time.Second)
	got, err = q.Consume(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, "2", got.Key)
	assert.NoError(t, q.Commit(got))
}

func Test_queue_Commit(t *testing.T) {
	q, _ := newQueue(testCluster, "/test/Produce")
	defer q.client.Delete(context.Background(), "/test/Produce", v3.WithPrefix())

	uMsg := &Message{
		Key:    "key",
		Value:  "value",
		Unique: true,
	}
	messageID, err := q.Produce(uMsg)
	assert.NoError(t, err)
	assert.Equal(t, messageID, uMsg.Key)

	got, err := q.Consume(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, uMsg.Key, got.Key)
	assert.NoError(t, q.Commit(got))
	assert.NoError(t, q.Commit(got))

	messageID, err = q.Produce(uMsg)
	assert.NoError(t, err)
	assert.Equal(t, messageID, uMsg.Key)
	got, err = q.Consume(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, uMsg.Key, got.Key)
	assert.NoError(t, q.Commit(got))
}
