package testintegration

import (
	"context"
	"testing"
	"time"

	"github.com/gagliardetto/solana-go"

	"github.com/dcaf-protocol/drip/pkg/repository"
	"github.com/dcaf-protocol/drip/pkg/repository/model"
	"github.com/dcaf-protocol/drip/pkg/repository/query"
	"github.com/dcaf-protocol/drip/pkg/test"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
	"github.com/test-go/testify/assert"
)

// TODO(Mocha): these tests all take a long time because each test fn creates a new DB and runs fresh migrations
// the db setup and migrations can be done once per file opposed to once per fn

//GetVaultByAddress(context.Context, string) (*model2.Vault, error)
//GetVaultsWithFilter(context.Context, *string, *string, *string) ([]*model2.Vault, error)
//GetProtoConfigs(context.Context, *string, *string) ([]*model2.ProtoConfig, error)
//GetVaultPeriods(context.Context, string, int, int, *string) ([]*model2.VaultPeriod, error)
//GetTokensWithSupportedTokenPair(context.Context, *string, bool) ([]*model2.Token, error)
//GetTokenPair(context.Context, string, string) (*model2.TokenPair, error)
//GetTokenPairByID(context.Context, string) (*model2.TokenPair, error)
//GetTokenPairs(context.Context, *string, *string) ([]*model2.TokenPair, error)
//GetTokenSwaps(context.Context, []string) ([]*model2.TokenSwap, error)
//GetTokenSwapsSortedByLiquidity(ctx context.Context, tokenPairIDs []string) ([]TokenSwapWithLiquidityRatio, error)
//GetTokenSwapForTokenAccount(context.Context, string) (*model2.TokenSwap, error)

//InternalGetVaultByAddress(ctx context.Context, pubkey string) (*model2.Vault, error)
//EnableVault(ctx context.Context, pubkey string) (*model2.Vault, error)

//nolint:funlen
func TestUpsertProtoConfigs(t *testing.T) {
	test.InjectDependencies(func(
		repo *query.Query,
		db *sqlx.DB,
	) {
		newRepository := repository.NewRepository(repo, db)

		t.Run("should insert protoConfig", func(t *testing.T) {
			defer truncateDB(t, db)

			protoConfig := model.ProtoConfig{
				Pubkey:               uuid.New().String()[0:4],
				Granularity:          1,
				TriggerDcaSpread:     5,
				BaseWithdrawalSpread: 10,
			}
			err := newRepository.UpsertProtoConfigs(context.Background(), &protoConfig)
			assert.NoError(t, err)

			var insertedConfig model.ProtoConfig
			err = db.Get(&insertedConfig, "select proto_config.* from proto_config where pubkey=$1", protoConfig.Pubkey)
			assert.NoError(t, err)
			assert.Equal(t, protoConfig.Pubkey, insertedConfig.Pubkey)
		})

		t.Run("should insert many protoConfigs", func(t *testing.T) {
			defer truncateDB(t, db)

			protoConfig1 := model.ProtoConfig{
				Pubkey:               uuid.New().String()[0:4],
				Granularity:          1,
				TriggerDcaSpread:     5,
				BaseWithdrawalSpread: 10,
			}
			protoConfig2 := model.ProtoConfig{
				Pubkey:               uuid.New().String()[0:4],
				Granularity:          1,
				TriggerDcaSpread:     5,
				BaseWithdrawalSpread: 10,
			}
			err := newRepository.UpsertProtoConfigs(context.Background(), &protoConfig1, &protoConfig2)
			assert.NoError(t, err)

			var insertedConfig model.ProtoConfig
			err = db.Get(&insertedConfig, "select proto_config.* from proto_config where pubkey=$1", protoConfig1.Pubkey)
			assert.NoError(t, err)
			assert.Equal(t, protoConfig1.Pubkey, insertedConfig.Pubkey)

			err = db.Get(&insertedConfig, "select proto_config.* from proto_config where pubkey=$1", protoConfig2.Pubkey)
			assert.NoError(t, err)
			assert.Equal(t, protoConfig2.Pubkey, insertedConfig.Pubkey)
		})

		t.Run("should update protoConfig", func(t *testing.T) {
			defer truncateDB(t, db)
			seededProtoConfig := seedProtoConfig(t, db, seedProtoConfigParams{})

			protoConfig := model.ProtoConfig{
				Pubkey:               seededProtoConfig.protoConfigPubkey,
				Granularity:          2,
				TriggerDcaSpread:     4,
				BaseWithdrawalSpread: 6,
			}
			err := newRepository.UpsertProtoConfigs(context.Background(), &protoConfig)
			assert.NoError(t, err)

			var updatedProtoConfig model.ProtoConfig
			err = db.Get(&updatedProtoConfig, "select proto_config.* from proto_config where pubkey=$1", seededProtoConfig.protoConfigPubkey)
			assert.NoError(t, err)
			assert.Equal(t, protoConfig.Pubkey, updatedProtoConfig.Pubkey)
			assert.Equal(t, updatedProtoConfig.Granularity, uint64(2))
			assert.Equal(t, updatedProtoConfig.TriggerDcaSpread, uint16(4))
			assert.Equal(t, updatedProtoConfig.BaseWithdrawalSpread, uint16(6))
		})

		t.Run("should update many protoConfigs", func(t *testing.T) {
			defer truncateDB(t, db)
			seededProtoConfig1 := seedProtoConfig(t, db, seedProtoConfigParams{})
			seededProtoConfig2 := seedProtoConfig(t, db, seedProtoConfigParams{})

			protoConfig1 := model.ProtoConfig{
				Pubkey:               seededProtoConfig1.protoConfigPubkey,
				Granularity:          7,
				TriggerDcaSpread:     8,
				BaseWithdrawalSpread: 9,
			}
			protoConfig2 := model.ProtoConfig{
				Pubkey:               seededProtoConfig2.protoConfigPubkey,
				Granularity:          10,
				TriggerDcaSpread:     11,
				BaseWithdrawalSpread: 12,
			}
			err := newRepository.UpsertProtoConfigs(context.Background(), &protoConfig1, &protoConfig2)
			assert.NoError(t, err)

			var updatedProtoConfig model.ProtoConfig
			err = db.Get(&updatedProtoConfig, "select proto_config.* from proto_config where pubkey=$1", seededProtoConfig1.protoConfigPubkey)
			assert.NoError(t, err)
			assert.Equal(t, protoConfig1.Pubkey, updatedProtoConfig.Pubkey)
			assert.Equal(t, updatedProtoConfig.Granularity, uint64(7))
			assert.Equal(t, updatedProtoConfig.TriggerDcaSpread, uint16(8))
			assert.Equal(t, updatedProtoConfig.BaseWithdrawalSpread, uint16(9))

			err = db.Get(&updatedProtoConfig, "select proto_config.* from proto_config where pubkey=$1", seededProtoConfig2.protoConfigPubkey)
			assert.NoError(t, err)
			assert.Equal(t, protoConfig2.Pubkey, updatedProtoConfig.Pubkey)
			assert.Equal(t, updatedProtoConfig.Granularity, uint64(10))
			assert.Equal(t, updatedProtoConfig.TriggerDcaSpread, uint16(11))
			assert.Equal(t, updatedProtoConfig.BaseWithdrawalSpread, uint16(12))
		})
	})
}

