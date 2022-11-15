package main

import (
	"context"
	"net"
	"net/http"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	pbhealth "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	pb "github.com/beihai0xff/pudding/api/gen/pudding/scheduler/v1"
	"github.com/beihai0xff/pudding/app/scheduler"
	"github.com/beihai0xff/pudding/app/scheduler/pkg/configs"
	"github.com/beihai0xff/pudding/pkg/log"
)

func startGrpcService(lis net.Listener) (*grpc.Server, *health.Server) {
	// register server
	server := grpc.NewServer(
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
		)),
	)

	// register scheduler server
	delay, realtime := newQueue()
	handler := scheduler.NewHandler(scheduler.New(configs.GetSchedulerConfig(), delay, realtime))
	pb.RegisterSchedulerServiceServer(server, handler)
	// register health check server
	healthcheck := health.NewServer()
	pbhealth.RegisterHealthServer(server, healthcheck)
	// Register reflection service on gRPC server.
	reflection.Register(server)

	go func() {
		// asynchronously inspect dependencies and toggle serving status as needed
		healthcheck.SetServingStatus(pb.SchedulerService_ServiceDesc.ServiceName, pbhealth.HealthCheckResponse_SERVING)
		log.Infof("server listening at %v", lis.Addr())
		if err := server.Serve(lis); err != nil {
			log.Fatalf("failed to start grpc serve: %v", err)
		}
	}()

	return server, healthcheck
}

func startHTTPService(grpcLis, httpLis net.Listener) *http.Server {
	conn, err := grpc.DialContext(
		context.Background(),
		grpcLis.Addr().String(),
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to dial server: %w", err)
	}

	// gRPC-Gateway httpServer
	gwmux := runtime.NewServeMux()
	err = pb.RegisterSchedulerServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalf("Failed to register gateway: %w", err)
	}

	// 定义HTTP server配置
	httpServer := &http.Server{
		Addr:    httpLis.Addr().String(),
		Handler: gwmux,
	}

	go func() {
		log.Infof("Serving gRPC-Gateway on %s", httpLis.Addr().String())
		if err = httpServer.ListenAndServe(); err != nil {
			log.Fatalf("Failed to serve gRPC-Gateway: %w", err)
		}
	}()

	return httpServer
}
