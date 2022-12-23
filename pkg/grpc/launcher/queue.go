package launcher

import (
	"time"

	"github.com/beihai0xff/pudding/app/broker/connector"
	"github.com/beihai0xff/pudding/app/broker/connector/pulsar_connector"
	"github.com/beihai0xff/pudding/app/broker/pkg/configs"
	"github.com/beihai0xff/pudding/app/broker/storage"
	"github.com/beihai0xff/pudding/app/broker/storage/redis_storage"
	conf "github.com/beihai0xff/pudding/configs"
	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/pkg/mq/pulsar"
	rdb "github.com/beihai0xff/pudding/pkg/redis"
)

func newQueue(config *conf.SchedulerConfig) (storage.DelayStorage, connector.RealTimeConnector) {
	return newDelayStorage(config.Broker), newConnector(config.Connector)
}

// NewDelayBroker create a new DelayStorage
func newDelayStorage(broker string) storage.DelayStorage {
	switch broker {
	case "redis":
		// parse Polling delay queue interval
		interval := configs.GetSchedulerConfig().TimeSliceInterval
		t, err := time.ParseDuration(interval)
		if err != nil {
			log.Fatalf("failed to parse '%s' to time.Duration: %v", interval, err)
		}
		log.Infof("timeSlice interval is: %d seconds", t.Seconds())
		return redis_storage.NewDelayQueue(rdb.New(configs.GetRedisConfig()), int64(t.Seconds()))
	default:
		log.Fatalf("unknown broker type: [%s]", broker)
	}
	return nil
}

// newConnector create a new RealTime Queue Connector
func newConnector(connector string) connector.RealTimeConnector {
	switch connector {
	case "pulsar":
		return pulsar_connector.NewRealTimeQueue(pulsar.New(configs.GetPulsarConfig()))
	default:
		log.Fatalf("unknown connector type: [%s]", connector)
	}
	return nil
}
