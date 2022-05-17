package api

import (
	"github.com/dcaf-protocol/drip/internal/configs"
	"github.com/dcaf-protocol/drip/internal/pkg/clients/solana"
	"github.com/dcaf-protocol/drip/internal/pkg/drip"
	"github.com/gorilla/schema"
)

const defaultLimit = 500

type Handler struct {
	decoder      *schema.Decoder
	solanaClient solana.Solana
	drip         drip.Drip
	env          configs.Environment
	port         int
}

func NewHandler(
	config *configs.AppConfig,
	solanaClient solana.Solana,
	drip drip.Drip,
) *Handler {
	return &Handler{
		decoder:      schema.NewDecoder(),
		solanaClient: solanaClient,
		drip:         drip,
		env:          config.Environment,
		port:         config.Port,
	}
}
