package testintegration

import (
	"context"
	"os"
	"testing"

	"github.com/dcaf-protocol/drip/pkg/repository"
	"github.com/dcaf-protocol/drip/pkg/repository/model"
	"github.com/dcaf-protocol/drip/pkg/repository/query"
	"github.com/dcaf-protocol/drip/pkg/test"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/test-go/testify/assert"
)

//nolint:funlen
func TestUpsertProtoConfigs(t *testing.T) {

	err := os.Setenv("IS_TEST_DB", "true")
	assert.NoError(t, err)

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

		t.Run("should insert proto config", func(t *testing.T) {
			defer cleanup()

			pubkey := "123"

			protoConfig := model.ProtoConfig{
				Pubkey:               pubkey,
				Granularity:          1,
				TriggerDcaSpread:     5,
				BaseWithdrawalSpread: 10,
			}
			err := newRepository.UpsertProtoConfigs(context.Background(), &protoConfig)
			assert.NoError(t, err)

			var insertedConfig model.ProtoConfig
			err = db.Get(&insertedConfig, "select proto_config.* from proto_config where pubkey=$1", pubkey)
			assert.NoError(t, err)
			assert.Equal(t, protoConfig.Pubkey, insertedConfig.Pubkey)
		})

		t.Run("should update proto config", func(t *testing.T) {
			defer cleanup()

			pubkey := uuid.New().String()[0:4]
			_, err := db.Exec("insert into proto_config(pubkey, granularity, trigger_dca_spread, base_withdrawal_spread) values($1, $2, $3, $4)",
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

		//UpsertTokens(context.Context, ...*model2.Token) error
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
	})
}
