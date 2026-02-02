// main.go
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/genjishimada/playtest-plotter/handler"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/chart", handler.ChartHandler)
	http.HandleFunc("/health", handler.HealthHandler)

	log.Printf("Starting server on :%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
