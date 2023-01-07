// Package configs provides config management
package configs

import (
	"github.com/beihai0xff/pudding/pkg/grpc/args"
)

const baseConfigPath = "server_config.base_config"

// BaseConfig server base Config
// support CommandLine
type BaseConfig struct {
	// GRPCPort grpc server port
	GRPCPort int `json:"grpc_port" yaml:"grpc_port" mapstructure:"grpc_port"`
	// HTTPPort http server port
	HTTPPort int `json:"http_port" yaml:"http_port" mapstructure:"http_port"`
	// CertPath tls cert file path, the file must contain PEM encoded data.
	CertPath string `json:"cert_path" yaml:"cert_path" mapstructure:"cert_path"`
	// KeyPath tls key file path, the file must contain PEM encoded data.
	KeyPath string `json:"key_path" yaml:"key_path" mapstructure:"key_path"`
	// NameServerURL name server url
	NameServerURL string `json:"name_server_url" yaml:"name_server_url" mapstructure:"name_server_url"`
	// Logger log config for output config message, do not use it
	Logger map[string]*LogConfig `json:"log_config" yaml:"log_config" mapstructure:"log_config"`
}

// SetFlags set flags to BaseConfig
// flags have the highest priority
func (c *BaseConfig) SetFlags() {
	// if flag changes or not set value in other config file, use flag value
	if *args.GRPCPort != args.DefaultGRPCPort || c.GRPCPort == 0 {
		c.GRPCPort = *args.GRPCPort
	}
	if *args.HTTPPort != args.DefaultHTTPPort || c.HTTPPort == 0 {
		c.HTTPPort = *args.HTTPPort
	}
	if *args.CertPath != args.DefaultCertPath || c.CertPath == "" {
		c.CertPath = *args.CertPath
	}
	if *args.KeyPath != args.DefaultKeyPath || c.KeyPath == "" {
		c.KeyPath = *args.KeyPath
	}
}

// BrokerConfig BrokerConfig Config
type BrokerConfig struct {
	// BaseConfig server base Config
	BaseConfig `json:"base_config" yaml:"base_config" mapstructure:"base_config"`

	// TimeSliceInterval broker loop time interval
	TimeSliceInterval string `json:"time_slice_interval" yaml:"time_slice_interval" mapstructure:"time_slice_interval"`
	// MessageTopic default message topic, if no topic set in message, use this topic
	MessageTopic string `json:"message_topic" yaml:"message_topic" mapstructure:"message_topic"`
	// TokenTopic TimeSlice token topic
	TokenTopic string `json:"token_topic" yaml:"token_topic" mapstructure:"token_topic"`
	// Broker type
	Broker string `json:"broker" yaml:"broker" mapstructure:"broker"`
	// Connector type
	Connector string `json:"connector" yaml:"connector" mapstructure:"connector"`
}

// TriggerConfig Trigger server config
type TriggerConfig struct {
	// BaseConfig server base Config
	BaseConfig `json:"base_config" yaml:"base_config" mapstructure:"base_config"`

	// WebhookPrefix is the prefix of webhook url.
	WebhookPrefix string `json:"webhook_prefix" yaml:"webhook_prefix" mapstructure:"webhook_prefix"`
	// SchedulerConsulURL is the scheduler consul connection url.
	SchedulerConsulURL string `json:"scheduler_consul_url" yaml:"scheduler_consul_url" mapstructure:"scheduler_consul_url"`
}
