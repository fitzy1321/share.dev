package supabaseclient

import (
	"os"

	"github.com/supabase-community/supabase-go"
)

// NewClient creates and returns a new Supabase client using environment variables SUPABASE_URL and SUPABASE_KEY.
func GetClient() (*supabase.Client, error) {
	supaURL := os.Getenv("SUPABASE_URL")
	supaKey := os.Getenv("SUPABASE_ANON_KEY")
	return supabase.NewClient(supaURL, supaKey, &supabase.ClientOptions{})
}
