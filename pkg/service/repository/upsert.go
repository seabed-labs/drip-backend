package repository

import (
	context "context"

	model2 "github.com/dcaf-labs/drip/pkg/service/repository/model"

	"gorm.io/gorm/clause"
)

func (d repositoryImpl) UpsertOrcaWhirlpools(ctx context.Context, whirlpools ...*model2.OrcaWhirlpool) error {
	return d.repo.OrcaWhirlpool.
		WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "pubkey"}, {Name: "token_pair_id"}},
			UpdateAll: true,
		}).
		Create(whirlpools...)
}

func (d repositoryImpl) UpsertTokenSwaps(ctx context.Context, tokenSwaps ...*model2.TokenSwap) error {
	return d.repo.TokenSwap.
		WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "pubkey"}, {Name: "token_a_mint"}, {Name: "token_b_mint"}},
			UpdateAll: true,
		}).
		Create(tokenSwaps...)
}

func (d repositoryImpl) UpsertVaultWhitelists(ctx context.Context, vaultWhiteLists ...*model2.VaultWhitelist) error {
	if len(vaultWhiteLists) == 0 {
		return nil
	}
	// Insert new vault whitelists or do no thing
	return d.repo.VaultWhitelist.
		WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "vault_pubkey"}, {Name: "token_swap_pubkey"}},
			DoNothing: true,
		}).
		Create(vaultWhiteLists...)
}

func (d repositoryImpl) InsertTokenPairs(ctx context.Context, tokenPairs ...*model2.TokenPair) error {
	return d.repo.TokenPair.
		WithContext(ctx).
		Clauses(clause.OnConflict{DoNothing: true}).
		Create(tokenPairs...)
}

func (d repositoryImpl) UpsertTokenAccountBalances(ctx context.Context, tokenAccountBalances ...*model2.TokenAccountBalance) error {
	return d.repo.TokenAccountBalance.
		WithContext(ctx).
		Clauses(clause.OnConflict{
			UpdateAll: true,
		}).
		Create(tokenAccountBalances...)
}

func (d repositoryImpl) UpsertProtoConfigs(ctx context.Context, protoConfigs ...*model2.ProtoConfig) error {
	return d.repo.ProtoConfig.
		WithContext(ctx).
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(protoConfigs...)
}

func (d repositoryImpl) UpsertTokens(ctx context.Context, tokens ...*model2.Token) error {
	return d.repo.Token.
		WithContext(ctx).
		Clauses(clause.OnConflict{UpdateAll: true}).
		Create(tokens...)
}

func (d repositoryImpl) UpsertVaults(ctx context.Context, vaults ...*model2.Vault) error {
	// Insert new vaults or update select fields on updates
	return d.repo.Vault.
		WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "pubkey"}},
			DoUpdates: clause.AssignmentColumns([]string{"last_dca_period", "drip_amount", "dca_activation_timestamp"}),
		}).
		Create(vaults...)
}

func (d repositoryImpl) UpsertVaultPeriods(ctx context.Context, vaultPeriods ...*model2.VaultPeriod) error {
	return d.repo.VaultPeriod.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Create(vaultPeriods...)
}

func (d repositoryImpl) UpsertPositions(ctx context.Context, positions ...*model2.Position) error {
	return d.repo.Position.WithContext(ctx).Clauses(clause.OnConflict{UpdateAll: true}).Create(positions...)
}
