package integrationtest

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	"github.com/dcaf-labs/drip/pkg/api/apispec"
	controller "github.com/dcaf-labs/drip/pkg/api/routes"
	"github.com/dcaf-labs/drip/pkg/integrationtest"
	"github.com/dcaf-labs/drip/pkg/service/processor"
	"github.com/dcaf-labs/drip/pkg/service/repository"
	"github.com/dcaf-labs/drip/pkg/service/repository/model"
	"github.com/dcaf-labs/drip/pkg/unittest"
	"github.com/gagliardetto/solana-go"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/test-go/testify/assert"
)

func TestHandler_GetV1AdminPositions(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration tests in short mode")
	}
	ctrl := gomock.NewController(t)
	mockConfig := unittest.GetMockMainnetProductionConfig(ctrl)
	e := echo.New()

	t.Run("should return one expanded position", func(t *testing.T) {
		integrationtest.InjectDependencies(
			&integrationtest.TestOptions{
				FixturePath: "./fixtures/get_v1_admin_expanded_positions_test_test1",
				AppConfig:   mockConfig,
			},
			func(
				handler *controller.Handler,
				processor processor.Processor,
				repo repository.Repository,
			) {
				tokenA, tokenB, protoConfig, vault, position := setupPosition(t, repo)
				c, rec := unittest.GetTestRequestRecorder(e, nil)
				params := apispec.GetV1AdminPositionsParams{
					Expand: &apispec.ExpandAdminPositionsQueryParam{"all"},
				}
				assert.NoError(t, handler.GetV1AdminPositions(c, params))
				assert.Equal(t, rec.Code, http.StatusOK)
				expandedPositions, err := apispec.ParseGetV1AdminPositionsResponse(rec.Result())
				assert.NoError(t, err)
				assert.NotNil(t, expandedPositions.JSON200)
				assert.Equal(t, len(*expandedPositions.JSON200), 1)

				assert.Equal(t, (*expandedPositions.JSON200)[0].Position.Pubkey, position.Pubkey)
				assert.Equal(t, (*expandedPositions.JSON200)[0].Position.Authority, position.Authority)
				// since we expanded everything, we expect Position.Vault (embed) to be overridden by the vault model
				assert.Equal(t, (*expandedPositions.JSON200)[0].Position.Vault, "")

				assert.Equal(t, (*expandedPositions.JSON200)[0].Vault.Pubkey, position.Vault)
				assert.Equal(t, (*expandedPositions.JSON200)[0].Vault.ProtoConfig, vault.ProtoConfig)
				assert.Equal(t, (*expandedPositions.JSON200)[0].Vault.TokenAAccount, vault.TokenAAccount)
				assert.Equal(t, (*expandedPositions.JSON200)[0].Vault.TokenAMint, vault.TokenAMint)
				assert.Equal(t, (*expandedPositions.JSON200)[0].Vault.TokenBAccount, vault.TokenBAccount)
				assert.Equal(t, (*expandedPositions.JSON200)[0].Vault.TokenBMint, vault.TokenBMint)
				assert.Equal(t, (*expandedPositions.JSON200)[0].Vault.TreasuryTokenBAccount, vault.TreasuryTokenBAccount)

				assert.Equal(t, (*expandedPositions.JSON200)[0].ProtoConfig.Pubkey, vault.ProtoConfig)
				assert.Equal(t, (*expandedPositions.JSON200)[0].ProtoConfig.Admin, protoConfig.Admin)
				assert.Equal(t, (*expandedPositions.JSON200)[0].ProtoConfig.Admin, protoConfig.Admin)

				assert.Equal(t, (*expandedPositions.JSON200)[0].TokenA.Pubkey, vault.TokenAMint)
				assert.Equal(t, (*expandedPositions.JSON200)[0].TokenA.Pubkey, tokenA.Pubkey)

				assert.Equal(t, (*expandedPositions.JSON200)[0].TokenB.Pubkey, vault.TokenBMint)
				assert.Equal(t, (*expandedPositions.JSON200)[0].TokenB.Pubkey, tokenB.Pubkey)
			})
	})
}

func setupPosition(t *testing.T, repo repository.Repository) (tokenA, tokenB *model.Token, protoConfig *model.ProtoConfig, vault *model.Vault, position *model.Position) {
	tokenAPubkey := "orcaEKTdK7LKz57vaAYr9QeNsVEPfiu6QeMU1kektZE"
	tokenBPubkey := "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v"
	protoConfigPubkey := "8V6NbywJAWBSmg6eGCYYZymRc89EVgH1Bw7JRTwcqSaJ"
	vaultPubkey := "8qdV9meyRk7Fwt9EkRveN7kWPBiNtmH9AJyYmAxXFXjY"
	vaultTokenAAccount := "76ZEX9d2FuN9Pu7F7pDZykz8gxwnS5FqdXbBmNvFhm6h"
	vaultTokenBAccount := "76ZEX9d2FuN9Pu7F7pDZykz8gxwnS5FqdXbBmNvFhm6h"
	treasuryTokenBAccount := "5cuSQkFYZa26XQn65McbKUVBFgViCDK5rAjAZwLZNKAd"

	tokenA = &model.Token{
		Pubkey:   tokenAPubkey,
		Decimals: 6,
	}
	tokenB = &model.Token{
		Pubkey:   tokenBPubkey,
		Decimals: 9,
	}
	assert.NoError(t, repo.UpsertTokens(context.Background(), tokenA, tokenB))
	tokenPairID := uuid.New()
	assert.NoError(t, repo.InsertTokenPairs(context.Background(), &model.TokenPair{
		ID:     tokenPairID.String(),
		TokenA: tokenAPubkey,
		TokenB: tokenBPubkey,
	}))
	protoConfig = &model.ProtoConfig{
		Pubkey:                  protoConfigPubkey,
		Granularity:             60,
		TokenADripTriggerSpread: 50,
		TokenBWithdrawalSpread:  50,
		Admin:                   solana.NewWallet().PublicKey().String(),
		TokenBReferralSpread:    50,
	}
	assert.NoError(t, repo.UpsertProtoConfigs(context.Background(), protoConfig))
	vault = &model.Vault{
		Pubkey:                vaultPubkey,
		ProtoConfig:           protoConfigPubkey,
		TokenAAccount:         vaultTokenAAccount,
		TokenBAccount:         vaultTokenBAccount,
		TreasuryTokenBAccount: treasuryTokenBAccount,
		TokenPairID:           tokenPairID.String(),
		TokenAMint:            tokenAPubkey,
		TokenBMint:            tokenBPubkey,
		Enabled:               true,
	}
	assert.NoError(t, repo.UpsertVaults(context.Background(), vault))
	position = &model.Position{
		Pubkey:                   solana.NewWallet().PublicKey().String(),
		Vault:                    vaultPubkey,
		Authority:                solana.NewWallet().PublicKey().String(),
		DepositedTokenAAmount:    9000,
		WithdrawnTokenBAmount:    1000,
		DepositTimestamp:         time.Now(),
		DcaPeriodIDBeforeDeposit: 10,
		NumberOfSwaps:            10,
		PeriodicDripAmount:       5,
		IsClosed:                 false,
	}
	assert.NoError(t, repo.UpsertPositions(context.Background(), position))
	return tokenA, tokenB, protoConfig, vault, position
}