//nolint:funlen
func TestUpsertUpsertTokens(t *testing.T) {
	test.InjectDependencies(func(
		repo *query.Query,
		db *sqlx.DB,
	) {
		symbol1 := "btc"
		symbol2 := "eth"
		newRepository := repository.NewRepository(repo, db)

		t.Run("should insert token", func(t *testing.T) {
			defer truncateDB(t, db)

			pubkey := uuid.New().String()[0:4]
			token := model.Token{
				Pubkey:   pubkey,
				Symbol:   &symbol1,
				Decimals: 0,
				IconURL:  nil,
			}
			err := newRepository.UpsertTokens(context.Background(), &token)
			assert.NoError(t, err)

			var insertedToken model.Token
			err = db.Get(&insertedToken, "select token.* from token where pubkey=$1", pubkey)
			assert.NoError(t, err)
			assert.Equal(t, token.Pubkey, insertedToken.Pubkey)
		})

		t.Run("should insert many tokens", func(t *testing.T) {
			defer truncateDB(t, db)

			pubkey1 := uuid.New().String()[0:4]
			pubkey2 := uuid.New().String()[0:4]
			token1 := model.Token{
				Pubkey:   pubkey1,
				Symbol:   &symbol1,
				Decimals: 0,
				IconURL:  nil,
			}
			token2 := model.Token{
				Pubkey:   pubkey2,
				Symbol:   &symbol2,
				Decimals: 0,
				IconURL:  nil,
			}
			err := newRepository.UpsertTokens(context.Background(), &token1, &token2)
			assert.NoError(t, err)

			var insertedToken model.Token
			err = db.Get(&insertedToken, "select token.* from token where pubkey=$1", pubkey1)
			assert.NoError(t, err)
			assert.Equal(t, token1.Pubkey, insertedToken.Pubkey)

			err = db.Get(&insertedToken, "select token.* from token where pubkey=$1", pubkey2)
			assert.NoError(t, err)
			assert.Equal(t, token2.Pubkey, insertedToken.Pubkey)
		})

		t.Run("should update token", func(t *testing.T) {
			defer truncateDB(t, db)
			seededTokens := seedTokens(t, db, seedTokensParams{})

			token := model.Token{
				Pubkey:   seededTokens.tokenAPubkey,
				Symbol:   &symbol2,
				Decimals: 0,
				IconURL:  nil,
			}
			err := newRepository.UpsertTokens(context.Background(), &token)
			assert.NoError(t, err)

			var updatedToken model.Token
			err = db.Get(&updatedToken, "select token.* from token where pubkey=$1", seededTokens.tokenAPubkey)
			assert.NoError(t, err)
			assert.Equal(t, token.Pubkey, updatedToken.Pubkey)
			assert.Equal(t, *updatedToken.Symbol, symbol2)
		})

		t.Run("should update many tokens", func(t *testing.T) {
			defer truncateDB(t, db)
			seededTokens := seedTokens(t, db, seedTokensParams{})

			symbol1 := "sol"
			token1 := model.Token{
				Pubkey:   seededTokens.tokenAPubkey,
				Symbol:   &symbol1,
				Decimals: 0,
				IconURL:  nil,
			}
			symbol2 := "ltc"
			token2 := model.Token{
				Pubkey:   seededTokens.tokenBPubkey,
				Symbol:   &symbol2,
				Decimals: 0,
				IconURL:  nil,
			}
			err := newRepository.UpsertTokens(context.Background(), &token1, &token2)
			assert.NoError(t, err)

			var updatedToken model.Token
			err = db.Get(&updatedToken, "select token.* from token where pubkey=$1", seededTokens.tokenAPubkey)
			assert.NoError(t, err)
			assert.Equal(t, token1.Pubkey, updatedToken.Pubkey)
			assert.Equal(t, *updatedToken.Symbol, symbol1)

			err = db.Get(&updatedToken, "select token.* from token where pubkey=$1", seededTokens.tokenBPubkey)
			assert.NoError(t, err)
			assert.Equal(t, token2.Pubkey, updatedToken.Pubkey)
			assert.Equal(t, *updatedToken.Symbol, symbol2)
		})
	})
}

