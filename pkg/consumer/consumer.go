package consumer

import (
	"context"

	"github.com/dcaf-labs/drip/pkg/service/clients/solana"
	"github.com/dcaf-labs/drip/pkg/service/config"
	"github.com/dcaf-labs/drip/pkg/service/processor"
	"go.uber.org/fx"
)

type DripProgramConsumer struct {
	client      solana.Solana
	processor   processor.Processor
	cancel      context.CancelFunc
	environment config.Environment
}

func Server(
	lifecycle fx.Lifecycle,
	client solana.Solana,
	processor processor.Processor,
	appConfig config.AppConfig,
) error {
	ctx, cancel := context.WithCancel(context.Background())
	dripProgramProcessor := DripProgramConsumer{
		client:      client,
		processor:   processor,
		cancel:      cancel,
		environment: appConfig.GetEnvironment(),
	}
	lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			if err := dripProgramProcessor.start(ctx); err != nil {
				return err
			}
			return nil
		},
		OnStop: func(_ context.Context) error {
			dripProgramProcessor.stop()
			return nil
		},
	})
	return nil
}

func (d *DripProgramConsumer) start(ctx context.Context) error {
	go d.processor.ProcessAccountUpdateQueue(ctx)
	go d.processor.ProcessTransactionUpateQueue(ctx)
	return nil
}

func (d *DripProgramConsumer) stop() {
	d.cancel()
}
