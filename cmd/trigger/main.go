package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/beihai0xff/pudding/app/trigger/pkg/configs"
	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/pkg/log/logger"
	"github.com/beihai0xff/pudding/pkg/resolver"
	"github.com/beihai0xff/pudding/pkg/shutdown"
	"github.com/beihai0xff/pudding/pkg/utils"
)

const serviceName = "pudding.trigger"

var (
	grpcPort = flag.Int("grpcPort", 50051, "The grpc server port")
	httpPort = flag.Int("httpPort", 8081, "The http server port")

	confPath = flag.String("config", "./config.yaml", "The server config file path")
	mysqlDSN = flag.String("mysql", "", "The server mysql dsn")
)

func main() {
	flag.Parse()

	configs.Init(*confPath, configs.WithMySQLDSN(*mysqlDSN))
	registerLogger()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(interrupt)

	// start server
	grpcServer, healthcheck, httpServer := startServer()
	// register service to consul
	rsv, serviceID := serviceRegistration()

	// block until a signal is received.
	sign := <-interrupt
	log.Warnf("received shutdown signal: %s", sign.String())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	shutdown.GracefulShutdown(ctx, shutdown.ResolverDeregister(rsv, serviceID),
		shutdown.HealthcheckServerShutdown(healthcheck),
		shutdown.HTTPServerShutdown(httpServer),
		shutdown.GRPCServerShutdown(grpcServer),
		shutdown.LogSync(),
	)
}

func serviceRegistration() (resolver.Resolver, string) {
	rsv, err := resolver.NewConsulResolver(configs.GetConsulURL())
	if err != nil {
		log.Fatalf("failed to create rsv: %v", err)
	}
	serviceID, err := rsv.Register(serviceName, utils.GetOutBoundIP(), *grpcPort)
	if err != nil {
		log.Fatalf("failed to register service: %v", err)
	}
	return rsv, serviceID
}

func registerLogger() {
	log.RegisterLogger(log.DefaultLoggerName, log.WithCallerSkip(1))
	log.RegisterLogger("gorm_log", log.WithCallerSkip(3))
	logger.GetGRPCLogger()
}
