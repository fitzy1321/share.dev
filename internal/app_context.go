package internal

import (
	"github.com/labstack/echo/v4"
	"github.com/supabase-community/supabase-go"
)

type AppContext struct {
	echo.Context
	Supabase *supabase.Client
}
