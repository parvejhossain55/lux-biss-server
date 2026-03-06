package giftcard

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
	var req CreateGiftcardRequest
	if errs := common.ValidateRequest(c, &req); errs != nil {
		common.BadRequest(c, "Validation failed", errs)
		return
	}

	giftcard, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		if appErr, ok := common.IsAppError(err); ok {
			c.JSON(appErr.StatusCode, common.Response{
				Success:   false,
				Message:   appErr.Message,
				RequestID: c.GetString("request_id"),
			})
			return
		}
		common.InternalError(c, "Failed to create giftcard")
		return
	}

	common.Created(c, "Giftcard created successfully", ToResponse(giftcard))
}

func (h *Handler) List(c *gin.Context) {
	pagination := common.NewPagination(c)

	giftcards, total, err := h.service.List(c.Request.Context(), pagination.PerPage, pagination.Offset)
	if err != nil {
		common.InternalError(c, "Failed to list giftcards")
		return
	}

	common.OKWithMeta(c, "Giftcards retrieved successfully", ToResponseList(giftcards), pagination.ToMeta(total))
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		common.BadRequest(c, "Invalid giftcard ID", nil)
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
		common.InternalError(c, "Failed to delete giftcard")
		return
	}

	common.NoContent(c)
}

func (h *Handler) Apply(c *gin.Context) {
	var req ApplyGiftcardRequest
	if errs := common.ValidateRequest(c, &req); errs != nil {
		common.BadRequest(c, "Validation failed", errs)
		return
	}

	userID := c.GetString("user_id")
	userEmail := c.GetString("user_email")

	giftcard, err := h.service.Apply(c.Request.Context(), &req, userID, userEmail)
	if err != nil {
		if appErr, ok := common.IsAppError(err); ok {
			c.JSON(appErr.StatusCode, common.Response{
				Success:   false,
				Message:   appErr.Message,
				RequestID: c.GetString("request_id"),
			})
			return
		}
		common.InternalError(c, "Failed to apply giftcard")
		return
	}

	common.OK(c, "Giftcard applied successfully", ToResponse(giftcard))
}

func (h *Handler) Verify(c *gin.Context) {
	var req VerifyGiftcardRequest
	if errs := common.ValidateRequest(c, &req); errs != nil {
		common.BadRequest(c, "Validation failed", errs)
		return
	}

	giftcard, err := h.service.Verify(c.Request.Context(), &req)
	if err != nil {
		if appErr, ok := common.IsAppError(err); ok {
			c.JSON(appErr.StatusCode, common.Response{
				Success:   false,
				Message:   appErr.Message,
				RequestID: c.GetString("request_id"),
			})
			return
		}
		common.InternalError(c, "Failed to verify giftcard")
		return
	}

	common.OK(c, "Giftcard verified successfully", ToResponse(giftcard))
}
