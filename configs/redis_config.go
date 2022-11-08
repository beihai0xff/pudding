package configs

// RedisConfig Redis 配置
type RedisConfig struct {
	RedisURL    string `json:"redisURL" yaml:"redisURL" mapstructure:"redisURL"`
	DialTimeout int    `json:"dialTimeout" yaml:"dialTimeout" mapstructure:"dialTimeout"`
}
