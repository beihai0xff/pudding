package cron

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/beihai0xff/pudding/api/gen/pudding/trigger/v1"
	"github.com/beihai0xff/pudding/api/gen/pudding/types/v1"
	"github.com/beihai0xff/pudding/app/trigger/entity"
)

type Handler struct {
	t *Trigger
	pb.UnimplementedCronTriggerServiceServer
}

func NewHandler(t *Trigger) *Handler {
	return &Handler{
		t:                                     t,
		UnimplementedCronTriggerServiceServer: pb.UnimplementedCronTriggerServiceServer{},
	}
}

func (h *Handler) Ping(context.Context, *emptypb.Empty) (*types.PingResponse, error) {
	return &types.PingResponse{Message: "pong"}, nil
}

func (h *Handler) FindOneByID(ctx context.Context, request *pb.FindOneByIDRequest) (*pb.FindOneByIDResponse, error) {
	e, err := h.t.FindByID(ctx, uint(request.Id))
	if err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &pb.FindOneByIDResponse{
		Code:    0,
		Message: "OK",
		Body:    h.convertTemplateEntityToPb(e),
	}, nil
}

func (h *Handler) PageQuery(ctx context.Context, request *pb.PageQueryRequest) (*pb.PageQueryResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (h *Handler) Register(ctx context.Context, request *pb.RegisterRequest) (*types.Response, error) {
	// TODO implement me
	panic("implement me")
}

func (h *Handler) UpdateStatus(ctx context.Context, request *pb.UpdateStatusRequest) (*types.Response, error) {
	if err := h.t.UpdateStatus(ctx, uint(request.Id), request.Status); err != nil {
		return nil, status.Errorf(codes.Internal, err.Error())
	}
	return &types.Response{
		Code:    0,
		Message: "OK",
	}, nil
}

func (h *Handler) convertTemplateEntityToPb(e *entity.CronTriggerTemplate) *pb.TriggerTemplate {
	return &pb.TriggerTemplate{
		Id:                uint64(e.ID),
		CronExpr:          e.CronExpr,
		Topic:             e.Topic,
		Payload:           e.Payload,
		LastExecutionTime: timestamppb.New(e.LastExecutionTime),
		ExceptedEndTime:   timestamppb.New(e.ExceptedEndTime),
		LoopedTimes:       e.LoopedTimes,
		Status:            e.Status,
	}
}
