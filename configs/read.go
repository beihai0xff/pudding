package configs

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type OptionFunc func(path, format string)

func Parse(filePath, format string, c interface{}, opt OptionFunc) {
	opt(filePath, format)
	if err := viper.Unmarshal(c); err != nil {
		panic(fmt.Errorf("failed to unmarshal config: %w", err))
	}
}

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
