package testintegration

import (
	"context"
	"testing"
	"time"

	"github.com/dcaf-protocol/drip/pkg/repository"
	"github.com/dcaf-protocol/drip/pkg/repository/model"
	"github.com/dcaf-protocol/drip/pkg/repository/query"
	"github.com/dcaf-protocol/drip/pkg/test"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/test-go/testify/assert"
)

// TODO(Mocha): these tests all take a long time because each test fn creates a new DB and runs fresh migrations
// the db setup and migrations can be done once per file opposed to once per fn

//UpsertVaultPeriods(context.Context, ...*model2.VaultPeriod) error
//UpsertPositions(context.Context, ...*model2.Position) error
//UpsertTokenSwaps(context.Context, ...*model2.TokenSwap) error
//UpsertTokenAccountBalances(context.Context, ...*model2.TokenAccountBalance) error

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
		cleanup := func() {
			_, err := db.Exec("DELETE from proto_config")
			assert.NoError(t, err)
		}
		cleanup()

		t.Run("should insert protoConfig", func(t *testing.T) {
			defer cleanup()

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
			defer cleanup()

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

		t.Run("should update proto config", func(t *testing.T) {
			defer cleanup()

			pubkey := uuid.New().String()[0:4]
			_, err := db.Exec(
				`insert into 
    						proto_config(pubkey, granularity, trigger_dca_spread, base_withdrawal_spread) 
							values($1, $2, $3, $4)`,
				pubkey, 1, 2, 3)
			assert.NoError(t, err)

			protoConfig := model.ProtoConfig{
				Pubkey:               pubkey,
				Granularity:          2,
				TriggerDcaSpread:     4,
				BaseWithdrawalSpread: 6,
			}
			err = newRepository.UpsertProtoConfigs(context.Background(), &protoConfig)
			assert.NoError(t, err)

			var updatedProtoConfig model.ProtoConfig
			err = db.Get(&updatedProtoConfig, "select proto_config.* from proto_config where pubkey=$1", pubkey)
			assert.NoError(t, err)
			assert.Equal(t, protoConfig.Pubkey, updatedProtoConfig.Pubkey)
			assert.Equal(t, updatedProtoConfig.Granularity, uint64(2))
			assert.Equal(t, updatedProtoConfig.TriggerDcaSpread, uint16(4))
			assert.Equal(t, updatedProtoConfig.BaseWithdrawalSpread, uint16(6))
		})

		t.Run("should update many protoConfigs", func(t *testing.T) {
			defer cleanup()

			pubkey1 := uuid.New().String()[0:4]
			_, err := db.Exec(
				`insert into 
    						proto_config(pubkey, granularity, trigger_dca_spread, base_withdrawal_spread) 
							values($1, $2, $3, $4)`,
				pubkey1, 1, 2, 3)
			assert.NoError(t, err)

			pubkey2 := uuid.New().String()[0:4]
			_, err = db.Exec(
				`insert into 
    						proto_config(pubkey, granularity, trigger_dca_spread, base_withdrawal_spread) 
							values($1, $2, $3, $4)`,
				pubkey2, 4, 5, 6)
			assert.NoError(t, err)

			protoConfig1 := model.ProtoConfig{
				Pubkey:               pubkey1,
				Granularity:          7,
				TriggerDcaSpread:     8,
				BaseWithdrawalSpread: 9,
			}
			protoConfig2 := model.ProtoConfig{
				Pubkey:               pubkey2,
				Granularity:          10,
				TriggerDcaSpread:     11,
				BaseWithdrawalSpread: 12,
			}
			err = newRepository.UpsertProtoConfigs(context.Background(), &protoConfig1, &protoConfig2)
			assert.NoError(t, err)

			var updatedProtoConfig model.ProtoConfig
			err = db.Get(&updatedProtoConfig, "select proto_config.* from proto_config where pubkey=$1", pubkey1)
			assert.NoError(t, err)
			assert.Equal(t, protoConfig1.Pubkey, updatedProtoConfig.Pubkey)
			assert.Equal(t, updatedProtoConfig.Granularity, uint64(7))
			assert.Equal(t, updatedProtoConfig.TriggerDcaSpread, uint16(8))
			assert.Equal(t, updatedProtoConfig.BaseWithdrawalSpread, uint16(9))

			err = db.Get(&updatedProtoConfig, "select proto_config.* from proto_config where pubkey=$1", pubkey2)
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
		newRepository := repository.NewRepository(repo, db)
		cleanup := func() {
			_, err := db.Exec("DELETE from token")
			assert.NoError(t, err)
		}
		cleanup()

		t.Run("should insert token", func(t *testing.T) {
			defer cleanup()

			pubkey := uuid.New().String()[0:4]
			symbol := "btc"
			token := model.Token{
				Pubkey:   pubkey,
				Symbol:   &symbol,
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
			defer cleanup()

			pubkey1 := uuid.New().String()[0:4]
			symbol1 := "btc"
			pubkey2 := uuid.New().String()[0:4]
			symbol2 := "eth"
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
			defer cleanup()

			pubkey := uuid.New().String()[0:4]
			_, err := db.Exec(
				`insert into 
    						token(pubkey, symbol, decimals, icon_url) 
							values($1, $2, $3, $4)`,
				pubkey, "btc", 2, nil)
			assert.NoError(t, err)

			symbol := "eth"
			token := model.Token{
				Pubkey:   pubkey,
				Symbol:   &symbol,
				Decimals: 0,
				IconURL:  nil,
			}
			err = newRepository.UpsertTokens(context.Background(), &token)
			assert.NoError(t, err)

			var updatedToken model.Token
			err = db.Get(&updatedToken, "select token.* from token where pubkey=$1", pubkey)
			assert.NoError(t, err)
			assert.Equal(t, token.Pubkey, updatedToken.Pubkey)
			assert.Equal(t, *updatedToken.Symbol, symbol)
		})

		t.Run("should update many tokens", func(t *testing.T) {
			defer cleanup()

			pubkey1 := uuid.New().String()[0:4]
			_, err := db.Exec(
				`insert into 
    						token(pubkey, symbol, decimals, icon_url) 
							values($1, $2, $3, $4)`,
				pubkey1, "btc", 2, nil)
			assert.NoError(t, err)

			pubkey2 := uuid.New().String()[0:4]
			_, err = db.Exec(
				`insert into 
    						token(pubkey, symbol, decimals, icon_url) 
							values($1, $2, $3, $4)`,
				pubkey2, "ltc", 2, nil)
			assert.NoError(t, err)

			symbol1 := "eth"
			token1 := model.Token{
				Pubkey:   pubkey1,
				Symbol:   &symbol1,
				Decimals: 0,
				IconURL:  nil,
			}
			symbol2 := "sol"
			token2 := model.Token{
				Pubkey:   pubkey2,
				Symbol:   &symbol2,
				Decimals: 0,
				IconURL:  nil,
			}
			err = newRepository.UpsertTokens(context.Background(), &token1, &token2)
			assert.NoError(t, err)

			var updatedToken model.Token
			err = db.Get(&updatedToken, "select token.* from token where pubkey=$1", pubkey1)
			assert.NoError(t, err)
			assert.Equal(t, token1.Pubkey, updatedToken.Pubkey)
			assert.Equal(t, *updatedToken.Symbol, symbol1)

			err = db.Get(&updatedToken, "select token.* from token where pubkey=$1", pubkey2)
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
		cleanup := func() {
			_, err := db.Exec("DELETE from token_pair")
			assert.NoError(t, err)
			_, err = db.Exec("DELETE from token")
			assert.NoError(t, err)
		}
		cleanup()

		t.Run("should fail to insert tokenPair if tokenA doesn't exit", func(t *testing.T) {
			defer cleanup()
			btcPubkey := uuid.New().String()
			_, err := db.Exec(
				`insert into 
    						token(pubkey, symbol, decimals, icon_url) 
							values
							    ($1, $2, $3, $4)`,
				btcPubkey, "btc", 2, nil,
			)
			assert.NoError(t, err)

			tokenPair := model.TokenPair{
				ID:     uuid.New().String(),
				TokenA: uuid.New().String(),
				TokenB: btcPubkey,
			}

			err = newRepository.InsertTokenPairs(context.Background(), &tokenPair)
			assert.Error(t, err)
		})

		t.Run("should fail to insert tokenPair if tokenB doesn't exit", func(t *testing.T) {
			defer cleanup()
			btcPubkey := uuid.New().String()
			_, err := db.Exec(
				`insert into 
    						token(pubkey, symbol, decimals, icon_url) 
							values
							    ($1, $2, $3, $4)`,
				btcPubkey, "btc", 2, nil,
			)
			assert.NoError(t, err)

			tokenPair := model.TokenPair{
				ID:     uuid.New().String(),
				TokenA: btcPubkey,
				TokenB: uuid.New().String(),
			}

			err = newRepository.InsertTokenPairs(context.Background(), &tokenPair)
			assert.Error(t, err)
		})

		t.Run("should insert tokenPair", func(t *testing.T) {
			defer cleanup()
			btcPubkey := uuid.New().String()
			ethPubkey := uuid.New().String()
			_, err := db.Exec(
				`insert into 
    						token(pubkey, symbol, decimals, icon_url) 
							values
							    ($1, $2, $3, $4),
							    ($5, $6, $7, $8)`,
				btcPubkey, "btc", 2, nil,
				ethPubkey, "eth", 2, nil,
			)
			assert.NoError(t, err)

			tokenPair := model.TokenPair{
				ID:     uuid.New().String(),
				TokenA: btcPubkey,
				TokenB: ethPubkey,
			}

			err = newRepository.InsertTokenPairs(context.Background(), &tokenPair)
			assert.NoError(t, err)

			var insertedTokenPair model.TokenPair
			err = db.Get(&insertedTokenPair, "select token_pair.* from token_pair where id=$1", tokenPair.ID)
			assert.NoError(t, err)
			assert.Equal(t, tokenPair.ID, insertedTokenPair.ID)
			assert.Equal(t, tokenPair.TokenA, insertedTokenPair.TokenA)
			assert.Equal(t, tokenPair.TokenB, insertedTokenPair.TokenB)
		})

		t.Run("should insert many tokenPairs", func(t *testing.T) {
			defer cleanup()
			btcPubkey := uuid.New().String()
			ethPubkey := uuid.New().String()
			_, err := db.Exec(
				`insert into 
    						token(pubkey, symbol, decimals, icon_url) 
							values
							    ($1, $2, $3, $4),
							    ($5, $6, $7, $8)`,
				btcPubkey, "btc", 2, nil,
				ethPubkey, "eth", 2, nil,
			)
			assert.NoError(t, err)

			tokenPair1 := model.TokenPair{
				ID:     uuid.New().String(),
				TokenA: btcPubkey,
				TokenB: ethPubkey,
			}

			tokenPair2 := model.TokenPair{
				ID:     uuid.New().String(),
				TokenA: ethPubkey,
				TokenB: btcPubkey,
			}

			err = newRepository.InsertTokenPairs(context.Background(), &tokenPair1, &tokenPair2)
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
			defer cleanup()
			btcPubkey := uuid.New().String()
			ethPubkey := uuid.New().String()
			_, err := db.Exec(
				`insert into 
    						token(pubkey, symbol, decimals, icon_url) 
							values
							    ($1, $2, $3, $4),
							    ($5, $6, $7, $8)`,
				btcPubkey, "btc", 2, nil,
				ethPubkey, "eth", 2, nil,
			)
			assert.NoError(t, err)
			originalTokenPair := model.TokenPair{
				ID:     uuid.New().String(),
				TokenA: btcPubkey,
				TokenB: ethPubkey,
			}

			err = newRepository.InsertTokenPairs(context.Background(), &originalTokenPair)
			assert.NoError(t, err)

			tokenPair := model.TokenPair{
				ID:     uuid.New().String(),
				TokenA: btcPubkey,
				TokenB: ethPubkey,
			}

			err = newRepository.InsertTokenPairs(context.Background(), &tokenPair)
			assert.NoError(t, err)

			var insertedTokenPair model.TokenPair
			err = db.Get(&insertedTokenPair, "select token_pair.* from token_pair where id=$1", originalTokenPair.ID)
			assert.NoError(t, err)
			assert.Equal(t, originalTokenPair.ID, insertedTokenPair.ID)
			assert.Equal(t, originalTokenPair.TokenA, insertedTokenPair.TokenA)
			assert.Equal(t, originalTokenPair.TokenB, insertedTokenPair.TokenB)
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
		cleanup := func() {
			_, err := db.Exec("truncate proto_config cascade")
			assert.NoError(t, err)
			_, err = db.Exec("truncate token_pair cascade")
			assert.NoError(t, err)
			_, err = db.Exec(" truncate token cascade")
			assert.NoError(t, err)
			_, err = db.Exec(" truncate vault cascade")
			assert.NoError(t, err)
		}
		cleanup()

		t.Run("should fail to insert vault when protoConfig is missing", func(t *testing.T) {
			defer cleanup()

			btcPubkey := uuid.New().String()
			ethPubkey := uuid.New().String()
			_, err := db.Exec(
				`insert into 
    						token(pubkey, symbol, decimals, icon_url) 
							values
							    ($1, $2, $3, $4),
							    ($5, $6, $7, $8)`,
				btcPubkey, "btc", 2, nil,
				ethPubkey, "eth", 2, nil,
			)
			assert.NoError(t, err)

			tokenPairID := uuid.New()
			_, err = db.Exec(
				`insert into 
    						token_pair(id, token_a, token_b) 
							values
							    ($1, $2, $3)`,
				tokenPairID.String(), btcPubkey, ethPubkey,
			)
			assert.NoError(t, err)

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
				TokenPairID:            tokenPairID.String(),
			}
			err = newRepository.UpsertVaults(context.Background(), &vault)
			assert.Error(t, err)
		})

		t.Run("should fail to insert vault when token pair is missing", func(t *testing.T) {
			defer cleanup()

			pubkey := uuid.New().String()[0:4]
			_, err := db.Exec(
				`insert into 
    						proto_config(pubkey, granularity, trigger_dca_spread, base_withdrawal_spread) 
							values($1, $2, $3, $4)`,
				pubkey, 1, 2, 3)
			assert.NoError(t, err)

			btcPubkey := uuid.New().String()
			ethPubkey := uuid.New().String()
			_, err = db.Exec(
				`insert into 
    						token(pubkey, symbol, decimals, icon_url) 
							values
							    ($1, $2, $3, $4),
							    ($5, $6, $7, $8)`,
				btcPubkey, "btc", 2, nil,
				ethPubkey, "eth", 2, nil,
			)
			assert.NoError(t, err)

			vault := model.Vault{
				Pubkey:                 uuid.New().String(),
				ProtoConfig:            pubkey,
				TokenAAccount:          uuid.New().String(),
				TokenBAccount:          uuid.New().String(),
				TreasuryTokenBAccount:  uuid.New().String(),
				LastDcaPeriod:          0,
				DripAmount:             0,
				DcaActivationTimestamp: time.Time{},
				Enabled:                false,
				TokenPairID:            uuid.New().String(),
			}
			err = newRepository.UpsertVaults(context.Background(), &vault)
			assert.Error(t, err)
		})

		t.Run("should insert vault", func(t *testing.T) {
			defer cleanup()

			pubkey := uuid.New().String()[0:4]
			_, err := db.Exec(
				`insert into 
    						proto_config(pubkey, granularity, trigger_dca_spread, base_withdrawal_spread) 
							values($1, $2, $3, $4)`,
				pubkey, 1, 2, 3)
			assert.NoError(t, err)

			btcPubkey := uuid.New().String()
			ethPubkey := uuid.New().String()
			_, err = db.Exec(
				`insert into 
    						token(pubkey, symbol, decimals, icon_url) 
							values
							    ($1, $2, $3, $4),
							    ($5, $6, $7, $8)`,
				btcPubkey, "btc", 2, nil,
				ethPubkey, "eth", 2, nil,
			)
			assert.NoError(t, err)

			tokenPairID := uuid.New()
			_, err = db.Exec(
				`insert into 
    						token_pair(id, token_a, token_b) 
							values
							    ($1, $2, $3)`,
				tokenPairID.String(), btcPubkey, ethPubkey,
			)
			assert.NoError(t, err)

			vault := model.Vault{
				Pubkey:                 uuid.New().String(),
				ProtoConfig:            pubkey,
				TokenAAccount:          uuid.New().String(),
				TokenBAccount:          uuid.New().String(),
				TreasuryTokenBAccount:  uuid.New().String(),
				LastDcaPeriod:          0,
				DripAmount:             0,
				DcaActivationTimestamp: time.Time{},
				Enabled:                false,
				TokenPairID:            tokenPairID.String(),
			}
			err = newRepository.UpsertVaults(context.Background(), &vault)
			assert.NoError(t, err)

			var insertedVault model.Vault
			err = db.Get(&insertedVault, "select vault.* from vault where pubkey=$1", vault.Pubkey)
			assert.NoError(t, err)
			assert.Equal(t, vault.Pubkey, insertedVault.Pubkey)
			assert.Equal(t, insertedVault.TokenPairID, tokenPairID.String())
		})

		t.Run("should insert many vaults", func(t *testing.T) {
			defer cleanup()

			pubkey := uuid.New().String()[0:4]
			_, err := db.Exec(
				`insert into 
    						proto_config(pubkey, granularity, trigger_dca_spread, base_withdrawal_spread) 
							values($1, $2, $3, $4)`,
				pubkey, 1, 2, 3)
			assert.NoError(t, err)

			btcPubkey := uuid.New().String()
			ethPubkey := uuid.New().String()
			_, err = db.Exec(
				`insert into 
    						token(pubkey, symbol, decimals, icon_url) 
							values
							    ($1, $2, $3, $4),
							    ($5, $6, $7, $8)`,
				btcPubkey, "btc", 2, nil,
				ethPubkey, "eth", 2, nil,
			)
			assert.NoError(t, err)

			tokenPairID := uuid.New()
			_, err = db.Exec(
				`insert into 
    						token_pair(id, token_a, token_b) 
							values
							    ($1, $2, $3)`,
				tokenPairID.String(), btcPubkey, ethPubkey,
			)
			assert.NoError(t, err)

			vault1 := model.Vault{
				Pubkey:                 uuid.New().String(),
				ProtoConfig:            pubkey,
				TokenAAccount:          uuid.New().String(),
				TokenBAccount:          uuid.New().String(),
				TreasuryTokenBAccount:  uuid.New().String(),
				LastDcaPeriod:          0,
				DripAmount:             0,
				DcaActivationTimestamp: time.Time{},
				Enabled:                false,
				TokenPairID:            tokenPairID.String(),
			}
			vault2 := model.Vault{
				Pubkey:                 uuid.New().String(),
				ProtoConfig:            pubkey,
				TokenAAccount:          uuid.New().String(),
				TokenBAccount:          uuid.New().String(),
				TreasuryTokenBAccount:  uuid.New().String(),
				LastDcaPeriod:          0,
				DripAmount:             0,
				DcaActivationTimestamp: time.Time{},
				Enabled:                false,
				TokenPairID:            tokenPairID.String(),
			}
			err = newRepository.UpsertVaults(context.Background(), &vault1, &vault2)
			assert.NoError(t, err)

			var insertedVault model.Vault
			err = db.Get(&insertedVault, "select vault.* from vault where pubkey=$1", vault1.Pubkey)
			assert.NoError(t, err)
			assert.Equal(t, vault1.Pubkey, insertedVault.Pubkey)
			assert.Equal(t, insertedVault.TokenPairID, tokenPairID.String())

			err = db.Get(&insertedVault, "select vault.* from vault where pubkey=$1", vault2.Pubkey)
			assert.NoError(t, err)
			assert.Equal(t, vault2.Pubkey, insertedVault.Pubkey)
			assert.Equal(t, insertedVault.TokenPairID, tokenPairID.String())
		})

		t.Run("should update vault", func(t *testing.T) {
			defer cleanup()

			protoConfigPubkey := uuid.New().String()[0:4]
			_, err := db.Exec(
				`insert into 
    						proto_config(pubkey, granularity, trigger_dca_spread, base_withdrawal_spread) 
							values($1, $2, $3, $4)`,
				protoConfigPubkey, 1, 2, 3)
			assert.NoError(t, err)

			btcPubkey := uuid.New().String()
			ethPubkey := uuid.New().String()
			_, err = db.Exec(
				`insert into 
    						token(pubkey, symbol, decimals, icon_url) 
							values
							    ($1, $2, $3, $4),
							    ($5, $6, $7, $8)`,
				btcPubkey, "btc", 2, nil,
				ethPubkey, "eth", 2, nil,
			)
			assert.NoError(t, err)

			tokenPairID := uuid.New()
			_, err = db.Exec(
				`insert into 
    						token_pair(id, token_a, token_b) 
							values
							    ($1, $2, $3)`,
				tokenPairID.String(), btcPubkey, ethPubkey,
			)
			assert.NoError(t, err)

			vaultPubkey := uuid.New().String()
			tokenAccountPubkey := uuid.New().String()
			_, err = db.Exec(
				`insert into 
    						vault(pubkey, proto_config, token_a_account, token_b_account, treasury_token_b_account, last_dca_period, drip_amount, dca_activation_timestamp, enabled, token_pair_id) 
							values
							    ($1, $2, $3, $4,$5, $6,$7,$8,$9,$10)`,
				vaultPubkey, protoConfigPubkey, tokenAccountPubkey, tokenAccountPubkey, tokenAccountPubkey, 0, 0, time.Time{}, false, tokenPairID.String(),
			)
			assert.NoError(t, err)
			vault := model.Vault{
				Pubkey:                 vaultPubkey,
				ProtoConfig:            protoConfigPubkey,
				TokenAAccount:          tokenAccountPubkey,
				TokenBAccount:          tokenAccountPubkey,
				TreasuryTokenBAccount:  tokenAccountPubkey,
				LastDcaPeriod:          1,
				DripAmount:             100,
				DcaActivationTimestamp: time.Now(),
				Enabled:                true,
				TokenPairID:            tokenPairID.String(),
			}
			err = newRepository.UpsertVaults(context.Background(), &vault)
			assert.NoError(t, err)

			var insertedVault model.Vault
			err = db.Get(&insertedVault, "select vault.* from vault where pubkey=$1", vaultPubkey)
			assert.NoError(t, err)
			assert.Equal(t, vault.Pubkey, insertedVault.Pubkey)
			assert.Equal(t, insertedVault.LastDcaPeriod, uint64(1))
			assert.Equal(t, insertedVault.DripAmount, uint64(100))
			assert.NotEqual(t, insertedVault.LastDcaPeriod, time.Time{})
			assert.Equal(t, insertedVault.Enabled, true)
		})

	})
}
