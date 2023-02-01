package redis_storage

import (
	"context"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"

	"github.com/beihai0xff/pudding/api/gen/pudding/types/v1"
	type2 "github.com/beihai0xff/pudding/app/broker/pkg/types"
	"github.com/beihai0xff/pudding/pkg/errno"
	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/pkg/msgpack"
	rdb "github.com/beihai0xff/pudding/pkg/redis"
)

const (
	timeSliceNameFormat = "%d~%d"
	zsetNameFormat      = "zset_timeSlice_%s_bucket_%d"
	hashtableNameFormat = "hashTable_timeSlice_%s_bucket_%d"
)

// DelayStorage is a delay queue based on redis
type DelayStorage struct {
	rdb *rdb.Client // Redis Client
	// interval timeSlice interval (Seconds)
	interval int64
	// key is timeSlice name, value is the bucket nums in the partition
	bucket map[string]int8
}

// NewDelayStorage create a new DelayStorage
func NewDelayStorage(client *rdb.Client, interval int64) *DelayStorage {
	return &DelayStorage{
		rdb:      client,
		interval: interval,
	}
}

// Produce produce a Message to DelayStorage
func (q *DelayStorage) Produce(ctx context.Context, msg *types.Message) error {
	// member := &redis.Z{Score: float64(msg.DeliverAt.Unix()), Member: msg.Key}

	timeSlice := q.getTimeSlice(msg.DeliverAt)
	return q.pushToZSet(ctx, timeSlice, msg)
}

func (q *DelayStorage) pushToZSet(ctx context.Context, timeSlice string, msg *types.Message) error {
	/*	err := q.rdb.ZAddNX(ctx, q.getZSetName(timeSlice), *member)
		if err != nil {
			return fmt.Errorf("pushToZSet failed: %v", err)
		}
	*/
	c, err := msgpack.Encode(msg)
	if err != nil {
		return fmt.Errorf("pushToZSet: failed to marshal message: %w", err)
	}

	count, err := pushScript.Run(ctx, q.rdb.GetClient(), []string{q.getZSetName(timeSlice),
		q.getHashtableName(timeSlice)}, msg.Key, c, msg.DeliverAt).Int()
	if err != nil {
		return fmt.Errorf("pushToZSet: failed to push message: %w", err)
	}
	if count == 0 {
		return errno.ErrDuplicateMessage
	}
	return nil
}

// Consume consume Messages from the queue
func (q *DelayStorage) Consume(ctx context.Context, now, batchSize int64,
	fn type2.HandleMessage) error {
	timeSlice := q.getTimeSlice(now)

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
			// use pointer to avoid copylocks: range var msg copies lock alarm
			// 处理消息
			err = fn(ctx, msg)
			if err != nil {
				log.Errorf("failed to handle message: %s, caused by: %v", msg.String(), err)
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

func (q *DelayStorage) getFromZSetByScore(timeSlice string, now, batchSize int64) ([]*types.Message, error) {
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

	res := make([]*types.Message, 0, len(zs))

	hashTable := q.getHashtableName(timeSlice)

	// 遍历每个 message key，根据 message key 获取 message body
	for i := range zs {
		// use pointer to avoid copylocks: range var msg copies lock alarm
		z := &zs[i]
		key := z.Member.(string)
		// 获取消息的 body
		body, err := q.rdb.HGet(context.Background(), hashTable, key)
		if err != nil {
			log.Errorf("failed to get message body from hashTable: %v", err)
			continue
		}
		msg := types.Message{}
		if err := msgpack.Decode(body, &msg); err != nil {
			log.Errorf("failed to decode message body: %v", err)
			continue
		}
		log.Debugf("get message from zset: %s", msg.String())
		res = append(res, &msg)
	}
	return res, nil
}

// Close the queue
func (q *DelayStorage) Close() error {
	return nil
}

// getTimeSlice get the time slice of the given time
// Left closed right open interval
// e.g. the given interval is 60, the range is [0, 60)、[60, 120)、[120, 180)...
// 59 => 0~60
// 60 => 60~120
// 61 => 60~120
func (q *DelayStorage) getTimeSlice(readyTime int64) string {
	startAt := (readyTime / q.interval) * q.interval
	endAt := startAt + q.interval
	return fmt.Sprintf(timeSliceNameFormat, startAt, endAt)
}

func (q *DelayStorage) getZSetName(timeSlice string) string {
	return fmt.Sprintf(zsetNameFormat, timeSlice, q.getBucket(timeSlice))
}

func (q *DelayStorage) getHashtableName(timeSlice string) string {
	return fmt.Sprintf(hashtableNameFormat, timeSlice, q.getBucket(timeSlice))
}

func (q *DelayStorage) getBucket(timeSlice string) int8 {
	return 1
}
