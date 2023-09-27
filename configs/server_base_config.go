// Package configs provides config management
// server_base_config.go contains the base config of a server
package configs

const (
	// DefaultHostDomain is the default server domain
	DefaultHostDomain = "localhost"
	// DefaultGRPCPort is the default port for gRPC server.
	DefaultGRPCPort = 50050
	// DefaultHTTPPort is the default port for HTTP server.
	DefaultHTTPPort = 8080
	// DefaultTLSEnable is enabled TLS
	DefaultTLSEnable = false
	// DefaultEnableAutocert is enabled TLS autocert
	DefaultEnableAutocert = false
	// DefaultCertPath is the default path for TLS certificate.
	DefaultCertPath = "./certs/pudding.pem"
	// DefaultKeyPath is the default path for TLS key.
	DefaultKeyPath = "./certs/pudding-key.pem"
	// DefaultNameServerURL is the default name server connection url
	DefaultNameServerURL = ""
)

// OptionFunc is the option function for config.
type OptionFunc func(map[string]any)

// BaseConfig server base Config
// support CommandLine
type BaseConfig struct {
	// server domain
	HostDomain string `json:"host_domain" yaml:"host_domain" mapstructure:"host_domain"`
	// GRPCPort grpc server port
	GRPCPort int `json:"grpc_port" yaml:"grpc_port" mapstructure:"grpc_port"`
	// HTTPPort http server port
	HTTPPort int `json:"http_port" yaml:"http_port" mapstructure:"http_port"`
	// NameServerURL name server url
	NameServerURL string `json:"name_server_url" yaml:"name_server_url" mapstructure:"name_server_url"`
	// TLS config
	TLS *TLS `json:"tls" yaml:"tls" mapstructure:"tls"`
}

// TLS config
type TLS struct {
	// Enable is enabled TLS, default is false
	Enable bool `json:"enable" yaml:"enable" mapstructure:"enable"`
	// AutoCert is enabled TLS autocert, default is false
	AutoCert bool `json:"autocert" yaml:"autocert" mapstructure:"autocert"`
	// CertPath tls cert file path, the file must contain PEM encoded data.
	CACert string `json:"ca_cert" yaml:"ca_cert" mapstructure:"ca_cert"`
	// ServerCert tls server cert file path, the file must contain PEM encoded data.
	ServerCert string `json:"server_cert" yaml:"server_cert" mapstructure:"server_cert"`
	// ServerKey tls server key file path, the file must contain PEM encoded data.
	ServerKey string `json:"server_key" yaml:"server_key" mapstructure:"server_key"`
	// ClientCert tls client cert file path, the file must contain PEM encoded data.
	ClientCert string `json:"client_cert" yaml:"client_cert" mapstructure:"client_cert"`
	// ClientKey tls client key file path, the file must contain PEM encoded data.
	ClientKey string `json:"client_key" yaml:"client_key" mapstructure:"client_key"`
}
