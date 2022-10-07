package redis_broker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/go-redis/redis/v9"

	"github.com/beihai0xff/pudding/pkg/errno"
	rdb "github.com/beihai0xff/pudding/pkg/redis"

	"github.com/beihai0xff/pudding/types"
)

type RedisDelayQueue struct {
	rdb *rdb.Client // Redis客户端
}

func (q *RedisDelayQueue) Produce(ctx context.Context, msg *types.Message) error {
	// 如果设置了 ReadyTime，则使用 ReadyTime
	var readyTime int64
	if !msg.ReadyTime.IsZero() {
		readyTime = msg.ReadyTime.Unix()
	} else {
		// 否则使用Delay
		readyTime = time.Now().Add(msg.Delay).Unix()
	}
	return q.pushToZSet(ctx, readyTime, msg)
}

func (q *RedisDelayQueue) NewConsumer(topic string, batchSize int, fn func(msg *types.Message) error) {
	for {
		// 批量获取已经准备好执行的消息
		messages, err := q.getFromZSetByScore(topic, batchSize)
		// 如果获取出错或者获取不到消息，则休眠一秒
		if err != nil || len(messages) == 0 {
			time.Sleep(time.Second)
			continue
		}
		// 遍历每个消息
		for _, msg := range messages {

			// 处理消息
			err = fn(&msg)
			if err != nil {
				log.Printf("failed to handle message: %+v, caused by: %v", msg, err)
				continue
			}
			// 如果消息处理成功，删除消息
			deleteScript.Run(context.Background(), q.rdb.GetClient(),
				[]string{q.topicZSet(topic), q.topicHashtable(topic)}, msg.Key)
		}
	}
}

func (q *RedisDelayQueue) topicZSet(topic string) string {
	return "zset_" + topic
}

func (q *RedisDelayQueue) topicHashtable(topic string) string {
	return "hashTable_" + topic
}

func (q *RedisDelayQueue) pushToZSet(ctx context.Context, readyTime int64, msg *types.Message) error {
	c, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("pushToZSet: failed to marshal message:%w", err)
	}
	success, err := pushScript.Run(ctx, q.rdb.GetClient(), []string{q.topicZSet(msg.Topic), q.topicHashtable(msg.Topic)},
		msg.Key, c, readyTime).Bool()
	if err != nil {
		return fmt.Errorf("pushToZSet: failed to push message:%w", err)
	}
	if !success {
		return errno.ErrDuplicateMessage
	}
	return nil
}

func (q *RedisDelayQueue) getFromZSetByScore(topic string, batchSize int) ([]types.Message, error) {
	// 批量获取已经准备好执行的消息
	zs, err := q.rdb.GetClient().ZRangeByScoreWithScores(context.Background(), q.topicZSet(topic), &redis.ZRangeBy{
		Min:   "-inf",
		Max:   strconv.FormatInt(time.Now().Unix(), 10),
		Count: int64(batchSize),
	}).Result()

	if err != nil {
		return nil, fmt.Errorf("failed to get messages from zset: %w", err)
	}

	if zs == nil || len(zs) == 0 {
		return nil, nil
	}

	res := make([]types.Message, len(zs))

	// 遍历每个 message key，根据 message key 获取 message body
	for _, z := range zs {
		key := z.Member.(string)
		// 获取消息的body
		body, err := q.rdb.GetClient().HGet(context.Background(), q.topicHashtable(topic), key).Bytes()
		if err != nil {
			continue
		}
		res = append(res, types.Message{
			Topic:     topic,
			Key:       key,
			Body:      body,
			ReadyTime: time.Unix(int64(z.Score), 0),
		})
	}
	return res, nil
}
