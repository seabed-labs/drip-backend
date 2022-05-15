package api

import (
	"github.com/dcaf-protocol/drip/internal/configs"
	"github.com/dcaf-protocol/drip/internal/pkg/clients/solana"
	"github.com/gorilla/schema"
)

type Handler struct {
	decoder      *schema.Decoder
	solanaClient solana.Solana
	env          configs.Environment
	port         int
}

func NewHandler(
	solanaClient solana.Solana,
	config *configs.AppConfig,
) *Handler {
	return &Handler{
		decoder:      schema.NewDecoder(),
		solanaClient: solanaClient,
		env:          config.Environment,
		port:         config.Port,
	}
}
