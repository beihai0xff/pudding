package redis_broker

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v9"

	rdb "github.com/beihai0xff/pudding/pkg/redis"
	"github.com/beihai0xff/pudding/types"
)

type RealTimeQueue struct {
	rdb *rdb.Client // Redis客户端
}

// Produce produce a Message to the queue in realtime
func (q *RealTimeQueue) Produce(ctx context.Context, msg *types.Message) error {
	return q.rdb.StreamSend(ctx, q.getTopicPartition(msg.Topic), msg.Body)
}

// NewConsumer consume Messages from the queue in real time
func (q *RealTimeQueue) NewConsumer(topic, group, consumerName string, batchSize int,
	fn func(msg *types.Message) error) {

	for {
		// 拉取已经投递却未被 ACK 的消息，保证消息至少被成功消费1次
		msgs, err := q.rdb.XGroupConsume(context.Background(), topic, group, consumerName, "0", batchSize)
		if err != nil {
			// TODO: 记录错误日志
		}
		q.handlerRealTimeMessage(msgs, topic, group, fn)
		if len(msgs) == batchSize {
			// 如果一次未拉取完未被 ACK 的消息，则继续拉取
			// 确保先消费完成未被 ACK 的消息
			continue
		}

		// 拉取新消息
		msgs, err = q.rdb.XGroupConsume(context.Background(), topic, group, consumerName, ">", batchSize)
		if err != nil {
			// TODO: 记录错误日志
		}
		q.handlerRealTimeMessage(msgs, topic, group, fn)

	}
}

// handlerRealTimeMessage handle Messages from the queue in real time
func (q *RealTimeQueue) handlerRealTimeMessage(msgs []redis.XMessage, topic, group string,
	fn func(msg *types.Message) error) {

	// 遍历处理消息
	for _, msg := range msgs {

		// TODO: 消费超过三次的消息，记录错误日志，并添加到死信队列
		m, err := types.GetMessageFromJSON(msg.Values["body"].([]byte))
		if err != nil {
			// TODO: 记录错误日志
			continue
		}

		// handle message
		if err := fn(m); err != nil {
			// TODO: 记录错误日志
			continue
		}

		// handle message success，ACK
		if err := q.rdb.XAck(context.Background(), topic, group, msg.ID); err != nil {
			// TODO: 记录错误日志
			continue
		}
	}

	return
}

func (q *RealTimeQueue) getTopicPartition(topic string) string {
	return fmt.Sprintf("stream_%s", topic)
}

// Close close the queue
func (q *RealTimeQueue) Close() error {
	return q.rdb.Close()
}
