package base

import (
	"context"

	"github.com/dcaf-labs/drip/pkg/service/configs"

	"github.com/dcaf-labs/drip/pkg/service/repository"
	"github.com/dcaf-labs/drip/pkg/service/repository/model"
)

type Base interface {
	GetBestTokenSwapsForVaults(ctx context.Context, vaults []*model.Vault) (map[string]*model.TokenSwap, error)
	GetBestOrcaWhirlpoolForVaults(ctx context.Context, vaults []*model.Vault) (map[string]*model.OrcaWhirlpool, error)
}

func NewBase(
	config *configs.AppConfig,
	repo repository.Repository,
) Base {
	return newBaseService(repo, config.Network)
}

type impl struct {
	repo    repository.Repository
	network configs.Network
}

func newBaseService(
	repo repository.Repository,
	network configs.Network,
) impl {
	return impl{
		repo:    repo,
		network: network,
	}
}
