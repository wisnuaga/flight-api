package router

import (
	"github.com/gin-gonic/gin"
	"github.com/wisnuaga/flight-api/internal/delivery/http/handler"
)

type Handlers struct {
	Health *handler.HealthHandler
	Flight *handler.FlightHandler
}

func Setup() *gin.Engine {
	r := gin.Default()

	// Init handlers
	handlers := initHandlers()
	registerRoutes(r, handlers)

	return r
}

func initHandlers() Handlers {
	return Handlers{
		Health: handler.NewHealthHandler(),
		Flight: handler.NewFlightHandler(),
	}
}

func registerRoutes(r *gin.Engine, h Handlers) {
	registerHealthRoutes(r, h.Health)
	registerFlightRoutes(r, h.Flight)
}
