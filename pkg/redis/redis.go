// Package redis 实现了统一的 Redis 客户端，并提供基础的分布式缓存与分布式锁功能封装
package redis

import (
	"context"
	"errors"
	"runtime"
	"sync"
	"time"

	"github.com/bsm/redislock"
	"github.com/go-redis/redis/v9"
	"github.com/go-redis/redis_rate/v10"

	"github.com/beihai0xff/pudding/pkg/configs"
)

var (

	// ErrConsumerGroupExists 该
	ErrConsumerGroupExists = errors.New("BUSYGROUP Consumer Group name already exists")
	clientOnce             sync.Once
	client                 *Client
)

// Client Redis 客户端
type Client struct {
	client *redis.Client
	config *configs.RedisConfig
	// Create a new locker client.
	locker *redislock.Client
}

// New 获取客户端
func New(c *configs.RedisConfig) *Client {
	clientOnce.Do(
		func() {
			opt, err := redis.ParseURL(c.RedisURL)
			if err != nil {
				panic(err)
			}

			opt.DialTimeout = time.Duration(c.DialTimeout) * time.Second
			opt.PoolSize = 40 * runtime.GOMAXPROCS(runtime.NumCPU())

			client = &Client{
				client: redis.NewClient(opt),
				config: c,
			}
			client.locker = redislock.New(client.client)
		})

	return client
}

func (c *Client) GetClient() *redis.Client {
	return c.client
}

/*
	JSON 相关 Command
*/

// Set 执行 Redis SET 命令，expireTime 时间单位为秒
func (c *Client) Set(ctx context.Context, key, value string, expireTime time.Duration) error {
	if key == "" || value == "" {
		return errors.New("redis SET key or value can't be empty")
	}
	return c.client.Set(ctx, key, value, expireTime).Err()
}

// Get 执行 Redis GET 命令
func (c *Client) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

/*
	ZSet 相关 Command
*/

// ZAddNX 执行 Redis ZAdd 命令
func (c *Client) ZAddNX(ctx context.Context, key string, members ...redis.Z) error {
	return c.client.ZAddNX(ctx, key, members...).Err()
}

// ZRangeByScore 执行 Redis ZRangeByScore 命令
func (c *Client) ZRangeByScore(ctx context.Context, key string, opt *redis.ZRangeBy) ([]redis.Z, error) {
	return c.client.ZRangeByScoreWithScores(ctx, key, opt).Result()
}

// ZRem 执行 Redis ZRem 命令
func (c *Client) ZRem(ctx context.Context, key string, members ...interface{}) error {
	return c.client.ZRem(ctx, key, members...).Err()
}

/*
	HashTable 相关 Command
*/

// HGet 执行 Redis HGet 命令
func (c *Client) HGet(ctx context.Context, key, field string) ([]byte, error) {
	return c.client.HGet(ctx, key, field).Bytes()
}

// Del 执行 Redis Del 命令
func (c *Client) Del(ctx context.Context, keys string) error {
	return c.client.Del(ctx, keys).Err()
}

/*
	Stream 相关 Command
*/

// StreamSend 向指定 Stream 发送消息
func (c *Client) StreamSend(ctx context.Context, streamName string, msg []byte) error {
	return c.client.XAdd(ctx, &redis.XAddArgs{
		Stream: streamName,
		MaxLen: 100000,
		Approx: true,
		ID:     "*",
		Values: []interface{}{"body", msg},
	}).Err()
}

// XGroupCreate start 参数表示该消费者组从哪个位置开始消费消息，可以指定为 ID 或 $，其中 $ 表示从最后一条消息开始消费。
func (c *Client) XGroupCreate(ctx context.Context, topic, group, start string) error {
	// 创建消费者组
	// 但是 XGroupCreate 这个命令不幂等，不能重复创建同一个消费者组，
	// 如果 group 已经存在了则会返回错误
	// 也不能在不存在的 stream 上创建 group
	err := c.client.XGroupCreateMkStream(ctx, topic, group, start).Err()
	if err != nil && errors.Is(err, ErrConsumerGroupExists) {
		return err
	}

	return nil
}

func (c *Client) XGroupConsume(ctx context.Context, topic, group, consumerName, id string,
	batchSize int) ([]redis.XMessage, error) {
	// 阻塞的获取消息
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

func (c *Client) XGroupDelConsumer(ctx context.Context, topic, group, consumerName string) error {
	_, err := c.client.XGroupDelConsumer(ctx, topic, group, consumerName).Result()

	return err
}

func (c *Client) XAck(ctx context.Context, topic, group string, ids ...string) error {
	return c.client.XAck(ctx, topic, group, ids...).Err()
}

/*
	DistributeLock 相关
*/

// GetDistributeLock 获取一个分布式锁
func (c *Client) GetDistributeLock(ctx context.Context, name string,
	expireTime time.Duration) (*redislock.Lock, error) {
	return c.locker.Obtain(ctx, name, expireTime, nil)
}

/*
	leaky bucket
*/

func (c *Client) GetLimiter() *redis_rate.Limiter {
	return redis_rate.NewLimiter(client.client)
}

// Close redis client
func (c *Client) Close() error {
	return c.client.Close()
}
