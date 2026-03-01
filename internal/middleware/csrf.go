package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/parvej/luxbiss_server/internal/config"
)

// CSRF validates the Origin/Referer header on state-changing requests (POST, PUT, PATCH, DELETE).
// Safe methods (GET, HEAD, OPTIONS) are always allowed through.
//
// This works alongside SameSite=Lax cookies as a defense-in-depth layer:
//   - SameSite=Lax → browser-side protection (browser won't send cookies cross-site)
//   - Origin validation → server-side protection (server rejects requests from unknown origins)
func CSRF(corsConfig *config.CORSConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Safe methods don't change state — always allow
		if isSafeMethod(c.Request.Method) {
			c.Next()
			return
		}

		// Non-browser requests (API clients, Postman, curl) don't send Origin.
		// These are safe because CSRF only exploits browsers.
		origin := c.GetHeader("Origin")
		referer := c.GetHeader("Referer")

		// If neither Origin nor Referer is present, this is likely
		// a non-browser client (Postman, curl, mobile app) — allow it.
		// Browsers ALWAYS send at least one of these on POST/PUT/DELETE.
		if origin == "" && referer == "" {
			c.Next()
			return
		}

		// Validate Origin header (preferred, more reliable)
		if origin != "" {
			if isAllowedOrigin(origin, corsConfig.AllowedOrigins) {
				c.Next()
				return
			}

			// In debug mode, allow localhost origins for development
			if gin.Mode() == gin.DebugMode && isLocalhostOrigin(origin) {
				c.Next()
				return
			}

			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success":    false,
				"message":    "CSRF validation failed: origin not allowed",
				"request_id": c.GetString("request_id"),
			})
			return
		}

		// Fallback to Referer header (some browsers strip Origin on same-origin)
		if referer != "" {
			if isAllowedReferer(referer, corsConfig.AllowedOrigins) {
				c.Next()
				return
			}

			// In debug mode, allow localhost referers for development
			if gin.Mode() == gin.DebugMode && isLocalhostOrigin(referer) {
				c.Next()
				return
			}

			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"success":    false,
				"message":    "CSRF validation failed: referer not allowed",
				"request_id": c.GetString("request_id"),
			})
			return
		}
	}
}

// isSafeMethod returns true for HTTP methods that should not change state
func isSafeMethod(method string) bool {
	return method == http.MethodGet ||
		method == http.MethodHead ||
		method == http.MethodOptions
}

// isAllowedOrigin checks if the origin matches any allowed origin
func isAllowedOrigin(origin string, allowedOrigins []string) bool {
	for _, allowed := range allowedOrigins {
		if strings.TrimSpace(allowed) == origin {
			return true
		}
	}
	return false
}

// isAllowedReferer checks if the referer starts with any allowed origin
func isAllowedReferer(referer string, allowedOrigins []string) bool {
	for _, allowed := range allowedOrigins {
		trimmed := strings.TrimSpace(allowed)
		if strings.HasPrefix(referer, trimmed) {
			return true
		}
	}
	return false
}

// isLocalhostOrigin checks if the origin is from localhost (for development)
func isLocalhostOrigin(origin string) bool {
	return strings.HasPrefix(origin, "http://localhost:") ||
		strings.HasPrefix(origin, "http://127.0.0.1:") ||
		origin == "http://localhost" ||
		origin == "http://127.0.0.1"
}
