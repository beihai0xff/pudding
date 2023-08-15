// Package main provides the main function of the trigger
package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "go.uber.org/automaxprocs"

	"github.com/beihai0xff/pudding/app/trigger/server"
	"github.com/beihai0xff/pudding/configs"
	"github.com/beihai0xff/pudding/pkg/autocert"
	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/pkg/shutdown"
)

func main() {
	// TODO: add flags
	// mysqlDSN := flag.String("mysql-url", "", "The server mysql connection dsn")
	// webhookPrefix := flag.String("webhook-prefix", "", "The server webhook prefix")
	flag.Parse()

	conf := configs.ParseTriggerConfig(*configs.GetConfigFilePath())
	autocert.New(conf.BaseConfig.HostDomain)

	server.RegisterLogger()

	interrupt := make(chan os.Signal, 1)

	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(interrupt)

	// start server
	grpcServer, healthcheck, httpServer := server.StartServer(conf)
	// register service to consul
	resolverPairs := server.RegisterResolver(conf)

	// block until a signal is received.
	sign := <-interrupt
	log.Warnf("received shutdown signal: %s", sign.String())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	shutdown.GracefulShutdown(ctx,
		shutdown.ResolverDeregister(resolverPairs...),
		shutdown.HealthServerShutdown(healthcheck),
		shutdown.HTTPServerShutdown(httpServer),
		shutdown.GRPCServerShutdown(grpcServer),
		shutdown.LogSync(),
	)
}
