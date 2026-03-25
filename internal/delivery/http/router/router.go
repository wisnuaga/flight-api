package router

import (
	"github.com/gin-gonic/gin"
	"github.com/wisnuaga/flight-api/internal/config"
	"github.com/wisnuaga/flight-api/internal/delivery/http/handler"
	"github.com/wisnuaga/flight-api/internal/repository/provider"
	"github.com/wisnuaga/flight-api/internal/usecase"
)

type Handlers struct {
	Health *handler.HealthHandler
	Flight *handler.FlightHandler
}

type Usecases struct {
	FlightUsecase usecase.FlightUsecase
}

func Setup(cfg *config.Config) *gin.Engine {
	r := gin.Default()

	// Init handlers
	usecases := initUsecases(cfg)
	handlers := initHandlers(&usecases)
	registerRoutes(r, handlers)

	return r
}

func initHandlers(usecases *Usecases) Handlers {
	return Handlers{
		Health: handler.NewHealthHandler(),
		Flight: handler.NewFlightHandler(&handler.FlightHandlerUsecases{
			FlightUsecase: usecases.FlightUsecase,
		}),
	}
}

func initUsecases(cfg *config.Config) Usecases {
	providerRegistry := provider.NewRegistry(provider.RegistryConfig{
		EnabledProviders: cfg.Providers,
		MockPath: map[string]string{
			"garuda": cfg.GarudaConfig.MockPath,
		},
	})

	return Usecases{
		FlightUsecase: usecase.NewFlightUsecase(providerRegistry),
	}
}

func registerRoutes(r *gin.Engine, h Handlers) {
	registerHealthRoutes(r, h.Health)
	registerFlightRoutes(r, h.Flight)
}