//nolint:funlen
func TestUpsertTokenPairs(t *testing.T) {
	test.InjectDependencies(func(
		repo *query.Query,
		db *sqlx.DB,
	) {
		newRepository := repository.NewRepository(repo, db)

		t.Run("should fail to insert tokenPair if tokenA doesn't exit", func(t *testing.T) {
			defer truncateDB(t, db)
			seededTokens := seedTokens(t, db, seedTokensParams{})

			tokenPair := model.TokenPair{
				ID:     uuid.New().String(),
				TokenA: uuid.New().String(),
				TokenB: seededTokens.tokenBPubkey,
			}

			err := newRepository.InsertTokenPairs(context.Background(), &tokenPair)
			assert.Error(t, err)
		})

		t.Run("should fail to insert tokenPair if tokenB doesn't exit", func(t *testing.T) {
			defer truncateDB(t, db)
			seededTokens := seedTokens(t, db, seedTokensParams{})

			tokenPair := model.TokenPair{
				ID:     uuid.New().String(),
				TokenA: seededTokens.tokenAPubkey,
				TokenB: uuid.New().String(),
			}

			err := newRepository.InsertTokenPairs(context.Background(), &tokenPair)
			assert.Error(t, err)
		})

		t.Run("should insert tokenPair", func(t *testing.T) {
			defer truncateDB(t, db)
			seededTokens := seedTokens(t, db, seedTokensParams{})

			tokenPair := model.TokenPair{
				ID:     uuid.New().String(),
				TokenA: seededTokens.tokenAPubkey,
				TokenB: seededTokens.tokenBPubkey,
			}

			err := newRepository.InsertTokenPairs(context.Background(), &tokenPair)
			assert.NoError(t, err)

			var insertedTokenPair model.TokenPair
			err = db.Get(&insertedTokenPair, "select token_pair.* from token_pair where id=$1", tokenPair.ID)
			assert.NoError(t, err)
			assert.Equal(t, tokenPair.ID, insertedTokenPair.ID)
			assert.Equal(t, tokenPair.TokenA, insertedTokenPair.TokenA)
			assert.Equal(t, tokenPair.TokenB, insertedTokenPair.TokenB)
		})

		t.Run("should insert many tokenPairs", func(t *testing.T) {
			defer truncateDB(t, db)
			seededTokens := seedTokens(t, db, seedTokensParams{})

			tokenPair1 := model.TokenPair{
				ID:     uuid.New().String(),
				TokenA: seededTokens.tokenAPubkey,
				TokenB: seededTokens.tokenBPubkey,
			}

			tokenPair2 := model.TokenPair{
				ID:     uuid.New().String(),
				TokenA: seededTokens.tokenBPubkey,
				TokenB: seededTokens.tokenAPubkey,
			}

			err := newRepository.InsertTokenPairs(context.Background(), &tokenPair1, &tokenPair2)
			assert.NoError(t, err)

			var insertedTokenPair model.TokenPair
			err = db.Get(&insertedTokenPair, "select token_pair.* from token_pair where id=$1", tokenPair1.ID)
			assert.NoError(t, err)
			assert.Equal(t, tokenPair1.ID, insertedTokenPair.ID)
			assert.Equal(t, tokenPair1.TokenA, insertedTokenPair.TokenA)
			assert.Equal(t, tokenPair1.TokenB, insertedTokenPair.TokenB)

			err = db.Get(&insertedTokenPair, "select token_pair.* from token_pair where id=$1", tokenPair2.ID)
			assert.NoError(t, err)
			assert.Equal(t, tokenPair2.ID, insertedTokenPair.ID)
			assert.Equal(t, tokenPair2.TokenA, insertedTokenPair.TokenA)
			assert.Equal(t, tokenPair2.TokenB, insertedTokenPair.TokenB)
		})

		t.Run("should not update tokenPair if it already exists", func(t *testing.T) {
			defer truncateDB(t, db)

			seededTokenPair := seedTokenPair(t, db, seedTokenPairParams{})

			tokenPair := model.TokenPair{
				ID:     uuid.New().String(),
				TokenA: seededTokenPair.tokenAPubkey,
				TokenB: seededTokenPair.tokenBPubkey,
			}

			err := newRepository.InsertTokenPairs(context.Background(), &tokenPair)
			assert.NoError(t, err)

			var insertedTokenPair model.TokenPair
			err = db.Get(&insertedTokenPair, "select token_pair.* from token_pair where id=$1", seededTokenPair.tokenPairID)
			assert.NoError(t, err)
			assert.Equal(t, seededTokenPair.tokenPairID, insertedTokenPair.ID)
			assert.Equal(t, seededTokenPair.tokenAPubkey, insertedTokenPair.TokenA)
			assert.Equal(t, seededTokenPair.tokenBPubkey, insertedTokenPair.TokenB)
		})
	})
}

