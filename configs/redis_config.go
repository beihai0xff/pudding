// Package configs provides config management
package configs

import "fmt"

const redisConfigPath = "redis_config"

// RedisConfig Redis config
type RedisConfig struct {
	// URL is the redis connection url
	URL string `json:"url" yaml:"url" mapstructure:"url"`
	// DialTimeout Dial timeout for establishing new connections.
	// Default is 5 seconds.
	DialTimeout int `json:"dial_timeout" yaml:"dial_timeout" mapstructure:"dial_timeout"`
}

// WithRedisURL set the redis url.
func WithRedisURL(url string) OptionFunc {
	return func(confMap map[string]interface{}) {
		if url != "" {
			confMap[fmt.Sprintf("%s.url", redisConfigPath)] = url
		}
	}
}
