package main

import (
	"fmt"
	"io"
	"net/http"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello world")
}

func main() {
	http.HandleFunc("/", mainHandler)

	fmt.Println("Listening on port 5050...")

	http.ListenAndServe(":5050", nil)
}
