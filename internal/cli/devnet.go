package cli

import (
	"context"

	"github.com/dcaf-labs/drip/pkg/service/config"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

func getBackfillDevnetAccountsCommand(i impl) *cli.Command {
	return &cli.Command{
		Name:        "devnet",
		Description: "Backfill a hardcoded list of devnet accounts",
		Action: func(cCtx *cli.Context) (err error) {
			return backfillDevnetAccounts(cCtx.Context, i)
		},
	}
}

func backfillDevnetAccounts(
	ctx context.Context, i impl,
) error {
	logrus.WithField("network", i.network).Infof("starting cli")
	if config.IsDevnetNetwork(i.network) {
		for _, address := range []string{
			"35WMYrE8E4vbmm4phVxkRDTK5gE8db2KUVGFTYduE1Uz",
			"3MiLbpHuKHDEnMUNpDPhJmMQPzcJL2Gp8kQdGMRHPcwP",
			"536VybxwBRjXP8JSQTjVsjMffkzgdRErZReCk1Ja7Psr",
			"Bgn6dKZhWRYAGmsaE34h5tqgPg3iYUmUvqz2MBqr3hYA",
			"BkMeyctybgU3yaboRrMjtgvqSXHaPBarYawcLejGeJg3",
			"CyxNwQH1WPcH7NkCULdj4Xf6aMU3S9ws7YAoyzAZkU6h",
			"Evy3P7kFy6epfbY3QZmwWiqWj1veZ87LQoCEWYzb3vW7",
			"EZZTdHW9rskAJs4HGPwcM2CFeJ5BpRSdFmfTjZAVktwh",
			"FcbTUmkEuizRjnr17iggHon9rAc2FeQ8uMdjxWBgFb58",
			"GUr5RGCrS1bxvsiAHrLkQvK1WS6QFRCCy7V72mkN4b7s",
			"H81dLAxwMFSy4HRqsDQWJq8BVtQW81sapENekRBSNUj7",
		} {
			log := logrus.WithField("address", address)
			if err := i.processor.UpsertTokenSwapByAddress(ctx, address); err != nil {
				log.WithError(err).Error("failed to cli tokenSwap")
			}
		}
		for _, address := range []string{
			"2w9DNJRFGmvN2wuVi3CtLT49cM8DWdKfzAGas1XDw3Ve",
			"5fkps3wttvX3ysprtWzLRuxajSkmdxEa12Ys8E4bMTPh",
			"ADPEtfPLmn5Nb92dm6MFUEmmeFxyMXiX85JRfN5e8xyo",
			"C5CBERnsLjFNDPC6xNtjyFR8HeDDcV9ZYKUgGUHNFKEE",
			"Dr75jCuqpkYCGs4ATp2v31sU6bzDKoCdxJNvxXUUgb4S",
			"GSFnjnJ7TdSsGWb6JgFhWakWrv8VGZUAghnY3EA8Xj46",
		} {
			log := logrus.WithField("address", address)
			if err := i.processor.UpsertWhirlpoolByAddress(ctx, address); err != nil {
				log.WithError(err).Error("failed to cli whirlpool")
			}
		}
	}
	if config.IsDevnetNetwork(i.network) && config.IsStagingEnvironment(i.env) {
		for _, address := range []string{
			"HhKqRDp2ooZPj8554rHtJLGiMsDrHy1zkZqidDM1h9kQ",
		} {
			log := logrus.WithField("address", address)
			if err := i.processor.UpsertPositionByAddress(ctx, address); err != nil {
				log.WithError(err).Error("failed to cli position")
			}
		}
	}
	if config.IsDevnetNetwork(i.network) && config.IsProductionEnvironment(i.env) {
		for _, address := range []string{
			"5VfSyiFenN99Nk3KTsuB93E6783cpB1rdJkjFdg9qSLK",
			"J3nPeD3VrP3i23pDgsG9uXiPURd7ptRXoixL8CJRQbRW",
			"BY2YSxzwZwPh7MAJ86hsbu1uop9SyhZWyKqfXtN6FNu4",
			"6mCv8tF2wxq3pjPaT7r7Qf9xLyTwQwWJMYncdJsatpDP",
			"EErEQN63Tubyq7zHRW9y4ndHukPs3hMTEq6zQG7LQETz",
			"Chn9T1M93piu89GnnPDzAsjHwKjoMKC8CCgX9wmtvUqp",
			"B89HUcrgyNRCAffb33v52NucUrYbMZtxfPzcW3EwzXWs",
			"Bqkkq8AsaAyhgL53zEazJ1wYMqiHjEdF7osA8XAREE2q",
		} {
			log := logrus.WithField("address", address)
			if err := i.processor.UpsertVaultByAddress(ctx, address); err != nil {
				log.WithError(err).Error("failed to cli vault")
			}
		}
	}
	logrus.Infof("done backfilling")
	return nil
}
