package api

import (
	"github.com/dcaf-protocol/drip/internal/pkg/solanaclient"
	"github.com/gorilla/schema"
)

type Handler struct {
	decoder      *schema.Decoder
	solanaClient *solanaclient.Solana
}

func NewHandler(
	solanaClient *solanaclient.Solana,
) *Handler {
	return &Handler{
		decoder:      schema.NewDecoder(),
		solanaClient: solanaClient,
	}
}
