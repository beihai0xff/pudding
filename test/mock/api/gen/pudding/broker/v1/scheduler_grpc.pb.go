// Package mock is a GoMock package.
package mock

import (
	context "context"

	grpc "google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"

	"github.com/beihai0xff/pudding/api/gen/pudding/broker/v1"
)

// MockSchedulerServiceClient is a mock of SchedulerServiceClient interface.
type MockSchedulerServiceClient struct {
}

// NewMockSchedulerServiceClient creates a new mock instance.
func NewMockSchedulerServiceClient() *MockSchedulerServiceClient {
	return &MockSchedulerServiceClient{}
}

// SendDelayMessage mocks base method.
func (m *MockSchedulerServiceClient) SendDelayMessage(ctx context.Context, in *broker.SendDelayMessageRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
