package main

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
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	pbhealth "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	_ "github.com/mbobakov/grpc-consul-resolver" // It's important

	"github.com/beihai0xff/pudding/api/gen/pudding/scheduler/v1"
	pb "github.com/beihai0xff/pudding/api/gen/pudding/trigger/v1"
	"github.com/beihai0xff/pudding/app/trigger/domain/cron"
	"github.com/beihai0xff/pudding/app/trigger/domain/webhook"
	"github.com/beihai0xff/pudding/app/trigger/pkg/configs"
	"github.com/beihai0xff/pudding/pkg/db/mysql"
	"github.com/beihai0xff/pudding/pkg/log"
)

func startServer() (*grpc.Server, *health.Server, *http.Server) {
	grpcLis, err := net.Listen("tcp", fmt.Sprintf(":%d", *grpcPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	httpLis, err := net.Listen("tcp", fmt.Sprintf(":%d", *httpPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer, healthcheck := startGrpcService(grpcLis)
	httpServer := startHTTPService(grpcLis, httpLis)
	return grpcServer, healthcheck, httpServer
}

func startGrpcService(lis net.Listener) (*grpc.Server, *health.Server) {
	log.Info("start grpc server ...")
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

	// register Trigger server
	db := mysql.New(configs.GetMySQLConfig())
	// create grpc dail
	conn, err := grpc.Dial(
		configs.GetSchedulerConsulURL(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		log.Fatalf("grpc Dial err: %v", err)
	}
	defer conn.Close()
	// create scheduler service client
	schedulerClient := scheduler.NewSchedulerServiceClient(conn)
	cronHandler := cron.NewHandler(cron.NewTrigger(db, schedulerClient))
	pb.RegisterCronTriggerServiceServer(server, cronHandler)
	webhookHandler := webhook.NewHandler(webhook.NewTrigger(db, schedulerClient))
	pb.RegisterWebhookTriggerServiceServer(server, webhookHandler)
	// register health check server
	healthcheck := health.NewServer()
	pbhealth.RegisterHealthServer(server, healthcheck)
	// Register reflection service on gRPC server.
	reflection.Register(server)

	go func() {
		// asynchronously inspect dependencies and toggle serving status as needed
		healthcheck.SetServingStatus(pb.CronTriggerService_ServiceDesc.ServiceName, pbhealth.HealthCheckResponse_SERVING)
		log.Infof("grpc server listening at %v", lis.Addr())
		if err := server.Serve(lis); err != nil {
			log.Fatalf("failed to start grpc serve: %v", err)
		}
	}()

	return server, healthcheck
}

func startHTTPService(grpcLis, httpLis net.Listener) *http.Server {
	log.Info("start http server ...")
	log.Infof(grpcLis.Addr().String())
	conn, err := grpc.DialContext(
		context.Background(),
		fmt.Sprintf("127.0.0.1:%d", grpcLis.Addr().(*net.TCPAddr).Port),
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("Failed to dial server: %v", err)
	}

	// gRPC-Gateway httpServer
	gwmux := runtime.NewServeMux()
	err = pb.RegisterCronTriggerServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalf("Failed to register gateway: %v", err)
	}

	// HTTP server config
	httpServer := &http.Server{
		Handler: gwmux,
	}

	go func() {
		time.Sleep(3 * time.Second)
		log.Infof("http server listening at %s", httpLis.Addr().String())
		if err = httpServer.Serve(httpLis); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				log.Info(err.Error())
				return
			}
			log.Fatalf("Failed to serve gRPC-Gateway: %v", err)
		}
	}()

	return httpServer
}
