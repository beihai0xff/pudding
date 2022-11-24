package cron

import (
	"context"

	"google.golang.org/protobuf/types/known/emptypb"

	pb "github.com/beihai0xff/pudding/api/gen/pudding/trigger/v1"
	"github.com/beihai0xff/pudding/api/gen/pudding/types/v1"
)

type Handler struct {
	pb.UnimplementedCronTriggerServiceServer
}

func (h Handler) Ping(ctx context.Context, empty *emptypb.Empty) (*pb.PingResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (h Handler) FindOneByID(ctx context.Context, request *pb.FindOneByIDRequest) (*pb.FindOneByIDResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (h Handler) PageQuery(ctx context.Context, request *pb.PageQueryRequest) (*pb.PageQueryResponse, error) {
	// TODO implement me
	panic("implement me")
}

func (h Handler) Register(ctx context.Context, request *pb.RegisterRequest) (*types.Response, error) {
	// TODO implement me
	panic("implement me")
}

func (h Handler) UpdateStatus(ctx context.Context, request *pb.UpdateStatusRequest) (*types.Response, error) {
	// TODO implement me
	panic("implement me")
}

func NewHandler() *Handler {
	return &Handler{}
}
