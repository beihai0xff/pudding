// Package configs provides config management
package configs

// ServerConfig server common Config
type ServerConfig struct {
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

// BrokerConfig BrokerConfig Config
type BrokerConfig struct {
	// ServerConfig server common Config
	ServerConfig

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
