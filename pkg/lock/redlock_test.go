package lock

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/bsm/redislock"

	rdb "github.com/beihai0xff/pudding/pkg/redis"
)

var testRedLockClient *RedLockClient

func TestMain(m *testing.M) {

	// initial Redis database
	rdb := rdb.NewMockRdb()

	testRedLockClient = NewRedLockClient(rdb)

	exitCode := m.Run()
	os.Exit(exitCode)
}

// Test_NewRedLock 测试 NewRedLock 方法
func Test_NewRedLock(t *testing.T) {

	lock, _ := testRedLockClient.NewRedLock(context.Background(), "DLockUsed", 10*time.Second)
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
		{"DistributeLock get unused locker", args{"DLockUnused"}, false, true},
		{"DistributeLock get released locker", args{"DLockUnused"}, false, false},
		{"DistributeLock get locked locker", args{"DLockUsed"}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			mutex, err := testRedLockClient.NewRedLock(ctx, tt.args.name, 2*time.Second)
			if (err != nil) != tt.wantDLockErr {
				t.Errorf("mutex Lock error = %v, wantErr %v", err, tt.wantDLockErr)
				return
			}
			if err == redislock.ErrNotObtained {
				fmt.Println(err)
			}
			if tt.unlock {
				mutex.Release(ctx)
			}
		})
	}
}