//nolint:funlen
func TestUpsertVaults(t *testing.T) {
	test.InjectDependencies(func(
		repo *query.Query,
		db *sqlx.DB,
	) {
		newRepository := repository.NewRepository(repo, db)

		t.Run("should fail to insert vault when protoConfig is missing", func(t *testing.T) {
			defer truncateDB(t, db)
			seededTokenPair := seedTokenPair(t, db, seedTokenPairParams{})

			vault := model.Vault{
				Pubkey:                 uuid.New().String(),
				ProtoConfig:            uuid.New().String(),
				TokenAAccount:          uuid.New().String(),
				TokenBAccount:          uuid.New().String(),
				TreasuryTokenBAccount:  uuid.New().String(),
				LastDcaPeriod:          0,
				DripAmount:             0,
				DcaActivationTimestamp: time.Time{},
				Enabled:                false,
				TokenPairID:            seededTokenPair.tokenPairID,
			}
			err := newRepository.UpsertVaults(context.Background(), &vault)
			assert.Error(t, err)
		})

		t.Run("should fail to insert vault when token pair is missing", func(t *testing.T) {
			defer truncateDB(t, db)

			seededProtoConfig := seedProtoConfig(t, db, seedProtoConfigParams{})
			vault := model.Vault{
				Pubkey:                 uuid.New().String(),
				ProtoConfig:            seededProtoConfig.protoConfigPubkey,
				TokenAAccount:          uuid.New().String(),
				TokenBAccount:          uuid.New().String(),
				TreasuryTokenBAccount:  uuid.New().String(),
				LastDcaPeriod:          0,
				DripAmount:             0,
				DcaActivationTimestamp: time.Time{},
				Enabled:                false,
				TokenPairID:            uuid.New().String(),
			}
			err := newRepository.UpsertVaults(context.Background(), &vault)
			assert.Error(t, err)
		})

		t.Run("should insert vault", func(t *testing.T) {
			defer truncateDB(t, db)

			seededProtoConfig := seedProtoConfig(t, db, seedProtoConfigParams{})
			seededTokenPair := seedTokenPair(t, db, seedTokenPairParams{})

			vault := model.Vault{
				Pubkey:                 uuid.New().String(),
				ProtoConfig:            seededProtoConfig.protoConfigPubkey,
				TokenAAccount:          uuid.New().String(),
				TokenBAccount:          uuid.New().String(),
				TreasuryTokenBAccount:  uuid.New().String(),
				LastDcaPeriod:          0,
				DripAmount:             0,
				DcaActivationTimestamp: time.Time{},
				Enabled:                false,
				TokenPairID:            seededTokenPair.tokenPairID,
			}
			err := newRepository.UpsertVaults(context.Background(), &vault)
			assert.NoError(t, err)

			var insertedVault model.Vault
			err = db.Get(&insertedVault, "select vault.* from vault where pubkey=$1", vault.Pubkey)
			assert.NoError(t, err)
			assert.Equal(t, vault.Pubkey, insertedVault.Pubkey)
			assert.Equal(t, insertedVault.TokenPairID, seededTokenPair.tokenPairID)
		})

		t.Run("should insert many vaults", func(t *testing.T) {
			defer truncateDB(t, db)

			seededProtoConfig := seedProtoConfig(t, db, seedProtoConfigParams{})
			seededTokenPair := seedTokenPair(t, db, seedTokenPairParams{})

			vault1 := model.Vault{
				Pubkey:                 uuid.New().String(),
				ProtoConfig:            seededProtoConfig.protoConfigPubkey,
				TokenAAccount:          uuid.New().String(),
				TokenBAccount:          uuid.New().String(),
				TreasuryTokenBAccount:  uuid.New().String(),
				LastDcaPeriod:          0,
				DripAmount:             0,
				DcaActivationTimestamp: time.Time{},
				Enabled:                false,
				TokenPairID:            seededTokenPair.tokenPairID,
			}
			vault2 := model.Vault{
				Pubkey:                 uuid.New().String(),
				ProtoConfig:            seededProtoConfig.protoConfigPubkey,
				TokenAAccount:          uuid.New().String(),
				TokenBAccount:          uuid.New().String(),
				TreasuryTokenBAccount:  uuid.New().String(),
				LastDcaPeriod:          0,
				DripAmount:             0,
				DcaActivationTimestamp: time.Time{},
				Enabled:                false,
				TokenPairID:            seededTokenPair.tokenPairID,
			}
			err := newRepository.UpsertVaults(context.Background(), &vault1, &vault2)
			assert.NoError(t, err)

			var insertedVault model.Vault
			err = db.Get(&insertedVault, "select vault.* from vault where pubkey=$1", vault1.Pubkey)
			assert.NoError(t, err)
			assert.Equal(t, vault1.Pubkey, insertedVault.Pubkey)
			assert.Equal(t, insertedVault.TokenPairID, seededTokenPair.tokenPairID)

			err = db.Get(&insertedVault, "select vault.* from vault where pubkey=$1", vault2.Pubkey)
			assert.NoError(t, err)
			assert.Equal(t, vault2.Pubkey, insertedVault.Pubkey)
			assert.Equal(t, insertedVault.TokenPairID, seededTokenPair.tokenPairID)
		})

		t.Run("should update vault", func(t *testing.T) {
			defer truncateDB(t, db)

			seedVault := seedVault(t, db, seedVaultParams{})
			vault := model.Vault{
				Pubkey:                 seedVault.vaultPubkey,
				ProtoConfig:            seedVault.protoConfigPubkey,
				TokenAAccount:          seedVault.tokenAAcount,
				TokenBAccount:          seedVault.tokenBAccount,
				TreasuryTokenBAccount:  seedVault.treasuryAccount,
				LastDcaPeriod:          1,
				DripAmount:             100,
				DcaActivationTimestamp: time.Now(),
				Enabled:                true,
				TokenPairID:            seedVault.tokenPairID,
			}
			err := newRepository.UpsertVaults(context.Background(), &vault)
			assert.NoError(t, err)

			var insertedVault model.Vault
			err = db.Get(&insertedVault, "select vault.* from vault where pubkey=$1", seedVault.vaultPubkey)
			assert.NoError(t, err)
			assert.Equal(t, vault.Pubkey, insertedVault.Pubkey)
			assert.Equal(t, insertedVault.LastDcaPeriod, uint64(1))
			assert.Equal(t, insertedVault.DripAmount, uint64(100))
			assert.NotEqual(t, insertedVault.LastDcaPeriod, time.Time{})
			assert.Equal(t, insertedVault.Enabled, true)
		})
	})
}

