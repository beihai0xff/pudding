package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	pbhealth "google.golang.org/grpc/health/grpc_health_v1"

	pb "github.com/beihai0xff/pudding/api/gen/pudding/scheduler/v1"
	"github.com/beihai0xff/pudding/app/scheduler/pkg/configs"
	"github.com/beihai0xff/pudding/pkg/lock"
	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/pkg/log/logger"
	"github.com/beihai0xff/pudding/pkg/redis"
	"github.com/beihai0xff/pudding/pkg/resolver"
	"github.com/beihai0xff/pudding/pkg/utils"
)

const serviceName = "pudding.scheduler"

var (
	grpcPort = flag.Int("grpcPort", 50051, "The grpc server grpcPort")
	httpPort = flag.Int("httpPort", 8081, "The http server grpcPort")

	confPath  = flag.String("config", "./config.yaml", "The server config file path")
	redisURL  = flag.String("redis", "", "The server redis url")
	pulsarURL = flag.String("pulsar", "", "The server pulsar url")
)

func main() {
	flag.Parse()

	configs.Init(*confPath, configs.WithRedisURL(*redisURL), configs.WithPulsarURL(*pulsarURL))
	registerLogger()

	// log.RegisterLogger("gorm_log", log.WithCallerSkip(3))
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(interrupt)

	grpcLis, err := net.Listen("tcp", fmt.Sprintf(":%d", *grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	httpLis, err := net.Listen("tcp", fmt.Sprintf(":%d", *httpPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer, healthcheck := startGrpcService(grpcLis)
	httpServer := startHTTPService(grpcLis, httpLis)

	// register service to consul
	resolver, err := resolver.NewConsulResolver(configs.GetConsulURL())
	if err != nil {
		log.Fatalf("failed to create resolver: %w", err)
	}
	serviceID, err := resolver.Register(serviceName, utils.GetOutBoundIP(), *grpcPort)
	if err != nil {
		log.Fatalf("failed to register service: %w", err)
	}
	// block until a signal is received.
	sign := <-interrupt
	log.Warn("received shutdown signal ", sign.String())

	_ = resolver.Deregister(serviceID)
	gracefulShutdown(grpcServer, healthcheck, httpServer)
}

func registerLogger() {
	log.RegisterLogger(log.DefaultLoggerName, log.WithCallerSkip(1))
	log.RegisterLogger(logger.PulsarLoggerName, log.WithCallerSkip(1))
	log.RegisterLogger(logger.GRPCLoggerName, log.WithCallerSkip(1))
	logger.GetGRPCLogger()
	rdb := redis.New(configs.GetRedisConfig())
	lock.Init(rdb)
}

func gracefulShutdown(s *grpc.Server, healthcheck *health.Server, httpServer *http.Server) {
	// 1. tell the load balancer to stop sending new requests
	healthcheck.SetServingStatus(pb.SchedulerService_ServiceDesc.ServiceName, pbhealth.HealthCheckResponse_NOT_SERVING)
	time.Sleep(3 * time.Second)
	// 2. stop accepting new HTTP requests and wait for existing HTTP requests to finish
	_ = httpServer.Shutdown(context.Background())
	// 3. stop accepting new connections and RPCs and blocks until all the pending RPCs are finished.
	s.GracefulStop()
	// 4. flushing any buffered log entries
	log.Sync()
}
