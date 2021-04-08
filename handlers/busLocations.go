package handlers

import (
	"errors"
	"server/utils"
	"server/controllers"
	"server/types"
	"strconv"
	"strings"
	"net/http"
	"fmt"
)

type busLocationHandler struct {}

// BusLocation handles all bus location requests (GET, OPTIONS)
func BusLocation(w http.ResponseWriter, r *http.Request) {
	busLocationHandler := busLocationHandler{}
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
		case http.MethodGet: busLocationHandler.get(w, r)
		default:
			w.Header().Set("Allow", strings.Join(acceptedMethods, ", "))
			w.WriteHeader(http.StatusMethodNotAllowed)

			fmt.Fprint(w, http.StatusText(http.StatusMethodNotAllowed))
	}
}

type getBusLocationBody struct {
	Buses  []types.Bus
}

// get is a GET route for getting bus locations within a bounds
func (*busLocationHandler) get(w http.ResponseWriter, r *http.Request) {
	urlQuery := r.URL.Query()

	topLeftCoordinate, err := parseCoordinate(urlQuery.Get("topLeft"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)

		return
	}

	bottomRightCoordinate, err := parseCoordinate(urlQuery.Get("bottomRight"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)

		return
	}

	if topLeftCoordinate.Longitude >= bottomRightCoordinate.Longitude {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Top left coordinate must be above bottom right coordinate")

		return
	}

	if topLeftCoordinate.Latitude <= bottomRightCoordinate.Latitude {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Top left coordinate must be to the left of bottom right coordinate")

		return
	}

	busLocations, err := controllers.GetBusLocations(topLeftCoordinate, bottomRightCoordinate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, http.StatusText(http.StatusInternalServerError))

		return
	}

	// Response ok
	response := getBusLocationBody{Buses: busLocations}
	compress := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")

	utils.SendJSONResponse(w, http.StatusOK, compress, response)
}

func parseCoordinate(coordinate string) (types.Coordinate, error) {
	splitCoordinates := strings.Split(coordinate, ",")

	if len(splitCoordinates) != 2 {
		return types.Coordinate{}, errors.New("Coordinate must only contain longitude and latitude")
	}

	longitude, err := strconv.ParseFloat(splitCoordinates[0], 32)
	if err != nil {
		return types.Coordinate{}, errors.New("Coordinate must be of type Coordinate(float32, float32)")
	}

	latitude, err := strconv.ParseFloat(splitCoordinates[1], 32)
	if err != nil {
		return types.Coordinate{}, errors.New("Coordinate must be of type Coordinate(float32, float32)")
	}

	return types.Coordinate{
			Longitude: float32(longitude),
			Latitude:  float32(latitude),
		},
		nil
}