package handlers

import (
	"net/http"

	"github.com/fitzy1321/share.dev/components"

	"github.com/supabase-community/supabase-go"
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
