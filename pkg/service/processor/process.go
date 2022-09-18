package processor

import (
	"context"
	"fmt"
	"runtime/debug"
	"strconv"
	"time"

	bin "github.com/gagliardetto/binary"

	"github.com/dcaf-labs/drip/pkg/clients/solana"
	"github.com/dcaf-labs/drip/pkg/repository"
	"github.com/dcaf-labs/drip/pkg/repository/model"
	"github.com/dcaf-labs/solana-go-clients/pkg/drip"
	"github.com/dcaf-labs/solana-go-clients/pkg/tokenswap"
	"github.com/dcaf-labs/solana-go-clients/pkg/whirlpool"
	solana2 "github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/sirupsen/logrus"
)

type Processor interface {
	UpsertProtoConfigByAddress(context.Context, string) error
	UpsertVaultByAddress(context.Context, string) error
	UpsertPositionByAddress(context.Context, string) error
	UpsertPosition(context.Context, string, drip.Position) error
	UpsertVaultPeriodByAddress(context.Context, string) error
	UpsertTokenSwapByAddress(context.Context, string) error
	UpsertWhirlpoolByAddress(context.Context, string) error
	UpsertTokenPair(context.Context, string, string) error
	UpsertTokenAccountBalanceByAddress(context.Context, string) error
	UpsertTokenAccountBalance(context.Context, string, token.Account) error

	Backfill(ctx context.Context, programID string, processor func(string, []byte))
	ProcessDripEvent(address string, data []byte)
	ProcessTokenSwapEvent(address string, data []byte)
	ProcessWhirlpoolEvent(address string, data []byte)
	ProcessTokenEvent(address string, data []byte)
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

func (p impl) Backfill(ctx context.Context, programID string, processor func(string, []byte)) {
	log := logrus.WithField("program", programID).WithField("func", "Backfill")
	accounts, err := p.client.GetProgramAccounts(ctx, programID)
	if err != nil {
		log.WithError(err).Error("failed to get accounts")
	}
	page, pageSize, total := 0, 50, len(accounts)
	start, end := paginate(page, pageSize, total)
	for start < end {
		log = log.
			WithField("page", page).
			WithField("pageSize", pageSize).
			WithField("total", total)
		log.Infof("backfilling program accounts")
		err := p.client.GetAccounts(ctx, accounts[start:end], func(address string, data []byte) {
			processor(address, data)
		})
		if err != nil {
			log.WithError(err).
				Error("failed to get accounts")
		}
		page++
		start, end = paginate(page, pageSize, total)
	}
}

func (p impl) ProcessDripEvent(address string, data []byte) {
	if pubkey, err := solana2.PublicKeyFromBase58(address); err != nil || pubkey.IsZero() {
		logrus.WithField("address", address).Info("skipping zero/invalid address")
	}
	ctx := context.Background()
	log := logrus.WithField("address", address)
	// log.Infof("received drip account update")
	defer func() {
		if r := recover(); r != nil {
			log.WithField("stack", debug.Stack()).Errorf("panic in processEvent")
		}
	}()
	var vaultPeriod drip.VaultPeriod
	if err := bin.NewBinDecoder(data).Decode(&vaultPeriod); err == nil {
		// log.Infof("decoded as vaultPeriod")
		if err := p.UpsertVaultPeriodByAddress(ctx, address); err != nil {
			log.WithError(err).Error("failed to upsert vaultPeriod")
		}
		return
	}
	var position drip.Position
	if err := bin.NewBinDecoder(data).Decode(&position); err == nil {
		// log.Infof("decoded as position")
		if err := p.UpsertPosition(ctx, address, position); err != nil {
			log.WithError(err).Error("failed to upsert position")
		}
		return
	}
	var vault drip.Vault
	if err := bin.NewBinDecoder(data).Decode(&vault); err == nil {
		// log.Infof("decoded as vault")
		if err := p.UpsertVaultByAddress(ctx, address); err != nil {
			log.WithError(err).Error("failed to upsert vault")
		}
		return
	}
	var protoConfig drip.VaultProtoConfig
	if err := bin.NewBinDecoder(data).Decode(&protoConfig); err == nil {
		// log.Infof("decoded as protoConfig")
		if err := p.UpsertProtoConfigByAddress(ctx, address); err != nil {
			log.WithError(err).Error("failed to upsert protoConfig")
		}
		return
	}
	log.Errorf("failed to decode drip account to known types")
}

func (p impl) ProcessTokenSwapEvent(address string, data []byte) {
	if pubkey, err := solana2.PublicKeyFromBase58(address); err != nil || pubkey.IsZero() {
		logrus.WithField("address", address).Info("skipping zero/invalid address")
	}
	ctx := context.Background()
	log := logrus.WithField("address", address)
	defer func() {
		if r := recover(); r != nil {
			log.WithField("stack", debug.Stack()).Errorf("panic in processTokenSwapEvent")
		}
	}()
	var tokenSwap tokenswap.TokenSwap
	err := bin.NewBinDecoder(data).Decode(&tokenSwap)
	if err == nil {
		if err := p.UpsertTokenSwapByAddress(ctx, address); err != nil {
			log.WithError(err).Error("failed to upsert tokenSwap")
		}
		return
	}
	// log.WithError(err).Errorf("failed to decode token swap account")
}

func (p impl) ProcessWhirlpoolEvent(address string, data []byte) {
	if pubkey, err := solana2.PublicKeyFromBase58(address); err != nil || pubkey.IsZero() {
		logrus.WithField("address", address).Info("skipping zero/invalid address")
	}
	ctx := context.Background()
	log := logrus.WithField("address", address)
	defer func() {
		if r := recover(); r != nil {
			log.WithField("stack", debug.Stack()).Errorf("panic in processWhirlpoolEvent")
		}
	}()
	var whirlpoolAccount whirlpool.Whirlpool
	err := bin.NewBinDecoder(data).Decode(&whirlpoolAccount)
	if err == nil {
		if err := p.UpsertWhirlpoolByAddress(ctx, address); err != nil {
			log.WithError(err).Error("failed to upsert tokenSwap")
		}
		return
	}
}

func (p impl) ProcessTokenEvent(address string, data []byte) {
	if pubkey, err := solana2.PublicKeyFromBase58(address); err != nil || pubkey.IsZero() {
		logrus.WithField("address", address).Info("skipping zero/invalid address")
	}
	ctx := context.Background()
	log := logrus.WithField("address", address)
	defer func() {
		if r := recover(); r != nil {
			log.WithField("stack", debug.Stack()).Errorf("panic in processTokenEvent")
		}
	}()
	var tokenAccount token.Account
	err := bin.NewBinDecoder(data).Decode(&tokenAccount)
	if err == nil {
		if err := p.UpsertTokenAccountBalance(ctx, address, tokenAccount); err != nil {
			log.WithError(err).Error("failed to upsert tokenAccountBalance")
		}
		return
	}
}

func (p impl) UpsertWhirlpoolByAddress(ctx context.Context, address string) error {
	var orcaWhirlpool whirlpool.Whirlpool
	whirlpoolPubkey := solana2.MustPublicKeyFromBase58(address)
	if err := p.client.GetAccount(ctx, address, &orcaWhirlpool); err != nil {
		return err
	}

	tokenPair, err := p.ensureTokenPair(ctx, orcaWhirlpool.TokenMintA.String(), orcaWhirlpool.TokenMintB.String())
	if err != nil {
		return err
	}

	inverseTokenPair, err := p.ensureTokenPair(ctx, orcaWhirlpool.TokenMintB.String(), orcaWhirlpool.TokenMintA.String())
	if err != nil {
		return err
	}

	protocolFeeOwedA, _ := decimal.NewFromString(strconv.FormatUint(orcaWhirlpool.ProtocolFeeOwedA, 10))
	protocolFeeOwedB, _ := decimal.NewFromString(strconv.FormatUint(orcaWhirlpool.ProtocolFeeOwedB, 10))
	rewardLastUpdatedTimestamp, _ := decimal.NewFromString(strconv.FormatUint(orcaWhirlpool.RewardLastUpdatedTimestamp, 10))
	liquidity, _ := decimal.NewFromString(orcaWhirlpool.Liquidity.String())
	sqrtPrice, _ := decimal.NewFromString(orcaWhirlpool.SqrtPrice.String())
	feeGrowthGlobalA, _ := decimal.NewFromString(orcaWhirlpool.FeeGrowthGlobalA.String())
	feeGrowthGlobalB, _ := decimal.NewFromString(orcaWhirlpool.FeeGrowthGlobalB.String())
	oracle, _, _ := solana2.FindProgramAddress([][]byte{
		[]byte("oracle"),
		whirlpoolPubkey[:],
	}, whirlpool.ProgramID)

	if err := p.repo.UpsertOrcaWhirlpools(ctx,
		&model.OrcaWhirlpool{
			ID:                         uuid.New().String(),
			TokenPairID:                tokenPair.ID,
			Pubkey:                     whirlpoolPubkey.String(),
			WhirlpoolsConfig:           orcaWhirlpool.WhirlpoolsConfig.String(),
			TokenMintA:                 orcaWhirlpool.TokenMintA.String(),
			TokenVaultA:                orcaWhirlpool.TokenVaultA.String(),
			TokenMintB:                 orcaWhirlpool.TokenMintB.String(),
			TokenVaultB:                orcaWhirlpool.TokenVaultB.String(),
			TickSpacing:                int32(orcaWhirlpool.TickSpacing),
			FeeRate:                    int32(orcaWhirlpool.FeeRate),
			ProtocolFeeRate:            int32(orcaWhirlpool.ProtocolFeeRate),
			ProtocolFeeOwedA:           protocolFeeOwedA,
			ProtocolFeeOwedB:           protocolFeeOwedB,
			RewardLastUpdatedTimestamp: rewardLastUpdatedTimestamp,
			TickCurrentIndex:           orcaWhirlpool.TickCurrentIndex,
			Liquidity:                  liquidity,
			SqrtPrice:                  sqrtPrice,
			FeeGrowthGlobalA:           feeGrowthGlobalA,
			FeeGrowthGlobalB:           feeGrowthGlobalB,
			Oracle:                     oracle.String(),
		},
		// The only inverse is the token pair ID
		// For token swap it makes sense to inverse the mints, but for whirlpool it doesn't
		&model.OrcaWhirlpool{
			ID:                         uuid.New().String(),
			TokenPairID:                inverseTokenPair.ID,
			Pubkey:                     whirlpoolPubkey.String(),
			WhirlpoolsConfig:           orcaWhirlpool.WhirlpoolsConfig.String(),
			TokenMintA:                 orcaWhirlpool.TokenMintA.String(),
			TokenVaultA:                orcaWhirlpool.TokenVaultA.String(),
			TokenMintB:                 orcaWhirlpool.TokenMintB.String(),
			TokenVaultB:                orcaWhirlpool.TokenVaultB.String(),
			TickSpacing:                int32(orcaWhirlpool.TickSpacing),
			FeeRate:                    int32(orcaWhirlpool.FeeRate),
			ProtocolFeeRate:            int32(orcaWhirlpool.ProtocolFeeRate),
			ProtocolFeeOwedA:           protocolFeeOwedA,
			ProtocolFeeOwedB:           protocolFeeOwedB,
			RewardLastUpdatedTimestamp: rewardLastUpdatedTimestamp,
			TickCurrentIndex:           orcaWhirlpool.TickCurrentIndex,
			Liquidity:                  liquidity,
			SqrtPrice:                  sqrtPrice,
			FeeGrowthGlobalA:           feeGrowthGlobalA,
			FeeGrowthGlobalB:           feeGrowthGlobalB,
			Oracle:                     oracle.String(),
		},
	); err != nil {
		return err
	}

	if err := p.UpsertTokenAccountBalanceByAddress(ctx, orcaWhirlpool.TokenVaultA.String()); err != nil {
		return err
	}
	if err := p.UpsertTokenAccountBalanceByAddress(ctx, orcaWhirlpool.TokenVaultB.String()); err != nil {
		return err
	}
	return nil
}

func (p impl) UpsertTokenSwapByAddress(ctx context.Context, address string) error {
	var tokenSwap tokenswap.TokenSwap
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
	if err := p.repo.UpsertTokenSwaps(ctx, &model.TokenSwap{
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
	if err := p.repo.UpsertTokenSwaps(ctx, &model.TokenSwap{
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
	if !p.shouldIngestTokenBalance(ctx, tokenAccount) {
		return nil
	}
	state := "initialized"
	if tokenAccount.State == token.Uninitialized {
		state = "uninitialized"
	} else if tokenAccount.State == token.Frozen {
		state = "frozen"
	}

	var tokenMint token.Mint
	if err := p.client.GetAccount(ctx, tokenAccount.Mint.String(), &tokenMint); err != nil {
		return fmt.Errorf("failed to GetAccount %s, err: %w", tokenAccount.Mint.String(), err)
	}
	if err := p.UpsertTokenByAddress(ctx, tokenAccount.Mint.String()); err != nil {
		return fmt.Errorf("failed to UpsertTokenByAddress %s, err: %w", tokenAccount.Mint.String(), err)
	}
	return p.repo.UpsertTokenAccountBalances(ctx, &model.TokenAccountBalance{
		Pubkey: address,
		Mint:   tokenAccount.Mint.String(),
		Owner:  tokenAccount.Owner.String(),
		Amount: tokenAccount.Amount,
		State:  state,
	})
}

func (p impl) UpsertTokenByAddress(ctx context.Context, mintAddress string) error {
	tokenMint, err := p.client.GetTokenMint(ctx, mintAddress)
	if err != nil {
		return fmt.Errorf("failed to GetTokenMint %s, err: %w", mintAddress, err)
	}
	tokenMetadataAccount, err := p.client.GetTokenMetadataAccount(ctx, mintAddress)
	if err != nil {
		if err.Error() != solana.ErrNotFound {
			logrus.WithError(err).WithField("mint", mintAddress).Error("failed to GetTokenMetadataAccount for mint")
		}
		tokenModel := model.Token{
			Pubkey:   mintAddress,
			Decimals: int16(tokenMint.Decimals),
		}
		return p.repo.UpsertTokens(ctx, &tokenModel)
	}
	tokenModel := model.Token{
		Pubkey:   mintAddress,
		Symbol:   &tokenMetadataAccount.Data.Symbol,
		Decimals: int16(tokenMint.Decimals),
		IconURL:  nil, // TODO(Mocha): fetch the image via the URI if the URI is present
	}
	return p.repo.UpsertTokens(ctx, &tokenModel)
}

func (p impl) UpsertProtoConfigByAddress(ctx context.Context, address string) error {
	var protoConfig drip.VaultProtoConfig
	if err := p.client.GetAccount(ctx, address, &protoConfig); err != nil {
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

	if err := p.repo.UpsertVaults(ctx, &model.Vault{
		Pubkey:                 address,
		ProtoConfig:            vaultAccount.ProtoConfig.String(),
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
	if err := p.UpsertTokenAccountBalanceByAddress(ctx, vaultAccount.TokenAAccount.String()); err != nil {
		return err
	}
	if err := p.UpsertTokenAccountBalanceByAddress(ctx, vaultAccount.TokenBAccount.String()); err != nil {
		return err
	}
	if err := p.UpsertTokenAccountBalanceByAddress(ctx, vaultAccount.TreasuryTokenBAccount.String()); err != nil {
		return err
	}
	return nil
}

func (p impl) UpsertPositionByAddress(ctx context.Context, address string) error {
	var position drip.Position
	if err := p.client.GetAccount(ctx, address, &position); err != nil {
		return err
	}
	return p.UpsertPosition(ctx, address, position)
}

func (p impl) UpsertPosition(ctx context.Context, address string, position drip.Position) error {
	vault, err := p.ensureVault(ctx, position.Vault.String())
	if err != nil {
		return fmt.Errorf("failed to ensureVault, err: %w", err)
	}
	tokenPair, err := p.repo.GetTokenPairByID(ctx, vault.TokenPairID)
	if err != nil {
		return fmt.Errorf("failed to GetTokenPairByID, err: %w", err)
	}
	// Get up to date token metadata info
	if err := p.UpsertTokenByAddress(ctx, tokenPair.TokenA); err != nil {
		return fmt.Errorf("failed to UpsertTokenByAddress, err: %w", err)
	}
	if err := p.UpsertTokenByAddress(ctx, tokenPair.TokenB); err != nil {
		return err
	}
	if err := p.UpsertTokenByAddress(ctx, position.PositionAuthority.String()); err != nil {
		return err
	}
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
		logrus.WithError(err).Error("failed to UpsertPositions in UpsertPosition")
		return err
	}
	largestAccounts, err := p.client.GetLargestTokenAccounts(ctx, position.PositionAuthority.String())
	if err != nil {
		return err
	}
	for _, account := range largestAccounts {
		if account == nil {
			continue
		}
		if err := p.UpsertTokenAccountBalanceByAddress(ctx, account.Address.String()); err != nil {
			logrus.WithError(err).Error("failed to UpsertTokenAccountBalanceByAddress in UpsertPosition")
		}
	}
	return nil
}

func (p impl) UpsertVaultPeriodByAddress(ctx context.Context, address string) error {
	return p.upsertVaultPeriodByAddress(ctx, address, true)
}

// upsertVaultPeriodByAddress: this is potentially a recursive call
// if shouldUpsertPrice is set to true, we will try and price the period, which depends on period[i-1]
// if shouldUpsertPrice is set to false, we will not upsert the price to 0
func (p impl) upsertVaultPeriodByAddress(ctx context.Context, address string, shouldUpsertPrice bool) error {
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
//	p[i] = twap[i]*i - twap[i-1]*(i-1) for i > 0
//  p[i] = twap[i] for i = 0
func (p impl) getVaultPeriodPriceBOverA(ctx context.Context, periodI drip.VaultPeriod) (decimal.Decimal, error) {
	twapI, err := decimal.NewFromString(periodI.Twap.String())
	if err != nil {
		return decimal.Decimal{}, fmt.Errorf("failed to get twapI decimal, err: %w", err)
	}
	twapI = twapI.Shift(-64)
	if periodI.PeriodId == 0 {
		return twapI, nil
	}
	periodIPrecedingAddress, err := p.client.GetVaultPeriodPDA(periodI.Vault.String(), int64(periodI.PeriodId-1))
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
	twapIPreceding = twapIPreceding.Shift(-64)

	periodIID, err := decimal.NewFromString(strconv.FormatUint(periodI.PeriodId, 10))
	if err != nil {
		return decimal.Decimal{}, fmt.Errorf("failed to periodIId decimal, err: %w", err)
	}
	periodIPrecedingID := periodIID.Sub(decimal.NewFromInt(1))
	rawPrice := twapI.Mul(periodIID).Sub(twapIPreceding.Mul(periodIPrecedingID))
	tokenA, tokenB, err := p.getTokensForVault(ctx, periodI.Vault.String())
	if err != nil {
		return decimal.Decimal{}, err
	}
	return rawPrice.
		Mul(decimal.NewFromInt(10).Pow(decimal.NewFromInt(int64(tokenA.Decimals)))).
		DivRound(decimal.NewFromInt(10).Pow(decimal.NewFromInt(int64(tokenB.Decimals))), 128), nil
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
	if err := p.UpsertTokenByAddress(ctx, tokenAAMint); err != nil {
		return err
	}
	if err := p.UpsertTokenByAddress(ctx, tokenBMint); err != nil {
		return err
	}
	return p.repo.InsertTokenPairs(ctx, &model.TokenPair{
		ID:     uuid.New().String(),
		TokenA: tokenAAMint,
		TokenB: tokenBMint,
	})
}

// todo: this should prob live in the repo layer
func (p impl) getTokensForVault(ctx context.Context, vaultAddress string) (*model.Token, *model.Token, error) {
	vault, err := p.ensureVault(ctx, vaultAddress)
	if err != nil {
		return nil, nil, err
	}
	tokenPair, err := p.repo.GetTokenPairByID(ctx, vault.TokenPairID)
	if err != nil {
		return nil, nil, err
	}
	tokens, err := p.repo.GetTokensByMints(ctx, []string{tokenPair.TokenA, tokenPair.TokenB})
	if err != nil {
		return nil, nil, err
	}
	if len(tokens) != 2 {
		return nil, nil, fmt.Errorf("invalid number of tokens return for GetTokensByMints for id: %s", vault.TokenPairID)
	}
	if tokens[0].Pubkey == tokenPair.TokenA {
		return tokens[0], tokens[1], nil
	}
	return tokens[1], tokens[0], nil
}

func (p impl) shouldIngestTokenBalance(ctx context.Context, tokenAccount token.Account) bool {
	if p.IsTokenSwapTokenAccount(ctx, tokenAccount.Owner.String()) ||
		p.isOrcaWhirlpoolTokenAccount(ctx, tokenAccount.Owner.String()) ||
		p.isVaultTokenAccount(ctx, tokenAccount.Owner.String()) ||
		p.isUserPositionTokenAccount(ctx, tokenAccount.Mint.String()) {
		return true
	}
	return false
}

func (p impl) IsTokenSwapTokenAccount(ctx context.Context, tokenAccountOwner string) bool {
	_, err := p.repo.GetTokenSwapByAddress(ctx, tokenAccountOwner)
	if err != nil {
		return false
	}
	if err != nil {
		if err.Error() != repository.ErrRecordNotFound {
			logrus.WithError(err).Error("failed to query for token swap")
		}
		return false
	}
	return true
}

func (p impl) isOrcaWhirlpoolTokenAccount(ctx context.Context, tokenAccountOwner string) bool {
	_, err := p.repo.GetOrcaWhirlpoolByAddress(ctx, tokenAccountOwner)
	if err != nil {
		if err.Error() != repository.ErrRecordNotFound {
			logrus.WithError(err).Error("failed to query for whirlpool")
		}
		return false
	}
	return true
}

func (p impl) isVaultTokenAccount(ctx context.Context, tokenAccountOwner string) bool {
	_, err := p.repo.AdminGetVaultByAddress(ctx, tokenAccountOwner)
	if err != nil {
		if err.Error() != repository.ErrRecordNotFound {
			logrus.WithError(err).Error("failed to query for vault")
		}
		return false
	}
	return true
}

func (p impl) isUserPositionTokenAccount(ctx context.Context, mint string) bool {
	_, err := p.repo.GetPositionByNFTMint(ctx, mint)
	if err != nil {
		if err.Error() != repository.ErrRecordNotFound {
			logrus.WithError(err).Error("failed to query for position")
		}
		return false
	}
	return true
}

// ensureTokenPair - if token pair exists return it, else upsert tokenPair and all needed tokenPair foreign keys
func (p impl) ensureTokenPair(ctx context.Context, tokenAAMint string, tokenBMint string) (*model.TokenPair, error) {
	tokenPair, err := p.repo.GetTokenPair(ctx, tokenAAMint, tokenBMint)
	if err != nil && err.Error() == repository.ErrRecordNotFound {
		if err := p.UpsertTokenPair(ctx, tokenAAMint, tokenBMint); err != nil {
			return nil, err
		}
		return p.repo.GetTokenPair(ctx, tokenAAMint, tokenBMint)
	}

	return tokenPair, err
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

func paginate(pageNum int, pageSize int, sliceLength int) (int, int) {
	start := pageNum * pageSize

	if start > sliceLength {
		start = sliceLength
	}

	end := start + pageSize
	if end > sliceLength {
		end = sliceLength
	}

	return start, end
}
