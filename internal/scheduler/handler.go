package scheduler

import (
	"context"

	pb "github.com/beihai0xff/pudding/api/scheduler/v1"
)

type Handler struct {
	s Scheduler
}

func (s *Handler) GetFeature(ctx context.Context, point *pb.PingReply) (*pb.PingRequest, error) {
	return nil, nil
}
