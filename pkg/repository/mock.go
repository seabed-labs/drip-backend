// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package repository is a generated GoMock package.
package repository

import (
	context "context"
	reflect "reflect"

	model "github.com/dcaf-labs/drip/pkg/repository/model"
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

// AdminGetVaultByAddress mocks base method.
func (m *MockRepository) AdminGetVaultByAddress(ctx context.Context, address string) (*model.Vault, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AdminGetVaultByAddress", ctx, address)
	ret0, _ := ret[0].(*model.Vault)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AdminGetVaultByAddress indicates an expected call of AdminGetVaultByAddress.
func (mr *MockRepositoryMockRecorder) AdminGetVaultByAddress(ctx, address interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AdminGetVaultByAddress", reflect.TypeOf((*MockRepository)(nil).AdminGetVaultByAddress), ctx, address)
}

// AdminGetVaults mocks base method.
func (m *MockRepository) AdminGetVaults(ctx context.Context, vaultFilterParams VaultFilterLikeParams, paginationParams PaginationParams) ([]*model.Vault, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AdminGetVaults", ctx, vaultFilterParams, paginationParams)
	ret0, _ := ret[0].([]*model.Vault)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AdminGetVaults indicates an expected call of AdminGetVaults.
func (mr *MockRepositoryMockRecorder) AdminGetVaults(ctx, vaultFilterParams, paginationParams interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AdminGetVaults", reflect.TypeOf((*MockRepository)(nil).AdminGetVaults), ctx, vaultFilterParams, paginationParams)
}

// AdminSetVaultEnabled mocks base method.
func (m *MockRepository) AdminSetVaultEnabled(ctx context.Context, pubkey string, enabled bool) (*model.Vault, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AdminSetVaultEnabled", ctx, pubkey, enabled)
	ret0, _ := ret[0].(*model.Vault)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AdminSetVaultEnabled indicates an expected call of AdminSetVaultEnabled.
func (mr *MockRepositoryMockRecorder) AdminSetVaultEnabled(ctx, pubkey, enabled interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AdminSetVaultEnabled", reflect.TypeOf((*MockRepository)(nil).AdminSetVaultEnabled), ctx, pubkey, enabled)
}

// GetActiveWallets mocks base method.
func (m *MockRepository) GetActiveWallets(ctx context.Context, params GetActiveWalletParams) ([]ActiveWallet, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetActiveWallets", ctx, params)
	ret0, _ := ret[0].([]ActiveWallet)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetActiveWallets indicates an expected call of GetActiveWallets.
func (mr *MockRepositoryMockRecorder) GetActiveWallets(ctx, params interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetActiveWallets", reflect.TypeOf((*MockRepository)(nil).GetActiveWallets), ctx, params)
}

// GetAdminPositions mocks base method.
func (m *MockRepository) GetAdminPositions(ctx context.Context, isVaultEnabled *bool, positionFilterParams PositionFilterParams, paginationParams PaginationParams) ([]*model.Position, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAdminPositions", ctx, isVaultEnabled, positionFilterParams, paginationParams)
	ret0, _ := ret[0].([]*model.Position)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAdminPositions indicates an expected call of GetAdminPositions.
func (mr *MockRepositoryMockRecorder) GetAdminPositions(ctx, isVaultEnabled, positionFilterParams, paginationParams interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAdminPositions", reflect.TypeOf((*MockRepository)(nil).GetAdminPositions), ctx, isVaultEnabled, positionFilterParams, paginationParams)
}

// GetOrcaWhirlpoolByAddress mocks base method.
func (m *MockRepository) GetOrcaWhirlpoolByAddress(ctx context.Context, address string) (*model.OrcaWhirlpool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrcaWhirlpoolByAddress", ctx, address)
	ret0, _ := ret[0].(*model.OrcaWhirlpool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrcaWhirlpoolByAddress indicates an expected call of GetOrcaWhirlpoolByAddress.
func (mr *MockRepositoryMockRecorder) GetOrcaWhirlpoolByAddress(ctx, address interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrcaWhirlpoolByAddress", reflect.TypeOf((*MockRepository)(nil).GetOrcaWhirlpoolByAddress), ctx, address)
}

// GetOrcaWhirlpoolsByTokenPairIDs mocks base method.
func (m *MockRepository) GetOrcaWhirlpoolsByTokenPairIDs(ctx context.Context, tokenPairIDs []string) ([]*model.OrcaWhirlpool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrcaWhirlpoolsByTokenPairIDs", ctx, tokenPairIDs)
	ret0, _ := ret[0].([]*model.OrcaWhirlpool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrcaWhirlpoolsByTokenPairIDs indicates an expected call of GetOrcaWhirlpoolsByTokenPairIDs.
func (mr *MockRepositoryMockRecorder) GetOrcaWhirlpoolsByTokenPairIDs(ctx, tokenPairIDs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrcaWhirlpoolsByTokenPairIDs", reflect.TypeOf((*MockRepository)(nil).GetOrcaWhirlpoolsByTokenPairIDs), ctx, tokenPairIDs)
}

// GetPositionByNFTMint mocks base method.
func (m *MockRepository) GetPositionByNFTMint(ctx context.Context, nftMint string) (*model.Position, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPositionByNFTMint", ctx, nftMint)
	ret0, _ := ret[0].(*model.Position)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPositionByNFTMint indicates an expected call of GetPositionByNFTMint.
func (mr *MockRepositoryMockRecorder) GetPositionByNFTMint(ctx, nftMint interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPositionByNFTMint", reflect.TypeOf((*MockRepository)(nil).GetPositionByNFTMint), ctx, nftMint)
}

// GetProtoConfigs mocks base method.
func (m *MockRepository) GetProtoConfigs(arg0 context.Context) ([]*model.ProtoConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProtoConfigs", arg0)
	ret0, _ := ret[0].([]*model.ProtoConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProtoConfigs indicates an expected call of GetProtoConfigs.
func (mr *MockRepositoryMockRecorder) GetProtoConfigs(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProtoConfigs", reflect.TypeOf((*MockRepository)(nil).GetProtoConfigs), arg0)
}

// GetProtoConfigsByAddresses mocks base method.
func (m *MockRepository) GetProtoConfigsByAddresses(ctx context.Context, pubkeys []string) ([]*model.ProtoConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProtoConfigsByAddresses", ctx, pubkeys)
	ret0, _ := ret[0].([]*model.ProtoConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProtoConfigsByAddresses indicates an expected call of GetProtoConfigsByAddresses.
func (mr *MockRepositoryMockRecorder) GetProtoConfigsByAddresses(ctx, pubkeys interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProtoConfigsByAddresses", reflect.TypeOf((*MockRepository)(nil).GetProtoConfigsByAddresses), ctx, pubkeys)
}

// GetTokenAccountBalancesByIDS mocks base method.
func (m *MockRepository) GetTokenAccountBalancesByIDS(arg0 context.Context, arg1 []string) ([]*model.TokenAccountBalance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTokenAccountBalancesByIDS", arg0, arg1)
	ret0, _ := ret[0].([]*model.TokenAccountBalance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTokenAccountBalancesByIDS indicates an expected call of GetTokenAccountBalancesByIDS.
func (mr *MockRepositoryMockRecorder) GetTokenAccountBalancesByIDS(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTokenAccountBalancesByIDS", reflect.TypeOf((*MockRepository)(nil).GetTokenAccountBalancesByIDS), arg0, arg1)
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
func (m *MockRepository) GetTokenPairs(arg0 context.Context, arg1 TokenPairFilterParams) ([]*model.TokenPair, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTokenPairs", arg0, arg1)
	ret0, _ := ret[0].([]*model.TokenPair)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTokenPairs indicates an expected call of GetTokenPairs.
func (mr *MockRepositoryMockRecorder) GetTokenPairs(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTokenPairs", reflect.TypeOf((*MockRepository)(nil).GetTokenPairs), arg0, arg1)
}

// GetTokenPairsByIDS mocks base method.
func (m *MockRepository) GetTokenPairsByIDS(arg0 context.Context, arg1 []string) ([]*model.TokenPair, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTokenPairsByIDS", arg0, arg1)
	ret0, _ := ret[0].([]*model.TokenPair)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTokenPairsByIDS indicates an expected call of GetTokenPairsByIDS.
func (mr *MockRepositoryMockRecorder) GetTokenPairsByIDS(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTokenPairsByIDS", reflect.TypeOf((*MockRepository)(nil).GetTokenPairsByIDS), arg0, arg1)
}

// GetTokenSwapByAddress mocks base method.
func (m *MockRepository) GetTokenSwapByAddress(arg0 context.Context, arg1 string) (*model.TokenSwap, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTokenSwapByAddress", arg0, arg1)
	ret0, _ := ret[0].(*model.TokenSwap)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTokenSwapByAddress indicates an expected call of GetTokenSwapByAddress.
func (mr *MockRepositoryMockRecorder) GetTokenSwapByAddress(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTokenSwapByAddress", reflect.TypeOf((*MockRepository)(nil).GetTokenSwapByAddress), arg0, arg1)
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

// GetTokenSwapsWithBalance mocks base method.
func (m *MockRepository) GetTokenSwapsWithBalance(ctx context.Context, tokenPairIDs []string) ([]TokenSwapWithBalance, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTokenSwapsWithBalance", ctx, tokenPairIDs)
	ret0, _ := ret[0].([]TokenSwapWithBalance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTokenSwapsWithBalance indicates an expected call of GetTokenSwapsWithBalance.
func (mr *MockRepositoryMockRecorder) GetTokenSwapsWithBalance(ctx, tokenPairIDs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTokenSwapsWithBalance", reflect.TypeOf((*MockRepository)(nil).GetTokenSwapsWithBalance), ctx, tokenPairIDs)
}

// GetTokensByMints mocks base method.
func (m *MockRepository) GetTokensByMints(ctx context.Context, mints []string) ([]*model.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTokensByMints", ctx, mints)
	ret0, _ := ret[0].([]*model.Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTokensByMints indicates an expected call of GetTokensByMints.
func (mr *MockRepositoryMockRecorder) GetTokensByMints(ctx, mints interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTokensByMints", reflect.TypeOf((*MockRepository)(nil).GetTokensByMints), ctx, mints)
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
func (m *MockRepository) GetVaultByAddress(arg0 context.Context, arg1 string) (*VaultWithTokenPair, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVaultByAddress", arg0, arg1)
	ret0, _ := ret[0].(*VaultWithTokenPair)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVaultByAddress indicates an expected call of GetVaultByAddress.
func (mr *MockRepositoryMockRecorder) GetVaultByAddress(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVaultByAddress", reflect.TypeOf((*MockRepository)(nil).GetVaultByAddress), arg0, arg1)
}

// GetVaultPeriods mocks base method.
func (m *MockRepository) GetVaultPeriods(arg0 context.Context, arg1 string, arg2 *string, arg3 PaginationParams) ([]*model.VaultPeriod, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVaultPeriods", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].([]*model.VaultPeriod)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVaultPeriods indicates an expected call of GetVaultPeriods.
func (mr *MockRepositoryMockRecorder) GetVaultPeriods(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVaultPeriods", reflect.TypeOf((*MockRepository)(nil).GetVaultPeriods), arg0, arg1, arg2, arg3)
}

// GetVaultWhitelistsByVaultAddress mocks base method.
func (m *MockRepository) GetVaultWhitelistsByVaultAddress(arg0 context.Context, arg1 []string) ([]*model.VaultWhitelist, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVaultWhitelistsByVaultAddress", arg0, arg1)
	ret0, _ := ret[0].([]*model.VaultWhitelist)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVaultWhitelistsByVaultAddress indicates an expected call of GetVaultWhitelistsByVaultAddress.
func (mr *MockRepositoryMockRecorder) GetVaultWhitelistsByVaultAddress(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVaultWhitelistsByVaultAddress", reflect.TypeOf((*MockRepository)(nil).GetVaultWhitelistsByVaultAddress), arg0, arg1)
}

// GetVaultsWithFilter mocks base method.
func (m *MockRepository) GetVaultsWithFilter(arg0 context.Context, arg1 VaultFilterParams) ([]*VaultWithTokenPair, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVaultsWithFilter", arg0, arg1)
	ret0, _ := ret[0].([]*VaultWithTokenPair)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVaultsWithFilter indicates an expected call of GetVaultsWithFilter.
func (mr *MockRepositoryMockRecorder) GetVaultsWithFilter(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVaultsWithFilter", reflect.TypeOf((*MockRepository)(nil).GetVaultsWithFilter), arg0, arg1)
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

// UpsertOrcaWhirlpools mocks base method.
func (m *MockRepository) UpsertOrcaWhirlpools(arg0 context.Context, arg1 ...*model.OrcaWhirlpool) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpsertOrcaWhirlpools", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertOrcaWhirlpools indicates an expected call of UpsertOrcaWhirlpools.
func (mr *MockRepositoryMockRecorder) UpsertOrcaWhirlpools(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertOrcaWhirlpools", reflect.TypeOf((*MockRepository)(nil).UpsertOrcaWhirlpools), varargs...)
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

// UpsertVaultWhitelists mocks base method.
func (m *MockRepository) UpsertVaultWhitelists(arg0 context.Context, arg1 ...*model.VaultWhitelist) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpsertVaultWhitelists", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertVaultWhitelists indicates an expected call of UpsertVaultWhitelists.
func (mr *MockRepositoryMockRecorder) UpsertVaultWhitelists(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertVaultWhitelists", reflect.TypeOf((*MockRepository)(nil).UpsertVaultWhitelists), varargs...)
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
