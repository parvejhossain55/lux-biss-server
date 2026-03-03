package auth

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/parvej/luxbiss_server/internal/common"
	"github.com/parvej/luxbiss_server/internal/logger"
)

type Handler struct {
	service       *Service
	cookieManager *CookieManager
	log           *logger.Logger
}

func NewHandler(service *Service, cookieManager *CookieManager, log *logger.Logger) *Handler {
	return &Handler{service: service, cookieManager: cookieManager, log: log}
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

	// Set tokens as HTTP-only cookies
	h.cookieManager.SetTokenCookies(c, resp.AccessToken, resp.RefreshToken)

	common.Created(c, "Registration successful", resp.User)
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

	// Set tokens as HTTP-only cookies
	h.cookieManager.SetTokenCookies(c, resp.AccessToken, resp.RefreshToken)

	common.OK(c, "Login successful", resp.User)
}

func (h *Handler) RefreshToken(c *gin.Context) {
	// Try to read refresh token from cookie first, then fallback to JSON body
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil || refreshToken == "" {
		var req RefreshTokenRequest
		if errs := common.ValidateRequest(c, &req); errs != nil {
			common.BadRequest(c, "Validation failed", errs)
			return
		}
		refreshToken = req.RefreshToken
	}

	resp, svcErr := h.service.RefreshToken(c.Request.Context(), refreshToken)
	if svcErr != nil {
		if appErr, ok := common.IsAppError(svcErr); ok {
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

	// Set new tokens as HTTP-only cookies
	h.cookieManager.SetTokenCookies(c, resp.AccessToken, resp.RefreshToken)

	common.OK(c, "Token refreshed successfully", resp.User)
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

	// Set tokens as HTTP-only cookies
	h.cookieManager.SetTokenCookies(c, resp.AccessToken, resp.RefreshToken)

	common.OK(c, "Google login successful", resp.User)
}

func (h *Handler) Logout(c *gin.Context) {
	// 1. Get token from cookie or header
	var tokenStr string
	if cookie, err := c.Cookie("access_token"); err == nil && cookie != "" {
		tokenStr = cookie
	} else {
		authHeader := c.GetHeader("Authorization")
		tokenStr = strings.TrimPrefix(authHeader, "Bearer ")
	}

	if tokenStr != "" {
		_ = h.service.Logout(c.Request.Context(), tokenStr)
	}

	h.cookieManager.ClearTokenCookies(c)
	common.OK(c, "Logged out successfully", nil)
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

func (h *Handler) VerifyOTP(c *gin.Context) {
	var req VerifyOTPRequest
	if errs := common.ValidateRequest(c, &req); errs != nil {
		common.BadRequest(c, "Validation failed", errs)
		return
	}

	err := h.service.VerifyOTP(c.Request.Context(), &req)
	if err != nil {
		if appErr, ok := common.IsAppError(err); ok {
			c.JSON(appErr.StatusCode, common.Response{
				Success:   false,
				Message:   appErr.Message,
				RequestID: c.GetString("request_id"),
			})
			return
		}
		common.InternalError(c, "Failed to verify OTP")
		return
	}

	common.OK(c, "OTP verified successfully", nil)
}
