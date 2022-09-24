package controller

import (
	"github.com/dcaf-labs/drip/pkg/service/base"
	"github.com/dcaf-labs/drip/pkg/service/clients/solana"
	"github.com/dcaf-labs/drip/pkg/service/configs"
	"github.com/dcaf-labs/drip/pkg/service/repository"
	"github.com/gorilla/schema"
)

const defaultLimit = 500

type Handler struct {
	decoder      *schema.Decoder
	base         base.Base
	solanaClient solana.Solana
	repo         repository.Repository
	network      configs.Network
	env          configs.Environment
	port         int
}

func NewHandler(
	config *configs.AppConfig,
	base base.Base,
	solanaClient solana.Solana,
	repo repository.Repository,
) *Handler {
	return &Handler{
		decoder:      schema.NewDecoder(),
		base:         base,
		solanaClient: solanaClient,
		repo:         repo,
		network:      config.Network,
		env:          config.Environment,
		port:         config.Port,
	}
}
