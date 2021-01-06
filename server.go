package main

import (
	"fmt"
	"io"
	"net/http"
	"server/handlers"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello World")
}

func main() {
	http.HandleFunc("/health-check", handlers.HealthCheckHandler)
	http.HandleFunc("/", mainHandler)

	fmt.Println("Listening on port 5050...")

	http.ListenAndServe(":5050", nil)
}
