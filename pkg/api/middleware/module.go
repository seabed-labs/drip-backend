package middleware

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/dcaf-labs/drip/pkg/apispec"
	"github.com/dcaf-labs/drip/pkg/configs"
	"github.com/dcaf-labs/drip/pkg/repository"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
	"google.golang.org/api/idtoken"
)

type Handler struct {
	googleClientID string
	repo           repository.Repository
	rateLimiter    *limiter.Limiter
}

func NewHandler(
	config *configs.AppConfig,
	repo repository.Repository,
) *Handler {
	// 10 requests / second
	rate, err := limiter.NewRateFromFormatted("10-S")
	if err != nil {
		panic(err)
	}
	store := memory.NewStore()
	rateLimiter := limiter.New(store, rate)

	return &Handler{
		googleClientID: config.GoogleClientID,
		repo:           repo,
		rateLimiter:    rateLimiter,
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
			return c.JSON(http.StatusUnauthorized, apispec.ErrorResponse{Error: "invalid token-id"})
		}
		logrus.
			WithField("email", payload.Claims["email"]).
			WithField("name", payload.Claims["name"]).
			Info("authorized")
		return next(c)
	}
}

func (h *Handler) RateLimit(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ip := c.RealIP()
		limiterCtx, err := h.rateLimiter.Get(c.Request().Context(), ip)
		log := logrus.
			WithField("ip", ip).
			WithField("url", c.Request().URL)
		if err != nil {
			log.
				WithError(err).
				Info("IPRateLimit - ipRateLimiter.Get")
			return c.JSON(http.StatusInternalServerError, echo.Map{
				"success": false,
				"message": err,
			})
		}

		h := c.Response().Header()
		h.Set("X-RateLimit-Limit", strconv.FormatInt(limiterCtx.Limit, 10))
		h.Set("X-RateLimit-Remaining", strconv.FormatInt(limiterCtx.Remaining, 10))
		h.Set("X-RateLimit-Reset", strconv.FormatInt(limiterCtx.Reset, 10))

		if limiterCtx.Reached {
			log.Info("Too Many Requests from")
			return c.JSON(http.StatusTooManyRequests, apispec.ErrorResponse{
				Error: "Too Many Requests on " + c.Request().URL.String(),
			})
		}
		return next(c)
	}
}
