package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"share.dev/handlers"
	"share.dev/handlers/api"
	"share.dev/internal"
)

func loadEnvFile() error {
	env := os.Getenv("GO_ENV")
	var envFile string
	switch env {
	case "development":
		envFile = ".env.local"
	case "test":
		envFile = ".env.test"
	case "production":
		// production is fly.io (for now)
		// no .env file, just leave
		return nil
	default:
		envFile = ".env"
		return nil
	}
	return godotenv.Load(envFile)
}

func main() {
	// Load env
	if err := loadEnvFile(); err != nil {
		log.Fatalf("error loading .env: %v", err)
	}
	client, err := internal.NewSupabase()
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

	e.Static("/static", "static")

	// Routes
	e.GET("/", handlers.IndexPage)

	// e.GET("/login", handlers.LoginPage)
	e.POST("/login", handlers.Login(client))

	// e.GET("/signup", handlers.SignupPage)
	// e.POST("/signup", handlers.Signup(client))

	e.GET("/logout", handlers.Logout(client))

	e.GET("/home", handlers.Home, handlers.AuthRequired)

	g := e.Group("/api", handlers.AuthRequired)
	g.GET("/feed", api.Feed)
	e.Logger.Fatal(e.Start(":8080"))
}