//nolint:funlen
func TestUpsertVaultPeriod(t *testing.T) {
	test.InjectDependencies(func(
		repo *query.Query,
		db *sqlx.DB,
	) {
		newRepository := repository.NewRepository(repo, db)

		t.Run("should fail to insert vaultPeriod when vault doesn't exist", func(t *testing.T) {
			defer truncateDB(t, db)

			vaultPeriod := model.VaultPeriod{
				Pubkey:   uuid.New().String(),
				Vault:    uuid.New().String(),
				PeriodID: 0,
				Twap:     decimal.NewFromInt(0),
				Dar:      0,
			}
			err := newRepository.UpsertVaultPeriods(context.Background(), &vaultPeriod)
			assert.Error(t, err)
		})

		t.Run("should insert vaultPeriod", func(t *testing.T) {
			defer truncateDB(t, db)

			seededVault := seedVault(t, db, seedVaultParams{})

			vaultPeriod := model.VaultPeriod{
				Pubkey:   uuid.New().String(),
				Vault:    seededVault.vaultPubkey,
				PeriodID: 0,
				Twap:     decimal.NewFromInt(0),
				Dar:      0,
			}
			err := newRepository.UpsertVaultPeriods(context.Background(), &vaultPeriod)
			assert.NoError(t, err)

			var insertedVaultPeriod model.VaultPeriod
			err = db.Get(&insertedVaultPeriod, "select vault_period.* from vault_period where pubkey=$1", vaultPeriod.Pubkey)
			assert.NoError(t, err)
			assert.Equal(t, vaultPeriod.Pubkey, insertedVaultPeriod.Pubkey)
		})

		t.Run("should insert many vaultPeriods", func(t *testing.T) {
			defer truncateDB(t, db)

			seededVault := seedVault(t, db, seedVaultParams{})

			vaultPeriod1 := model.VaultPeriod{
				Pubkey:   uuid.New().String(),
				Vault:    seededVault.vaultPubkey,
				PeriodID: 0,
				Twap:     decimal.NewFromInt(0),
				Dar:      0,
			}
			vaultPeriod2 := model.VaultPeriod{
				Pubkey:   uuid.New().String(),
				Vault:    seededVault.vaultPubkey,
				PeriodID: 1,
				Twap:     decimal.NewFromInt(10),
				Dar:      0,
			}
			err := newRepository.UpsertVaultPeriods(context.Background(), &vaultPeriod1, &vaultPeriod2)
			assert.NoError(t, err)

			var insertedVaultPeriod model.VaultPeriod
			err = db.Get(&insertedVaultPeriod, "select vault_period.* from vault_period where pubkey=$1", vaultPeriod1.Pubkey)
			assert.NoError(t, err)
			assert.Equal(t, vaultPeriod1.Pubkey, insertedVaultPeriod.Pubkey)
			assert.Equal(t, vaultPeriod1.Twap, decimal.NewFromInt(0))

			err = db.Get(&insertedVaultPeriod, "select vault_period.* from vault_period where pubkey=$1", vaultPeriod2.Pubkey)
			assert.NoError(t, err)
			assert.Equal(t, vaultPeriod2.Pubkey, insertedVaultPeriod.Pubkey)
			assert.Equal(t, vaultPeriod2.Twap, decimal.NewFromInt(10))
		})

		t.Run("should update vaultPeriod", func(t *testing.T) {
			defer truncateDB(t, db)
			seededVaultPeriod1 := seedVaultPeriod(t, db, seedVaultPeriodParams{})
			vaultPeriod1 := model.VaultPeriod{
				Pubkey:   seededVaultPeriod1.vaultPeriodPubkey,
				Vault:    seededVaultPeriod1.vaultPubkey,
				PeriodID: 1,
				Twap:     decimal.NewFromInt(0),
				Dar:      0,
			}
			err := newRepository.UpsertVaultPeriods(context.Background(), &vaultPeriod1)
			assert.NoError(t, err)

			var insertedVaultPeriod model.VaultPeriod
			err = db.Get(&insertedVaultPeriod, "select vault_period.* from vault_period where pubkey=$1", vaultPeriod1.Pubkey)
			assert.NoError(t, err)
			assert.Equal(t, vaultPeriod1.Pubkey, insertedVaultPeriod.Pubkey)
			assert.Equal(t, vaultPeriod1.PeriodID, uint64(1))
			assert.Equal(t, vaultPeriod1.Twap, decimal.NewFromInt(0))
		})

		t.Run("should update many vaultPeriods", func(t *testing.T) {
			defer truncateDB(t, db)
			seededVaultPeriod1 := seedVaultPeriod(t, db, seedVaultPeriodParams{})
			vaultPeriod1 := model.VaultPeriod{
				Pubkey:   seededVaultPeriod1.vaultPeriodPubkey,
				Vault:    seededVaultPeriod1.vaultPubkey,
				PeriodID: 1,
				Twap:     decimal.NewFromInt(0),
				Dar:      0,
			}

			seededVaultPeriod2 := seedVaultPeriod(t, db, seedVaultPeriodParams{
				seedVaultParams: seedVaultParams{
					seedTokenPairParams: seedTokenPairParams{
						seedTokensParams: seedTokensParams{
							tokenAPubkey: &seededVaultPeriod1.tokenAPubkey,
							tokenBPubkey: &seededVaultPeriod1.tokenBPubkey,
						},
						tokenPairID: &seededVaultPeriod1.tokenPairID,
					},
				},
				vaultPubkey: nil,
			})
			vaultPeriod2 := model.VaultPeriod{
				Pubkey:   seededVaultPeriod2.vaultPeriodPubkey,
				Vault:    seededVaultPeriod2.vaultPubkey,
				PeriodID: 2,
				Twap:     decimal.NewFromInt(10),
				Dar:      0,
			}

			err := newRepository.UpsertVaultPeriods(context.Background(), &vaultPeriod1, &vaultPeriod2)
			assert.NoError(t, err)

			var insertedVaultPeriod model.VaultPeriod
			err = db.Get(&insertedVaultPeriod, "select vault_period.* from vault_period where pubkey=$1", vaultPeriod1.Pubkey)
			assert.NoError(t, err)
			assert.Equal(t, vaultPeriod1.Pubkey, insertedVaultPeriod.Pubkey)
			assert.Equal(t, vaultPeriod1.PeriodID, uint64(1))
			assert.Equal(t, vaultPeriod1.Twap, decimal.NewFromInt(0))

			err = db.Get(&insertedVaultPeriod, "select vault_period.* from vault_period where pubkey=$1", vaultPeriod2.Pubkey)
			assert.NoError(t, err)
			assert.Equal(t, vaultPeriod2.Pubkey, insertedVaultPeriod.Pubkey)
			assert.Equal(t, vaultPeriod2.PeriodID, uint64(2))
			assert.Equal(t, vaultPeriod2.Twap, decimal.NewFromInt(10))
		})
	})
}

