package redis

import (
	"context"
	"os"
	"reflect"
	"strconv"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	c *Client
)

func TestMain(m *testing.M) {

	// initial Redis Client
	c = NewMockRdb()

	exitCode := m.Run()
	os.Exit(exitCode)
}

// TestClient_Set
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
		{"Set normal test data", args{"SetKey0", "SetValue0", 60 * time.Second}, false},
		{"Set expire time is zero", args{"SetKey1", "SetValue1", 0}, false},
		{"Set key is empty", args{"", "SetValue", 60 * time.Second}, true},
		{"Set value is empty", args{"SetKey", "", 60 * time.Second}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := c.Set(context.Background(), tt.args.key, tt.args.value, tt.args.expireTime); (err != nil) != tt.wantErr {
				t.Errorf("Set() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestClient_Get test Get function
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
		{"Get exist data in Redis", args{"GetKey"}, "GetValue", false},
		{"Get non-exist data in Redis", args{"GetKey100"}, "", true},
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
		{"HGet exist data in Redis", args{table, key}, []byte("GetValue"), false},
		{"HGet non-exist data in Redis", args{"no", key}, nil, true},
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
		{"ZSet exist data in Redis", args{context.Background(), zset, &redis.ZRangeBy{
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
		{"delete all data in map", args{context.Background(), table}, false},
		{"delete non-exist key in map", args{context.Background(), "unknow key"}, false},
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
