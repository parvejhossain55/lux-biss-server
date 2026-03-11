package transaction

import (
	"strconv"

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
	requestingUserID := c.GetString("user_id")
	requestingRole := c.GetString("user_role")

	var req CreateTransactionRequest
	if errs := common.ValidateRequest(c, &req); errs != nil {
		common.BadRequest(c, "Validation failed", errs)
		return
	}

	tx, err := h.service.Create(c.Request.Context(), &req, requestingUserID, requestingRole)
	if err != nil {
		if appErr, ok := common.IsAppError(err); ok {
			c.JSON(appErr.StatusCode, common.Response{
				Success:   false,
				Message:   appErr.Message,
				RequestID: c.GetString("request_id"),
			})
			return
		}
		common.InternalError(c, "Failed to create transaction")
		return
	}

	common.Created(c, "Transaction created successfully", ToResponse(tx))
}

func (h *Handler) GetByID(c *gin.Context) {
	requestingUserID := c.GetString("user_id")
	requestingRole := c.GetString("user_role")

	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		common.BadRequest(c, "Invalid transaction ID", nil)
		return
	}

	tx, err := h.service.GetByID(c.Request.Context(), id, requestingUserID, requestingRole)
	if err != nil {
		if appErr, ok := common.IsAppError(err); ok {
			c.JSON(appErr.StatusCode, common.Response{
				Success:   false,
				Message:   appErr.Message,
				RequestID: c.GetString("request_id"),
			})
			return
		}
		common.InternalError(c, "Failed to get transaction")
		return
	}

	common.OK(c, "Transaction retrieved successfully", ToResponse(tx))
}

func (h *Handler) List(c *gin.Context) {
	requestingUserID := c.GetString("user_id")
	requestingRole := c.GetString("user_role")

	pagination := common.NewPagination(c)

	userID := c.Query("user_id")
	txType := c.Query("type")
	status := c.Query("status")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "DESC")

	txs, total, err := h.service.List(c.Request.Context(), userID, txType, status, pagination.PerPage, pagination.Offset, sortBy, sortOrder, requestingUserID, requestingRole)
	if err != nil {
		common.InternalError(c, "Failed to list transactions")
		return
	}

	common.OKWithMeta(c, "Transactions retrieved successfully", ToResponseList(txs), pagination.ToMeta(total))
}

func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		common.BadRequest(c, "Invalid transaction ID", nil)
		return
	}

	var req UpdateTransactionRequest
	if errs := common.ValidateRequest(c, &req); errs != nil {
		common.BadRequest(c, "Validation failed", errs)
		return
	}

	requestingRole := c.GetString("user_role")
	tx, err := h.service.Update(c.Request.Context(), id, &req, requestingRole)
	if err != nil {
		if appErr, ok := common.IsAppError(err); ok {
			c.JSON(appErr.StatusCode, common.Response{
				Success:   false,
				Message:   appErr.Message,
				RequestID: c.GetString("request_id"),
			})
			return
		}
		common.InternalError(c, "Failed to update transaction")
		return
	}

	common.OK(c, "Transaction updated successfully", ToResponse(tx))
}

// Only admin can delete full records
func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		common.BadRequest(c, "Invalid transaction ID", nil)
		return
	}

	requestingRole := c.GetString("user_role")
	if err := h.service.Delete(c.Request.Context(), id, requestingRole); err != nil {
		if appErr, ok := common.IsAppError(err); ok {
			c.JSON(appErr.StatusCode, common.Response{
				Success:   false,
				Message:   appErr.Message,
				RequestID: c.GetString("request_id"),
			})
			return
		}
		common.InternalError(c, "Failed to delete transaction")
		return
	}

	common.NoContent(c)
}

func (h *Handler) GetSummary(c *gin.Context) {
	requestingUserID := c.GetString("user_id")
	requestingRole := c.GetString("user_role")

	targetUserID := c.Query("user_id")

	// Regular users can only see their own summary
	if requestingRole != "admin" {
		targetUserID = requestingUserID
	} else if requestingRole == "admin" && targetUserID == "" {
		targetUserID = requestingUserID // fallback if no ID provided by admin
	}

	daysStr := c.DefaultQuery("days", "7")
	days, err := strconv.Atoi(daysStr)
	if err != nil || days <= 0 {
		days = 7
	}

	summary, err := h.service.GetSummary(c.Request.Context(), targetUserID, days)
	if err != nil {
		common.InternalError(c, "Failed to compute transaction summary")
		return
	}

	common.OK(c, "Transaction summary computed", summary)
}

func (h *Handler) Invest(c *gin.Context) {
	userID := c.GetString("user_id")

	var req InvestRequest
	if errs := common.ValidateRequest(c, &req); errs != nil {
		common.BadRequest(c, "Validation failed", errs)
		return
	}

	if err := h.service.Invest(c.Request.Context(), userID, &req); err != nil {
		if appErr, ok := common.IsAppError(err); ok {
			c.JSON(appErr.StatusCode, common.Response{
				Success:   false,
				Message:   appErr.Message,
				RequestID: c.GetString("request_id"),
			})
			return
		}
		common.InternalError(c, err.Error())
		return
	}

	common.OK(c, "Investment completed successfully", nil)
}
