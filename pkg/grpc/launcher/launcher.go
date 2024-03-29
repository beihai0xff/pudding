// Package launcher provides a launcher to start gRPC server, health server and grpc gateway server.
package launcher

import (
	"context"
	btls "crypto/tls"
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
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	pbhealth "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	pb "github.com/beihai0xff/pudding/api/gen/pudding/broker/v1"

	"github.com/beihai0xff/pudding/pkg/tls"

	"github.com/beihai0xff/pudding/configs"
	"github.com/beihai0xff/pudding/pkg/grpc/swagger"
	"github.com/beihai0xff/pudding/pkg/log"
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
	// This allows the gRPC server to be introspected by clients using reflection
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

// createGRPCServer creates a gRPC server
func createGRPCServer(config *configs.BaseConfig) *grpc.Server {
	options := []grpc.ServerOption{
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
			unaryServerRequestLog(),
			// define open telemetry MeterProvider
			otelgrpc.UnaryServerInterceptor(otelgrpc.WithMeterProvider(otel.GetMeterProvider())),
		)),
	}

	if config.TLS.Enable {
		// define TLS configuration
		cred := credentials.NewTLS(tls.GetTLSConfig(tls.ConfigTypeServer))
		options = append(options, grpc.Creds(cred))
	}

	// init grpc server
	return grpc.NewServer(options...)
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
	httpServer := &http.Server{
		Addr:        fmt.Sprintf(":%d", config.HTTPPort),
		Handler:     Handler(gwmux, WithRequestLog, WithRedirectToHTTPS),
		ReadTimeout: 10 * time.Second,
		IdleTimeout: 30 * time.Second,
	}

	if config.TLS.Enable {
		c := tls.GetTLSConfig(tls.ConfigTypeServer).Clone()
		c.ClientAuth = btls.NoClientCert

		httpServer.TLSConfig = c
	}

	go func() {
		log.Infof("http server listening at %v", httpServer.Addr)

		if config.TLS.Enable {
			if err := httpServer.ListenAndServeTLS(config.TLS.ServerCert, config.TLS.ServerKey); err != nil {
				if errors.Is(err, http.ErrServerClosed) {
					log.Info("https server closed")
					return
				}

				log.Fatalf("Failed to serve https server: %v", err)
			}
		} else {
			if err := httpServer.ListenAndServe(); err != nil {
				if errors.Is(err, http.ErrServerClosed) {
					log.Info("http server closed")
					return
				}

				log.Fatalf("Failed to serve http server: %v", err)
			}
		}
	}()

	return httpServer
}

// createGRPCLocalClient creates a gRPC client to the local gRPC server.
// It is used by the gRPC-Gateway.
func createGRPCLocalClient(config *configs.BaseConfig) *grpc.ClientConn {
	log.Debugf("creating local gRPC client...")

	options := []grpc.DialOption{
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
			Timeout:             time.Second,      // wait 1 second for ping ack before considering the connection dead
			PermitWithoutStream: true,             // send pings even without active streams
		}),
	}

	if config.TLS.Enable {
		cred := credentials.NewTLS(tls.GetTLSConfig(tls.ConfigTypeClient))
		options = append(options, grpc.WithTransportCredentials(cred))
	} else {
		options = append(options, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	conn, err := grpc.DialContext(
		context.Background(),
		fmt.Sprintf("localhost:%d", config.GRPCPort),
		options...,
	)
	if err != nil {
		log.Panicf("Failed to dial server: %v", err)
	}

	log.Debug("create gRPC local client success")

	return conn
}
