package handlers

import (
	"net/http"

	"github.com/supabase-community/supabase-go"
	"share.dev/components"
)

func Index(client *supabase.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		component := components.IndexPage("Hello Golang!")
		component.Render(r.Context(), w)
	}
}

func Login(client *supabase.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}
}
