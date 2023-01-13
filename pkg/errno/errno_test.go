package errno

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestInternalError(t *testing.T) {
	// test no detail
	err := InternalError("can not find trigger by id")
	assert.EqualError(t, err, status.New(codes.Internal, "can not find trigger by id").Err().Error())
	sta, ok := status.FromError(err)
	assert.Equal(t, true, ok)
	assert.Equal(t, []interface{}{}, sta.Details())

	// test has one detail
	detail1 := errdetails.ErrorInfo{
		Reason:   "test_reason",
		Domain:   "detail1",
		Metadata: map[string]string{"id": "0"},
	}
	err = InternalError("can not find trigger by id", &detail1)
	assert.EqualError(t, err, status.New(codes.Internal, "can not find trigger by id").Err().Error())
	sta, ok = status.FromError(err)
	if ok {
		assert.Equal(t, detail1.String(), sta.Details()[0].(*errdetails.ErrorInfo).String())
	} else {
		t.Errorf("can not convert error to status")
	}

	// test has one detail
	detail2 := detail1
	detail2.Reason = "detail2"
	err = InternalError("can not find trigger by id", &detail1, &detail2)
	assert.EqualError(t, err, status.New(codes.Internal, "can not find trigger by id").Err().Error())
	sta, ok = status.FromError(err)
	if ok {
		assert.Equal(t, detail1.String(), sta.Details()[0].(*errdetails.ErrorInfo).String())
		assert.Equal(t, detail2.String(), sta.Details()[1].(*errdetails.ErrorInfo).String())
	} else {
		t.Errorf("can not convert error to status")
	}

}

func TestBadRequest(t *testing.T) {
	// test no detail
	err := BadRequest("can not find trigger by id")
	assert.EqualError(t, err, status.New(codes.InvalidArgument, "can not find trigger by id").Err().Error())
	sta, ok := status.FromError(err)
	assert.Equal(t, true, ok)
	assert.Equal(t, []interface{}{}, sta.Details())

	// test has one detail
	detail1 := errdetails.ErrorInfo{
		Reason:   "test_reason",
		Domain:   "detail1",
		Metadata: map[string]string{"id": "0"},
	}
	err = BadRequest("can not find trigger by id", &detail1)
	assert.EqualError(t, err, status.New(codes.InvalidArgument, "can not find trigger by id").Err().Error())
	sta, ok = status.FromError(err)
	if ok {
		assert.Equal(t, detail1.String(), sta.Details()[0].(*errdetails.ErrorInfo).String())
	} else {
		t.Errorf("can not convert error to status")
	}

	// test has one detail
	detail2 := detail1
	detail2.Reason = "detail2"
	err = BadRequest("can not find trigger by id", &detail1, &detail2)
	assert.EqualError(t, err, status.New(codes.InvalidArgument, "can not find trigger by id").Err().Error())
	sta, ok = status.FromError(err)
	if ok {
		assert.Equal(t, detail1.String(), sta.Details()[0].(*errdetails.ErrorInfo).String())
		assert.Equal(t, detail2.String(), sta.Details()[1].(*errdetails.ErrorInfo).String())
	} else {
		t.Errorf("can not convert error to status")
	}

}
