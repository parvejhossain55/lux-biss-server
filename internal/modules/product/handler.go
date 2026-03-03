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

	common.NoContent(c)
}

func (h *Handler) ListLevels(c *gin.Context) {
	levels, err := h.service.ListLevels(c.Request.Context())
	if err != nil {
		common.InternalError(c, "Failed to list levels")
		return
	}

	common.OK(c, "Levels retrieved successfully", ToLevelResponseList(levels))
}

func (h *Handler) ListStepsByLevel(c *gin.Context) {
	levelIDStr := c.Param("level_id")
	levelID, err := strconv.ParseUint(levelIDStr, 10, 32)
	if err != nil {
		common.BadRequest(c, "Invalid level ID", nil)
		return
	}

	steps, err := h.service.ListStepsByLevel(c.Request.Context(), uint(levelID))
	if err != nil {
		common.InternalError(c, "Failed to list steps")
		return
	}

	common.OK(c, "Steps retrieved successfully", ToStepResponseList(steps))
}
