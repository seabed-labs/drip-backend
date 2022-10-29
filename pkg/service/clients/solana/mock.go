// Code generated by MockGen. DO NOT EDIT.
// Source: client.go

// Package solana is a generated GoMock package.
package solana

import (
	context "context"
	reflect "reflect"

	config "github.com/dcaf-labs/drip/pkg/service/config"
	token_metadata "github.com/gagliardetto/metaplex-go/clients/token-metadata"
	solana "github.com/gagliardetto/solana-go"
	token "github.com/gagliardetto/solana-go/programs/token"
	rpc "github.com/gagliardetto/solana-go/rpc"
	gomock "github.com/golang/mock/gomock"
)

// MockSolana is a mock of Solana interface.
type MockSolana struct {
	ctrl     *gomock.Controller
	recorder *MockSolanaMockRecorder
}

// MockSolanaMockRecorder is the mock recorder for MockSolana.
type MockSolanaMockRecorder struct {
	mock *MockSolana
}

// NewMockSolana creates a new mock instance.
func NewMockSolana(ctrl *gomock.Controller) *MockSolana {
	mock := &MockSolana{ctrl: ctrl}
	mock.recorder = &MockSolanaMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSolana) EXPECT() *MockSolanaMockRecorder {
	return m.recorder
}

// GetAccount mocks base method.
func (m *MockSolana) GetAccount(arg0 context.Context, arg1 string, arg2 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccount", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetAccount indicates an expected call of GetAccount.
func (mr *MockSolanaMockRecorder) GetAccount(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccount", reflect.TypeOf((*MockSolana)(nil).GetAccount), arg0, arg1, arg2)
}

// GetAccountInfo mocks base method.
func (m *MockSolana) GetAccountInfo(arg0 context.Context, arg1 string) (*rpc.GetAccountInfoResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccountInfo", arg0, arg1)
	ret0, _ := ret[0].(*rpc.GetAccountInfoResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccountInfo indicates an expected call of GetAccountInfo.
func (mr *MockSolanaMockRecorder) GetAccountInfo(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccountInfo", reflect.TypeOf((*MockSolana)(nil).GetAccountInfo), arg0, arg1)
}

// GetAccounts mocks base method.
func (m *MockSolana) GetAccounts(arg0 context.Context, arg1 []string, arg2 func(string, []byte)) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccounts", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetAccounts indicates an expected call of GetAccounts.
func (mr *MockSolanaMockRecorder) GetAccounts(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccounts", reflect.TypeOf((*MockSolana)(nil).GetAccounts), arg0, arg1, arg2)
}

// GetLargestTokenAccounts mocks base method.
func (m *MockSolana) GetLargestTokenAccounts(ctx context.Context, mint string) ([]*rpc.TokenLargestAccountsResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLargestTokenAccounts", ctx, mint)
	ret0, _ := ret[0].([]*rpc.TokenLargestAccountsResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLargestTokenAccounts indicates an expected call of GetLargestTokenAccounts.
func (mr *MockSolanaMockRecorder) GetLargestTokenAccounts(ctx, mint interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLargestTokenAccounts", reflect.TypeOf((*MockSolana)(nil).GetLargestTokenAccounts), ctx, mint)
}

// GetNetwork mocks base method.
func (m *MockSolana) GetNetwork() config.Network {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNetwork")
	ret0, _ := ret[0].(config.Network)
	return ret0
}

// GetNetwork indicates an expected call of GetNetwork.
func (mr *MockSolanaMockRecorder) GetNetwork() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNetwork", reflect.TypeOf((*MockSolana)(nil).GetNetwork))
}

// GetProgramAccounts mocks base method.
func (m *MockSolana) GetProgramAccounts(arg0 context.Context, arg1 string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetProgramAccounts", arg0, arg1)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetProgramAccounts indicates an expected call of GetProgramAccounts.
func (mr *MockSolanaMockRecorder) GetProgramAccounts(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetProgramAccounts", reflect.TypeOf((*MockSolana)(nil).GetProgramAccounts), arg0, arg1)
}

// GetTokenMetadataAccount mocks base method.
func (m *MockSolana) GetTokenMetadataAccount(ctx context.Context, mintAddress string) (token_metadata.Metadata, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTokenMetadataAccount", ctx, mintAddress)
	ret0, _ := ret[0].(token_metadata.Metadata)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTokenMetadataAccount indicates an expected call of GetTokenMetadataAccount.
func (mr *MockSolanaMockRecorder) GetTokenMetadataAccount(ctx, mintAddress interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTokenMetadataAccount", reflect.TypeOf((*MockSolana)(nil).GetTokenMetadataAccount), ctx, mintAddress)
}

// GetTokenMint mocks base method.
func (m *MockSolana) GetTokenMint(ctx context.Context, mintAddress string) (token.Mint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTokenMint", ctx, mintAddress)
	ret0, _ := ret[0].(token.Mint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetTokenMint indicates an expected call of GetTokenMint.
func (mr *MockSolanaMockRecorder) GetTokenMint(ctx, mintAddress interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTokenMint", reflect.TypeOf((*MockSolana)(nil).GetTokenMint), ctx, mintAddress)
}

// GetUserBalances mocks base method.
func (m *MockSolana) GetUserBalances(arg0 context.Context, arg1 string) (*rpc.GetTokenAccountsResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserBalances", arg0, arg1)
	ret0, _ := ret[0].(*rpc.GetTokenAccountsResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserBalances indicates an expected call of GetUserBalances.
func (mr *MockSolanaMockRecorder) GetUserBalances(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserBalances", reflect.TypeOf((*MockSolana)(nil).GetUserBalances), arg0, arg1)
}

// GetVersion mocks base method.
func (m *MockSolana) GetVersion(arg0 context.Context) (*rpc.GetVersionResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVersion", arg0)
	ret0, _ := ret[0].(*rpc.GetVersionResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVersion indicates an expected call of GetVersion.
func (mr *MockSolanaMockRecorder) GetVersion(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVersion", reflect.TypeOf((*MockSolana)(nil).GetVersion), arg0)
}

// GetWalletPubKey mocks base method.
func (m *MockSolana) GetWalletPubKey() solana.PublicKey {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWalletPubKey")
	ret0, _ := ret[0].(solana.PublicKey)
	return ret0
}

// GetWalletPubKey indicates an expected call of GetWalletPubKey.
func (mr *MockSolanaMockRecorder) GetWalletPubKey() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWalletPubKey", reflect.TypeOf((*MockSolana)(nil).GetWalletPubKey))
}

// MintToWallet mocks base method.
func (m *MockSolana) MintToWallet(arg0 context.Context, arg1, arg2 string, arg3 uint64) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MintToWallet", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// MintToWallet indicates an expected call of MintToWallet.
func (mr *MockSolanaMockRecorder) MintToWallet(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MintToWallet", reflect.TypeOf((*MockSolana)(nil).MintToWallet), arg0, arg1, arg2, arg3)
}

// ProgramSubscribe mocks base method.
func (m *MockSolana) ProgramSubscribe(arg0 context.Context, arg1 string, arg2 func(string, []byte) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProgramSubscribe", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// ProgramSubscribe indicates an expected call of ProgramSubscribe.
func (mr *MockSolanaMockRecorder) ProgramSubscribe(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProgramSubscribe", reflect.TypeOf((*MockSolana)(nil).ProgramSubscribe), arg0, arg1, arg2)
}

// getWalletPrivKey mocks base method.
func (m *MockSolana) getWalletPrivKey() solana.PrivateKey {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "getWalletPrivKey")
	ret0, _ := ret[0].(solana.PrivateKey)
	return ret0
}

// getWalletPrivKey indicates an expected call of getWalletPrivKey.
func (mr *MockSolanaMockRecorder) getWalletPrivKey() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "getWalletPrivKey", reflect.TypeOf((*MockSolana)(nil).getWalletPrivKey))
}

// signAndBroadcast mocks base method.
func (m *MockSolana) signAndBroadcast(arg0 context.Context, arg1 rpc.CommitmentType, arg2 ...solana.Instruction) (string, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "signAndBroadcast", varargs...)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// signAndBroadcast indicates an expected call of signAndBroadcast.
func (mr *MockSolanaMockRecorder) signAndBroadcast(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "signAndBroadcast", reflect.TypeOf((*MockSolana)(nil).signAndBroadcast), varargs...)
}
