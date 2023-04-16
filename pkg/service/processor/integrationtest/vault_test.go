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
		integrationtest.TestWithInjectedDependencies(
			&integrationtest.TestOptions{
				FixturePath: "./fixtures/test1",
				AppConfig:   unittest.GetMockMainnetProductionConfig(ctrl),
			},
			func(
				processor processor.Processor,
				repo repository.Repository,
			) {
				protoConfigAddress := "BxEimfWJshFBXSuqttnjdabud98h3BTnmvbzCQkZXWZK"
				// protoConfig: https://explorer.solana.com/address/BxEimfWJshFBXSuqttnjdabud98h3BTnmvbzCQkZXWZK
				assert.NoError(t, processor.UpsertProtoConfigByAddress(context.Background(), protoConfigAddress))

				protoConfig, err := repo.GetProtoConfigByAddress(context.Background(), protoConfigAddress)
				assert.NoError(t, err)
				assert.NotNil(t, protoConfig)

				assert.Equal(t, "BxEimfWJshFBXSuqttnjdabud98h3BTnmvbzCQkZXWZK", protoConfig.Pubkey)
				assert.Equal(t, "JC5NuYudj4Dd8vFDS8re1spg3PZEZ3PZyBRef852vz5R", protoConfig.Admin)
				assert.Equal(t, uint64(3600), protoConfig.Granularity)
				assert.Equal(t, uint16(35), protoConfig.TokenADripTriggerSpread)
				assert.Equal(t, uint16(5), protoConfig.TokenBWithdrawalSpread)
				assert.Equal(t, uint16(5), protoConfig.TokenBReferralSpread)
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
		integrationtest.TestWithInjectedDependencies(
			&integrationtest.TestOptions{
				FixturePath: "./fixtures/test2",
				AppConfig:   unittest.GetMockMainnetProductionConfig(ctrl),
			},
			func(
				processor processor.Processor,
				repo repository.Repository,
			) {
				vaultAddress := "HJcGW3iQGvLpnaJ7LATnrBdw66MTC1hzmmQMKN8pgVhH"
				// vault:https://explorer.solana.com/address/HJcGW3iQGvLpnaJ7LATnrBdw66MTC1hzmmQMKN8pgVhH/anchor-account
				assert.NoError(t, processor.UpsertVaultByAddress(context.Background(), vaultAddress))

				// vault should exist
				vault, err := repo.AdminGetVaultByAddress(context.Background(), vaultAddress)
				assert.NoError(t, err)
				assert.NotNil(t, vault)

				// by extension proto config should exist
				protoConfig, err := repo.GetProtoConfigByAddress(context.Background(), "BxEimfWJshFBXSuqttnjdabud98h3BTnmvbzCQkZXWZK")
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
				assert.Equal(t, "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v", tokenPair.TokenA)
				assert.Equal(t, "So11111111111111111111111111111111111111112", tokenPair.TokenB)
				// if the line below needs to be updated, add the field assertion above
				assert.Equal(t, reflect.TypeOf(*tokenPair).NumField(), 3)

				assert.Equal(t, vaultAddress, vault.Pubkey)
				assert.Equal(t, "BxEimfWJshFBXSuqttnjdabud98h3BTnmvbzCQkZXWZK", vault.ProtoConfig)
				assert.Equal(t, "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v", vault.TokenAMint)
				assert.Equal(t, "So11111111111111111111111111111111111111112", vault.TokenBMint)
				assert.Equal(t, "2JxQvnZcXNLugZTdLxCohWvRa24juDyfPY75w3HxMxUv", vault.TokenAAccount)
				assert.Equal(t, "2AvGt9sXPxdVfJWq9mYn2umjnKLfKTUHEfchPmwADEn2", vault.TokenBAccount)
				assert.Equal(t, "6i5YrPuJWReB9XCxaoB4ghga2PShFFWRCYYXPXhu4j1W", vault.TreasuryTokenBAccount)
				assert.Equal(t, uint64(0xb55), vault.LastDcaPeriod)
				assert.Equal(t, uint64(0x1bebc), vault.DripAmount)
				assert.Equal(t, int64(1681664400), vault.DcaActivationTimestamp.Unix())
				assert.Equal(t, int32(500), vault.MaxSlippageBps)
				assert.Equal(t, false, vault.Enabled)
				assert.Equal(t, tokenPair.ID, vault.TokenPairID)
				// if the line below needs to be updated, add the field assertion above
				assert.Equal(t, 13, reflect.TypeOf(*vault).NumField())
			})
	})
}
