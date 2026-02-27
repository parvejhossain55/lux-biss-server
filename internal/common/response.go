package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Meta struct {
	Page    int   `json:"page"`
	PerPage int   `json:"per_page"`
	Total   int64 `json:"total"`
}

type Response struct {
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Meta      *Meta       `json:"meta,omitempty"`
	Errors    interface{} `json:"errors,omitempty"`
	RequestID string      `json:"request_id,omitempty"`
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func OK(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Success:   true,
		Message:   message,
		Data:      data,
		RequestID: c.GetString("request_id"),
	})
}

func OKWithMeta(c *gin.Context, message string, data interface{}, meta *Meta) {
	c.JSON(http.StatusOK, Response{
		Success:   true,
		Message:   message,
		Data:      data,
		Meta:      meta,
		RequestID: c.GetString("request_id"),
	})
}

func Created(c *gin.Context, message string, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Success:   true,
		Message:   message,
		Data:      data,
		RequestID: c.GetString("request_id"),
	})
}

func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

func BadRequest(c *gin.Context, message string, errs interface{}) {
	c.JSON(http.StatusBadRequest, Response{
		Success:   false,
		Message:   message,
		Errors:    errs,
		RequestID: c.GetString("request_id"),
	})
}

func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, Response{
		Success:   false,
		Message:   message,
		RequestID: c.GetString("request_id"),
	})
}

func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, Response{
		Success:   false,
		Message:   message,
		RequestID: c.GetString("request_id"),
	})
}

func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, Response{
		Success:   false,
		Message:   message,
		RequestID: c.GetString("request_id"),
	})
}

func Conflict(c *gin.Context, message string) {
	c.JSON(http.StatusConflict, Response{
		Success:   false,
		Message:   message,
		RequestID: c.GetString("request_id"),
	})
}

func InternalError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, Response{
		Success:   false,
		Message:   message,
		RequestID: c.GetString("request_id"),
	})
}
