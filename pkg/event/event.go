package event

import (
	"context"
	"fmt"
	"math/rand"
	"sync"

	"github.com/dcaf-labs/drip/pkg/clients/solana"
	"github.com/dcaf-labs/drip/pkg/configs"
	"github.com/dcaf-labs/drip/pkg/service/processor"
	"github.com/dcaf-labs/solana-go-clients/pkg/drip"
	"github.com/dcaf-labs/solana-go-clients/pkg/tokenswap"
	"github.com/dcaf-labs/solana-go-clients/pkg/whirlpool"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

const backfillEvery = 3600

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
) error {
	ctx, cancel := context.WithCancel(context.Background())
	dripProgramProcessor := DripProgramProcessor{
		client:      client,
		processor:   processor,
		cancel:      cancel,
		environment: config.Environment,
	}
	job := cron.New()
	if _, err := job.AddFunc(fmt.Sprintf("@every %ds", backfillEvery), dripProgramProcessor.runBackfill); err != nil {
		logrus.WithError(err).Error("failed to addFunc to cronJob")
		return err
	}
	lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go func() {
				_ = dripProgramProcessor.start(ctx)
			}()
			job.Start()
			return nil
		},
		OnStop: func(_ context.Context) error {
			dripProgramProcessor.stop()
			_ = job.Stop()
			return nil
		},
	})
	go dripProgramProcessor.runBackfill()
	return nil
}

func (d DripProgramProcessor) start(ctx context.Context) error {
	// Track Drip accounts
	if err := d.client.ProgramSubscribe(ctx, drip.ProgramID.String(), d.processor.ProcessDripEvent); err != nil {
		return err
	}

	// Track token_swap program accounts
	if err := d.client.ProgramSubscribe(ctx, tokenswap.ProgramID.String(), d.processor.ProcessTokenSwapEvent); err != nil {
		return err
	}
	// Don't need to constantly backfill these, just do it once
	go d.processor.Backfill(context.Background(), tokenswap.ProgramID.String(), d.processor.ProcessTokenSwapEvent)

	// Track orca_whirlpool program accounts
	if err := d.client.ProgramSubscribe(ctx, whirlpool.ProgramID.String(), d.processor.ProcessWhirlpoolEvent); err != nil {
		return err
	}
	// Don't need to constantly backfill these, just do it once
	go d.processor.Backfill(context.Background(), whirlpool.ProgramID.String(), d.processor.ProcessWhirlpoolEvent)

	// Track Balance Updates Live
	if err := d.client.ProgramSubscribe(ctx, token.ProgramID.String(), d.processor.ProcessTokenEvent); err != nil {
		return err
	}
	return nil
}

func (d DripProgramProcessor) runBackfill() {
	id := rand.Int()
	log := logrus.WithField("id", id)
	log.Info("starting backfill")
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		d.processor.Backfill(context.Background(), drip.ProgramID.String(), d.processor.ProcessDripEvent)
	}()
	wg.Wait()
	log.Info("done backfill")
}

func (d DripProgramProcessor) stop() {
	d.cancel()
}
