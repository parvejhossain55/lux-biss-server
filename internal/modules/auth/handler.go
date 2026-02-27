package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/parvej/luxbiss_server/internal/common"
	"github.com/parvej/luxbiss_server/internal/logger"
)

type Handler struct {
	service *Service
	log     *logger.Logger
}

func NewHandler(service *Service, log *logger.Logger) *Handler {
	return &Handler{service: service, log: log}
}

func (h *Handler) Register(c *gin.Context) {
	var req RegisterRequest
	if errs := common.ValidateRequest(c, &req); errs != nil {
		common.BadRequest(c, "Validation failed", errs)
		return
	}

	resp, err := h.service.Register(c.Request.Context(), &req)
	if err != nil {
		if appErr, ok := common.IsAppError(err); ok {
			c.JSON(appErr.StatusCode, common.Response{
				Success:   false,
				Message:   appErr.Message,
				RequestID: c.GetString("request_id"),
			})
			return
		}
		common.InternalError(c, "Failed to register")
		return
	}

	common.Created(c, "Registration successful", resp)
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if errs := common.ValidateRequest(c, &req); errs != nil {
		common.BadRequest(c, "Validation failed", errs)
		return
	}

	resp, err := h.service.Login(c.Request.Context(), &req)
	if err != nil {
		if appErr, ok := common.IsAppError(err); ok {
			c.JSON(appErr.StatusCode, common.Response{
				Success:   false,
				Message:   appErr.Message,
				RequestID: c.GetString("request_id"),
			})
			return
		}
		common.InternalError(c, "Failed to login")
		return
	}

	common.OK(c, "Login successful", resp)
}

func (h *Handler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if errs := common.ValidateRequest(c, &req); errs != nil {
		common.BadRequest(c, "Validation failed", errs)
		return
	}

	resp, err := h.service.RefreshToken(c.Request.Context(), &req)
	if err != nil {
		if appErr, ok := common.IsAppError(err); ok {
			c.JSON(appErr.StatusCode, common.Response{
				Success:   false,
				Message:   appErr.Message,
				RequestID: c.GetString("request_id"),
			})
			return
		}
		common.InternalError(c, "Failed to refresh token")
		return
	}

	common.OK(c, "Token refreshed successfully", resp)
}

func (h *Handler) GoogleLogin(c *gin.Context) {
	var req GoogleOAuthRequest
	if errs := common.ValidateRequest(c, &req); errs != nil {
		common.BadRequest(c, "Validation failed", errs)
		return
	}

	resp, err := h.service.GoogleLogin(c.Request.Context(), &req)
	if err != nil {
		if appErr, ok := common.IsAppError(err); ok {
			c.JSON(appErr.StatusCode, common.Response{
				Success:   false,
				Message:   appErr.Message,
				RequestID: c.GetString("request_id"),
			})
			return
		}
		common.InternalError(c, "Failed to login with Google")
		return
	}

	common.OK(c, "Google login successful", resp)
}

func (h *Handler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if errs := common.ValidateRequest(c, &req); errs != nil {
		common.BadRequest(c, "Validation failed", errs)
		return
	}

	err := h.service.ForgotPassword(c.Request.Context(), &req)
	if err != nil {
		if appErr, ok := common.IsAppError(err); ok {
			c.JSON(appErr.StatusCode, common.Response{
				Success:   false,
				Message:   appErr.Message,
				RequestID: c.GetString("request_id"),
			})
			return
		}
		common.InternalError(c, "Failed to process forgot password request")
		return
	}

	common.OK(c, "If an account exists with this email, you will receive an OTP", nil)
}

func (h *Handler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if errs := common.ValidateRequest(c, &req); errs != nil {
		common.BadRequest(c, "Validation failed", errs)
		return
	}

	err := h.service.ResetPassword(c.Request.Context(), &req)
	if err != nil {
		if appErr, ok := common.IsAppError(err); ok {
			c.JSON(appErr.StatusCode, common.Response{
				Success:   false,
				Message:   appErr.Message,
				RequestID: c.GetString("request_id"),
			})
			return
		}
		common.InternalError(c, "Failed to reset password")
		return
	}

	common.OK(c, "Password reset successful", nil)
}
