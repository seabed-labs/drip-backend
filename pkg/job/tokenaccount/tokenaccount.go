package tokenaccount

import (
	"context"
	"time"

	"github.com/dcaf-labs/drip/pkg/service/processor"
	"github.com/dcaf-labs/drip/pkg/service/repository"
	"github.com/dcaf-labs/drip/pkg/service/repository/model"
	"github.com/go-co-op/gocron"
	"github.com/samber/lo"
	"github.com/sirupsen/logrus"
	"go.uber.org/fx"
)

type TokenJob interface {
	UpsertAllVaultTokensAccountAccounts(ctx context.Context)
}

type impl struct {
	repo      repository.Repository
	processor processor.Processor
}

func NewTokenAccountJob(
	lifecycle fx.Lifecycle,
	repo repository.Repository,
	processor processor.Processor,
) (TokenJob, error) {
	impl := impl{
		repo:      repo,
		processor: processor,
	}
	s := gocron.NewScheduler(time.UTC)
	ctx, cancel := context.WithCancel(context.Background())
	if _, err := s.Every(30).Minutes().Do(impl.UpsertAllVaultTokensAccountAccounts, ctx); err != nil {
		cancel()
		return nil, err
	}
	s.StartImmediately().StartAsync()
	lifecycle.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			cancel()
			s.Stop()
			return nil
		},
	})
	return impl, nil
}

func (i impl) UpsertAllVaultTokensAccountAccounts(ctx context.Context) {
	logrus.Info("starting UpsertAllVaultTokensAccountAccounts")
	vaults, err := i.repo.GetVaultsWithFilter(ctx, nil, nil, nil)
	if err != nil {
		logrus.WithError(err).Error("failed to GetVaultsWithFilter")
		return
	}
	vaultTokenAccounts := lo.FlatMap[*model.Vault, string](vaults, func(vault *model.Vault, _ int) []string {
		return []string{vault.TokenAAccount, vault.TokenBAccount, vault.TreasuryTokenBAccount}
	})
	if err := i.processor.UpsertTokenAccountsByAddresses(ctx, vaultTokenAccounts...); err != nil {
		logrus.WithError(err).WithField("len(tokens", len(vaultTokenAccounts)).Error("failed to UpsertTokenAccountsByAddresses")
	}
	logrus.Info("done UpsertAllVaultTokensAccountAccounts")
}
