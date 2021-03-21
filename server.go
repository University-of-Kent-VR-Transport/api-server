package main

import (
	"fmt"
	"net/http"
	"server/handlers"
	"os"
)

func main() {
	if _, isPresent := os.LookupEnv("DFT_SECRET"); isPresent == false {
		fmt.Fprintln(os.Stderr, "No DFT_SECRET provided")
		os.Exit(1)
	}
	if _, isPresent := os.LookupEnv("MAPBOX_TOKEN"); isPresent == false {
		fmt.Fprintln(os.Stderr, "No MAPBOX_TOKEN provided")
		os.Exit(1)
	}

	router := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./static"))
	router.Handle("/static/", http.StripPrefix("/static", fileServer))

	router.Handle("/docs", http.RedirectHandler("https://github.com/University-of-Kent-VR-Transport/api-server/tree/master/docs", 301))
	router.Handle("/download", http.RedirectHandler("https://github.com/University-of-Kent-VR-Transport/vr-client/releases", 301))

	router.HandleFunc("/api/get-bus-locations", handlers.BoundingBoxHandler)
	router.HandleFunc("/health-check", handlers.HealthCheckHandler)
	router.HandleFunc("/", handlers.IndexHandler)

	fmt.Println("Listening on port 5050...")

	if err := http.ListenAndServe(":5050", router); err != nil {
		fmt.Fprintln(os.Stderr, "Service crashed")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
