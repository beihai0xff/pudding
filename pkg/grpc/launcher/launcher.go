// Package launcher provides a launcher to start gRPC server, health server and grpc gateway server.
package launcher

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/health"
	pbhealth "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	pb "github.com/beihai0xff/pudding/api/gen/pudding/broker/v1"
	"github.com/beihai0xff/pudding/configs"
	"github.com/beihai0xff/pudding/pkg/grpc/swagger"
	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/pkg/log/logger"
	"github.com/beihai0xff/pudding/pkg/otel"
)

// StartServiceFunc is a function that starts a service.
// Note that the service name is passed by reference and can be modified.
type StartServiceFunc func(server *grpc.Server, serviceName *string) error

func getListen(port int) net.Listener {
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Panicf("failed to listen: %v", err)
	}
	return listen
}

// StartGRPCServer starts the gRPC server with the given service.
// It serves the given gRPC service and the gRPC-healthz.
func StartGRPCServer(config *configs.BaseConfig, opts ...StartServiceFunc) (
	*grpc.Server, *health.Server) {

	log.Info("starting grpc server ...")

	server := createGRPCServer(config)

	// register health server
	healthServer := health.NewServer()
	pbhealth.RegisterHealthServer(server, healthServer)

	for _, opt := range opts {
		serviceName := ""
		if err := opt(server, &serviceName); err != nil {
			log.Panicf("failed to start service [%s]", serviceName)
		}
		// asynchronously inspect dependencies and toggle serving status as needed
		healthServer.SetServingStatus(serviceName, pbhealth.HealthCheckResponse_SERVING)
	}
	healthServer.Resume()
	// RegisterGRPC reflection service on gRPC server.
	// 提供该服务器端上可公开使用的 gRPC 服务的信息，
	// 服务反射向客户端提供了服务端注册的服务的信息，因此客户端不需要预编译服务定义就能与服务端交互
	// 通过此方式支持 grpcCRUL
	reflection.Register(server)

	grpcLis := getListen(config.GRPCPort)
	go func() {
		log.Infof("grpc server listening at %v", grpcLis.Addr())
		if err := server.Serve(grpcLis); err != nil {
			log.Panicf("failed to start grpc serve: %v", err)
		}
	}()

	return server, healthServer
}

// StartHTTPServer starts the HTTP server with the given service.
// It serves the gRPC-gateway, gRPC-healthz and the swagger UI.
// StartHTTPServer must be called after StartGRPCServer,
// because it uses the same listener, and HTTP server base on gRPC Gateway.
func StartHTTPServer(config *configs.BaseConfig, healthEndpointPath, swaggerEndpointPath string) *http.Server {
	log.Info("starting http server ...")

	conn := createGRPCLocalClient(config)

	// gRPC-Gateway httpServer
	gwmux := runtime.NewServeMux(runtime.WithHealthEndpointAt(pbhealth.NewHealthClient(conn), healthEndpointPath))
	swagger.RegisterHandler(gwmux, swaggerEndpointPath)
	otel.RegisterHandler(gwmux)

	if err := pb.RegisterSchedulerServiceHandler(context.Background(), gwmux, conn); err != nil {
		log.Panicf("Failed to register gateway: %v", err)
	}

	// define HTTP server configuration
	httpLis := getListen(config.HTTPPort)
	httpServer := &http.Server{
		Addr:        httpLis.Addr().String(),
		Handler:     gwmux,
		ReadTimeout: 10 * time.Second,
		IdleTimeout: 30 * time.Second,
	}

	go func() {
		log.Infof("http server listening at %v", httpLis.Addr())
		if err := httpServer.ServeTLS(httpLis, config.CertPath, config.KeyPath); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				log.Info("http server closed")
				return
			}
			log.Fatalf("Failed to serve gRPC-Gateway: %v", err)
		}
	}()

	return httpServer
}

// createGRPCLocalClient creates a gRPC client to the local gRPC server.
// It is used by the gRPC-Gateway.
func createGRPCLocalClient(config *configs.BaseConfig) *grpc.ClientConn {
	cred, err := credentials.NewClientTLSFromFile(config.CertPath, "localhost")
	if err != nil {
		log.Panicf("Failed to load Server TLS: %v", err)
	}
	conn, err := grpc.DialContext(
		context.Background(),
		// net.JoinHostPort("localhost", grpcLis.Addr().(*net.TCPAddr).Port),
		fmt.Sprintf("localhost:%d", config.GRPCPort),
		grpc.WithBlock(),
		grpc.WithTransportCredentials(cred),
	)
	if err != nil {
		log.Panicf("Failed to dial server: %v", err)
	}
	return conn
}

// createGRPCServer creates a gRPC server
func createGRPCServer(config *configs.BaseConfig) *grpc.Server {
	grpclog.SetLoggerV2(logger.GetGRPCLogger())

	// define TLS configuration
	cred, err := credentials.NewServerTLSFromFile(config.CertPath, config.KeyPath)
	if err != nil {
		log.Panicf("Failed to load Server TLS: %v", err)
	}

	// init grpc server
	return grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    time.Minute,
			Timeout: 5 * time.Second,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             1,
			PermitWithoutStream: true,
		}),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			grpc_recovery.UnaryServerInterceptor(),
			grpc_validator.UnaryServerInterceptor(),
			// define open telemetry MeterProvider
			otelgrpc.UnaryServerInterceptor(otelgrpc.WithMeterProvider(otel.GetMeterProvider())),
		)),
		grpc.Creds(cred),
	)
}
