package bootstrap

import (
	"github.com/wisnuaga/flight-api/internal/command"
	"github.com/wisnuaga/flight-api/internal/config"
	"github.com/wisnuaga/flight-api/internal/domain/entity"
	infraprovider "github.com/wisnuaga/flight-api/internal/infra/provider"
	"github.com/wisnuaga/flight-api/internal/port"
	"github.com/wisnuaga/flight-api/internal/usecase"
	"github.com/wisnuaga/flight-api/pkg/cache"
)

// App serves as the Composition Root container, holding all fully instantiated
// application boundaries across the program.
type App struct {
	FlightUsecase port.FlightUsecase
}

// NewApp bootstraps and manually injects all dependencies recursively.
func NewApp(cfg *config.Config) *App {
	// 1. Build Infrastructure Resources
	flightCache := cache.NewMemoryCache[[]*entity.Flight]()
	registry := infraprovider.NewRegistry(cfg)

	// 2. Build Application Commands
	filterCmd := command.NewFlightFilterCommand()
	sortCmd := command.NewFlightSortCommand()

	// 3. Build Core Usecases with deeply injected contexts
	flightUsecase := usecase.NewFlightUsecase(
		registry.GetProviders(),
		flightCache,
		filterCmd,
		sortCmd,
	)

	return &App{
		FlightUsecase: flightUsecase,
	}
}
