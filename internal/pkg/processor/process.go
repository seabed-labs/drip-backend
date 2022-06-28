package processor

import (
	"context"
	"time"

	"github.com/dcaf-protocol/drip/internal/pkg/clients/solana"
	"github.com/dcaf-protocol/drip/internal/pkg/clients/solana/drip"
	"github.com/dcaf-protocol/drip/internal/pkg/clients/solana/token_swap"
	"github.com/dcaf-protocol/drip/internal/pkg/repository"
	"github.com/dcaf-protocol/drip/internal/pkg/repository/model"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Processor interface {
	UpsertProtoConfigByAddress(context.Context, string) error
	UpsertVaultByAddress(context.Context, string) error
	UpsertPositionByAddress(context.Context, string) error
	UpsertVaultPeriodByAddress(context.Context, string) error
	UpsertTokenSwapByAddress(context.Context, string) error
	UpsertTokenPair(context.Context, string, string) error
}

type impl struct {
	repo   repository.Repository
	client solana.Solana
}

func NewProcessor(
	repo repository.Repository,
	client solana.Solana,
) Processor {
	return impl{
		repo:   repo,
		client: client,
	}
}

func (p impl) UpsertTokenSwapByAddress(ctx context.Context, address string) error {
	var tokenSwap token_swap.TokenSwap
	if err := p.client.GetAccount(ctx, address, &tokenSwap); err != nil {
		return err
	}
	var tokenLPMint token.Mint
	if err := p.client.GetAccount(ctx, tokenSwap.TokenPool.String(), &tokenLPMint); err != nil {
		return err
	}
	tokenPair, err := p.ensureTokenPair(ctx, tokenSwap.MintA.String(), tokenSwap.MintB.String())
	if err != nil {
		return err
	}
	return p.repo.UpsertTokenSwaps(ctx, &model.TokenSwap{
		Pubkey:        address,
		Mint:          tokenSwap.TokenPool.String(),
		Authority:     tokenLPMint.MintAuthority.String(),
		FeeAccount:    tokenSwap.FeeAccount.String(),
		TokenAAccount: tokenSwap.TokenAccountA.String(),
		TokenBAccount: tokenSwap.TokenAccountB.String(),
		Pair:          tokenPair.ID,
	})
}

func (p impl) UpsertProtoConfigByAddress(ctx context.Context, address string) error {
	var protoConfig drip.VaultProtoConfig
	if err := p.client.GetAccount(ctx, address, &protoConfig); err != nil {
		return err
	}
	return p.repo.UpsertProtoConfigs(ctx, &model.ProtoConfig{
		Pubkey:               address,
		Granularity:          protoConfig.Granularity,
		TriggerDcaSpread:     protoConfig.TriggerDcaSpread,
		BaseWithdrawalSpread: protoConfig.BaseWithdrawalSpread,
	})
}

func (p impl) UpsertVaultByAddress(ctx context.Context, address string) error {
	var vaultAccount drip.Vault
	if err := p.client.GetAccount(ctx, address, &vaultAccount); err != nil {
		return err
	}
	if err := p.UpsertProtoConfigByAddress(ctx, vaultAccount.ProtoConfig.String()); err != nil {
		return nil
	}
	tokenPair, err := p.ensureTokenPair(ctx, vaultAccount.TokenAMint.String(), vaultAccount.TokenBMint.String())
	if err != nil {
		return err
	}
	// TODO(Mocha): If exists - backfill vaultPeriods in goRoutine, if not exists - upsert
	return p.repo.UpsertVaults(ctx, &model.Vault{
		Pubkey:                 address,
		ProtoConfig:            vaultAccount.ProtoConfig.String(),
		TokenAAccount:          vaultAccount.TokenAAccount.String(),
		TokenBAccount:          vaultAccount.TokenBAccount.String(),
		TreasuryTokenBAccount:  vaultAccount.TreasuryTokenBAccount.String(),
		LastDcaPeriod:          vaultAccount.LastDcaPeriod,
		DripAmount:             vaultAccount.DripAmount,
		DcaActivationTimestamp: time.Unix(vaultAccount.DcaActivationTimestamp, 0),
		Enabled:                false,
		TokenPairID:            tokenPair.ID,
	})
}

func (p impl) UpsertPositionByAddress(ctx context.Context, address string) error {
	var position drip.Position
	if err := p.client.GetAccount(ctx, address, &position); err != nil {
		return err
	}
	if _, err := p.ensureVault(ctx, position.Vault.String()); err != nil {
		return err
	}
	return p.repo.UpsertPositions(ctx, &model.Position{
		Pubkey:                   address,
		Vault:                    position.Vault.String(),
		Authority:                position.PositionAuthority.String(),
		DepositedTokenAAmount:    position.DepositedTokenAAmount,
		WithdrawnTokenBAmount:    position.WithdrawnTokenBAmount,
		DepositTimestamp:         time.Unix(position.DepositTimestamp, 0),
		DcaPeriodIDBeforeDeposit: position.DcaPeriodIdBeforeDeposit,
		NumberOfSwaps:            position.NumberOfSwaps,
		PeriodicDripAmount:       position.PeriodicDripAmount,
		IsClosed:                 position.IsClosed,
	})
}

func (p impl) UpsertVaultPeriodByAddress(ctx context.Context, address string) error {
	var vaultPeriodAccount drip.VaultPeriod
	if err := p.client.GetAccount(ctx, address, &vaultPeriodAccount); err != nil {
		return err
	}
	twap, err := decimal.NewFromString(vaultPeriodAccount.Twap.String())
	if err != nil {
		return err
	}
	if _, err := p.ensureVault(ctx, vaultPeriodAccount.Vault.String()); err != nil {
		return err
	}
	return p.repo.UpsertVaultPeriods(ctx, &model.VaultPeriod{
		Pubkey:   address,
		Vault:    vaultPeriodAccount.Vault.String(),
		PeriodID: vaultPeriodAccount.PeriodId,
		Twap:     twap,
		Dar:      vaultPeriodAccount.Dar,
	})
}

func (p impl) UpsertTokenPair(ctx context.Context, tokenAAMint string, tokenBMint string) error {
	var tokenA token.Mint
	if err := p.client.GetAccount(ctx, tokenAAMint, &tokenA); err != nil {
		return err
	}
	var tokenB token.Mint
	if err := p.client.GetAccount(ctx, tokenBMint, &tokenB); err != nil {
		return err
	}
	if err := p.repo.UpsertTokens(ctx,
		&model.Token{
			Pubkey:   tokenAAMint,
			Symbol:   nil,
			Decimals: int16(tokenA.Decimals),
			IconURL:  nil,
		}, &model.Token{
			Pubkey:   tokenBMint,
			Symbol:   nil,
			Decimals: int16(tokenB.Decimals),
			IconURL:  nil,
		}); err != nil {
		return err
	}
	return p.repo.UpsertTokenPairs(ctx, &model.TokenPair{
		ID:     uuid.New().String(),
		TokenA: tokenAAMint,
		TokenB: tokenBMint,
	})
}

// ensureTokenPair - if token pair exists return it, else upsert tokenPair and all needed tokenPair foreign keys
func (p impl) ensureTokenPair(ctx context.Context, tokenAAMint string, tokenBMint string) (*model.TokenPair, error) {
	tokenPair, err := p.repo.GetTokenPair(ctx, tokenAAMint, tokenBMint)
	if err != nil && err.Error() == "record not found" {
		if err := p.UpsertTokenPair(ctx, tokenAAMint, tokenBMint); err != nil {
			return nil, err
		}
		return p.repo.GetTokenPair(ctx, tokenAAMint, tokenBMint)
	}

	return tokenPair, err
}

// ensureVault - if vault exists return it , else upsert vault and all needed vault foreign keys
func (p impl) ensureVault(ctx context.Context, address string) (*model.Vault, error) {
	vault, err := p.repo.GetVaultByAddress(ctx, address)
	if err != nil && err.Error() == "record not found" {
		if err := p.UpsertVaultByAddress(ctx, address); err != nil {
			return nil, err
		}
		return p.repo.GetVaultByAddress(ctx, address)
	}
	return vault, err
}
