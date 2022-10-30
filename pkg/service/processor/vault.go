package processor

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/dcaf-labs/drip/pkg/service/utils"
	"github.com/shopspring/decimal"

	"github.com/dcaf-labs/drip/pkg/service/repository"
	"github.com/dcaf-labs/drip/pkg/service/repository/model"
	"github.com/dcaf-labs/solana-go-clients/pkg/drip"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

func (p impl) UpsertProtoConfigByAddress(ctx context.Context, address string) error {
	var protoConfig drip.VaultProtoConfig
	if err := p.solanaClient.GetAccount(ctx, address, &protoConfig); err != nil {
		return err
	}
	return p.repo.UpsertProtoConfigs(ctx, &model.ProtoConfig{
		Pubkey:                  address,
		Admin:                   protoConfig.Admin.String(),
		Granularity:             protoConfig.Granularity,
		TokenADripTriggerSpread: protoConfig.TokenADripTriggerSpread,
		TokenBWithdrawalSpread:  protoConfig.TokenBWithdrawalSpread,
		TokenBReferralSpread:    protoConfig.TokenBReferralSpread,
	})
}

func (p impl) UpsertVaultByAddress(ctx context.Context, address string) error {
	var vaultAccount drip.Vault
	if err := p.solanaClient.GetAccount(ctx, address, &vaultAccount); err != nil {
		return err
	}
	if err := p.UpsertProtoConfigByAddress(ctx, vaultAccount.ProtoConfig.String()); err != nil {
		return err
	}
	tokenPair, err := p.ensureTokenPair(ctx, vaultAccount.TokenAMint.String(), vaultAccount.TokenBMint.String())
	if err != nil {
		return err
	}

	if err := p.repo.UpsertVaults(ctx, &model.Vault{
		Pubkey:                 address,
		ProtoConfig:            vaultAccount.ProtoConfig.String(),
		TokenAMint:             vaultAccount.TokenAMint.String(),
		TokenBMint:             vaultAccount.TokenBMint.String(),
		TokenAAccount:          vaultAccount.TokenAAccount.String(),
		TokenBAccount:          vaultAccount.TokenBAccount.String(),
		TreasuryTokenBAccount:  vaultAccount.TreasuryTokenBAccount.String(),
		LastDcaPeriod:          vaultAccount.LastDripPeriod,
		DripAmount:             vaultAccount.DripAmount,
		DcaActivationTimestamp: time.Unix(vaultAccount.DripActivationTimestamp, 0),
		Enabled:                false,
		TokenPairID:            tokenPair.ID,
		MaxSlippageBps:         int32(vaultAccount.MaxSlippageBps),
	}); err != nil {
		return err
	}

	var vaultWhitelists []*model.VaultWhitelist
	for i := range vaultAccount.WhitelistedSwaps {
		whitelistedSwap := vaultAccount.WhitelistedSwaps[i]
		if whitelistedSwap.IsZero() {
			continue
		}
		vaultWhitelists = append(vaultWhitelists, &model.VaultWhitelist{
			ID:              uuid.New().String(),
			VaultPubkey:     address,
			TokenSwapPubkey: whitelistedSwap.String(),
		})
	}
	if err := p.repo.UpsertVaultWhitelists(ctx, vaultWhitelists...); err != nil {
		logrus.
			WithField("vault", address).
			WithError(err).
			Error("failed to insert vaultWhitelists")
	}
	if err := p.UpsertTokenAccountByAddress(ctx, vaultAccount.TokenAAccount.String()); err != nil {
		return err
	}
	if err := p.UpsertTokenAccountByAddress(ctx, vaultAccount.TokenBAccount.String()); err != nil {
		return err
	}
	if err := p.UpsertTokenAccountByAddress(ctx, vaultAccount.TreasuryTokenBAccount.String()); err != nil {
		return err
	}
	return nil
}

func (p impl) UpsertVaultPeriodByAddress(ctx context.Context, address string) error {
	return p.upsertVaultPeriodByAddress(ctx, address, true)
}

// returns vault tokens in order of tokenA, tokenB
func (p impl) getTokensForVault(ctx context.Context, vaultAddress string) (*model.Token, *model.Token, error) {
	vault, err := p.ensureVault(ctx, vaultAddress)
	if err != nil {
		return nil, nil, err
	}
	tokens, err := p.repo.GetTokensByAddresses(ctx, []string{vault.TokenAMint, vault.TokenBMint})
	if err != nil {
		return nil, nil, err
	}
	if len(tokens) != 2 {
		return nil, nil, fmt.Errorf("invalid number of tokens return for GetTokensByAddresses for id: %s", vault.TokenPairID)
	}
	if tokens[0].Pubkey == vault.TokenAMint {
		return tokens[0], tokens[1], nil
	}
	return tokens[1], tokens[0], nil
}

// ensureVault - if vault exists return it , else upsert vault and all needed vault foreign keys
func (p impl) ensureVault(ctx context.Context, address string) (*model.Vault, error) {
	vault, err := p.repo.AdminGetVaultByAddress(ctx, address)
	if err != nil && err.Error() == repository.ErrRecordNotFound {
		if err := p.UpsertVaultByAddress(ctx, address); err != nil {
			return nil, err
		}
		return p.repo.AdminGetVaultByAddress(ctx, address)
	}
	return vault, err
}

