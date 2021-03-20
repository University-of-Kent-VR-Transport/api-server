package main

import (
	"fmt"
	"net/http"
	"server/handlers"
	"os"
)

func main() {
	http.HandleFunc("/health-check", handlers.HealthCheckHandler)
	http.HandleFunc("/", http.NotFound)

	fmt.Println("Listening on port 5050...")

	if err := http.ListenAndServe(":5050", nil); err != nil {
		fmt.Fprintln(os.Stderr, "Service crashed")
		fmt.Fprintln(os.Stderr, err)
	}
}
