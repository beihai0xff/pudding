// Package configs provides config management
// server_base_config.go contains the base config of a server
package configs

import "github.com/beihai0xff/pudding/pkg/log"

const (
	// baseConfigPath server base config path
	baseConfigPath = "server_config.base_config"

	// DefaultHostDomain is the default server domain
	DefaultHostDomain = "localhost"
	// DefaultGRPCPort is the default port for gRPC server.
	DefaultGRPCPort = 50050
	// DefaultHTTPPort is the default port for HTTP server.
	DefaultHTTPPort = 8080
	// DefaultEnableTLS is enabled TLS
	DefaultEnableTLS = false
	// DefaultCertPath is the default path for TLS certificate.
	DefaultCertPath = "./certs/pudding.pem"
	// DefaultKeyPath is the default path for TLS key.
	DefaultKeyPath = "./certs/pudding-key.pem"
	// DefaultNameServerURL is the default name server connection url
	DefaultNameServerURL = ""
)

// OptionFunc is the option function for config.
type OptionFunc func(map[string]interface{})

// BaseConfig server base Config
// support CommandLine
type BaseConfig struct {
	// server domain
	HostDomain string `json:"host_domain" yaml:"host_domain" mapstructure:"host_domain"`
	// GRPCPort grpc server port
	GRPCPort int `json:"grpc_port" yaml:"grpc_port" mapstructure:"grpc_port"`
	// HTTPPort http server port
	HTTPPort int `json:"http_port" yaml:"http_port" mapstructure:"http_port"`
	// EnableTLS is enabled TLS, default is false
	EnableTLS bool `json:"enable_tls" yaml:"enable_tls" mapstructure:"enable_tls"`
	// CertPath tls cert file path, the file must contain PEM encoded data.
	CertPath string `json:"cert_path" yaml:"cert_path" mapstructure:"cert_path"`
	// KeyPath tls key file path, the file must contain PEM encoded data.
	KeyPath string `json:"key_path" yaml:"key_path" mapstructure:"key_path"`
	// NameServerURL name server url
	NameServerURL string `json:"name_server_url" yaml:"name_server_url" mapstructure:"name_server_url"`
	// Logger log config for output config message, do not use it
	Logger []log.Config `json:"log_config" yaml:"log_config" mapstructure:"log_config"`
}
