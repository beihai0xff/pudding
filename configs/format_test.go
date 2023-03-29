package configs

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnmarshalToStruct(t *testing.T) {
	Parse("../test/data/config.test.yaml", ConfigFormatYAML, ReadFromFile)
	type args struct {
		path string
		c    *BrokerConfig
	}
	tests := []struct {
		name       string
		args       args
		wantConfig *BrokerConfig
		wantErr    bool
	}{
		{
			name: "happy_path",
			args: args{"", &BrokerConfig{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := UnmarshalToStruct(tt.args.path, tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalToStruct() error = %v, wantErr %v", err, tt.wantErr)
			}
			assert.Equal(t, "redis", tt.args.c.ServerConfig.Broker)
			buf, _ := JSONFormat(tt.args.c)
			fmt.Println(buf.String())
		})
	}
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
            "grpc_port": 50051,
            "http_port": 8080,
            "cert_path": "",
            "key_path": "",
            "name_server_url": "",
            "log_config": [
                {
                    "log_name": "default",
                    "writers": [
                        "stdout"
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
                        "stdout"
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
        "connector": "kafka"
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
