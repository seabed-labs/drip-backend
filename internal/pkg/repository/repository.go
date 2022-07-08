package repository

import (
	"context"

	"github.com/google/uuid"

	"github.com/lib/pq"

	"github.com/jmoiron/sqlx"

	"gorm.io/gorm/clause"

	"github.com/dcaf-protocol/drip/internal/pkg/clients/solana"
	"github.com/dcaf-protocol/drip/internal/pkg/repository/model"
	"github.com/dcaf-protocol/drip/internal/pkg/repository/query"
)

type TokenSwapWithLiquidityRatio struct {
	model.TokenSwap
	LiquidityRatio float64 `json:"liquidityRatio" db:"liquidity_ratio"`
}

type Repository interface {
	UpsertProtoConfigs(context.Context, ...*model.ProtoConfig) error
	UpsertTokens(context.Context, ...*model.Token) error
	UpsertTokenPairs(context.Context, ...*model.TokenPair) error
	UpsertVaults(context.Context, ...*model.Vault) error
	UpsertVaultPeriods(context.Context, ...*model.VaultPeriod) error
	UpsertPositions(context.Context, ...*model.Position) error
	UpsertTokenSwaps(context.Context, ...*model.TokenSwap) error
	UpsertTokenAccountBalances(context.Context, ...*model.TokenAccountBalance) error

	GetVaultByAddress(context.Context, string) (*model.Vault, error)
	GetVaultsWithFilter(context.Context, *string, *string, *string) ([]*model.Vault, error)
	GetProtoConfigs(context.Context, *string, *string) ([]*model.ProtoConfig, error)
	GetVaultPeriods(context.Context, string, int, int, *string) ([]*model.VaultPeriod, error)
	GetTokensWithSupportedTokenPair(context.Context, *string, bool) ([]*model.Token, error)
	GetTokenPair(context.Context, string, string) (*model.TokenPair, error)
	GetTokenPairByID(context.Context, string) (*model.TokenPair, error)
	GetTokenPairs(context.Context, *string, *string) ([]*model.TokenPair, error)
	GetTokenSwaps(context.Context, []string) ([]*model.TokenSwap, error)
	GetTokenSwapsSortedByLiquidity(ctx context.Context, tokenPairIDs []string) ([]TokenSwapWithLiquidityRatio, error)
	GetTokenSwapForTokenAccount(context.Context, string) (*model.TokenSwap, error)

	InternalGetVaultByAddress(ctx context.Context, pubkey string) (*model.Vault, error)
	EnableVault(ctx context.Context, pubkey string) (*model.Vault, error)
}

type repositoryImpl struct {
	client solana.Solana
	repo   *query.Query
	db     *sqlx.DB
}

func (d repositoryImpl) InternalGetVaultByAddress(ctx context.Context, pubkey string) (*model.Vault, error) {
	return d.repo.
		Vault.WithContext(ctx).
		Where(d.repo.Vault.Pubkey.Eq(pubkey)).
		First()
}

func (d repositoryImpl) EnableVault(ctx context.Context, vaultPubkey string) (*model.Vault, error) {
	var res model.Vault
	_, err := d.repo.Vault.
		WithContext(ctx).
		Returning(&res, res.GetAllColumns()...).
		Where(d.repo.Vault.Pubkey.Eq(vaultPubkey)).
		Update(d.repo.Vault.Enabled, true)
	return &res, err
}

func NewRepository(
	client solana.Solana,
	repo *query.Query,
	db *sqlx.DB,
) Repository {
	return repositoryImpl{
		client: client,
		repo:   repo,
		db:     db,
	}
}

func (d repositoryImpl) UpsertTokenSwaps(ctx context.Context, tokenSwaps ...*model.TokenSwap) error {
	return d.repo.TokenSwap.
		WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "pubkey"}, {Name: "token_a_mint"}, {Name: "token_b_mint"}},
			UpdateAll: true,
		}).
		Create(tokenSwaps...)
}

func (d repositoryImpl) UpsertTokenAccountBalances(ctx context.Context, tokenAccountBalances ...*model.TokenAccountBalance) error {
	return d.repo.TokenAccountBalance.
		WithContext(ctx).
		Clauses(clause.OnConflict{
			UpdateAll: true,
		}).
		Create(tokenAccountBalances...)
}

func (d repositoryImpl) UpsertProtoConfigs(ctx context.Context, protoConfigs ...*model.ProtoConfig) error {
	return d.repo.ProtoConfig.
		WithContext(ctx).
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(protoConfigs...)
}

func (d repositoryImpl) UpsertTokens(ctx context.Context, tokens ...*model.Token) error {
	return d.repo.Token.
		WithContext(ctx).
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(tokens...)
}

func (d repositoryImpl) UpsertTokenPairs(ctx context.Context, tokenPairs ...*model.TokenPair) error {
	return d.repo.TokenPair.
		WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(tokenPairs...)
}

