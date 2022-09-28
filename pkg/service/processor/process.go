package processor

import (
	"context"
	"runtime/debug"
	"strconv"
	"sync"
	"time"

	"gorm.io/gorm"

	"github.com/dcaf-labs/drip/pkg/service/utils"

	"github.com/dcaf-labs/drip/pkg/service/repository/model"

	"github.com/dcaf-labs/drip/pkg/service/alert"
	"github.com/dcaf-labs/drip/pkg/service/clients/solana"
	"github.com/dcaf-labs/drip/pkg/service/clients/tokenregistry"
	"github.com/dcaf-labs/drip/pkg/service/repository"
	"github.com/dcaf-labs/solana-go-clients/pkg/drip"
	"github.com/dcaf-labs/solana-go-clients/pkg/tokenswap"
	"github.com/dcaf-labs/solana-go-clients/pkg/whirlpool"
	bin "github.com/gagliardetto/binary"
	solana2 "github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/sirupsen/logrus"
)

const processConcurrency = 100

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

	BackfillProgramOwnedAccounts(ctx context.Context, logId string, programID string)
	AddItemToUpdateQueueCallback(ctx context.Context, programId string) func(string, []byte)
	ProcessAccountUpdateQueue(ctx context.Context)
	ProcessDripEvent(address string, data []byte)
	ProcessTokenSwapEvent(address string, data []byte)
	ProcessWhirlpoolEvent(address string, data []byte)
	ProcessTokenEvent(address string, data []byte)
}

type impl struct {
	repo                repository.Repository
	accountUpdateQueue  repository.AccountUpdateQueue
	tokenRegistryClient tokenregistry.TokenRegistry
	solanaClient        solana.Solana
	alertService        alert.Service
}

func NewProcessor(
	repo repository.Repository,
	accountUpdateQueue repository.AccountUpdateQueue,
	client solana.Solana,
	tokenRegistryClient tokenregistry.TokenRegistry,
	alertService alert.Service,
) Processor {
	return impl{
		repo:                repo,
		accountUpdateQueue:  accountUpdateQueue,
		solanaClient:        client,
		tokenRegistryClient: tokenRegistryClient,
		alertService:        alertService,
	}
}

func (p impl) BackfillProgramOwnedAccounts(ctx context.Context, logId string, programID string) {
	log := logrus.WithField("id", logId).WithField("program", programID).WithField("func", "BackfillProgramOwnedAccounts")
	accounts, err := p.solanaClient.GetProgramAccounts(ctx, programID)
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
		for i := start; i < end; i++ {
			if err := p.accountUpdateQueue.AddItem(ctx, &model.AccountUpdateQueueItem{
				Pubkey:    accounts[i],
				ProgramID: programID,
				Time:      utils.GetTimePtr(time.Now()),
				Priority:  utils.GetIntPtr(3),
			}); err != nil {
				log.WithError(err).Error("failed to add backfill account to queue")
			}
		}
		page++
		start, end = paginate(page, pageSize, total)
	}
}

func (p impl) AddItemToUpdateQueueCallback(ctx context.Context, programId string) func(string, []byte) {
	return func(account string, data []byte) {
		priority := int32(3)
		if programId == drip.ProgramID.String() {
			priority = 1
		} else if programId == whirlpool.ProgramID.String() || programId == token.ProgramID.String() {
			priority = 2
		}
		if err := p.accountUpdateQueue.AddItem(ctx, &model.AccountUpdateQueueItem{
			Pubkey:    account,
			ProgramID: programId,
			Time:      utils.GetTimePtr(time.Now()),
			Priority:  &priority,
		}); err != nil {
			logrus.
				WithError(err).
				WithField("programId", programId).
				WithField("account", account).
				Error("failed to add queue item")
		}
	}
}

func (p impl) ProcessAccountUpdateQueue(ctx context.Context) {
	var wg sync.WaitGroup
	ch := make(chan *model.AccountUpdateQueueItem)
	defer func() {
		close(ch)
		logrus.Info("exiting ProcessAccountUpdateQueue...")
		wg.Wait()
	}()

	for i := 0; i < processConcurrency; i++ {
		wg.Add(1)
		go p.processAccountUpdateQueueItemWorker(ctx, strconv.FormatInt(int64(i), 10), &wg, ch)
	}

	for {
		select {
		case <-ctx.Done():
			return
		default:
			queueItem, err := p.accountUpdateQueue.Pop(ctx)
			if err != nil && err == gorm.ErrRecordNotFound {
				continue
			} else if queueItem == nil {
				logrus.WithError(err).Error("failed to get next queue item")
				continue
			}
			//if depth, err := p.accountUpdateQueue.Depth(ctx); err != nil {
			//	logrus.WithError(err).Error("failed to get queue depth")
			//} else {
			//	logrus.WithField("depth", depth).Infof("queue depth")
			//}
			ch <- queueItem
		}
	}
}

func (p impl) processAccountUpdateQueueItemWorker(ctx context.Context, id string, wg *sync.WaitGroup, queueCh chan *model.AccountUpdateQueueItem) {
	defer wg.Done()
	logrus.Info("spawned processAccountUpdateQueueItemWorker goroutine")
	for {
		select {
		case <-ctx.Done():
			logrus.Info("exiting processAccountUpdateQueueItemWorker")
			return
		case queueItem := <-queueCh:
			p.processAccountUpdateQueueItem(ctx, id, queueItem)
		}
	}

}

func (p impl) processAccountUpdateQueueItem(ctx context.Context, id string, queueItem *model.AccountUpdateQueueItem) {
	log := logrus.WithField("id", id).WithField("pubkey", queueItem.Pubkey).WithField("programId", queueItem.ProgramID)
	shouldReQeue := false
	defer func() {
		if !shouldReQeue {
			return
		}
		if err := p.accountUpdateQueue.AddItem(ctx, queueItem); err != nil {
			log.WithError(err).Error("failed to add item to queue")
		}
	}()

	accountInfo, err := p.solanaClient.GetAccountInfo(ctx, queueItem.Pubkey)
	if err != nil && err.Error() == solana.ErrNotFound {
		log.WithError(err).Error("failed to get accountInfo")
		return
	} else if err != nil {
		log.WithError(err).Error("failed to get accountInfo, re-queueing")
		shouldReQeue = true
		return
	} else if accountInfo == nil || accountInfo.Value == nil || accountInfo.Value.Data == nil {
		log.Info("account is empty")
		return
	}
	switch queueItem.ProgramID {
	case drip.ProgramID.String():
		p.ProcessDripEvent(queueItem.Pubkey, accountInfo.Value.Data.GetBinary())
	case whirlpool.ProgramID.String():
		p.ProcessWhirlpoolEvent(queueItem.Pubkey, accountInfo.Value.Data.GetBinary())
	case tokenswap.ProgramID.String():
		p.ProcessTokenSwapEvent(queueItem.Pubkey, accountInfo.Value.Data.GetBinary())
	case token.ProgramID.String():
		p.ProcessTokenEvent(queueItem.Pubkey, accountInfo.Value.Data.GetBinary())
	default:
		logrus.WithField("programId", queueItem.ProgramID).Error("invalid programID")
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
			debug.PrintStack()
			log.WithField("stack", string(debug.Stack())).Errorf("panic in processEvent")
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
