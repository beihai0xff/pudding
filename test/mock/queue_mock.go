// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/beihai0xff/pudding/app/scheduler/broker (interfaces: DelayQueue,RealTimeQueue)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	types "github.com/beihai0xff/pudding/types"
	gomock "github.com/golang/mock/gomock"
)

// MockDelayQueue is a mock of DelayQueue interface.
type MockDelayQueue struct {
	ctrl     *gomock.Controller
	recorder *MockDelayQueueMockRecorder
}

// MockDelayQueueMockRecorder is the mock recorder for MockDelayQueue.
type MockDelayQueueMockRecorder struct {
	mock *MockDelayQueue
}

// NewMockDelayQueue creates a new mock instance.
func NewMockDelayQueue(ctrl *gomock.Controller) *MockDelayQueue {
	mock := &MockDelayQueue{ctrl: ctrl}
	mock.recorder = &MockDelayQueueMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDelayQueue) EXPECT() *MockDelayQueueMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockDelayQueue) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockDelayQueueMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockDelayQueue)(nil).Close))
}

// Consume mocks base method.
func (m *MockDelayQueue) Consume(arg0 context.Context, arg1 string, arg2, arg3 int64, arg4 types.HandleMessage) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Consume", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].(error)
	return ret0
}

// Consume indicates an expected call of Consume.
func (mr *MockDelayQueueMockRecorder) Consume(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Consume", reflect.TypeOf((*MockDelayQueue)(nil).Consume), arg0, arg1, arg2, arg3, arg4)
}

// Produce mocks base method.
func (m *MockDelayQueue) Produce(arg0 context.Context, arg1 string, arg2 *types.Message) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Produce", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Produce indicates an expected call of Produce.
func (mr *MockDelayQueueMockRecorder) Produce(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Produce", reflect.TypeOf((*MockDelayQueue)(nil).Produce), arg0, arg1, arg2)
}

// MockRealTimeQueue is a mock of RealTimeQueue interface.
type MockRealTimeQueue struct {
	ctrl     *gomock.Controller
	recorder *MockRealTimeQueueMockRecorder
}

// MockRealTimeQueueMockRecorder is the mock recorder for MockRealTimeQueue.
type MockRealTimeQueueMockRecorder struct {
	mock *MockRealTimeQueue
}

// NewMockRealTimeQueue creates a new mock instance.
func NewMockRealTimeQueue(ctrl *gomock.Controller) *MockRealTimeQueue {
	mock := &MockRealTimeQueue{ctrl: ctrl}
	mock.recorder = &MockRealTimeQueueMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRealTimeQueue) EXPECT() *MockRealTimeQueueMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockRealTimeQueue) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockRealTimeQueueMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockRealTimeQueue)(nil).Close))
}

// NewConsumer mocks base method.
func (m *MockRealTimeQueue) NewConsumer(arg0, arg1 string, arg2 int, arg3 types.HandleMessage) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewConsumer", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// NewConsumer indicates an expected call of NewConsumer.
func (mr *MockRealTimeQueueMockRecorder) NewConsumer(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewConsumer", reflect.TypeOf((*MockRealTimeQueue)(nil).NewConsumer), arg0, arg1, arg2, arg3)
}

// Produce mocks base method.
func (m *MockRealTimeQueue) Produce(arg0 context.Context, arg1 *types.Message) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Produce", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Produce indicates an expected call of Produce.
func (mr *MockRealTimeQueueMockRecorder) Produce(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Produce", reflect.TypeOf((*MockRealTimeQueue)(nil).Produce), arg0, arg1)
}
