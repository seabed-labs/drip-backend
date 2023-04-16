package producer

import (
	"context"
	"encoding/json"
	"errors"
	"math"
	"math/rand"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/dcaf-labs/drip/pkg/service/clients/solana"
	"github.com/dcaf-labs/drip/pkg/service/config"
	"github.com/dcaf-labs/drip/pkg/service/repository"
	"github.com/dcaf-labs/drip/pkg/service/repository/model"
	"github.com/dcaf-labs/drip/pkg/service/utils"
	drip "github.com/dcaf-labs/solana-drip-go/pkg/v1"
	"github.com/dcaf-labs/solana-go-clients/pkg/whirlpool"
	solana2 "github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/gagliardetto/solana-go/rpc/jsonrpc"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

// Defined
const DEFAULTMAXRETRY = 3

var ExceededMaxRetries = errors.New("exceeded max retries")

const backfillEvery = time.Hour * 12

type DripProgramProducer struct {
	client                     solana.Solana
	txProcessingCheckpointRepo repository.TransactionProcessingCheckpointRepository
	accountUpdateQueue         repository.AccountUpdateQueue
	transactionUpdateQueue     repository.TransactionUpdateQueue
	cancel                     context.CancelFunc
	environment                config.Environment
	txBackfillStartSlot        uint64
}

func Server(
	lifecycle fx.Lifecycle,
	client solana.Solana,
	txProcessingCheckpointRepo repository.TransactionProcessingCheckpointRepository,
	accountUpdateQueue repository.AccountUpdateQueue,
	transactionUpdateQueue repository.TransactionUpdateQueue,
	appConfig config.AppConfig,
) error {
	ctx, cancel := context.WithCancel(context.Background())
	dripProgramProducer := DripProgramProducer{
		client:                     client,
		txProcessingCheckpointRepo: txProcessingCheckpointRepo,
		accountUpdateQueue:         accountUpdateQueue,
		transactionUpdateQueue:     transactionUpdateQueue,
		cancel:                     cancel,
		environment:                appConfig.GetEnvironment(),
		txBackfillStartSlot:        145626971,
	}
	lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			if err := dripProgramProducer.start(ctx); err != nil {
				return err
			}
			return nil
		},
		OnStop: func(_ context.Context) error {
			dripProgramProducer.stop()
			return nil
		},
	})
	return nil
}

func (d *DripProgramProducer) start(ctx context.Context) error {
	logrus.Info("starting producer")
	// Track Drip accounts
	// We track drip accounts live
	if err := d.client.ProgramSubscribe(ctx, drip.ProgramID.String(), d.AddItemToAccountUpdateQueueCallback(ctx, drip.ProgramID.String())); err != nil {
		return err
	}

	// In staging, we manually backfill tokenswaps and whirlpools so that we can limit the # of rows in the DB
	if config.IsProductionEnvironment(d.environment) {
		// Track token_swap program accounts
		//if err := d.client.ProgramSubscribe(ctx, tokenswap.ProgramID.String(), d.processor.AddItemToAccountUpdateQueueCallback(ctx, tokenswap.ProgramID.String())); err != nil {
		//	return err
		//}

		// Track orca_whirlpool program accounts
		if err := d.client.ProgramSubscribe(ctx, whirlpool.ProgramID.String(), d.AddItemToAccountUpdateQueueCallback(ctx, whirlpool.ProgramID.String())); err != nil {
			return err
		}
	}

	go d.backfillAccounts(ctx)
	go d.pollTransactions(ctx)
	return nil
}

func (d *DripProgramProducer) stop() {
	d.cancel()
}

func (d *DripProgramProducer) backfillAccounts(ctx context.Context) {
	if d.environment == config.StagingEnv {
		logrus.Infof("skipping backfill on staging")
		return
	}
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			logrus.WithField("stack", string(debug.Stack())).Errorf("panic in backfillAccounts")
		}
		time.AfterFunc(backfillEvery, func() {
			if ctx.Err() != nil {
				return
			}
			d.backfillAccounts(ctx)
		})
	}()
	id := strconv.FormatInt(int64(rand.Int()), 10)
	log := logrus.WithField("id", id)
	log.Info("starting backfill")
	d.BackfillProgramOwnedAccounts(ctx, id, drip.ProgramID.String())
	log.Info("done backfill")
}

