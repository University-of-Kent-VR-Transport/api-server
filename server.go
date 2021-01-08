package main

import (
	"fmt"
	"log"
	"net/http"
	"server/handlers"
)

func main() {
	http.HandleFunc("/health-check", handlers.HealthCheckHandler)
	http.HandleFunc("/", http.NotFound)

	fmt.Println("Listening on port 5050...")

	log.Fatal(http.ListenAndServe(":5050", nil))
}
