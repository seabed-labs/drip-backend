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
func (m *MockRepository) AdminGetVaults(ctx context.Context, vaultFilterParams VaultFilterParams, paginationParams PaginationParams) ([]*model.Vault, error) {
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

// AdminGetVaultsByAddresses mocks base method.
func (m *MockRepository) AdminGetVaultsByAddresses(ctx context.Context, addresses ...string) ([]*model.Vault, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range addresses {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "AdminGetVaultsByAddresses", varargs...)
	ret0, _ := ret[0].([]*model.Vault)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AdminGetVaultsByAddresses indicates an expected call of AdminGetVaultsByAddresses.
func (mr *MockRepositoryMockRecorder) AdminGetVaultsByAddresses(ctx interface{}, addresses ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, addresses...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AdminGetVaultsByAddresses", reflect.TypeOf((*MockRepository)(nil).AdminGetVaultsByAddresses), varargs...)
}

// AdminGetVaultsByTokenPairID mocks base method.
func (m *MockRepository) AdminGetVaultsByTokenPairID(ctx context.Context, tokenPairID string) ([]*model.Vault, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AdminGetVaultsByTokenPairID", ctx, tokenPairID)
	ret0, _ := ret[0].([]*model.Vault)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AdminGetVaultsByTokenPairID indicates an expected call of AdminGetVaultsByTokenPairID.
func (mr *MockRepositoryMockRecorder) AdminGetVaultsByTokenPairID(ctx, tokenPairID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AdminGetVaultsByTokenPairID", reflect.TypeOf((*MockRepository)(nil).AdminGetVaultsByTokenPairID), ctx, tokenPairID)
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

// GetAllSupportTokens mocks base method.
func (m *MockRepository) GetAllSupportTokens(ctx context.Context) ([]*model.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllSupportTokens", ctx)
	ret0, _ := ret[0].([]*model.Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllSupportTokens indicates an expected call of GetAllSupportTokens.
func (mr *MockRepositoryMockRecorder) GetAllSupportTokens(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllSupportTokens", reflect.TypeOf((*MockRepository)(nil).GetAllSupportTokens), ctx)
}

// GetAllSupportedTokenAs mocks base method.
func (m *MockRepository) GetAllSupportedTokenAs(ctx context.Context) ([]*model.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllSupportedTokenAs", ctx)
	ret0, _ := ret[0].([]*model.Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllSupportedTokenAs indicates an expected call of GetAllSupportedTokenAs.
func (mr *MockRepositoryMockRecorder) GetAllSupportedTokenAs(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllSupportedTokenAs", reflect.TypeOf((*MockRepository)(nil).GetAllSupportedTokenAs), ctx)
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

// GetOrcaWhirlpoolDeltaBQuote mocks base method.
func (m *MockRepository) GetOrcaWhirlpoolDeltaBQuote(ctx context.Context, vaultPubkey, whirlpoolPubkey string) (*model.OrcaWhirlpoolDeltaBQuote, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrcaWhirlpoolDeltaBQuote", ctx, vaultPubkey, whirlpoolPubkey)
	ret0, _ := ret[0].(*model.OrcaWhirlpoolDeltaBQuote)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrcaWhirlpoolDeltaBQuote indicates an expected call of GetOrcaWhirlpoolDeltaBQuote.
func (mr *MockRepositoryMockRecorder) GetOrcaWhirlpoolDeltaBQuote(ctx, vaultPubkey, whirlpoolPubkey interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrcaWhirlpoolDeltaBQuote", reflect.TypeOf((*MockRepository)(nil).GetOrcaWhirlpoolDeltaBQuote), ctx, vaultPubkey, whirlpoolPubkey)
}

// GetOrcaWhirlpoolDeltaBQuoteByVaultAddresses mocks base method.
func (m *MockRepository) GetOrcaWhirlpoolDeltaBQuoteByVaultAddresses(ctx context.Context, vaultPubkeys ...string) ([]*model.OrcaWhirlpoolDeltaBQuote, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range vaultPubkeys {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetOrcaWhirlpoolDeltaBQuoteByVaultAddresses", varargs...)
	ret0, _ := ret[0].([]*model.OrcaWhirlpoolDeltaBQuote)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrcaWhirlpoolDeltaBQuoteByVaultAddresses indicates an expected call of GetOrcaWhirlpoolDeltaBQuoteByVaultAddresses.
func (mr *MockRepositoryMockRecorder) GetOrcaWhirlpoolDeltaBQuoteByVaultAddresses(ctx interface{}, vaultPubkeys ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, vaultPubkeys...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrcaWhirlpoolDeltaBQuoteByVaultAddresses", reflect.TypeOf((*MockRepository)(nil).GetOrcaWhirlpoolDeltaBQuoteByVaultAddresses), varargs...)
}

// GetOrcaWhirlpoolsByTokenPairIDs mocks base method.
func (m *MockRepository) GetOrcaWhirlpoolsByTokenPairIDs(ctx context.Context, tokenPairIDs ...string) ([]*model.OrcaWhirlpool, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range tokenPairIDs {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetOrcaWhirlpoolsByTokenPairIDs", varargs...)
	ret0, _ := ret[0].([]*model.OrcaWhirlpool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrcaWhirlpoolsByTokenPairIDs indicates an expected call of GetOrcaWhirlpoolsByTokenPairIDs.
func (mr *MockRepositoryMockRecorder) GetOrcaWhirlpoolsByTokenPairIDs(ctx interface{}, tokenPairIDs ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, tokenPairIDs...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrcaWhirlpoolsByTokenPairIDs", reflect.TypeOf((*MockRepository)(nil).GetOrcaWhirlpoolsByTokenPairIDs), varargs...)
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
func (m *MockRepository) GetProtoConfigs(ctx context.Context, filterParams ProtoConfigParams) ([]*model.ProtoConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProtoConfigs", ctx, filterParams)
	ret0, _ := ret[0].([]*model.ProtoConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProtoConfigs indicates an expected call of GetProtoConfigs.
func (mr *MockRepositoryMockRecorder) GetProtoConfigs(ctx, filterParams interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProtoConfigs", reflect.TypeOf((*MockRepository)(nil).GetProtoConfigs), ctx, filterParams)
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

// GetSPLTokenSwapsByTokenPairIDs mocks base method.
func (m *MockRepository) GetSPLTokenSwapsByTokenPairIDs(ctx context.Context, tokenPairIDs ...string) ([]*model.TokenSwap, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range tokenPairIDs {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetSPLTokenSwapsByTokenPairIDs", varargs...)
	ret0, _ := ret[0].([]*model.TokenSwap)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSPLTokenSwapsByTokenPairIDs indicates an expected call of GetSPLTokenSwapsByTokenPairIDs.
func (mr *MockRepositoryMockRecorder) GetSPLTokenSwapsByTokenPairIDs(ctx interface{}, tokenPairIDs ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, tokenPairIDs...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSPLTokenSwapsByTokenPairIDs", reflect.TypeOf((*MockRepository)(nil).GetSPLTokenSwapsByTokenPairIDs), varargs...)
}

// GetSupportedTokenAs mocks base method.
func (m *MockRepository) GetSupportedTokenAs(ctx context.Context, givenTokenBMint *string) ([]*model.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSupportedTokenAs", ctx, givenTokenBMint)
	ret0, _ := ret[0].([]*model.Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSupportedTokenAs indicates an expected call of GetSupportedTokenAs.
func (mr *MockRepositoryMockRecorder) GetSupportedTokenAs(ctx, givenTokenBMint interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSupportedTokenAs", reflect.TypeOf((*MockRepository)(nil).GetSupportedTokenAs), ctx, givenTokenBMint)
}

// GetSupportedTokenBs mocks base method.
func (m *MockRepository) GetSupportedTokenBs(ctx context.Context, givenTokenAMint string) ([]*model.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSupportedTokenBs", ctx, givenTokenAMint)
	ret0, _ := ret[0].([]*model.Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetSupportedTokenBs indicates an expected call of GetSupportedTokenBs.
func (mr *MockRepositoryMockRecorder) GetSupportedTokenBs(ctx, givenTokenAMint interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSupportedTokenBs", reflect.TypeOf((*MockRepository)(nil).GetSupportedTokenBs), ctx, givenTokenAMint)
}

// GetTokenAccountBalancesByAddresses mocks base method.
func (m *MockRepository) GetTokenAccountBalancesByAddresses(ctx context.Context, addresses ...string) ([]*model.TokenAccountBalance, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range addresses {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetTokenAccountBalancesByAddresses", varargs...)
	ret0, _ := ret[0].([]*model.TokenAccountBalance)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTokenAccountBalancesByAddresses indicates an expected call of GetTokenAccountBalancesByAddresses.
func (mr *MockRepositoryMockRecorder) GetTokenAccountBalancesByAddresses(ctx interface{}, addresses ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, addresses...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTokenAccountBalancesByAddresses", reflect.TypeOf((*MockRepository)(nil).GetTokenAccountBalancesByAddresses), varargs...)
}

// GetTokenByMint mocks base method.
func (m *MockRepository) GetTokenByMint(ctx context.Context, mint string) (*model.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTokenByMint", ctx, mint)
	ret0, _ := ret[0].(*model.Token)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTokenByMint indicates an expected call of GetTokenByMint.
func (mr *MockRepositoryMockRecorder) GetTokenByMint(ctx, mint interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTokenByMint", reflect.TypeOf((*MockRepository)(nil).GetTokenByMint), ctx, mint)
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

// GetVaultPeriodByAddress mocks base method.
func (m *MockRepository) GetVaultPeriodByAddress(ctx context.Context, address string) (*model.VaultPeriod, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVaultPeriodByAddress", ctx, address)
	ret0, _ := ret[0].(*model.VaultPeriod)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVaultPeriodByAddress indicates an expected call of GetVaultPeriodByAddress.
func (mr *MockRepositoryMockRecorder) GetVaultPeriodByAddress(ctx, address interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVaultPeriodByAddress", reflect.TypeOf((*MockRepository)(nil).GetVaultPeriodByAddress), ctx, address)
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

// UpsertOrcaWhirlpoolDeltaBQuotes mocks base method.
func (m *MockRepository) UpsertOrcaWhirlpoolDeltaBQuotes(ctx context.Context, quotes ...*model.OrcaWhirlpoolDeltaBQuote) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range quotes {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UpsertOrcaWhirlpoolDeltaBQuotes", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpsertOrcaWhirlpoolDeltaBQuotes indicates an expected call of UpsertOrcaWhirlpoolDeltaBQuotes.
func (mr *MockRepositoryMockRecorder) UpsertOrcaWhirlpoolDeltaBQuotes(ctx interface{}, quotes ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, quotes...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertOrcaWhirlpoolDeltaBQuotes", reflect.TypeOf((*MockRepository)(nil).UpsertOrcaWhirlpoolDeltaBQuotes), varargs...)
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