// upsertVaultPeriodByAddress: this is potentially a recursive call
// if shouldUpsertPrice is set to true, we will try and price period[i], which will try to ensure period[i-1]
// if shouldUpsertPrice is set to false, we will not calculate a price and will upsert it to 0
func (p impl) upsertVaultPeriodByAddress(ctx context.Context, address string, shouldUpsertPrice bool) error {
	var vaultPeriodAccount drip.VaultPeriod
	if err := p.solanaClient.GetAccount(ctx, address, &vaultPeriodAccount); err != nil {
		return err
	}
	twap, err := decimal.NewFromString(vaultPeriodAccount.Twap.String())
	if err != nil {
		return err
	}
	if _, err := p.ensureVault(ctx, vaultPeriodAccount.Vault.String()); err != nil {
		return err
	}
	var priceBOverA decimal.Decimal
	if shouldUpsertPrice {
		priceBOverA, err = p.getVaultPeriodPriceBOverA(ctx, vaultPeriodAccount)
		if err != nil {
			logrus.
				WithField("vaultPeriodAddress", address).
				WithField("vaultPeriodId", vaultPeriodAccount.PeriodId).
				WithError(err).Errorf("failed to getVaultPeriodPriceBOverA")
			return err
		}
	}
	return p.repo.UpsertVaultPeriods(ctx, &model.VaultPeriod{
		Pubkey:      address,
		Vault:       vaultPeriodAccount.Vault.String(),
		PeriodID:    vaultPeriodAccount.PeriodId,
		Twap:        twap,
		Dar:         vaultPeriodAccount.Dar,
		PriceBOverA: priceBOverA,
	})
}

// getVaultPeriodPriceBOverA calculate and return normalized price of b over a
// in the following, twap[x] is the normalized twap value (not the x64 value stored on chain)
//
//	p[i] = twap[i]*i - twap[i-1]*(i-1) for i > 0
//	p[i] = twap[i] for i = 0
func (p impl) getVaultPeriodPriceBOverA(ctx context.Context, periodI drip.VaultPeriod) (decimal.Decimal, error) {
	tokenA, tokenB, err := p.getTokensForVault(ctx, periodI.Vault.String())
	if err != nil {
		return decimal.Decimal{}, err
	}
	twapI, err := decimal.NewFromString(periodI.Twap.String())
	if err != nil {
		return decimal.Decimal{}, fmt.Errorf("failed to get twapI decimal, err: %w", err)
	}
	if periodI.PeriodId == 0 {
		return normalizePrice(twapI, decimal.NewFromInt(int64(tokenA.Decimals)), decimal.NewFromInt(int64(tokenB.Decimals))), nil
	}
	periodIPrecedingAddress, err := utils.GetVaultPeriodPDA(periodI.Vault.String(), int64(periodI.PeriodId-1))
	if err != nil {
		return decimal.Decimal{}, fmt.Errorf("failed to GetVaultPeriodPDA, err: %w", err)
	}
	periodIPPreceding, err := p.ensureVaultPeriod(ctx, periodIPrecedingAddress)
	if err != nil {
		return decimal.Decimal{}, fmt.Errorf("failed to ensureVaultPeriod, err: %w", err)
	}
	twapIPreceding, err := decimal.NewFromString(periodIPPreceding.Twap.String())
	if err != nil {
		return decimal.Decimal{}, fmt.Errorf("failed to get twapIPreceeding decimal, err: %w", err)
	}

	periodIID, err := decimal.NewFromString(strconv.FormatUint(periodI.PeriodId, 10))
	if err != nil {
		return decimal.Decimal{}, fmt.Errorf("failed to periodIId decimal, err: %w", err)
	}
	periodIPrecedingID := periodIID.Sub(decimal.NewFromInt(1))
	// average price from period I to period I-1
	rawPrice := twapI.Mul(periodIID).Sub(twapIPreceding.Mul(periodIPrecedingID))
	return normalizePrice(rawPrice, decimal.NewFromInt(int64(tokenA.Decimals)), decimal.NewFromInt(int64(tokenB.Decimals))), nil
}

func normalizePrice(rawPrice, tokenADecimals, tokenBDecimals decimal.Decimal) decimal.Decimal {
	return rawPrice.Div(decimal.NewFromInt(2).Pow(decimal.NewFromInt(64))).
		Mul(decimal.NewFromInt(10).Pow(tokenADecimals)).
		DivRound(decimal.NewFromInt(10).Pow(tokenBDecimals), 64)
}

// ensureVaultPeriod - if vaultPeriod exists return it , else upsert vaultPeriods with a price of 0
func (p impl) ensureVaultPeriod(ctx context.Context, address string) (*model.VaultPeriod, error) {
	vaultPeriod, err := p.repo.GetVaultPeriodByAddress(ctx, address)
	if err != nil && err.Error() == repository.ErrRecordNotFound {
		if err := p.upsertVaultPeriodByAddress(ctx, address, false); err != nil {
			return nil, err
		}
		return p.repo.GetVaultPeriodByAddress(ctx, address)
	}
	return vaultPeriod, err
}
