package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/beihai0xff/pudding/app/broker/pkg/configs"
	"github.com/beihai0xff/pudding/app/broker/server"
	// nolint:revive
	. "github.com/beihai0xff/pudding/pkg/grpc/args"
	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/pkg/shutdown"
)

var (
	redisURL  = flag.String("redis_url", "", "The server redis url")
	pulsarURL = flag.String("pulsar_url", "", "The server pulsar url")
)

func main() {
	flag.Parse()

	configs.Init(*ConfigPath, configs.WithRedisURL(*redisURL), configs.WithPulsarURL(*pulsarURL))
	server.RegisterLogger()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(interrupt)

	// start server
	grpcServer, healthcheck, httpServer := server.StartServer()
	// register service to consul
	resolverPairs := server.RegisterResolver(*GRPCPort, *HTTPPort)

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