func (d *DripProgramProducer) BackfillProgramOwnedAccounts(ctx context.Context, logId string, programID string) {
	log := logrus.WithField("id", logId).WithField("program", programID).WithField("func", "BackfillProgramOwnedAccounts")
	accounts, err := d.client.GetProgramAccounts(ctx, programID)
	if err != nil {
		log.WithError(err).Error("failed to get accounts")
	}
	page, pageSize, total := 0, 50, len(accounts)
	start, end := utils.Paginate(page, pageSize, total)
	for start < end {
		log = log.
			WithField("page", page).
			WithField("pageSize", pageSize).
			WithField("total", total)
		log.Infof("backfilling program accounts")
		for i := start; i < end; i++ {
			if err := d.accountUpdateQueue.AddAccountUpdateQueueItem(ctx, &model.AccountUpdateQueueItem{
				Pubkey:    accounts[i],
				ProgramID: programID,
				Time:      utils.GetTimePtr(time.Now()),
				// Hardcode priority to 2 so that we don't block live drip updates (priority 1)
				Priority: utils.GetInt32Ptr(2),
				Try:      0,
				MaxTry:   utils.GetInt32Ptr(10),
			}); err != nil {
				log.WithError(err).Error("failed to add backfill account to queue")
			}
		}
		page++
		start, end = utils.Paginate(page, pageSize, total)
	}
}

func (d *DripProgramProducer) pollTransactions(ctx context.Context) {
	slot, err := d.backfillCheckpointSlot(ctx)
	if err != nil {
		logrus.WithError(err).Error("failed to fetch transaction checkpoint, starting at checkpoint slot +1")
		slot = slot + 1
	}

	for {
		select {
		case <-time.After(100 * time.Millisecond):
			if err := d.processSlotWithRetry(ctx, slot, 0, DEFAULTMAXRETRY); err != nil {
				logrus.WithError(err).Error("failed to produce block with retry, skipping...")
			}
			slot = slot + 1
		case <-ctx.Done():
			logrus.Info("context ended, exiting...")
		}
	}
}

func (d *DripProgramProducer) backfillCheckpointSlot(ctx context.Context) (uint64, error) {
	slot, skipTillTxSignature := func() (uint64, string) {
		checkpoint := d.txProcessingCheckpointRepo.GetLatestTransactionCheckpoint(ctx)
		if checkpoint != nil {
			return checkpoint.Slot, checkpoint.Signature
		}
		return d.txBackfillStartSlot, ""
	}()
	log := logrus.WithField("slot", slot).WithField("skipTillTxSignature", skipTillTxSignature)
	if skipTillTxSignature == "" {
		log.Info("nothing to producer")
		return slot, nil
	}

	block, err := d.client.GetBlock(ctx, slot)
	if err != nil {
		log.WithError(err).Error("failed to get last checkpoint slot block")
		return 0, err
	}
	if block == nil {
		log.Info("block is nil (not confirmed), nothing to producer")
		return slot, nil
	}

	// skip all transactions until we find a transaction with the same signature as skipTillTxSignature
	shouldSkip := true
	transactions := lo.Filter(block.Transactions, func(tx rpc.TransactionWithMeta, idx int) bool {
		if !shouldSkip {
			return true
		}
		if tx.MustGetTransaction().Signatures[0].String() == skipTillTxSignature {
			shouldSkip = false
		}
		return false
	})
	if err := d.pushTransactionsToQueue(ctx, transactions, logrus.WithField("slot", slot), slot); err != nil {
		log.WithError(err).Error("failed to producer checkpoint slot")
		return 0, err
	}
	return slot + 1, nil
}

func (p DripProgramProducer) processSlotWithRetry(
	ctx context.Context,
	slot uint64,
	try int16,
	maxTry int16,
) error {
	log := logrus.WithField("slot", slot).WithField("try", try).WithField("maxTry", maxTry)
	if try > 0 && try < maxTry {
		sleepDuration := time.Duration(math.Pow(2, float64(try)) * float64(time.Second))
		log.WithField("sleepDuration", sleepDuration).Warning("sleeping before retry...")
		time.Sleep(sleepDuration)
		log.WithField("sleepDuration", sleepDuration).Warning("woke up and retrying...")
	}
	if try > maxTry {
		log.WithError(ExceededMaxRetries).Errorf("failed to processSlotWithRetry")
		return ExceededMaxRetries
	}
	log.Info("fetching block")
	block, err := p.client.GetBlock(ctx, slot)
	if err != nil {
		switch t := err.(type) {
		case *jsonrpc.RPCError:
			switch t.Code {
			case -32009:
				fallthrough
			case -32007:
				log.WithError(err).Error("skipping slot without retry...")
				return nil
			default:
				log.WithError(err).Warning("unhandled JSONRPCError, retrying...")
				return p.processSlotWithRetry(ctx, slot, try+1, maxTry)
			}
		default:
			log.WithError(err).Error("failed to get block, retrying...")
		}
	}
	if block == nil {
		log.Info("block is nil (not confirmed), retrying...")
		return p.processSlotWithRetry(ctx, slot, try+1, maxTry)
	}
	if err := p.pushTransactionsToQueue(ctx, block.Transactions, log, slot); err != nil {
		return err
	}
	log.Info("processed slot")
	return nil
}

