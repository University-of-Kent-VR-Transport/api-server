package handlers

import (
	"server/utils"
	"log"
	"strconv"
	"os"
	"strings"
	"server/controllers"
	"server/models"
	"net/http"
	"fmt"
)

type busStopHandler struct {}

// BusStop handles all bus stop requests (GET, PUT, OPTIONS)
func BusStop(w http.ResponseWriter, r *http.Request) {
	busStopHandler := busStopHandler{}
	acceptedMethods := []string{
		http.MethodGet,
		http.MethodPut,
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
		case http.MethodPut: busStopHandler.put(w, r)
		case http.MethodGet: busStopHandler.get(w, r)
		default:
			w.Header().Set("Allow", strings.Join(acceptedMethods, ", "))
			w.WriteHeader(http.StatusMethodNotAllowed)

			fmt.Fprint(w, http.StatusText(http.StatusMethodNotAllowed))
	}
}

type getBusStopBody struct {
	BusStops  []models.BusStop
}

// get is a GET route for getting bus stops within a bounds
func (*busStopHandler) get(w http.ResponseWriter, r *http.Request) {
	urlQuery := r.URL.Query()

	minLongitude, err := strconv.ParseFloat(urlQuery.Get("minLongitude"), 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "minLongitude must be of type float32")

		return
	}

	minLatitude, err := strconv.ParseFloat(urlQuery.Get("minLatitude"), 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "minLatitude must be of type float32")

		return
	}

	maxLongitude, err := strconv.ParseFloat(urlQuery.Get("maxLongitude"), 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "maxLongitude must be of type float32")

		return
	}

	maxLatitude, err := strconv.ParseFloat(urlQuery.Get("maxLatitude"), 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "maxLatitude must be of type float32")

		return
	}

	busStops, err := models.GetBusStopWithinBounds(
		float32(minLongitude), float32(minLatitude),
		float32(maxLongitude), float32(maxLatitude),
	)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, http.StatusText(http.StatusInternalServerError))

		return
	}

	// Response ok
	response := getBusStopBody{BusStops: busStops}
	compress := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")

	utils.SendJSONResponse(w, http.StatusOK, compress, response)
}

type putBusStopBody struct {
	Job  models.BackgroundJob
}

// put is a PUT route for updating the bus stop database from DFT
func (*busStopHandler) put(w http.ResponseWriter, r *http.Request) {
	// Check auth token is correct
	authorizationHeader := r.Header.Get("Authorization")
	adminToken := os.Getenv("ADMIN_TOKEN")
	if authorizationHeader != "Bearer " + adminToken {
		log.Println("Unauthorized request to PUT /api/bus-stops")

		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, http.StatusText(http.StatusUnauthorized))

		return
	}

	// Update bus stops
	jobStatus, err := controllers.UpdateBusStops()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, http.StatusText(http.StatusInternalServerError))

		return
	}

	// Request accepted
	response := putBusStopBody{ Job: jobStatus }
	compress := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")

	utils.SendJSONResponse(w, http.StatusAccepted, compress, response)
}