package main

import (
	"fmt"
	"log"
	"time"

	"github.com/wisnuaga/flight-api/internal/config"
	"github.com/wisnuaga/flight-api/internal/delivery/http/router"
	"github.com/wisnuaga/flight-api/internal/infra/cache"
	infraprovider "github.com/wisnuaga/flight-api/internal/infra/provider"
	"github.com/wisnuaga/flight-api/internal/usecase"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Build infrastructure
	flightCache := cache.NewMemoryCache(5 * time.Minute)
	registry := infraprovider.NewRegistry(cfg)

	// Build usecase with injected dependencies
	flightUsecase := usecase.NewFlightUsecase(registry.GetProviders(), flightCache)

	// Setup HTTP router (pure routing, no DI logic)
	r := router.Setup(flightUsecase)

	port := fmt.Sprintf(":%s", cfg.Service.Port)
	log.Printf("Server is starting and listening on port %s...\n", port)

	if err := r.Run(port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
