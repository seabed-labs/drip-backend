package event

import (
	"context"
	"runtime/debug"

	"github.com/dcaf-protocol/drip/pkg/clients/solana"
	drip2 "github.com/dcaf-protocol/drip/pkg/clients/solana/drip"
	token_swap2 "github.com/dcaf-protocol/drip/pkg/clients/solana/token_swap"
	"github.com/dcaf-protocol/drip/pkg/configs"
	"github.com/dcaf-protocol/drip/pkg/processor"

	bin "github.com/gagliardetto/binary"
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
	if err := d.client.ProgramSubscribe(ctx, drip2.ProgramID.String(), d.processDripEvent); err != nil {
		return err
	}
	go d.Backfill(context.Background(), drip2.ProgramID.String(), d.processDripEvent)

	for _, swapProgram := range []string{
		token_swap2.ProgramID.String(),
		"9W959DqEETiGZocYWCQPaJ6sBmUzgfxXfqGeTEdp3aQP", // orca swap v2
		"DjVE6JNiYqPL2QXyCUUh8rNjHrbz9hXHNYt99MQ59qw1", // orca swap v1
		"SSwapUtytfBdBn1b9NUGG6foMVPtcWgpRU32HToDUZr",  // Saros AMM
		"PSwapMdSai8tjrEXcxFeQth87xC4rRsa4VA5mhGhXkP",  // Penguin Swap
	} {
		if err := d.client.ProgramSubscribe(ctx, swapProgram, d.processTokenSwapEvent); err != nil {
			return err
		}
		go d.Backfill(context.Background(), swapProgram, d.processTokenSwapEvent)
	}
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
	ctx := context.Background()
	log := logrus.WithField("address", address)
	//log.Infof("received drip account update")
	defer func() {
		if r := recover(); r != nil {
			log.WithField("stack", debug.Stack()).Errorf("panic in processEvent")
		}
	}()
	var vaultPeriod drip2.VaultPeriod
	if err := bin.NewBinDecoder(data).Decode(&vaultPeriod); err == nil {
		//log.Infof("decoded as vaultPeriod")
		if err := d.processor.UpsertVaultPeriodByAddress(ctx, address); err != nil {
			log.WithError(err).Error("failed to upsert vaultPeriod")
		}
		return
	}
	var position drip2.Position
	if err := bin.NewBinDecoder(data).Decode(&position); err == nil {
		//log.Infof("decoded as position")
		if err := d.processor.UpsertPositionByAddress(ctx, address); err != nil {
			log.WithError(err).Error("failed to upsert position")
		}
		return
	}
	var vault drip2.Vault
	if err := bin.NewBinDecoder(data).Decode(&vault); err == nil {
		//log.Infof("decoded as vault")
		if err := d.processor.UpsertVaultByAddress(ctx, address); err != nil {
			log.WithError(err).Error("failed to upsert vault")
		}
		return
	}
	var protoConfig drip2.VaultProtoConfig
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
	ctx := context.Background()
	log := logrus.WithField("address", address)
	//log.Infof("received token swap account update")
	defer func() {
		if r := recover(); r != nil {
			log.WithField("stack", debug.Stack()).Errorf("panic in processTokenSwapEvent")
		}
	}()
	var tokenSwap token_swap2.TokenSwap
	err := bin.NewBinDecoder(data).Decode(&tokenSwap)
	if err == nil {
		//log.Infof("decoded as tokenSwap")
		if err := d.processor.UpsertTokenSwapByAddress(ctx, address); err != nil {
			log.WithError(err).Error("failed to upsert tokenSwap")
		}
		return
	}
	log.WithError(err).Errorf("failed to decode token swap account")
}

func (d DripProgramProcessor) processTokenEvent(address string, data []byte) {
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
		if err := d.processor.UpsertTokenAccountBalance(ctx, address, tokenAccount); err != nil {
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
