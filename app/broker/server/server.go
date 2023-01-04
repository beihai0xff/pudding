// Package server provides the start and dependency registration of the broker server
package server

import (
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"

	pb "github.com/beihai0xff/pudding/api/gen/pudding/broker/v1"
	"github.com/beihai0xff/pudding/app/broker/pkg/configs"
	"github.com/beihai0xff/pudding/pkg/grpc/launcher"
	"github.com/beihai0xff/pudding/pkg/grpc/resolver"
	"github.com/beihai0xff/pudding/pkg/lock"
	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/pkg/log/logger"
	"github.com/beihai0xff/pudding/pkg/redis"
	"github.com/beihai0xff/pudding/pkg/utils"
)

const (
	// custom http endpoint prifix path
	httpPrefix = "/pudding/broker"
)

var (
	// healthEndpointPath health check http endpoint path.
	healthEndpointPath = utils.GetHealthEndpointPath(httpPrefix)
	// swaggerEndpointPath Swagger ui http endpoint path.
	swaggerEndpointPath = utils.GetSwaggerEndpointPath(httpPrefix)
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
func RegisterResolver(grpcPort, httpPort int) []*resolver.Pair {
	consulURL := configs.GetNameServerURL()

	pairs := []*resolver.Pair{
		resolver.GRPCRegistration(pb.SchedulerService_ServiceDesc.ServiceName,
			grpcPort, resolver.WithConsulResolver(consulURL)),
		resolver.HTTPRegistration(healthEndpointPath,
			httpPort, resolver.WithConsulResolver(consulURL)),
	}
	return pairs
}

// StartServer starts the server.
func StartServer() (*grpc.Server, *health.Server, *http.Server) {
	grpcServer, healthcheck := launcher.StartGRPCServer(startSchedulerService)
	httpServer := launcher.StartHTTPServer(healthEndpointPath, swaggerEndpointPath)
	return grpcServer, healthcheck, httpServer
}
