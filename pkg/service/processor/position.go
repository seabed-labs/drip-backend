package processor

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/dcaf-labs/drip/pkg/service/alert"
	"github.com/dcaf-labs/drip/pkg/service/repository/model"
	"github.com/dcaf-labs/solana-go-clients/pkg/drip"
	"github.com/sirupsen/logrus"
)

func (p impl) UpsertPositionByAddress(ctx context.Context, address string) error {
	var position drip.Position
	if err := p.solanaClient.GetAccount(ctx, address, &position); err != nil {
		return err
	}
	return p.UpsertPosition(ctx, address, position)
}

func (p impl) UpsertPosition(ctx context.Context, address string, position drip.Position) error {
	log := logrus.WithField("address", address).WithField("operation", "UpsertPosition")
	vault, err := p.ensureVault(ctx, position.Vault.String())
	if err != nil {
		return fmt.Errorf("failed to ensureVault, err: %w", err)
	}
	// Get up to date token metadata info
	if err := p.UpsertTokenByAddress(ctx, vault.TokenAMint); err != nil {
		return fmt.Errorf("failed to UpsertTokenByAddress, err: %w", err)
	}
	if err := p.UpsertTokenByAddress(ctx, vault.TokenBMint); err != nil {
		return err
	}
	if err := p.UpsertTokenByAddress(ctx, position.PositionAuthority.String()); err != nil {
		return err
	}
	//shouldAlert := func() bool {
	//	_, err := p.repo.GetPositionByAddress(ctx, address)
	//	return err != nil && err.Error() == repository.ErrRecordNotFound
	//}()
	shouldAlert := true

	if err := p.repo.UpsertPositions(ctx, &model.Position{
		Pubkey:                   address,
		Vault:                    position.Vault.String(),
		Authority:                position.PositionAuthority.String(),
		DepositedTokenAAmount:    position.DepositedTokenAAmount,
		WithdrawnTokenBAmount:    position.WithdrawnTokenBAmount,
		DepositTimestamp:         time.Unix(position.DepositTimestamp, 0),
		DcaPeriodIDBeforeDeposit: position.DripPeriodIdBeforeDeposit,
		NumberOfSwaps:            position.NumberOfSwaps,
		PeriodicDripAmount:       position.PeriodicDripAmount,
		IsClosed:                 position.IsClosed,
	}); err != nil {
		log.WithError(err).Error("failed to UpsertPositions in UpsertPosition")
		return err
	}
	largestAccounts, err := p.solanaClient.GetLargestTokenAccounts(ctx, position.PositionAuthority.String())
	if err != nil {
		return err
	}
	for _, account := range largestAccounts {
		if account == nil {
			continue
		}
		if err := p.UpsertTokenAccountBalanceByAddress(ctx, account.Address.String()); err != nil {
			log.WithError(err).Error("failed to UpsertTokenAccountBalanceByAddress in UpsertPosition")
		}
	}
	if shouldAlert {
		if err := p.sendNewPositionAlert(ctx, address); err != nil {
			log.WithError(err).Error("failed to sendNewPositionAlert")
		}

	}
	return nil
}

func (p impl) sendNewPositionAlert(ctx context.Context, positionAddr string) error {
	position, err := p.repo.GetPositionByAddress(ctx, positionAddr)
	if err != nil {
		return err
	}
	vault, err := p.repo.GetVaultByAddress(ctx, position.Vault)
	if err != nil {
		return err
	}
	protoConfig, err := p.repo.GetProtoConfigByAddress(ctx, vault.ProtoConfig)
	if err != nil {
		return err
	}
	tokenA, err := p.repo.GetTokenByAddress(ctx, vault.TokenAMint)
	if err != nil {
		return err
	}
	tokenB, err := p.repo.GetTokenByAddress(ctx, vault.TokenBMint)
	if err != nil {
		return err
	}
	balances, err := p.repo.GetActiveTokenAccountBalancesByMint(ctx, position.Authority)
	if err != nil {
		return err
	}
	if len(balances) != 1 {
		return fmt.Errorf("invalid number of balnaces returned for mint %s, len(balances %d", position.Authority, len(balances))
	}
	scaledTokenADepositAmount := float64(position.DepositedTokenAAmount) / math.Pow(10, float64(tokenA.Decimals))
	scaledPeriodicDripAmount := float64(position.PeriodicDripAmount) / math.Pow(10, float64(tokenA.Decimals))

	newPositionAlert := alert.NewPositionAlert{
		TokenASymbol:              tokenA.Symbol,
		TokenAIconURL:             tokenA.IconURL,
		TokenAMint:                tokenA.Pubkey,
		TokenBSymbol:              tokenB.Symbol,
		TokenBIconURL:             tokenB.IconURL,
		TokenBMint:                tokenB.Pubkey,
		ScaledTokenADepositAmount: scaledTokenADepositAmount,
		Granularity:               protoConfig.Granularity,
		ScaledDripAmount:          scaledPeriodicDripAmount,
		NumberOfSwaps:             position.NumberOfSwaps,
		Owner:                     balances[0].Owner,
		Position:                  positionAddr,
	}

	return p.alertService.SendNewPositionAlert(ctx, newPositionAlert)
}
