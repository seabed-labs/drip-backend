package integrationtest

import (
	"context"
	"reflect"
	"testing"

	"github.com/dcaf-labs/drip/pkg/unittest"
	"github.com/golang/mock/gomock"

	"github.com/dcaf-labs/drip/pkg/integrationtest"
	"github.com/dcaf-labs/drip/pkg/service/processor"
	"github.com/dcaf-labs/drip/pkg/service/repository"
	"github.com/test-go/testify/assert"
)

func TestHandler_UpsertProtoConfigByAddress(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration tests in short mode")
	}
	ctrl := gomock.NewController(t)
	t.Run("should upsert vault proto config", func(t *testing.T) {
		integrationtest.InjectDependencies(
			&integrationtest.TestOptions{
				FixturePath: "./fixtures/test1",
				AppConfig:   unittest.GetMockDevnetStagingConfig(ctrl),
			},
			func(
				processor processor.Processor,
				repo repository.Repository,
			) {
				protoConfigAddress := "Et3bqQq32LPkrndf8gU9gRqfL4S13ubdUuiqBE1jjrgr"
				// protoConfig: https://explorer.solana.com/address/Et3bqQq32LPkrndf8gU9gRqfL4S13ubdUuiqBE1jjrgr?cluster=devnet
				assert.NoError(t, processor.UpsertProtoConfigByAddress(context.Background(), protoConfigAddress))

				protoConfig, err := repo.GetProtoConfigByAddress(context.Background(), protoConfigAddress)
				assert.NoError(t, err)
				assert.NotNil(t, protoConfig)

				assert.Equal(t, "Et3bqQq32LPkrndf8gU9gRqfL4S13ubdUuiqBE1jjrgr", protoConfig.Pubkey)
				assert.Equal(t, "3CTkqdcjzn1ptnNYUcFq4Sk1Smy91kz5p9JhgJBHGe3e", protoConfig.Admin)
				assert.Equal(t, uint64(60), protoConfig.Granularity)
				assert.Equal(t, uint16(50), protoConfig.TokenADripTriggerSpread)
				assert.Equal(t, uint16(25), protoConfig.TokenBWithdrawalSpread)
				assert.Equal(t, uint16(25), protoConfig.TokenBReferralSpread)
				// if the line below needs to be updated, add the field assertion above
				assert.Equal(t, reflect.TypeOf(*protoConfig).NumField(), 6)
			})
	})
}

func TestHandler_UpsertVaultByAddress(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration tests in short mode")
	}
	ctrl := gomock.NewController(t)
	t.Run("should upsert vault and all related accounts", func(t *testing.T) {
		integrationtest.InjectDependencies(
			&integrationtest.TestOptions{
				FixturePath: "./fixtures/test2",
				AppConfig:   unittest.GetMockDevnetStagingConfig(ctrl),
			},
			func(
				processor processor.Processor,
				repo repository.Repository,
			) {
				vaultAddress := "BwJj7DYyMR1xMnWK1PGKPLi5u2ZP5EDBFBkPAAv4UDP8"
				// protoConfig: https://explorer.solana.com/address/BwJj7DYyMR1xMnWK1PGKPLi5u2ZP5EDBFBkPAAv4UDP8?cluster=devnet
				assert.NoError(t, processor.UpsertVaultByAddress(context.Background(), vaultAddress))

				// vault should exist
				vault, err := repo.AdminGetVaultByAddress(context.Background(), "BwJj7DYyMR1xMnWK1PGKPLi5u2ZP5EDBFBkPAAv4UDP8")
				assert.NoError(t, err)
				assert.NotNil(t, vault)

				// by extension proto config should exist
				protoConfig, err := repo.GetProtoConfigByAddress(context.Background(), "Et3bqQq32LPkrndf8gU9gRqfL4S13ubdUuiqBE1jjrgr")
				assert.NoError(t, err)
				assert.NotNil(t, protoConfig)

				// by extension vault token accounts should exist
				tokenAccounts, err := repo.GetTokenAccountsByAddresses(context.Background(), vault.TokenAAccount, vault.TokenBAccount, vault.TreasuryTokenBAccount)
				assert.NoError(t, err)
				assert.Len(t, tokenAccounts, 3)

				// by extension vault tokens should exist
				mints, err := repo.GetTokensByAddresses(context.Background(), vault.TokenAMint, vault.TokenBMint)
				assert.NoError(t, err)
				assert.Len(t, mints, 2)

				// by extension token pair should exist
				tokenPair, err := repo.GetTokenPair(context.Background(), vault.TokenAMint, vault.TokenBMint)
				assert.NoError(t, err)
				assert.NotNil(t, tokenPair)
				assert.Equal(t, "H9gBUJs5Kc5zyiKRTzZcYom4Hpj9VPHLy4VzExTVPgxa", tokenPair.TokenA)
				assert.Equal(t, "7ihthG4cFydyDnuA3zmJrX13ePGpLcANf3tHLmKLPN7M", tokenPair.TokenB)
				// if the line below needs to be updated, add the field assertion above
				assert.Equal(t, reflect.TypeOf(*tokenPair).NumField(), 3)

				assert.Equal(t, "BwJj7DYyMR1xMnWK1PGKPLi5u2ZP5EDBFBkPAAv4UDP8", vault.Pubkey)
				assert.Equal(t, "Et3bqQq32LPkrndf8gU9gRqfL4S13ubdUuiqBE1jjrgr", vault.ProtoConfig)
				assert.Equal(t, "H9gBUJs5Kc5zyiKRTzZcYom4Hpj9VPHLy4VzExTVPgxa", vault.TokenAMint)
				assert.Equal(t, "7ihthG4cFydyDnuA3zmJrX13ePGpLcANf3tHLmKLPN7M", vault.TokenBMint)
				assert.Equal(t, "HPTSZFzxKUsJKPrxSyYh6TuGQV7qfBqpW6c4A8DERxFr", vault.TokenAAccount)
				assert.Equal(t, "9Swein4rvYYN1MJyFBjfTdvx8Lh3wzabwtWKddPWdhrP", vault.TokenBAccount)
				assert.Equal(t, "6smBUZt2e7Dz9o6hWoXmY78n1rpXrgpNwNFtocZMU5QN", vault.TreasuryTokenBAccount)
				assert.Equal(t, uint64(1847), vault.LastDcaPeriod)
				assert.Equal(t, uint64(75250000), vault.DripAmount)
				assert.Equal(t, int64(1664673300), vault.DcaActivationTimestamp.Unix())
				assert.Equal(t, int32(1000), vault.MaxSlippageBps)
				assert.Equal(t, false, vault.Enabled)
				assert.Equal(t, tokenPair.ID, vault.TokenPairID)
				// if the line below needs to be updated, add the field assertion above
				assert.Equal(t, reflect.TypeOf(*vault).NumField(), 13)
			})
	})
}
