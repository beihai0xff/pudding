package scheduler

import (
	"context"

	pb "github.com/beihai0xff/pudding/api/scheduler/v1"

	"github.com/beihai0xff/pudding/configs"
	"github.com/beihai0xff/pudding/pkg/lock"
)

type Handler struct {
	s Scheduler
}

func NewHandler() {
	_ = New(configs.GetSchedulerConfig())
	lock.Init()

}

func (s *Handler) GetFeature(ctx context.Context, point *pb.PingReply) (*pb.PingRequest, error) {
	return nil, nil
}
