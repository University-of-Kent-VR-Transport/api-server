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

	http.HandleFunc("/api/get-bus-locations", handlers.BoundingBoxHandler)
	http.HandleFunc("/health-check", handlers.HealthCheckHandler)
	http.NotFoundHandler()

	fmt.Println("Listening on port 5050...")

	if err := http.ListenAndServe(":5050", nil); err != nil {
		fmt.Fprintln(os.Stderr, "Service crashed")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
