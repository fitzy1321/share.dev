package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/fitzy1321/share.dev/internal/handlers"

	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"
)

func initSupaClient(envs map[string]string) *supabase.Client {
	API_KEY, ok := envs["SUPABASE_ANON_KEY"]
	if !ok {
		log.Fatalln("This API needs env key 'SUPABASE_ANON_KEY' to run!")
	}
	client, err := supabase.NewClient("https://cmdvpbcuqjxljfewsxng.supabase.co", API_KEY, &supabase.ClientOptions{})
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
	envs, err := godotenv.Read()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var _ = initSupaClient(envs)

	http.HandleFunc("/", handlers.Index)

	http.HandleFunc("/time", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Current time: " + time.Now().Format(time.RFC1123)))
	})

	fmt.Println("Started http server on port :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalln(err)
	}
}
