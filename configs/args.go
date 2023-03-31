// Package configs provides config management
// args.go provides some common command line arguments.
package configs

import (
	"flag"
)

var (
	// ConfigPath config file path
	ConfigPath = flag.String("config_file_path", "./config.yaml", "The server config file path")

	// GRPCPort grpc server port
	GRPCPort = flag.Int("grpc_port", DefaultGRPCPort, "The grpc server grpcPort")
	// HTTPPort http server port
	HTTPPort = flag.Int("http_port", DefaultHTTPPort, "The http server grpcPort")

	// CertPath tls cert file path, the file must contain PEM encoded data.
	CertPath = flag.String("cert_path", DefaultCertPath, "The tls cert file path")
	// KeyPath tls key file path, the file must contain PEM encoded data.
	KeyPath = flag.String("key_path", DefaultKeyPath, "The tls key file path")
	// NameServerURL name server url
	NameServerURL = flag.String("name_server_url", DefaultNameServerURL, "The name server connection url")
)
