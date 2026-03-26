package router

import (
	"github.com/gin-gonic/gin"
	"github.com/wisnuaga/flight-api/internal/delivery/http/handler"
	"github.com/wisnuaga/flight-api/internal/port"
)

// Setup builds the Gin engine with all routes registered.
// All dependencies must be constructed and injected by the caller (cmd/main.go).
func Setup(flightUsecase port.FlightUsecase) *gin.Engine {
	r := gin.Default()

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
