package processor

import (
	"context"
	"time"

	"github.com/dcaf-labs/drip/pkg/service/alert"
	"github.com/dcaf-labs/drip/pkg/service/clients/coingecko"
	"github.com/dcaf-labs/drip/pkg/service/clients/orcawhirlpool"
	"github.com/dcaf-labs/drip/pkg/service/clients/solana"
	"github.com/dcaf-labs/drip/pkg/service/clients/tokenregistry"
	"github.com/dcaf-labs/drip/pkg/service/ixparser"
	"github.com/dcaf-labs/drip/pkg/service/repository"
	queuerepository "github.com/dcaf-labs/drip/pkg/service/repository/queue"
	"github.com/gagliardetto/solana-go/rpc"
)

const processConcurrency = 10
const POLLFREQUENCY = 1 * time.Second

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

	ProcessAccountUpdateQueue(ctx context.Context)
	ProcessAccount(ctx context.Context, pubkey string, programId string) error

	ProcessTransactionUpdateQueue(ctx context.Context)
	ProcessTransaction(ctx context.Context, txRaw rpc.GetTransactionResult) error
}

type impl struct {
	repo                   repository.Repository
	accountUpdateQueue     queuerepository.AccountUpdateQueue
	transactionUpdateQueue queuerepository.TransactionUpdateQueue
	tokenRegistryClient    tokenregistry.TokenRegistry
	orcaWhirlpoolClient    orcawhirlpool.OrcaWhirlpoolClient
	solanaClient           solana.Solana
	ixParser               ixparser.IxParser
	alertService           alert.Service
	coinGeckoClient        coingecko.CoinGeckoClient
}

func NewProcessor(
	repo repository.Repository,
	accountUpdateQueue queuerepository.AccountUpdateQueue,
	transactionUpdateQueue queuerepository.TransactionUpdateQueue,
	client solana.Solana,
	tokenRegistryClient tokenregistry.TokenRegistry,
	orcaWhirlpoolClient orcawhirlpool.OrcaWhirlpoolClient,
	ixParser ixparser.IxParser,
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
		ixParser:               ixParser,
		alertService:           alertService,
		coinGeckoClient:        coinGeckoClient,
	}
}