//nolint:funlen
func TestUpsertPositions(t *testing.T) {
	test.InjectDependencies(func(
		repo *query.Query,
		db *sqlx.DB,
	) {
		newRepository := repository.NewRepository(repo, db)

		t.Run("should fail to insert position if vault doesn't exist", func(t *testing.T) {
			defer truncateDB(t, db)

			position := model.Position{
				Pubkey:                   uuid.New().String(),
				Vault:                    uuid.New().String(),
				Authority:                uuid.New().String(),
				DepositedTokenAAmount:    1,
				WithdrawnTokenBAmount:    2,
				DepositTimestamp:         time.Now(),
				DcaPeriodIDBeforeDeposit: 3,
				NumberOfSwaps:            4,
				PeriodicDripAmount:       5,
				IsClosed:                 false,
			}
			err := newRepository.UpsertPositions(context.Background(), &position)
			assert.Error(t, err)
		})

		t.Run("should insert multiple positions", func(t *testing.T) {
			defer truncateDB(t, db)
			seededVault := seedVault(t, db, seedVaultParams{})

			position1 := model.Position{
				Pubkey:                   uuid.New().String(),
				Vault:                    seededVault.vaultPubkey,
				Authority:                uuid.New().String(),
				DepositedTokenAAmount:    1,
				WithdrawnTokenBAmount:    2,
				DepositTimestamp:         time.Now(),
				DcaPeriodIDBeforeDeposit: 3,
				NumberOfSwaps:            4,
				PeriodicDripAmount:       5,
				IsClosed:                 false,
			}
			position2 := model.Position{
				Pubkey:                   uuid.New().String(),
				Vault:                    seededVault.vaultPubkey,
				Authority:                "",
				DepositedTokenAAmount:    2,
				WithdrawnTokenBAmount:    3,
				DepositTimestamp:         time.Now(),
				DcaPeriodIDBeforeDeposit: 4,
				NumberOfSwaps:            5,
				PeriodicDripAmount:       6,
				IsClosed:                 false,
			}
			err := newRepository.UpsertPositions(context.Background(), &position1, &position2)
			assert.NoError(t, err)
			var insertedPosition model.Position
			err = db.Get(&insertedPosition, "select position.* from position where pubkey=$1", position1.Pubkey)
			assert.NoError(t, err)
			assert.Equal(t, position1.Vault, insertedPosition.Vault)
			assert.Equal(t, position1.DepositedTokenAAmount, insertedPosition.DepositedTokenAAmount)
			assert.Equal(t, position1.WithdrawnTokenBAmount, insertedPosition.WithdrawnTokenBAmount)
			assert.Equal(t, position1.DcaPeriodIDBeforeDeposit, insertedPosition.DcaPeriodIDBeforeDeposit)
			assert.Equal(t, position1.NumberOfSwaps, insertedPosition.NumberOfSwaps)
			assert.Equal(t, position1.PeriodicDripAmount, insertedPosition.PeriodicDripAmount)
			assert.NotEqual(t, insertedPosition.DepositTimestamp, time.Time{})

			err = db.Get(&insertedPosition, "select position.* from position where pubkey=$1", position2.Pubkey)
			assert.NoError(t, err)
			assert.Equal(t, position2.Vault, insertedPosition.Vault)
			assert.Equal(t, position2.DepositedTokenAAmount, insertedPosition.DepositedTokenAAmount)
			assert.Equal(t, position2.WithdrawnTokenBAmount, insertedPosition.WithdrawnTokenBAmount)
			assert.Equal(t, position2.DcaPeriodIDBeforeDeposit, insertedPosition.DcaPeriodIDBeforeDeposit)
			assert.Equal(t, position2.NumberOfSwaps, insertedPosition.NumberOfSwaps)
			assert.Equal(t, position2.PeriodicDripAmount, insertedPosition.PeriodicDripAmount)
			assert.NotEqual(t, insertedPosition.DepositTimestamp, time.Time{})
		})

		t.Run("should update position", func(t *testing.T) {
			defer truncateDB(t, db)
			seededPosition := seedPosition(t, db, seedPositionParams{})

			position := model.Position{
				Pubkey:                   seededPosition.positionPubkey,
				Vault:                    seededPosition.vaultPubkey,
				Authority:                uuid.New().String(),
				DepositedTokenAAmount:    1,
				WithdrawnTokenBAmount:    2,
				DepositTimestamp:         time.Now(),
				DcaPeriodIDBeforeDeposit: 3,
				NumberOfSwaps:            4,
				PeriodicDripAmount:       5,
				IsClosed:                 false,
			}
			err := newRepository.UpsertPositions(context.Background(), &position)
			assert.NoError(t, err)

			var updatedPosition model.Position
			err = db.Get(&updatedPosition, "select position.* from position where pubkey=$1", position.Pubkey)
			assert.NoError(t, err)
			assert.Equal(t, position.Vault, updatedPosition.Vault)
			assert.Equal(t, position.DepositedTokenAAmount, updatedPosition.DepositedTokenAAmount)
			assert.Equal(t, position.WithdrawnTokenBAmount, updatedPosition.WithdrawnTokenBAmount)
			assert.Equal(t, position.DcaPeriodIDBeforeDeposit, updatedPosition.DcaPeriodIDBeforeDeposit)
			assert.Equal(t, position.NumberOfSwaps, updatedPosition.NumberOfSwaps)
			assert.Equal(t, position.PeriodicDripAmount, updatedPosition.PeriodicDripAmount)
			assert.NotEqual(t, updatedPosition.DepositTimestamp, time.Time{})
		})

		t.Run("should update many positions", func(t *testing.T) {
			defer truncateDB(t, db)
			seededPosition1 := seedPosition(t, db, seedPositionParams{})
			seededPosition2 := seedPosition(t, db, seedPositionParams{
				seedVaultParams: seedVaultParams{
					seedProtoConfigParams: seedProtoConfigParams{protoConfigPubkey: &seededPosition1.protoConfigPubkey},
					seedTokenPairParams: seedTokenPairParams{
						seedTokensParams: seedTokensParams{
							tokenAPubkey: &seededPosition1.tokenAPubkey,
							tokenBPubkey: &seededPosition1.tokenBPubkey,
						},
						tokenPairID: &seededPosition1.tokenPairID,
					},
					tokenAAcount:    &seededPosition1.tokenAAcount,
					tokenBAccount:   &seededPosition1.tokenBAccount,
					treasuryAccount: &seededPosition1.treasuryAccount,
				},
			})
			position1 := model.Position{
				Pubkey:                   seededPosition1.positionPubkey,
				Vault:                    seededPosition1.vaultPubkey,
				Authority:                uuid.New().String(),
				DepositedTokenAAmount:    1,
				WithdrawnTokenBAmount:    2,
				DepositTimestamp:         time.Now(),
				DcaPeriodIDBeforeDeposit: 3,
				NumberOfSwaps:            4,
				PeriodicDripAmount:       5,
				IsClosed:                 false,
			}

			position2 := model.Position{
				Pubkey:                   seededPosition2.positionPubkey,
				Vault:                    seededPosition2.vaultPubkey,
				Authority:                uuid.New().String(),
				DepositedTokenAAmount:    2,
				WithdrawnTokenBAmount:    3,
				DepositTimestamp:         time.Now(),
				DcaPeriodIDBeforeDeposit: 4,
				NumberOfSwaps:            5,
				PeriodicDripAmount:       6,
				IsClosed:                 false,
			}
			err := newRepository.UpsertPositions(context.Background(), &position1, &position2)
			assert.NoError(t, err)

			var updatedPosition model.Position
			err = db.Get(&updatedPosition, "select position.* from position where pubkey=$1", position1.Pubkey)
			assert.NoError(t, err)
			assert.Equal(t, position1.Vault, updatedPosition.Vault)
			assert.Equal(t, position1.DepositedTokenAAmount, updatedPosition.DepositedTokenAAmount)
			assert.Equal(t, position1.WithdrawnTokenBAmount, updatedPosition.WithdrawnTokenBAmount)
			assert.Equal(t, position1.DcaPeriodIDBeforeDeposit, updatedPosition.DcaPeriodIDBeforeDeposit)
			assert.Equal(t, position1.NumberOfSwaps, updatedPosition.NumberOfSwaps)
			assert.Equal(t, position1.PeriodicDripAmount, updatedPosition.PeriodicDripAmount)
			assert.NotEqual(t, updatedPosition.DepositTimestamp, time.Time{})

			err = db.Get(&updatedPosition, "select position.* from position where pubkey=$1", position2.Pubkey)
			assert.NoError(t, err)
			assert.Equal(t, position2.Vault, updatedPosition.Vault)
			assert.Equal(t, position2.DepositedTokenAAmount, updatedPosition.DepositedTokenAAmount)
			assert.Equal(t, position2.WithdrawnTokenBAmount, updatedPosition.WithdrawnTokenBAmount)
			assert.Equal(t, position2.DcaPeriodIDBeforeDeposit, updatedPosition.DcaPeriodIDBeforeDeposit)
			assert.Equal(t, position2.NumberOfSwaps, updatedPosition.NumberOfSwaps)
			assert.Equal(t, position2.PeriodicDripAmount, updatedPosition.PeriodicDripAmount)
			assert.NotEqual(t, updatedPosition.DepositTimestamp, time.Time{})
		})
	})
}

