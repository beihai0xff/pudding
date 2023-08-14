// Package configs provides config management
// server_trigger_config.go contains the config of trigger server
package configs

import "github.com/beihai0xff/pudding/pkg/log"

// TriggerConfig Trigger server config
type TriggerConfig struct {
	// BaseConfig server base Config
	BaseConfig `json:"server_config" yaml:"server_config" mapstructure:"server_config"`
	// ServerConfig server config
	// use same struct tag merge BaseConfig to ServerConfig
	//nolint:govet,revive
	ServerConfig struct {
		// WebhookPrefix is the prefix of webhook url.
		WebhookPrefix string `json:"webhook_prefix" yaml:"webhook_prefix" mapstructure:"webhook_prefix"`
		// SchedulerConsulURL is the scheduler consul connection url.
		SchedulerConsulURL string `json:"scheduler_consul_url" yaml:"scheduler_consul_url" mapstructure:"scheduler_consul_url"`
	} `json:"server_config" yaml:"server_config" mapstructure:"server_config"`

	// Logger log config for output config message
	Logger []log.Config `json:"log_config" yaml:"log_config" mapstructure:"log_config"`
	// MySQLConfig config
	MySQLConfig *MySQLConfig `json:"mysql_config" yaml:"mysql_config" mapstructure:"mysql_config"`
}

// ParseTriggerConfig read the config from the given configPath.
func ParseTriggerConfig(configPath string, opts ...OptionFunc) *TriggerConfig {
	if err := Parse(configPath, ConfigFormatYAML, ReadFromFile, opts...); err != nil {
		panic(err)
	}

	// unmarshal all config to TriggerConfig
	var c = TriggerConfig{}
	if err := UnmarshalToStruct("", &c); err != nil {
		panic(err)
	}

	return &c
}
