// Package kafka_connector provides a connector to Kafka
package kafka_connector

import (
	"context"

	"github.com/google/uuid"

	"github.com/beihai0xff/pudding/api/gen/pudding/types/v1"
	type2 "github.com/beihai0xff/pudding/app/broker/pkg/types"
	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/pkg/mq/kafka"
)

// Connector is a kafka connector
type Connector struct {
	client    kafka.Client
	consumers map[string]kafka.Consumer
}

// NewConnector create a new kafka connector
func NewConnector(client kafka.Client) *Connector {
	return &Connector{
		client:    client,
		consumers: map[string]kafka.Consumer{},
	}
}

// Produce produce a Message to the queue in realtime
func (c *Connector) Produce(ctx context.Context, msg *types.Message) error {
	_, err := c.client.SendMessage(ctx, &kafka.Message{
		Topic: msg.Topic,
		Key:   []byte(msg.Key),
		Value: msg.Payload,
	})

	if err != nil {
		log.Errorf("SendMessage [%s] error: %v", msg.Key, err)
		return err
	}

	return nil
}

// NewConsumer consume Messages from the queue in real time
func (c *Connector) NewConsumer(topic, group string, batchSize int, fn type2.HandleMessage) error {
	ctx := context.Background()
	handler := func(ctx context.Context, kafkaMessage *kafka.Message) error {
		msg := types.Message{
			Topic:   kafkaMessage.Topic,
			Key:     string(kafkaMessage.Key),
			Payload: kafkaMessage.Value,
		}

		return fn(ctx, &msg)
	}

	consumer, err := c.client.NewConsumer(ctx, topic, group, handler)
	if err != nil {
		log.Errorf("NewConsumer error: %v", err)
		return err
	}

	consumer.Run(ctx)

	c.consumers[topic+group+uuid.NewString()] = consumer

	return nil
}

// Close the connector
func (c *Connector) Close() error {
	for consumerName, consumer := range c.consumers {
		if err := consumer.Close(); err != nil {
			log.Errorf("consumer [%s] close error: %v", consumerName, err)
		}
	}

	return c.client.Close()
}
