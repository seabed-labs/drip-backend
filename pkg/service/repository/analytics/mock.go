// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package repository is a generated GoMock package.
package repository

import (
	context "context"
	reflect "reflect"

	model "github.com/dcaf-labs/drip/pkg/service/repository/model"
	gomock "github.com/golang/mock/gomock"
)

// MockAnalyticsRepository is a mock of AnalyticsRepository interface.
type MockAnalyticsRepository struct {
	ctrl     *gomock.Controller
	recorder *MockAnalyticsRepositoryMockRecorder
}

// MockAnalyticsRepositoryMockRecorder is the mock recorder for MockAnalyticsRepository.
type MockAnalyticsRepositoryMockRecorder struct {
	mock *MockAnalyticsRepository
}

// NewMockAnalyticsRepository creates a new mock instance.
func NewMockAnalyticsRepository(ctrl *gomock.Controller) *MockAnalyticsRepository {
	mock := &MockAnalyticsRepository{ctrl: ctrl}
	mock.recorder = &MockAnalyticsRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAnalyticsRepository) EXPECT() *MockAnalyticsRepositoryMockRecorder {
	return m.recorder
}

// GetCurrentTVL mocks base method.
func (m *MockAnalyticsRepository) GetCurrentTVL(ctx context.Context) (*model.CurrentTVL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCurrentTVL", ctx)
	ret0, _ := ret[0].(*model.CurrentTVL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCurrentTVL indicates an expected call of GetCurrentTVL.
func (mr *MockAnalyticsRepositoryMockRecorder) GetCurrentTVL(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCurrentTVL", reflect.TypeOf((*MockAnalyticsRepository)(nil).GetCurrentTVL), ctx)
}

// GetDepositMetricBySignature mocks base method.
func (m *MockAnalyticsRepository) GetDepositMetricBySignature(ctx context.Context, signature string) (*model.DepositMetric, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDepositMetricBySignature", ctx, signature)
	ret0, _ := ret[0].(*model.DepositMetric)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDepositMetricBySignature indicates an expected call of GetDepositMetricBySignature.
func (mr *MockAnalyticsRepositoryMockRecorder) GetDepositMetricBySignature(ctx, signature interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDepositMetricBySignature", reflect.TypeOf((*MockAnalyticsRepository)(nil).GetDepositMetricBySignature), ctx, signature)
}

// GetDripMetricBySignature mocks base method.
func (m *MockAnalyticsRepository) GetDripMetricBySignature(ctx context.Context, signature string) (*model.DripMetric, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDripMetricBySignature", ctx, signature)
	ret0, _ := ret[0].(*model.DripMetric)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDripMetricBySignature indicates an expected call of GetDripMetricBySignature.
func (mr *MockAnalyticsRepositoryMockRecorder) GetDripMetricBySignature(ctx, signature interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDripMetricBySignature", reflect.TypeOf((*MockAnalyticsRepository)(nil).GetDripMetricBySignature), ctx, signature)
}

// GetLifeTimeDepositNormalizedToCurrentPrice mocks base method.
func (m *MockAnalyticsRepository) GetLifeTimeDepositNormalizedToCurrentPrice(ctx context.Context) (*model.LifeTimeDeposit, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLifeTimeDepositNormalizedToCurrentPrice", ctx)
	ret0, _ := ret[0].(*model.LifeTimeDeposit)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLifeTimeDepositNormalizedToCurrentPrice indicates an expected call of GetLifeTimeDepositNormalizedToCurrentPrice.
func (mr *MockAnalyticsRepositoryMockRecorder) GetLifeTimeDepositNormalizedToCurrentPrice(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLifeTimeDepositNormalizedToCurrentPrice", reflect.TypeOf((*MockAnalyticsRepository)(nil).GetLifeTimeDepositNormalizedToCurrentPrice), ctx)
}

// GetLifeTimeVolumeNormalizedToCurrentPrice mocks base method.
func (m *MockAnalyticsRepository) GetLifeTimeVolumeNormalizedToCurrentPrice(ctx context.Context) (*model.LifeTimeVolume, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLifeTimeVolumeNormalizedToCurrentPrice", ctx)
	ret0, _ := ret[0].(*model.LifeTimeVolume)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLifeTimeVolumeNormalizedToCurrentPrice indicates an expected call of GetLifeTimeVolumeNormalizedToCurrentPrice.
func (mr *MockAnalyticsRepositoryMockRecorder) GetLifeTimeVolumeNormalizedToCurrentPrice(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLifeTimeVolumeNormalizedToCurrentPrice", reflect.TypeOf((*MockAnalyticsRepository)(nil).GetLifeTimeVolumeNormalizedToCurrentPrice), ctx)
}

// GetLifeTimeWithdrawalNormalizedToCurrentPrice mocks base method.
func (m *MockAnalyticsRepository) GetLifeTimeWithdrawalNormalizedToCurrentPrice(ctx context.Context) (*model.LifeTimeWithdrawal, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLifeTimeWithdrawalNormalizedToCurrentPrice", ctx)
	ret0, _ := ret[0].(*model.LifeTimeWithdrawal)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLifeTimeWithdrawalNormalizedToCurrentPrice indicates an expected call of GetLifeTimeWithdrawalNormalizedToCurrentPrice.
func (mr *MockAnalyticsRepositoryMockRecorder) GetLifeTimeWithdrawalNormalizedToCurrentPrice(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLifeTimeWithdrawalNormalizedToCurrentPrice", reflect.TypeOf((*MockAnalyticsRepository)(nil).GetLifeTimeWithdrawalNormalizedToCurrentPrice), ctx)
}

// GetUniqueDepositorCount mocks base method.
func (m *MockAnalyticsRepository) GetUniqueDepositorCount(ctx context.Context) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUniqueDepositorCount", ctx)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUniqueDepositorCount indicates an expected call of GetUniqueDepositorCount.
func (mr *MockAnalyticsRepositoryMockRecorder) GetUniqueDepositorCount(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUniqueDepositorCount", reflect.TypeOf((*MockAnalyticsRepository)(nil).GetUniqueDepositorCount), ctx)
}

// GetWithdrawalMetricBySignature mocks base method.
func (m *MockAnalyticsRepository) GetWithdrawalMetricBySignature(ctx context.Context, signature string) (*model.WithdrawalMetric, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWithdrawalMetricBySignature", ctx, signature)
	ret0, _ := ret[0].(*model.WithdrawalMetric)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWithdrawalMetricBySignature indicates an expected call of GetWithdrawalMetricBySignature.
func (mr *MockAnalyticsRepositoryMockRecorder) GetWithdrawalMetricBySignature(ctx, signature interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWithdrawalMetricBySignature", reflect.TypeOf((*MockAnalyticsRepository)(nil).GetWithdrawalMetricBySignature), ctx, signature)
}
