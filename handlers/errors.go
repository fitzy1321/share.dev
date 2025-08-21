package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"share.dev/templates"
)

func CustomErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	var message string

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code
		if m, ok := he.Message.(string); ok {
			message = m
		} else {
			message = he.Error()
		}
	} else {
		message = err.Error()
	}

	// Log internal server errors
	if code >= 500 {
		c.Logger().Error(err)
	}

	// Fallback plain text for non-HTML requests (optional)
	if c.Request().Header.Get("Accept") != "text/html" {
		c.String(code, message)
		return
	}

	// Render templ error page
	if renderErr := templates.ErrorPage(code, message).Render(c.Request().Context(), c.Response().Writer); renderErr != nil {
		c.Logger().Error(renderErr)
		c.String(http.StatusInternalServerError, "Internal Server Error")
	}
}
