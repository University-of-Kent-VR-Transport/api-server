package handlers

import (
	"server/models"
	"server/utils"
	"server/controllers"
	"strconv"
	"strings"
	"net/http"
	"fmt"
	"os"
)

type busRouteHandler struct {}

// BusRoutes handles all bus routes requests (GET, OPTIONS)
func BusRoutes(w http.ResponseWriter, r *http.Request) {
	busRouteHandler := busRouteHandler{}
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
		case http.MethodGet: busRouteHandler.get(w, r)
		case http.MethodPut: busRouteHandler.update(w, r)
		default:
			w.Header().Set("Allow", strings.Join(acceptedMethods, ", "))
			w.WriteHeader(http.StatusMethodNotAllowed)

			fmt.Fprint(w, http.StatusText(http.StatusMethodNotAllowed))
	}
}

type getBusRouteBody struct {
	Route  models.Route
}

// get is a GET route for getting a bus route by 
func (*busRouteHandler) get(w http.ResponseWriter, r *http.Request) {
	urlQuery := r.URL.Query()

	lineName := urlQuery.Get("lineName")
	if lineName == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "You must provide a lineName")
		
		return
	}

	direction := strings.ToUpper(urlQuery.Get("direction"))
	if direction != "INBOUND" && direction != "OUTBOUND" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Println(direction)
		fmt.Fprint(w, "Direction must be \"INBOUND\" or \"OUTBOUND\"")
		
		return
	}

	operatorID := urlQuery.Get("operatorID")
	if operatorID == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "You must provide an operatorID")
		
		return
	}

	route, err := controllers.GetRoute(lineName, direction, operatorID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		fmt.Fprint(w, http.StatusText(http.StatusInternalServerError))

		return
	}

	if route.LineID == "" {
		w.WriteHeader(http.StatusNotFound)

		fmt.Fprint(w, http.StatusText(http.StatusNotFound))

		return
	}

	// Response ok
	response := getBusRouteBody{ Route: route }
	compress := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")

	// Set caching header to 6 hours
	w.Header().Set("Cache-Control", "public, maxage=21600")

	// Send json response
	utils.SendJSONResponse(w, http.StatusOK, compress, response)
}

// update is a UPDATE route for updating routes with a datasetID. The route is
// protected by an admin token
func (*busRouteHandler) update(w http.ResponseWriter, r *http.Request) {
	authorizationHeader := r.Header.Get("Authorization")
	adminToken := os.Getenv("ADMIN_TOKEN")
	if authorizationHeader != "Bearer " + adminToken {
		fmt.Fprintln(os.Stderr, "Unauthorized request to GET /api/bus-routes")

		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, http.StatusText(http.StatusUnauthorized))

		return
	}

	urlPath := strings.Split(r.URL.EscapedPath(), "/")

	datasetID, err := strconv.ParseUint(urlPath[3], 10, 32)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, "Dataset ID must be a positive integer")

		return
	}

	job, err := controllers.UpdateRoute(uint(datasetID), &http.Client{}, &models.BusRoutes{})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)

		fmt.Fprint(w, http.StatusText(http.StatusInternalServerError))

		return
	}

	// Request accepted
	response := putBusStopBody{ Job: job }
	compress := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")

	utils.SendJSONResponse(w, http.StatusAccepted, compress, response)
}