package errno

import (
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func InternalError(message string, details ...proto.Message) error {
	stat, _ := status.New(codes.Internal, message).WithDetails(details...)
	return stat.Err()
}

func BadRequest(message string, details ...proto.Message) error {
	stat, _ := status.New(codes.InvalidArgument, message).WithDetails(details...)
	return stat.Err()
}
