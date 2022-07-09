// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package repository is a generated GoMock package.
package repository

import (
	context "context"
	reflect "reflect"

	model "github.com/dcaf-protocol/drip/pkg/repository/model"
	gomock "github.com/golang/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// EnableVault mocks base method.
func (m *MockRepository) EnableVault(ctx context.Context, pubkey string) (*model.Vault, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "EnableVault", ctx, pubkey)
	ret0, _ := ret[0].(*model.Vault)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// EnableVault indicates an expected call of EnableVault.
func (mr *MockRepositoryMockRecorder) EnableVault(ctx, pubkey interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "EnableVault", reflect.TypeOf((*MockRepository)(nil).EnableVault), ctx, pubkey)
}

// GetProtoConfigs mocks base method.
func (m *MockRepository) GetProtoConfigs(arg0 context.Context, arg1, arg2 *string) ([]*model.ProtoConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProtoConfigs", arg0, arg1, arg2)
	ret0, _ := ret[0].([]*model.ProtoConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProtoConfigs indicates an expected call of GetProtoConfigs.
func (mr *MockRepositoryMockRecorder) GetProtoConfigs(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProtoConfigs", reflect.TypeOf((*MockRepository)(nil).GetProtoConfigs), arg0, arg1, arg2)
}

// GetTokenPair mocks base method.
func (m *MockRepository) GetTokenPair(arg0 context.Context, arg1, arg2 string) (*model.TokenPair, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTokenPair", arg0, arg1, arg2)
	ret0, _ := ret[0].(*model.TokenPair)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTokenPair indicates an expected call of GetTokenPair.
func (mr *MockRepositoryMockRecorder) GetTokenPair(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTokenPair", reflect.TypeOf((*MockRepository)(nil).GetTokenPair), arg0, arg1, arg2)
}

// GetTokenPairByID mocks base method.
func (m *MockRepository) GetTokenPairByID(arg0 context.Context, arg1 string) (*model.TokenPair, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTokenPairByID", arg0, arg1)
	ret0, _ := ret[0].(*model.TokenPair)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTokenPairByID indicates an expected call of GetTokenPairByID.
func (mr *MockRepositoryMockRecorder) GetTokenPairByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTokenPairByID", reflect.TypeOf((*MockRepository)(nil).GetTokenPairByID), arg0, arg1)
}

// GetTokenPairs mocks base method.
func (m *MockRepository) GetTokenPairs(arg0 context.Context, arg1, arg2 *string) ([]*model.TokenPair, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTokenPairs", arg0, arg1, arg2)
	ret0, _ := ret[0].([]*model.TokenPair)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTokenPairs indicates an expected call of GetTokenPairs.
func (mr *MockRepositoryMockRecorder) GetTokenPairs(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTokenPairs", reflect.TypeOf((*MockRepository)(nil).GetTokenPairs), arg0, arg1, arg2)
}

// GetTokenSwapForTokenAccount mocks base method.
func (m *MockRepository) GetTokenSwapForTokenAccount(arg0 context.Context, arg1 string) (*model.TokenSwap, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTokenSwapForTokenAccount", arg0, arg1)
	ret0, _ := ret[0].(*model.TokenSwap)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTokenSwapForTokenAccount indicates an expected call of GetTokenSwapForTokenAccount.
func (mr *MockRepositoryMockRecorder) GetTokenSwapForTokenAccount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTokenSwapForTokenAccount", reflect.TypeOf((*MockRepository)(nil).GetTokenSwapForTokenAccount), arg0, arg1)
}

// GetTokenSwaps mocks base method.
func (m *MockRepository) GetTokenSwaps(arg0 context.Context, arg1 []string) ([]*model.TokenSwap, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTokenSwaps", arg0, arg1)
	ret0, _ := ret[0].([]*model.TokenSwap)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTokenSwaps indicates an expected call of GetTokenSwaps.
func (mr *MockRepositoryMockRecorder) GetTokenSwaps(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTokenSwaps", reflect.TypeOf((*MockRepository)(nil).GetTokenSwaps), arg0, arg1)
}

// GetTokenSwapsSortedByLiquidity mocks base method.
func (m *MockRepository) GetTokenSwapsSortedByLiquidity(ctx context.Context, tokenPairIDs []string) ([]TokenSwapWithLiquidityRatio, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTokenSwapsSortedByLiquidity", ctx, tokenPairIDs)
	ret0, _ := ret[0].([]TokenSwapWithLiquidityRatio)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTokenSwapsSortedByLiquidity indicates an expected call of GetTokenSwapsSortedByLiquidity.
func (mr *MockRepositoryMockRecorder) GetTokenSwapsSortedByLiquidity(ctx, tokenPairIDs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTokenSwapsSortedByLiquidity", reflect.TypeOf((*MockRepository)(nil).GetTokenSwapsSortedByLiquidity), ctx, tokenPairIDs)
}

// GetTokensWithSupportedTokenPair mocks base method.
func (m *MockRepository) GetTokensWithSupportedTokenPair(arg0 context.Context, arg1 *string, arg2 bool) ([]*model.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTokensWithSupportedTokenPair", arg0, arg1, arg2)
	ret0, _ := ret[0].([]*model.Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTokensWithSupportedTokenPair indicates an expected call of GetTokensWithSupportedTokenPair.
func (mr *MockRepositoryMockRecorder) GetTokensWithSupportedTokenPair(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTokensWithSupportedTokenPair", reflect.TypeOf((*MockRepository)(nil).GetTokensWithSupportedTokenPair), arg0, arg1, arg2)
}

// GetVaultByAddress mocks base method.
func (m *MockRepository) GetVaultByAddress(arg0 context.Context, arg1 string) (*model.Vault, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVaultByAddress", arg0, arg1)
	ret0, _ := ret[0].(*model.Vault)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVaultByAddress indicates an expected call of GetVaultByAddress.
func (mr *MockRepositoryMockRecorder) GetVaultByAddress(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVaultByAddress", reflect.TypeOf((*MockRepository)(nil).GetVaultByAddress), arg0, arg1)
}

// GetVaultPeriods mocks base method.
func (m *MockRepository) GetVaultPeriods(arg0 context.Context, arg1 string, arg2, arg3 int, arg4 *string) ([]*model.VaultPeriod, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVaultPeriods", arg0, arg1, arg2, arg3, arg4)
	ret0, _ := ret[0].([]*model.VaultPeriod)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVaultPeriods indicates an expected call of GetVaultPeriods.
func (mr *MockRepositoryMockRecorder) GetVaultPeriods(arg0, arg1, arg2, arg3, arg4 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVaultPeriods", reflect.TypeOf((*MockRepository)(nil).GetVaultPeriods), arg0, arg1, arg2, arg3, arg4)
}

// GetVaultsWithFilter mocks base method.
func (m *MockRepository) GetVaultsWithFilter(arg0 context.Context, arg1, arg2, arg3 *string) ([]*model.Vault, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVaultsWithFilter", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].([]*model.Vault)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVaultsWithFilter indicates an expected call of GetVaultsWithFilter.
func (mr *MockRepositoryMockRecorder) GetVaultsWithFilter(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVaultsWithFilter", reflect.TypeOf((*MockRepository)(nil).GetVaultsWithFilter), arg0, arg1, arg2, arg3)
}

// InsertTokenPairs mocks base method.
func (m *MockRepository) InsertTokenPairs(arg0 context.Context, arg1 ...*model.TokenPair) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "InsertTokenPairs", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertTokenPairs indicates an expected call of InsertTokenPairs.
func (mr *MockRepositoryMockRecorder) InsertTokenPairs(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertTokenPairs", reflect.TypeOf((*MockRepository)(nil).InsertTokenPairs), varargs...)
}

// InternalGetVaultByAddress mocks base method.
func (m *MockRepository) InternalGetVaultByAddress(ctx context.Context, pubkey string) (*model.Vault, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InternalGetVaultByAddress", ctx, pubkey)
	ret0, _ := ret[0].(*model.Vault)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InternalGetVaultByAddress indicates an expected call of InternalGetVaultByAddress.
func (mr *MockRepositoryMockRecorder) InternalGetVaultByAddress(ctx, pubkey interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InternalGetVaultByAddress", reflect.TypeOf((*MockRepository)(nil).InternalGetVaultByAddress), ctx, pubkey)
}

// UpsertPositions mocks base method.
func (m *MockRepository) UpsertPositions(arg0 context.Context, arg1 ...*model.Position) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpsertPositions", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertPositions indicates an expected call of UpsertPositions.
func (mr *MockRepositoryMockRecorder) UpsertPositions(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertPositions", reflect.TypeOf((*MockRepository)(nil).UpsertPositions), varargs...)
}

// UpsertProtoConfigs mocks base method.
func (m *MockRepository) UpsertProtoConfigs(arg0 context.Context, arg1 ...*model.ProtoConfig) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpsertProtoConfigs", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertProtoConfigs indicates an expected call of UpsertProtoConfigs.
func (mr *MockRepositoryMockRecorder) UpsertProtoConfigs(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertProtoConfigs", reflect.TypeOf((*MockRepository)(nil).UpsertProtoConfigs), varargs...)
}

// UpsertTokenAccountBalances mocks base method.
func (m *MockRepository) UpsertTokenAccountBalances(arg0 context.Context, arg1 ...*model.TokenAccountBalance) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpsertTokenAccountBalances", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertTokenAccountBalances indicates an expected call of UpsertTokenAccountBalances.
func (mr *MockRepositoryMockRecorder) UpsertTokenAccountBalances(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertTokenAccountBalances", reflect.TypeOf((*MockRepository)(nil).UpsertTokenAccountBalances), varargs...)
}

// UpsertTokenSwaps mocks base method.
func (m *MockRepository) UpsertTokenSwaps(arg0 context.Context, arg1 ...*model.TokenSwap) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpsertTokenSwaps", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertTokenSwaps indicates an expected call of UpsertTokenSwaps.
func (mr *MockRepositoryMockRecorder) UpsertTokenSwaps(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertTokenSwaps", reflect.TypeOf((*MockRepository)(nil).UpsertTokenSwaps), varargs...)
}

// UpsertTokens mocks base method.
func (m *MockRepository) UpsertTokens(arg0 context.Context, arg1 ...*model.Token) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpsertTokens", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertTokens indicates an expected call of UpsertTokens.
func (mr *MockRepositoryMockRecorder) UpsertTokens(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertTokens", reflect.TypeOf((*MockRepository)(nil).UpsertTokens), varargs...)
}

// UpsertVaultPeriods mocks base method.
func (m *MockRepository) UpsertVaultPeriods(arg0 context.Context, arg1 ...*model.VaultPeriod) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpsertVaultPeriods", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertVaultPeriods indicates an expected call of UpsertVaultPeriods.
func (mr *MockRepositoryMockRecorder) UpsertVaultPeriods(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertVaultPeriods", reflect.TypeOf((*MockRepository)(nil).UpsertVaultPeriods), varargs...)
}

// UpsertVaults mocks base method.
func (m *MockRepository) UpsertVaults(arg0 context.Context, arg1 ...*model.Vault) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpsertVaults", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertVaults indicates an expected call of UpsertVaults.
func (mr *MockRepositoryMockRecorder) UpsertVaults(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertVaults", reflect.TypeOf((*MockRepository)(nil).UpsertVaults), varargs...)
}
