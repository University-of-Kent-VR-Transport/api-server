package main

import (
	"server/util"
	"fmt"
	"net/http"
	"server/handlers"
	"os"
)

func main() {
	if !util.VerifyEnvSet() {
		os.Exit(1)
	}

	router := http.NewServeMux()

	// file server
	fileServer := http.FileServer(http.Dir("./public"))
	router.Handle("/static/", http.StripPrefix("/static", fileServer))

	router.Handle("/docs", http.RedirectHandler("https://github.com/University-of-Kent-VR-Transport/api-server/tree/master/docs", 301))
	router.Handle("/download", http.RedirectHandler("https://github.com/University-of-Kent-VR-Transport/vr-client/releases", 301))

	router.HandleFunc("/api/get-bus-locations", handlers.BoundingBoxHandler)
	router.HandleFunc("/api/bus-stop", handlers.BusStops)
	router.HandleFunc("/api/update-naptan", handlers.UpdateBusStops)
	router.HandleFunc("/api/job/", handlers.GetBackgroundJobStatus)

	router.HandleFunc("/health-check", handlers.HealthCheckHandler)
	router.HandleFunc("/", handlers.IndexHandler)

	fmt.Println("Listening on port 5050...")

	if err := http.ListenAndServe(":5050", router); err != nil {
		fmt.Fprintln(os.Stderr, "Service crashed")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
