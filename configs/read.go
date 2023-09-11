// Package configs provides config management
package configs

import (
	"flag"
	"fmt"
	"strings"

	kjson "github.com/knadh/koanf/parsers/json"
	kyaml "github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/confmap"
	kconsul "github.com/knadh/koanf/providers/consul"
	kenv "github.com/knadh/koanf/providers/env"
	kfile "github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"

	"github.com/beihai0xff/pudding/configs/provider"
	"github.com/beihai0xff/pudding/pkg/log"
)

// Global koanf instance. Use defaultDelim as the key path delimiter. This can be "/" or any character.
var k = koanf.NewWithConf(koanf.Conf{
	Delim:       defaultDelim,
	StrictMerge: true,
})

// ParserFunc Parse config option func
type ParserFunc func(path string, parser koanf.Parser) error

const (
	defaultDelim = "."

	// ConfigFormatYAML config format type yaml
	ConfigFormatYAML = "yaml"
	// ConfigFormatJSON config format type json
	ConfigFormatJSON = "json"
)

// Parse load config from given filePath and format
func Parse(configPath, format string, reader ParserFunc, opts ...OptionFunc) error {
	// First, read config from given config read func, such as file, consul, etc.
	var parser koanf.Parser

	switch format {
	case ConfigFormatYAML:
		parser = kyaml.Parser()
	case ConfigFormatJSON:
		parser = kjson.Parser()
	default:
		return fmt.Errorf("unsupported config format: %s", format)
	}

	if err := reader(configPath, parser); err != nil {
		return err
	}

	// Second, read config from environment variables,
	// Parse environment variables and merge into the loaded config.
	// "PUDDING" is the prefix to filter the env vars by.
	// "." is the delimiter used to represent the key hierarchy in env vars
	// The (optional, or can be nil) function can be used to transform
	// the env var names, for instance, to lowercase them
	if err := k.Load(kenv.Provider("PUDDING_", defaultDelim, func(s string) string {
		return strings.ReplaceAll(strings.ToLower(
			strings.TrimPrefix(s, "PUDDING_")), "_", defaultDelim)
	}), nil); err != nil {
		return fmt.Errorf("error loading config from env: %w", err)
	}

	// third, read config from cli arguments
	if err := k.Load(provider.ProviderWithKey(flag.CommandLine, defaultDelim, k, cliCallback), nil); err != nil {
		panic(err)
	}

	// Finally, read config from given option func
	configMap := map[string]any{}
	for _, opt := range opts {
		opt(configMap)
	}

	if err := k.Load(confmap.Provider(configMap, defaultDelim), nil); err != nil {
		panic(err)
	}

	k.Print()

	return nil
}

func cliCallback(key string, value flag.Value) (string, any) {
	getter, ok := value.(flag.Getter)
	if !ok {
		log.Warnf("flag %s does not implement flag.Getter, skip it", key)
		return "", ""
	}

	return strings.ReplaceAll(fmt.Sprintf(serverConfigPath, key), "-", "_"), getter.Get()
}

// ReadFromFile read config from filePath with format
func ReadFromFile(filePath string, parser koanf.Parser) error {
	// Find and read the config file
	if err := k.Load(kfile.Provider(filePath), parser); err != nil {
		// Config file was found but another error was produced
		return fmt.Errorf("error loading config from file [%s]: %w", filePath, err)
	}

	return nil
}

// ReadFromConsul read config from consul with format
func ReadFromConsul(configPath string, parser koanf.Parser) error {
	// Find and read the config file
	if err := k.Load(kconsul.Provider(kconsul.Config{}), parser); err != nil {
		return fmt.Errorf("error loading config from consul: %w", err)
	}

	return nil
}
