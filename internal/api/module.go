package api

import (
	"github.com/dcaf-protocol/drip/internal/configs"
	"github.com/dcaf-protocol/drip/internal/pkg/client"
	"github.com/gorilla/schema"
)

type Handler struct {
	decoder      *schema.Decoder
	solanaClient client.Solana
	env          configs.Environment
	port         int
}

func NewHandler(
	solanaClient client.Solana,
	config *configs.Config,
) *Handler {
	return &Handler{
		decoder:      schema.NewDecoder(),
		solanaClient: solanaClient,
		env:          config.Environment,
		port:         config.Port,
	}
}
