package pulsar

import (
	"context"
	"fmt"
	"os"

	"github.com/apache/pulsar-client-go/pulsar"

	"github.com/beihai0xff/pudding/pkg/configs"
	"github.com/beihai0xff/pudding/pkg/log"
)

// HandleMessage is the function type for handling message
type HandleMessage func(ctx context.Context, msg pulsar.Message) error

type Client struct {
	client    pulsar.Client
	producers map[string]pulsar.Producer
	consumers map[string]pulsar.Consumer
}

func New(config *configs.PulsarConfig) *Client {
	// create pulsar client
	c, err := pulsar.NewClient(pulsar.ClientOptions{
		URL: config.PulsarURL,
	})

	if err != nil {
		log.Errorf("create pulsar client failed: %v", err)
		panic(err)
	}

	// create wrapper client
	client := &Client{
		client:    c,
		producers: make(map[string]pulsar.Producer),
		consumers: make(map[string]pulsar.Consumer),
	}

	// create producers
	for _, pc := range config.ProducersConfig {
		producer, err := c.CreateProducer(pc)
		if err != nil {
			log.Errorf("create pulsar Producer %s failed: %v", pc.Topic, err)
			panic(err)
		}
		client.producers[pc.Topic] = producer
	}

	return client
}

func (c *Client) Produce(ctx context.Context, topic string, msg *pulsar.ProducerMessage) error {
	produce, ok := c.producers[topic]
	if !ok {
		return fmt.Errorf("producer for topic [%s] not exists", topic)
	}
	_, err := produce.Send(ctx, msg)
	return err
}

func (c *Client) NewConsumer(topic, group string, fn HandleMessage) error {
	if topic == "" || group == "" {
		return fmt.Errorf("topic and group can not be empty")
	}
	consumerName := c.getConsumerName(topic, group)

	// check consumer exist
	// if consumer already exists, return error
	if _, ok := c.consumers[topic]; ok {
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

	defer func() {
		// when consumer created, add it to the map
		c.consumers[topic] = consumer
	}()

	// wrap the ack function
	ack := func(msg pulsar.Message) {
		if err := consumer.Ack(msg); err != nil {
			log.Errorf("ack message failed: %v, message msgId: %#v -- content: '%s'\n",
				err, msg.ID(), string(msg.Payload()))
		}
	}

	go func() {
		for {

			ctx := context.Background()
			// receive message, block if queue is empty
			msg, err := consumer.Receive(ctx)
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

			if err := fn(ctx, msg); err != nil {
				log.Errorf("handle message failed: %v, message msgId: %#v -- content: '%s'\n",
					err, msg.ID(), string(msg.Payload()))
				continue
			}

			// Acknowledge successful processing of the message
			ack(msg)
		}
	}()

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
	//  close all producers
	for _, producer := range c.producers {
		producer.Flush()
		producer.Close()
	}
	// close all consumers
	for _, consumer := range c.consumers {
		consumer.Close()
	}

	// close the client
	c.client.Close()
}
