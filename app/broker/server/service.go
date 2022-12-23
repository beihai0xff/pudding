package server

import (
	"fmt"
	"net"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"

	pb "github.com/beihai0xff/pudding/api/gen/pudding/broker/v1"
	"github.com/beihai0xff/pudding/app/broker"
	"github.com/beihai0xff/pudding/app/broker/pkg/configs"
	"github.com/beihai0xff/pudding/pkg/grpc/launcher"
	resolver2 "github.com/beihai0xff/pudding/pkg/grpc/resolver"
	"github.com/beihai0xff/pudding/pkg/lock"
	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/pkg/log/logger"
	"github.com/beihai0xff/pudding/pkg/redis"
)

const (
	httpPrefix          = "/pudding/broker"
	healthEndpointPath  = httpPrefix + "/healthz"
	swaggerEndpointPath = httpPrefix + "/swagger"
)

// RegisterLogger registers the logger to the resolver.
func RegisterLogger() {
	log.RegisterLogger(log.DefaultLoggerName, log.WithCallerSkip(1))
	log.RegisterLogger(logger.PulsarLoggerName, log.WithCallerSkip(1))
	log.RegisterLogger(logger.GRPCLoggerName, log.WithCallerSkip(1))
	logger.GetGRPCLogger()
	rdb := redis.New(configs.GetRedisConfig())
	lock.Init(rdb)
}

// RegisterResolver registers the service to the resolver.
func RegisterResolver(grpcPort, httpPort int) []*resolver2.Pair {
	var pairs []*resolver2.Pair
	consulURL := configs.GetConsulURL()

	pairs = append(pairs, resolver2.GRPCRegistration(pb.SchedulerService_ServiceDesc.ServiceName,
		grpcPort, resolver2.WithConsulResolver(consulURL)))
	pairs = append(pairs, resolver2.HTTPRegistration(healthEndpointPath,
		httpPort, resolver2.WithConsulResolver(consulURL)))
	return pairs
}

// StartServer starts the server.
func StartServer(grpcPort, httpPort int) (*grpc.Server, *health.Server, *http.Server) {
	grpcLis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	httpLis, err := net.Listen("tcp", fmt.Sprintf(":%d", httpPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer, healthcheck := launcher.StartGRPCService(grpcLis, startSchedulerService)
	httpServer := launcher.StartHTTPService(grpcLis, httpLis, healthEndpointPath, swaggerEndpointPath)
	return grpcServer, healthcheck, httpServer
}

func startSchedulerService(server *grpc.Server, serviceName *string) error {
	// set serviceName
	// use string point to store serviceName, so that we can return it to the caller
	// this is not a good way to do this, but it works
	*serviceName = pb.SchedulerService_ServiceDesc.ServiceName

	// Initialize dependencies
	schedulerConfig := configs.GetSchedulerConfig()
	delay, realtime := newQueue(schedulerConfig)
	s := broker.New(schedulerConfig, delay, realtime)
	s.Run()
	handler := broker.NewHandler(s)

	// register scheduler server
	pb.RegisterSchedulerServiceServer(server, handler)

	return nil
}
