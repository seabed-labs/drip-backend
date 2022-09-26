package event

import (
	"context"
	"math/rand"
	"strconv"
	"time"

	"github.com/dcaf-labs/drip/pkg/service/clients/solana"
	"github.com/dcaf-labs/drip/pkg/service/configs"
	"github.com/dcaf-labs/drip/pkg/service/processor"
	"github.com/dcaf-labs/solana-go-clients/pkg/drip"
	"github.com/dcaf-labs/solana-go-clients/pkg/tokenswap"
	"github.com/dcaf-labs/solana-go-clients/pkg/whirlpool"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

const backfillEvery = time.Hour * 1

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
	if err := d.client.ProgramSubscribe(ctx, drip.ProgramID.String(), d.processor.ProcessDripEvent); err != nil {
		return err
	}

	// In staging, we manually backfill tokenswaps and whirlpools so that we can limit the # of rows in the DB
	if configs.IsProd(d.environment) {
		// Track token_swap program accounts
		if err := d.client.ProgramSubscribe(ctx, tokenswap.ProgramID.String(), d.processor.ProcessTokenSwapEvent); err != nil {
			return err
		}
		// Don't need to constantly backfill these, just do it once
		//go d.processor.Backfill(context.Background(), tokenswap.ProgramID.String(), d.processor.ProcessTokenSwapEvent)

		// Track orca_whirlpool program accounts
		if err := d.client.ProgramSubscribe(ctx, whirlpool.ProgramID.String(), d.processor.ProcessWhirlpoolEvent); err != nil {
			return err
		}
		// Don't need to constantly backfill these, just do it once
		//go d.processor.Backfill(context.Background(), whirlpool.ProgramID.String(), d.processor.ProcessWhirlpoolEvent)
	}

	// Track Balance Updates Live
	// Too many messages... need to implement an actual queue before we re-enable this
	//if err := d.client.ProgramSubscribe(ctx, token.ProgramID.String(), d.processor.ProcessTokenEvent); err != nil {
	//	return err
	//}

	go d.runBackfill(ctx)
	return nil
}

func (d *DripProgramProcessor) stop() {
	d.cancel()
}

func (d *DripProgramProcessor) runBackfill(ctx context.Context) {
	defer func() {
		time.AfterFunc(backfillEvery, func() {
			if ctx.Err() == context.Canceled {
				return
			}
			d.runBackfill(ctx)
		})
	}()
	id := strconv.FormatInt(int64(rand.Int()), 10)
	log := logrus.WithField("id", id)
	log.Info("starting backfill")
	d.processor.Backfill(ctx, id, drip.ProgramID.String(), d.processor.ProcessDripEvent)
	log.Info("done backfill")
}
