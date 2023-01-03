package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/beihai0xff/pudding/app/trigger/pkg/configs"
	"github.com/beihai0xff/pudding/app/trigger/server"
	"github.com/beihai0xff/pudding/pkg/grpc/args"
	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/pkg/shutdown"
)

var (
	mysqlDSN = flag.String("mysql", "", "The server mysql dsn")

	webhookPrefix = flag.String("webhook_prefix", "", "The server webhook prefix")
)

func main() {
	flag.Parse()

	configs.Init(*args.ConfigPath, configs.WithMySQLDSN(*mysqlDSN), configs.WithWebhookPrefix(*webhookPrefix))
	server.RegisterLogger()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(interrupt)

	// start server
	grpcServer, healthcheck, httpServer := server.StartServer()
	// register service to consul
	resolverPairs := server.RegisterResolver(*args.GRPCPort, *args.HTTPPort)

	// block until a signal is received.
	sign := <-interrupt
	log.Warnf("received shutdown signal: %s", sign.String())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	shutdown.GracefulShutdown(ctx,
		shutdown.ResolverDeregister(resolverPairs...),
		shutdown.HealthcheckServerShutdown(healthcheck),
		shutdown.HTTPServerShutdown(httpServer),
		shutdown.GRPCServerShutdown(grpcServer),
		shutdown.LogSync(),
	)
}