func (d DripProgramProducer) pushTransactionsToQueue(
	ctx context.Context,
	transactions []rpc.TransactionWithMeta,
	log *logrus.Entry,
	slot uint64,
) error {
	txPushCount := 0
	// Filter out all tx's with errors
	transactions = lo.Filter(transactions, func(tx rpc.TransactionWithMeta, index int) bool {
		return tx.Meta.Err == nil
	})
	for i, tx := range transactions {
		// we need to populate the slot as it's not provided via the rpc api for get block
		tx.Slot = slot
		parsedTx := tx.MustGetTransaction()
		if len(parsedTx.Signatures) == 0 {
			log.WithField("txIndex", i).Warning("no signatures for transaction, skipping...")
			continue
		}
		// we should only push transactions into the Q if at least 1 ix contains a program defined in the config
		shouldPushTx := lo.SomeBy(parsedTx.Message.Instructions, func(ix solana2.CompiledInstruction) bool {
			programId, _ := parsedTx.Message.Account(ix.ProgramIDIndex)
			return drip.ProgramID.String() == programId.String()
		})
		txSignature := parsedTx.Signatures[0].String()
		if shouldPushTx {
			log.WithField("transactionSignature", txSignature).Info("pushing tx to queue...")
			bytes, err := json.Marshal(tx)
			if err != nil {
				log.WithError(err).Error("failed to marshal tx")
				return err
			}
			if err := d.AddItemToTransactionUpdate(ctx, txSignature, string(bytes)); err != nil {
				log.WithError(err).Error("failed to insert data into queue...")
				return err
			} else {
				log.WithField("transactionSignature", txSignature).Info("pushed tx to queue...")
				txPushCount = txPushCount + 1
			}
		}
		if err := d.txProcessingCheckpointRepo.UpsertTransactionProcessingCheckpoint(ctx, slot, txSignature); err != nil {
			log.WithError(err).Error("failed to insert metadata...")
			return err
		}
	}
	log.WithField("txPushCount", txPushCount).Info("finished pushing transactions")
	return nil
}

func (d DripProgramProducer) AddItemToAccountUpdateQueueCallback(ctx context.Context, programId string) func(string, []byte) error {
	return func(account string, data []byte) error {
		priority := int32(3)
		if programId == drip.ProgramID.String() {
			priority = 1
		} else if programId == whirlpool.ProgramID.String() || programId == token.ProgramID.String() {
			priority = 2
		}
		if err := d.accountUpdateQueue.AddAccountUpdateQueueItem(ctx, &model.AccountUpdateQueueItem{
			Pubkey:    account,
			ProgramID: programId,
			Time:      utils.GetTimePtr(time.Now()),
			Priority:  &priority,
			MaxTry:    utils.GetInt32Ptr(10),
		}); err != nil {
			logrus.
				WithError(err).
				WithField("programId", programId).
				WithField("account", account).
				Error("failed to add queue item")
			return err
		}
		return nil
	}
}

func (d DripProgramProducer) AddItemToTransactionUpdate(ctx context.Context, signature string, txJSON string) error {
	priority := int32(1)
	if err := d.transactionUpdateQueue.AddTransactionUpdateQueueItem(ctx, &model.TransactionUpdateQueueItem{
		Signature: signature,
		TxJSON:    txJSON,
		Time:      utils.GetTimePtr(time.Now()),
		Priority:  &priority,
		MaxTry:    utils.GetInt32Ptr(10),
	}); err != nil {
		logrus.
			WithError(err).
			WithField("signature", signature).
			Error("failed to add queue item")
		return err
	}
	return nil
}
