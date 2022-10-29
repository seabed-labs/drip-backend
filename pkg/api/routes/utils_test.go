package controller

import (
	"reflect"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/dcaf-labs/drip/pkg/service/config"
	"github.com/dcaf-labs/drip/pkg/service/repository/model"
	"github.com/gagliardetto/solana-go"
	"github.com/google/uuid"
	"github.com/test-go/testify/assert"
)

func TestUtils(t *testing.T) {

	t.Run("getServerURL should return correct api URL", func(t *testing.T) {
		assert.True(t, strings.Contains(getServerURL(config.NilNetwork, config.StagingEnv, 0), "localhost"))
		assert.True(t, strings.Contains(getServerURL(config.LocalNetwork, config.StagingEnv, 0), "localhost"))
		assert.True(t, strings.Contains(getServerURL(config.DevnetNetwork, config.StagingEnv, 0), "devnet"))
		assert.True(t, strings.Contains(getServerURL(config.MainnetNetwork, config.StagingEnv, 0), "mainnet"))
	})

	t.Run("hasValue should return true", func(t *testing.T) {
		assert.True(t, hasValue([]string{"1", "2"}, "1"))
		assert.True(t, hasValue([]string{"1", "2"}, "2"))
	})

	t.Run("hasValue should return false", func(t *testing.T) {
		assert.False(t, hasValue([]string{"1", "2"}, "3"))
		assert.False(t, hasValue([]string{}, ""))
	})

	t.Run("vaultModelToAPI should return correct apiSpec model", func(t *testing.T) {
		dbModel := model.Vault{
			DcaActivationTimestamp: time.Unix(1667076855, 0),
			Pubkey:                 solana.NewWallet().PublicKey().String(),
			ProtoConfig:            solana.NewWallet().PublicKey().String(),
			TokenAAccount:          solana.NewWallet().PublicKey().String(),
			TokenBAccount:          solana.NewWallet().PublicKey().String(),
			TreasuryTokenBAccount:  solana.NewWallet().PublicKey().String(),
			LastDcaPeriod:          10,
			DripAmount:             99,
			Enabled:                true,
			TokenPairID:            uuid.New().String(),
			MaxSlippageBps:         50,
			TokenAMint:             solana.NewWallet().PublicKey().String(),
			TokenBMint:             solana.NewWallet().PublicKey().String(),
		}
		apiModel := vaultModelToAPI(&dbModel)
		assert.Equal(t, strconv.FormatInt(dbModel.DcaActivationTimestamp.Unix(), 10), apiModel.DcaActivationTimestamp)
		assert.Equal(t, dbModel.Pubkey, apiModel.Pubkey)
		assert.Equal(t, dbModel.ProtoConfig, apiModel.ProtoConfig)
		assert.Equal(t, dbModel.TokenAAccount, apiModel.TokenAAccount)
		assert.Equal(t, dbModel.TokenBAccount, apiModel.TokenBAccount)
		assert.Equal(t, dbModel.TreasuryTokenBAccount, apiModel.TreasuryTokenBAccount)
		assert.Equal(t, dbModel.TokenAMint, apiModel.TokenAMint)
		assert.Equal(t, dbModel.TokenBMint, apiModel.TokenBMint)
		assert.Equal(t, strconv.FormatUint(dbModel.LastDcaPeriod, 10), apiModel.LastDcaPeriod)
		assert.Equal(t, strconv.FormatUint(dbModel.DripAmount, 10), apiModel.DripAmount)
		assert.Equal(t, dbModel.Enabled, apiModel.Enabled)

		// if the line below needs to be updated, add the field assertion above
		assert.Equal(t, reflect.TypeOf(apiModel).NumField(), 11)
	})

}
