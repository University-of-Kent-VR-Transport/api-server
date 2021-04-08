package main

import (
	"server/utils"
	"server/handlers"
	"fmt"
	"net/http"
	"os"
)

const serverDocs = "https://github.com/University-of-Kent-VR-Transport/api-server/tree/master/docs"
const clientReleases = "https://github.com/University-of-Kent-VR-Transport/vr-client/releases"

func main() {
	if !utils.VerifyEnvSet() {
		os.Exit(1)
	}

	router := http.NewServeMux()

	// file server
	fileServer := http.FileServer(http.Dir("./public"))
	router.Handle("/static/", http.StripPrefix("/static", fileServer))

	// redirect
	router.Handle("/docs", http.RedirectHandler(serverDocs, 301))
	router.Handle("/download", http.RedirectHandler(clientReleases, 301))

	// api routes
	router.HandleFunc("/api/bus-locations", handlers.BusLocation)
	router.HandleFunc("/api/bus-stops", handlers.BusStop)
	router.HandleFunc("/api/job", handlers.BackgroundJob)
	router.HandleFunc("/api/job/", handlers.BackgroundJob)
	router.HandleFunc("/api/health-check", handlers.HealthCheck)

	// html routes
	router.HandleFunc("/", handlers.Index)

	fmt.Println("Listening on port 5050...")

	if err := http.ListenAndServe(":5050", router); err != nil {
		fmt.Fprintln(os.Stderr, "Service crashed", err)
		os.Exit(1)
	}
}
