// Package redis 实现了统一的 Redis 客户端，并提供基础的分布式缓存与分布式锁功能封装
package redis

import (
	"context"
	"errors"
	"runtime"
	"sync"
	"time"

	"github.com/go-redis/redis/v9"

	"github.com/beihai0xff/pudding/configs"
)

var (

	// ErrConsumerGroupExists 该
	ErrConsumerGroupExists = errors.New("BUSYGROUP Consumer Group name already exists")
	clientOnce             sync.Once
	client                 *Client
)

// Client Redis client wrapper
type Client struct {
	client *redis.Client
	config *configs.RedisConfig
}

// New create a new redis client
func New(c *configs.RedisConfig) *Client {
	clientOnce.Do(
		func() {
			opt, err := redis.ParseURL(c.URL)
			if err != nil {
				panic(err)
			}

			opt.DialTimeout = time.Duration(c.DialTimeout) * time.Second
			opt.PoolSize = 40 * runtime.GOMAXPROCS(runtime.NumCPU())

			client = &Client{
				client: redis.NewClient(opt),
				config: c,
			}
		})

	return client
}

// GetClient get redis client
func (c *Client) GetClient() *redis.Client {
	return c.client
}

/*
	KV related Command
*/

// Set executes the Redis SET command, expireTime time unit is seconds
// If expireTime is 0, it means not to set the expiration time
func (c *Client) Set(ctx context.Context, key, value string, expireTime time.Duration) error {
	if key == "" || value == "" {
		return errors.New("redis SET key or value can't be empty")
	}
	return c.client.Set(ctx, key, value, expireTime).Err()
}

// Get executes the Redis GET command
// it returns the value of key. If the key does not exist the special value nil is returned.
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

/*
	ZSet related Command
*/

// ZAddNX executes the Redis ZAddNX command
// If the member already exists,
// the score is updated and the element reinserted at the right position to ensure the correct ordering.
func (c *Client) ZAddNX(ctx context.Context, key string, members ...redis.Z) error {
	return c.client.ZAddNX(ctx, key, members...).Err()
}

// ZRangeByScore executes the Redis ZRangeByScore command
func (c *Client) ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) ([]redis.Z, error) {
	return c.client.ZRangeByScoreWithScores(ctx, key, opt).Result()
}

// ZRem executes the Redis ZRem command
func (c *Client) ZRem(ctx context.Context, key string, members ...interface{}) error {
	return c.client.ZRem(ctx, key, members...).Err()
}

/*
	HashTable related Command
*/

// HGet executes the Redis HGet command
func (c *Client) HGet(ctx context.Context, key, field string) ([]byte, error) {
	return c.client.HGet(ctx, key, field).Bytes()
}

// Del executes the Redis Del 命令
func (c *Client) Del(ctx context.Context, keys string) error {
	return c.client.Del(ctx, keys).Err()
}

/*
	Stream related Command
*/

// StreamSend send message to stream
func (c *Client) StreamSend(ctx context.Context, streamName string, msg []byte) error {
	return c.client.XAdd(ctx, &redis.XAddArgs{
		Stream: streamName,
		MaxLen: 100000,
		Approx: true,
		ID:     "*",
		Values: []interface{}{"body", msg},
	}).Err()
}

// XGroupCreate creates a new consumer group
// start argument is the ID of the first message to consume
// start can be specified as ID or $, where $ means from the last message
func (c *Client) XGroupCreate(ctx context.Context, topic, group, start string) error {
	// 创建消费者组
	// 但是 XGroupCreate 这个命令不幂等，不能重复创建同一个消费者组，
	// 如果 group 已经存在了则会返回错误
	// 也不能在不存在的 stream 上创建 group
	err := c.client.XGroupCreateMkStream(ctx, topic, group, start).Err()
	if err != nil && !errors.Is(err, ErrConsumerGroupExists) {
		return err
	}

	return nil
}

// XGroupConsume consumes messages from a stream
func (c *Client) XGroupConsume(ctx context.Context, topic, group, consumerName, id string,
	batchSize int) ([]redis.XMessage, error) {
	// get messages from stream, it will block until there is a message
	result, err := c.client.XReadGroup(ctx, &redis.XReadGroupArgs{
		Group:    group,
		Consumer: consumerName,
		Streams:  []string{topic, id},
		Count:    int64(batchSize),
	}).Result()
	if err != nil {
		return nil, err
	}

	return result[0].Messages, nil
}

// XGroupDelConsumer delete consumer from consumer group
func (c *Client) XGroupDelConsumer(ctx context.Context, topic, group, consumerName string) error {
	_, err := c.client.XGroupDelConsumer(ctx, topic, group, consumerName).Result()

	return err
}

func (c *Client) XAck(ctx context.Context, topic, group string, ids ...string) error {
	return c.client.XAck(ctx, topic, group, ids...).Err()
}

// Close redis client
func (c *Client) Close() error {
	return c.client.Close()
}
