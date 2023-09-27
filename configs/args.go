// Package configs provides config management
// args.go provides some common command line arguments.
package configs

import (
	"flag"
)

const (
	// serverConfigPath server config path
	serverConfigPath = "server_config.%s"
)

var (
	configFilepath = flag.String("config-filepath", "./config.yaml", "The server config file path")
	_              = flag.String("host-domain", DefaultHostDomain, "The server host domain")
	_              = flag.Int("grpc-port", DefaultGRPCPort, "The grpc server port")
	_              = flag.Int("http-port", DefaultHTTPPort, "The http server port")
	_              = flag.String("name-server-url", DefaultNameServerURL, "The name server connection url")
)

// GetConfigFilePath get the config file path
func GetConfigFilePath() *string {
	return configFilepath
}
