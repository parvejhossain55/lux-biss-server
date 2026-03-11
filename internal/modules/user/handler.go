package user

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
	var req CreateUserRequest
	if errs := common.ValidateRequest(c, &req); errs != nil {
		common.BadRequest(c, "Validation failed", errs)
		return
	}

	user, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		if appErr, ok := common.IsAppError(err); ok {
			c.JSON(appErr.StatusCode, common.Response{
				Success:   false,
				Message:   appErr.Message,
				RequestID: c.GetString("request_id"),
			})
			return
		}
		common.InternalError(c, "Failed to create user")
		return
	}

	common.Created(c, "User created successfully", ToResponse(user))
}

func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		common.BadRequest(c, "Invalid user ID", nil)
		return
	}

	userID := c.GetString("user_id")
	userRole := c.GetString("user_role")

	if userRole != RoleAdmin && id != userID {
		common.Forbidden(c, "You don't have permission to access this user")
		return
	}

	user, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if appErr, ok := common.IsAppError(err); ok {
			c.JSON(appErr.StatusCode, common.Response{
				Success:   false,
				Message:   appErr.Message,
				RequestID: c.GetString("request_id"),
			})
			return
		}
		common.InternalError(c, "Failed to get user")
		return
	}

	common.OK(c, "User retrieved successfully", ToResponse(user))
}

func (h *Handler) List(c *gin.Context) {
	pagination := common.NewPagination(c)

	users, total, err := h.service.List(c.Request.Context(), pagination.PerPage, pagination.Offset)
	if err != nil {
		common.InternalError(c, "Failed to list users")
		return
	}

	common.OKWithMeta(c, "Users retrieved successfully", ToResponseList(users), pagination.ToMeta(total))
}

func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		common.BadRequest(c, "Invalid user ID", nil)
		return
	}

	userID := c.GetString("user_id")
	userRole := c.GetString("user_role")

	if userRole != RoleAdmin && id != userID {
		common.Forbidden(c, "You don't have permission to update this user")
		return
	}

	var req UpdateUserRequest
	if errs := common.ValidateRequest(c, &req); errs != nil {
		common.BadRequest(c, "Validation failed", errs)
		return
	}

	if userRole != RoleAdmin && (req.Role != nil || req.Status != nil) {
		common.Forbidden(c, "Only admins can update role or status")
		return
	}

	user, err := h.service.Update(c.Request.Context(), id, &req)
	if err != nil {
		if appErr, ok := common.IsAppError(err); ok {
			c.JSON(appErr.StatusCode, common.Response{
				Success:   false,
				Message:   appErr.Message,
				RequestID: c.GetString("request_id"),
			})
			return
		}
		common.InternalError(c, "Failed to update user")
		return
	}

	common.OK(c, "User updated successfully", ToResponse(user))
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		common.BadRequest(c, "Invalid user ID", nil)
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
		common.InternalError(c, "Failed to delete user")
		return
	}

	common.NoContent(c)
}

func (h *Handler) ApproveHoldBalance(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		common.BadRequest(c, "Invalid user ID", nil)
		return
	}

	user, err := h.service.ApproveHoldBalance(c.Request.Context(), id)
	if err != nil {
		if appErr, ok := common.IsAppError(err); ok {
			c.JSON(appErr.StatusCode, common.Response{
				Success:   false,
				Message:   appErr.Message,
				RequestID: c.GetString("request_id"),
			})
			return
		}
		common.InternalError(c, "Failed to approve hold balance")
		return
	}

	common.OK(c, "Hold balance approved successfully", ToResponse(user))
}

func (h *Handler) GetMe(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		common.Unauthorized(c, "Authentication required")
		return
	}

	user, err := h.service.GetByID(c.Request.Context(), userID)
	if err != nil {
		if appErr, ok := common.IsAppError(err); ok {
			c.JSON(appErr.StatusCode, common.Response{
				Success:   false,
				Message:   appErr.Message,
				RequestID: c.GetString("request_id"),
			})
			return
		}
		common.InternalError(c, "Failed to get user")
		return
	}

	common.OK(c, "User retrieved successfully", ToResponse(user))
}
