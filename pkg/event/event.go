package event

import (
	"context"
	"runtime/debug"

	"github.com/dcaf-labs/drip/pkg/clients/solana"
	"github.com/dcaf-labs/drip/pkg/configs"
	"github.com/dcaf-labs/drip/pkg/service/processor"
	"github.com/dcaf-labs/solana-go-clients/pkg/drip"
	"github.com/dcaf-labs/solana-go-clients/pkg/tokenswap"
	"github.com/dcaf-labs/solana-go-clients/pkg/whirlpool"
	bin "github.com/gagliardetto/binary"
	solana2 "github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

type DripProgramProcessor struct {
	client      solana.Solana
	processor   processor.Processor
	cancel      context.CancelFunc
	environment configs.Environment
}

func Server(
	lifecycle fx.Lifecycle,
	client solana.Solana,
	processor processor.Processor,
	config *configs.AppConfig,
) {
	ctx, cancel := context.WithCancel(context.Background())
	dripProgramProcessor := DripProgramProcessor{
		client:      client,
		processor:   processor,
		cancel:      cancel,
		environment: config.Environment,
	}
	lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go func() {
				_ = dripProgramProcessor.start(ctx)
			}()
			return nil
		},
		OnStop: func(_ context.Context) error {
			dripProgramProcessor.stop()
			return nil
		},
	})
}

func (d DripProgramProcessor) start(ctx context.Context) error {
	// These "token swaps" have potentially different interfaces then the token_swap program
	//"9W959DqEETiGZocYWCQPaJ6sBmUzgfxXfqGeTEdp3aQP", // orca swap v2
	//"DjVE6JNiYqPL2QXyCUUh8rNjHrbz9hXHNYt99MQ59qw1", // orca swap v1
	//"SSwapUtytfBdBn1b9NUGG6foMVPtcWgpRU32HToDUZr",  // Saros AMM
	//"PSwapMdSai8tjrEXcxFeQth87xC4rRsa4VA5mhGhXkP",  // Penguin Swap

	// Track Drip accounts
	if err := d.client.ProgramSubscribe(ctx, drip.ProgramID.String(), d.processDripEvent); err != nil {
		return err
	}
	go d.Backfill(context.Background(), drip.ProgramID.String(), d.processDripEvent)

	// Track token_swap program accounts
	if err := d.client.ProgramSubscribe(ctx, tokenswap.ProgramID.String(), d.processTokenSwapEvent); err != nil {
		return err
	}
	go d.Backfill(context.Background(), tokenswap.ProgramID.String(), d.processTokenSwapEvent)

	// Track orca_whirlpool program accounts
	// TODO(Mocha): remove processor
	if err := d.client.ProgramSubscribe(ctx, whirlpool.ProgramID.String(), d.processWhirlpoolEvent); err != nil {
		return err
	}
	go d.Backfill(context.Background(), whirlpool.ProgramID.String(), d.processWhirlpoolEvent)

	// Track Balance Updates Live
	if err := d.client.ProgramSubscribe(ctx, token.ProgramID.String(), d.processTokenEvent); err != nil {
		return err
	}
	return nil
}

func (d DripProgramProcessor) stop() {
	d.cancel()
}

func (d DripProgramProcessor) Backfill(ctx context.Context, programID string, processor func(string, []byte)) {
	log := logrus.WithField("program", programID).WithField("func", "Backfill")
	accounts, err := d.client.GetProgramAccounts(ctx, programID)
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
		err := d.client.GetAccounts(ctx, accounts[start:end], func(address string, data []byte) {
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

func (d DripProgramProcessor) processDripEvent(address string, data []byte) {
	if pubkey, err := solana2.PublicKeyFromBase58(address); err != nil || pubkey.IsZero() {
		logrus.WithField("address", address).Info("skipping zero/invalid address")
	}
	ctx := context.Background()
	log := logrus.WithField("address", address)
	//log.Infof("received drip account update")
	defer func() {
		if r := recover(); r != nil {
			log.WithField("stack", debug.Stack()).Errorf("panic in processEvent")
		}
	}()
	var vaultPeriod drip.VaultPeriod
	if err := bin.NewBinDecoder(data).Decode(&vaultPeriod); err == nil {
		//log.Infof("decoded as vaultPeriod")
		if err := d.processor.UpsertVaultPeriodByAddress(ctx, address); err != nil {
			log.WithError(err).Error("failed to upsert vaultPeriod")
		}
		return
	}
	var position drip.Position
	if err := bin.NewBinDecoder(data).Decode(&position); err == nil {
		//log.Infof("decoded as position")
		if err := d.processor.UpsertPositionByAddress(ctx, address); err != nil {
			log.WithError(err).Error("failed to upsert position")
		}
		return
	}
	var vault drip.Vault
	if err := bin.NewBinDecoder(data).Decode(&vault); err == nil {
		//log.Infof("decoded as vault")
		if err := d.processor.UpsertVaultByAddress(ctx, address); err != nil {
			log.WithError(err).Error("failed to upsert vault")
		}
		return
	}
	var protoConfig drip.VaultProtoConfig
	if err := bin.NewBinDecoder(data).Decode(&protoConfig); err == nil {
		//log.Infof("decoded as protoConfig")
		if err := d.processor.UpsertProtoConfigByAddress(ctx, address); err != nil {
			log.WithError(err).Error("failed to upsert protoConfig")
		}
		return
	}
	log.Errorf("failed to decode drip account to known types")
}

func (d DripProgramProcessor) processTokenSwapEvent(address string, data []byte) {
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
		if err := d.processor.UpsertTokenSwapByAddress(ctx, address); err != nil {
			log.WithError(err).Error("failed to upsert tokenSwap")
		}
		return
	}
	//log.WithError(err).Errorf("failed to decode token swap account")
}

func (d DripProgramProcessor) processWhirlpoolEvent(address string, data []byte) {
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
		if err := d.processor.UpsertWhirlpoolByAddress(ctx, address); err != nil {
			log.WithError(err).Error("failed to upsert tokenSwap")
		}
		return
	}
}

func (d DripProgramProcessor) processTokenEvent(address string, data []byte) {
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
		if err := d.processor.UpsertTokenAccountBalance(ctx, address, tokenAccount, false); err != nil {
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
