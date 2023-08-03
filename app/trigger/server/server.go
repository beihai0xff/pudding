// Package server provides the start and dependency registration of the trigger server
package server

import (
	"net/http"

	_ "github.com/mbobakov/grpc-consul-resolver" // It's important
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"

	pb "github.com/beihai0xff/pudding/api/gen/pudding/trigger/v1"
	"github.com/beihai0xff/pudding/configs"
	"github.com/beihai0xff/pudding/pkg/grpc/launcher"
	"github.com/beihai0xff/pudding/pkg/grpc/resolver"
	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/pkg/log/logger"
	"github.com/beihai0xff/pudding/pkg/utils"
)

const (
	// custom http endpoint prifix path
	httpPrefix = "/pudding/trigger"
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
	// gorm log need 3 skip
	log.RegisterLogger(logger.BackendLoggerName, log.WithCallerSkip(3))
}

// RegisterResolver registers the service to the resolver.
func RegisterResolver(conf *configs.TriggerConfig) []*resolver.Pair {
	baseConfig := conf.ServerConfig.BaseConfig
	consulURL := baseConfig.NameServerURL

	return []*resolver.Pair{
		resolver.GRPCRegistration(pb.CronTriggerService_ServiceDesc.ServiceName,
			baseConfig.GRPCPort, resolver.WithConsulResolver(consulURL)),
		resolver.GRPCRegistration(pb.WebhookTriggerService_ServiceDesc.ServiceName,
			baseConfig.GRPCPort, resolver.WithConsulResolver(consulURL)),
		resolver.HTTPRegistration(healthEndpointPath,
			baseConfig.HTTPPort, resolver.WithConsulResolver(consulURL)),
	}
}

// StartServer starts the server.
func StartServer(conf *configs.TriggerConfig) (*grpc.Server, *health.Server, *http.Server) {
	baseConfig := conf.ServerConfig.BaseConfig
	grpcServer, healthcheck := launcher.StartGRPCServer(&baseConfig,
		withCronTriggerService(conf), withWebhookTriggerService(conf))
	httpServer := launcher.StartHTTPServer(&baseConfig, healthEndpointPath, swaggerEndpointPath)

	return grpcServer, healthcheck, httpServer
}
