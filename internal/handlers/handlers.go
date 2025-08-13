package handlers

import (
	"fmt"
	"html/template"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	// http.ServeFile(w, r, "index.html")
	tmpl, err := template.ParseFiles("views/index.html")
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Message string
	}{
		Message: "Hello, World! Look at me mama, NO FRAMEWORKS!",
	}

	tmpl.Execute(w, data)
}
