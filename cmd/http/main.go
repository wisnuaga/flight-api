package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/wisnuaga/flight-api/internal/config"
	"github.com/wisnuaga/flight-api/internal/delivery/http/router"
)

func main() {
	// Load Configurations
	cfg := config.LoadConfig()

	// Init HTTP mux
	mux := http.NewServeMux()

	// Setup routes & handlers
	router.Setup(mux)

	port := fmt.Sprintf(":%s", cfg.Service.Port)
	log.Printf("Server is starting and listening on port %s...\n", port)

	// Start server
	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}
