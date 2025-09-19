package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"share.dev/routes"
	"share.dev/templates/components"
	"share.dev/templates/pages"
)

func getCSRFToken(c echo.Context) string {
	csrfToken, ok := c.Get(middleware.DefaultCSRFConfig.ContextKey).(string)
	if !ok {
		c.Logger().Fatal("CSRF Token not found")
	}
	return csrfToken
}

func IndexPage(c echo.Context) error {
	_, ok := c.Cookie(accessTokenCookie)
	if ok == nil { // if there IS an accessTokenCookie, redirect to dashboard
		c.Redirect(http.StatusPermanentRedirect, routes.MainPage)
	}
	csrf := getCSRFToken(c)
	return pages.IndexPage(csrf).Render(context.Background(), c.Response().Writer)
}

func MainPage(c echo.Context) error {
	userEmail, _ := c.Get("user_email").(string)
	return pages.MainPage(userEmail).Render(context.Background(), c.Response().Writer)
}

func CheckEmailPage(c echo.Context) error {
	return pages.CheckEmailPage().Render(context.Background(), c.Response().Writer)
}

func AuthFormGet(c echo.Context) error {
	if h := c.Request().Header.Get("HX-Request"); h == "" {
		c.Error(errors.New("Not an HTMX request!"))
	}
	csrf := getCSRFToken(c)
	return components.AuthForm(csrf).Render(context.Background(), c.Response().Writer)
}
