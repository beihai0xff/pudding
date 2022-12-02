package cron

import (
	"context"
	"fmt"
	"strconv"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/beihai0xff/pudding/api/gen/pudding/trigger/v1"
	"github.com/beihai0xff/pudding/api/gen/pudding/types/v1"
	"github.com/beihai0xff/pudding/app/trigger/entity"
	"github.com/beihai0xff/pudding/pkg/errno"
)

const cronTriggerDomain = "pudding.trigger.cron"

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

func (h *Handler) FindOneByID(ctx context.Context, req *pb.FindOneByIDRequest) (*pb.CronFindOneByIDResponse, error) {
	if req.Id <= 0 {
		return nil, errno.BadRequest("Invalid ID", &errdetails.BadRequest_FieldViolation{
			Field:       "ID",
			Description: fmt.Sprintf("ID [%d] should be greater than zero", req.Id),
		})
	}
	e, err := h.t.FindByID(ctx, uint(req.Id))
	if err != nil {
		return nil, errno.InternalError("can not find trigger by id", &errdetails.ErrorInfo{
			Reason:   err.Error(),
			Domain:   cronTriggerDomain,
			Metadata: map[string]string{"id": strconv.FormatUint(req.Id, 10)},
		})
	}
	return &pb.CronFindOneByIDResponse{
		Body: h.convertTemplateEntityToPb(e),
	}, nil
}

func (h *Handler) PageQuery(ctx context.Context, req *pb.PageQueryTemplateRequest) (*pb.CronPageQueryResponse, error) {
	p := entity.PageQuery{
		Offset: int(req.Offset),
		Limit:  int(req.Limit),
	}
	res, count, err := h.t.PageQuery(ctx, &p, req.Status)
	if err != nil {
		return nil, errno.InternalError("can not pageQuery cron trigger", &errdetails.ErrorInfo{
			Reason: err.Error(),
			Domain: cronTriggerDomain,
			Metadata: map[string]string{"offset": strconv.FormatUint(req.Offset, 10),
				"limit": strconv.FormatUint(req.Limit, 10)},
		})
	}

	return &pb.CronPageQueryResponse{
		Count: uint64(count),
		Body:  h.convertTemplateEntitySliceToPb(res),
	}, nil
}

func (h *Handler) Register(ctx context.Context, req *pb.CronTriggerServiceRegisterRequest) (*emptypb.Empty, error) {
	e := &entity.CronTriggerTemplate{
		CronExpr: req.CronExpr,
		Topic:    req.Topic,
		Payload:  req.Payload,
	}
	if err := h.t.Register(ctx, e); err != nil {
		return nil, errno.InternalError("can not register trigger", &errdetails.ErrorInfo{
			Reason:   err.Error(),
			Domain:   cronTriggerDomain,
			Metadata: map[string]string{"request body": req.String()},
		})
	}
	return &emptypb.Empty{}, nil
}

func (h *Handler) UpdateStatus(ctx context.Context, req *pb.UpdateStatusRequest) (*pb.UpdateStatusResponse, error) {
	// check params
	if req.Id <= 0 {
		return nil, errno.BadRequest("Invalid ID", &errdetails.BadRequest_FieldViolation{
			Field:       "ID",
			Description: fmt.Sprintf("ID [%d] should be greater than zero", req.Id),
		})
	}
	if req.Status > pb.TriggerStatus_MAX_AGE || req.Status <= pb.TriggerStatus_UNKNOWN_UNSPECIFIED {
		return nil, errno.BadRequest("Invalid status code", &errdetails.BadRequest_FieldViolation{
			Field:       "status",
			Description: fmt.Sprintf("Invalid status code [%d], please use proto define status code", req.Status),
		})
	}

	// update status
	rowsAffected, err := h.t.UpdateStatus(ctx, uint(req.Id), req.Status)
	if err != nil {
		return nil, errno.InternalError("can not update trigger by id", &errdetails.ErrorInfo{
			Reason: err.Error(),
			Domain: cronTriggerDomain,
			Metadata: map[string]string{"id": strconv.FormatUint(req.Id, 10),
				"status": strconv.FormatInt(int64(req.Status), 10)},
		})
	}

	return &pb.UpdateStatusResponse{RowsAffected: rowsAffected}, nil
}

func (h *Handler) convertTemplateEntityToPb(e *entity.CronTriggerTemplate) *pb.CronTriggerTemplate {
	return &pb.CronTriggerTemplate{
		Id:                uint64(e.ID),
		CronExpr:          e.CronExpr,
		Topic:             e.Topic,
		Payload:           e.Payload,
		LastExecutionTime: timestamppb.New(e.LastExecutionTime),
		LoopedTimes:       e.LoopedTimes,
		ExceptedEndTime:   timestamppb.New(e.ExceptedEndTime),
		ExceptedLoopTimes: e.ExceptedLoopTimes,
		Status:            e.Status,
	}
}

func (h *Handler) convertTemplateEntitySliceToPb(es []*entity.CronTriggerTemplate) []*pb.CronTriggerTemplate {
	res := make([]*pb.CronTriggerTemplate, len(es))
	for _, e := range es {
		res = append(res, h.convertTemplateEntityToPb(e))
	}
	return res
}
