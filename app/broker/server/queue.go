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
func newDelayStorage(conf *configs.BrokerConfig) storage.DelayStorage {
	broker := conf.ServerConfig.Broker
	switch broker {
	case "redis":
		// parse Polling delay queue interval
		t, err := time.ParseDuration(conf.ServerConfig.TimeSliceInterval)
		if err != nil {
			log.Fatalf("failed to parse '%s' to time.Duration: %v", conf.ServerConfig.TimeSliceInterval, err)
		}

		log.Infof("timeSlice interval is: %f seconds", t.Seconds())

		return redis_storage.NewDelayStorage(rdb.New(&conf.RedisConfig), uint64(t.Seconds()))
	default:
		log.Fatalf("unknown broker type: [%s]", broker)
	}

	return nil
}

// newConnector create a new RealTime Queue Connector
func newConnector(conf *configs.BrokerConfig) connector.RealTimeConnector {
	connectorName := conf.ServerConfig.Connector
	switch connectorName {
	case "pulsar":
		log.Fatalf("pulsar connectorName is not implemented yet")
	case "kafka":
		return kafka_connector.NewConnector(kafka.New(&conf.KafkaConfig))
	case "redis":
		return redis_connector.NewConnector(rdb.New(&conf.RedisConfig))
	default:
		log.Fatalf("unknown connectorType type: [%s]", connectorName)
	}

	return nil
}
