// Package args provides some common command line arguments.
package args

import "flag"

const (
	// DefaultGRPCPort is the default port for gRPC server.
	DefaultGRPCPort = 50050
	// DefaultHTTPPort is the default port for HTTP server.
	DefaultHTTPPort = 8081
	// DefaultCertPath is the default path for TLS certificate.
	DefaultCertPath = "./certs/pudding.pem"
	// DefaultKeyPath is the default path for TLS key.
	DefaultKeyPath = "./certs/pudding-key.pem"
)

var (
	// ConfigPath config file path
	ConfigPath = flag.String("config", "./config.yaml", "The server config file path")

	// GRPCPort grpc server port
	GRPCPort = flag.Int("grpcPort", DefaultGRPCPort, "The grpc server grpcPort")
	// HTTPPort http server port
	HTTPPort = flag.Int("httpPort", DefaultHTTPPort, "The http server grpcPort")

	// CertPath tls cert file path, the file must contain PEM encoded data.
	CertPath = flag.String("cert_path", DefaultCertPath, "The tls cert file path")
	// KeyPath tls key file path, the file must contain PEM encoded data.
	KeyPath = flag.String("key_path", DefaultKeyPath, "The tls key file path")
)
