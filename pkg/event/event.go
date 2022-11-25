package event

import (
	"context"
	"math/rand"
	"runtime/debug"
	"strconv"
	"time"

	"github.com/dcaf-labs/drip/pkg/service/clients/solana"
	"github.com/dcaf-labs/drip/pkg/service/config"
	"github.com/dcaf-labs/drip/pkg/service/processor"
	"github.com/dcaf-labs/solana-go-clients/pkg/drip"
	"github.com/dcaf-labs/solana-go-clients/pkg/whirlpool"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

const backfillEvery = time.Hour * 12

type DripProgramProcessor struct {
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
	dripProgramProcessor := DripProgramProcessor{
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

func (d *DripProgramProcessor) start(ctx context.Context) error {
	// Track Drip accounts
	// We track drip accounts live
	if err := d.client.ProgramSubscribe(ctx, drip.ProgramID.String(), d.processor.AddItemToUpdateQueueCallback(ctx, drip.ProgramID.String())); err != nil {
		return err
	}

	// In staging, we manually backfill tokenswaps and whirlpools so that we can limit the # of rows in the DB
	if config.IsProductionEnvironment(d.environment) {
		// Track token_swap program accounts
		//if err := d.client.ProgramSubscribe(ctx, tokenswap.ProgramID.String(), d.processor.AddItemToUpdateQueueCallback(ctx, tokenswap.ProgramID.String())); err != nil {
		//	return err
		//}

		// Track orca_whirlpool program accounts
		if err := d.client.ProgramSubscribe(ctx, whirlpool.ProgramID.String(), d.processor.AddItemToUpdateQueueCallback(ctx, whirlpool.ProgramID.String())); err != nil {
			return err
		}
	}

	go d.runBackfill(ctx)
	go d.processor.ProcessAccountUpdateQueue(ctx)
	return nil
}

func (d *DripProgramProcessor) stop() {
	d.cancel()
}

func (d *DripProgramProcessor) runBackfill(ctx context.Context) {
	if d.environment == config.StagingEnv {
		logrus.Infof("skipping backfill on staging")
		return
	}
	defer func() {
		if r := recover(); r != nil {
			debug.PrintStack()
			logrus.WithField("stack", string(debug.Stack())).Errorf("panic in runBackfill")
		}
		time.AfterFunc(backfillEvery, func() {
			if ctx.Err() != nil {
				return
			}
			d.runBackfill(ctx)
		})
	}()
	id := strconv.FormatInt(int64(rand.Int()), 10)
	log := logrus.WithField("id", id)
	log.Info("starting backfill")
	d.processor.BackfillProgramOwnedAccounts(ctx, id, drip.ProgramID.String())
	log.Info("done backfill")
}
