package handlers

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/supabase-community/gotrue-go/types"
	"github.com/supabase-community/supabase-go"
)

const (
	accessTokenCookie  = "__Host_aid"
	refreshTokenCookie = "__Host_rtk"
)

type Credentials struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

func Signup(client *supabase.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		var creds Credentials
		if err := c.Bind(&creds); err != nil {
			return c.String(http.StatusBadRequest, "Invalid signup data")
		}

		session, err := client.Auth.Signup(types.SignupRequest{
			Email:    creds.Email,
			Password: creds.Password,
			// Optionally,
			// EmailRedirectTo: "https://yourapp.com/welcome", // for email confirmation redirect
		})

		if err != nil || session == nil {
			return c.String(http.StatusBadRequest, "Signup failed: "+err.Error())
		}

		setAuthCookies(c, session.AccessToken, session.RefreshToken)
		return c.Redirect(http.StatusSeeOther, "/dashboard")
	}
}

func Login(client *supabase.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		var creds Credentials
		if err := c.Bind(&creds); err != nil {
			return c.String(http.StatusBadRequest, "Invalid login data")
		}

		session, err := client.Auth.SignInWithEmailPassword(creds.Email, creds.Password)
		if err != nil || session == nil {
			return c.String(http.StatusUnauthorized, "Login failed: "+err.Error())
		}

		setAuthCookies(c, session.AccessToken, session.RefreshToken)
		return c.Redirect(http.StatusSeeOther, "/dashboard")
	}
}

// Logout calls Supabase Auth.Logout and clears auth cookies
func Logout(client *supabase.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		accessCookie, err := c.Cookie(accessTokenCookie)
		if err == nil && accessCookie.Value != "" {
			err = client.Auth.Logout()
			if err != nil {
				c.Logger().Error("Problem with logout", err)
			}
		}

		clearAuthCookies(c)

		return c.Redirect(http.StatusSeeOther, "/login")
	}
}

func setAuthCookies(c echo.Context, accessToken, refreshToken string) {
	c.SetCookie(&http.Cookie{
		Name:     accessTokenCookie,
		Value:    accessToken,
		Path:     "/",
		Expires:  time.Now().Add(60 * time.Minute),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
	c.SetCookie(&http.Cookie{
		Name:     refreshTokenCookie,
		Value:    refreshToken,
		Path:     "/",
		Expires:  time.Now().Add(7 * 24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}

func clearAuthCookies(c echo.Context) {
	c.SetCookie(&http.Cookie{
		Name:     accessTokenCookie,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
	c.SetCookie(&http.Cookie{
		Name:     refreshTokenCookie,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
	})
}
