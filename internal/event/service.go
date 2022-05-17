package event

import (
	"context"
	"time"

	"github.com/dcaf-protocol/drip/internal/pkg/repository/model"
	"github.com/shopspring/decimal"

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
	return d.client.ProgramSubscribe(ctx, dca_vault.ProgramID.String(), d.processEvent)
}

func (d DripProgramProcessor) stop() {
	d.cancel()
}

//
func (d DripProgramProcessor) processEvent(address string, data []byte) {
	logrus.WithField("address", address).Infof("received drip account update")
	var vaultPeriod dca_vault.VaultPeriod
	if err := bin.NewBinDecoder(data).Decode(&vaultPeriod); err == nil {
		logrus.WithField("address", address).Infof("decoded as vaultPeriod")
		twap, err := decimal.NewFromString(vaultPeriod.Twap.String())
		if err != nil {
			logrus.WithError(err).Errorf("failed to decode twap as decimal")
			return
		}
		if err := d.processor.UpsertVaultPeriods(context.Background(), &model.VaultPeriod{
			Pubkey:   address,
			Vault:    vaultPeriod.Vault.String(),
			PeriodID: vaultPeriod.PeriodId,
			Twap:     twap,
			Dar:      vaultPeriod.Dar,
		}); err != nil {
			logrus.WithError(err).Errorf("failed to upsert vault period")
			return
		}
		return
	}
	var position dca_vault.Position
	if err := bin.NewBinDecoder(data).Decode(&position); err == nil {
		logrus.WithField("address", address).Infof("decoded as position")
		if err := d.processor.UpsertPositions(context.Background(), &model.Position{
			Pubkey:                   address,
			Vault:                    position.Vault.String(),
			Authority:                position.PositionAuthority.String(),
			DepositedTokenAAmount:    position.DepositedTokenAAmount,
			WithdrawnTokenBAmount:    position.WithdrawnTokenBAmount,
			DepositTimestamp:         time.Unix(position.DepositTimestamp, 0),
			DcaPeriodIDBeforeDeposit: position.DcaPeriodIdBeforeDeposit,
			NumberOfSwaps:            position.NumberOfSwaps,
			PeriodicDripAmount:       position.PeriodicDripAmount,
			IsClosed:                 position.IsClosed,
		}); err != nil {
			logrus.WithError(err).Errorf("failed to upsert position")
			return
		}
		return
	}
	// TODO(mocha): we don't want to process vault changes, we don't want to ingest random vaults for security purposes
	var vault dca_vault.Vault
	if err := bin.NewBinDecoder(data).Decode(&vault); err == nil {
		logrus.WithField("address", address).Infof("decoded as vault")
		//if err := d.processor.UpsertVaults(context.Background(), &model.Vault{
		//	Pubkey:                 address,
		//	ProtoConfig:            vault.ProtoConfig.String(),
		//	TokenAMint:             vault.TokenAMint.String(),
		//	TokenBMint:             vault.TokenBMint.String(),
		//	TokenAAccount:          vault.TokenAAccount.String(),
		//	TokenBAccount:          vault.TokenBAccount.String(),
		//	TreasuryTokenBAccount:  vault.TreasuryTokenBAccount.String(),
		//	LastDcaPeriod:          vault.LastDcaPeriod,
		//	DripAmount:             vault.DripAmount,
		//	DcaActivationTimestamp: time.Unix(vault.DcaActivationTimestamp, 0),
		//	Enabled:                false,
		//}); err != nil {
		//	logrus.WithError(err).Errorf("failed to upsert vault")
		//	return
		//}
		//return
	}
	var protoConfig dca_vault.VaultProtoConfig
	if err := bin.NewBinDecoder(data).Decode(&protoConfig); err == nil {
		logrus.WithField("address", address).Infof("decoded as protoConfig")
		if err := d.processor.UpsertProtoConfigs(context.Background(), &model.ProtoConfig{
			Pubkey:               address,
			Granularity:          protoConfig.Granularity,
			TriggerDcaSpread:     protoConfig.TriggerDcaSpread,
			BaseWithdrawalSpread: protoConfig.BaseWithdrawalSpread,
		}); err != nil {
			logrus.WithError(err).Errorf("failed to upsert proto config")
			return
		}
		return
	}
	logrus.WithField("address", address).Errorf("failed to decode account")
}
