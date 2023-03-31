// Package kafka implements a kafka client.
package kafka

import (
	"context"
	"fmt"
	"io"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
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

// Client kafka client interface
type Client interface {
	SendMessage(ctx context.Context, msg *Message) (string, error)
	NewConsumer(ctx context.Context, topic, group string,
		handler func(context.Context, *Message) error) (Consumer, error)
	CreateTopic(ctx context.Context, topic string, numPartitions, replicationFactor int) error
	Close() error
}

// client kafka client
type client struct {
	*configs.KafkaConfig
	mutex  sync.Mutex
	writer *kafka.Writer
	logger *logger.MessageLogger
}

// New create a kafka client
func New(conf *configs.KafkaConfig) Client {
	return newClient(conf)
}

func newClient(conf *configs.KafkaConfig) *client {
	l := logger.NewMessageLogger()
	return &client{
		KafkaConfig: conf,
		logger:      l,
		mutex:       sync.Mutex{},
		writer: &kafka.Writer{
			Addr: kafka.TCP(conf.Address...),
			// the same key will be sent to the same partition
			Balancer: &kafka.CRC32Balancer{},
			// the minimum amount of time to wait before sending a batch of messages
			BatchTimeout: time.Duration(conf.ProducerBatchTimeout) * time.Millisecond,
			BatchSize:    conf.BatchSize,
			RequiredAcks: kafka.RequireAll,
			Async:        false,
			Logger:       kafka.LoggerFunc(l.RecordMessageInfoLog),
			ErrorLogger:  kafka.LoggerFunc(l.RecordMessageErrorLog),
		},
	}
}

// SendMessage send kafka message
func (c *client) SendMessage(ctx context.Context, msg *Message) (string, error) {
	const retries = 3
	var err error
	for i := 0; i < retries; i++ {
		if err = c.writer.WriteMessages(ctx, *msg); err != nil {
			if errors.Is(err, kafka.LeaderNotAvailable) || errors.Is(err, context.DeadlineExceeded) {
				time.Sleep(time.Millisecond * 250)
				continue
			}
		} else {
			log.Debugf("send message [%s] to kafka success for times %d: %v", msg.Key, i, err)
			break
		}
		log.Warnf("send message [%s] to kafka failed for times %d: %v", msg.Key, i, err)
	}

	if err != nil {
		return "", errors.Wrapf(err, "failed to write messages, topic=%s, address=%s", msg.Topic, c.Address)
	}

	log.Debugf("send message to kafka success: %s", msg)
	return buildKafkaMsgID(msg), nil
}

// NewConsumer create a new consumer
func (c *client) NewConsumer(ctx context.Context, topic, group string,
	handler func(context.Context, *Message) error) (Consumer, error) {
	reader := kafka.NewReader(*c.getReaderConfig(topic, group, c.KafkaConfig))

	kafkaConsumer := &consumer{
		reader:  reader,
		name:    fmt.Sprintf("%s-%s-%s-%s", topic, group, utils.GetOutBoundIP(), uuid.NewString()),
		wg:      sync.WaitGroup{},
		logger:  c.logger,
		handler: handler,
	}

	return kafkaConsumer, nil
}

func (c *client) getReaderConfig(topic, group string, config *configs.KafkaConfig) *kafka.ReaderConfig {
	return &kafka.ReaderConfig{
		Brokers:  c.Address,
		Topic:    topic,
		GroupID:  group,
		MinBytes: 1,
		MaxBytes: 10 * 1024,
		MaxWait:  time.Duration(config.ConsumerMaxWaitTime) * time.Millisecond,
		// if the broker has no offset for the consumer group, start with the LastOffset
		StartOffset:      kafka.FirstOffset,
		ReadBatchTimeout: 10 * time.Second,
		Logger:           kafka.LoggerFunc(c.logger.RecordMessageInfoLog),
		ErrorLogger:      kafka.LoggerFunc(c.logger.RecordMessageErrorLog),
	}
}

// CreateTopic close kafka client
func (c *client) CreateTopic(ctx context.Context, topic string, numPartitions, replicationFactor int) error {
	if c.Address == nil || len(c.Address) == 0 {
		return errors.New("kafka address is empty")
	}

	conn, err := kafka.DialContext(ctx, c.Network, c.Address[0])
	if err != nil {
		log.Errorf("failed to dial kafka, err=%v", err)
		return err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		log.Errorf("failed to get controller, err=%v", err)
		return err
	}
	controllerConn, err := kafka.DialContext(ctx, c.Network,
		net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		log.Errorf("failed to dial controller, err=%v", err)
		return err
	}
	defer controllerConn.Close()

	topicConfig := kafka.TopicConfig{Topic: topic, NumPartitions: numPartitions, ReplicationFactor: replicationFactor}
	if err = controllerConn.CreateTopics(topicConfig); err != nil {
		log.Errorf("failed to create topic, err=%v", err)
		return err
	}

	// need wait for a while to make sure the topic is created
	c.waitForTopic(ctx, topic)
	log.Infof("create topic [%s] success", topic)
	return nil
}

// Block until topic exists.
func (c *client) waitForTopic(ctx context.Context, topic string) {
	for {
		select {
		case <-ctx.Done():
			log.Errorf("reached deadline before verifying topic existence")
		default:
		}

		cli := &kafka.Client{
			Addr:    c.writer.Addr,
			Timeout: 5 * time.Second,
		}

		response, err := cli.Metadata(ctx, &kafka.MetadataRequest{
			Addr:   cli.Addr,
			Topics: []string{topic},
		})
		if err != nil {
			log.Errorf("waitForTopic: error listing topics: %v", err)
		}

		// Find a topic which has at least 1 partition in the metadata response
		for _, top := range response.Topics {
			if top.Name != topic {
				continue
			}

			numPartitions := len(top.Partitions)
			log.Debugf("waitForTopic: found topic %s with %d partitions", topic, numPartitions)

			if numPartitions > 0 {
				return
			}
		}

		log.Debugf("retrying after 1s")
		time.Sleep(time.Second)
		continue
	}
}

// Close close kafka client
func (c *client) Close() error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	return c.writer.Close()
}

