// Package configs provides config management
// server_base_config.go contains the base config of a server
package configs

import "github.com/beihai0xff/pudding/pkg/grpc/args"

const baseConfigPath = "server_config.base_config"

// OptionFunc is the option function for config.
type OptionFunc func(map[string]interface{})

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

var baseConfig *BaseConfig

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
	baseConfig = c
}
