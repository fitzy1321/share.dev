package handlers

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"share.dev/templates"
)

func getCSRFToken(c echo.Context) string {
	csrfToken, ok := c.Get(middleware.DefaultCSRFConfig.ContextKey).(string)
	if !ok {
		c.Logger().Fatal("CSRF Token not found")
	}
	return csrfToken
}

func IndexPage(c echo.Context) error {
	csrf := getCSRFToken(c)
	return templates.IndexPage(csrf).Render(context.Background(), c.Response().Writer)
}

func LoginPage(c echo.Context) error {
	csrfToken := getCSRFToken(c)
	return templates.LoginPage(csrfToken).Render(context.Background(), c.Response().Writer)
}

func SignupPage(c echo.Context) error {
	csrfToken := getCSRFToken(c)
	return templates.Signup(csrfToken).Render(context.Background(), c.Response().Writer)
}

func Dashboard(c echo.Context) error {
	userEmail, _ := c.Get("user_email").(string)
	return templates.Dashboard(userEmail).Render(context.Background(), c.Response().Writer)
}
