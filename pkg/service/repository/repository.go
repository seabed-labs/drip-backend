package repository

import (
	"context"

	model2 "github.com/dcaf-labs/drip/pkg/service/repository/model"
	"github.com/dcaf-labs/drip/pkg/service/repository/query"

	"github.com/jmoiron/sqlx"
)

const ErrRecordNotFound = "record not found"

// TODO(Mocha): clean this up, likely as separate repo file
type Repository interface {
	InsertTokenPairs(context.Context, ...*model2.TokenPair) error
	UpsertProtoConfigs(context.Context, ...*model2.ProtoConfig) error
	UpsertTokens(context.Context, ...*model2.Token) error
	UpsertVaults(context.Context, ...*model2.Vault) error
	UpsertVaultWhitelists(context.Context, ...*model2.VaultWhitelist) error
	UpsertVaultPeriods(context.Context, ...*model2.VaultPeriod) error
	UpsertPositions(context.Context, ...*model2.Position) error
	UpsertTokenSwaps(context.Context, ...*model2.TokenSwap) error
	UpsertOrcaWhirlpools(context.Context, ...*model2.OrcaWhirlpool) error
	UpsertTokenAccountBalances(context.Context, ...*model2.TokenAccountBalance) error

	GetVaultByAddress(context.Context, string) (*model2.Vault, error)
	GetVaultsWithFilter(context.Context, *string, *string, *string) ([]*model2.Vault, error)

	GetVaultWhitelistsByVaultAddress(context.Context, []string) ([]*model2.VaultWhitelist, error)

	GetProtoConfigs(ctx context.Context, filterParams ProtoConfigParams) ([]*model2.ProtoConfig, error)
	GetProtoConfigsByAddresses(ctx context.Context, pubkeys []string) ([]*model2.ProtoConfig, error)

	GetVaultPeriodByAddress(ctx context.Context, address string) (*model2.VaultPeriod, error)
	GetVaultPeriods(context.Context, string, *string, PaginationParams) ([]*model2.VaultPeriod, error)

	GetAllSupportTokens(ctx context.Context) ([]*model2.Token, error)
	GetAllSupportedTokenAs(ctx context.Context) ([]*model2.Token, error)
	GetSupportedTokenAs(ctx context.Context, givenTokenBMint *string) ([]*model2.Token, error)
	GetSupportedTokenBs(ctx context.Context, givenTokenAMint string) ([]*model2.Token, error)

	GetTokenPairByID(context.Context, string) (*model2.TokenPair, error)
	GetTokenPair(context.Context, string, string) (*model2.TokenPair, error)
	GetTokenPairsByIDS(context.Context, []string) ([]*model2.TokenPair, error)
	GetTokenPairs(context.Context, *string, *string) ([]*model2.TokenPair, error)
	GetTokensByMints(ctx context.Context, mints []string) ([]*model2.Token, error)

	GetTokenSwapByAddress(context.Context, string) (*model2.TokenSwap, error)
	GetTokenSwaps(context.Context, []string) ([]*model2.TokenSwap, error)
	GetTokenSwapsWithBalance(ctx context.Context, tokenPairIDs []string) ([]TokenSwapWithBalance, error)

	GetOrcaWhirlpoolsByTokenPairIDs(ctx context.Context, tokenPairIDs []string) ([]*model2.OrcaWhirlpool, error)
	GetOrcaWhirlpoolByAddress(ctx context.Context, address string) (*model2.OrcaWhirlpool, error)

	GetPositionByNFTMint(ctx context.Context, nftMint string) (*model2.Position, error)
	GetAdminPositions(ctx context.Context, isVaultEnabled *bool, positionFilterParams PositionFilterParams, paginationParams PaginationParams) ([]*model2.Position, error)

	GetTokenAccountBalancesByIDS(context.Context, []string) ([]*model2.TokenAccountBalance, error)

	AdminSetVaultEnabled(ctx context.Context, pubkey string, enabled bool) (*model2.Vault, error)
	AdminGetVaults(ctx context.Context, vaultFilterParams VaultFilterParams, paginationParams PaginationParams) ([]*model2.Vault, error)
	AdminGetVaultByAddress(ctx context.Context, address string) (*model2.Vault, error)
	AdminGetVaultsByAddresses(ctx context.Context, addresses ...string) ([]*model2.Vault, error)

	GetActiveWallets(ctx context.Context, params GetActiveWalletParams) ([]ActiveWallet, error)
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
	model2.TokenSwap
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
