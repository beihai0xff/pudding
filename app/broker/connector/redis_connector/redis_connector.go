// Package redis_connector implements a connector with redis
package redis_connector

import (
	"context"
	"time"

	"github.com/go-redis/redis/v9"

	"github.com/beihai0xff/pudding/api/gen/pudding/types/v1"
	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/pkg/msgpack"
	rdb "github.com/beihai0xff/pudding/pkg/redis"
	type2 "github.com/beihai0xff/pudding/types"
)

type RealTimeQueue struct {
	rdb          *rdb.Client // Redis客户端
	consumerName string
}

// Produce produce a Message to the queue in realtime
func (q *RealTimeQueue) Produce(ctx context.Context, msg *types.Message) error {
	b, err := msgpack.Encode(msg)
	if err != nil {
		log.Errorf("msgpack encode message error: %v", err)
		return err
	}
	return q.rdb.StreamSend(ctx, msg.Topic, b)
}

// NewConsumer consume Messages from the queue in real time
func (q *RealTimeQueue) NewConsumer(topic, group string, batchSize int, fn type2.HandleMessage) error {
	go func() {
		for {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			// 拉取已经投递却未被 ACK 的消息，保证消息至少被成功消费1次
			msgs, err := q.rdb.XGroupConsume(ctx, topic, group, q.consumerName, "0", batchSize)
			if err != nil {
				log.Errorf("XGroupConsume unack message error: %v", err)
			}

			if len(msgs) == batchSize {
				// 如果一次未拉取完未被 ACK 的消息，则继续拉取
				// 确保先消费完成未被 ACK 的消息
				continue
			}

			// 拉取新消息
			msgs, err = q.rdb.XGroupConsume(ctx, topic, group, q.consumerName, ">", batchSize)
			if err != nil {
				log.Errorf("XGroupConsume message error: %v", err)
			}
			q.handlerRealTimeMessage(ctx, msgs, topic, group, fn)

			cancel()
		}
	}()

	return nil
}

// handlerRealTimeMessage handle Messages from the queue in real time
func (q *RealTimeQueue) handlerRealTimeMessage(ctx context.Context, msgs []redis.XMessage, topic, group string,
	fn type2.HandleMessage) {
	// 遍历处理消息
	for _, msg := range msgs {
		var m *types.Message
		// TODO: 消费超过三次的消息，记录错误日志，并添加到死信队列
		err := msgpack.Decode(msg.Values["body"].([]byte), m)
		if err != nil {
			// TODO: 记录错误日志
			continue
		}

		// handle message
		if err := fn(ctx, m); err != nil {
			// TODO: 记录错误日志
			continue
		}

		// handle message success，ACK
		if err := q.rdb.XAck(ctx, topic, group, msg.ID); err != nil {
			// TODO: 记录错误日志
			continue
		}
	}
}

// Close close the queue
func (q *RealTimeQueue) Close() error {
	return q.rdb.Close()
}
