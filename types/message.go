package types

import (
	"context"
	"encoding/json"
)

type HandleMessage func(ctx context.Context, msg *Message) error

// Message 消息
type Message struct {
	Topic     string // Message Topic
	Partition int    // Message Partition
	Key       string // Message Key
	Payload   []byte // Message Payload
	Delay     int64  // Message Delay Time (Seconds)
	ReadyTime int64  // Message Ready Time（now + delay, Unix Timestamp, Seconds）
}

func GetMessageFromJSON(j []byte) (*Message, error) {
	var m Message
	if err := json.Unmarshal(j, &m); err != nil {
		return nil, err
	}

	return &m, nil
}
