package handlers

import (
	"io"
	"fmt"
	"os"
	"net/http"
	"database/sql"
	_ "github.com/lib/pq"
)

// HealthCheckHandler checks the health of the service and other services used
// by this service. It response to the HTTP request.
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	w.Header().Set("Content-Type", "application/json")

	// check database health
	connectionString, _ := os.LookupEnv("DATABASE_URL")

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)

		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"database": false}`)

		return
	}
	defer db.Close()

	if !testDatabaseConnection(db) {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"database": false}`)

		return
	}

	w.WriteHeader(http.StatusOK)

	io.WriteString(w, `{"database": true}`)
}

type sqlDB interface {
	Ping() error
}

func testDatabaseConnection(db sqlDB) bool {
	if err := db.Ping(); err != nil {
		fmt.Fprintln(os.Stderr, err)

		return false
	}

	return true
}
