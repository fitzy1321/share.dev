package api

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"share.dev/models"
)

// GET /api/feed/?page=0;limit=20
func Feed(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	// Defaults
	if page == 0 {
		page = 1
	}
	if limit == 0 {
		limit = 100
	}
	posts := []*models.Post{}
	return c.JSON(http.StatusOK, posts)
}
