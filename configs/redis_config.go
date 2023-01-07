// Package configs provides config management
package configs

// RedisConfig Redis 配置
type RedisConfig struct {
	// URL is the redis connection url
	URL string `json:"url" yaml:"url" mapstructure:"url"`
	// DialTimeout Dial timeout for establishing new connections.
	// Default is 5 seconds.
	DialTimeout int `json:"dial_timeout" yaml:"dial_timeout" mapstructure:"dial_timeout"`
}
