package handlers

import (
	"server/utils"
	"compress/gzip"
	"fmt"
	"net/http"
	"text/template"
	"os"
	"strings"
)

type indexHandler struct {}

const contentTypeJson = "application/json; charset=utf-8"
const contentTypeHtml = "text/html; charset=utf-8"

// Index is the root of the server and sends to the client the html index page
func Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)

		return
	}

	indexHandler := indexHandler{}
	acceptedMethods := []string{http.MethodGet,http.MethodOptions}

	if r.Method == http.MethodOptions {
		utils.OptionsResponse(w, acceptedMethods, contentTypeHtml)

		return
	}

	// Check content type of HTML is accepted by client
	acceptHeader := r.Header.Get("Accept")
	if !(strings.Contains(acceptHeader, "*/*") ||
		strings.Contains(acceptHeader, "text/html")) {
		w.WriteHeader(http.StatusNotAcceptable)

		fmt.Fprint(w, contentTypeHtml)

		return
	}

	switch method := r.Method; method {
		case http.MethodGet: indexHandler.get(w, r)
		default:
			w.Header().Set("Allow", strings.Join(acceptedMethods, ", "))
			w.WriteHeader(http.StatusMethodNotAllowed)

			fmt.Fprint(w, http.StatusText(http.StatusMethodNotAllowed))
	}
}

type mapBoxConfig struct {
	AccessToken string
	Style       string
}

// get sends the client the html index page
func (*indexHandler) get(w http.ResponseWriter, r *http.Request) {
	parsedTemplate, err := template.ParseFiles("views/index.html")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	var config = mapBoxConfig{
		AccessToken: os.Getenv("MAPBOX_TOKEN"),
		Style:       "mapbox://styles/mapbox/light-v10",
	}

	w.Header().Set("Content-Type", contentTypeHtml)

	if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		w.WriteHeader(http.StatusOK)

		if err := parsedTemplate.Execute(w, config); err != nil {
			fmt.Fprint(os.Stderr, "Error encoding (W) get root response", err)
		}

		return
	}

	w.Header().Set("Content-Encoding", "gzip")
	w.WriteHeader(http.StatusOK)

	gz := gzip.NewWriter(w)
	defer gz.Close()

	if err := parsedTemplate.Execute(gz, config); err != nil {
		fmt.Fprint(os.Stderr, "Error encoding (GZ) get root response", err)
	}
}