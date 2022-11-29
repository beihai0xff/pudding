package configs

import conf "github.com/beihai0xff/pudding/configs"

var c = &Config{
	ConsulURL: "",
}

// Config is the config for trigger module.
type Config struct {

	// MySQL config
	MySQL *conf.MySQLConfig `json:"mysql_config" yaml:"mysql_config" mapstructure:"mysql_config"`

	// Logger log config for output config message, do not use it
	Logger    map[string]*conf.LogConfig `json:"log_config" yaml:"log_config" mapstructure:"log_config"`
	ConsulURL string                     `json:"consul_url" yaml:"consul_url" mapstructure:"consul_url"`
}

// GetMySQLConfig returns the MySQL config.
func GetMySQLConfig() *conf.MySQLConfig {
	return c.MySQL
}

// GetConsulURL returns the consul url.
func GetConsulURL() string {
	return c.ConsulURL
}
