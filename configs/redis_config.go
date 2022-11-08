package configs

// RedisConfig Redis 配置
type RedisConfig struct {
	URL         string `json:"url" yaml:"url" mapstructure:"url"`
	DialTimeout int    `json:"dial_timeout" yaml:"dial_timeout" mapstructure:"dial_timeout"`
}
