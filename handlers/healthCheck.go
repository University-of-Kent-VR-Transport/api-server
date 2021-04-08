package handlers

import (
	"server/controllers"
	"server/utils"
	"strings"
	"fmt"
	"os"
	"net/http"
	"database/sql"
	_ "github.com/lib/pq"
)

type healthCheckHandler struct {}

// HealthCheckHandler checks the health of the service and other services used
// by this service. It responses to the HTTP request.
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	healthCheckHandler := healthCheckHandler{}
	acceptedMethods := []string{
		http.MethodGet,
		http.MethodOptions,
	}

	if r.Method == http.MethodOptions {
		utils.OptionsResponse(w, acceptedMethods, contentTypeJson)

		return
	}

	// Check content type of JSON is accepted by client
	acceptHeader := r.Header.Get("Accept")
	if !(strings.Contains(acceptHeader, "*/*") ||
		strings.Contains(acceptHeader, "application/json")) {
		w.WriteHeader(http.StatusNotAcceptable)

		fmt.Fprint(w, contentTypeJson)

		return
	}

	switch method := r.Method; method {
		case http.MethodGet: healthCheckHandler.get(w, r)
		default:
			w.Header().Set("Allow", strings.Join(acceptedMethods, ", "))
			w.WriteHeader(http.StatusMethodNotAllowed)

			fmt.Fprint(w, http.StatusText(http.StatusMethodNotAllowed))
	}
}

// get is a GET route for getting the health status of the service
func (*healthCheckHandler) get(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// check database health
	connectionString := os.Getenv("DATABASE_URL")

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{"database": false}`)

		return
	}
	defer db.Close()

	if !controllers.HealthCheck(db) {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, `{"database": false}`)

		return
	}

	w.WriteHeader(http.StatusOK)

	fmt.Fprint(w, `{"database": true}`)
}
