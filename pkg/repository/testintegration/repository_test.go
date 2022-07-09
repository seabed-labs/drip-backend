package testintegration

import (
	"context"
	"testing"

	"github.com/dcaf-protocol/drip/pkg/repository"
	"github.com/dcaf-protocol/drip/pkg/repository/model"
	"github.com/dcaf-protocol/drip/pkg/repository/query"
	"github.com/dcaf-protocol/drip/pkg/test"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/test-go/testify/assert"
)

//UpsertTokenPairs(context.Context, ...*model2.TokenPair) error
//UpsertVaults(context.Context, ...*model2.Vault) error
//UpsertVaultPeriods(context.Context, ...*model2.VaultPeriod) error
//UpsertPositions(context.Context, ...*model2.Position) error
//UpsertTokenSwaps(context.Context, ...*model2.TokenSwap) error
//UpsertTokenAccountBalances(context.Context, ...*model2.TokenAccountBalance) error
//
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
//
//InternalGetVaultByAddress(ctx context.Context, pubkey string) (*model2.Vault, error)
//EnableVault(ctx context.Context, pubkey string) (*model2.Vault, error)
//originalBalance, err := solClient.GetBalance(context.Background(), walletProvider.Wallet.PublicKey(), rpc.CommitmentConfirmed)
//assert.NoError(t, err)
//balance, err := walletPkg.InitTestWallet(solClient, walletProvider)
//assert.NoError(t, err)
//assert.NotZero(t, balance)
//assert.NotEqual(t, originalBalance, balance)

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
