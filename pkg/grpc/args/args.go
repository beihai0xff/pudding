package args

import "flag"

var (
	// GRPCPort grpc server port
	GRPCPort = flag.Int("grpcPort", 50051, "The grpc server grpcPort")
	// HTTPPort http server port
	HTTPPort = flag.Int("httpPort", 8081, "The http server grpcPort")

	// ConfigPath config file path
	ConfigPath = flag.String("config", "./config.yaml", "The server config file path")

	// CertPath tls cert file path, the file must contain PEM encoded data.
	CertPath = flag.String("cert_path", "./certs/pudding.pem", "The tls cert file path")
	// KeyPath tls key file path, the file must contain PEM encoded data.
	KeyPath = flag.String("key_path", "./certs/pudding-key.pem", "The tls key file path")
)
