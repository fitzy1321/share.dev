package internal

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/supabase-community/supabase-go"
)

type AppContext struct {
	echo.Context
	Supabase *supabase.Client
}

// NewClient creates and returns a new Supabase client using environment variables SUPABASE_URL and SUPABASE_KEY.
func NewSupabase() (*supabase.Client, error) {
	supaURL := os.Getenv("SUPABASE_URL")
	supaKey := os.Getenv("SUPABASE_ANON_KEY")
	return supabase.NewClient(supaURL, supaKey, &supabase.ClientOptions{})
}
