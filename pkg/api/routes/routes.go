package controller

import (
	"github.com/dcaf-labs/drip/pkg/service/base"
	"github.com/dcaf-labs/drip/pkg/service/clients/solana"
	"github.com/dcaf-labs/drip/pkg/service/config"
	"github.com/dcaf-labs/drip/pkg/service/repository"
	analyticsrepository "github.com/dcaf-labs/drip/pkg/service/repository/analytics"
	"github.com/gorilla/schema"
)

const defaultLimit = 500

type Handler struct {
	decoder       *schema.Decoder
	base          base.Base
	solanaClient  solana.Solana
	repo          repository.Repository
	analyticsRepo analyticsrepository.AnalyticsRepository
	network       config.Network
	env           config.Environment
	port          int
}

func NewHandler(
	appConfig config.AppConfig,
	solanaClient solana.Solana,
	base base.Base,
	repo repository.Repository,
	analyticsRepo analyticsrepository.AnalyticsRepository,
) *Handler {
	return &Handler{
		decoder:       schema.NewDecoder(),
		solanaClient:  solanaClient,
		base:          base,
		repo:          repo,
		analyticsRepo: analyticsRepo,
		network:       appConfig.GetNetwork(),
		env:           appConfig.GetEnvironment(),
		port:          appConfig.GetServerPort(),
	}
}
