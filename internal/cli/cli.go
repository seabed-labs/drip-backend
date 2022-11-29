package cli

import (
	"os"

	"github.com/dcaf-labs/drip/pkg/service/config"

	"github.com/dcaf-labs/drip/pkg/service/processor"
	"github.com/urfave/cli/v2"
)

type impl struct {
	network   config.Network
	env       config.Environment
	processor processor.Processor
}

func RunCLI(
	appConfig config.AppConfig,
	processor processor.Processor,
) error {
	i := impl{
		network:   appConfig.GetNetwork(),
		env:       appConfig.GetEnvironment(),
		processor: processor,
	}
	cliApp := &cli.App{
		Name:  "Drip API CLI",
		Usage: "To be used to execute one-time/in-frequent scripts and utilities",
		Commands: cli.Commands{
			&cli.Command{
				Name: "backfill",
				Subcommands: cli.Commands{
					getBackfillTokenCommand(i),
					getBackfillDevnetAccountsCommand(i),
				},
			},
		},
	}
	return cliApp.Run(os.Args)
}
