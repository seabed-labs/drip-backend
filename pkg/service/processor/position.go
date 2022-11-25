package processor

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/gagliardetto/solana-go/rpc"

	"github.com/samber/lo"

	"github.com/dcaf-labs/drip/pkg/service/repository"

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
	shouldAlert := func() bool {
		_, err := p.repo.GetPositionByAddress(ctx, address)
		return err != nil && err.Error() == repository.ErrRecordNotFound
	}()

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
	allTokenAccountPubkeys := func() (ret []string) {
		// Update existing account balances for this position if any (relevant when backfilling)
		if existingBalances, err := p.repo.GetActiveTokenAccountsByMint(ctx, position.PositionAuthority.String()); err != nil {
			log.WithError(err).Warning("failed to GetActiveTokenAccountsByMint")
		} else {
			tokenAccountPubkeys := lo.Map[*model.TokenAccount, string](existingBalances, func(tokenAccount *model.TokenAccount, _ int) string {
				return tokenAccount.Pubkey
			})
			ret = append(ret, tokenAccountPubkeys...)
		}
		// Add new token account balance if any
		if largestAccounts, err := p.solanaClient.GetLargestTokenAccounts(ctx, position.PositionAuthority.String()); err != nil {
			log.WithError(err).Error("failed to GetLargestTokenAccounts")
		} else {
			tokenAccountPubkeys := lo.FilterMap[*rpc.TokenLargestAccountsResult, string](largestAccounts, func(account *rpc.TokenLargestAccountsResult, _ int) (string, bool) {
				if account == nil {
					return "", false
				}
				return account.Address.String(), true
			})
			ret = append(ret, tokenAccountPubkeys...)
		}
		ret = lo.FindUniquesBy[string, string](ret, func(pubkey string) string { return pubkey })
		return ret
	}()
	if err := p.UpsertTokenAccountsByAddresses(ctx, allTokenAccountPubkeys...); err != nil {
		log.WithError(err).Warning("failed to UpsertTokenAccountsByAddresses")
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
	vault, err := p.ensureVault(ctx, position.Vault)
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
	balances, err := p.repo.GetActiveTokenAccountsByMint(ctx, position.Authority)
	if err != nil {
		return err
	}
	owner := "unknown"
	if len(balances) == 1 {
		owner = balances[0].Owner
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
		Owner:                     owner,
		Position:                  positionAddr,
	}

	return p.alertService.SendNewPositionAlert(ctx, newPositionAlert)
}
