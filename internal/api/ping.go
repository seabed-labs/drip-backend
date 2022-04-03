package api

import (
	"github.com/gin-gonic/gin"
)

type PingResponse struct {
	Message string `json:"message" example:"pong"`
}

// Ping godoc
// @Summary  Health Check
// @Schemes
// @Description  do ping
// @Tags         info
// @Accept       json
// @Produce      json
// @Success      200  {object}  api.PingResponse
// @Router       / [get]
// @Router       /ping [get]
func (h Handler) Ping(
	c *gin.Context,
) {
	c.JSON(200, PingResponse{
		Message: "pong",
	})
}
