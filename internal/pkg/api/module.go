package api

import (
	"github.com/dcaf-protocol/drip/internal/pkg/clients/solana"
	"github.com/dcaf-protocol/drip/internal/pkg/configs"
	"github.com/gorilla/schema"
)

type Handler struct {
	decoder      *schema.Decoder
	solanaClient solana.Solana
	env          configs.Environment
	vaultConfigs []configs.VaultConfig
	port         int
}

func NewHandler(
	solanaClient solana.Solana,
	config *configs.Config,
) *Handler {
	return &Handler{
		decoder:      schema.NewDecoder(),
		solanaClient: solanaClient,
		env:          config.Environment,
		vaultConfigs: config.VaultConfigs,
		port:         config.Port,
	}
}