//nolint:funlen
func TestUpsertTokenSwaps(t *testing.T) {
	test.InjectDependencies(func(
		repo *query.Query,
		db *sqlx.DB,
	) {
		newRepository := repository.NewRepository(repo, db)

		t.Run("should insert tokenSwap", func(t *testing.T) {
			defer truncateDB(t, db)
			seededTokenPair := seedTokenPair(t, db, seedTokenPairParams{})

			tokenSwap := model.TokenSwap{
				ID:            uuid.New().String(),
				Pubkey:        solana.NewWallet().PublicKey().String(),
				Mint:          solana.NewWallet().PublicKey().String(),
				Authority:     "",
				FeeAccount:    "",
				TokenAAccount: solana.NewWallet().PublicKey().String(),
				TokenBAccount: solana.NewWallet().PublicKey().String(),
				TokenPairID:   seededTokenPair.tokenPairID,
				TokenAMint:    seededTokenPair.tokenAPubkey,
				TokenBMint:    seededTokenPair.tokenBPubkey,
			}
			err := newRepository.UpsertTokenSwaps(context.Background(), &tokenSwap)
			assert.NoError(t, err)

			var insertedTokenSwap model.TokenSwap
			err = db.Get(&insertedTokenSwap, "select token_swap.* from token_swap where id=$1", tokenSwap.ID)
			assert.NoError(t, err)
			assert.Equal(t, tokenSwap.Pubkey, insertedTokenSwap.Pubkey)
		})

		t.Run("should insert multiple tokenSwaps", func(t *testing.T) {
			defer truncateDB(t, db)
			seededTokenPair1 := seedTokenPair(t, db, seedTokenPairParams{})
			seededTokenPair2 := seedTokenPair(t, db, seedTokenPairParams{})
			tokenSwap1 := model.TokenSwap{
				ID:            uuid.New().String(),
				Pubkey:        solana.NewWallet().PublicKey().String(),
				Mint:          solana.NewWallet().PublicKey().String(),
				Authority:     solana.NewWallet().PublicKey().String(),
				FeeAccount:    solana.NewWallet().PublicKey().String(),
				TokenAAccount: solana.NewWallet().PublicKey().String(),
				TokenBAccount: solana.NewWallet().PublicKey().String(),
				TokenPairID:   seededTokenPair1.tokenPairID,
				TokenAMint:    seededTokenPair1.tokenAPubkey,
				TokenBMint:    seededTokenPair1.tokenBPubkey,
			}
			tokenSwap2 := model.TokenSwap{
				ID:            uuid.New().String(),
				Pubkey:        solana.NewWallet().PublicKey().String(),
				Mint:          solana.NewWallet().PublicKey().String(),
				Authority:     solana.NewWallet().PublicKey().String(),
				FeeAccount:    solana.NewWallet().PublicKey().String(),
				TokenAAccount: solana.NewWallet().PublicKey().String(),
				TokenBAccount: solana.NewWallet().PublicKey().String(),
				TokenPairID:   seededTokenPair2.tokenPairID,
				TokenAMint:    seededTokenPair2.tokenAPubkey,
				TokenBMint:    seededTokenPair2.tokenBPubkey,
			}
			err := newRepository.UpsertTokenSwaps(context.Background(), &tokenSwap1, &tokenSwap2)
			assert.NoError(t, err)

			var insertedTokenSwap model.TokenSwap
			err = db.Get(&insertedTokenSwap, "select token_swap.* from token_swap where id=$1", tokenSwap1.ID)
			assert.NoError(t, err)
			assert.Equal(t, tokenSwap1.Pubkey, insertedTokenSwap.Pubkey)

			err = db.Get(&insertedTokenSwap, "select token_swap.* from token_swap where id=$1", tokenSwap2.ID)
			assert.NoError(t, err)
			assert.Equal(t, tokenSwap2.Pubkey, insertedTokenSwap.Pubkey)
		})

		t.Run("should update tokenSwap if (swap, tokenAMint, tokenBMint) is violated", func(t *testing.T) {
			defer truncateDB(t, db)
			seededTokenSwap := seedTokenSwap(t, db, seedTokenSwapParams{})
			var seededTokenSwapModel model.TokenSwap
			err := db.Get(&seededTokenSwapModel, "select token_swap.* from token_swap where id=$1", seededTokenSwap.tokenSwapID)
			assert.NoError(t, err)

			tokenSwap := model.TokenSwap{
				ID:            seededTokenSwapModel.ID,
				Pubkey:        seededTokenSwapModel.Pubkey,
				Mint:          solana.NewWallet().PublicKey().String(),
				Authority:     solana.NewWallet().PublicKey().String(),
				FeeAccount:    solana.NewWallet().PublicKey().String(),
				TokenAAccount: solana.NewWallet().PublicKey().String(),
				TokenBAccount: solana.NewWallet().PublicKey().String(),
				TokenPairID:   seededTokenSwapModel.TokenPairID,
				TokenAMint:    seededTokenSwapModel.TokenAMint,
				TokenBMint:    seededTokenSwapModel.TokenBMint,
			}
			err = newRepository.UpsertTokenSwaps(context.Background(), &tokenSwap)
			assert.NoError(t, err)

			var insertedTokenSwap model.TokenSwap
			err = db.Get(&insertedTokenSwap, "select token_swap.* from token_swap where id=$1", tokenSwap.ID)
			assert.NoError(t, err)
			assert.Equal(t, tokenSwap.Pubkey, insertedTokenSwap.Pubkey)
			assert.Equal(t, tokenSwap.Mint, insertedTokenSwap.Mint)
			assert.Equal(t, tokenSwap.Authority, insertedTokenSwap.Authority)
			assert.Equal(t, tokenSwap.FeeAccount, insertedTokenSwap.FeeAccount)
			assert.Equal(t, tokenSwap.TokenAAccount, insertedTokenSwap.TokenAAccount)
			assert.Equal(t, tokenSwap.TokenBAccount, insertedTokenSwap.TokenBAccount)
			assert.Equal(t, tokenSwap.TokenAMint, insertedTokenSwap.TokenAMint)
			assert.Equal(t, tokenSwap.TokenBMint, insertedTokenSwap.TokenBMint)
		})
	})
}

