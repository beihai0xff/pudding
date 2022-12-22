package broker

import (
	"context"
	"fmt"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/beihai0xff/pudding/api/gen/pudding/broker/v1"
	"github.com/beihai0xff/pudding/api/gen/pudding/types/v1"
	"github.com/beihai0xff/pudding/pkg/errno"
)

const puddingScheduler = "pudding.scheduler"

type Handler struct {
	s Scheduler
	pb.UnimplementedSchedulerServiceServer
}

func NewHandler(s Scheduler) *Handler {
	return &Handler{s: s, UnimplementedSchedulerServiceServer: pb.UnimplementedSchedulerServiceServer{}}
}

func (s *Handler) SendDelayMessage(ctx context.Context, req *pb.SendDelayMessageRequest) (*emptypb.Empty, error) {
	msg := s.convPBToMessage(req)
	if msg.DeliverAt <= 0 && msg.DeliverAfter <= 0 {
		return nil, errno.BadRequest("deliver_at and deliver_after can't be both zero",
			&errdetails.BadRequest_FieldViolation{
				Field:       "deliver_after",
				Description: fmt.Sprintf("deliver_after [%d] should be greater than zero", req.DeliverAfter),
			},
			&errdetails.BadRequest_FieldViolation{
				Field:       "deliver_at",
				Description: fmt.Sprintf("deliver_at [%d] should be greater than zero", req.DeliverAt),
			})
	}

	if err := s.s.Produce(ctx, msg); err != nil {
		return nil, errno.InternalError("can not produce message", &errdetails.ErrorInfo{
			Reason:   err.Error(),
			Domain:   puddingScheduler,
			Metadata: map[string]string{"body": req.String()},
		})
	}
	return &emptypb.Empty{}, nil
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
