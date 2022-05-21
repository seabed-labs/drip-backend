package processor

import (
	"context"
	"time"

	"github.com/dcaf-protocol/drip/internal/pkg/clients/solana"
	"github.com/dcaf-protocol/drip/internal/pkg/clients/solana/dca_vault"
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
}

type ProcessorImpl struct {
	repo   repository.Repository
	client solana.Solana
}

func NewProcessor(
	repo repository.Repository,
	client solana.Solana,
) Processor {
	return ProcessorImpl{
		repo:   repo,
		client: client,
	}
}

func (p ProcessorImpl) UpsertProtoConfigByAddress(ctx context.Context, address string) error {
	var protoConfig dca_vault.VaultProtoConfig
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

func (p ProcessorImpl) UpsertVaultByAddress(ctx context.Context, address string) error {
	var vaultAccount dca_vault.Vault
	if err := p.client.GetAccount(ctx, address, &vaultAccount); err != nil {
		return err
	}
	if err := p.UpsertProtoConfigByAddress(ctx, vaultAccount.ProtoConfig.String()); err != nil {
		return nil
	}
	var tokenA token.Mint
	if err := p.client.GetAccount(ctx, vaultAccount.TokenAMint.String(), &tokenA); err != nil {
		return err
	}
	var tokenB token.Mint
	if err := p.client.GetAccount(ctx, vaultAccount.TokenBMint.String(), &tokenB); err != nil {
		return err
	}
	if err := p.repo.UpsertTokens(ctx,
		&model.Token{
			Pubkey:   vaultAccount.TokenAMint.String(),
			Symbol:   nil,
			Decimals: int16(tokenA.Decimals),
			IconURL:  nil,
		}, &model.Token{
			Pubkey:   vaultAccount.TokenBMint.String(),
			Symbol:   nil,
			Decimals: int16(tokenB.Decimals),
			IconURL:  nil,
		}); err != nil {
		return err
	}
	if err := p.repo.UpsertTokenPairs(ctx, &model.TokenPair{
		ID:     uuid.New().String(),
		TokenA: vaultAccount.TokenAMint.String(),
		TokenB: vaultAccount.TokenBMint.String(),
	}); err != nil {
		return err
	}
	tokenPair, err := p.repo.GetTokenPair(ctx, vaultAccount.TokenAMint.String(), vaultAccount.TokenBMint.String())
	if err != nil {
		return err
	}
	// TODO(Mocha): If exists - backfill vaultPeriods in goRoutine
	// if not exists, upsert
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

func (p ProcessorImpl) UpsertPositionByAddress(ctx context.Context, address string) error {
	var position dca_vault.Position
	if err := p.client.GetAccount(ctx, address, &position); err != nil {
		return err
	}
	if err := p.ensureVault(ctx, position.Vault.String()); err != nil {
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

func (p ProcessorImpl) UpsertVaultPeriodByAddress(ctx context.Context, address string) error {
	var vaultPeriodAccount dca_vault.VaultPeriod
	if err := p.client.GetAccount(ctx, address, &vaultPeriodAccount); err != nil {
		return err
	}
	twap, err := decimal.NewFromString(vaultPeriodAccount.Twap.String())
	if err != nil {
		return err
	}
	if err := p.ensureVault(ctx, vaultPeriodAccount.Vault.String()); err != nil {
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

// ensureVault - if vault exists return, else upsert vault
func (p ProcessorImpl) ensureVault(ctx context.Context, address string) error {
	_, err := p.repo.GetVaultByAddress(ctx, address)
	if err != nil && err.Error() == "record not found" {
		if err := p.UpsertVaultByAddress(ctx, address); err != nil {
			return err
		}
		return nil
	}
	return err
}
