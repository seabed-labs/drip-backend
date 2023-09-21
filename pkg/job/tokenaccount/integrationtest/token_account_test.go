package integrationtest

import (
	"context"
	"testing"

	"github.com/dcaf-labs/drip/pkg/integrationtest"
	"github.com/dcaf-labs/drip/pkg/job/tokenaccount"
	"github.com/dcaf-labs/drip/pkg/service/clients/coingecko"
	"github.com/dcaf-labs/drip/pkg/service/processor"
	"github.com/dcaf-labs/drip/pkg/service/repository"
	"github.com/dcaf-labs/drip/pkg/service/repository/model"
	"github.com/dcaf-labs/drip/pkg/unittest"
	"github.com/gagliardetto/solana-go"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/test-go/testify/assert"
	"go.uber.org/fx"
)

type dummyInterface struct {
}

func (d dummyInterface) Append(_ fx.Hook) {

}
func TestTokenAccountJob(t *testing.T) {
	t.Skip("skipping test...")
	//if testing.Short() {
	//	t.Skip("skipping integration tests in short mode")
	//}
	ctrl := gomock.NewController(t)
	t.Run("should update 2 token accounts with metadata", func(t *testing.T) {
		mockConfig := unittest.GetMockMainnetProductionConfig(ctrl)
		integrationtest.TestWithInjectedDependencies(
			&integrationtest.TestOptions{
				FixturePath: "./fixtures/test1",
				AppConfig:   mockConfig,
			},
			func(
				repo repository.Repository,
				processor processor.Processor,
				coinGeckoClient coingecko.CoinGeckoClient,
			) {
				setupVaults(t, repo)

				tokenAccountJob, err := tokenaccount.NewTokenAccountJob(dummyInterface{}, repo, processor)
				assert.NoError(t, err)
				tokenAccountJob.UpsertAllVaultTokensAccountAccounts(context.Background())

				tokenAccount, err := repo.GetTokenAccountsByAddresses(context.Background(), "9UFPCaLbsst45ytt4ksUoAnL66nic3AicpGN6yVKW9Az")
				assert.NoError(t, err)
				assert.Len(t, tokenAccount, 1)
				assert.Equal(t, "9UFPCaLbsst45ytt4ksUoAnL66nic3AicpGN6yVKW9Az", tokenAccount[0].Pubkey)
				assert.Equal(t, uint64(0x6), tokenAccount[0].Amount)
				assert.Equal(t, "initialized", tokenAccount[0].State)
				assert.Equal(t, "J4WLLfYFJJpyBcpDeRGCfwPQgN36Ai1dfmLiRxgxSXhu", tokenAccount[0].Owner)
				assert.Equal(t, "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v", tokenAccount[0].Mint)

				tokenAccount, err = repo.GetTokenAccountsByAddresses(context.Background(), "3W8k2S2wUkZVBmeFFS19QfMLkbbxABzjiFsYHLfZtSC4")
				assert.NoError(t, err)
				assert.Len(t, tokenAccount, 1)
				assert.Equal(t, "3W8k2S2wUkZVBmeFFS19QfMLkbbxABzjiFsYHLfZtSC4", tokenAccount[0].Pubkey)
				assert.Equal(t, uint64(0x1e05d68d7f), tokenAccount[0].Amount)
				assert.Equal(t, "initialized", tokenAccount[0].State)
				assert.Equal(t, "J4WLLfYFJJpyBcpDeRGCfwPQgN36Ai1dfmLiRxgxSXhu", tokenAccount[0].Owner)
				assert.Equal(t, "7xKXtg2CW87d97TXJSDpbD5jBkheTqA83TZRuJosgAsU", tokenAccount[0].Mint)

				tokenAccount, err = repo.GetTokenAccountsByAddresses(context.Background(), "5bmpKycNDC5hkS3eGZDwYctJKJF6BtCF2ebYoVFjN3er")
				assert.NoError(t, err)
				assert.Len(t, tokenAccount, 1)
				assert.Equal(t, "5bmpKycNDC5hkS3eGZDwYctJKJF6BtCF2ebYoVFjN3er", tokenAccount[0].Pubkey)
				assert.Equal(t, uint64(0x3a639915), tokenAccount[0].Amount)
				assert.Equal(t, "initialized", tokenAccount[0].State)
				assert.Equal(t, "JC5NuYudj4Dd8vFDS8re1spg3PZEZ3PZyBRef852vz5R", tokenAccount[0].Owner)
				assert.Equal(t, "7xKXtg2CW87d97TXJSDpbD5jBkheTqA83TZRuJosgAsU", tokenAccount[0].Mint)

			})
	})
}

func setupVaults(t *testing.T, repo repository.Repository) {
	tokenA := "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"
	tokenB := "7xKXtg2CW87d97TXJSDpbD5jBkheTqA83TZRuJosgAsU"
	protoConfig := "BK7m7aEb5SrJBSAzMBGumnRhRvF3C7j7fssUrPKPTwxG"
	assert.NoError(t, repo.UpsertTokens(context.Background(), &model.Token{
		Pubkey:   tokenA,
		Decimals: 6,
	}, &model.Token{
		Pubkey:   tokenB,
		Decimals: 9,
	}))
	tokenPairID := uuid.New()
	assert.NoError(t, repo.InsertTokenPairs(context.Background(), &model.TokenPair{
		ID:     tokenPairID.String(),
		TokenA: tokenA,
		TokenB: tokenB,
	}))
	assert.NoError(t, repo.UpsertProtoConfigs(context.Background(), &model.ProtoConfig{
		Pubkey:                  protoConfig,
		Granularity:             60,
		TokenADripTriggerSpread: 50,
		TokenBWithdrawalSpread:  50,
		Admin:                   solana.NewWallet().PublicKey().String(),
		TokenBReferralSpread:    50,
	}))
	assert.NoError(t, repo.UpsertVaults(context.Background(), &model.Vault{
		Pubkey:                "J4WLLfYFJJpyBcpDeRGCfwPQgN36Ai1dfmLiRxgxSXhu",
		ProtoConfig:           protoConfig,
		TokenAAccount:         "9UFPCaLbsst45ytt4ksUoAnL66nic3AicpGN6yVKW9Az",
		TokenBAccount:         "3W8k2S2wUkZVBmeFFS19QfMLkbbxABzjiFsYHLfZtSC4",
		TreasuryTokenBAccount: "5bmpKycNDC5hkS3eGZDwYctJKJF6BtCF2ebYoVFjN3er",
		TokenPairID:           tokenPairID.String(),
		TokenAMint:            tokenA,
		TokenBMint:            tokenB,
		Enabled:               true,
	}))
}
