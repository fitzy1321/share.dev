package internal

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/supabase-community/supabase-go"
)

const (
	Dev         = "dev"
	Development = "development"
	Test        = "test"
	Prod        = "prod"
	Production  = "production"
)

type CustomContext struct {
	echo.Context
	Env string
}

func (c *CustomContext) IsDev() bool {
	return c.Env == Dev || c.Env == Development
}

func (c *CustomContext) IsTest() bool {
	return c.Env == Test
}

func (c *CustomContext) IsProd() bool {
	return c.Env == Prod || c.Env == Production
}

func LoadEnvFile() (string, error) {
	env, ok := os.LookupEnv("GO_ENV")
	if !ok {
		return "", godotenv.Load(".env.local")
	}

	var envFile string
	switch env {
	case Dev, Development:
		envFile = ".env.local"
	case Test:
		envFile = ".env.test"
	case Prod, Production:
		// production is fly.io (for now)
		// no .env file, just leave
		return env, nil
	default:
		return "", errors.New("unknown $GO_ENV key found")
	}

	return env, godotenv.Load(envFile)
}

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
