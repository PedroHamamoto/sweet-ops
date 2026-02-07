package main

import (
	"log"
	"net/http"

	"sweet-ops/internal/http/routes"
)

func main() {
	r := routes.NewRouter()
	log.Println("Starting the application on :8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
