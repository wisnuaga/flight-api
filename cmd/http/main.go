package main

import (
	"fmt"
	"log"

	"github.com/wisnuaga/flight-api/internal/config"
	"github.com/wisnuaga/flight-api/internal/delivery/http/router"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Setup HTTP router (handles full composition root wiring internally)
	r := router.Setup(cfg)

	port := fmt.Sprintf(":%s", cfg.Service.Port)
	log.Printf("Server is starting and listening on port %s...\n", port)

	if err := r.Run(port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
