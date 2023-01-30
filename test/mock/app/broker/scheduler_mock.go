// Code generated by MockGen. DO NOT EDIT.
// Source: app/broker/scheduler.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"

	types "github.com/beihai0xff/pudding/api/gen/pudding/types/v1"
	types0 "github.com/beihai0xff/pudding/app/broker/pkg/types"
)

// MockScheduler is a mock of Scheduler interface.
type MockScheduler struct {
	ctrl     *gomock.Controller
	recorder *MockSchedulerMockRecorder
}

// MockSchedulerMockRecorder is the mock recorder for MockScheduler.
type MockSchedulerMockRecorder struct {
	mock *MockScheduler
}

// NewMockScheduler creates a new mock instance.
func NewMockScheduler(ctrl *gomock.Controller) *MockScheduler {
	mock := &MockScheduler{ctrl: ctrl}
	mock.recorder = &MockSchedulerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockScheduler) EXPECT() *MockSchedulerMockRecorder {
	return m.recorder
}

// NewConsumer mocks base method.
func (m *MockScheduler) NewConsumer(topic, group string, batchSize int, fn types0.HandleMessage) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewConsumer", topic, group, batchSize, fn)
	ret0, _ := ret[0].(error)
	return ret0
}

// NewConsumer indicates an expected call of NewConsumer.
func (mr *MockSchedulerMockRecorder) NewConsumer(topic, group, batchSize, fn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewConsumer", reflect.TypeOf((*MockScheduler)(nil).NewConsumer), topic, group, batchSize, fn)
}

// Produce mocks base method.
func (m *MockScheduler) Produce(ctx context.Context, msg *types.Message) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Produce", ctx, msg)
	ret0, _ := ret[0].(error)
	return ret0
}

// Produce indicates an expected call of Produce.
func (mr *MockSchedulerMockRecorder) Produce(ctx, msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Produce", reflect.TypeOf((*MockScheduler)(nil).Produce), ctx, msg)
}

// Run mocks base method.
func (m *MockScheduler) Run() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Run")
}

// Run indicates an expected call of Run.
func (mr *MockSchedulerMockRecorder) Run() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockScheduler)(nil).Run))
}
