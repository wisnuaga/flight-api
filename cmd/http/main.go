package main

import (
	"fmt"
	"log"

	"github.com/wisnuaga/flight-api/internal/command"
	"github.com/wisnuaga/flight-api/internal/config"
	"github.com/wisnuaga/flight-api/internal/delivery/http/router"
	"github.com/wisnuaga/flight-api/internal/domain/entity"
	infraprovider "github.com/wisnuaga/flight-api/internal/infra/provider"
	"github.com/wisnuaga/flight-api/internal/usecase"
	"github.com/wisnuaga/flight-api/pkg/cache"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Build infrastructure
	flightCache := cache.NewMemoryCache[[]*entity.Flight]()
	registry := infraprovider.NewRegistry(cfg)

	// Build application commands
	filterCmd := command.NewFlightFilterCommand()
	sortCmd := command.NewFlightSortCommand()

	// Build usecase with injected dependencies
	flightUsecase := usecase.NewFlightUsecase(registry.GetProviders(), flightCache, filterCmd, sortCmd)

	// Setup HTTP router (pure routing, no DI logic)
	r := router.Setup(flightUsecase)

	port := fmt.Sprintf(":%s", cfg.Service.Port)
	log.Printf("Server is starting and listening on port %s...\n", port)

	if err := r.Run(port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
