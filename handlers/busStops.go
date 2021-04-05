package handlers

import (
	"fmt"
	"server/models"
	"compress/gzip"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type busStopResponse struct {
	BusStops  []models.BusStop
}

func BusStops(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprint(w, http.StatusText(http.StatusMethodNotAllowed))

		return
	}

	if !(strings.Contains(r.Header.Get("Accept"), "*/*") ||
		strings.Contains(r.Header.Get("Accept"), "application/json")) {
		w.WriteHeader(http.StatusNotAcceptable)
		fmt.Fprint(w, "application/json")

		return
	}

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

	response := busStopResponse{BusStops: busStops}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Error while encoding bus stop results: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	w.Header().Set("Content-Encoding", "gzip")
	w.WriteHeader(http.StatusOK)

	gz := gzip.NewWriter(w)
	if err := json.NewEncoder(gz).Encode(response); err != nil {
		log.Printf("Error while compressing bus stop results: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
	gz.Close()
}