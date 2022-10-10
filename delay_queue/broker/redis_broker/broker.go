package redis_broker

import (
	"fmt"

	rdb "github.com/beihai0xff/pudding/pkg/redis"
)

type RedisDelayQueue struct {
	rdb *rdb.Client // Redis客户端
}

// Close the queue
func (q *RedisDelayQueue) Close() error {
	return q.rdb.Close()
}

func (q *RedisDelayQueue) topicZSet(topic string, partition int) string {
	return fmt.Sprintf("zset_%s:%d", topic, partition)
}

func (q *RedisDelayQueue) topicHashtable(topic string, partition int) string {
	return fmt.Sprintf("hashTable_%s:%d", topic, partition)
}
