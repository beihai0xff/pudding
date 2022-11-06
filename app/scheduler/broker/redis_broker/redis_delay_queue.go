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

const (
	zsetNameFormat      = "zset_timeSlice_%s_bucket_%d"
	hashtableNameFormat = "hashTable_timeSlice_%s_bucket_%d"
)

type DelayQueue struct {
	rdb *rdb.Client // Redis Client
	// key is timeSlice name, value is the bucket nums in the partition
	bucket map[string]int8
}

func NewDelayQueue(rdb *rdb.Client) *DelayQueue {
	return &DelayQueue{
		rdb: rdb,
	}
}

func (q *DelayQueue) Produce(ctx context.Context, timeSlice string, msg *types.Message) error {
	// member := &redis.Z{Score: float64(msg.ReadyTime.Unix()), Member: msg.Key}
	log.Debugf("produce message: %+v", msg)
	return q.pushToZSet(ctx, timeSlice, msg)
}

func (q *DelayQueue) pushToZSet(ctx context.Context, timeSlice string, msg *types.Message) error {
	/*	err := q.rdb.ZAddNX(ctx, q.getZSetName(timeSlice), *member)
		if err != nil {
			return fmt.Errorf("pushToZSet failed: %w", err)
		}
	*/
	c, err := msgpack.Encode(msg)
	if err != nil {
		return fmt.Errorf("pushToZSet: failed to marshal message:%w", err)
	}

	count, err := pushScript.Run(ctx, q.rdb.GetClient(), []string{q.getZSetName(timeSlice),
		q.getHashtableName(timeSlice)}, msg.Key, c, msg.ReadyTime).Int()
	if err != nil {
		return fmt.Errorf("pushToZSet: failed to push message:%w", err)
	}
	if count == 0 {
		return errno.ErrDuplicateMessage
	}
	return nil
}

func (q *DelayQueue) Consume(ctx context.Context, timeSlice string, now, batchSize int64,
	fn types.HandleMessage) error {

	for {
		// batch get messages which are ready to execute
		messages, err := q.getFromZSetByScore(timeSlice, now, batchSize)
		// if you get error return directly
		if err != nil {
			return err
		}

		// if no data in the timeSlice, break the loop
		if len(messages) == 0 {
			break
		}

		// iterate all messages
		for _, msg := range messages {

			// 处理消息
			err = fn(ctx, &msg)
			if err != nil {
				log.Errorf("failed to handle message: %+v, caused by: %w", msg, err)
				continue
			}

			// delete message from zset and hash table
			if err := deleteScript.Run(ctx, q.rdb.GetClient(), []string{q.getZSetName(timeSlice),
				q.getHashtableName(timeSlice)}, msg.Key).Err(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (q *DelayQueue) getFromZSetByScore(timeSlice string, now, batchSize int64) ([]types.Message, error) {
	// 批量获取已经准备好执行的消息
	zs, err := q.rdb.ZRangeByScore(context.Background(), q.getZSetName(timeSlice), &redis.ZRangeBy{
		Min:    strconv.FormatInt(now, 10),
		Max:    strconv.FormatInt(now, 10),
		Offset: 0,
		Count:  batchSize,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get messages from zset: %w", err)
	}

	if len(zs) == 0 {
		return nil, nil
	}

	res := make([]types.Message, 0, len(zs))

	hashTable := q.getHashtableName(timeSlice)

	// 遍历每个 message key，根据 message key 获取 message body
	for _, z := range zs {
		key := z.Member.(string)
		// 获取消息的 body
		body, err := q.rdb.HGet(context.Background(), hashTable, key)
		if err != nil {
			log.Errorf("failed to get message body from hashTable: %w", err)
			continue
		}
		msg := &types.Message{}
		if err := msgpack.Decode(body, msg); err != nil {
			log.Errorf("failed to decode message body: %w", err)
			continue
		}
		log.Debugf("get message from zset: %+v", msg)
		res = append(res, *msg)
	}
	return res, nil
}

// Close the queue
func (q *DelayQueue) Close() error {
	return nil
}

func (q *DelayQueue) getZSetName(timeSlice string) string {
	return fmt.Sprintf(zsetNameFormat, timeSlice, q.getBucket(timeSlice))
}

func (q *DelayQueue) getHashtableName(timeSlice string) string {
	return fmt.Sprintf(hashtableNameFormat, timeSlice, q.getBucket(timeSlice))
}

func (q *DelayQueue) getBucket(timeSlice string) int8 {

	return 1
}
