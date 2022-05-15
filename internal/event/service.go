package event

import (
	"context"

	"github.com/dcaf-protocol/drip/internal/pkg/clients/solana"
	"github.com/dcaf-protocol/drip/internal/pkg/clients/solana/dca_vault"
	"github.com/dcaf-protocol/drip/internal/pkg/repository"
	bin "github.com/gagliardetto/binary"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

type DripProgramProcessor struct {
	repo   *repository.Query
	client solana.Solana
	cancel context.CancelFunc
}

func NewDripProgramProcessor(
	lifecycle fx.Lifecycle,
	repo *repository.Query,
	client solana.Solana,
) *DripProgramProcessor {
	ctx, cancel := context.WithCancel(context.Background())
	processor := DripProgramProcessor{
		repo:   repo,
		client: client,
		cancel: cancel,
	}
	lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			return processor.start(ctx)
		},
		OnStop: func(_ context.Context) error {
			processor.stop()
			return nil
		},
	})
	return &processor
}

func (d DripProgramProcessor) start(ctx context.Context) error {
	return d.client.ProgramSubscribe(ctx, dca_vault.ProgramID.String(), d.processEvent)
}

func (d DripProgramProcessor) stop() {
	d.cancel()
}

func (d DripProgramProcessor) processEvent(address string, data []byte) {
	logrus.WithField("address", address).Infof("received drip account update")
	var vaultPeriod dca_vault.VaultPeriod
	if err := bin.NewBinDecoder(data).Decode(&vaultPeriod); err == nil {
		logrus.WithField("address", address).Infof("decoded as vaultPeriod")
		return
	}
	var position dca_vault.Position
	if err := bin.NewBinDecoder(data).Decode(&position); err == nil {
		logrus.WithField("address", address).Infof("decoded as position")
		return
	}
	var vault dca_vault.Vault
	if err := bin.NewBinDecoder(data).Decode(&vault); err == nil {
		logrus.WithField("address", address).Infof("decoded as vault")
		return
	}
	var protoConfig dca_vault.VaultProtoConfig
	if err := bin.NewBinDecoder(data).Decode(&protoConfig); err == nil {
		logrus.WithField("address", address).Infof("decoded as protoConfig")
		return
	}
	logrus.WithField("address", address).Errorf("failed to decode account")
}

//t.Run("ProgramSubscribe should subscribe to drip", func(t *testing.T) {
//	timeout := time.Minute * 30
//	ctx, cancel := context.WithTimeout(context.Background(), timeout)
//	defer cancel()
//	err := client.ProgramSubscribe(ctx, dca_vault.ProgramID.String(), func(address string, data []byte) {
//		var vaultPeriod dca_vault.VaultPeriod
//		logrus.WithField("address", address).Infof("got message")
//		if err := bin.NewBinDecoder(data).Decode(&vaultPeriod); err != nil {
//			logrus.
//				WithError(err).
//				Errorf("failed to decode as vault period")
//		} else {
//			logrus.
//				WithField("vaultPeriod", vaultPeriod).
//				Infof("decoded vault period")
//		}
//		assert.NotEmpty(t, data)
//	})
//	assert.NoError(t, err)
//	select {
//	case <-time.After(timeout):
//		break
//	}
//})
