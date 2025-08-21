package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"

	"share.dev/handlers"
	supabaseclient "share.dev/supabase"
)

func loadEnvFile() error {
	env := os.Getenv("GO_ENV")
	var envFile string
	switch env {
	case "production":
		envFile = ".env"
	case "test":
		envFile = ".env.test"
	default:
		envFile = ".env.local"
	}
	return godotenv.Load(envFile)
}

func main() {
	// Load env
	if err := loadEnvFile(); err != nil {
		log.Fatalf("error loading .env: %v", err)
	}
	client, err := supabaseclient.GetClient()
	if err != nil {
		log.Fatalln("Error preparing supabase client:", err)
	}

	e := echo.New()

	e.Static("/static", "static")

	e.HTTPErrorHandler = handlers.CustomErrorHandler

	// Security and CSRF middleware
	e.Use(handlers.RateLimitMiddleware)
	e.Use(handlers.CSRFMiddleware)
	e.Use(handlers.SecurityHeaders)

	// Setup Supabase

	// Routes using a-h/templ
	e.GET("/", handlers.IndexPage)
	e.GET("/login", handlers.LoginPage)
	e.POST("/login", handlers.Login(client))
	// e.GET("/signup", handlers.SignupPage)
	// e.POST("/signup", handlers.Signup(client))
	e.GET("/dashboard", handlers.Dashboard, handlers.AuthRequired)
	e.GET("/logout", handlers.Logout(client))

	e.Logger.Fatal(e.Start(":8080"))
}
