package handlers

import (
	"fmt"
	"html"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/supabase-community/gotrue-go/types"
	"github.com/supabase-community/supabase-go"
	"share.dev/internal"
	"share.dev/routes"
)

const (
	accessTokenCookie  = "__Host_aid"
	refreshTokenCookie = "__Host_rtk"
)

type Credentials struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

// Add this func to your routes
func AuthRequired(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie(accessTokenCookie)
		if err != nil {
			return c.Redirect(http.StatusSeeOther, routes.Login)
		}

		token, _, err := new(jwt.Parser).ParseUnverified(cookie.Value, jwt.MapClaims{})
		if err != nil {
			return c.Redirect(http.StatusSeeOther, routes.Login)
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Redirect(http.StatusSeeOther, routes.Login)
		}

		email, _ := claims["email"].(string)
		if email == "" {
			return c.Redirect(http.StatusSeeOther, routes.Login)
		}

		c.Set("user_email", email)

		return next(c)
	}
}

//* Auth Related Routes

// POST /login formdata
func Login(client *supabase.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Only Htmx requests allowed
		if c.Request().Header.Get("HX-Request") != "true" {
			return c.JSON(http.StatusBadRequest, "Invalid request")
		}

		// Get form data manually since you're using form fields "login" and "password"
		loginField := c.FormValue("login") // This handles email or username
		password := c.FormValue("password")

		// Validate input
		if loginField == "" || password == "" {
			errorHTML := `
				<div class="notification is-danger is-light">
					<button class="delete" onclick="this.parentElement.style.display='none'"></button>
					<strong>Validation Error:</strong> Email/username and password are required.
				</div>
			`
			return c.HTML(400, errorHTML)
		}

		// Attempt authentication with Supabase
		session, err := client.Auth.SignInWithEmailPassword(loginField, password)
		if err != nil {
			errorHTML := fmt.Sprintf(`
				<div class="notification is-danger is-light">
					<button class="delete" onclick="this.parentElement.style.display='none'"></button>
					<strong>Login Failed:</strong> %s
				</div>
			`, html.EscapeString(err.Error()))
			return c.HTML(401, errorHTML)
		}

		if session == nil {
			errorHTML := `
				<div class="notification is-danger is-light">
					<button class="delete" onclick="this.parentElement.style.display='none'"></button>
					<strong>Login Failed:</strong> Authentication session could not be created.
				</div>
			`
			return c.HTML(401, errorHTML)
		}

		// Set authentication cookies
		setAuthCookies(c, session.AccessToken, session.RefreshToken)

		// SUCCESS: Full page redirect using HTMX
		c.Response().Header().Set("HX-Redirect", routes.MainPage)
		return c.NoContent(200)
	}
}

// POST /logout ~ calls Supabase Auth.Logout and clears auth cookies
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

		return c.Redirect(http.StatusSeeOther, routes.IndexPage)
	}
}

// POST /signup formdata
func Signup(client *supabase.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		// Only Htmx requests allowed

		if c.Request().Header.Get("HX-Request") != "true" {
			return c.JSON(http.StatusBadRequest, "Invalid request")
		}

		var creds Credentials
		if err := c.Bind(&creds); err != nil {
			return c.String(http.StatusBadRequest, "Invalid signup data")
		}
		cc, ok := c.(internal.CustomContext)
		if !ok {
			return echo.NewHTTPError(500, "an error occurred when casting the echo.Context to the CustomContext object, in Signup func.")
		}
		signupResp, err := client.Auth.Signup(types.SignupRequest{
			Email:    creds.Email,
			Password: creds.Password,
			// Optionally,
			// EmailRedirectTo: "https://yourapp.com/welcome", // for email confirmation redirect
		})

		if err != nil || signupResp == nil {
			return echo.NewHTTPError(http.StatusConflict, "Signup failed: "+err.Error())
		}

		if cc.IsDev() {
			setAuthCookies(c, signupResp.AccessToken, signupResp.RefreshToken)
			return c.Redirect(http.StatusSeeOther, routes.MainPage)
		}

		return c.Redirect(http.StatusOK, routes.CheckEmailPage)
	}
}

func Verify(client *supabase.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		token := c.QueryParam("token")
		email := c.QueryParam("email")

		some, err := client.Auth.VerifyForUser(types.VerifyForUserRequest{
			Type:  types.VerificationTypeSignup,
			Token: token,
			Email: email,
		})

		if err != nil {
			return echo.NewHTTPError(http.StatusConflict, "Verify failed: "+err.Error())
		}

		setAuthCookies(c, some.AccessToken, some.RefreshToken)
		return c.Redirect(http.StatusOK, routes.MainPage)
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

// func ValidateSignup(email, pass string) error {
// 	if email != "" {

// 		email = bluemonday.StrictPolicy().Sanitize(email)
// 		if email == "" {
// 			return errors.New("Something went wrong sanitizing email string")
// 		}
// 	}
// 	if email != "" {
// 		pass = bluemonday.StrictPolicy().Sanitize(pass)
// 		if pass == "" {
// 			return errors.New("Something went wront sanitizing password string")
// 		}
// 	}
// 	return nil
// }
