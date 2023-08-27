package cluster

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.etcd.io/etcd/api/v3/mvccpb"
	v3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"

	"github.com/beihai0xff/pudding/pkg/log"
)

// Queue is a distributed queue implementation based on etcd
type Queue interface {
	// Produce adds an element to the queue.
	Produce(m *Message) (string, error)
	// Consume returns the first element in the queue.
	Consume(ctx context.Context) (*Message, error)
	// Commit removes the element from the queue.
	Commit(msg *Message) error
}

// Message is the unit of the queue
type Message struct {
	// Key is the unique identifier of the message
	Key string
	// Value is the content of the message
	Value string
	// the key is unique in the queue
	Unique bool
}

// queue implements a single-reader, multi-writer distributed queue.
// /namespace/topic/msg/key
// /namespace/topic/acked/key_rev
// /namespace/topic/unacked/key_rev
// /namespace/topic/consumerID
type queue struct {
	client        *v3.Client
	topic         string
	consumerMutex Mutex
	consumerID    string
}

// Queue creates a new queue.
func (c *cluster) Queue(topic string) (Queue, error) {
	return newQueue(c, topic)
}

func newQueue(cluster *cluster, topic string) (*queue, error) {
	q := queue{
		client:     cluster.client,
		topic:      topic,
		consumerID: uuid.New().String(),
	}

	m, err := cluster.Mutex(q.getConsumerLockPath(), 5*time.Second)
	if err != nil {
		return nil, err
	}

	q.consumerMutex = m

	return &q, nil
}

func (q *queue) Produce(m *Message) (string, error) {
	messageID := m.Key
	if !m.Unique {
		// the message do not need unique, add an uuid to make it won't conflict with other messages
		messageID = fmt.Sprintf("%s/%s", messageID, uuid.New().String())
	}

	// if the messageID is non-unique, an error will be returned
	if err := putNewKV(q.client.KV, q.getMsgPath()+messageID, m.Value, v3.NoLease); err != nil {
		return "", err
	}

	return messageID, nil
}

// Consume returns Produce()'d elements in FIFO order. If the
// queue is empty, Consume blocks until elements are available.
//
//nolint:gocyclo
func (q *queue) Consume(ctx context.Context) (*Message, error) {
	for {
		if held, err := q.consumerMutex.IsHeld(); err != nil {
			return nil, err
		} else if !held {
			// wait for the lock, only one consumer can consume the message
			if err := q.consumerMutex.Lock(ctx); err != nil {
				return nil, err
			}
		}

		resp, err := q.client.KV.Get(ctx, q.getMsgPath(), v3.WithFirstRev()...)
		if err != nil {
			return nil, err
		}

		if len(resp.Kvs) != 0 {
			kv := resp.Kvs[0]

			rev, err := q.tryUnackMessage(kv)
			if err != nil {
				return nil, err
			}

			if rev != 0 {
				_, err := waitKeyEvents(q.client, q.getUnAckedPath(), rev, []mvccpb.Event_EventType{mvccpb.DELETE})
				if err != nil {
					return nil, err
				}

				continue
			}

			return q.message(kv), nil
		}

		log.Infof("no message in the topic [%s] , wait for new message", q.topic)

		ev, err := waitPrefixEvents(
			q.client,
			q.getMsgPath(),
			resp.Header.Revision,
			[]mvccpb.Event_EventType{mvccpb.PUT})
		if err != nil {
			return nil, err
		}

		// If watch type is PUT, it indicates new data has been stored to the key.
		return q.message(ev.Kv), nil
	}
}

func (q *queue) Commit(msg *Message) error {
	key := fmt.Sprintf("%s%s", q.getMsgPath(), msg.Key)
	_, err := concurrency.NewSTM(q.client, func(stm concurrency.STM) error {
		rev := stm.Rev(key)
		// the message has been consumed
		if rev == 0 {
			return nil
		}
		stm.Del(q.getUnAckedPath())
		stm.Del(key)
		return nil
	})

	return err
}

// tryUnackMessage try to lock the message
// if the message has been locked, return false
func (q *queue) tryUnackMessage(kv *mvccpb.KeyValue) (int64, error) {
	var rev int64

	_, err := concurrency.NewSTM(q.client, func(stm concurrency.STM) error {
		rev = stm.Rev(q.getUnAckedPath())
		if rev != 0 {
			v := unAckedValue{}
			v.parse(stm.Get(q.getUnAckedPath()))
			// the message has been locked by the consumer
			// need to wait for the consumer to commit the message
			if v.ConsumerID == q.consumerID {
				// return the rev of the message
				return nil
			}
			// the message has been locked by other consumer
			// but the consumer has not commit the message and the consumer is dead
			// so we can lock the message and reconsume it
			rev = 0
		}

		v := unAckedValue{
			ConsumerID: q.consumerID,
			Rev:        kv.CreateRevision,
		}
		stm.Put(q.getUnAckedPath(), v.jsonFormat())
		return nil
	})

	return rev, err
}

func (q *queue) message(kv *mvccpb.KeyValue) *Message {
	return &Message{
		Key:   strings.TrimPrefix(string(kv.Key), q.getMsgPath()),
		Value: string(kv.Value),
	}
}

func (q *queue) getMsgPath() string {
	return fmt.Sprintf("%s/msg/", q.topic)
}

func (q *queue) getUnAckedPath() string {
	return fmt.Sprintf("%s/unacked", q.topic)
}

func (q *queue) getConsumerLockPath() string {
	return fmt.Sprintf("%s/consumer", q.topic)
}

type unAckedValue struct {
	ConsumerID string `json:"consumer_id"`
	Rev        int64  `json:"rev"`
}

func (v *unAckedValue) jsonFormat() string {
	b, _ := json.Marshal(v)
	return string(b)
}

func (v *unAckedValue) parse(s string) {
	_ = json.Unmarshal([]byte(s), v)
}
