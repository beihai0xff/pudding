// Package redis_connector implements a connector with redis
package redis_connector

import (
	"context"

	"github.com/redis/go-redis/v9"

	"github.com/beihai0xff/pudding/api/gen/pudding/types/v1"
	type2 "github.com/beihai0xff/pudding/app/broker/pkg/types"
	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/pkg/msgpack"
	rdb "github.com/beihai0xff/pudding/pkg/redis"
)

// Connector is a redis connector
type Connector struct {
	rdb          *rdb.Client // Redis客户端
	consumerName string
}

// NewConnector create a new redis connector
func NewConnector(client *rdb.Client) *Connector {
	return &Connector{
		rdb:          client,
		consumerName: "pudding",
	}
}

// Produce produce a Message to the queue in realtime
func (q *Connector) Produce(ctx context.Context, msg *types.Message) error {
	b, err := msgpack.Encode(msg)
	if err != nil {
		log.Errorf("msgpack encode message error: %v", err)

		return err
	}

	return q.rdb.StreamSend(ctx, msg.Topic, b)
}

// NewConsumer consume Messages from the queue in real time
func (q *Connector) NewConsumer(topic, group string, batchSize int, fn type2.HandleMessage) error {
	go func() {
		for {
			ctx := context.Background()
			// pull unack message, ensure message is consumed at least once
			msgs, err := q.rdb.XGroupConsume(ctx, topic, group, q.consumerName, "0", batchSize)
			if err != nil {
				log.Errorf("XGroupConsume unack message error: %v", err)
			}

			if len(msgs) == batchSize {
				// if unack message is not consumed, continue to consume
				continue
			}

			// pull new messages
			msgs, err = q.rdb.XGroupConsume(ctx, topic, group, q.consumerName, ">", batchSize)
			if err != nil {
				log.Errorf("XGroupConsume message error: %v", err)
			}

			q.handlerRealTimeMessage(ctx, msgs, topic, group, fn)
		}
	}()

	return nil
}

// handlerRealTimeMessage handle Messages from the queue in real time
func (q *Connector) handlerRealTimeMessage(ctx context.Context, msgs []redis.XMessage, topic, group string,
	fn type2.HandleMessage) {
	for _, msg := range msgs {
		var m *types.Message
		if err := msgpack.Decode(msg.Values["body"].([]byte), m); err != nil {
			log.Errorf("decode message error: %v", err)
			continue
		}

		// handle message
		if err := fn(ctx, m); err != nil {
			log.Errorf("handle message error: %v", err)
			continue
		}

		// handle message success，ACK
		if err := q.rdb.XAck(ctx, topic, group, msg.ID); err != nil {
			log.Errorf("XAck message error: %v", err)
			continue
		}
	}
}

// Close close the queue
func (q *Connector) Close() error {
	return q.rdb.Close()
}
