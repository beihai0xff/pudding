// Package configs provides config management
// args.go provides some common command line arguments.
package configs

import (
	"flag"
	"fmt"
	"strings"
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
	_              = flag.Bool("enable-tls", DefaultEnableTLS, "Is enable server tls")
	_              = flag.String("cert-path", DefaultCertPath, "The TLS cert file path")
	_              = flag.String("key-path", DefaultKeyPath, "The TLS key file path")
	_              = flag.String("name-server-url", DefaultNameServerURL, "The name server connection url")
)

func flagProvider(mp map[string]interface{}) {
	// It visits only those flags that have been set.
	flag.Visit(func(f *flag.Flag) {
		key := strings.ReplaceAll(fmt.Sprintf(serverConfigPath, f.Name), "-", "_")
		mp[key] = f.Value.(flag.Getter).Get()
	})
}

// GetConfigFilePath get the config file path
func GetConfigFilePath() *string {
	return configFilepath
}
