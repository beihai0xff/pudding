package configs

// RedisConfig Redis 配置
type RedisConfig struct {
	RedisURL    string `json:"redisURL"`
	Network     string `json:"network"`
	MaxIdle     int    `json:"maxIdle"`
	IdleTimeout int    `json:"idleTimeout"`
}

func GetRedisConfig() *RedisConfig {
	return &RedisConfig{
		Network: "tcp",
	}
}
