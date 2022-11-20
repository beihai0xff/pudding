package pulsar_connector

import (
	"testing"

	"github.com/apache/pulsar-client-go/pulsar"
	"github.com/stretchr/testify/assert"

	"github.com/beihai0xff/pudding/types"
)

func Test_convertToPulsarProducerMessage(t *testing.T) {
	type args struct {
		msg *types.Message
	}
	tests := []struct {
		name string
		args args
		want *pulsar.ProducerMessage
	}{
		{"test", args{&types.Message{Payload: []byte("hello"), Key: "key"}}, &pulsar.ProducerMessage{Payload: []byte("hello"), Key: "key"}},
		{"test_key_nil", args{&types.Message{Payload: []byte("hello"), Key: ""}}, &pulsar.ProducerMessage{Payload: []byte("hello"), Key: ""}},
		{"test_payload_nil", args{&types.Message{Payload: nil, Key: "key"}}, &pulsar.ProducerMessage{Payload: nil, Key: "key"}},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, convertToPulsarProducerMessage(tt.args.msg))
	}
}
