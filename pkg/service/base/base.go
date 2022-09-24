package base

import (
	"github.com/dcaf-labs/drip/pkg/service/clients/solana"
	"github.com/dcaf-labs/drip/pkg/service/configs"
	"github.com/dcaf-labs/drip/pkg/service/repository"
)

type Base struct {
	solanaClient solana.Solana
	repo         repository.Repository
	network      configs.Network
	env          configs.Environment
}
