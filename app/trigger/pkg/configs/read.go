package configs

import (
	"bytes"
	"encoding/json"
	"log"

	conf "github.com/beihai0xff/pudding/configs"
)

// JSON returns the json format of the config
func (c *Config) JSON() []byte {
	b, err := json.Marshal(c)
	if err != nil {
		log.Panicf("marshal config failed: %v", err)
		return nil
	}

	return b
}

// Init initializes the config
func Init(filePath string, opts ...OptionFunc) {
	conf.Parse(filePath, "yaml", c, conf.ReadFromFile)
	c.ServerConfig.SetFlags()

	for _, opt := range opts {
		opt(c)
	}

	var str bytes.Buffer
	_ = json.Indent(&str, c.JSON(), "", "    ")
	log.Printf("pudding trigger config:\n %s \n", str.String())
}

// OptionFunc is the type of option function
type OptionFunc func(config *Config)

// WithMySQLDSN sets the MySQL dsn
func WithMySQLDSN(dsn string) OptionFunc {
	return func(config *Config) {
		if dsn != "" {
			config.MySQL.DSN = dsn
		}
	}
}

// WithWebhookPrefix sets the webhook prefix
func WithWebhookPrefix(webhookPrefix string) OptionFunc {
	return func(config *Config) {
		if webhookPrefix != "" {
			config.ServerConfig.WebhookPrefix = webhookPrefix
		}
	}
}
