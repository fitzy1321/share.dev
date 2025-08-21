package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"sync"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"golang.org/x/time/rate"
)

var visitors = make(map[string]*rate.Limiter)
var mu sync.Mutex

func getVisitor(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := visitors[ip]
	if !exists {
		limiter = rate.NewLimiter(5, 10) // 5 req/sec with burst 10
		visitors[ip] = limiter
	}
	return limiter
}

func RateLimitMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		ip := c.RealIP()
		limiter := getVisitor(ip)

		if !limiter.Allow() {
			return c.JSON(http.StatusTooManyRequests, map[string]string{
				"error": "Too many requests",
			})
		}
		return next(c)
	}
}

const csrfCookieName = "__Host_csrf"

func generateCSRFToken() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func CSRFMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenCookie, err := c.Cookie(csrfCookieName)
		if err != nil {
			token := generateCSRFToken()
			c.SetCookie(&http.Cookie{
				Name:     csrfCookieName,
				Value:    token,
				Path:     "/",
				HttpOnly: true,
				Secure:   true,
				SameSite: http.SameSiteStrictMode,
			})
			c.Set("csrf_token", token)
		} else {
			c.Set("csrf_token", tokenCookie.Value)
		}

		method := c.Request().Method
		if method == http.MethodPost || method == http.MethodPut || method == http.MethodDelete {
			formToken := c.FormValue("csrf_token")
			headerToken := c.Request().Header.Get("X-CSRF-Token")
			expected := c.Get("csrf_token").(string)
			if formToken != expected && headerToken != expected {
				return c.Render(http.StatusForbidden, "error.html", map[string]interface{}{
					"title": "Forbidden",
					"error": "CSRF validation failed",
				})
			}
		}
		return next(c)
	}
}

func AuthRequired(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie(accessTokenCookie)
		if err != nil {
			return c.Redirect(http.StatusSeeOther, "/login")
		}

		token, _, err := new(jwt.Parser).ParseUnverified(cookie.Value, jwt.MapClaims{})
		if err != nil {
			return c.Redirect(http.StatusSeeOther, "/login")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Redirect(http.StatusSeeOther, "/login")
		}

		email, _ := claims["email"].(string)
		if email == "" {
			return c.Redirect(http.StatusSeeOther, "/login")
		}

		c.Set("user_email", email)

		return next(c)
	}
}

func SecurityHeaders(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		h := c.Response().Header()
		h.Set("Content-Security-Policy", "default-src 'self'; style-src 'self' https://cdn.jsdelivr.net; script-src 'self' https://unpkg.com;")
		h.Set("X-Frame-Options", "DENY")
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("Referrer-Policy", "strict-origin-when-cross-origin")
		h.Set("X-XSS-Protection", "1; mode=block")
		return next(c)
	}
}
