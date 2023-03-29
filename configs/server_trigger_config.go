// Package configs provides config management
// server_trigger_config.go contains the config of trigger server
package configs

import (
	"github.com/knadh/koanf/providers/confmap"
)

// TriggerConfig Trigger server config
type TriggerConfig struct {
	ServerConfig struct {
		// BaseConfig server base Config
		BaseConfig `json:"base_config" yaml:"base_config" mapstructure:"base_config"`

		// WebhookPrefix is the prefix of webhook url.
		WebhookPrefix string `json:"webhook_prefix" yaml:"webhook_prefix" mapstructure:"webhook_prefix"`
		// SchedulerConsulURL is the scheduler consul connection url.
		SchedulerConsulURL string `json:"scheduler_consul_url" yaml:"scheduler_consul_url" mapstructure:"scheduler_consul_url"`
	} `json:"server_config" yaml:"server_config" mapstructure:"server_config"`

	// MySQLConfig config
	MySQLConfig *MySQLConfig `json:"mysql_config" yaml:"mysql_config" mapstructure:"mysql_config"`
}

// ParseTriggerConfig read the config from the given configPath.
func ParseTriggerConfig(configPath string, opts ...OptionFunc) *TriggerConfig {
	if err := Parse(configPath, ConfigFormatYAML, ReadFromFile); err != nil {
		panic(err)
	}

	var configMap map[string]interface{}
	for _, opt := range opts {
		opt(configMap)
	}

	if err := k.Load(confmap.Provider(configMap, defaultDelim), nil); err != nil {
		panic(err)
	}

	// unmarshal all config to BrokerConfig{}
	var c = TriggerConfig{}
	if err := UnmarshalToStruct("", &c); err != nil {
		panic(err)
	}

	return &c
}
