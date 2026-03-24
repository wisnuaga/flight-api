package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/wisnuaga/flight-api/internal/config"
)

func main() {
	// Load Configurations
	cfg := config.LoadConfig()

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	port := fmt.Sprintf(":%s", cfg.Service.Port)
	log.Printf("Server is starting and listening on port %s...\n", port)

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}
