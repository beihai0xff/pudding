package scheduler

import (
	"context"
	"errors"

	pb "github.com/beihai0xff/pudding/api/gen/scheduler/v1"
	"github.com/beihai0xff/pudding/types"
)

type Handler struct {
	s Scheduler
	pb.UnimplementedSchedulerServiceServer
}

func NewHandler(s Scheduler) *Handler {
	return &Handler{s: s}
}

func (s *Handler) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	if req.Message != "ping" {
		return nil, errors.New("invalid message")
	}
	return &pb.PingResponse{Message: "pong"}, nil
}

func (s *Handler) SendDelayMessage(ctx context.Context, req *pb.SendDelayMessageRequest) (*pb.SendDelayMessageResponse,
	error) {
	if err := s.s.Produce(ctx, s.convPBToMessage(req)); err != nil {
		return &pb.SendDelayMessageResponse{}, nil
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
