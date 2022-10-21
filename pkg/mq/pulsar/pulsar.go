package pulsar

import (
	"context"
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

func (c *Client) Close() {
	for _, producer := range c.producers {
		producer.Close()
	}

	c.client.Close()
}
