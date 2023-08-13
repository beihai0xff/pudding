// Package main provides the main function of pudding broker
package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "go.uber.org/automaxprocs"

	"github.com/beihai0xff/pudding/app/broker/server"
	"github.com/beihai0xff/pudding/configs"
	"github.com/beihai0xff/pudding/pkg/autocert"
	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/pkg/shutdown"
)

func main() {
	flag := configs.GetConfigFlagSet()
	redisURL := flag.String("redis-url", "", "The redis connection url")

	configs.ParseFlag()

	conf := configs.ParseBrokerConfig(*configs.GetConfigFilePath(), configs.WithRedisURL(*redisURL))
	autocert.New(conf.ServerConfig.HostDomain)

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
