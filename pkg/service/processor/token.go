package processor

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/dcaf-labs/drip/pkg/service/clients/coingecko"
	"github.com/dcaf-labs/drip/pkg/service/repository"
	"github.com/dcaf-labs/drip/pkg/service/repository/model"
	"github.com/dcaf-labs/drip/pkg/service/utils"
	bin "github.com/gagliardetto/binary"
	tokenmetadata "github.com/gagliardetto/metaplex-go/clients/token-metadata"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
)

func (p impl) UpsertTokenPair(ctx context.Context, tokenAAMint string, tokenBMint string) error {
	var tokenA token.Mint
	if err := p.solanaClient.GetAccount(ctx, tokenAAMint, &tokenA); err != nil {
		return err
	}
	var tokenB token.Mint
	if err := p.solanaClient.GetAccount(ctx, tokenBMint, &tokenB); err != nil {
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

func (p impl) UpsertTokenByAddress(ctx context.Context, mintAddress string) error {
	return p.UpsertTokensByAddresses(ctx, mintAddress)
}

func (p impl) UpsertTokensByAddresses(ctx context.Context, addresses ...string) error {
	return utils.DoForPaginatedBatch(coingecko.CoinsMarketsPathLimit, len(addresses),
		func(start, end int) error { return p.upsertTokensByAddresses(ctx, addresses[start:end]...) },
		func(err error) error { return err },
	)
}

// process batch
func (p impl) upsertTokensByAddresses(ctx context.Context, addresses ...string) error {
	tokenMintsByAddress := make(map[string]token.Mint)
	var tokenMints []token.Mint
	if err := p.solanaClient.GetAccounts(ctx, addresses, func(address string, data []byte) {
		var tokenMint token.Mint
		if err := bin.NewBinDecoder(data).Decode(&tokenMint); err != nil {
			logrus.
				WithError(err).
				WithField("address", address).
				Errorf("failed to decode tokenMint")
		} else {
			tokenMints = append(tokenMints, tokenMint)
			tokenMintsByAddress[address] = tokenMint
		}
	}); err != nil {
		return err
	}

	metadataAddresses := lo.FilterMap[string, string](addresses, func(mintAddress string, _ int) (string, bool) {
		tokenMetadataAddress, err := utils.GetTokenMetadataPDA(mintAddress)
		if err != nil {
			logrus.
				WithError(err).
				WithField("address", mintAddress).
				Errorf("failed to GetTokenMetadataPDA")
			return "", false
		}
		return tokenMetadataAddress, true
	})

	tokenMetadataByAddress := make(map[string]tokenmetadata.Metadata)
	var tokenMetadatas []tokenmetadata.Metadata
	if err := p.solanaClient.GetAccounts(ctx, metadataAddresses, func(address string, data []byte) {
		var tokenMetadata tokenmetadata.Metadata
		if err := bin.NewBorshDecoder(data).Decode(&tokenMetadata); err != nil {
			logrus.
				WithError(err).
				WithField("address", address).
				Errorf("failed to decode tokenMetadata")
		} else {
			tokenMetadatas = append(tokenMetadatas, tokenMetadata)
			tokenMetadataByAddress[address] = tokenMetadata
		}
	}); err != nil {
		logrus.WithError(err).Info("failed to GetAccounts for token metadata accounts, continuing...")
	}

	// create base
	tokensByAddress, err := func() (map[string]*model.Token, error) {
		tokens, err := p.repo.GetTokensByAddresses(ctx, addresses...)
		if err != nil {
			return nil, err
		}
		return lo.KeyBy[string, *model.Token](tokens, func(tokenModel *model.Token) string {
			return tokenModel.Pubkey
		}), nil
	}() //nolint:errcheck
	if err != nil {
		return err
	}
	// add new tokens
	for address, tokenMint := range tokenMintsByAddress {
		if _, ok := tokensByAddress[address]; ok {
			continue
		}
		tokensByAddress[address] = &model.Token{
			Pubkey:   address,
			Decimals: int16(tokenMint.Decimals),
		}
	}
	// populate via token metadata
	for address := range tokensByAddress {
		tokenModel := tokensByAddress[address]
		tokenMetadataAddress, _ := utils.GetTokenMetadataPDA(address)
		if tokenMetadata, ok := tokenMetadataByAddress[tokenMetadataAddress]; ok {
			tokenModel.Symbol = utils.GetStringPtr(strings.Trim(tokenMetadata.Data.Symbol, "\u0000"))
			tokenModel.Name = utils.GetStringPtr(strings.Trim(tokenMetadata.Data.Name, "\u0000"))
		}
		tokensByAddress[address] = tokenModel
	}
	// populate via coin gecko metadata
	cgCoinsList := func() coingecko.CoinsListResponse {
		cgCoinsList, err := p.coinGeckoClient.GetSolanaCoinsList(ctx)
		if err != nil {
			logrus.WithError(err).Error("failed to get GetSolanaCoinsList")
			return nil
		}
		return cgCoinsList
	}()
	cgCoinsByAddress := func() map[string]coingecko.CoinResponse {
		cgCoinsByAddress := lo.KeyBy[string, coingecko.CoinResponse](cgCoinsList, func(coin coingecko.CoinResponse) string {
			return *coin.Platforms.Solana
		})
		return cgCoinsByAddress
	}()
	// cg 1: populate cg token ID
	var coingeckoIDs []string
	for address := range tokensByAddress {
		token := tokensByAddress[address]
		cgMetdata, ok := cgCoinsByAddress[address]
		if !ok {
			continue
		}
		token.CoinGeckoID = utils.GetStringPtr(cgMetdata.ID)
		coingeckoIDs = append(coingeckoIDs, cgMetdata.ID)
		if token.Symbol == nil || *token.Symbol == "" {
			token.Symbol = utils.GetStringPtr(cgMetdata.Symbol)
		}
		if token.Name == nil || *token.Name == "" {
			token.Name = utils.GetStringPtr(cgMetdata.Name)
		}
		tokensByAddress[address] = token
	}
	// Sort to make the api-request deterministic, makes for less flaky tests via the api replay
	sort.Strings(coingeckoIDs)
	// cg 2: populate market metadata
	cgTokenMarketDataByCGID := func() map[string]coingecko.CoinGeckoTokenMarketPriceResponse {
		tokenPrices, err := p.coinGeckoClient.GetMarketPriceForTokens(ctx, coingeckoIDs...)
		if err != nil {
			logrus.WithError(err).Error("failed to get GetMarketPriceForTokens")
			return nil
		}
		res := make(map[string]coingecko.CoinGeckoTokenMarketPriceResponse)
		for _, tokenPrice := range tokenPrices {
			res[tokenPrice.ID] = tokenPrice
		}
		return res
	}()
	for address := range tokensByAddress {
		token := tokensByAddress[address]
		if token.CoinGeckoID == nil || *token.CoinGeckoID == "" {
			continue
		}
		if marketData, ok := cgTokenMarketDataByCGID[*token.CoinGeckoID]; ok {
			token.MarketCapRank = marketData.MarketCapRank
			token.UIMarketPriceUsd = utils.GetFloat64Ptr(marketData.CurrentPrice)
		}
		tokensByAddress[address] = token
	}
	// note: this makes n network calls only if iconURL doesn't exist
	// todo: maybe this should live in a separate 1-time back-fill script?
	for address := range tokensByAddress {
		token := tokensByAddress[address]
		if token.IconURL != nil && *token.IconURL != "" {
			continue
		}
		if coinGeckoMeta, err := p.coinGeckoClient.GetCoinGeckoMetadata(ctx, token.Pubkey); err != nil {
			logrus.WithError(err).Error("failed to GetCoinGeckoMetadata")
		} else {
			token.IconURL = coinGeckoMeta.Image.Small
		}
		tokensByAddress[address] = token
	}
	return p.repo.UpsertTokens(ctx, lo.Values[string, *model.Token](tokensByAddress)...)
}

func (p impl) UpsertTokenAccountsByAddresses(ctx context.Context, addresses ...string) error {
	return utils.DoForPaginatedBatch(50, len(addresses),
		func(start, end int) error { return p.upsertTokenAccountsByAddresses(ctx, addresses[start:end]...) },
		func(err error) error { return err },
	)
}

// process batch
func (p impl) upsertTokenAccountsByAddresses(ctx context.Context, addresses ...string) error {
	if len(addresses) == 0 {
		return nil
	}
	//tokenAccounts := make(map[string]token.Mint)
	var tokenAccountModels []*model.TokenAccount
	if err := p.solanaClient.GetAccounts(ctx, addresses, func(address string, data []byte) {
		var tokenAccount token.Account
		if err := bin.NewBinDecoder(data).Decode(&tokenAccount); err != nil {
			logrus.
				WithError(err).
				WithField("address", address).
				Errorf("failed to decode tokenAccount")
		} else {
			tokenAccountModels = append(tokenAccountModels, &model.TokenAccount{
				Pubkey: address,
				Mint:   tokenAccount.Mint.String(),
				Owner:  tokenAccount.Owner.String(),
				Amount: tokenAccount.Amount,
				State:  getTokenAccountState(tokenAccount.State),
			})
		}
	}); err != nil {
		return err
	}
	if len(tokenAccountModels) == 0 {
		return nil
	}
	return p.repo.UpsertTokenAccounts(ctx, tokenAccountModels...)
}

func (p impl) UpsertTokenAccount(ctx context.Context, address string, tokenAccount token.Account) error {
	state := getTokenAccountState(tokenAccount.State)

	var tokenMint token.Mint
	if err := p.solanaClient.GetAccount(ctx, tokenAccount.Mint.String(), &tokenMint); err != nil {
		return fmt.Errorf("failed to GetAccount %s, err: %w", tokenAccount.Mint.String(), err)
	}
	if err := p.UpsertTokenByAddress(ctx, tokenAccount.Mint.String()); err != nil {
		return fmt.Errorf("failed to UpsertTokenByAddress %s, err: %w", tokenAccount.Mint.String(), err)
	}
	return p.repo.UpsertTokenAccounts(ctx, &model.TokenAccount{
		Pubkey: address,
		Mint:   tokenAccount.Mint.String(),
		Owner:  tokenAccount.Owner.String(),
		Amount: tokenAccount.Amount,
		State:  state,
	})
}

func getTokenAccountState(state token.AccountState) string {
	if state == token.Uninitialized {
		return "uninitialized"
	} else if state == token.Frozen {
		return "frozen"
	}
	return "initialized"
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
