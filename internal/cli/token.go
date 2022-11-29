package cli

import (
	"github.com/hashicorp/go-multierror"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func getBackfillTokenCommand(i impl) *cli.Command {
	var tokenMints cli.StringSlice
	tokenMints.Value()
	return &cli.Command{
		Name: "token",
		Action: func(cCtx *cli.Context) (err error) {
			for _, mintAddress := range tokenMints.Value() {
				if upsertErr := i.processor.UpsertTokenByAddress(cCtx.Context, mintAddress); upsertErr != nil {
					err = multierror.Append(err, upsertErr)
				}
			}
			logrus.Info("done")
			return err
		},
		Flags: []cli.Flag{
			&cli.StringSliceFlag{Name: "tokenMints", Aliases: []string{"N"}, Required: true, Destination: &tokenMints},
		},
	}
}
