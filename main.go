package main

import (
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"share.dev/handlers"
	"share.dev/handlers/api"
	"share.dev/internal"
	"share.dev/routes"
)

func main() {
	// Load env
	env, err := internal.LoadEnvFile()
	if err != nil || env == "" {
		log.Fatalf("error loading .env: %v", err)
	}
	client, err := internal.NewSupabase()
	if err != nil {
		log.Fatalln("Error preparing supabase client:", err)
	}

	e := echo.New()

	e.HTTPErrorHandler = handlers.CustomErrorHandler

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &internal.CustomContext{Context: c, Env: env}
			return next(cc)
		}
	})
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Security and CSRF middleware
	// e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
	e.Use(middleware.CSRF())
	e.Use(middleware.Secure())

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
	e.GET(routes.Verify, handlers.Verify(client))

	// API / Data routes
	g := e.Group("/api", handlers.AuthRequired)
	g.GET("/feed", api.Feed)

	// Start Server, and log results
	e.Logger.Fatal(e.Start(":8080"))
}
