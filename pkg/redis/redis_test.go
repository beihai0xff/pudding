package redis

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/beihai0xff/pudding/configs"
)

// redis_test.go 测试文件对 Redis 客户端对外暴漏的方法进行了功能测试，连接的是 dev 环境的数据库 。
// 下面的单元测试也可以作为使用范例参考

var (
	c = GetRDB(configs.GetRedisConfig())
)

func TestMain(m *testing.M) {
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
	_ = c.Set(context.Background(), "GetKey", "GetValue", 60*time.Second)
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"Get 测试已经存储在 Redis的数据", args{"GetKey"}, "GetValue", false},
		{"Get 测试未存储在 Redis的数据", args{"GetKey100"}, "", true},
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

// TestClient_GetDistributeLock 测试 GetDistributeLock 方法
func TestClient_GetDistributeLock(t *testing.T) {
	lock, _ := c.GetDistributeLock(context.Background(), "DLockUsed", 10*time.Second)
	defer lock.Release(context.Background())
	type args struct {
		name string
	}
	tests := []struct {
		name         string
		args         args
		wantDLockErr bool
		unlock       bool
	}{
		{"DistributeLock 测试获取未被使用的锁", args{"DLockUnused"}, false, true},
		{"DistributeLock 测试获取已经释放的锁", args{"DLockUnused"}, false, false},
		{"DistributeLock 测试获取已经使用的锁", args{"DLockUsed"}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mutex, err := c.GetDistributeLock(ctx, tt.args.name, 2*time.Second)
			if (err != nil) != tt.wantDLockErr {
				t.Errorf("mutex Lock error = %v, wantErr %v", err, tt.wantDLockErr)
				return
			}
			if err == ErrNotObtained {
				fmt.Println(err)
			}
			if tt.unlock {
				mutex.Release(ctx)
			}
		})
	}
}
