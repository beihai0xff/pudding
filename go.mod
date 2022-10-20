module github.com/beihai0xff/pudding

go 1.18

require (
	github.com/alicebob/miniredis/v2 v2.23.0
	github.com/bsm/redislock v0.8.1
	github.com/go-redis/redis/v9 v9.0.0-rc.1
	github.com/go-redis/redis_rate/v10 v10.0.0-00010101000000-000000000000
	github.com/google/uuid v1.3.0
	github.com/klauspost/compress v1.15.11
	github.com/natefinch/lumberjack v2.0.0+incompatible
	github.com/stretchr/testify v1.8.0
	github.com/vmihailenco/msgpack/v5 v5.3.5
	go.uber.org/zap v1.23.0
	gopkg.in/yaml.v3 v3.0.1
)

replace github.com/go-redis/redis_rate/v10 => github.com/beihai0xff/redis_rate/v10 v10.0.0-20221018024645-a6a5d00e2135

require (
	github.com/BurntSushi/toml v1.2.0 // indirect
	github.com/alicebob/gopher-json v0.0.0-20200520072559-a9ecdc9d1d3a // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/vmihailenco/tagparser/v2 v2.0.0 // indirect
	github.com/yuin/gopher-lua v0.0.0-20210529063254-f4c35e4016d9 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/multierr v1.6.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
)
