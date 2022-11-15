package scheduler

import (
	"context"

	codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/beihai0xff/pudding/api/gen/pudding/scheduler/v1"
	"github.com/beihai0xff/pudding/types"
)

type Handler struct {
	s Scheduler
	pb.UnimplementedSchedulerServiceServer
}

func NewHandler(s Scheduler) *Handler {
	return &Handler{s: s}
}

func (s *Handler) Ping(ctx context.Context, req *emptypb.Empty) (*pb.PingResponse, error) {
	return &pb.PingResponse{Message: "pong"}, nil
}

func (s *Handler) SendDelayMessage(ctx context.Context, req *pb.SendDelayMessageRequest) (*pb.SendDelayMessageResponse,
	error) {
	msg := s.convPBToMessage(req)
	if msg.DeliverAt == 0 && msg.DeliverAfter == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "DeliverAt and DeliverAfter can't be both zero")
	}

	if err := s.s.Produce(ctx, msg); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.SendDelayMessageResponse{}, nil
}

func (s *Handler) convPBToMessage(req *pb.SendDelayMessageRequest) *types.Message {
	return &types.Message{
		Topic:        req.GetTopic(),
		Key:          req.GetKey(),
		Payload:      req.GetPayload(),
		DeliverAfter: req.GetDeliverAfter(),
		DeliverAt:    req.GetDeliverAt(),
	}
}
