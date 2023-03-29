// Package server provides the start and dependency registration of the broker server
package server

import (
	"time"

	"github.com/beihai0xff/pudding/app/broker/connector"
	"github.com/beihai0xff/pudding/app/broker/connector/kafka_connector"
	"github.com/beihai0xff/pudding/app/broker/connector/redis_connector"
	"github.com/beihai0xff/pudding/app/broker/storage"
	"github.com/beihai0xff/pudding/app/broker/storage/redis_storage"
	"github.com/beihai0xff/pudding/configs"
	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/pkg/mq/kafka"
	rdb "github.com/beihai0xff/pudding/pkg/redis"
)

func newQueue(config *configs.BrokerConfig) (storage.DelayStorage, connector.RealTimeConnector) {
	return newDelayStorage(config), newConnector(config)
}

// newDelayStorage create a new DelayStorage
func newDelayStorage(config *configs.BrokerConfig) storage.DelayStorage {
	broker := config.ServerConfig.Broker
	switch broker {
	case "redis":
		// parse Polling delay queue interval
		interval := config.ServerConfig.TimeSliceInterval
		t, err := time.ParseDuration(interval)
		if err != nil {
			log.Fatalf("failed to parse '%s' to time.Duration: %v", interval, err)
		}
		log.Infof("timeSlice interval is: %f seconds", t.Seconds())
		return redis_storage.NewDelayStorage(rdb.New(config.RedisConfig), uint64(t.Seconds()))
	default:
		log.Fatalf("unknown broker type: [%s]", broker)
	}
	return nil
}

// newConnector create a new RealTime Queue Connector
func newConnector(config *configs.BrokerConfig) connector.RealTimeConnector {
	connectorName := config.ServerConfig.Connector
	switch connectorName {
	case "pulsar":
		log.Fatalf("pulsar connectorName is not implemented yet")
	case "kafka":
		return kafka_connector.NewConnector(kafka.New(config.KafkaConfig))
	case "redis":
		return redis_connector.NewConnector(rdb.New(config.RedisConfig))
	default:
		log.Fatalf("unknown connectorType type: [%s]", connectorName)
	}
	return nil
}
