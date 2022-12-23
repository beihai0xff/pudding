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
	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/pkg/shutdown"
)

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
	server.RegisterLogger()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(interrupt)

	// start server
	grpcServer, healthcheck, httpServer := server.StartServer(*grpcPort, *httpPort)
	// register service to consul
	resolverPairs := server.RegisterResolver(*grpcPort, *httpPort)

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
