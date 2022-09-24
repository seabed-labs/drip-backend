package base

import (
	"context"

	"github.com/dcaf-labs/drip/pkg/service/repository"
	"github.com/dcaf-labs/drip/pkg/service/repository/model"
)

type Base interface {
	GetBestTokenSwapsForVaults(ctx context.Context, vaults []*model.Vault) (map[string]*model.TokenSwap, error)
}

func NewBase(
	repo repository.Repository,
) Base {
	return newBaseService(repo)
}

type impl struct {
	repo repository.Repository
}

func newBaseService(
	repo repository.Repository,
) impl {
	return impl{
		repo: repo,
	}
}
