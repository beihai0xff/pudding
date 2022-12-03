package webhook

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

const webhookTriggerDomain = "pudding.trigger.webhook"

type Handler struct {
	t *Trigger
	pb.UnimplementedWebhookTriggerServiceServer
}

func NewHandler(t *Trigger) *Handler {
	return &Handler{
		t:                                        t,
		UnimplementedWebhookTriggerServiceServer: pb.UnimplementedWebhookTriggerServiceServer{},
	}
}

func (h *Handler) Ping(context.Context, *emptypb.Empty) (*types.PingResponse, error) {
	return &types.PingResponse{Message: "pong"}, nil
}

func (h *Handler) FindOneByID(ctx context.Context, req *pb.FindOneByIDRequest) (*pb.WebhookFindOneByIDResponse, error) {
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
			Domain:   webhookTriggerDomain,
			Metadata: map[string]string{"id": strconv.FormatUint(req.Id, 10)},
		})
	}
	return &pb.WebhookFindOneByIDResponse{
		Body: h.convertTemplateEntityToPb(e),
	}, nil
}

func (h *Handler) PageQueryTemplate(ctx context.Context, req *pb.PageQueryTemplateRequest) (*pb.WebhookPageQueryResponse, error) {
	p := entity.PageQuery{
		Offset: int(req.Offset),
		Limit:  int(req.Limit),
	}
	res, count, err := h.t.PageQuery(ctx, &p, req.Status)
	if err != nil {
		return nil, errno.InternalError("can not pageQuery Webhook trigger", &errdetails.ErrorInfo{
			Reason: err.Error(),
			Domain: webhookTriggerDomain,
			Metadata: map[string]string{"offset": strconv.FormatUint(req.Offset, 10),
				"limit": strconv.FormatUint(req.Limit, 10)},
		})
	}

	return &pb.WebhookPageQueryResponse{
		Count: uint64(count),
		Body:  h.convertTemplateEntitySliceToPb(res),
	}, nil
}

func (h *Handler) Register(ctx context.Context,
	req *pb.WebhookTriggerServiceRegisterRequest) (*pb.WebhookRegisterResponse, error) {

	e := &entity.WebhookTriggerTemplate{
		Topic:   req.Topic,
		Payload: req.Payload,
	}
	if err := h.t.Register(ctx, e); err != nil {
		return nil, errno.InternalError("can not register trigger", &errdetails.ErrorInfo{
			Reason:   err.Error(),
			Domain:   webhookTriggerDomain,
			Metadata: map[string]string{"request body": req.String()},
		})
	}
	return &pb.WebhookRegisterResponse{
		Url: h.t.genWebhookURL(e.ID),
	}, nil
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
			Domain: webhookTriggerDomain,
			Metadata: map[string]string{"id": strconv.FormatUint(req.Id, 10),
				"status": strconv.FormatInt(int64(req.Status), 10)},
		})
	}

	return &pb.UpdateStatusResponse{RowsAffected: rowsAffected}, nil
}

func (h *Handler) convertTemplateEntityToPb(e *entity.WebhookTriggerTemplate) *pb.WebhookTriggerTemplate {
	return &pb.WebhookTriggerTemplate{
		Id:                uint64(e.ID),
		Topic:             e.Topic,
		Payload:           e.Payload,
		LoopedTimes:       e.LoopedTimes,
		ExceptedEndTime:   timestamppb.New(e.ExceptedEndTime),
		ExceptedLoopTimes: e.ExceptedLoopTimes,
		Status:            e.Status,
	}
}

func (h *Handler) convertTemplateEntitySliceToPb(es []*entity.WebhookTriggerTemplate) []*pb.WebhookTriggerTemplate {
	res := make([]*pb.WebhookTriggerTemplate, len(es))
	for _, e := range es {
		res = append(res, h.convertTemplateEntityToPb(e))
	}
	return res
}
