package redis

import (
	"github.com/alicebob/miniredis/v2"

	"github.com/beihai0xff/pudding/configs"
)

var (
	table, zset = "hash_test", "zset_test"
	key, value  = "GetKey", "GetValue"
)

// NewMockRdb get Redis mock client
func NewMockRdb() *Client {
	// initial Redis DB
	s, _ := miniredis.Run()

	s.ZAdd(zset, 1, "a")
	s.ZAdd(zset, 2, "b")
	s.ZAdd(zset, 3, "c")

	s.HSet(table, key, value)

	c := &configs.RedisConfig{
		RedisURL:    "redis://" + s.Addr(),
		DialTimeout: 5,
	}

	return New(c)
}
