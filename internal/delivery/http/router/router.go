package router

import (
	"github.com/gin-gonic/gin"
	"github.com/wisnuaga/flight-api/internal/command"
	"github.com/wisnuaga/flight-api/internal/config"
	"github.com/wisnuaga/flight-api/internal/delivery/http/handler"
	"github.com/wisnuaga/flight-api/internal/domain/entity"
	infraprovider "github.com/wisnuaga/flight-api/internal/infra/provider"
	"github.com/wisnuaga/flight-api/internal/usecase"
	"github.com/wisnuaga/flight-api/pkg/cache"
)

// Setup builds the Gin engine with all routes registered.
// For this slim API, it serves as our composition root.
func Setup(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// 1. Build Infrastructure
	flightCache := cache.NewMemoryCache[[]*entity.Flight]()
	registry := infraprovider.NewRegistry(cfg)

	// 2. Build Commands
	filterCmd := command.NewFlightFilterCommand()
	sortCmd := command.NewFlightSortCommand()

	// 3. Build Core Usecases
	flightUsecase := usecase.NewFlightUsecase(
		registry.GetProviders(),
		flightCache,
		filterCmd,
		sortCmd,
	)

	// 4. Build Endpoints
	handlers := Handlers{
		Health: handler.NewHealthHandler(),
		Flight: handler.NewFlightHandler(&handler.FlightHandlerUsecases{
			FlightUsecase: flightUsecase,
		}),
	}

	registerRoutes(r, handlers)
	return r
}

type Handlers struct {
	Health *handler.HealthHandler
	Flight *handler.FlightHandler
}

func registerRoutes(r *gin.Engine, h Handlers) {
	registerHealthRoutes(r, h.Health)
	registerFlightRoutes(r, h.Flight)
}
