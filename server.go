package main

import (
	"fmt"
	"net/http"
	"os"
	"log"
	"server/handlers"
)

func main() {
	if _, isPresent := os.LookupEnv("MAPBOX_TOKEN"); isPresent == false {
		log.Panic("No MAPBOX_TOKEN provided. Exiting...")
}

	router := http.NewServeMux()

	router.HandleFunc("/", handlers.IndexHandler)

	fmt.Println("Listening on port 5050...")

	if err := http.ListenAndServe(":5050", router); err != nil {
		fmt.Println(err)
	}
}
