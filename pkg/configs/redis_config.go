package configs

// RedisConfig Redis 配置
type RedisConfig struct {
	RedisURL    string `json:"redisURL" yaml:"redisURL"`
	DialTimeout int    `json:"dialTimeout" yaml:"dialTimeout"`
}

func GetRedisConfig() *RedisConfig {
	return c.Redis
}
