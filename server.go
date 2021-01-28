package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello world")
}

func main() {
	if _, isPresent := os.LookupEnv("DFT_SECRET"); isPresent == false {
		log.Panic("No DFT_SECRET provided. Exiting...")
	}

	http.HandleFunc("/", mainHandler)

	fmt.Println("Listening on port 5050...")

	http.ListenAndServe(":5050", nil)
}
