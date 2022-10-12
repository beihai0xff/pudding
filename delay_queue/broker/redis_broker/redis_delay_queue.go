package redis_broker

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/beihai0xff/pudding/pkg/errno"
	"github.com/beihai0xff/pudding/pkg/log"
	rdb "github.com/beihai0xff/pudding/pkg/redis"
	"github.com/beihai0xff/pudding/types"

	"github.com/go-redis/redis/v9"
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
		// 否则使用 Delay
		readyTime = time.Now().Add(msg.Delay).Unix()
	}
	return q.pushToZSet(ctx, readyTime, msg)
}

func (q *RedisDelayQueue) pushToZSet(ctx context.Context, readyTime int64, msg *types.Message) error {
	c, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal message failed: %w", err)
	}

	success, err := pushScript.Run(ctx, q.rdb.GetClient(), []string{q.topicZSet(msg.Topic, msg.Partition),
		q.topicHashtable(msg.Topic, msg.Partition)}, msg.Key, c, readyTime).Bool()
	if err != nil {
		return fmt.Errorf("pushToZSet failed: %w", err)
	}
	if !success {
		return errno.ErrDuplicateMessage
	}
	return nil
}

func (q *RedisDelayQueue) NewConsumer(topic string, partition, batchSize int, fn func(msg *types.Message) error) {
	for {
		// 批量获取已经准备好执行的消息
		messages, err := q.getFromZSetByScore(topic, batchSize, partition)
		// 如果获取出错或者获取不到消息，则休眠一秒
		if err != nil || len(messages) == 0 {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		// 遍历每个消息
		for _, msg := range messages {

			// 处理消息
			err = fn(&msg)
			if err != nil {
				log.Errorf("failed to handle message: %+v, caused by: %v", msg, err)
				continue
			}
			// 如果消息处理成功，删除消息
			deleteScript.Run(context.Background(), q.rdb.GetClient(),
				[]string{q.topicZSet(topic, partition), q.topicHashtable(topic, partition)}, msg.Key)
		}
	}
}

func (q *RedisDelayQueue) getFromZSetByScore(topic string, batchSize, partition int) ([]types.Message, error) {
	// 批量获取已经准备好执行的消息
	zs, err := q.rdb.ZRangeByScore(context.Background(), q.topicZSet(topic, partition), &redis.ZRangeBy{
		Min:    "-inf",
		Max:    strconv.FormatInt(time.Now().Unix(), 10),
		Offset: 0,
		Count:  int64(batchSize),
	})

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
		// 获取消息的 body
		body, err := q.rdb.HGet(context.Background(), q.topicHashtable(topic, partition), key)
		if err != nil {
			// TODO: 记录错误日志
			continue
		}
		msg, err := types.GetMessageFromJSON(body)
		if err != nil {
			// TODO: 记录错误日志
			continue
		}
		res = append(res, *msg)
	}
	return res, nil
}

// Close Close the queue
func (q *RedisDelayQueue) Close() error {
	return q.rdb.Close()
}

func (q *RedisDelayQueue) topicZSet(topic string, partition int) string {
	return fmt.Sprintf("zset_%s:%d", topic, partition)
}

func (q *RedisDelayQueue) topicHashtable(topic string, partition int) string {
	return fmt.Sprintf("hashTable_%s:%d", topic, partition)
}
