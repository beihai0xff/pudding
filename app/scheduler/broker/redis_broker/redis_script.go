package redis_broker

import "github.com/go-redis/redis/v9"

var (
	pushScript = redis.NewScript(`
-- KEYS[1]: ZSet topic
-- KEYS[2]: Hashtable topic
-- ARGV[1]: Message Key
-- ARGV[2]: Message
-- ARGV[3]: Message Ready Time（now + delay）
local getTopicPartition = KEYS[1]
local getHashtableName = KEYS[2]
local key = ARGV[1]
local message = ARGV[2]
local readyTime = tonumber(ARGV[3])
-- add Message Key and ReadyTime to zset
local count = redis.call("zadd", getTopicPartition, readyTime, key)
-- Message already exists
if count == 0 then
   return 0
end
-- add Message Content to hashtable
redis.call("hsetnx", getHashtableName, key, message)
return 1
`)

	deleteScript = redis.NewScript(`
-- KEYS[1]: zset Name
-- KEYS[2]: hashtable Name
-- ARGV[1]: Message Key
local zset = KEYS[1]
local hashtable = KEYS[2]
local key = ARGV[1]
-- 删除zset和hash关于这条消息的内容
redis.call("zrem", zset, key)
redis.call("hdel", hashtable, key)
return 1
`)
)
