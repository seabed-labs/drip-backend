package repository

import (
	"context"

	"github.com/dcaf-labs/drip/pkg/repository/model"
	"github.com/dcaf-labs/drip/pkg/repository/query"
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
	UpsertTokenSwaps(context.Context, ...*model.TokenSwap) error
	UpsertOrcaWhirlpools(context.Context, ...*model.OrcaWhirlpool) error
	UpsertTokenAccountBalances(context.Context, ...*model.TokenAccountBalance) error

	GetVaultByAddress(context.Context, string) (*model.Vault, error)
	GetVaultsWithFilter(context.Context, *string, *string, *string) ([]*model.Vault, error)

	GetVaultWhitelistsByVaultAddress(context.Context, []string) ([]*model.VaultWhitelist, error)

	GetProtoConfigs(ctx context.Context, filterParams ProtoConfigParams) ([]*model.ProtoConfig, error)
	GetProtoConfigsByAddresses(ctx context.Context, pubkeys []string) ([]*model.ProtoConfig, error)

	GetVaultPeriodByAddress(ctx context.Context, address string) (*model.VaultPeriod, error)
	GetVaultPeriods(context.Context, string, *string, PaginationParams) ([]*model.VaultPeriod, error)

	GetAllSupportTokens(ctx context.Context) ([]*model.Token, error)
	GetAllSupportedTokenAs(ctx context.Context) ([]*model.Token, error)
	GetSupportedTokenAs(ctx context.Context, givenTokenBMint *string) ([]*model.Token, error)
	GetSupportedTokenBs(ctx context.Context, givenTokenAMint string) ([]*model.Token, error)

	GetTokenPairByID(context.Context, string) (*model.TokenPair, error)
	GetTokenPair(context.Context, string, string) (*model.TokenPair, error)
	GetTokenPairsByIDS(context.Context, []string) ([]*model.TokenPair, error)
	GetTokenPairs(context.Context, *string, *string) ([]*model.TokenPair, error)
	GetTokensByMints(ctx context.Context, mints []string) ([]*model.Token, error)

	GetTokenSwapByAddress(context.Context, string) (*model.TokenSwap, error)
	GetTokenSwaps(context.Context, []string) ([]*model.TokenSwap, error)
	GetTokenSwapsWithBalance(ctx context.Context, tokenPairIDs []string) ([]TokenSwapWithBalance, error)

	GetOrcaWhirlpoolsByTokenPairIDs(ctx context.Context, tokenPairIDs []string) ([]*model.OrcaWhirlpool, error)
	GetOrcaWhirlpoolByAddress(ctx context.Context, address string) (*model.OrcaWhirlpool, error)

	GetPositionByNFTMint(ctx context.Context, nftMint string) (*model.Position, error)
	GetAdminPositions(ctx context.Context, isVaultEnabled *bool, positionFilterParams PositionFilterParams, paginationParams PaginationParams) ([]*model.Position, error)

	GetTokenAccountBalancesByIDS(context.Context, []string) ([]*model.TokenAccountBalance, error)

	AdminSetVaultEnabled(ctx context.Context, pubkey string, enabled bool) (*model.Vault, error)
	AdminGetVaults(ctx context.Context, vaultFilterParams VaultFilterParams, paginationParams PaginationParams) ([]*model.Vault, error)
	AdminGetVaultByAddress(ctx context.Context, address string) (*model.Vault, error)
	AdminGetVaultsByAddresses(ctx context.Context, addresses ...string) ([]*model.Vault, error)

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
