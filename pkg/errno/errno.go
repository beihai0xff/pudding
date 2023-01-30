// Package errno implements errors returned by gRPC. These errors are
// serialized and transmitted on the wire between server and client, and allow
// for additional data to be transmitted via the Details field in the status
package errno

import (
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// InternalError returns an Internal error representing message and details.
func InternalError(message string, details ...proto.Message) error {
	stat, _ := status.New(codes.Internal, message).WithDetails(details...)
	return stat.Err()
}

// BadRequest returns an Bad Request error representing message and details.
func BadRequest(message string, details ...proto.Message) error {
	stat, _ := status.New(codes.InvalidArgument, message).WithDetails(details...)
	return stat.Err()
}