func (d repositoryImpl) UpsertVaults(ctx context.Context, vaults ...*model.Vault) error {
	return d.repo.Vault.
		WithContext(ctx).
		Clauses(clause.OnConflict{
			UpdateAll: true,
		}).
		Create(vaults...)
}

func (d repositoryImpl) UpsertVaultPeriods(ctx context.Context, vaultPeriods ...*model.VaultPeriod) error {
	return d.repo.VaultPeriod.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Create(vaultPeriods...)
}

func (d repositoryImpl) UpsertPositions(ctx context.Context, positions ...*model.Position) error {
	return d.repo.Position.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Create(positions...)
}

func (d repositoryImpl) GetTokenPair(ctx context.Context, tokenA string, tokenB string) (*model.TokenPair, error) {
	return d.repo.TokenPair.WithContext(ctx).
		Where(d.repo.TokenPair.TokenA.Eq(tokenA)).
		Where(d.repo.TokenPair.TokenB.Eq(tokenB)).
		First()
}

func (d repositoryImpl) GetTokenPairByID(ctx context.Context, id string) (*model.TokenPair, error) {
	return d.repo.TokenPair.WithContext(ctx).
		Where(d.repo.TokenPair.ID.Eq(id)).
		First()
}

func (d repositoryImpl) GetTokenPairs(ctx context.Context, tokenAMint *string, tokenBMint *string) ([]*model.TokenPair, error) {
	stmt := d.repo.TokenPair.WithContext(ctx)
	if tokenAMint != nil {
		stmt = stmt.Where(d.repo.TokenPair.TokenA.Eq(*tokenAMint))
	}
	if tokenBMint != nil {
		stmt = stmt.Where(d.repo.TokenPair.TokenB.Eq(*tokenBMint))
	}
	return stmt.Find()
}

func (d repositoryImpl) GetTokenSwaps(ctx context.Context, tokenPairID []string) ([]*model.TokenSwap, error) {
	stmt := d.repo.TokenSwap.WithContext(ctx)
	if len(tokenPairID) > 0 {
		stmt = stmt.Where(d.repo.TokenSwap.TokenPairID.In(tokenPairID...))
	}
	return stmt.Find()
}

func (d repositoryImpl) GetTokenSwapsSortedByLiquidity(ctx context.Context, tokenPairIDs []string) ([]TokenSwapWithLiquidityRatio, error) {
	var tokenSwaps []TokenSwapWithLiquidityRatio
	// TODO(Mocha): No clue how to do this in gorm-gen
	if len(tokenPairIDs) > 0 {
		var temp []uuid.UUID
		for _, tokenPairID := range tokenPairIDs {
			tokenPairUUID, _ := uuid.Parse(tokenPairID)
			temp = append(temp, tokenPairUUID)
		}
		// We should sort by liquidity ratio ascending, so that the largest ratio is at the end of the list
		if err := d.db.SelectContext(ctx,
			&tokenSwaps,
			`SELECT token_swap.*, token_account_b_balance.amount/token_account_a_balance.amount as liquidity_ratio
				FROM token_swap
				JOIN vault
				ON vault.token_pair_id = token_swap.token_pair_id
				JOIN token_account_balance token_account_a_balance
				ON token_account_a_balance.pubkey = token_swap.token_a_account
				JOIN token_account_balance token_account_b_balance
				ON token_account_b_balance.pubkey = token_swap.token_b_account
				WHERE token_account_a_balance.amount != 0
				AND token_account_b_balance.amount != 0
				AND vault.enabled = true
				AND token_swap.token_pair_id=ANY($1)
				ORDER BY token_swap.token_pair_id desc, liquidity_ratio asc;`,
			pq.Array(temp),
		); err != nil {
			return nil, err
		}
	} else {
		if err := d.db.SelectContext(ctx,
			&tokenSwaps,
			`SELECT token_swap.*, token_account_b_balance.amount/token_account_a_balance.amount as liquidity_ratio
				FROM token_swap
				JOIN vault
				ON vault.token_pair_id = token_swap.token_pair_id
				JOIN token_account_balance token_account_a_balance
				ON token_account_a_balance.pubkey = token_swap.token_a_account
				JOIN token_account_balance token_account_b_balance
				ON token_account_b_balance.pubkey = token_swap.token_b_account
				WHERE token_account_a_balance.amount != 0
				AND token_account_b_balance.amount != 0
				AND vault.enabled = true
				ORDER BY liquidity_ratio asc;`,
		); err != nil {
			return nil, err
		}
	}
	return tokenSwaps, nil
}

func (d repositoryImpl) GetTokenSwapForTokenAccount(ctx context.Context, tokenAccount string) (*model.TokenSwap, error) {
	return d.repo.
		TokenSwap.
		WithContext(ctx).
		Where(d.repo.TokenSwap.TokenAAccount.Eq(tokenAccount)).
		Or(d.repo.TokenSwap.TokenBAccount.Eq(tokenAccount)).
		First()
}

