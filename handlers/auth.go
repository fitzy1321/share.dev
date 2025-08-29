package handlers

import (
	"net/http"
	"time"

	"github.com/go-playground/validator"
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
	Email    string `form:"email" validate:"required,email"`
	Password string `form:"password" validate:"required,min=8"`
}

var validate *validator.Validate = validator.New()

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

func Login(client *supabase.Client) echo.HandlerFunc {
	return func(c echo.Context) error {
		creds := new(Credentials)
		if err := c.Bind(creds); err != nil {
			return c.String(http.StatusBadRequest, "Invalid input")
		}
		// Validate input
		if creds.Email == "" || creds.Password == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Email and Password are required.")
		}

		if err := validate.Struct(creds); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Data Vadliation Error:"+err.Error())
		}

		// TODO: more validation and sanitization needed

		// Attempt authentication with Supabase
		session, err := client.Auth.SignInWithEmailPassword(creds.Email, creds.Password)
		if err != nil {
			return echo.NewHTTPError(http.StatusConflict, "Problem with supabase login: "+err.Error())
		}

		if session == nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "Authentication session could not be created")
		}

		// Set authentication cookies
		setAuthCookies(c, session.AccessToken, session.RefreshToken)
		return c.Redirect(http.StatusOK, routes.MainPage)
	}
}

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
		creds := new(Credentials)
		if err := c.Bind(creds); err != nil {
			return c.String(http.StatusBadRequest, "Invalid input")
		}
		// Validate input
		if creds.Email == "" || creds.Password == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "Email and Password are required.")
		}

		if err := validate.Struct(creds); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Data Vadliation Error:"+err.Error())
		}

		// TODO: more validation and sanitization needed

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
		if err := validate.Var(token, "required"); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Something wront with token:"+err.Error())
		}
		if err := validate.Var(email, "required,email"); err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "Somethint wrong with the email:"+err.Error())
		}

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
