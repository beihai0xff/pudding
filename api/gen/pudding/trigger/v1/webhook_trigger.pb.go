// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.29.0
// 	protoc        (unknown)
// source: pudding/trigger/v1/webhook_trigger.proto

package trigger

import (
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2/options"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// WebhookFindOneByIDRequest The FindOneByID Response message.
type WebhookFindOneByIDResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Body *WebhookTriggerTemplate `protobuf:"bytes,1,opt,name=body,proto3" json:"body,omitempty"`
}

func (x *WebhookFindOneByIDResponse) Reset() {
	*x = WebhookFindOneByIDResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pudding_trigger_v1_webhook_trigger_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WebhookFindOneByIDResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WebhookFindOneByIDResponse) ProtoMessage() {}

func (x *WebhookFindOneByIDResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pudding_trigger_v1_webhook_trigger_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WebhookFindOneByIDResponse.ProtoReflect.Descriptor instead.
func (*WebhookFindOneByIDResponse) Descriptor() ([]byte, []int) {
	return file_pudding_trigger_v1_webhook_trigger_proto_rawDescGZIP(), []int{0}
}

func (x *WebhookFindOneByIDResponse) GetBody() *WebhookTriggerTemplate {
	if x != nil {
		return x.Body
	}
	return nil
}

type WebhookPageQueryResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Count uint64                    `protobuf:"varint,1,opt,name=count,proto3" json:"count,omitempty"`
	Body  []*WebhookTriggerTemplate `protobuf:"bytes,2,rep,name=body,proto3" json:"body,omitempty"`
}

func (x *WebhookPageQueryResponse) Reset() {
	*x = WebhookPageQueryResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pudding_trigger_v1_webhook_trigger_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WebhookPageQueryResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WebhookPageQueryResponse) ProtoMessage() {}

func (x *WebhookPageQueryResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pudding_trigger_v1_webhook_trigger_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WebhookPageQueryResponse.ProtoReflect.Descriptor instead.
func (*WebhookPageQueryResponse) Descriptor() ([]byte, []int) {
	return file_pudding_trigger_v1_webhook_trigger_proto_rawDescGZIP(), []int{1}
}

func (x *WebhookPageQueryResponse) GetCount() uint64 {
	if x != nil {
		return x.Count
	}
	return 0
}

func (x *WebhookPageQueryResponse) GetBody() []*WebhookTriggerTemplate {
	if x != nil {
		return x.Body
	}
	return nil
}

type WebhookTriggerServiceRegisterRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Topic             string                 `protobuf:"bytes,1,opt,name=topic,proto3" json:"topic,omitempty"`
	Payload           []byte                 `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
	DeliverAfter      uint64                 `protobuf:"varint,3,opt,name=deliver_after,json=deliverAfter,proto3" json:"deliver_after,omitempty"`
	ExceptedEndTime   *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=excepted_end_time,json=exceptedEndTime,proto3" json:"excepted_end_time,omitempty"`
	ExceptedLoopTimes uint64                 `protobuf:"varint,5,opt,name=excepted_loop_times,json=exceptedLoopTimes,proto3" json:"excepted_loop_times,omitempty"`
}

func (x *WebhookTriggerServiceRegisterRequest) Reset() {
	*x = WebhookTriggerServiceRegisterRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pudding_trigger_v1_webhook_trigger_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WebhookTriggerServiceRegisterRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WebhookTriggerServiceRegisterRequest) ProtoMessage() {}

func (x *WebhookTriggerServiceRegisterRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pudding_trigger_v1_webhook_trigger_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WebhookTriggerServiceRegisterRequest.ProtoReflect.Descriptor instead.
func (*WebhookTriggerServiceRegisterRequest) Descriptor() ([]byte, []int) {
	return file_pudding_trigger_v1_webhook_trigger_proto_rawDescGZIP(), []int{2}
}

func (x *WebhookTriggerServiceRegisterRequest) GetTopic() string {
	if x != nil {
		return x.Topic
	}
	return ""
}

func (x *WebhookTriggerServiceRegisterRequest) GetPayload() []byte {
	if x != nil {
		return x.Payload
	}
	return nil
}

func (x *WebhookTriggerServiceRegisterRequest) GetDeliverAfter() uint64 {
	if x != nil {
		return x.DeliverAfter
	}
	return 0
}

func (x *WebhookTriggerServiceRegisterRequest) GetExceptedEndTime() *timestamppb.Timestamp {
	if x != nil {
		return x.ExceptedEndTime
	}
	return nil
}

func (x *WebhookTriggerServiceRegisterRequest) GetExceptedLoopTimes() uint64 {
	if x != nil {
		return x.ExceptedLoopTimes
	}
	return 0
}

type WebhookRegisterResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// url is webhook url
	Url string `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
}

func (x *WebhookRegisterResponse) Reset() {
	*x = WebhookRegisterResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pudding_trigger_v1_webhook_trigger_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WebhookRegisterResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WebhookRegisterResponse) ProtoMessage() {}

func (x *WebhookRegisterResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pudding_trigger_v1_webhook_trigger_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WebhookRegisterResponse.ProtoReflect.Descriptor instead.
func (*WebhookRegisterResponse) Descriptor() ([]byte, []int) {
	return file_pudding_trigger_v1_webhook_trigger_proto_rawDescGZIP(), []int{3}
}

func (x *WebhookRegisterResponse) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

type WebhookTriggerServiceCallRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// trigger template ID
	Id uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *WebhookTriggerServiceCallRequest) Reset() {
	*x = WebhookTriggerServiceCallRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pudding_trigger_v1_webhook_trigger_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WebhookTriggerServiceCallRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WebhookTriggerServiceCallRequest) ProtoMessage() {}

func (x *WebhookTriggerServiceCallRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pudding_trigger_v1_webhook_trigger_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WebhookTriggerServiceCallRequest.ProtoReflect.Descriptor instead.
func (*WebhookTriggerServiceCallRequest) Descriptor() ([]byte, []int) {
	return file_pudding_trigger_v1_webhook_trigger_proto_rawDescGZIP(), []int{4}
}

func (x *WebhookTriggerServiceCallRequest) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type WebhookTriggerServiceCallResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// message_id delay message id
	MessageKey string `protobuf:"bytes,1,opt,name=message_key,json=messageKey,proto3" json:"message_key,omitempty"`
}

func (x *WebhookTriggerServiceCallResponse) Reset() {
	*x = WebhookTriggerServiceCallResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pudding_trigger_v1_webhook_trigger_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WebhookTriggerServiceCallResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WebhookTriggerServiceCallResponse) ProtoMessage() {}

func (x *WebhookTriggerServiceCallResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pudding_trigger_v1_webhook_trigger_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WebhookTriggerServiceCallResponse.ProtoReflect.Descriptor instead.
func (*WebhookTriggerServiceCallResponse) Descriptor() ([]byte, []int) {
	return file_pudding_trigger_v1_webhook_trigger_proto_rawDescGZIP(), []int{5}
}

func (x *WebhookTriggerServiceCallResponse) GetMessageKey() string {
	if x != nil {
		return x.MessageKey
	}
	return ""
}

var File_pudding_trigger_v1_webhook_trigger_proto protoreflect.FileDescriptor

var file_pudding_trigger_v1_webhook_trigger_proto_rawDesc = []byte{
	0x0a, 0x28, 0x70, 0x75, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x2f, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65,
	0x72, 0x2f, 0x76, 0x31, 0x2f, 0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x5f, 0x74, 0x72, 0x69,
	0x67, 0x67, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x12, 0x70, 0x75, 0x64, 0x64,
	0x69, 0x6e, 0x67, 0x2e, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x1a, 0x1c,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1f, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69,
	0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x63, 0x2d, 0x67, 0x65, 0x6e, 0x2d, 0x6f, 0x70, 0x65, 0x6e, 0x61, 0x70,
	0x69, 0x76, 0x32, 0x2f, 0x6f, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2f, 0x61, 0x6e, 0x6e, 0x6f,
	0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x70,
	0x75, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x2f, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x2f, 0x76,
	0x31, 0x2f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x5c, 0x0a,
	0x1a, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x46, 0x69, 0x6e, 0x64, 0x4f, 0x6e, 0x65, 0x42,
	0x79, 0x49, 0x44, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3e, 0x0a, 0x04, 0x62,
	0x6f, 0x64, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2a, 0x2e, 0x70, 0x75, 0x64, 0x64,
	0x69, 0x6e, 0x67, 0x2e, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x57,
	0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x54, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x54, 0x65, 0x6d,
	0x70, 0x6c, 0x61, 0x74, 0x65, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x22, 0x70, 0x0a, 0x18, 0x57,
	0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x50, 0x61, 0x67, 0x65, 0x51, 0x75, 0x65, 0x72, 0x79, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x3e, 0x0a,
	0x04, 0x62, 0x6f, 0x64, 0x79, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2a, 0x2e, 0x70, 0x75,
	0x64, 0x64, 0x69, 0x6e, 0x67, 0x2e, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x2e, 0x76, 0x31,
	0x2e, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x54, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x54,
	0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x52, 0x04, 0x62, 0x6f, 0x64, 0x79, 0x22, 0xf3, 0x01,
	0x0a, 0x24, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x54, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x6f, 0x70, 0x69, 0x63, 0x12, 0x18, 0x0a, 0x07,
	0x70, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x70,
	0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x23, 0x0a, 0x0d, 0x64, 0x65, 0x6c, 0x69, 0x76, 0x65,
	0x72, 0x5f, 0x61, 0x66, 0x74, 0x65, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0c, 0x64,
	0x65, 0x6c, 0x69, 0x76, 0x65, 0x72, 0x41, 0x66, 0x74, 0x65, 0x72, 0x12, 0x46, 0x0a, 0x11, 0x65,
	0x78, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x5f, 0x65, 0x6e, 0x64, 0x5f, 0x74, 0x69, 0x6d, 0x65,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61,
	0x6d, 0x70, 0x52, 0x0f, 0x65, 0x78, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x45, 0x6e, 0x64, 0x54,
	0x69, 0x6d, 0x65, 0x12, 0x2e, 0x0a, 0x13, 0x65, 0x78, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x5f,
	0x6c, 0x6f, 0x6f, 0x70, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x18, 0x05, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x11, 0x65, 0x78, 0x63, 0x65, 0x70, 0x74, 0x65, 0x64, 0x4c, 0x6f, 0x6f, 0x70, 0x54, 0x69,
	0x6d, 0x65, 0x73, 0x22, 0x2b, 0x0a, 0x17, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x52, 0x65,
	0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x10,
	0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c,
	0x22, 0x32, 0x0a, 0x20, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x54, 0x72, 0x69, 0x67, 0x67,
	0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x43, 0x61, 0x6c, 0x6c, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x02, 0x69, 0x64, 0x22, 0x44, 0x0a, 0x21, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x54,
	0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x43, 0x61, 0x6c,
	0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x6d, 0x65, 0x73,
	0x73, 0x61, 0x67, 0x65, 0x5f, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a,
	0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x4b, 0x65, 0x79, 0x32, 0xa2, 0x0a, 0x0a, 0x15, 0x57,
	0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x54, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x53, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x12, 0xf7, 0x01, 0x0a, 0x0b, 0x46, 0x69, 0x6e, 0x64, 0x4f, 0x6e, 0x65,
	0x42, 0x79, 0x49, 0x44, 0x12, 0x26, 0x2e, 0x70, 0x75, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x2e, 0x74,
	0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x69, 0x6e, 0x64, 0x4f, 0x6e,
	0x65, 0x42, 0x79, 0x49, 0x44, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2e, 0x2e, 0x70,
	0x75, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x2e, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x2e, 0x76,
	0x31, 0x2e, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x46, 0x69, 0x6e, 0x64, 0x4f, 0x6e, 0x65,
	0x42, 0x79, 0x49, 0x44, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x8f, 0x01, 0x92,
	0x41, 0x5b, 0x0a, 0x0f, 0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x5f, 0x74, 0x72, 0x69, 0x67,
	0x67, 0x65, 0x72, 0x12, 0x23, 0x66, 0x69, 0x6e, 0x64, 0x20, 0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f,
	0x6b, 0x20, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x20, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61,
	0x74, 0x65, 0x20, 0x62, 0x79, 0x20, 0x69, 0x64, 0x1a, 0x23, 0x66, 0x69, 0x6e, 0x64, 0x20, 0x77,
	0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x20, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x20, 0x74,
	0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x20, 0x62, 0x79, 0x20, 0x69, 0x64, 0x82, 0xd3, 0xe4,
	0x93, 0x02, 0x2b, 0x12, 0x29, 0x2f, 0x70, 0x75, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x2f, 0x74, 0x72,
	0x69, 0x67, 0x67, 0x65, 0x72, 0x2f, 0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x2f, 0x76, 0x31,
	0x2f, 0x66, 0x69, 0x6e, 0x64, 0x5f, 0x6f, 0x6e, 0x65, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x12, 0xa3,
	0x02, 0x0a, 0x11, 0x50, 0x61, 0x67, 0x65, 0x51, 0x75, 0x65, 0x72, 0x79, 0x54, 0x65, 0x6d, 0x70,
	0x6c, 0x61, 0x74, 0x65, 0x12, 0x2c, 0x2e, 0x70, 0x75, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x2e, 0x74,
	0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x61, 0x67, 0x65, 0x51, 0x75,
	0x65, 0x72, 0x79, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x2c, 0x2e, 0x70, 0x75, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x2e, 0x74, 0x72, 0x69,
	0x67, 0x67, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x50,
	0x61, 0x67, 0x65, 0x51, 0x75, 0x65, 0x72, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0xb1, 0x01, 0x92, 0x41, 0x5d, 0x0a, 0x0f, 0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x5f,
	0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x12, 0x24, 0x70, 0x61, 0x67, 0x65, 0x20, 0x71, 0x75,
	0x65, 0x72, 0x79, 0x20, 0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x20, 0x74, 0x72, 0x69, 0x67,
	0x67, 0x65, 0x72, 0x20, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x73, 0x1a, 0x24, 0x70,
	0x61, 0x67, 0x65, 0x20, 0x71, 0x75, 0x65, 0x72, 0x79, 0x20, 0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f,
	0x6b, 0x20, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x20, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61,
	0x74, 0x65, 0x73, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x4b, 0x12, 0x49, 0x2f, 0x70, 0x75, 0x64, 0x64,
	0x69, 0x6e, 0x67, 0x2f, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x2f, 0x77, 0x65, 0x62, 0x68,
	0x6f, 0x6f, 0x6b, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x71, 0x75, 0x65, 0x72,
	0x79, 0x2f, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x2f, 0x7b, 0x6f, 0x66, 0x66, 0x73,
	0x65, 0x74, 0x7d, 0x2f, 0x7b, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x7d, 0x2f, 0x7b, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x7d, 0x12, 0xfd, 0x01, 0x0a, 0x08, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65,
	0x72, 0x12, 0x38, 0x2e, 0x70, 0x75, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x2e, 0x74, 0x72, 0x69, 0x67,
	0x67, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x54, 0x72,
	0x69, 0x67, 0x67, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x52, 0x65, 0x67, 0x69,
	0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x2b, 0x2e, 0x70, 0x75,
	0x64, 0x64, 0x69, 0x6e, 0x67, 0x2e, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x2e, 0x76, 0x31,
	0x2e, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x89, 0x01, 0x92, 0x41, 0x57, 0x0a, 0x0f,
	0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x5f, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x12,
	0x21, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x20, 0x61, 0x20, 0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f,
	0x6b, 0x20, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x20, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61,
	0x74, 0x65, 0x1a, 0x21, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65, 0x20, 0x61, 0x20, 0x77, 0x65, 0x62,
	0x68, 0x6f, 0x6f, 0x6b, 0x20, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x20, 0x74, 0x65, 0x6d,
	0x70, 0x6c, 0x61, 0x74, 0x65, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x29, 0x3a, 0x01, 0x2a, 0x1a, 0x24,
	0x2f, 0x70, 0x75, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x2f, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72,
	0x2f, 0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x2f, 0x76, 0x31, 0x2f, 0x72, 0x65, 0x67, 0x69,
	0x73, 0x74, 0x65, 0x72, 0x12, 0xfc, 0x01, 0x0a, 0x0c, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x53,
	0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x27, 0x2e, 0x70, 0x75, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x2e,
	0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74,
	0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x28,
	0x2e, 0x70, 0x75, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x2e, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72,
	0x2e, 0x76, 0x31, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x98, 0x01, 0x92, 0x41, 0x61, 0x0a, 0x0f,
	0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x5f, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x12,
	0x26, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x20, 0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x20,
	0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x20, 0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65,
	0x20, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x1a, 0x26, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x20,
	0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x20, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x20,
	0x74, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x20, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x82,
	0xd3, 0xe4, 0x93, 0x02, 0x2e, 0x3a, 0x01, 0x2a, 0x22, 0x29, 0x2f, 0x70, 0x75, 0x64, 0x64, 0x69,
	0x6e, 0x67, 0x2f, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x2f, 0x77, 0x65, 0x62, 0x68, 0x6f,
	0x6f, 0x6b, 0x2f, 0x76, 0x31, 0x2f, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x5f, 0x73, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x12, 0xe9, 0x01, 0x0a, 0x04, 0x43, 0x61, 0x6c, 0x6c, 0x12, 0x34, 0x2e, 0x70,
	0x75, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x2e, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x2e, 0x76,
	0x31, 0x2e, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x54, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72,
	0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x43, 0x61, 0x6c, 0x6c, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x35, 0x2e, 0x70, 0x75, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x2e, 0x74, 0x72, 0x69,
	0x67, 0x67, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x54,
	0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x43, 0x61, 0x6c,
	0x6c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x74, 0x92, 0x41, 0x41, 0x0a, 0x0f,
	0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x5f, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x12,
	0x16, 0x63, 0x61, 0x6c, 0x6c, 0x20, 0x61, 0x20, 0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x20,
	0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x1a, 0x16, 0x63, 0x61, 0x6c, 0x6c, 0x20, 0x61, 0x20,
	0x77, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b, 0x20, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x82,
	0xd3, 0xe4, 0x93, 0x02, 0x2a, 0x3a, 0x01, 0x2a, 0x1a, 0x25, 0x2f, 0x70, 0x75, 0x64, 0x64, 0x69,
	0x6e, 0x67, 0x2f, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x2f, 0x77, 0x65, 0x62, 0x68, 0x6f,
	0x6f, 0x6b, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x61, 0x6c, 0x6c, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x42,
	0xe6, 0x01, 0x92, 0x41, 0x0a, 0x12, 0x05, 0x32, 0x03, 0x31, 0x2e, 0x30, 0x2a, 0x01, 0x02, 0x0a,
	0x16, 0x63, 0x6f, 0x6d, 0x2e, 0x70, 0x75, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x2e, 0x74, 0x72, 0x69,
	0x67, 0x67, 0x65, 0x72, 0x2e, 0x76, 0x31, 0x42, 0x13, 0x57, 0x65, 0x62, 0x68, 0x6f, 0x6f, 0x6b,
	0x54, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x40,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x62, 0x65, 0x69, 0x68, 0x61,
	0x69, 0x30, 0x78, 0x66, 0x66, 0x2f, 0x70, 0x75, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x2f, 0x61, 0x70,
	0x69, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x75, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x2f, 0x74, 0x72,
	0x69, 0x67, 0x67, 0x65, 0x72, 0x2f, 0x76, 0x31, 0x3b, 0x74, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72,
	0xa2, 0x02, 0x03, 0x50, 0x54, 0x58, 0xaa, 0x02, 0x12, 0x50, 0x75, 0x64, 0x64, 0x69, 0x6e, 0x67,
	0x2e, 0x54, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x12, 0x50, 0x75,
	0x64, 0x64, 0x69, 0x6e, 0x67, 0x5c, 0x54, 0x72, 0x69, 0x67, 0x67, 0x65, 0x72, 0x5c, 0x56, 0x31,
	0xe2, 0x02, 0x1e, 0x50, 0x75, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x5c, 0x54, 0x72, 0x69, 0x67, 0x67,
	0x65, 0x72, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0xea, 0x02, 0x14, 0x50, 0x75, 0x64, 0x64, 0x69, 0x6e, 0x67, 0x3a, 0x3a, 0x54, 0x72, 0x69,
	0x67, 0x67, 0x65, 0x72, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pudding_trigger_v1_webhook_trigger_proto_rawDescOnce sync.Once
	file_pudding_trigger_v1_webhook_trigger_proto_rawDescData = file_pudding_trigger_v1_webhook_trigger_proto_rawDesc
)

func file_pudding_trigger_v1_webhook_trigger_proto_rawDescGZIP() []byte {
	file_pudding_trigger_v1_webhook_trigger_proto_rawDescOnce.Do(func() {
		file_pudding_trigger_v1_webhook_trigger_proto_rawDescData = protoimpl.X.CompressGZIP(file_pudding_trigger_v1_webhook_trigger_proto_rawDescData)
	})
	return file_pudding_trigger_v1_webhook_trigger_proto_rawDescData
}

var file_pudding_trigger_v1_webhook_trigger_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_pudding_trigger_v1_webhook_trigger_proto_goTypes = []interface{}{
	(*WebhookFindOneByIDResponse)(nil),           // 0: pudding.trigger.v1.WebhookFindOneByIDResponse
	(*WebhookPageQueryResponse)(nil),             // 1: pudding.trigger.v1.WebhookPageQueryResponse
	(*WebhookTriggerServiceRegisterRequest)(nil), // 2: pudding.trigger.v1.WebhookTriggerServiceRegisterRequest
	(*WebhookRegisterResponse)(nil),              // 3: pudding.trigger.v1.WebhookRegisterResponse
	(*WebhookTriggerServiceCallRequest)(nil),     // 4: pudding.trigger.v1.WebhookTriggerServiceCallRequest
	(*WebhookTriggerServiceCallResponse)(nil),    // 5: pudding.trigger.v1.WebhookTriggerServiceCallResponse
	(*WebhookTriggerTemplate)(nil),               // 6: pudding.trigger.v1.WebhookTriggerTemplate
	(*timestamppb.Timestamp)(nil),                // 7: google.protobuf.Timestamp
	(*FindOneByIDRequest)(nil),                   // 8: pudding.trigger.v1.FindOneByIDRequest
	(*PageQueryTemplateRequest)(nil),             // 9: pudding.trigger.v1.PageQueryTemplateRequest
	(*UpdateStatusRequest)(nil),                  // 10: pudding.trigger.v1.UpdateStatusRequest
	(*UpdateStatusResponse)(nil),                 // 11: pudding.trigger.v1.UpdateStatusResponse
}
var file_pudding_trigger_v1_webhook_trigger_proto_depIdxs = []int32{
	6,  // 0: pudding.trigger.v1.WebhookFindOneByIDResponse.body:type_name -> pudding.trigger.v1.WebhookTriggerTemplate
	6,  // 1: pudding.trigger.v1.WebhookPageQueryResponse.body:type_name -> pudding.trigger.v1.WebhookTriggerTemplate
	7,  // 2: pudding.trigger.v1.WebhookTriggerServiceRegisterRequest.excepted_end_time:type_name -> google.protobuf.Timestamp
	8,  // 3: pudding.trigger.v1.WebhookTriggerService.FindOneByID:input_type -> pudding.trigger.v1.FindOneByIDRequest
	9,  // 4: pudding.trigger.v1.WebhookTriggerService.PageQueryTemplate:input_type -> pudding.trigger.v1.PageQueryTemplateRequest
	2,  // 5: pudding.trigger.v1.WebhookTriggerService.Register:input_type -> pudding.trigger.v1.WebhookTriggerServiceRegisterRequest
	10, // 6: pudding.trigger.v1.WebhookTriggerService.UpdateStatus:input_type -> pudding.trigger.v1.UpdateStatusRequest
	4,  // 7: pudding.trigger.v1.WebhookTriggerService.Call:input_type -> pudding.trigger.v1.WebhookTriggerServiceCallRequest
	0,  // 8: pudding.trigger.v1.WebhookTriggerService.FindOneByID:output_type -> pudding.trigger.v1.WebhookFindOneByIDResponse
	1,  // 9: pudding.trigger.v1.WebhookTriggerService.PageQueryTemplate:output_type -> pudding.trigger.v1.WebhookPageQueryResponse
	3,  // 10: pudding.trigger.v1.WebhookTriggerService.Register:output_type -> pudding.trigger.v1.WebhookRegisterResponse
	11, // 11: pudding.trigger.v1.WebhookTriggerService.UpdateStatus:output_type -> pudding.trigger.v1.UpdateStatusResponse
	5,  // 12: pudding.trigger.v1.WebhookTriggerService.Call:output_type -> pudding.trigger.v1.WebhookTriggerServiceCallResponse
	8,  // [8:13] is the sub-list for method output_type
	3,  // [3:8] is the sub-list for method input_type
	3,  // [3:3] is the sub-list for extension type_name
	3,  // [3:3] is the sub-list for extension extendee
	0,  // [0:3] is the sub-list for field type_name
}

func init() { file_pudding_trigger_v1_webhook_trigger_proto_init() }
func file_pudding_trigger_v1_webhook_trigger_proto_init() {
	if File_pudding_trigger_v1_webhook_trigger_proto != nil {
		return
	}
	file_pudding_trigger_v1_types_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_pudding_trigger_v1_webhook_trigger_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WebhookFindOneByIDResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pudding_trigger_v1_webhook_trigger_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WebhookPageQueryResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pudding_trigger_v1_webhook_trigger_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WebhookTriggerServiceRegisterRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pudding_trigger_v1_webhook_trigger_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WebhookRegisterResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pudding_trigger_v1_webhook_trigger_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WebhookTriggerServiceCallRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pudding_trigger_v1_webhook_trigger_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WebhookTriggerServiceCallResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pudding_trigger_v1_webhook_trigger_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pudding_trigger_v1_webhook_trigger_proto_goTypes,
		DependencyIndexes: file_pudding_trigger_v1_webhook_trigger_proto_depIdxs,
		MessageInfos:      file_pudding_trigger_v1_webhook_trigger_proto_msgTypes,
	}.Build()
	File_pudding_trigger_v1_webhook_trigger_proto = out.File
	file_pudding_trigger_v1_webhook_trigger_proto_rawDesc = nil
	file_pudding_trigger_v1_webhook_trigger_proto_goTypes = nil
	file_pudding_trigger_v1_webhook_trigger_proto_depIdxs = nil
}
