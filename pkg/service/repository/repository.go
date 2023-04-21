package repository

import (
	"context"

	"github.com/dcaf-labs/drip/pkg/service/repository/model"
	"github.com/dcaf-labs/drip/pkg/service/repository/query"
	"github.com/jmoiron/sqlx"
)

const ErrRecordNotFound = "record not found"

// TODO(Mocha): clean this up, likely as separate repo file
type Repository interface {
	InsertTokenPairs(context.Context, ...*model.TokenPair) error
	UpsertProtoConfigs(context.Context, ...*model.ProtoConfig) error
	UpsertTokens(context.Context, ...*model.Token) error
	UpsertVaults(context.Context, ...*model.Vault) error
	UpsertVaultWhitelists(context.Context, ...*model.VaultWhitelist) error
	UpsertVaultPeriods(context.Context, ...*model.VaultPeriod) error
	UpsertPositions(context.Context, ...*model.Position) error
	UpsertTokenAccounts(context.Context, ...*model.TokenAccount) error
	UpsertTokenSwaps(context.Context, ...*model.TokenSwap) error
	UpsertOrcaWhirlpools(context.Context, ...*model.OrcaWhirlpool) error
	UpsertOrcaWhirlpoolDeltaBQuotes(ctx context.Context, quotes ...*model.OrcaWhirlpoolDeltaBQuote) error
	UpsertDepositMetric(ctx context.Context, metrics ...*model.DepositMetric) error
	UpsertDripMetric(ctx context.Context, metrics ...*model.DripMetric) error
	UpsertWithdrawMetric(ctx context.Context, metrics ...*model.WithdrawalMetric) error

	// These will only returned "enabled=true" vaults
	GetVaultByAddress(context.Context, string) (*model.Vault, error)
	GetVaultsWithFilter(context.Context, *string, *string, *string) ([]*model.Vault, error)

	GetVaultWhitelistsByVaultAddress(context.Context, []string) ([]*model.VaultWhitelist, error)

	GetProtoConfigs(ctx context.Context, filterParams ProtoConfigParams) ([]*model.ProtoConfig, error)
	GetProtoConfigByAddress(ctx context.Context, pubkey string) (*model.ProtoConfig, error)
	GetProtoConfigsByAddresses(ctx context.Context, pubkeys []string) ([]*model.ProtoConfig, error)

	GetVaultPeriodByAddress(ctx context.Context, address string) (*model.VaultPeriod, error)
	GetVaultPeriods(context.Context, string, *string, PaginationParams) ([]*model.VaultPeriod, error)

	GetAllSupportedTokens(ctx context.Context) ([]*model.Token, error)
	GetAllSupportedTokenAs(ctx context.Context) ([]*model.Token, error)
	GetSupportedTokenAs(ctx context.Context, givenTokenBMint *string) ([]*model.Token, error)
	GetSupportedTokenBs(ctx context.Context, givenTokenAMint string) ([]*model.Token, error)

	GetTokenPair(context.Context, string, string) (*model.TokenPair, error)
	GetTokenByAddress(ctx context.Context, mint string) (*model.Token, error)
	GetTokensByAddresses(ctx context.Context, mints ...string) ([]*model.Token, error)

	GetTokenSwapByAddress(context.Context, string) (*model.TokenSwap, error)
	GetSPLTokenSwapsByTokenPairIDs(ctx context.Context, tokenPairIDs ...string) ([]*model.TokenSwap, error)

	GetOrcaWhirlpoolsByTokenPairIDs(ctx context.Context, tokenPairIDs ...string) ([]*model.OrcaWhirlpool, error)
	GetOrcaWhirlpoolByAddress(ctx context.Context, address string) (*model.OrcaWhirlpool, error)
	GetOrcaWhirlpoolDeltaBQuote(ctx context.Context, vaultPubkey, whirlpoolPubkey string) (*model.OrcaWhirlpoolDeltaBQuote, error)
	GetOrcaWhirlpoolDeltaBQuoteByVaultAddresses(ctx context.Context, vaultPubkeys ...string) ([]*model.OrcaWhirlpoolDeltaBQuote, error)

	GetPositionByAddress(ctx context.Context, address string) (*model.Position, error)
	GetPositionByNFTMint(ctx context.Context, nftMint string) (*model.Position, error)
	GetAdminPositions(ctx context.Context, isVaultEnabled *bool, positionFilterParams PositionFilterParams, paginationParams PaginationParams) ([]*model.Position, error)

	GetTokenAccountsByAddresses(ctx context.Context, addresses ...string) ([]*model.TokenAccount, error)
	GetActiveTokenAccountsByMint(context.Context, string) ([]*model.TokenAccount, error)

	AdminSetVaultEnabled(ctx context.Context, pubkey string, enabled bool) (*model.Vault, error)
	AdminGetVaults(ctx context.Context, vaultFilterParams VaultFilterParams, paginationParams PaginationParams) ([]*model.Vault, error)
	AdminGetVaultByAddress(ctx context.Context, address string) (*model.Vault, error)
	AdminGetVaultByTreasuryTokenBAccount(ctx context.Context, pubkey string) (*model.Vault, error)
	AdminGetVaultsByAddresses(ctx context.Context, addresses ...string) ([]*model.Vault, error)
	AdminGetVaultsByTokenPairID(ctx context.Context, tokenPairID string) ([]*model.Vault, error)
	GetActiveWallets(ctx context.Context, params GetActiveWalletParams) ([]ActiveWallet, error)

	GetCurrentTVL(ctx context.Context) (*model.CurrentTVL, error)
	GetDepositMetricBySignature(ctx context.Context, signature string) (*model.DepositMetric, error)
	GetDripMetricBySignature(ctx context.Context, signature string) (*model.DripMetric, error)
	GetWithdrawalMetricBySignature(ctx context.Context, signature string) (*model.WithdrawalMetric, error)
}

type repositoryImpl struct {
	repo *query.Query
	db   *sqlx.DB
}

func NewRepository(
	repo *query.Query,
	db *sqlx.DB,
) Repository {
	return repositoryImpl{
		repo: repo,
		db:   db,
	}
}

type TokenSwapWithBalance struct {
	model.TokenSwap
	TokenABalanceAmount uint64 `json:"tokenAccountABalanceAmount" db:"token_account_a_balance_amount"`
	TokenBBalanceAmount uint64 `json:"tokenAccountBBalanceAmount" db:"token_account_b_balance_amount"`
}

type ProtoConfigParams struct {
	TokenA *string
	TokenB *string
}

type GetActiveWalletParams struct {
	PositionIsClosed *bool
	Owner            *string
	Vault            *string
}

type ActiveWallet struct {
	Owner         string `json:"owner" db:"owner"`
	PositionCount int    `json:"position_count" db:"position_count"`
}

type PaginationParams struct {
	Limit  *int
	Offset *int
}

type VaultFilterParams struct {
	IsEnabled        *bool
	TokenA           *string
	TokenB           *string
	Vault            *string
	VaultProtoConfig *string
}

type PositionFilterParams struct {
	IsClosed *bool
	Wallet   *string
}
