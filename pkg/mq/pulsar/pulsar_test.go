package pulsar

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/stretchr/testify/assert"
)

var c *Client

func TestMain(m *testing.M) {
	c = NewMockPulsar()

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
		wantErr assert.ErrorAssertionFunc
	}{
		{"test_topic_1", args{context.Background(), "test_topic_1", &pulsar.ProducerMessage{Payload: []byte("hello"), Key: "key"}}, assert.NoError},
		{"test_topic_2", args{context.Background(), "test_topic_2", &pulsar.ProducerMessage{Payload: []byte("hello"), Key: "key"}}, assert.NoError},
		{"test_topic_not_exist", args{context.Background(), "test_not_exist", &pulsar.ProducerMessage{Payload: []byte("hello"), Key: "key"}}, assert.Error},
	}
	for _, tt := range tests {
		tt.wantErr(t, c.Produce(tt.args.ctx, tt.args.topic, tt.args.msg), fmt.Sprintf("ProducerMessage(%v)", tt.args.topic))

	}
}

func TestClient_NewConsumer(t *testing.T) {
	// produce test data
	TestClient_Produce(t)

	handle := HandleMessage(
		func(ctx context.Context, msg pulsar.Message) error {
			assert.Equal(t, "hello", string(msg.Payload()))
			return nil
		},
	)

	type args struct {
		topic string
		group string
		fn    HandleMessage
	}
	tests := []struct {
		name    string
		args    args
		wantErr assert.ErrorAssertionFunc
	}{
		{"test_topic_1", args{"test_topic_1", "test_group", handle}, assert.NoError},
		{"test_repeat_consumer", args{"test_topic_1", "test_group", handle}, assert.Error},
		{"test_topic_2", args{"test_topic_2", "test_group", handle}, assert.NoError},
		{"test_topic_not_exist", args{"test_not_exist", "test_group2", handle}, assert.NoError},
	}
	for _, tt := range tests {
		tt.wantErr(t, c.NewConsumer(tt.args.topic, tt.args.group, tt.args.fn), fmt.Sprintf("NewConsumer(%v, %v)", tt.args.topic, tt.args.group))
		time.Sleep(500 * time.Millisecond)
	}
}
