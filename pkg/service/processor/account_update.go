package processor

import (
	"context"
	"runtime/debug"
	"strconv"
	"sync"
	"time"

	"github.com/dcaf-labs/drip/pkg/service/repository/model"
	drip "github.com/dcaf-labs/solana-drip-go/pkg/v1"
	"github.com/dcaf-labs/solana-go-clients/pkg/tokenswap"
	"github.com/dcaf-labs/solana-go-clients/pkg/whirlpool"
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

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

	ticker := time.NewTicker(POLLFREQUENCY)
	for {
		queueItem, err := p.accountUpdateQueue.PopAccountUpdateQueueItem(ctx)
		if err != nil && err == gorm.ErrRecordNotFound {
		} else if err != nil {
			logrus.WithError(err).Error("failed to fetch account from queue")
		} else if queueItem == nil {
			logrus.WithError(err).Error("failed to get next queue item")
		} else {
			ch <- queueItem
		}
		select {
		case <-ticker.C:
			continue
		case <-ctx.Done():
			return
		}
	}
}

func (p impl) processAccountUpdateQueueItemWorker(ctx context.Context, id string, wg *sync.WaitGroup, queueCh chan *model.AccountUpdateQueueItem) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			logrus.Info("exiting processAccountUpdateQueueItemWorker")
			return
		case queueItem := <-queueCh:
			if err := p.ProcessAccount(ctx, queueItem.Pubkey, queueItem.ProgramID); err != nil {
				if requeueErr := p.accountUpdateQueue.ReQueueAccountUpdateQueueItem(ctx, queueItem, err.Error()); requeueErr != nil {
					logrus.WithField("pubkey", queueItem.Pubkey).WithError(requeueErr).Error("failed to add item to queue")
				}
			}
		}
	}
}

func (p impl) ProcessAccount(ctx context.Context, pubkey string, programId string) error {
	log := logrus.WithField("pubkey", pubkey).WithField("programId", programId)
	accountInfo, err := p.solanaClient.GetAccountInfo(ctx, pubkey)
	if err != nil || accountInfo == nil || accountInfo.Value == nil || accountInfo.Value.Data == nil {
		log.WithError(err).Error("error or invalid account")
		return err
	}
	switch programId {
	case drip.ProgramID.String():
		err = p.processDripEvent(pubkey, accountInfo.Value.Data.GetBinary())
	case whirlpool.ProgramID.String():
		err = p.processWhirlpoolEvent(pubkey, accountInfo.Value.Data.GetBinary())
	case tokenswap.ProgramID.String():
		err = p.processTokenSwapEvent(pubkey, accountInfo.Value.Data.GetBinary())
	default:
		log.Error("invalid programID")
	}
	return err
}

func (p impl) processDripEvent(address string, data []byte) error {
	if pubkey, err := solana.PublicKeyFromBase58(address); err != nil || pubkey.IsZero() {
		logrus.WithField("address", address).Info("skipping zero/invalid address")
		return nil
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
			return err
		}
		return nil
	}
	var position drip.Position
	if err := bin.NewBinDecoder(data).Decode(&position); err == nil {
		// log.Infof("decoded as position")
		if err := p.UpsertPosition(ctx, address, position); err != nil {
			log.WithError(err).Error("failed to upsert position")
			return err
		}
		return nil
	}
	var vault drip.Vault
	if err := bin.NewBinDecoder(data).Decode(&vault); err == nil {
		// log.Infof("decoded as vault")
		if err := p.UpsertVaultByAddress(ctx, address); err != nil {
			log.WithError(err).Error("failed to upsert vault")
			return err
		}
		return nil
	}
	var protoConfig drip.VaultProtoConfig
	if err := bin.NewBinDecoder(data).Decode(&protoConfig); err == nil {
		// log.Infof("decoded as protoConfig")
		if err := p.UpsertProtoConfigByAddress(ctx, address); err != nil {
			log.WithError(err).Error("failed to upsert protoConfig")
			return err
		}
		return nil
	}
	log.Errorf("failed to decode drip account to known types")
	return nil
}

func (p impl) processTokenSwapEvent(address string, data []byte) error {
	if pubkey, err := solana.PublicKeyFromBase58(address); err != nil || pubkey.IsZero() {
		logrus.WithField("address", address).Info("skipping zero/invalid address")
		return nil
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
			return err
		}
		return nil
	}
	return nil
	// log.WithError(err).Errorf("failed to decode token swap account")
}

func (p impl) processWhirlpoolEvent(address string, data []byte) error {
	if pubkey, err := solana.PublicKeyFromBase58(address); err != nil || pubkey.IsZero() {
		logrus.WithField("address", address).Info("skipping zero/invalid address")
		return nil
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
			log.WithError(err).Error("failed to upsert whirlpool")
			return err
		}
		return nil
	}
	return nil
}
