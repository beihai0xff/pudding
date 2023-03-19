// Package kafka implements a kafka client.
package kafka

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"

	"github.com/beihai0xff/pudding/configs"
	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/pkg/log/logger"
	"github.com/beihai0xff/pudding/pkg/utils"
)

// Message kafka message body
type Message = kafka.Message

var (
	kafkaClientOnce sync.Once
	kafkaClient     *client
)

// Client kafka client interface
type Client interface {
	SendMessage(ctx context.Context, msg *kafka.Message) (string, error)
	NewConsumer(ctx context.Context, topic, group string,
		handle func(context.Context, *Message) error) (Consumer, error)
	Close() error
}

// client kafka client
type client struct {
	*configs.KafkaConfig
	mutex  sync.Mutex
	writer *kafka.Writer
	logger *logger.MessageLogger
}

// NewClient create a kafka client
func NewClient(config *configs.KafkaConfig) Client {
	kafkaClientOnce.Do(func() {
		l := logger.NewMessageLogger()
		kafkaClient = &client{
			KafkaConfig: config,
			logger:      l,
			writer: &kafka.Writer{
				Addr: kafka.TCP(getAddress(config.Host, config.Port)),
				// the same key will be sent to the same partition
				Balancer: &kafka.CRC32Balancer{},
				// the minimum amount of time to wait before sending a batch of messages
				BatchTimeout: time.Duration(config.ProducerBatchTimeout) * time.Millisecond,
				BatchSize:    config.BatchSize,
				Logger:       kafka.LoggerFunc(l.RecordMessageInfoLog),
				ErrorLogger:  kafka.LoggerFunc(l.RecordMessageErrorLog),
			},
		}
	})
	return kafkaClient
}

// SendMessage send kafka message
func (c *client) SendMessage(ctx context.Context, msg *kafka.Message) (string, error) {
	if err := c.writer.WriteMessages(ctx, *msg); err != nil {
		return "", errors.Wrapf(err, "failed to write messages, topic=%s, host=%s, port=%d",
			msg.Topic, c.Host, c.Port)
	}

	return buildKafkaMsgID(msg), nil
}

// NewConsumer create a new consumer
func (c *client) NewConsumer(ctx context.Context, topic, group string,
	handle func(context.Context, *Message) error) (Consumer, error) {
	reader := kafka.NewReader(*c.getReaderConfig(topic, group, c.KafkaConfig))

	kafkaConsumer := &consumer{
		reader: reader,
		name:   fmt.Sprintf("%s-%s-%s-%s", topic, group, utils.GetOutBoundIP(), uuid.NewString()),
		mutex:  &sync.Mutex{},
		logger: c.logger,
		handle: handle,
	}

	return kafkaConsumer, nil
}

// Close close kafka client
func (c *client) Close() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.writer.Close()
}

func getAddress(host string, port int) string {
	return fmt.Sprintf("%s:%d", host, port)
}

func (c *client) getReaderConfig(topic, group string, config *configs.KafkaConfig) *kafka.ReaderConfig {
	return &kafka.ReaderConfig{
		Brokers:  []string{getAddress(config.Host, config.Port)},
		Topic:    topic,
		GroupID:  group,
		MinBytes: 1,
		MaxBytes: 10e6,
		MaxWait:  time.Duration(config.ConsumerMaxWaitTime) * time.Millisecond,
		// if the broker has no offset for the consumer group, start with the FirstOffset
		StartOffset: kafka.FirstOffset,
		Logger:      kafka.LoggerFunc(c.logger.RecordMessageInfoLog),
		ErrorLogger: kafka.LoggerFunc(c.logger.RecordMessageErrorLog),
	}
}

// Consumer kafka consumer interface
type Consumer interface {
	Close() error
	Run(ctx context.Context)
}

// consumer kafka consumer
type consumer struct {
	name string
	// mutex is used to protect the isClosed field
	mutex    *sync.Mutex
	reader   *kafka.Reader
	logger   *logger.MessageLogger
	isClosed bool
	handle   func(context.Context, *Message) error
}

// Run start a goroutine to consume kafka message
func (c *consumer) Run(ctx context.Context) {
	go c.worker(ctx)
}

// worker start a goroutine to consume kafka message
func (c *consumer) worker(ctx context.Context) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	for {
		if c.isClosed {
			break
		}
		// read kafka message in blocking mode
		msg, err := c.reader.FetchMessage(ctx)
		if err != nil {
			// if reader has been closed, break the loop
			if err == io.EOF {
				break
			}
			log.Errorf("Failed to read kafka msg, caused by %s", err)
			time.Sleep(time.Second)
			continue
		}

		c.logger.RecordMessageInfoLog("Success to read kafka msg: [%s] \n"+
			"\tTopic: %s\n\tpartition %d: %d",
			msg.Value, msg.Topic, msg.Partition, msg.Offset)

		// process kafka Message
		c.handleMsg(&msg)

		// commit kafka Message
		c.commitMsg(&msg)
	}
}

// handleMsg process kafka message
func (c *consumer) handleMsg(msg *kafka.Message) {
	if err := c.handle(context.Background(), msg); err != nil {
		log.Errorf("Failed to handle kafka msg: [%s] , caused by %s \n"+
			"\tTopic: %s\n\tpartition %d: %d",
			msg.Value, err, msg.Topic, msg.Partition, msg.Offset)
		return
	}
}

// handleMsg process kafka message
func (c *consumer) commitMsg(msg *kafka.Message) {
	ctx, Cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer func() {
		Cancel()
	}()

	// commit 消息
	if err := c.reader.CommitMessages(ctx, *msg); err != nil {
		// commit 失败
		log.Errorf("Failed to commit kafka msg: [%s], caused by %s \n"+
			"\tTopic: %s\n\tpartition %d: %d",
			msg.Value, err, msg.Topic, msg.Partition, msg.Offset)
		return
	}
}

// Close close kafka consumer
// we wrap the reader.Close() to make it compatible with kafka.Close()
func (c *consumer) Close() error {
	c.isClosed = true
	c.mutex.Lock()
	if err := c.reader.Close(); err != nil {
		return err
	}
	log.Infof("%s reader Closed", c.name)
	// wait for the worker() goroutine to exit

	defer c.mutex.Unlock()
	return nil
}
