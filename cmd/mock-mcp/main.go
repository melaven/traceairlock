package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status": "ok", "message": "Hit the mock server", "path": "%s"}`, r.URL.Path)
	})

	log.Println("Mock server starting on :9090")

	if err := http.ListenAndServe(":9090", nil); err != nil {
		log.Fatalf("Server crashed: %v", err)
	}
}
