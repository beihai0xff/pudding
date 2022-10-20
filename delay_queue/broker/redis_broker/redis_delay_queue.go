package redis_broker

import (
	"context"
	"fmt"
	"strconv"

	"github.com/go-redis/redis/v9"

	"github.com/beihai0xff/pudding/pkg/errno"
	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/pkg/msgpack"
	rdb "github.com/beihai0xff/pudding/pkg/redis"
	"github.com/beihai0xff/pudding/types"
)

type DelayQueue struct {
	rdb *rdb.Client // Redis Client
	// key is partition, value is the bucket nums in the partition
	bucket map[string]int8
}

func NewDelayQueue(rdb *rdb.Client) *DelayQueue {
	return &DelayQueue{
		rdb: rdb,
	}
}

func (q *DelayQueue) Produce(ctx context.Context, quantum string, msg *types.Message) error {
	// member := &redis.Z{Score: float64(msg.ReadyTime.Unix()), Member: msg.Key}
	return q.pushToZSet(ctx, quantum, msg)
}

func (q *DelayQueue) pushToZSet(ctx context.Context, quantum string, msg *types.Message) error {
	/*	err := q.rdb.ZAddNX(ctx, q.getZSetName(quantum), *member)
		if err != nil {
			return fmt.Errorf("pushToZSet failed: %w", err)
		}
	*/
	c, err := msgpack.Encode(msg)
	if err != nil {
		return fmt.Errorf("pushToZSet: failed to marshal message:%w", err)
	}

	count, err := pushScript.Run(ctx, q.rdb.GetClient(), []string{q.getZSetName(quantum),
		q.getHashtableName(quantum)}, msg.Key, c, msg.ReadyTime).Int()
	if err != nil {
		return fmt.Errorf("pushToZSet: failed to push message:%w", err)
	}
	if count == 0 {
		return errno.ErrDuplicateMessage
	}
	return nil
}

func (q *DelayQueue) Consume(ctx context.Context, quantum string, now, batchSize int64,
	fn types.HandleMessage) error {

	for {
		// batch get messages which are ready to execute
		messages, err := q.getFromZSetByScore(quantum, now, batchSize)
		// if you get error return directly
		if err != nil {
			return err
		}

		// if no data in the quantum, break the loop
		if messages == nil || len(messages) == 0 {
			break
		}

		// iterate all messages
		for _, msg := range messages {

			// 处理消息
			err = fn(ctx, &msg)
			if err != nil {
				log.Errorf("failed to handle message: %+v, caused by: %v", msg, err)
				continue
			}

			// delete message from zset and hash table
			if err := deleteScript.Run(ctx, q.rdb.GetClient(), []string{q.getZSetName(quantum),
				q.getHashtableName(quantum)}, msg.Key).Err(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (q *DelayQueue) getFromZSetByScore(quantum string, now, batchSize int64) ([]types.Message, error) {
	// 批量获取已经准备好执行的消息
	zs, err := q.rdb.ZRangeByScore(context.Background(), q.getZSetName(quantum), &redis.ZRangeBy{
		Min:    "-inf",
		Max:    strconv.FormatInt(now, 10),
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

	hashTable := q.getHashtableName(quantum)

	// 遍历每个 message key，根据 message key 获取 message body
	for _, z := range zs {
		key := z.Member.(string)
		// 获取消息的 body
		body, err := q.rdb.HGet(context.Background(), hashTable, key)
		if err != nil {
			log.Errorf("failed to get message body from hashTable: %v", err)
			continue
		}
		msg := &types.Message{}
		if err := msgpack.Decode(body, msg); err != nil {
			log.Errorf("failed to decode message body: %v", err)
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

func (q *DelayQueue) getZSetName(quantum string) string {
	return fmt.Sprintf("zset_quantum_%s_bucket_%8d", quantum, q.getBucket(quantum))
}

func (q *DelayQueue) getBucket(quantum string) int8 {
	buckets := q.bucket[quantum]
	if buckets <= 0 {
		q.bucket[quantum] = 1
		return 1
	}

	return 1
}

func (q *DelayQueue) getHashtableName(quantum string) string {
	return fmt.Sprintf("hashTable_quantum_%s_bucket_%8d", quantum, q.getBucket(quantum))
}
