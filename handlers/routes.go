package handlers

import (
	"context"

	"github.com/labstack/echo/v4"
	"share.dev/templates"
)

func IndexPage(c echo.Context) error {
	// csrfToken, _ := c.Get("csrf_token").(string)
	// return templates.IndexPage(csrfToken).Render(context.Background(), c.Response().Writer)
	return templates.IndexPage().Render(context.Background(), c.Response().Writer)

}

func LoginPage(c echo.Context) error {
	csrfToken, _ := c.Get("csrf_token").(string)
	return templates.LoginPage(csrfToken).Render(context.Background(), c.Response().Writer)
}

func SignupPage(c echo.Context) error {
	csrfToken, _ := c.Get("csrf_token").(string)
	return templates.Signup(csrfToken).Render(context.Background(), c.Response().Writer)
}

func Dashboard(c echo.Context) error {
	userEmail := c.Get("user_email").(string)
	return templates.Dashboard(userEmail).Render(context.Background(), c.Response().Writer)
}
