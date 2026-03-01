package auth

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/parvej/luxbiss_server/internal/config"
)

// CookieManager handles setting and clearing auth cookies
type CookieManager struct {
	domain          string
	secure          bool
	sameSite        http.SameSite
	path            string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

// NewCookieManager creates a CookieManager from config
func NewCookieManager(cookieCfg *config.CookieConfig, jwtCfg *config.JWTConfig) *CookieManager {
	sameSite := http.SameSiteLaxMode
	switch strings.ToLower(cookieCfg.SameSite) {
	case "strict":
		sameSite = http.SameSiteStrictMode
	case "none":
		sameSite = http.SameSiteNoneMode
	case "lax":
		sameSite = http.SameSiteLaxMode
	}

	return &CookieManager{
		domain:          cookieCfg.Domain,
		secure:          cookieCfg.Secure,
		sameSite:        sameSite,
		path:            cookieCfg.Path,
		accessTokenTTL:  jwtCfg.AccessTokenTTL,
		refreshTokenTTL: jwtCfg.RefreshTokenTTL,
	}
}

// SetTokenCookies sets both access and refresh token cookies on the response
func (cm *CookieManager) SetTokenCookies(c *gin.Context, accessToken, refreshToken string) {
	// Access token cookie — sent with every request
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		MaxAge:   int(cm.accessTokenTTL.Seconds()),
		Path:     cm.path,
		Domain:   cm.domain,
		Secure:   cm.secure,
		HttpOnly: true,
		SameSite: cm.sameSite,
	})

	// Refresh token cookie — restricted to refresh endpoint only
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		MaxAge:   int(cm.refreshTokenTTL.Seconds()),
		Path:     "/api/v1/auth/refresh",
		Domain:   cm.domain,
		Secure:   cm.secure,
		HttpOnly: true,
		SameSite: cm.sameSite,
	})
}

// ClearTokenCookies clears both auth cookies (used for logout)
func (cm *CookieManager) ClearTokenCookies(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		MaxAge:   -1,
		Path:     cm.path,
		Domain:   cm.domain,
		Secure:   cm.secure,
		HttpOnly: true,
		SameSite: cm.sameSite,
	})

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		MaxAge:   -1,
		Path:     "/api/v1/auth/refresh",
		Domain:   cm.domain,
		Secure:   cm.secure,
		HttpOnly: true,
		SameSite: cm.sameSite,
	})
}
