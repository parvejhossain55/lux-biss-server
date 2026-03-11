package wallet

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

func (h *Handler) Create(c *gin.Context) {
	var req CreateWalletRequest
	if errs := common.ValidateRequest(c, &req); errs != nil {
		common.BadRequest(c, "Validation failed", errs)
		return
	}

	wallet, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		if appErr, ok := common.IsAppError(err); ok {
			c.JSON(appErr.StatusCode, common.Response{
				Success:   false,
				Message:   appErr.Message,
				RequestID: c.GetString("request_id"),
			})
			return
		}
		common.InternalError(c, "Failed to create wallet")
		return
	}

	common.Created(c, "Wallet created successfully", ToResponse(wallet))
}

func (h *Handler) List(c *gin.Context) {
	wallets, err := h.service.List(c.Request.Context())
	if err != nil {
		common.InternalError(c, "Failed to list wallets")
		return
	}

	common.OK(c, "Wallets retrieved successfully", ToResponseList(wallets))
}

func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		common.BadRequest(c, "Invalid wallet ID", nil)
		return
	}

	var req UpdateWalletRequest
	if errs := common.ValidateRequest(c, &req); errs != nil {
		common.BadRequest(c, "Validation failed", errs)
		return
	}

	wallet, err := h.service.Update(c.Request.Context(), id, &req)
	if err != nil {
		if appErr, ok := common.IsAppError(err); ok {
			c.JSON(appErr.StatusCode, common.Response{
				Success:   false,
				Message:   appErr.Message,
				RequestID: c.GetString("request_id"),
			})
			return
		}
		common.InternalError(c, "Failed to update wallet")
		return
	}

	common.OK(c, "Wallet updated successfully", ToResponse(wallet))
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		common.BadRequest(c, "Invalid wallet ID", nil)
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		if appErr, ok := common.IsAppError(err); ok {
			c.JSON(appErr.StatusCode, common.Response{
				Success:   false,
				Message:   appErr.Message,
				RequestID: c.GetString("request_id"),
			})
			return
		}
		common.InternalError(c, "Failed to delete wallet")
		return
	}

	common.NoContent(c)
}
