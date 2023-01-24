// Package configs defines the config of the pudding.
package configs

import conf "github.com/beihai0xff/pudding/configs"

var c = &Config{}

// Config is the config for trigger module.`
type Config struct {
	// ServerConfig server config
	ServerConfig *conf.TriggerConfig `json:"server_config" yaml:"server_config" mapstructure:"server_config"`

	// MySQL config
	MySQL *conf.MySQLConfig `json:"mysql_config" yaml:"mysql_config" mapstructure:"mysql_config"`
}

// GetMySQLConfig returns the MySQL config.
func GetMySQLConfig() *conf.MySQLConfig {
	return c.MySQL
}

// GetNameServerURL returns the name server url.
func GetNameServerURL() string {
	return c.ServerConfig.NameServerURL
}

// GetServerConfig returns the scheduler config.
func GetServerConfig() *conf.TriggerConfig {
	return c.ServerConfig
}

// GetWebhookPrefix returns the webhook prefix.
func GetWebhookPrefix() string {
	return c.ServerConfig.WebhookPrefix
}

// GetSchedulerConsulURL returns the scheduler consul connection url.
func GetSchedulerConsulURL() string {
	return c.ServerConfig.SchedulerConsulURL
}
