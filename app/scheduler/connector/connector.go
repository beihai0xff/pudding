package connector

import (
	"context"

	"github.com/beihai0xff/pudding/app/scheduler/connector/pulsar_connector"
	"github.com/beihai0xff/pudding/app/scheduler/pkg/configs"
	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/pkg/mq/pulsar"
	"github.com/beihai0xff/pudding/types"
)

// nolint:lll
//go:generate mockgen -destination=../../../test/mock/connector_mock.go -package=mock github.com/beihai0xff/pudding/app/scheduler/connector RealTimeConnector

// RealTimeConnector is a connector which can send messages to the realtime queue
// the realtime queue can store or consume messages in realtime
type RealTimeConnector interface {
	// Produce produce a Message to the queue in real time
	Produce(ctx context.Context, msg *types.Message) error
	// NewConsumer new a consumer to consume Messages from the realtime queue in background
	NewConsumer(topic, group string, batchSize int, fn types.HandleMessage) error
	// Close the queue
	Close() error
}

// NewConnector create a new RealTime Queue Connector
func NewConnector(connector string) RealTimeConnector {
	switch connector {
	case "pulsar":
		return pulsar_connector.NewRealTimeQueue(pulsar.New(configs.GetPulsarConfig()))
	default:
		log.Fatalf("unknown connector type: [%s]", connector)
	}
	return nil

}
