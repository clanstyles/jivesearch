// Sample instant demonstrates how to run a simple instant answers server.
package main

import (
	"encoding/json"
	"fmt"
	"jivesearch/instant"
	"log"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	sol := instant.Detect(r)

	if err := json.NewEncoder(w).Encode(sol); err != nil {
		http.Error(w, http.StatusText(500), 500)
	}
}

func favHandler(w http.ResponseWriter, r *http.Request) {}

func main() {
	port := 8000
	http.HandleFunc("/", handler)
	http.HandleFunc("/favicon.ico", favHandler)
	log.Printf("Listening at http://localhost:%d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
