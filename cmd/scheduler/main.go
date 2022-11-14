package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
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
	registerLogger()

	// log.RegisterLogger("gorm_log", log.WithCallerSkip(3))
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(interrupt)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server, healthcheck := newServer()

	log.Infof("server listening at %v", lis.Addr())
	if err := server.Serve(lis); err != nil {
		panic(fmt.Errorf("failed to serve: %w", err))
	}

	// block until a signal is received.
	sign := <-interrupt
	log.Warn("received shutdown signal ", sign.String())

	gracefulShutdown(server, healthcheck)
}

func registerLogger() {
	log.RegisterLogger(log.DefaultLoggerName, log.WithCallerSkip(1))
	log.RegisterLogger(logger.PulsarLoggerName, log.WithCallerSkip(1))
	log.RegisterLogger(logger.GRPCLoggerName, log.WithCallerSkip(1))
	logger.GetGRPCLogger()
}

func newQueue() (broker.DelayQueue, broker.RealTimeQueue) {
	rdb := redis.New(configs.GetRedisConfig())

	pulsarClient := pulsar.New(configs.GetPulsarConfig())

	lock.Init(rdb)

	return broker.NewDelayQueue(rdb), broker.NewRealTimeQueue(pulsarClient)
}

func newServer() (*grpc.Server, *health.Server) {
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
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(),
			grpc_validator.UnaryServerInterceptor(),
		)),
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

	return server, healthcheck
}

func gracefulShutdown(s *grpc.Server, healthcheck *health.Server) {
	// 1. tell the load balancer to stop sending new requests
	healthcheck.SetServingStatus(pb.SchedulerService_ServiceDesc.ServiceName, pbhealth.HealthCheckResponse_NOT_SERVING)
	time.Sleep(3 * time.Second)
	// 2. stop accepting new connections and RPCs and blocks until all the pending RPCs are finished.
	s.GracefulStop()
	// 3. flushing any buffered log entries
	log.Sync()
}