func (d repositoryImpl) GetTokensWithSupportedTokenPair(ctx context.Context, tokenMint *string, supportedTokenA bool) ([]*model.Token, error) {
	stmt := d.repo.Token.WithContext(ctx).Distinct(d.repo.Token.ALL)
	if tokenMint != nil {
		if supportedTokenA {
			stmt = stmt.
				Join(d.repo.TokenPair, d.repo.TokenPair.TokenB.EqCol(d.repo.Token.Pubkey)).
				Join(d.repo.Vault, d.repo.Vault.TokenPairID.EqCol(d.repo.TokenPair.ID)).
				Where(d.repo.Vault.Enabled.Is(true)).
				Where(d.repo.TokenPair.TokenA.Eq(*tokenMint))
		} else {
			stmt = stmt.
				Join(d.repo.TokenPair, d.repo.TokenPair.TokenA.EqCol(d.repo.Token.Pubkey)).
				Join(d.repo.Vault, d.repo.Vault.TokenPairID.EqCol(d.repo.TokenPair.ID)).
				Where(d.repo.Vault.Enabled.Is(true)).
				Where(d.repo.TokenPair.TokenB.Eq(*tokenMint))
		}
	}
	return stmt.Find()
}

func (d repositoryImpl) GetTokensWithSupportedTokenB(ctx context.Context, tokenBMint *string) ([]*model.Token, error) {
	stmt := d.repo.Token.WithContext(ctx).Distinct(d.repo.Token.ALL)
	if tokenBMint != nil {
		stmt = stmt.
			Join(d.repo.TokenPair, d.repo.TokenPair.TokenB.EqCol(d.repo.Token.Pubkey)).
			Join(d.repo.Vault, d.repo.Vault.TokenPairID.EqCol(d.repo.TokenPair.ID)).
			Where(d.repo.Vault.Enabled.Is(true))
	}
	return stmt.Find()
}

func (d repositoryImpl) GetVaultsWithFilter(ctx context.Context, tokenAMint, tokenBMint, protoConfig *string) ([]*model.Vault, error) {
	stmt := d.repo.Vault.WithContext(ctx)
	if tokenAMint != nil || tokenBMint != nil {
		stmt = stmt.Join(d.repo.Vault, d.repo.TokenPair.ID.EqCol(d.repo.Vault.TokenPairID))
	}
	if tokenAMint != nil {
		stmt = stmt.Where(d.repo.TokenPair.TokenA.Eq(*tokenAMint))
	}
	if tokenBMint != nil {
		stmt = stmt.Where(d.repo.TokenPair.TokenB.Eq(*tokenBMint))
	}
	if protoConfig != nil {
		stmt = stmt.Where(d.repo.Vault.ProtoConfig.Eq(*protoConfig))
	}
	stmt = stmt.Where(d.repo.Vault.Enabled.Is(true))
	return stmt.Find()
}

func (d repositoryImpl) GetVaultByAddress(ctx context.Context, address string) (*model.Vault, error) {
	return d.repo.
		Vault.WithContext(ctx).
		Where(d.repo.Vault.Pubkey.Eq(address)).
		Where(d.repo.Vault.Enabled.Is(true)).
		First()
}

func (d repositoryImpl) GetProtoConfigs(ctx context.Context, tokenAMint *string, tokenBMint *string) ([]*model.ProtoConfig, error) {
	stmt := d.repo.ProtoConfig.WithContext(ctx)
	stmt = stmt.Join(d.repo.Vault, d.repo.ProtoConfig.Pubkey.EqCol(d.repo.Vault.ProtoConfig))

	if tokenAMint != nil || tokenBMint != nil {
		stmt = stmt.Join(d.repo.Vault, d.repo.TokenPair.ID.EqCol(d.repo.Vault.TokenPairID))
	}
	if tokenAMint != nil {
		stmt = stmt.Where(d.repo.TokenPair.TokenA.Eq(*tokenAMint))
	}
	if tokenBMint != nil {
		stmt = stmt.Where(d.repo.TokenPair.TokenB.Eq(*tokenBMint))
	}
	stmt = stmt.Where(d.repo.Vault.Enabled.Is(true))
	return stmt.Find()
}

func (d repositoryImpl) GetVaultPeriods(ctx context.Context, vault string, limit int, offset int, vaultPeriod *string) ([]*model.VaultPeriod, error) {
	stmt := d.repo.
		VaultPeriod.WithContext(ctx).
		Join(d.repo.Vault, d.repo.VaultPeriod.Vault.EqCol(d.repo.Vault.Pubkey)).
		Where(d.repo.VaultPeriod.Vault.Eq(vault)).
		Where(d.repo.Vault.Enabled.Is(true))
	if vaultPeriod != nil {
		stmt = stmt.Where(d.repo.VaultPeriod.Pubkey.Eq(*vaultPeriod))
	}
	return stmt.Limit(limit).Offset(offset).Find()
}
