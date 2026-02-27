package health

import (
	"github.com/gin-gonic/gin"
	"github.com/parvej/luxbiss_server/internal/common"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) HealthCheck(c *gin.Context) {
	common.OK(c, "Server is healthy", gin.H{
		"status": "ok",
	})
}

func (h *Handler) ReadinessCheck(c *gin.Context) {
	common.OK(c, "Server is ready", gin.H{
		"status": "ready",
	})
}
