package configs

import (
	"bytes"
	"encoding/json"
	"log"

	conf "github.com/beihai0xff/pudding/configs"
)

func (c *Config) JSON() []byte {
	b, err := json.Marshal(c)
	if err != nil {
		log.Panicf("marshal config failed: %v", err)
		return nil
	}

	return b
}

func Init(filePath string, opts ...OptionFunc) {
	conf.Parse(filePath, "yaml", c, conf.ReadFromFile)

	for _, opt := range opts {
		opt(c)
	}

	var str bytes.Buffer
	_ = json.Indent(&str, c.JSON(), "", "    ")
	log.Printf("pudding trigger config:\n %s \n", str.String())
}

type OptionFunc func(config *Config)

func WithMySQLDSN(dsn string) OptionFunc {
	return func(config *Config) {
		if dsn != "" {
			config.MySQL.DSN = dsn
		}
	}
}

func WithWebhookPrefix(webhookPrefix string) OptionFunc {
	return func(config *Config) {
		if webhookPrefix != "" {
			config.WebhookPrefix = webhookPrefix
		}
	}
}
