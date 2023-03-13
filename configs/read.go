// Package configs provides config management
package configs

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// ParseOptionFunc Parse config option func
type ParseOptionFunc func(path, format string)

const (
	ConfigFormatYAML = "yaml"
)

// Parse ReadConfig read config from filePath with format
func Parse(filePath, format string, c interface{}, opt ParseOptionFunc) {
	opt(filePath, format)
	if err := viper.Unmarshal(c); err != nil {
		panic(fmt.Errorf("failed to unmarshal config: %w", err))
	}
}

// ReadFromFile read config from filePath with format
func ReadFromFile(filePath, format string) {
	viper.SetConfigFile(filePath)
	viper.SetConfigType(format)

	// Find and read the config file
	if err := viper.ReadInConfig(); err != nil { // Handle errors reading the config file
		if errors.As(err, &viper.ConfigFileNotFoundError{}) {
			panic(fmt.Errorf("config file not found: %w", err))
		} else {
			// Config file was found but another error was produced
			panic(fmt.Errorf("failed to read config file: %w", err))
		}
	}
}

// ReadFromConsul read config from consul with format
func ReadFromConsul(filePath, format string) {
	path := strings.SplitN(filePath, "/", 2)
	if err := viper.AddRemoteProvider("consul", path[0], path[1]); err != nil {
		panic(fmt.Errorf("failed to add remote provider: %w", err))
	}
	viper.SetConfigType(format)

	// Find and read the config file
	if err := viper.ReadRemoteConfig(); err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("failed to read config file: %w", err))
	}
}
