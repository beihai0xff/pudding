package configs

// RedisConfig Redis 配置
type RedisConfig struct {
	Network     string `json:"network"`
	Address     string `json:"address"`
	Database    int    `json:"database"`
	Password    string `json:"password"`
	MaxIdle     int    `json:"maxIdle"`
	IdleTimeout int    `json:"idleTimeout"`
}

func GetRedisConfig() *RedisConfig {
	return &RedisConfig{
		Network: "tcp",
		Address: "",
	}
}
