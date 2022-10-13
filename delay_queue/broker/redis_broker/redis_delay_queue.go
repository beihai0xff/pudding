package redis_broker

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v9"

	"github.com/beihai0xff/pudding/pkg/configs"
	"github.com/beihai0xff/pudding/pkg/log"
	rdb "github.com/beihai0xff/pudding/pkg/redis"
	"github.com/beihai0xff/pudding/types"
)

type DelayQueue struct {
	rdb *rdb.Client // Redis Client
	// key is partition, value is the bucket ID in the partition
	bucket map[int64]int8
}

func NewDelayQueue() *DelayQueue {
	return &DelayQueue{
		rdb: rdb.NewRDB(configs.GetRedisConfig()),
	}
}

func (q *DelayQueue) Produce(ctx context.Context, partition int64, msg *types.Message) error {
	member := &redis.Z{Score: float64(msg.ReadyTime.Unix()), Member: msg.Key}
	return q.pushToZSet(ctx, partition, member)
}

func (q *DelayQueue) pushToZSet(ctx context.Context, partition int64, member *redis.Z) error {
	err := q.rdb.ZAddNX(ctx, q.getZSet(partition), *member)
	if err != nil {
		return fmt.Errorf("pushToZSet failed: %w", err)
	}

	return nil
}

func (q *DelayQueue) Consume(ctx context.Context, partition, batchSize int64,
	fn func(msg *types.Message) error) error {

	// 批量获取已经准备好执行的消息
	messages, err := q.getFromZSetByScore(partition, batchSize)
	// 如果获取出错或者获取不到消息，则直接返回
	if err != nil || len(messages) == 0 {
		return err
	}

	zset := q.getZSet(partition)
	// 遍历每个消息
	for _, msg := range messages {

		// 处理消息
		err = fn(&msg)
		if err != nil {
			log.Errorf("failed to handle message: %+v, caused by: %v", msg, err)
			continue
		}

		// 如果消息处理成功，删除消息
		if err := q.rdb.ZRem(ctx, zset, msg.Key); err != nil {
			return err
		}
	}

	return nil
}

func (q *DelayQueue) getFromZSetByScore(partition, batchSize int64) ([]types.Message, error) {
	// 批量获取已经准备好执行的消息
	zs, err := q.rdb.ZRangeByScore(context.Background(), q.getZSet(partition), &redis.ZRangeBy{
		Min:    "-inf",
		Max:    strconv.FormatInt(time.Now().Unix(), 10),
		Offset: 0,
		Count:  batchSize,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get messages from zset: %w", err)
	}

	if zs == nil || len(zs) == 0 {
		return nil, nil
	}

	res := make([]types.Message, len(zs))

	hashTable := q.getHashtable(partition)

	// 遍历每个 message key，根据 message key 获取 message body
	for _, z := range zs {
		key := z.Member.(string)
		// 获取消息的 body
		body, err := q.rdb.HGet(context.Background(), hashTable, key)
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

// Close the queue
func (q *DelayQueue) Close() error {
	return nil
}

func (q *DelayQueue) getZSet(partition int64) string {

	return fmt.Sprintf("zset_partition_%d_bucket_%8d", partition, q.getBucket(partition))
}

func (q *DelayQueue) getBucket(partition int64) int8 {
	buckets := q.bucket[partition]
	if buckets <= 0 {
		q.bucket[partition] = 1
		return 1
	}

	if buckets < 3 {
		return buckets
	}

	return 1
}

func (q *DelayQueue) getHashtable(partition int64) string {
	return fmt.Sprintf("hashTable_partition_%d", partition)
}
