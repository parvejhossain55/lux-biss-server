package admindashboard

import (
	"github.com/gin-gonic/gin"
	"github.com/parvej/luxbiss_server/internal/common"
	"github.com/parvej/luxbiss_server/internal/logger"
)

type Handler struct {
	service Service
	log     *logger.Logger
}

func NewHandler(service Service, log *logger.Logger) *Handler {
	return &Handler{service: service, log: log}
}

func (h *Handler) GetStats(c *gin.Context) {
	stats, err := h.service.GetStats(c.Request.Context())
	if err != nil {
		common.InternalError(c, "Failed to retrieve admin dashboard stats")
		return
	}

	common.OK(c, "Admin dashboard stats retrieved successfully", stats)
}

func (h *Handler) GetRecentActivity(c *gin.Context) {
	limit := 100 // default limit
	activities, err := h.service.GetRecentActivity(c.Request.Context(), limit)
	if err != nil {
		common.InternalError(c, "Failed to retrieve recent activity")
		return
	}

	common.OK(c, "Recent activity retrieved successfully", activities)
}
