package producer

import (
	"context"
	"encoding/json"
	"math/rand"
	"runtime/debug"
	"strconv"
	"time"

	solanaClient "github.com/dcaf-labs/drip/pkg/service/clients/solana"
	"github.com/dcaf-labs/drip/pkg/service/config"
	"github.com/dcaf-labs/drip/pkg/service/repository/model"
	queuerepository "github.com/dcaf-labs/drip/pkg/service/repository/queue"
	transactioncheckpointrepository "github.com/dcaf-labs/drip/pkg/service/repository/transactioncheckpoint"
	"github.com/dcaf-labs/drip/pkg/service/utils"
	drip "github.com/dcaf-labs/solana-drip-go/pkg/v1"
	"github.com/dcaf-labs/solana-go-clients/pkg/whirlpool"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

// TXPOLLFREQUENCY how often to foll for transactions
const TXPOLLFREQUENCY = 30 * time.Minute
const backfillEvery = time.Hour * 12

type DripProgramProducer struct {
	client                     solanaClient.Solana
	txProcessingCheckpointRepo transactioncheckpointrepository.TransactionProcessingCheckpointRepository
	accountUpdateQueue         queuerepository.AccountUpdateQueue
	transactionUpdateQueue     queuerepository.TransactionUpdateQueue
	cancel                     context.CancelFunc
	environment                config.Environment
	txBackfillStartSlot        uint64
	shouldBackfillDripAccounts bool
}

func Server(
	lifecycle fx.Lifecycle,
	client solanaClient.Solana,
	txProcessingCheckpointRepo transactioncheckpointrepository.TransactionProcessingCheckpointRepository,
	accountUpdateQueue queuerepository.AccountUpdateQueue,
	transactionUpdateQueue queuerepository.TransactionUpdateQueue,
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
		shouldBackfillDripAccounts: appConfig.GetShouldBackfillDripAccounts(),
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
	//if config.IsProductionEnvironment(d.environment) {
	// Track token_swap program accounts
	//if err := d.client.ProgramSubscribe(ctx, tokenswap.ProgramID.String(), d.processor.AddItemToAccountUpdateQueueCallback(ctx, tokenswap.ProgramID.String())); err != nil {
	//	return err
	//}

	// Track orca_whirlpool program accounts
	//if err := d.client.ProgramSubscribe(ctx, whirlpool.ProgramID.String(), d.AddItemToAccountUpdateQueueCallback(ctx, whirlpool.ProgramID.String())); err != nil {
	//	return err
	//}
	//}
	if d.shouldBackfillDripAccounts {
		go d.backfillAccounts(ctx)
	}
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
	logrus.Info("polling transactions")
	ticker := time.NewTicker(TXPOLLFREQUENCY)
	for {
		if err := d.processFromLastCheckpointSlot(ctx); err != nil {
			logrus.WithError(err).Error("failed to processFromLastCheckpointSlot, retrying in the next tick...")
		}
		select {
		case <-ticker.C:
			continue
		case <-ctx.Done():
			ticker.Stop()
			logrus.Info("context ended, exiting...")
			return
		}
	}
}

func (d *DripProgramProducer) processFromLastCheckpointSlot(ctx context.Context) error {
	checkpoint := d.txProcessingCheckpointRepo.GetLatestTransactionCheckpoint(ctx)
	var untilSignature solana.Signature
	var beforeSignature solana.Signature
	if checkpoint != nil {
		untilSignature = solana.MustSignatureFromBase58(checkpoint.Signature)
	}
	log := logrus.WithField("untilSignature", untilSignature.String())
	log.WithField("programId", drip.ProgramID.String()).Info("starting processing")
	defer func() {
		log.Info("done processing")
	}()
	// do while loop until txSignatures is not empty
	total := 0
	for {
		log = log.WithField("beforeSignature", beforeSignature.String())
		txSignatures, err := d.client.GetSignaturesForAddress(ctx, drip.ProgramID.String(), untilSignature, beforeSignature, nil)
		if err != nil {
			log.WithError(err).Error("failed to GetSignaturesForAddress")
			return err
		}
		total += len(txSignatures)
		log = log.WithField("len(txSignatures)", len(txSignatures))
		log.Info("got signatures")
		// insert oldest to newest
		txSignatures = lo.Reverse(txSignatures)
		for i := range txSignatures {
			txSignature := txSignatures[i]
			txPushLog := log.WithField("transactionSignature", txSignature.Signature.String())
			tx, err := d.client.GetTransaction(ctx, txSignature.Signature.String())
			if err != nil {
				txPushLog.WithError(err).Error("failed to GetTransaction")
				return err
			}
			if err := d.AddItemToTransactionUpdate(ctx, txSignature.Signature.String(), *tx); err != nil {
				txPushLog.WithError(err).Error("failed to insert data into queue...")
				return err
			}
			if err := d.txProcessingCheckpointRepo.UpsertTransactionProcessingCheckpoint(ctx, txSignature.Slot, txSignature.Signature.String()); err != nil {
				txPushLog.WithError(err).Error("failed to insert metadata...")
				return err
			}
		}
		if len(txSignatures) > 0 {
			beforeSignature = txSignatures[0].Signature
		} else {
			log = log.WithField("totalSignatures", total)
			return nil
		}
	}
}

func (d *DripProgramProducer) AddItemToAccountUpdateQueueCallback(ctx context.Context, programId string) func(string, []byte) error {
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

func (d *DripProgramProducer) AddItemToTransactionUpdate(ctx context.Context, signature string, tx rpc.GetTransactionResult) error {
	bytes, err := json.Marshal(tx)
	if err != nil {
		logrus.WithError(err).Error("failed to marshal tx")
		return err
	}
	priority := int32(1)
	if err := d.transactionUpdateQueue.AddTransactionUpdateQueueItem(ctx, &model.TransactionUpdateQueueItem{
		Signature: signature,
		TxJSON:    string(bytes),
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
