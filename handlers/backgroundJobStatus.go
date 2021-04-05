package handlers

import (
	"encoding/json"
	"strconv"
	"fmt"
	"server/models"
	"strings"
	"net/http"
	"log"
	"os"
	"database/sql"
	_ "github.com/lib/pq"
)

type getBackgroundJobStatusResponse struct {
	Job  models.BackgroundJob
}

func GetBackgroundJobStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, http.StatusText(http.StatusMethodNotAllowed))

		return
	}

	urlPath := strings.Split(r.URL.EscapedPath(), "/")
	jobID, err := strconv.ParseUint(urlPath[3], 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Job ID must be a positive integer")

		return
	}

	connectionString, _ := os.LookupEnv("DATABASE_URL")

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Println("Couldn't connect to db", err)

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, http.StatusText(http.StatusInternalServerError))

		return
	}
	defer db.Close()

	job, err := models.GetBackgroundJob(uint(jobID), db)
	if err != nil {
		if err == sql.ErrNoRows {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "No job found with ID of %v", jobID)
	
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, http.StatusText(http.StatusInternalServerError))

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := getBackgroundJobStatusResponse{Job: job}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error while encoding background job response: %s", err.Error())
	}
}