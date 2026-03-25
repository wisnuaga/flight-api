package router

import (
	"net/http"

	"github.com/wisnuaga/flight-api/internal/delivery/http/handler"
)

type Handlers struct {
	Flight *handler.FlightHandler
}

func Setup(mux *http.ServeMux) {
	handlers := initHandlers()
	registerRoutes(mux, handlers)
}

func initHandlers() Handlers {
	return Handlers{
		Flight: handler.NewFlightHandler(),
	}
}

func registerRoutes(mux *http.ServeMux, h Handlers) {
	registerHealthRoutes(mux)
	registerFlightRoutes(mux, h.Flight)
}
