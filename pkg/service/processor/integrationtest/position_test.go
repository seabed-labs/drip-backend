package integrationtest

import (
	"context"
	"reflect"
	"testing"

	"github.com/dcaf-labs/drip/pkg/integrationtest"
	"github.com/dcaf-labs/drip/pkg/service/processor"
	"github.com/dcaf-labs/drip/pkg/service/repository"
	"github.com/test-go/testify/assert"
)

func TestHandler_UpsertPositionByAddress(t *testing.T) {
	integrationtest.InjectDependencies(
		&integrationtest.APIRecorderOptions{
			Path: "./fixtures/upsert-position-by-address",
		},
		func(
			processor processor.Processor,
			repo repository.Repository,
		) {
			positionAddress := "46kbd7mjJhtxWW19wyh1uUgbFx8PktRBDgsY2BM9doRA"
			// position: https://explorer.solana.com/address/46kbd7mjJhtxWW19wyh1uUgbFx8PktRBDgsY2BM9doRA/anchor-account?cluster=devnet
			assert.NoError(t, processor.UpsertPositionByAddress(context.Background(), positionAddress))
			// vault should exist
			vault, err := repo.AdminGetVaultByAddress(context.Background(), "BwJj7DYyMR1xMnWK1PGKPLi5u2ZP5EDBFBkPAAv4UDP8")
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
			mints, err := repo.GetTokensByAddresses(context.Background(), []string{vault.TokenAMint, vault.TokenBMint})
			assert.NoError(t, err)
			assert.Len(t, mints, 2)
			// position should exist
			position, err := repo.GetPositionByAddress(context.Background(), positionAddress)
			assert.NoError(t, err)
			assert.NotNil(t, position)

			assert.Equal(t, "46kbd7mjJhtxWW19wyh1uUgbFx8PktRBDgsY2BM9doRA", position.Pubkey)
			assert.Equal(t, "BwJj7DYyMR1xMnWK1PGKPLi5u2ZP5EDBFBkPAAv4UDP8", position.Vault)
			assert.Equal(t, "JB4DSzoQ5SyxLETtGUrC5s1MvS6Y78PTmzczszZ7ZPtB", position.Authority)
			assert.Equal(t, uint64(101000000), position.DepositedTokenAAmount)
			assert.Equal(t, uint64(0), position.WithdrawnTokenBAmount)
			assert.Equal(t, int64(1667137701), position.DepositTimestamp.Unix())
			assert.Equal(t, uint64(1847), position.DcaPeriodIDBeforeDeposit)
			assert.Equal(t, uint64(4), position.NumberOfSwaps)
			assert.Equal(t, uint64(25250000), position.PeriodicDripAmount)
			assert.Equal(t, false, position.IsClosed)
			// if the line below needs to be updated, add the field assertion above
			assert.Equal(t, reflect.TypeOf(*position).NumField(), 10)

			// by extension position nft token account should exist
			positionNFTTokenAccounts, err := repo.GetActiveTokenAccountsByMint(context.Background(), position.Authority)
			assert.NoError(t, err)
			assert.Len(t, positionNFTTokenAccounts, 1)

			assert.Equal(t, "6tZikLFcRCvbpQca8bg9XrCNcnzwhfepBKWHcj8sLfjc", positionNFTTokenAccounts[0].Pubkey)
			assert.Equal(t, "JB4DSzoQ5SyxLETtGUrC5s1MvS6Y78PTmzczszZ7ZPtB", positionNFTTokenAccounts[0].Mint)
			assert.Equal(t, "3CTkqdcjzn1ptnNYUcFq4Sk1Smy91kz5p9JhgJBHGe3e", positionNFTTokenAccounts[0].Owner)
			assert.Equal(t, uint64(1), positionNFTTokenAccounts[0].Amount)
			assert.Equal(t, "initialized", positionNFTTokenAccounts[0].State)
			// if the line below needs to be updated, add the field assertion above
			assert.Equal(t, reflect.TypeOf(*positionNFTTokenAccounts[0]).NumField(), 5)
		})
}
