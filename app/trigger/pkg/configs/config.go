package configs

import conf "github.com/beihai0xff/pudding/configs"

var c = &Config{
	ConsulURL: "",
}

// Config is the config for trigger module.`
type Config struct {

	// MySQL config
	MySQL *conf.MySQLConfig `json:"mysql_config" yaml:"mysql_config" mapstructure:"mysql_config"`
	// Logger log config for output config message, do not use it
	Logger map[string]*conf.LogConfig `json:"log_config" yaml:"log_config" mapstructure:"log_config"`
	// ConsulURL is the consul connection url.
	ConsulURL string `json:"consul_url" yaml:"consul_url" mapstructure:"consul_url"`
	// WebhookPrefix is the prefix of webhook url.
	WebhookPrefix string `json:"webhook_prefix" yaml:"webhook_prefix" mapstructure:"webhook_prefix"`
	// SchedulerConsulURL is the scheduler consul connection url.
	SchedulerConsulURL string `json:"scheduler_consul_url" yaml:"scheduler_consul_url" mapstructure:"scheduler_consul_url"`
}

// GetMySQLConfig returns the MySQL config.
func GetMySQLConfig() *conf.MySQLConfig {
	return c.MySQL
}

// GetConsulURL returns the consul url.
func GetConsulURL() string {
	return c.ConsulURL
}

// GetWebhookPrefix returns the webhook prefix.
func GetWebhookPrefix() string {
	return c.WebhookPrefix
}

// GetSchedulerConsulURL returns the scheduler consul connection url.
func GetSchedulerConsulURL() string {
	return c.SchedulerConsulURL
}
