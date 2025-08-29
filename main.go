package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/supabase-community/supabase-go"
	"share.dev/handlers"
	"share.dev/handlers/api"
	"share.dev/routes"
)

func loadEnvFile() error {
	env, ok := os.LookupEnv("GO_ENV")
	if !ok {
		return godotenv.Load(".env.local")
	}

	var envFile string
	switch env {
	case "development":
		envFile = ".env.local"
	case "test":
		envFile = ".env.test"
	case "production", "prod":
		// production is fly.io (for now)
		// no .env file, just leave
		return nil
	default:
		envFile = ".env"
	}
	return godotenv.Load(envFile)
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

func main() {
	// Load env
	if err := loadEnvFile(); err != nil {
		log.Fatalf("error loading .env: %v", err)
	}
	client, err := NewSupabase()
	if err != nil {
		log.Fatalln("Error preparing supabase client:", err)
	}

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Security and CSRF middleware
	// e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
	e.Use(middleware.CSRF())
	e.Use(middleware.Secure())

	e.HTTPErrorHandler = handlers.CustomErrorHandler

	e.Static(routes.Static, "static")

	// Routes
	// Pages
	e.GET(routes.IndexPage, handlers.IndexPage)
	e.GET(routes.CheckEmailPage, handlers.CheckEmailPage)
	e.GET(routes.MainPage, handlers.MainPage, handlers.AuthRequired)

	// Auth
	e.POST(routes.Login, handlers.Login(client))
	e.GET(routes.Logout, handlers.Logout(client))
	e.POST(routes.Signup, handlers.Signup(client))
	e.GET(routes.ConfirmEmail, handlers.ConfirmEmail(client))

	// API / Data routes
	g := e.Group("/api", handlers.AuthRequired)
	g.GET("/feed", api.Feed)

	// Start Server, and log results
	e.Logger.Fatal(e.Start(":8080"))
}
