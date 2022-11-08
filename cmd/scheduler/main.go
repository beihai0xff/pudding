package main

import (
	"flag"
	"fmt"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	pb "github.com/beihai0xff/pudding/api/scheduler/v1"
	"github.com/beihai0xff/pudding/app/scheduler"
	"github.com/beihai0xff/pudding/app/scheduler/broker"
	"github.com/beihai0xff/pudding/app/scheduler/pkg/configs"
	"github.com/beihai0xff/pudding/pkg/lock"
	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/pkg/mq/pulsar"
	"github.com/beihai0xff/pudding/pkg/redis"
)

var (
	port     = flag.Int("port", 50051, "The server port")
	confPath = flag.String("config", "./config.yaml", "The server config file path")
)

func main() {
	flag.Parse()
	configs.Init(*confPath)

	log.RegisterLogger("default", log.WithCallerSkip(1))
	log.RegisterLogger("pulsar_log", log.WithCallerSkip(1))

	// log.RegisterLogger("gorm_log", log.WithCallerSkip(3))

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	delay, realtime := newQueue()

	s := grpc.NewServer()
	pb.RegisterPuddingServer(s, scheduler.NewHandler(scheduler.New(configs.GetSchedulerConfig(), delay, realtime)))
	// Register reflection service on gRPC server.
	reflection.Register(s)
	log.Infof("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func newQueue() (broker.DelayQueue, broker.RealTimeQueue) {
	rdb := redis.New(configs.GetRedisConfig())

	pulsarClient := pulsar.New(configs.GetPulsarConfig())

	lock.Init(rdb)

	return broker.NewDelayQueue(rdb), broker.NewRealTimeQueue(pulsarClient)
}
