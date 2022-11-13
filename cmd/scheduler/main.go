package main

import (
	"flag"
	"fmt"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	pbhealth "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	pb "github.com/beihai0xff/pudding/api/gen/scheduler/v1"
	"github.com/beihai0xff/pudding/app/scheduler"
	"github.com/beihai0xff/pudding/app/scheduler/broker"
	"github.com/beihai0xff/pudding/app/scheduler/pkg/configs"
	"github.com/beihai0xff/pudding/pkg/lock"
	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/pkg/log/logger"
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

	log.RegisterLogger(log.DefaultLoggerName, log.WithCallerSkip(1))
	log.RegisterLogger(logger.PulsarLoggerName, log.WithCallerSkip(1))
	log.RegisterLogger(logger.GRPCLoggerName, log.WithCallerSkip(1))
	logger.GetGRPCLogger()

	// log.RegisterLogger("gorm_log", log.WithCallerSkip(3))

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// register server
	server := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    time.Minute,
			Timeout: 5 * time.Second,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             1,
			PermitWithoutStream: true,
		}),
	)

	// register scheduler server
	delay, realtime := newQueue()
	handler := scheduler.NewHandler(scheduler.New(configs.GetSchedulerConfig(), delay, realtime))
	pb.RegisterSchedulerServiceServer(server, handler)
	// register health check server
	healthcheck := health.NewServer()
	pbhealth.RegisterHealthServer(server, healthcheck)
	// Register reflection service on gRPC server.
	reflection.Register(server)

	go func() {
		// asynchronously inspect dependencies and toggle serving status as needed
		healthcheck.SetServingStatus(pb.SchedulerService_ServiceDesc.ServiceName, pbhealth.HealthCheckResponse_SERVING)
	}()

	log.Infof("server listening at %v", lis.Addr())
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func newQueue() (broker.DelayQueue, broker.RealTimeQueue) {
	rdb := redis.New(configs.GetRedisConfig())

	pulsarClient := pulsar.New(configs.GetPulsarConfig())

	lock.Init(rdb)

	return broker.NewDelayQueue(rdb), broker.NewRealTimeQueue(pulsarClient)
}
