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
	"github.com/beihai0xff/pudding/pkg/log"
	pulsar2 "github.com/beihai0xff/pudding/pkg/mq/pulsar"
	"github.com/beihai0xff/pudding/pkg/redis"
)

var (
	port = flag.Int("port", 50051, "The server port")
	conf = flag.String("config", ".././config.yaml", "The server config file path")
)

func main() {
	flag.Parse()
	configs.Init(*conf)

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

	pulsar := pulsar2.New(configs.GetPulsarConfig())

	return broker.NewDelayQueue(rdb), broker.NewRealTimeQueue(pulsar)
}
