package types

import "time"

// Message 消息
type Message struct {
	Topic     string        // Message Topic
	Key       string        // Message Key
	Body      []byte        // Message Body
	Partition int           // Message Partition
	Delay     time.Duration // Message Delay Time (Seconds)
	ReadyTime time.Time     // Message Ready Time（now + delay, Unix Timestamp, Seconds）
}
