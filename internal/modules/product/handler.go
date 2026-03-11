package product

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
	var req CreateProductRequest
	if errs := common.ValidateRequest(c, &req); errs != nil {
		common.BadRequest(c, "Validation failed", errs)
		return
	}

	product, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		if appErr, ok := common.IsAppError(err); ok {
			c.JSON(appErr.StatusCode, common.Response{
				Success:   false,
				Message:   appErr.Message,
				RequestID: c.GetString("request_id"),
			})
			return
		}
		common.InternalError(c, "Failed to create product")
		return
	}

	common.Created(c, "Product created successfully", ToResponse(product))
}

func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		common.BadRequest(c, "Invalid product ID", nil)
		return
	}

	product, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if appErr, ok := common.IsAppError(err); ok {
			c.JSON(appErr.StatusCode, common.Response{
				Success:   false,
				Message:   appErr.Message,
				RequestID: c.GetString("request_id"),
			})
			return
		}
		common.InternalError(c, "Failed to get product")
		return
	}

	common.OK(c, "Product retrieved successfully", ToResponse(product))
}

func (h *Handler) List(c *gin.Context) {
	pagination := common.NewPagination(c)
	sortBy := c.Query("sort_by")
	order := c.Query("order")

	// Map JSON field names to DB column names if needed
	sortFieldMapping := map[string]string{
		"price":      "price",
		"rating":     "rating",
		"name":       "name",
		"created_at": "created_at",
	}

	dbSortBy := ""
	if sortBy != "" {
		if field, ok := sortFieldMapping[sortBy]; ok {
			dbSortBy = field
		}
	}

	// Normalize order
	if order != "desc" {
		order = "asc"
	}

	levelIDStr := c.Query("level_id")
	var levelID uint
	if levelIDStr != "" {
		if val, err := strconv.ParseUint(levelIDStr, 10, 32); err == nil {
			levelID = uint(val)
		}
	}

	stepIDStr := c.Query("step_id")
	var stepID uint
	if stepIDStr != "" {
		if val, err := strconv.ParseUint(stepIDStr, 10, 32); err == nil {
			stepID = uint(val)
		}
	}

	products, total, err := h.service.List(c.Request.Context(), pagination.PerPage, pagination.Offset, dbSortBy, order, levelID, stepID)
	if err != nil {
		common.InternalError(c, "Failed to list products")
		return
	}

	common.OKWithMeta(c, "Products retrieved successfully", ToResponseList(products), pagination.ToMeta(total))
}

func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		common.BadRequest(c, "Invalid product ID", nil)
		return
	}

	var req UpdateProductRequest
	if errs := common.ValidateRequest(c, &req); errs != nil {
		common.BadRequest(c, "Validation failed", errs)
		return
	}

	product, err := h.service.Update(c.Request.Context(), id, &req)
	if err != nil {
		if appErr, ok := common.IsAppError(err); ok {
			c.JSON(appErr.StatusCode, common.Response{
				Success:   false,
				Message:   appErr.Message,
				RequestID: c.GetString("request_id"),
			})
			return
		}
		common.InternalError(c, "Failed to update product")
		return
	}

	common.OK(c, "Product updated successfully", ToResponse(product))
}

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if _, err := uuid.Parse(id); err != nil {
		common.BadRequest(c, "Invalid product ID", nil)
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
		common.InternalError(c, "Failed to delete product")
		return
	}

	common.OK(c, "Deleted successfully", nil)
}

func (h *Handler) ListLevels(c *gin.Context) {
	pagination := common.NewPagination(c)
	levels, total, err := h.service.ListLevels(c.Request.Context(), pagination.PerPage, pagination.Offset)
	if err != nil {
		common.InternalError(c, "Failed to list levels")
		return
	}

	common.OKWithMeta(c, "Levels retrieved successfully", ToLevelResponseList(levels), pagination.ToMeta(total))
}

func (h *Handler) CreateLevel(c *gin.Context) {
	var req CreateLevelRequest
	if errs := common.ValidateRequest(c, &req); errs != nil {
		common.BadRequest(c, "Validation failed", errs)
		return
	}

	level, err := h.service.CreateLevel(c.Request.Context(), &req)
	if err != nil {
		common.InternalError(c, "Failed to create level")
		return
	}

	common.Created(c, "Level created successfully", level)
}

func (h *Handler) UpdateLevel(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		common.BadRequest(c, "Invalid level ID", nil)
		return
	}

	var req UpdateLevelRequest
	if errs := common.ValidateRequest(c, &req); errs != nil {
		common.BadRequest(c, "Validation failed", errs)
		return
	}

	level, err := h.service.UpdateLevel(c.Request.Context(), uint(id), &req)
	if err != nil {
		common.InternalError(c, "Failed to update level")
		return
	}

	common.OK(c, "Level updated successfully", level)
}

func (h *Handler) DeleteLevel(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		common.BadRequest(c, "Invalid level ID", nil)
		return
	}

	if err := h.service.DeleteLevel(c.Request.Context(), uint(id)); err != nil {
		common.InternalError(c, "Failed to delete level")
		return
	}

	common.OK(c, "Deleted successfully", nil)
}

func (h *Handler) ListStepsByLevel(c *gin.Context) {
	levelIDStr := c.Param("level_id")
	levelID, err := strconv.ParseUint(levelIDStr, 10, 32)
	if err != nil {
		common.BadRequest(c, "Invalid level ID", nil)
		return
	}

	pagination := common.NewPagination(c)
	steps, total, err := h.service.ListStepsByLevel(c.Request.Context(), uint(levelID), pagination.PerPage, pagination.Offset)
	if err != nil {
		common.InternalError(c, "Failed to list steps")
		return
	}

	common.OKWithMeta(c, "Steps retrieved successfully", ToStepResponseList(steps), pagination.ToMeta(total))
}

func (h *Handler) CreateStep(c *gin.Context) {
	var req CreateStepRequest
	if errs := common.ValidateRequest(c, &req); errs != nil {
		common.BadRequest(c, "Validation failed", errs)
		return
	}

	step, err := h.service.CreateStep(c.Request.Context(), &req)
	if err != nil {
		common.InternalError(c, "Failed to create step")
		return
	}

	common.Created(c, "Step created successfully", step)
}

func (h *Handler) UpdateStep(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		common.BadRequest(c, "Invalid step ID", nil)
		return
	}

	var req UpdateStepRequest
	if errs := common.ValidateRequest(c, &req); errs != nil {
		common.BadRequest(c, "Validation failed", errs)
		return
	}

	step, err := h.service.UpdateStep(c.Request.Context(), uint(id), &req)
	if err != nil {
		common.InternalError(c, "Failed to update step")
		return
	}

	common.OK(c, "Step updated successfully", step)
}

func (h *Handler) DeleteStep(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		common.BadRequest(c, "Invalid step ID", nil)
		return
	}

	if err := h.service.DeleteStep(c.Request.Context(), uint(id)); err != nil {
		common.InternalError(c, "Failed to delete step")
		return
	}

	common.OK(c, "Deleted successfully", nil)
}
