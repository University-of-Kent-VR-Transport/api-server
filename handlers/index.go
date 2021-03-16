package handlers

import (
	"fmt"
	"net/http"
	"text/template"
	"os"
	"strings"
	"io"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	acceptHeader := r.Header.Get("Accept")
	if !strings.Contains(acceptHeader, "*/*") && !strings.Contains(acceptHeader, "text/html") {
		w.WriteHeader(http.StatusNotAcceptable)
		io.WriteString(w, "text/html")

		return
	}

	parsedTemplate, err := template.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	type mapBoxConfig struct {
		AccessToken string
		Style       string
	}

	var config = mapBoxConfig{
		AccessToken: os.Getenv("MAPBOX_TOKEN"),
		Style:       "mapbox://styles/mapbox/light-v10",
	}

	if err := parsedTemplate.Execute(w, config); err != nil {
		fmt.Fprintln(os.Stderr, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}