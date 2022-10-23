package pulsar

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/apache/pulsar-client-go/pulsar"

	"github.com/beihai0xff/pudding/pkg/configs"
	"github.com/beihai0xff/pudding/pkg/log"
)

var (
	clientOnce sync.Once
	client     *Client
)

type Client struct {
	client    pulsar.Client
	producers map[string]pulsar.Producer
	consumers map[string]pulsar.Consumer
}

func New(config *configs.PulsarConfig) *Client {
	clientOnce.Do(
		func() {
			c, err := pulsar.NewClient(pulsar.ClientOptions{
				URL: config.PulsarURL,
			})

			if err != nil {
				log.Errorf("create pulsar client failed: %v", err)
				panic(err)
			}

			client = &Client{
				client:    c,
				producers: make(map[string]pulsar.Producer),
				consumers: make(map[string]pulsar.Consumer),
			}

			for _, pc := range config.ProducersConfig {
				producer, err := c.CreateProducer(pc)
				if err != nil {
					log.Errorf("create pulsar Producer %s failed: %v", pc.Topic, err)
					panic(err)
				}
				client.producers[pc.Topic] = producer
			}
		})

	return client
}

func (c *Client) Produce(ctx context.Context, topic string, msg *pulsar.ProducerMessage) error {
	_, err := c.producers[topic].Send(ctx, msg)
	return err
}

func (c *Client) NewConsumer(topic, group string, fn func(msg pulsar.Message) error) error {
	consumerName := c.getConsumerName(topic, group)

	// check consumer exist
	// if consumer already exists, return error
	if c.consumers[topic] != nil {
		return fmt.Errorf("consumer for group [%s] topic [%s]  already exists", group, topic)
	}

	// if consumer not exists, create a new consumer
	consumer, err := c.client.Subscribe(pulsar.ConsumerOptions{
		Topic:            topic,
		SubscriptionName: group,
		Type:             pulsar.Shared,
		Name:             consumerName,
	})
	if err != nil {
		return fmt.Errorf("create pulsar Consumer %s failed: %v", group, err)
	}

	// wrap the ack function
	ack := func(msg pulsar.Message) {
		if err := consumer.Ack(msg); err != nil {
			log.Errorf("ack message failed: %v, message msgId: %#v -- content: '%s'\n",
				err, msg.ID(), string(msg.Payload()))
		}
	}

	go func() {
		for {
			// receive message, block if queue is empty
			msg, err := consumer.Receive(context.Background())
			if err != nil {
				log.Errorf("receive message failed: %v, message msgId: %#v -- content: '%s'\n",
					err, msg.ID(), string(msg.Payload()))
			}

			if msg.RedeliveryCount() > 3 {
				log.Errorf("message redelivery count exceed 3, message msgId: %#v -- content: '%s'\n",
					msg.ID(), string(msg.Payload()))
				ack(msg)
				continue
			}

			log.Debugf("Received message msgId: %#v -- content: '%s'\n",
				msg.ID(), string(msg.Payload()))

			if err := fn(msg); err != nil {
				log.Errorf("handle message failed: %v, message msgId: %#v -- content: '%s'\n",
					err, msg.ID(), string(msg.Payload()))
				continue
			}

			// Acknowledge successful processing of the message
			ack(msg)
		}
	}()

	// when consumer created, add it to the map
	c.consumers[topic] = consumer
	return nil
}

func (c *Client) getConsumerName(topic, group string) string {
	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = "unknown"
	}
	return fmt.Sprintf("%s-%s-%s", topic, group, hostname)
}

func (c *Client) Close() {
	for _, producer := range c.producers {
		producer.Close()
	}

	c.client.Close()
}
