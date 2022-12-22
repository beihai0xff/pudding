package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "github.com/beihai0xff/pudding/api/gen/pudding/trigger/v1"
	"github.com/beihai0xff/pudding/app/trigger/pkg/configs"
	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/pkg/log/logger"
	"github.com/beihai0xff/pudding/pkg/resolver"
	"github.com/beihai0xff/pudding/pkg/shutdown"
)

var (
	grpcPort = flag.Int("grpcPort", 50051, "The grpc server port")
	httpPort = flag.Int("httpPort", 8081, "The http server port")

	confPath = flag.String("config", "./config.yaml", "The server config file path")
	mysqlDSN = flag.String("mysql", "", "The server mysql dsn")

	webhookPrefix = flag.String("webhook_prefix", "", "The server webhook prefix")
)

func main() {
	flag.Parse()

	configs.Init(*confPath, configs.WithMySQLDSN(*mysqlDSN), configs.WithWebhookPrefix(*webhookPrefix))
	registerLogger()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(interrupt)

	// start server
	grpcServer, healthcheck, httpServer := startServer()
	// register grpc service to consul
	cronRsv, cronServiceID := resolver.GRPCRegistration(pb.CronTriggerService_ServiceDesc.ServiceName,
		*grpcPort, resolver.WithConsulResolver(configs.GetConsulURL()))
	webhookRsv, webhookServiceID := resolver.GRPCRegistration(pb.WebhookTriggerService_ServiceDesc.ServiceName,
		*grpcPort, resolver.WithConsulResolver(configs.GetConsulURL()))

	// block until a signal is received.
	sign := <-interrupt
	log.Warnf("received shutdown signal: %s", sign.String())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	shutdown.GracefulShutdown(ctx, shutdown.ResolverDeregister(cronRsv, cronServiceID),
		shutdown.ResolverDeregister(webhookRsv, webhookServiceID),
		shutdown.HealthcheckServerShutdown(healthcheck),
		shutdown.HTTPServerShutdown(httpServer),
		shutdown.GRPCServerShutdown(grpcServer),
		shutdown.LogSync(),
	)
}

func registerLogger() {
	log.RegisterLogger(log.DefaultLoggerName, log.WithCallerSkip(1))
	log.RegisterLogger("gorm_log", log.WithCallerSkip(3))
	logger.GetGRPCLogger()
}
