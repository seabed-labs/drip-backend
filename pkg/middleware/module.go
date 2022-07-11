package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/dcaf-protocol/drip/pkg/clients/solana"
	"github.com/dcaf-protocol/drip/pkg/configs"
	"github.com/dcaf-protocol/drip/pkg/repository"
	swagger "github.com/dcaf-protocol/drip/pkg/swagger"
	"github.com/labstack/echo/v4"
	"google.golang.org/api/idtoken"
)

type Handler struct {
	googleClientID string
	repo           repository.Repository
}

func NewHandler(
	config *configs.AppConfig,
	solanaClient solana.Solana,
	repo repository.Repository,
) *Handler {
	return &Handler{
		googleClientID: config.GoogleClientID,
		repo:           repo,
	}
}

func (h *Handler) ValidateAccessToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if !strings.Contains(c.Request().RequestURI, "admin") {
			return next(c)
		}
		accessToken := c.Request().Header.Get("token-id")
		payload, err := idtoken.Validate(context.Background(), accessToken, h.googleClientID)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, swagger.ErrorResponse{Error: "invalid token-id"})
		}
		logrus.
			WithField("email", payload.Claims["email"]).
			WithField("name", payload.Claims["name"]).
			Info("authorized")
		return next(c)
	}
}
