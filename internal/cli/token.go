package cli

import (
	"github.com/hashicorp/go-multierror"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func getBackfillTokenCommand(i impl) *cli.Command {
	var tokenMints cli.StringSlice
	return &cli.Command{
		Name: "token",
		Action: func(cCtx *cli.Context) (err error) {
			tokenMintAddresses := tokenMints.Value()
			logrus.WithField("len(tokenMintAddresses)", len(tokenMintAddresses)).Info("starting token backfill")
			for _, mintAddress := range tokenMintAddresses {
				if upsertErr := i.processor.UpsertTokenByAddress(cCtx.Context, mintAddress); upsertErr != nil {
					err = multierror.Append(err, upsertErr)
				}
				logrus.WithField("mint", mintAddress).Info("backfilled token")
			}
			logrus.Info("done")
			return err
		},
		Flags: []cli.Flag{
			&cli.StringSliceFlag{Name: "tokenMints", Required: true, Destination: &tokenMints},
		},
	}
}
