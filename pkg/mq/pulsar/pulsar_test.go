package pulsar

import (
	"context"
	"os"
	"testing"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/stretchr/testify/assert"

	"github.com/beihai0xff/pudding/pkg/configs"
)

var c *Client

func TestMain(m *testing.M) {
	config := &configs.PulsarConfig{
		PulsarURL:   "pulsar://localhost:6650",
		DialTimeout: 10,
		ProducersConfig: []pulsar.ProducerOptions{
			{Topic: "test_topic_1"},
			{Topic: "test_topic_2"},
		},
	}
	c = New(config)
	exitCode := m.Run()
	// 退出
	os.Exit(exitCode)
}

func TestClient_getConsumerName(t *testing.T) {
	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = "unknown"
	}

	type args struct {
		topic string
		group string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"test", args{"test_topic_1", "test_group"}, "test_topic_1-test_group-" + hostname},
		{"test", args{"test_topic_2", "test_group"}, "test_topic_2-test_group-" + hostname},
		{"test-topic-and-group-empty", args{"", ""}, "--" + hostname},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, c.getConsumerName(tt.args.topic, tt.args.group))
	}
}

func TestClient_Produce(t *testing.T) {
	type args struct {
		ctx   context.Context
		topic string
		msg   *pulsar.ProducerMessage
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"test_topic_1", args{context.Background(), "test_topic_1", &pulsar.ProducerMessage{Payload: []byte("hello"), Key: "key"}}, false},
		{"test_topic_2", args{context.Background(), "test_topic_2", &pulsar.ProducerMessage{Payload: []byte("hello"), Key: "key"}}, false},
		{"test_not_exist", args{context.Background(), "test_not_exist", &pulsar.ProducerMessage{Payload: []byte("hello"), Key: "key"}}, true},
	}
	for _, tt := range tests {
		err := c.Produce(tt.args.ctx, tt.args.topic, tt.args.msg)

		assert.Equal(t, tt.wantErr, err != nil)

	}
}
