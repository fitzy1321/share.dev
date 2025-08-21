package internal

import (
	"log"
	"os"

	"github.com/supabase-community/supabase-go"
)

// import "github.com/labstack/echo/v4"

// type AppContext struct {
// 	echo.Context
// 	Supabase *supabase.Client
// }

// NewClient creates and returns a new Supabase client using environment variables SUPABASE_URL and SUPABASE_KEY.
func NewSupabase() (*supabase.Client, error) {
	supaURL, ok := os.LookupEnv("SUPABASE_URL")
	if !ok {
		log.Fatalln("'SUPABASE_URL' env key not found")
	}
	supaKey, ok := os.LookupEnv("SUPABASE_ANON_KEY")
	if !ok {
		log.Fatalln("'SUPABASE_ANON_KEY' env key not found")
	}
	return supabase.NewClient(supaURL, supaKey, &supabase.ClientOptions{})
}
