// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/beihai0xff/pudding/internal/scheduler (interfaces: Scheduler)

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	types "github.com/beihai0xff/pudding/types"
	gomock "github.com/golang/mock/gomock"
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
func (m *MockScheduler) NewConsumer(arg0, arg1 string, arg2 int, arg3 types.HandleMessage) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewConsumer", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// NewConsumer indicates an expected call of NewConsumer.
func (mr *MockSchedulerMockRecorder) NewConsumer(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewConsumer", reflect.TypeOf((*MockScheduler)(nil).NewConsumer), arg0, arg1, arg2, arg3)
}

// Produce mocks base method.
func (m *MockScheduler) Produce(arg0 context.Context, arg1 *types.Message) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Produce", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Produce indicates an expected call of Produce.
func (mr *MockSchedulerMockRecorder) Produce(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Produce", reflect.TypeOf((*MockScheduler)(nil).Produce), arg0, arg1)
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
