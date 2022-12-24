package launcher

import (
	"context"
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

	pb "github.com/beihai0xff/pudding/api/gen/pudding/broker/v1"
	"github.com/beihai0xff/pudding/pkg/log"
	"github.com/beihai0xff/pudding/pkg/swagger"
)

type StartServiceFunc func(server *grpc.Server, serviceName *string) error

func StartGRPCService(grpcLis net.Listener, opts ...StartServiceFunc) (*grpc.Server, *health.Server) {
	log.Info("starting grpc server ...")
	// init grpc server
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

	// register health check server
	healthcheckServer := health.NewServer()
	pbhealth.RegisterHealthServer(server, healthcheckServer)

	for _, opt := range opts {
		serviceName := ""
		if err := opt(server, &serviceName); err != nil {
			log.Fatalf("failed to start service ")
		}
		// asynchronously inspect dependencies and toggle serving status as needed
		healthcheckServer.SetServingStatus(serviceName, pbhealth.HealthCheckResponse_SERVING)
	}

	// RegisterGRPC reflection service on gRPC server.
	// 提供该服务器端上可公开使用的 gRPC 服务的信息，
	// 服务反射向客户端提供了服务端注册的服务的信息，因此客户端不需要预编译服务定义就能与服务端交互
	// 通过此方式支持 grpcCRUL
	reflection.Register(server)

	go func() {

		log.Infof("grpc server listening at %v", grpcLis.Addr())
		if err := server.Serve(grpcLis); err != nil {
			log.Fatalf("failed to start grpc serve: %v", err)
		}
	}()

	return server, healthcheckServer
}

// StartHTTPService starts the HTTP service.
// It serves the gRPC-gateway, gRPC-healthz and the swagger UI.
func StartHTTPService(grpcLis, httpLis net.Listener, healthEndpointPath, swaggerEndpointPath string) *http.Server {
	log.Info("starting http server ...")
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
	gwmux := runtime.NewServeMux(runtime.WithHealthEndpointAt(pbhealth.NewHealthClient(conn), healthEndpointPath))
	swagger.RegisterHandler(gwmux, swaggerEndpointPath)

	err = pb.RegisterSchedulerServiceHandler(context.Background(), gwmux, conn)
	if err != nil {
		log.Fatalf("Failed to register gateway: %v", err)
	}

	// define HTTP server configuration
	httpServer := &http.Server{
		Handler: gwmux,
	}

	go func() {
		log.Infof("http server listening at %v", httpLis.Addr())
		if err = httpServer.Serve(httpLis); err != nil {
			if err == http.ErrServerClosed {
				log.Info("http server closed")
				return
			}
			log.Fatalf("Failed to serve gRPC-Gateway: %v", err)
		}
	}()

	return httpServer
}
