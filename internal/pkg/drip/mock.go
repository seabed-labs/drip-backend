// Code generated by MockGen. DO NOT EDIT.
// Source: drip.go

// Package mock_drip is a generated GoMock package.
package drip

import (
	context "context"
	reflect "reflect"

	model "github.com/dcaf-protocol/drip/internal/pkg/repository/model"
	gomock "github.com/golang/mock/gomock"
)

// MockDrip is a mock of Drip interface.
type MockDrip struct {
	ctrl     *gomock.Controller
	recorder *MockDripMockRecorder
}

// MockDripMockRecorder is the mock recorder for MockDrip.
type MockDripMockRecorder struct {
	mock *MockDrip
}

// NewMockDrip creates a new mock instance.
func NewMockDrip(ctrl *gomock.Controller) *MockDrip {
	mock := &MockDrip{ctrl: ctrl}
	mock.recorder = &MockDripMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDrip) EXPECT() *MockDripMockRecorder {
	return m.recorder
}

// GetProtoConfigs mocks base method.
func (m *MockDrip) GetProtoConfigs(arg0 context.Context, arg1, arg2 *string) ([]*model.ProtoConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProtoConfigs", arg0, arg1, arg2)
	ret0, _ := ret[0].([]*model.ProtoConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProtoConfigs indicates an expected call of GetProtoConfigs.
func (mr *MockDripMockRecorder) GetProtoConfigs(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProtoConfigs", reflect.TypeOf((*MockDrip)(nil).GetProtoConfigs), arg0, arg1, arg2)
}

// GetTokenPair mocks base method.
func (m *MockDrip) GetTokenPair(arg0 context.Context, arg1 string) (*model.TokenPair, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTokenPair", arg0, arg1)
	ret0, _ := ret[0].(*model.TokenPair)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTokenPair indicates an expected call of GetTokenPair.
func (mr *MockDripMockRecorder) GetTokenPair(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTokenPair", reflect.TypeOf((*MockDrip)(nil).GetTokenPair), arg0, arg1)
}

// GetTokens mocks base method.
func (m *MockDrip) GetTokens(arg0 context.Context, arg1, arg2 *string) ([]*model.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTokens", arg0, arg1, arg2)
	ret0, _ := ret[0].([]*model.Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTokens indicates an expected call of GetTokens.
func (mr *MockDripMockRecorder) GetTokens(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTokens", reflect.TypeOf((*MockDrip)(nil).GetTokens), arg0, arg1, arg2)
}

// GetVaultPeriods mocks base method.
func (m *MockDrip) GetVaultPeriods(arg0 context.Context, arg1 string, arg2, arg3 int, arg4 *string) ([]*model.VaultPeriod, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVaultPeriods", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].([]*model.VaultPeriod)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVaultPeriods indicates an expected call of GetVaultPeriods.
func (mr *MockDripMockRecorder) GetVaultPeriods(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVaultPeriods", reflect.TypeOf((*MockDrip)(nil).GetVaultPeriods), arg0, arg1, arg2, arg3, arg4)
}

// GetVaults mocks base method.
func (m *MockDrip) GetVaults(arg0 context.Context, arg1, arg2, arg3 *string) ([]*model.Vault, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVaults", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].([]*model.Vault)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVaults indicates an expected call of GetVaults.
func (mr *MockDripMockRecorder) GetVaults(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVaults", reflect.TypeOf((*MockDrip)(nil).GetVaults), arg0, arg1, arg2, arg3)
}
