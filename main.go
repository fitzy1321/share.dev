package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/fitzy1321/share.dev/handlers"

	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"
)

// func initSupaClient(envs map[string]string) *supabase.Client {
func initSupaClient() *supabase.Client {
	SUPABASE_URL, ok := os.LookupEnv("SUPABASE_URL")
	if !ok || SUPABASE_URL == "" {
		log.Fatalln("'SUPABASE_URL env key is required")
	}

	SUPABASE_ANON_KEY, ok := os.LookupEnv("SUPABASE_ANON_KEY")
	if !ok || SUPABASE_ANON_KEY == "" {
		log.Fatalln("This API needs env key 'SUPABASE_ANON_KEY' to run!")
	}

	client, err := supabase.NewClient(SUPABASE_URL, SUPABASE_ANON_KEY, &supabase.ClientOptions{})
	if err != nil {
		log.Fatalln("Error occurred setting up Supabase client: ", err)
	}

	_, err = client.Auth.HealthCheck()
	if err != nil {
		log.Fatalln("Supabase Auth Health Check Error: ", err)
	}
	return client
}

func main() {
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

	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatalf("Error loading %s file", envFile)
	}

	client := initSupaClient()

	// server static resources
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", handlers.Index(client))
	http.HandleFunc("/login", handlers.Login(client))

	http.HandleFunc("/time", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Current time: " + time.Now().Format(time.RFC1123)))
	})

	fmt.Println("Started http server on port :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalln(err)
	}
}
