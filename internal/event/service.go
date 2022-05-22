package event

import (
	"context"
	"runtime/debug"

	"github.com/dcaf-protocol/drip/internal/pkg/clients/solana/token_swap"

	"github.com/dcaf-protocol/drip/internal/pkg/clients/solana"

	"github.com/dcaf-protocol/drip/internal/pkg/clients/solana/dca_vault"

	"github.com/dcaf-protocol/drip/internal/pkg/processor"

	bin "github.com/gagliardetto/binary"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

type DripProgramProcessor struct {
	client    solana.Solana
	processor processor.Processor
	cancel    context.CancelFunc
}

func NewDripProgramProcessor(
	lifecycle fx.Lifecycle,
	client solana.Solana,
	processor processor.Processor,
) *DripProgramProcessor {
	ctx, cancel := context.WithCancel(context.Background())
	dripProgramProcessor := DripProgramProcessor{
		client:    client,
		processor: processor,
		cancel:    cancel,
	}
	lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			return dripProgramProcessor.start(ctx)
		},
		OnStop: func(_ context.Context) error {
			dripProgramProcessor.stop()
			return nil
		},
	})
	return &dripProgramProcessor
}

func (d DripProgramProcessor) start(ctx context.Context) error {
	// TODO(Mocha): the program ID's should be in a config since they will change
	if err := d.client.ProgramSubscribe(ctx, dca_vault.ProgramID.String(), d.processDripEvent); err != nil {
		return err
	}
	if err := d.client.ProgramSubscribe(ctx, token_swap.ProgramID.String(), d.processTokenSwapEvent); err != nil {
		return err
	}
	return nil
}

func (d DripProgramProcessor) stop() {
	d.cancel()
}

// TODO(Mocha): Backfill using GetProgramAccounts
func (d DripProgramProcessor) processDripEvent(address string, data []byte) {
	ctx := context.Background()
	log := logrus.WithField("address", address)
	log.Infof("received drip account update")
	defer func() {
		if r := recover(); r != nil {
			log.WithField("stack", debug.Stack()).Errorf("panic in processEvent")
		}
	}()
	var vaultPeriod dca_vault.VaultPeriod
	if err := bin.NewBinDecoder(data).Decode(&vaultPeriod); err == nil {
		log.Infof("decoded as vaultPeriod")
		if err := d.processor.UpsertVaultPeriodByAddress(ctx, address); err != nil {
			log.WithError(err).Error("failed to upsert vaultPeriod")
		}
		return
	}
	var position dca_vault.Position
	if err := bin.NewBinDecoder(data).Decode(&position); err == nil {
		log.Infof("decoded as position")
		if err := d.processor.UpsertPositionByAddress(ctx, address); err != nil {
			log.WithError(err).Error("failed to upsert position")
		}
		return
	}
	var vault dca_vault.Vault
	if err := bin.NewBinDecoder(data).Decode(&vault); err == nil {
		log.Infof("decoded as vault")
		if err := d.processor.UpsertVaultByAddress(ctx, address); err != nil {
			log.WithError(err).Error("failed to upsert vault")
		}
		return
	}
	var protoConfig dca_vault.VaultProtoConfig
	if err := bin.NewBinDecoder(data).Decode(&protoConfig); err == nil {
		log.Infof("decoded as protoConfig")
		if err := d.processor.UpsertProtoConfigByAddress(ctx, address); err != nil {
			log.WithError(err).Error("failed to upsert protoConfig")
		}
		return
	}
	log.Errorf("failed to decode account")
}

func (d DripProgramProcessor) processTokenSwapEvent(address string, data []byte) {
	ctx := context.Background()
	log := logrus.WithField("address", address)
	log.Infof("received token swap account update")
	defer func() {
		if r := recover(); r != nil {
			log.WithField("stack", debug.Stack()).Errorf("panic in processTokenSwapEvent")
		}
	}()
	var tokenSwap token_swap.TokenSwap
	err := bin.NewBinDecoder(data).Decode(&tokenSwap)
	if err == nil {
		log.Infof("decoded as tokenSwap")
		if err := d.processor.UpsertTokenSwapByAddress(ctx, address); err != nil {
			log.WithError(err).Error("failed to upsert tokenSwap")
		}
		return
	}
	log.WithError(err).Errorf("failed to decode token swap account")
}
