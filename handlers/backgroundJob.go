package handlers

import (
	"server/controllers"
	"server/utils"
	"strconv"
	"fmt"
	"server/models"
	"strings"
	"net/http"
	"log"
	"os"
)

type backgroundJob struct {}

// BackgroundJob takes all bus background job requests (GET, OPTIONS). It is
// protected by an admin auth token
func BackgroundJob(w http.ResponseWriter, r *http.Request) {
	authorizationHeader := r.Header.Get("Authorization")
	adminToken := os.Getenv("ADMIN_TOKEN")
	if authorizationHeader != "Bearer " + adminToken {
		log.Println("Unauthorized request to GET /api/job")

		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, http.StatusText(http.StatusUnauthorized))

		return
	}

	backgroundJob := backgroundJob{}
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
		case http.MethodGet: backgroundJob.get(w, r)
		default:
			w.Header().Set("Allow", strings.Join(acceptedMethods, ", "))
			w.WriteHeader(http.StatusMethodNotAllowed)

			fmt.Fprint(w, http.StatusText(http.StatusMethodNotAllowed))
	}
}

type getBackgroundJobBody struct {
	Job  models.BackgroundJob
}

// get is a GET route for getting a background job by an ID
func (*backgroundJob) get(w http.ResponseWriter, r *http.Request) {
	urlPath := strings.Split(r.URL.EscapedPath(), "/")
	jobID, err := strconv.ParseUint(urlPath[3], 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Job ID must be a positive integer")

		return
	}

	job, found, err := controllers.GetBackgroundJobStatus(uint(jobID))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		fmt.Fprint(w, http.StatusText(http.StatusInternalServerError))

		return
	}

	if !found {
		w.WriteHeader(http.StatusNotFound)

		fmt.Fprint(w, "No job found")

		return
	}

	// Request ok
	response := getBackgroundJobBody{Job: job}
	compress := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")

	utils.SendJSONResponse(w, http.StatusOK, compress, response)
}