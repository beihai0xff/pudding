package configs

import (
	"fmt"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

type formatConfig struct {
	Name         string `json:"name" yaml:"name" mapstructure:"name"`
	ServerConfig struct {
		Broker     string `json:"broker" yaml:"broker" mapstructure:"broker"`
		BaseConfig struct {
			HostDomain string      `json:"host_domain" yaml:"host_domain" mapstructure:"host_domain"`
			GRPCPort   int         `json:"grpc_port" yaml:"grpc_port" mapstructure:"grpc_port"`
			HTTPPort   int         `json:"http_port" yaml:"http_port" mapstructure:"http_port"`
			EnableTLS  bool        `json:"enable_tls" yaml:"enable_tls" mapstructure:"enable_tls"`
			Logger     []LogConfig `json:"log_config" yaml:"log_config" mapstructure:"log_config"`
		} `json:"base_config" yaml:"base_config" mapstructure:"base_config"`
	} `json:"server_config" yaml:"server_config"  mapstructure:"server_config"`
}

func TestUnmarshalToStruct(t *testing.T) {
	assert.NoError(t, Parse("../test/data/config.format.yaml", ConfigFormatYAML, ReadFromFile))

	// happy_path
	format := formatConfig{}
	err := UnmarshalToStruct("", &format)
	assert.NoError(t, err)
	assert.Equal(t, "redis", format.ServerConfig.Broker)
	buf, _ := JSONFormat(format)
	assert.Equal(t, testJSONFormat, buf.String())

	// get_logger_configs
	var logConfig []LogConfig
	err = UnmarshalToStruct("server_config.base_config.log_config", &logConfig)
	assert.NoError(t, err)
	v, _ := lo.Find(logConfig, func(conf LogConfig) bool {
		return conf.LogName == "default"
	})
	assert.Equal(t, LogConfig{LogName: "default", Writers: []string{OutputConsole}, Level: "debug", Format: EncoderTypeJSON}, v)
}

func TestJSONFormat(t *testing.T) {
	assert.NoError(t, Parse("../test/data/config.test.yaml", ConfigFormatYAML, ReadFromFile))
	config := BrokerConfig{}
	_ = UnmarshalToStruct("", &config)

	tests := []struct {
		name    string
		args    interface{}
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{"happy_path", &config, testJSONFormat, assert.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := JSONFormat(tt.args)
			if !tt.wantErr(t, err, fmt.Sprintf("JSONFormat(%v)", tt.args)) {
				return
			}
			assert.Equalf(t, tt.want, got.String(), "JSONFormat(%v)", tt.args)
		})
	}
}

var testJSONFormat = `{
    "name": "pudding",
    "server_config": {
        "broker": "redis",
        "base_config": {
            "host_domain": "localhost",
            "grpc_port": 50051,
            "http_port": 8080,
            "enable_tls": false,
            "log_config": [
                {
                    "log_name": "default",
                    "writers": [
                        "console"
                    ],
                    "file_config": {
                        "filepath": "",
                        "max_age": 0,
                        "max_backups": 0,
                        "compress": false,
                        "max_size": 0
                    },
                    "format": "json",
                    "level": "debug",
                    "caller_skip": 0
                },
                {
                    "log_name": "kafka_log",
                    "writers": [
                        "console"
                    ],
                    "file_config": {
                        "filepath": "",
                        "max_age": 0,
                        "max_backups": 0,
                        "compress": false,
                        "max_size": 0
                    },
                    "format": "json",
                    "level": "debug",
                    "caller_skip": 0
                }
            ]
        }
    }
}`
