package kafka

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"

	"github.com/beihai0xff/pudding/configs"
	"github.com/beihai0xff/pudding/pkg/log"
	test_utils "github.com/beihai0xff/pudding/test/mock/utils"
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestNew(t *testing.T) {
	type args struct {
		config *configs.KafkaConfig
	}
	tests := []struct {
		name string
		args args
		want Client
	}{
		{"happy-path", args{test_utils.TestKafkaConfig}, Client(&client{})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.IsType(t, tt.want, New(tt.args.config), "New(%v)", tt.args.config)
		})
	}
}

func Test_newClient(t *testing.T) {
	type args struct {
		config *configs.KafkaConfig
	}
	tests := []struct {
		name     string
		args     args
		wantType *client
	}{
		{"happy-path", args{test_utils.TestKafkaConfig}, &client{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newClient(tt.args.config)
			assert.IsType(t, tt.wantType, c)
			c.Close()
		})
	}
}

func Test_client_SendMessage(t *testing.T) {
	c := newClient(test_utils.TestKafkaConfig)
	defer c.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	c.CreateTopic(ctx, "test-SendMessage-topic", 1, 1)

	type args struct {
		ctx context.Context
		msg *Message
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{"happy-path", args{context.Background(), &Message{
			Topic: "test-SendMessage-topic",
			Key:   []byte("test-key"),
			Value: []byte("test"),
		},
		}, "topic:test-SendMessage-topic-partition:0-offset:0", assert.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.SendMessage(tt.args.ctx, tt.args.msg)
			if !tt.wantErr(t, err, fmt.Sprintf("SendMessage(%v, %v)", tt.args.ctx, tt.args.msg)) {
				return
			}
			assert.Equalf(t, tt.want, got, "SendMessage(%v, %v)", tt.args.ctx, tt.args.msg)
		})
	}
}

func Test_client_getReaderConfig(t *testing.T) {
	c := newClient(test_utils.TestKafkaConfig)
	defer c.Close()
	type args struct {
		topic  string
		group  string
		config *configs.KafkaConfig
	}
	tests := []struct {
		name string
		args args
		want *kafka.ReaderConfig
	}{
		{
			name: "happy-path",
			args: args{
				topic:  "test-topic",
				group:  "test-consumer-group",
				config: test_utils.TestKafkaConfig,
			},
			want: &kafka.ReaderConfig{
				Brokers:          c.Address,
				Topic:            "test-topic",
				GroupID:          "test-consumer-group",
				MinBytes:         1,
				MaxBytes:         10 * 1024,
				MaxWait:          time.Duration(test_utils.TestKafkaConfig.ConsumerMaxWaitTime) * time.Millisecond,
				StartOffset:      kafka.FirstOffset,
				ReadBatchTimeout: 10 * time.Second,
				CommitInterval:   0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			readerConfig := c.getReaderConfig(tt.args.topic, tt.args.group, tt.args.config)
			assert.Equal(t, tt.want.Brokers, readerConfig.Brokers)
			assert.Equal(t, tt.want.Topic, readerConfig.Topic)
			assert.Equal(t, tt.want.GroupID, readerConfig.GroupID)
			assert.Equal(t, tt.want.MaxBytes, readerConfig.MaxBytes)
			assert.Equal(t, tt.want.MaxWait, readerConfig.MaxWait)
			assert.Equal(t, tt.want.StartOffset, readerConfig.StartOffset)

			assert.Equal(t, tt.want.MinBytes, readerConfig.MinBytes)
			assert.Equal(t, tt.want.ReadBatchTimeout, readerConfig.ReadBatchTimeout)
			assert.Equal(t, tt.want.CommitInterval, readerConfig.CommitInterval)
		})
	}
}

func Test_client_Close(t *testing.T) {
	c := newClient(test_utils.TestKafkaConfig)
	defer c.Close()
	tests := []struct {
		name    string
		wantErr assert.ErrorAssertionFunc
	}{
		{"close non-closed-client", assert.NoError},
		{"close closed-client", assert.NoError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.wantErr(t, c.Close(), fmt.Sprintf("Close()"))
		})
	}
}

func Test_client_NewConsumer(t *testing.T) {
	c := newClient(test_utils.TestKafkaConfig)
	defer c.Close()

	type args struct {
		ctx    context.Context
		topic  string
		group  string
		handle func(context.Context, *Message) error
	}
	tests := []struct {
		name    string
		args    args
		want    Consumer
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "happy-path",
			args: args{
				ctx:    context.Background(),
				topic:  "test-consumer-topic",
				group:  "test-consumer-group",
				handle: func(ctx context.Context, msg *Message) error { return nil },
			},
			want:    Consumer(&consumer{}),
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			consumer, err := c.NewConsumer(tt.args.ctx, tt.args.topic, tt.args.group, tt.args.handle)
			if !tt.wantErr(t, err) {
			}
			assert.IsType(t, tt.want, consumer)
		})
	}
}

var count = 0

func Test_consumer_Run(t *testing.T) {
	c := newClient(test_utils.TestKafkaConfig)
	defer c.Close()

	ctx := context.Background()
	c.CreateTopic(ctx, "test-consumer-topic", 1, 1)

	handler := func(ctx context.Context, msg *Message) error {
		log.Infof("success handle msg: %v", msg)
		assert.Equal(t, msg.Value, []byte("test-consumer-value"))
		log.Infof("%d", count)
		count++
		return nil
	}
	consumerInst, _ := c.NewConsumer(ctx, "test-consumer-topic", "test-consumer-group", handler)
	defer consumerInst.Close()

	for i := 0; i < 10; i++ {
		c.SendMessage(ctx, &Message{
			Topic: "test-consumer-topic",
			Value: []byte("test-consumer-value"),
		})
	}

	consumerInst.Run(ctx)

	time.Sleep(5 * time.Second)
	assert.Equal(t, 10, count)
	consumerInst.Close()
}
