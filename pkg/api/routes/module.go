package controller

import (
	"github.com/dcaf-labs/drip/pkg/clients/solana"
	"github.com/dcaf-labs/drip/pkg/configs"
	"github.com/dcaf-labs/drip/pkg/repository"
	"github.com/gorilla/schema"
)

const defaultLimit = 500

type Handler struct {
	decoder      *schema.Decoder
	solanaClient solana.Solana
	repo         repository.Repository
	env          configs.Environment
	port         int
}

func NewHandler(
	config *configs.AppConfig,
	solanaClient solana.Solana,
	repo repository.Repository,
) *Handler {
	return &Handler{
		decoder:      schema.NewDecoder(),
		solanaClient: solanaClient,
		repo:         repo,
		env:          config.Environment,
		port:         config.Port,
	}
}
