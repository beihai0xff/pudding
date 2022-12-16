// Package mock is a GoMock package.
package mock

import (
	context "context"

	grpc "google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"

	scheduler "github.com/beihai0xff/pudding/api/gen/pudding/scheduler/v1"
	types "github.com/beihai0xff/pudding/api/gen/pudding/types/v1"
)

// MockSchedulerServiceClient is a mock of SchedulerServiceClient interface.
type MockSchedulerServiceClient struct {
}

// NewMockSchedulerServiceClient creates a new mock instance.
func NewMockSchedulerServiceClient() *MockSchedulerServiceClient {
	return &MockSchedulerServiceClient{}
}

// Ping mocks base method.
func (m *MockSchedulerServiceClient) Ping(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*types.PingResponse, error) {
	return &types.PingResponse{Message: "pong"}, nil
}

// SendDelayMessage mocks base method.
func (m *MockSchedulerServiceClient) SendDelayMessage(ctx context.Context, in *scheduler.SendDelayMessageRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
