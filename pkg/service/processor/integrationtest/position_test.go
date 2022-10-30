package integrationtest

import (
	"context"
	"testing"

	"github.com/dcaf-labs/drip/pkg/integrationtest"
	"github.com/dcaf-labs/drip/pkg/service/processor"
	"github.com/dcaf-labs/drip/pkg/service/repository"
	"github.com/test-go/testify/assert"
)

// TODO: We should mock out api calls with mock data or an api replay
func TestHandler_UpsertPositionByAddress(t *testing.T) {
	// this takes 21 seconds ☠️
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
			mints, err := repo.GetTokensByMints(context.Background(), []string{vault.TokenAMint, vault.TokenBMint})
			assert.NoError(t, err)
			assert.Len(t, mints, 2)
			// position should exist
			position, err := repo.GetPositionByAddress(context.Background(), positionAddress)
			assert.NoError(t, err)
			assert.NotNil(t, position)
			// by extension position nft token account should exist
			positionNFTTokenAccounts, err := repo.GetActiveTokenAccountsByMint(context.Background(), position.Authority)
			assert.NoError(t, err)
			assert.Len(t, positionNFTTokenAccounts, 1)
		})
}
