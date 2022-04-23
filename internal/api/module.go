package api

import (
	"github.com/dcaf-protocol/drip/internal/configs"
	"github.com/dcaf-protocol/drip/internal/pkg/solanaclient"
	"github.com/gorilla/schema"
)

type Handler struct {
	decoder      *schema.Decoder
	solanaClient *solanaclient.Solana
	env          configs.Environment
	port         int
}

func NewHandler(
	solanaClient *solanaclient.Solana,
	config *configs.Config,
) *Handler {
	return &Handler{
		decoder:      schema.NewDecoder(),
		solanaClient: solanaClient,
		env:          config.Environment,
		port:         config.Port,
	}
}
