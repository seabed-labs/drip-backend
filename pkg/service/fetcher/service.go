package fetcher

import (
	"github.com/dcaf-labs/drip/pkg/clients/solana"
	"github.com/dcaf-labs/drip/pkg/repository"
)

type Fetcher interface {
}

type impl struct {
	repo   repository.Repository
	client solana.Solana
}

func NewFetcher(
	repo repository.Repository,
	client solana.Solana,
) Fetcher {
	return impl{
		repo:   repo,
		client: client,
	}
}
