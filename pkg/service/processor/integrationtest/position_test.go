package integrationtest

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/dcaf-labs/drip/pkg/integrationtest"
	"github.com/dcaf-labs/drip/pkg/service/processor"
	"github.com/dcaf-labs/drip/pkg/service/repository"
	"github.com/dcaf-labs/drip/pkg/unittest"
	"github.com/golang/mock/gomock"
	"github.com/test-go/testify/assert"
)

func TestHandler_UpsertPositionByAddress(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration tests in short mode")
	}
	ctrl := gomock.NewController(t)
	t.Run("should upsert position and all related accounts", func(t *testing.T) {
		integrationtest.InjectDependencies(
			&integrationtest.TestOptions{
				FixturePath: "./fixtures/test3",
				AppConfig:   unittest.GetMockMainnetProductionConfig(ctrl),
			},
			func(
				processor processor.Processor,
				repo repository.Repository,
			) {
				positionAddress := "3ob8nYeffN5TUBDLJsjsVLMmf8GhkPkpUdAcJctEWCEZ"
				// position: https://explorer.solana.com/address/3ob8nYeffN5TUBDLJsjsVLMmf8GhkPkpUdAcJctEWCEZ
				assert.NoError(t, processor.UpsertPositionByAddress(context.Background(), positionAddress))
				// vault should exist
				vault, err := repo.AdminGetVaultByAddress(context.Background(), "HJcGW3iQGvLpnaJ7LATnrBdw66MTC1hzmmQMKN8pgVhH")
				assert.NoError(t, err)
				assert.NotNil(t, vault)
				// by extension proto config should exist
				protoConfig, err := repo.GetProtoConfigByAddress(context.Background(), vault.ProtoConfig)
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
				// position should exist
				position, err := repo.GetPositionByAddress(context.Background(), positionAddress)
				assert.NoError(t, err)
				assert.NotNil(t, position)

				assert.Equal(t, "3ob8nYeffN5TUBDLJsjsVLMmf8GhkPkpUdAcJctEWCEZ", position.Pubkey)
				assert.Equal(t, "HJcGW3iQGvLpnaJ7LATnrBdw66MTC1hzmmQMKN8pgVhH", position.Vault)
				assert.Equal(t, "7vRcWTNo7Fz6LhhGTqaeKv3m63uE3hzDPMVdDhNccA9C", position.Authority)
				assert.Equal(t, uint64(0x23c34600), position.DepositedTokenAAmount)
				assert.Equal(t, uint64(0), position.WithdrawnTokenBAmount)
				assert.Equal(t, time.Unix(1669882164, 0).Unix(), position.DepositTimestamp.Unix())
				assert.Equal(t, uint64(0x21), position.DcaPeriodIDBeforeDeposit)
				assert.Equal(t, uint64(0x3f), position.NumberOfSwaps)
				assert.Equal(t, uint64(0x915261), position.PeriodicDripAmount)
				assert.Equal(t, false, position.IsClosed)
				// if the line below needs to be updated, add the field assertion above
				assert.Equal(t, reflect.TypeOf(*position).NumField(), 10)

				// by extension position nft token account should exist
				positionNFTTokenAccounts, err := repo.GetActiveTokenAccountsByMint(context.Background(), position.Authority)
				assert.NoError(t, err)
				assert.Len(t, positionNFTTokenAccounts, 1)

				assert.Equal(t, "9dAqHvbRPK7WXCxSQf6e79fCAQJseS9WdE4D3U2LhBGB", positionNFTTokenAccounts[0].Pubkey)
				assert.Equal(t, "7vRcWTNo7Fz6LhhGTqaeKv3m63uE3hzDPMVdDhNccA9C", positionNFTTokenAccounts[0].Mint)
				assert.Equal(t, "GYvcAPtJKo9ierGgWhryrhrEKCrmZuPmYf1obMLq2uBK", positionNFTTokenAccounts[0].Owner)
				assert.Equal(t, uint64(1), positionNFTTokenAccounts[0].Amount)
				assert.Equal(t, "initialized", positionNFTTokenAccounts[0].State)
				// if the line below needs to be updated, add the field assertion above
				assert.Equal(t, reflect.TypeOf(*positionNFTTokenAccounts[0]).NumField(), 5)
			})
	})
}
