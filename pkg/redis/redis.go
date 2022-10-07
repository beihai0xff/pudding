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

	"github.com/beihai0xff/pudding/configs"
)

var (
	// ErrNotObtained 分布式锁加锁失败：该锁已经被占用。配合 GetDistributeLock 使用
	ErrNotObtained = redislock.ErrNotObtained
	ClientOnce     sync.Once
	client         *Client
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
			client = &Client{
				client: redis.NewClient(&redis.Options{
					Addr:            c.Address,
					Password:        c.Password,
					DB:              c.Database, // use DB
					Network:         c.Network,
					ConnMaxIdleTime: time.Duration(c.IdleTimeout) * time.Second,
					MinIdleConns:    c.MaxIdle,
					PoolSize:        runtime.NumCPU() * 40,
				}),
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

// GetPipeline 获取一个 Pipeline 对象
func (c *Client) GetPipeline() redis.Pipeliner {
	return c.client.TxPipeline()
}

// GetDistributeLock 获取一个分布式锁
func (c *Client) GetDistributeLock(ctx context.Context, name string,
	expireTime time.Duration) (*redislock.Lock, error) {
	// Create a new lock client.
	locker := redislock.New(c.client)

	// Try to obtain lock.
	return locker.Obtain(ctx, name, expireTime, nil)
}
