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

	fileServer := http.FileServer(http.Dir("./static"))
	router.Handle("/static/", http.StripPrefix("/static", fileServer))

	router.Handle("/docs", http.RedirectHandler("https://github.com/University-of-Kent-VR-Transport/api-server/tree/master/docs", 301))
	router.Handle("/download", http.RedirectHandler("https://github.com/University-of-Kent-VR-Transport/vr-client/releases", 301))

	router.HandleFunc("/", handlers.IndexHandler)

	fmt.Println("Listening on port 5050...")

	if err := http.ListenAndServe(":5050", router); err != nil {
		fmt.Println(err)
	}
}