//nolint:funlen
func TestUpsertTokenAccountBalances(t *testing.T) {
	test.InjectDependencies(func(
		repo *query.Query,
		db *sqlx.DB,
	) {
		newRepository := repository.NewRepository(repo, db)

		t.Run("should insert tokenAccountBalance", func(t *testing.T) {
			defer truncateDB(t, db)
			seededTokenPair := seedTokenPair(t, db, seedTokenPairParams{})

			tokenAccountBalance := model.TokenAccountBalance{
				Pubkey: solana.NewWallet().PublicKey().String(),
				Mint:   seededTokenPair.tokenAPubkey,
				Owner:  solana.NewWallet().PublicKey().String(),
				Amount: 10,
				State:  "initialized",
			}
			err := newRepository.UpsertTokenAccountBalances(context.Background(), &tokenAccountBalance)
			assert.NoError(t, err)

			var insertedTokenAccountBalance model.TokenAccountBalance
			err = db.Get(&insertedTokenAccountBalance, "select token_account_balance.* from token_account_balance where pubkey=$1", tokenAccountBalance.Pubkey)
			assert.NoError(t, err)
			assert.Equal(t, tokenAccountBalance.Pubkey, insertedTokenAccountBalance.Pubkey)
			assert.Equal(t, tokenAccountBalance.Mint, insertedTokenAccountBalance.Mint)
			assert.Equal(t, tokenAccountBalance.Owner, insertedTokenAccountBalance.Owner)
			assert.Equal(t, tokenAccountBalance.Amount, insertedTokenAccountBalance.Amount)
			assert.Equal(t, tokenAccountBalance.State, insertedTokenAccountBalance.State)
		})

		t.Run("should insert multiple tokenAccountBalances", func(t *testing.T) {
			defer truncateDB(t, db)
			seededTokenPair := seedTokenPair(t, db, seedTokenPairParams{})
			tokenAccountBalance1 := model.TokenAccountBalance{
				Pubkey: solana.NewWallet().PublicKey().String(),
				Mint:   seededTokenPair.tokenAPubkey,
				Owner:  solana.NewWallet().PublicKey().String(),
				Amount: 10,
				State:  "initialized",
			}
			tokenAccountBalance2 := model.TokenAccountBalance{
				Pubkey: solana.NewWallet().PublicKey().String(),
				Mint:   seededTokenPair.tokenBPubkey,
				Owner:  solana.NewWallet().PublicKey().String(),
				Amount: 6,
				State:  "initialized",
			}
			err := newRepository.UpsertTokenAccountBalances(context.Background(), &tokenAccountBalance1, &tokenAccountBalance2)
			assert.NoError(t, err)

			var insertedTokenAccountBalance model.TokenAccountBalance
			err = db.Get(&insertedTokenAccountBalance, "select token_account_balance.* from token_account_balance where pubkey=$1", tokenAccountBalance1.Pubkey)
			assert.NoError(t, err)
			assert.Equal(t, tokenAccountBalance1.Pubkey, insertedTokenAccountBalance.Pubkey)
			assert.Equal(t, tokenAccountBalance1.Mint, insertedTokenAccountBalance.Mint)
			assert.Equal(t, tokenAccountBalance1.Owner, insertedTokenAccountBalance.Owner)
			assert.Equal(t, tokenAccountBalance1.Amount, insertedTokenAccountBalance.Amount)
			assert.Equal(t, tokenAccountBalance1.State, insertedTokenAccountBalance.State)

			err = db.Get(&insertedTokenAccountBalance, "select token_account_balance.* from token_account_balance where pubkey=$1", tokenAccountBalance2.Pubkey)
			assert.NoError(t, err)
			assert.Equal(t, tokenAccountBalance2.Pubkey, insertedTokenAccountBalance.Pubkey)
			assert.Equal(t, tokenAccountBalance2.Mint, insertedTokenAccountBalance.Mint)
			assert.Equal(t, tokenAccountBalance2.Owner, insertedTokenAccountBalance.Owner)
			assert.Equal(t, tokenAccountBalance2.Amount, insertedTokenAccountBalance.Amount)
			assert.Equal(t, tokenAccountBalance2.State, insertedTokenAccountBalance.State)
		})

		t.Run("should update tokenAccountBalance", func(t *testing.T) {
			defer truncateDB(t, db)
			seededTokenAccountBalance := seedTokenAccountBalance(t, db, seedTokenAccountBalanceParams{})

			tokenAccountBalance := model.TokenAccountBalance{
				Pubkey: seededTokenAccountBalance.tokenAccountBalancePubkey,
				Mint:   seededTokenAccountBalance.tokenAPubkey,
				Owner:  solana.NewWallet().PublicKey().String(),
				Amount: 1000,
				State:  "initialized",
			}
			err := newRepository.UpsertTokenAccountBalances(context.Background(), &tokenAccountBalance)
			assert.NoError(t, err)

			var insertedTokenAccountBalance model.TokenAccountBalance
			err = db.Get(&insertedTokenAccountBalance, "select token_account_balance.* from token_account_balance where pubkey=$1", tokenAccountBalance.Pubkey)
			assert.NoError(t, err)
			assert.Equal(t, tokenAccountBalance.Pubkey, insertedTokenAccountBalance.Pubkey)
			assert.Equal(t, tokenAccountBalance.Mint, insertedTokenAccountBalance.Mint)
			assert.Equal(t, tokenAccountBalance.Owner, insertedTokenAccountBalance.Owner)
			assert.Equal(t, tokenAccountBalance.Amount, insertedTokenAccountBalance.Amount)
			assert.Equal(t, tokenAccountBalance.State, insertedTokenAccountBalance.State)
		})
	})
}
