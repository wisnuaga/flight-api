package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/wisnuaga/flight-api/internal/config"
	"github.com/wisnuaga/flight-api/internal/router"
)

func main() {
	// Load Configurations
	cfg := config.LoadConfig()

	mux := http.NewServeMux()

	router.Setup(mux)

	port := fmt.Sprintf(":%s", cfg.Service.Port)
	log.Printf("Server is starting and listening on port %s...\n", port)

	if err := http.ListenAndServe(port, mux); err != nil {
		log.Fatalf("Failed to start server: %v\n", err)
	}
}
