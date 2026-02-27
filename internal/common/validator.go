package common

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func RegisterCustomValidators(registerFunc func(*validator.Validate)) {
	registerFunc(validate)
}

func ValidateRequest(c *gin.Context, req interface{}) []FieldError {
	if err := c.ShouldBindJSON(req); err != nil {
		return []FieldError{
			{Field: "body", Message: "Invalid request body"},
		}
	}

	if err := validate.Struct(req); err != nil {
		var fieldErrors []FieldError
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, e := range validationErrors {
				fieldErrors = append(fieldErrors, FieldError{
					Field:   toSnakeCase(e.Field()),
					Message: formatValidationError(e),
				})
			}
		}
		return fieldErrors
	}

	return nil
}

func formatValidationError(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", toSnakeCase(e.Field()))
	case "email":
		return "must be a valid email address"
	case "min":
		return fmt.Sprintf("must be at least %s characters", e.Param())
	case "max":
		return fmt.Sprintf("must be at most %s characters", e.Param())
	case "oneof":
		return fmt.Sprintf("must be one of: %s", e.Param())
	case "url":
		return "must be a valid URL"
	case "uuid":
		return "must be a valid UUID"
	default:
		return fmt.Sprintf("failed on '%s' validation", e.Tag())
	}
}

func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if i > 0 && r >= 'A' && r <= 'Z' {
			result.WriteRune('_')
		}
		result.WriteRune(r)
	}
	return strings.ToLower(result.String())
}
