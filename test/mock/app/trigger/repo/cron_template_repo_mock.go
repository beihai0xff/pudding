// Code generated by MockGen. DO NOT EDIT.
// Source: app/trigger/repo/cron_template_repo.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"
	time "time"

	trigger "github.com/beihai0xff/pudding/api/gen/pudding/trigger/v1"
	constants "github.com/beihai0xff/pudding/app/trigger/pkg/constants"
	repo "github.com/beihai0xff/pudding/app/trigger/repo"
	po "github.com/beihai0xff/pudding/app/trigger/repo/po"
	gomock "github.com/golang/mock/gomock"
)

// MockCronTemplateDAO is a mock of CronTemplateDAO interface.
type MockCronTemplateDAO struct {
	ctrl     *gomock.Controller
	recorder *MockCronTemplateDAOMockRecorder
}

// MockCronTemplateDAOMockRecorder is the mock recorder for MockCronTemplateDAO.
type MockCronTemplateDAOMockRecorder struct {
	mock *MockCronTemplateDAO
}

// NewMockCronTemplateDAO creates a new mock instance.
func NewMockCronTemplateDAO(ctrl *gomock.Controller) *MockCronTemplateDAO {
	mock := &MockCronTemplateDAO{ctrl: ctrl}
	mock.recorder = &MockCronTemplateDAOMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCronTemplateDAO) EXPECT() *MockCronTemplateDAOMockRecorder {
	return m.recorder
}

// BatchHandleRecords mocks base method.
func (m *MockCronTemplateDAO) BatchHandleRecords(ctx context.Context, t time.Time, batchSize int, f repo.CronTempHandler) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BatchHandleRecords", ctx, t, batchSize, f)
	ret0, _ := ret[0].(error)
	return ret0
}

// BatchHandleRecords indicates an expected call of BatchHandleRecords.
func (mr *MockCronTemplateDAOMockRecorder) BatchHandleRecords(ctx, t, batchSize, f interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BatchHandleRecords", reflect.TypeOf((*MockCronTemplateDAO)(nil).BatchHandleRecords), ctx, t, batchSize, f)
}

// FindByID mocks base method.
func (m *MockCronTemplateDAO) FindByID(ctx context.Context, id uint) (*po.CronTriggerTemplate, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", ctx, id)
	ret0, _ := ret[0].(*po.CronTriggerTemplate)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockCronTemplateDAOMockRecorder) FindByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockCronTemplateDAO)(nil).FindByID), ctx, id)
}

// Insert mocks base method.
func (m *MockCronTemplateDAO) Insert(ctx context.Context, e *po.CronTriggerTemplate) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Insert", ctx, e)
	ret0, _ := ret[0].(error)
	return ret0
}

// Insert indicates an expected call of Insert.
func (mr *MockCronTemplateDAOMockRecorder) Insert(ctx, e interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Insert", reflect.TypeOf((*MockCronTemplateDAO)(nil).Insert), ctx, e)
}

// PageQuery mocks base method.
func (m *MockCronTemplateDAO) PageQuery(ctx context.Context, p *constants.PageQuery, status trigger.TriggerStatus) ([]*po.CronTriggerTemplate, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PageQuery", ctx, p, status)
	ret0, _ := ret[0].([]*po.CronTriggerTemplate)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// PageQuery indicates an expected call of PageQuery.
func (mr *MockCronTemplateDAOMockRecorder) PageQuery(ctx, p, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PageQuery", reflect.TypeOf((*MockCronTemplateDAO)(nil).PageQuery), ctx, p, status)
}

// UpdateStatus mocks base method.
func (m *MockCronTemplateDAO) UpdateStatus(ctx context.Context, id uint, status trigger.TriggerStatus) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateStatus", ctx, id, status)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateStatus indicates an expected call of UpdateStatus.
func (mr *MockCronTemplateDAOMockRecorder) UpdateStatus(ctx, id, status interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateStatus", reflect.TypeOf((*MockCronTemplateDAO)(nil).UpdateStatus), ctx, id, status)
}
