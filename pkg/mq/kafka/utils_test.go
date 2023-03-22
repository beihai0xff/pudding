package kafka

import (
	"testing"

	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
)

func Test_buildKafkaMsgID(t *testing.T) {
	type args struct {
		message *kafka.Message
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"happy-path", args{&Message{Topic: "test", Partition: 1, Offset: 2}}, "topic:test-partition:1-offset:2"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, buildKafkaMsgID(tt.args.message), "they should be equal")
		})
	}
}
