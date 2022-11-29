package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/beihai0xff/pudding/app/scheduler/pkg/configs"
	"github.com/beihai0xff/pudding/pkg/lock"
	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/pkg/log/logger"
	"github.com/beihai0xff/pudding/pkg/redis"
	"github.com/beihai0xff/pudding/pkg/resolver"
	"github.com/beihai0xff/pudding/pkg/shutdown"
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
	rsv, err := resolver.NewConsulResolver(configs.GetConsulURL())
	if err != nil {
		log.Fatalf("failed to create rsv: %w", err)
	}
	serviceID, err := rsv.Register(serviceName, utils.GetOutBoundIP(), *grpcPort)
	if err != nil {
		log.Fatalf("failed to register service: %w", err)
	}
	// block until a signal is received.
	sign := <-interrupt
	log.Warn("received shutdown signal ", sign.String())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	shutdown.GracefulShutdown(ctx, shutdown.ResolverDeregister(rsv, serviceID),
		shutdown.HealthcheckServerShutdown(healthcheck),
		shutdown.HTTPServerShutdown(httpServer),
		shutdown.GRPCServerShutdown(grpcServer),
		shutdown.LogSync(),
	)
}

func registerLogger() {
	log.RegisterLogger(log.DefaultLoggerName, log.WithCallerSkip(1))
	log.RegisterLogger(logger.PulsarLoggerName, log.WithCallerSkip(1))
	log.RegisterLogger(logger.GRPCLoggerName, log.WithCallerSkip(1))
	logger.GetGRPCLogger()
	rdb := redis.New(configs.GetRedisConfig())
	lock.Init(rdb)
}