// Consumer kafka consumer interface
type Consumer interface {
	Close() error
	Run(ctx context.Context)
}

// consumer kafka consumer
type consumer struct {
	name    string
	wg      sync.WaitGroup
	reader  *kafka.Reader
	logger  *logger.MessageLogger
	closed  atomic.Bool
	handler func(context.Context, *Message) error
}

// Run start a goroutine to consume kafka message
func (c *consumer) Run(ctx context.Context) {
	c.wg.Add(1)
	go c.worker(ctx)
}

// worker start a goroutine to consume kafka message
func (c *consumer) worker(ctx context.Context) {
	for {
		if c.closed.Load() {
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
	c.wg.Done()
}

// handleMsg process kafka message
func (c *consumer) handleMsg(msg *kafka.Message) {
	if err := c.handler(context.Background(), msg); err != nil {
		log.Errorf("Failed to handler kafka msg: [%s] , caused by %s \n"+
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

	// commit message
	if err := c.reader.CommitMessages(ctx, *msg); err != nil {
		// commit failed
		log.Errorf("Failed to commit kafka msg: [%s], caused by %s \n"+
			"\tTopic: %s\n\tpartition %d: %d",
			msg.Value, err, msg.Topic, msg.Partition, msg.Offset)
		return
	}
}

// Close close kafka consumer
// we wrap the reader.Close() to make it compatible with kafka.Close()
func (c *consumer) Close() error {
	c.closed.Store(true)
	if err := c.reader.Close(); err != nil {
		return err
	}
	c.wg.Wait()
	log.Infof("%s reader Closed", c.name)
	return nil
}
