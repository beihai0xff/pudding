// Package configs provides config management
// args.go provides some common command line arguments.
package configs

import (
	"flag"
	"os"
)

var (
	configFlagSet  = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	configFilepath = new(string)
	baseConfig     = &BaseConfig{}
)

func init() {
	configFilepath = configFlagSet.String("config-filepath", "./config.yaml", "The server config file path")

	configFlagSet.StringVar(&baseConfig.HostDomain, "host-domain", DefaultHostDomain, "The server host domain")
	configFlagSet.IntVar(&baseConfig.GRPCPort, "grpc-port", DefaultGRPCPort, "The grpc server port")
	configFlagSet.IntVar(&baseConfig.HTTPPort, "http-port", DefaultHTTPPort, "The http server port")
	configFlagSet.BoolVar(&baseConfig.EnableTLS, "enable-tls", DefaultEnableTLS, "Is enable server tls")
	configFlagSet.StringVar(&baseConfig.CertPath, "tls-cert-path", DefaultCertPath, "The TLS cert file path")
	configFlagSet.StringVar(&baseConfig.KeyPath, "tls-key-path", DefaultKeyPath, "The TLS key file path")
	configFlagSet.StringVar(&baseConfig.NameServerURL, "name-server-url", DefaultNameServerURL, "The name server connection url")
}

// GetConfigFlagSet get the config flag set
func GetConfigFlagSet() *flag.FlagSet {
	return configFlagSet
}

// ParseFlag parse the command line arguments
func ParseFlag() {
	// Ignore errors; CommandLine is set for ExitOnError.
	_ = configFlagSet.Parse(os.Args[1:])
}

// GetConfigFilePath get the config file path
func GetConfigFilePath() *string {
	return configFilepath
}
