package configs

import (
	"fmt"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestUnmarshalToStruct(t *testing.T) {
	Parse("../test/data/config.test.yaml", ConfigFormatYAML, ReadFromFile)

	// happy_path
	brokerConfig := BrokerConfig{}
	err := UnmarshalToStruct("", &brokerConfig)
	assert.NoError(t, err)
	assert.Equal(t, "redis", brokerConfig.ServerConfig.Broker)
	buf, _ := JSONFormat(brokerConfig)
	assert.Equal(t, testJSONFormat, buf.String())

	// get_logger_configs
	var logConfig []LogConfig
	err = UnmarshalToStruct("server_config.base_config.log_config", &logConfig)
	assert.NoError(t, err)
	v, _ := lo.Find(logConfig, func(conf LogConfig) bool {
		return conf.LogName == "default"
	})
	assert.Equal(t, LogConfig{LogName: "default", Writers: []string{OutputConsole}, Level: "debug", Format: EncoderTypeJSON}, v)
	fmt.Println(logConfig)
}

func TestJSONFormat(t *testing.T) {
	Parse("../test/data/config.test.yaml", ConfigFormatYAML, ReadFromFile)
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
    "server_config": {
        "base_config": {
            "host_domain": "localhost",
            "grpc_port": 50051,
            "http_port": 8080,
            "enable_tls": false,
            "cert_path": "",
            "key_path": "",
            "name_server_url": "",
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
        },
        "time_slice_interval": "",
        "message_topic": "",
        "token_topic": "",
        "broker": "redis",
        "connector": "kafka",
        "etcd_urls": null
    },
    "redis_config": {
        "url": "redis://default:default@localhost:6379/11",
        "dial_timeout": 20
    },
    "pulsar_config": {
        "url": "pulsar://localhost:6650",
        "connection_timeout": 10,
        "producers_config": [
            {
                "topic": "token",
                "batching_max_publish_delay": 10,
                "batching_max_messages": 100,
                "batching_max_size": 5
            },
            {
                "topic": "test_topic_1",
                "batching_max_publish_delay": 10,
                "batching_max_messages": 100,
                "batching_max_size": 5
            },
            {
                "topic": "test_topic_2",
                "batching_max_publish_delay": 10,
                "batching_max_messages": 100,
                "batching_max_size": 5
            }
        ]
    },
    "kafka_config": null
}`
