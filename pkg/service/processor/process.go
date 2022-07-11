package processor

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/dcaf-protocol/drip/pkg/clients/solana"
	"github.com/dcaf-protocol/drip/pkg/clients/solana/drip"
	"github.com/dcaf-protocol/drip/pkg/clients/solana/token_swap"
	"github.com/dcaf-protocol/drip/pkg/repository"
	model2 "github.com/dcaf-protocol/drip/pkg/repository/model"

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
	UpsertTokenAccountBalanceByAddress(context.Context, string) error
	UpsertTokenAccountBalance(context.Context, string, token.Account) error
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

	// Add swap A -> B
	tokenPair, err := p.ensureTokenPair(ctx, tokenSwap.MintA.String(), tokenSwap.MintB.String())
	if err != nil {
		return err
	}
	if err := p.repo.UpsertTokenSwaps(ctx, &model2.TokenSwap{
		ID:            uuid.New().String(),
		Pubkey:        address,
		Mint:          tokenSwap.TokenPool.String(),
		Authority:     tokenLPMint.MintAuthority.String(),
		FeeAccount:    tokenSwap.FeeAccount.String(),
		TokenAMint:    tokenSwap.MintA.String(),
		TokenAAccount: tokenSwap.TokenAccountA.String(),
		TokenBMint:    tokenSwap.MintB.String(),
		TokenBAccount: tokenSwap.TokenAccountB.String(),
		TokenPairID:   tokenPair.ID,
	}); err != nil {
		return err
	}
	// Add swap B -> A
	tokenPairInverse, err := p.ensureTokenPair(ctx, tokenSwap.MintB.String(), tokenSwap.MintA.String())
	if err != nil {
		return err
	}
	if err := p.repo.UpsertTokenSwaps(ctx, &model2.TokenSwap{
		ID:            uuid.New().String(),
		Pubkey:        address,
		Mint:          tokenSwap.TokenPool.String(),
		Authority:     tokenLPMint.MintAuthority.String(),
		FeeAccount:    tokenSwap.FeeAccount.String(),
		TokenAMint:    tokenSwap.MintB.String(),
		TokenAAccount: tokenSwap.TokenAccountB.String(),
		TokenBMint:    tokenSwap.MintA.String(),
		TokenBAccount: tokenSwap.TokenAccountA.String(),
		TokenPairID:   tokenPairInverse.ID,
	}); err != nil {
		return err
	}

	// Upsert balances
	if err := p.UpsertTokenAccountBalanceByAddress(ctx, tokenSwap.TokenAccountA.String()); err != nil {
		return err
	}
	if err := p.UpsertTokenAccountBalanceByAddress(ctx, tokenSwap.TokenAccountB.String()); err != nil {
		return err
	}
	return nil
}

func (p impl) UpsertTokenAccountBalanceByAddress(ctx context.Context, address string) error {
	var tokenAccount token.Account
	if err := p.client.GetAccount(ctx, address, &tokenAccount); err != nil {
		return err
	}
	return p.UpsertTokenAccountBalance(ctx, address, tokenAccount)
}

func (p impl) UpsertTokenAccountBalance(ctx context.Context, address string, tokenAccount token.Account) error {
	isTokenSwapTokenAccount, _ := p.IsTokenSwapTokenAccount(ctx, address)
	isUserPositionNFTTokenAccount, _ := p.IsUserPositionTokenAccount(ctx, tokenAccount.Mint.String())
	if !isTokenSwapTokenAccount && !isUserPositionNFTTokenAccount {
		return nil
	}
	if isUserPositionNFTTokenAccount {
		logrus.
			WithField("mint", tokenAccount.Mint.String()).
			Info("recording user position token swap/creation")
	}
	state := "initialized"
	if tokenAccount.State == token.Uninitialized {
		state = "uninitialized"
	} else if tokenAccount.State == token.Frozen {
		state = "frozen"
	}

	var tokenMint token.Mint
	if err := p.client.GetAccount(ctx, tokenAccount.Mint.String(), &tokenMint); err != nil {
		return err
	}
	// TODO(Mocha): If this is a drip nft token, we can decorate the symbol with a deterministic name
	// with tokenA, tokenB, start, end
	tokenModel := model2.Token{
		Pubkey:   tokenAccount.Mint.String(),
		Symbol:   nil,
		Decimals: int16(tokenMint.Decimals),
		IconURL:  nil,
	}
	if err := p.repo.UpsertTokens(ctx, &tokenModel); err != nil {
		logrus.WithError(err).Error("failed to upsert tokens")
		return err
	}
	return p.repo.UpsertTokenAccountBalances(ctx, &model2.TokenAccountBalance{
		Pubkey: address,
		Mint:   tokenAccount.Mint.String(),
		Owner:  tokenAccount.Owner.String(),
		Amount: tokenAccount.Amount,
		State:  state,
	})
}

func (p impl) UpsertProtoConfigByAddress(ctx context.Context, address string) error {
	var protoConfig drip.VaultProtoConfig
	if err := p.client.GetAccount(ctx, address, &protoConfig); err != nil {
		return err
	}
	return p.repo.UpsertProtoConfigs(ctx, &model2.ProtoConfig{
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
	return p.repo.UpsertVaults(ctx, &model2.Vault{
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
	return p.repo.UpsertPositions(ctx, &model2.Position{
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
	return p.repo.UpsertVaultPeriods(ctx, &model2.VaultPeriod{
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
		&model2.Token{
			Pubkey:   tokenAAMint,
			Symbol:   nil,
			Decimals: int16(tokenA.Decimals),
			IconURL:  nil,
		}, &model2.Token{
			Pubkey:   tokenBMint,
			Symbol:   nil,
			Decimals: int16(tokenB.Decimals),
			IconURL:  nil,
		}); err != nil {
		return err
	}
	return p.repo.InsertTokenPairs(ctx, &model2.TokenPair{
		ID:     uuid.New().String(),
		TokenA: tokenAAMint,
		TokenB: tokenBMint,
	})
}

func (p impl) IsTokenSwapTokenAccount(ctx context.Context, tokenAccount string) (bool, error) {
	tokenSwap, err := p.repo.GetTokenSwapForTokenAccount(ctx, tokenAccount)
	if err != nil {
		return false, err
	}
	if tokenSwap == nil {
		return false, nil
	}
	return true, nil
}

func (p impl) IsUserPositionTokenAccount(ctx context.Context, mint string) (bool, error) {
	position, err := p.repo.GetPositionByNFTMint(ctx, mint)
	if err != nil {
		return false, err
	}
	if position == nil {
		return false, nil
	}
	return true, nil
}

// ensureTokenPair - if token pair exists return it, else upsert tokenPair and all needed tokenPair foreign keys
func (p impl) ensureTokenPair(ctx context.Context, tokenAAMint string, tokenBMint string) (*model2.TokenPair, error) {
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
func (p impl) ensureVault(ctx context.Context, address string) (*model2.Vault, error) {
	vault, err := p.repo.GetVaultByAddress(ctx, address)
	if err != nil && err.Error() == "record not found" {
		if err := p.UpsertVaultByAddress(ctx, address); err != nil {
			return nil, err
		}
		return p.repo.GetVaultByAddress(ctx, address)
	}
	return vault, err
}