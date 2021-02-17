package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"server/handlers"
)

func main() {
	if _, isPresent := os.LookupEnv("DFT_SECRET"); isPresent == false {
		log.Panic("No DFT_SECRET provided. Exiting...")
	}

	http.HandleFunc("/api/get-bus-locations", handlers.BoundingBoxHandler)
	http.NotFoundHandler()

	fmt.Println("Listening on port 5050...")

	http.ListenAndServe(":5050", nil)
}
