package redis

import (
	"context"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/go-redis/redis/v9"
)

var (
	c *Client
)

// redis_test.go 测试文件对 Redis 客户端对外暴漏的方法进行了功能测试，连接的是 dev 环境的数据库 。
// 下面的单元测试也可以作为使用范例参考

func TestMain(m *testing.M) {

	// initial Redis Client
	c = NewMockRdb()

	exitCode := m.Run()
	// 退出
	os.Exit(exitCode)
}

// TestClient_Set 测试 Set 方法
func TestClient_Set(t *testing.T) {

	type args struct {
		key        string
		value      string
		expireTime time.Duration
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Set 测试正常数据", args{"SetKey0", "SetValue0", 60 * time.Second}, false},
		{"Set 测试数据过期时间为0", args{"SetKey1", "SetValue1", 0}, false},
		{"Set 测试 key 为 empty", args{"", "SetValue", 60 * time.Second}, true},
		{"Set 测试 value 为 empty", args{"SetKey", "", 60 * time.Second}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := c.Set(context.Background(), tt.args.key, tt.args.value, tt.args.expireTime); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestClient_Get 测试 Get 方法
func TestClient_Get(t *testing.T) {
	c.Set(context.Background(), "GetKey", "GetValue", 60*time.Second)
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"Get 测试已存储在 Redis 的数据", args{"GetKey"}, "GetValue", false},
		{"Get 测试未存储在 Redis 的数据", args{"GetKey100"}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.Get(context.Background(), tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_HGet(t *testing.T) {

	type args struct {
		key   string
		field string
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{"HGet 测试已经存储在 Redis 的数据", args{table, key}, []byte("GetValue"), false},
		{"HGet 测试未存储在 Redis 的数据", args{"no", key}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.HGet(context.Background(), tt.args.key, tt.args.field)
			if (err != nil) != tt.wantErr {
				t.Errorf("HGet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("HGet() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_ZRangeByScore(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
		opt *redis.ZRangeBy
	}
	tests := []struct {
		name    string
		args    args
		want    []redis.Z
		wantErr bool
	}{
		{"ZSet 获取已经存储在 Redis 的数据", args{context.Background(), zset, &redis.ZRangeBy{
			Min:    "-inf",
			Max:    strconv.FormatInt(time.Now().Unix(), 10),
			Offset: 0,
			Count:  10,
		}}, []redis.Z{{1, "a"}, {2, "b"}, {3, "c"}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.ZRangeByScore(tt.args.ctx, tt.args.key, tt.args.opt)
			if (err != nil) != tt.wantErr {
				t.Errorf("ZRangeByScore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ZRangeByScore() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_DelMap(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"删除整个哈希表 的数据", args{context.Background(), table}, false},
		{"删除不存在的 key", args{context.Background(), "unknow key"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := c.Del(tt.args.ctx, tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("Del() error = %v, wantErr %v", err, tt.wantErr)
			}
			if res, _ := c.HGet(tt.args.ctx, tt.args.key, key); res != nil {
				t.Errorf("the key already delete, but got: %v", res)
			}
		})
	}
}
