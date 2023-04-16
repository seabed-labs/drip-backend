package processor

import (
	"context"

	"github.com/dcaf-labs/drip/pkg/service/alert"
	"github.com/dcaf-labs/drip/pkg/service/clients/coingecko"
	"github.com/dcaf-labs/drip/pkg/service/clients/orcawhirlpool"
	"github.com/dcaf-labs/drip/pkg/service/clients/solana"
	"github.com/dcaf-labs/drip/pkg/service/clients/tokenregistry"
	"github.com/dcaf-labs/drip/pkg/service/ixparser"
	"github.com/dcaf-labs/drip/pkg/service/repository"
	drip "github.com/dcaf-labs/solana-drip-go/pkg/v1"
	"github.com/gagliardetto/solana-go/programs/token"
)

const processConcurrency = 10

type Processor interface {
	UpsertProtoConfigByAddress(context.Context, string) error
	UpsertVaultPeriodByAddress(context.Context, string) error
	UpsertVaultByAddress(context.Context, string) error
	UpsertPositionByAddress(context.Context, string) error
	UpsertTokenSwapByAddress(context.Context, string) error
	UpsertWhirlpoolByAddress(context.Context, string) error
	UpsertTokenByAddress(ctx context.Context, mintAddress string) error

	UpsertTokensByAddresses(ctx context.Context, addresses ...string) error
	UpsertTokenAccountsByAddresses(ctx context.Context, addresses ...string) error

	UpsertPosition(context.Context, string, drip.Position) error
	UpsertTokenPair(context.Context, string, string) error
	UpsertTokenAccount(context.Context, string, token.Account) error

	ProcessAccountUpdateQueue(ctx context.Context)
	ProcessDripEvent(address string, data []byte) error
	ProcessTokenSwapEvent(address string, data []byte) error
	ProcessWhirlpoolEvent(address string, data []byte) error

	ProcessTransactionUpateQueue(ctx context.Context)
}

type impl struct {
	repo                   repository.Repository
	accountUpdateQueue     repository.AccountUpdateQueue
	transactionUpdateQueue repository.TransactionUpdateQueue
	tokenRegistryClient    tokenregistry.TokenRegistry
	orcaWhirlpoolClient    orcawhirlpool.OrcaWhirlpoolClient
	solanaClient           solana.Solana
	dripIxParser           ixparser.IxParser
	alertService           alert.Service
	coinGeckoClient        coingecko.CoinGeckoClient
}

func NewProcessor(
	repo repository.Repository,
	accountUpdateQueue repository.AccountUpdateQueue,
	transactionUpdateQueue repository.TransactionUpdateQueue,
	client solana.Solana,
	tokenRegistryClient tokenregistry.TokenRegistry,
	orcaWhirlpoolClient orcawhirlpool.OrcaWhirlpoolClient,
	dripIxParser ixparser.IxParser,
	alertService alert.Service,
	coinGeckoClient coingecko.CoinGeckoClient,
) Processor {
	return impl{
		repo:                   repo,
		accountUpdateQueue:     accountUpdateQueue,
		transactionUpdateQueue: transactionUpdateQueue,
		solanaClient:           client,
		tokenRegistryClient:    tokenRegistryClient,
		orcaWhirlpoolClient:    orcaWhirlpoolClient,
		dripIxParser:           dripIxParser,
		alertService:           alertService,
		coinGeckoClient:        coinGeckoClient,
	}
}
