package types

import "time"

// Message 消息
type Message struct {
	Topic     string        // Message Topic
	Partition int           // Message Partition
	Key       string        // Message Key
	Body      []byte        // Message Body
	Delay     time.Duration // Message Delay Time (Seconds)
	ReadyTime time.Time     // Message Ready Time（now + delay, Unix Timestamp, Seconds）
}
