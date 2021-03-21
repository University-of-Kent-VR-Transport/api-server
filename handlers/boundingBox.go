package handlers

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"server/models"
	"server/transformers"
	"server/types"
	"strconv"
	"strings"
)

type boundingBoxResponse struct {
	Buses  []types.Bus
}


// BoundingBoxHandler retieives all the buses in a given bounding box
func BoundingBoxHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		w.WriteHeader(http.StatusMethodNotAllowed)

		return
	}

	if !(strings.Contains(r.Header.Get("Accept"), "*/*") ||
		strings.Contains(r.Header.Get("Accept"), "application/json")) {
		w.WriteHeader(http.StatusNotAcceptable)
		io.WriteString(w, "application/json")

		return
	}

	urlQuery := r.URL.Query()
	topLeft := urlQuery.Get("topLeft")
	bottomRight := urlQuery.Get("bottomRight")

	topLeftCoordinate, err := parseCoordinate(topLeft)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, err.Error())

		return
	}

	bottomRightCoordinate, err := parseCoordinate(bottomRight)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, err.Error())

		return
	}

	if topLeftCoordinate.Longitude >= bottomRightCoordinate.Longitude {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Top left coordinate must be above bottom right coordinate")

		return
	}

	if topLeftCoordinate.Latitude <= bottomRightCoordinate.Latitude {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, "Top left coordinate must be to the left of bottom right coordinate")

		return
	}

	resp, err := models.GetBusLocation(topLeftCoordinate, bottomRightCoordinate)
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)

		return
	}

	buses := transformers.Bus(resp)
	var response boundingBoxResponse = boundingBoxResponse{Buses: buses}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Printf("Error while encoding bus result: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
		}

		return
	}

	w.Header().Set("Content-Encoding", "gzip")
	w.WriteHeader(http.StatusOK)

	gz := gzip.NewWriter(w)
	if err := json.NewEncoder(gz).Encode(response); err != nil {
		log.Printf("Error while compressing bus result: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
	gz.Close()
}

func parseCoordinate(coordinate string) (types.Coordinate, error) {
	splitCoordinates := strings.Split(coordinate, ",")

	if len(splitCoordinates) != 2 {
		return types.Coordinate{}, fmt.Errorf("Expect coordinate to contain longitude and latitude but received %d values", len(splitCoordinates))
	}

	longitude, err := strconv.ParseFloat(splitCoordinates[0], 32)
	if err != nil {
		return types.Coordinate{}, err
	}

	latitude, err := strconv.ParseFloat(splitCoordinates[1], 32)
	if err != nil {
		return types.Coordinate{}, err
	}

	return types.Coordinate{
			Longitude: float32(longitude),
			Latitude:  float32(latitude),
		},
		nil
}
