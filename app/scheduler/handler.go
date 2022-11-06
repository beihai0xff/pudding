package scheduler

import (
	"context"
	"errors"

	pb "github.com/beihai0xff/pudding/api/scheduler/v1"
	"github.com/beihai0xff/pudding/types"
)

type Handler struct {
	s Scheduler
	pb.UnimplementedPuddingServer
}

func NewHandler(s Scheduler) *Handler {
	return &Handler{s: s}
}

func (s *Handler) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingReply, error) {
	if req.Message != "ping" {
		return nil, errors.New("invalid message")
	}
	return &pb.PingReply{Message: "pong"}, nil
}

func (s *Handler) SendDelayMessage(ctx context.Context, req *pb.SendMessageRequest) (*pb.SendMessageReply, error) {
	if err := s.s.Produce(ctx, s.convPBToMessage(req)); err != nil {
		return &pb.SendMessageReply{}, nil
	}
	return &pb.SendMessageReply{}, nil
}

func (s *Handler) convPBToMessage(req *pb.SendMessageRequest) *types.Message {
	return &types.Message{
		Topic:     req.GetTopic(),
		Key:       req.GetKey(),
		Payload:   req.GetPayload(),
		Delay:     req.GetDelay(),
		ReadyTime: req.GetReadyTime(),
	}
}
