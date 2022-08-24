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
	InsertTokenPairs(ctx context.Context, tokenPairs ...*model.TokenPair) error
	UpsertProtoConfigs(ctx context.Context, protoConfigs ...*model.ProtoConfig) error
	UpsertTokens(ctx context.Context, tokens ...*model.Token) error
	UpsertVaults(ctx context.Context, vaults ...*model.Vault) error
	UpsertVaultWhitelists(ctx context.Context, vaultWhitelists ...*model.VaultWhitelist) error
	UpsertVaultPeriods(ctx context.Context, vaultPeriods ...*model.VaultPeriod) error
	UpsertPositions(ctx context.Context, positions ...*model.Position) error
	UpsertTokenSwaps(ctx context.Context, tokenSwaps ...*model.TokenSwap) error
	UpsertOrcaWhirlpools(ctx context.Context, whirlpools ...*model.OrcaWhirlpool) error
	UpsertTokenAccountBalances(ctx context.Context, tokenAccountBalances ...*model.TokenAccountBalance) error

	GetVaultByAddress(ctx context.Context, adress string) (*VaultWithTokenPair, error)
	GetVaultsWithFilter(ctx context.Context, params VaultFilterParams) ([]*VaultWithTokenPair, error)
	SearchVaultsWithFilter(ctx context.Context, params VaultSearchFilterParams) ([]*VaultWithTokenPair, error)

	GetVaultWhitelistsForVaults(ctx context.Context, vaultAddresses ...string) ([]*model.VaultWhitelist, error)

	GetProtoConfigs(ctx context.Context) ([]*model.ProtoConfig, error)
	GetProtoConfigsByAddresses(ctx context.Context, addresses ...string) ([]*model.ProtoConfig, error)

	GetVaultPeriods(ctx context.Context, vault string, vaultPeriodId *string, paginationParams PaginationParams) ([]*model.VaultPeriod, error)

	GetTokenPairByID(ctx context.Context, id string) (*model.TokenPair, error)
	GetTokenPairsByIDS(ctx context.Context, ids []string) ([]*model.TokenPair, error)
	GetTokenPairs(ctx context.Context, params TokenPairFilterParams) ([]*model.TokenPair, error)

	GetTokensByMints(ctx context.Context, mints []string) ([]*model.Token, error)

	GetTokenSwapByAddress(ctx context.Context, address string) (*model.TokenSwap, error)
	GetTokenSwapsByAddresses(ctx context.Context, addresses []string) ([]*model.TokenSwap, error)
	GetTokenSwapsByTokenPairIDsWithBalance(ctx context.Context, tokenPairIDs []string) ([]TokenSwapWithBalance, error)

	GetOrcaWhirlpoolByAddress(ctx context.Context, address string) (*model.OrcaWhirlpool, error)
	GetOrcaWhirlpoolsByTokenPairIDs(ctx context.Context, tokenPairIDs []string) ([]*model.OrcaWhirlpool, error)

	GetPositionByNFTMint(ctx context.Context, nftMint string) (*model.Position, error)
	GetPositionsWithFilter(ctx context.Context, isVaultEnabled *bool, positionFilterParams PositionFilterParams, paginationParams PaginationParams) ([]*model.Position, error)

	GetTokenAccountBalancesByIDS(ctx context.Context, ids []string) ([]*model.TokenAccountBalance, error)

	GetActiveWallets(ctx context.Context, params GetActiveWalletParams) ([]ActiveWallet, error)

	AdminSetVaultEnabled(ctx context.Context, pubkey string, enabled bool) (*model.Vault, error)
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

type VaultFilterLikeParams struct {
	IsEnabled   *bool
	LikeTokenA  *string
	LikeTokenB  *string
	LikeVault   *string
	ProtoConfig *string
}

type VaultFilterParams struct {
	IsEnabled   *bool
	TokenA      *string
	TokenB      *string
	ProtoConfig *string
}

type VaultSearchFilterParams struct {
	IsEnabled       *bool
	LikeTokenA      *string
	LikeTokenB      *string
	LikeVault       *string
	LikeProtoConfig *string
}

type VaultWithTokenPair struct {
	model.Vault
	model.TokenPair
}

type PositionFilterParams struct {
	IsClosed *bool
	Wallet   *string
}

type TokenPairFilterParams struct {
	TokenA *string
	TokenB *string
}
