// Package redis 实现了统一的 Redis 客户端，并提供基础的分布式缓存与分布式锁功能封装
package redis

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/bsm/redislock"
	"github.com/go-redis/redis/v9"

	"github.com/beihai0xff/pudding/configs"
)

var (
	// ErrNotObtained 分布式锁加锁失败：该锁已经被占用。配合 GetDistributeLock 使用
	ErrNotObtained = redislock.ErrNotObtained
	// ErrConsumerGroupExists 该
	ErrConsumerGroupExists = errors.New("BUSYGROUP Consumer Group name already exists")
	ClientOnce             sync.Once
	client                 *Client
)

// Client Redis 客户端
type Client struct {
	client *redis.Client
	Config *configs.RedisConfig
}

// GetRDB 获取客户端
func GetRDB(c *configs.RedisConfig) *Client {
	ClientOnce.Do(
		func() {
			opt, err := redis.ParseURL(c.RedisURL)
			if err != nil {
				panic(err)
			}

			client = &Client{
				client: redis.NewClient(opt),
				Config: c,
			}
		})

	return client
}

func (c *Client) GetClient() *redis.Client {
	return c.client
}

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

// HGet 执行 Redis HGet 命令
func (c *Client) HGet(ctx context.Context, key, field string) ([]byte, error) {
	return c.client.HGet(ctx, key, field).Bytes()
}

// GetDistributeLock 获取一个分布式锁
func (c *Client) GetDistributeLock(ctx context.Context, name string,
	expireTime time.Duration) (*redislock.Lock, error) {
	// Create a new lock client.
	locker := redislock.New(c.client)

	// Try to obtain lock.
	return locker.Obtain(ctx, name, expireTime, nil)
}

// Close redis client
func (c *Client) Close() error {
	return c.client.Close()
}

func (c *Client) SteamSend(ctx context.Context, streamName string, msg []byte) error {
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
	if err != nil && err != ErrConsumerGroupExists {
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
