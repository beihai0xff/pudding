// Package launcher provides a launcher to start gRPC server, health server and grpc gateway server.
// interceptor.go provides interceptor to handle gRPC requests.
package launcher

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/samber/lo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"

	"github.com/beihai0xff/pudding/pkg/log/logger"
)

var grpcBlackList = []string{"/grpc.health.v1.Health/Check"}

// unaryServerRequestLog logs the grpc request.
func unaryServerRequestLog() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		rsp, err := handler(ctx, req)
		if !lo.Contains(grpcBlackList, info.FullMethod) {
			recordGRPCRequestLog(ctx, req, rsp, info, err)
		}

		return rsp, err
	}
}

func recordGRPCRequestLog(ctx context.Context, req, rsp interface{}, info *grpc.UnaryServerInfo, err error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		logger.GetGRPCLogger().Errorf("failed to get peer from context")
		return
	}
	request := map[string]interface{}{
		"method":      info.FullMethod,
		"remote_addr": p.Addr.String(),
		"request":     req,
		"response":    rsp,
	}

	if err != nil {
		s, ok := status.FromError(err)
		if !ok {
			logger.GetGRPCLogger().Errorf("failed to get grpc status from error")
			return
		}
		request["response"] = fmt.Sprintf("%s, details:%+v", s.String(), s.Details())
		b, _ := json.Marshal(request)
		logger.GetGRPCLogger().Error(string(b))
		return
	}
	b, _ := json.Marshal(request)
	logger.GetGRPCLogger().Info(string(b))
}
