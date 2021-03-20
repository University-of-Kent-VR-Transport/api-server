package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello world")
}

func main() {
	http.HandleFunc("/", mainHandler)

	fmt.Println("Listening on port 5050...")

	if err := http.ListenAndServe(":5050", nil); err != nil {
		fmt.Fprintln(os.Stderr, "Service crashed")
		fmt.Fprintln(os.Stderr, err)
	}
}
