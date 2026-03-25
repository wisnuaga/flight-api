package main

import (
	"fmt"
	"log"

	"github.com/wisnuaga/flight-api/internal/config"
	"github.com/wisnuaga/flight-api/internal/delivery/http/router"
)

func main() {
	// Load Configurations
	cfg := config.LoadConfig()

	// Setup router (Gin)
	r := router.Setup()

	port := fmt.Sprintf(":%s", cfg.Service.Port)
	log.Printf("Server is starting and listening on port %s...\n", port)

	// Start server
	if err := r.Run(port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
