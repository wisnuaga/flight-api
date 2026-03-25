package router

import (
	"net/http"

	"github.com/wisnuaga/flight-api/internal/delivery/http/handler"
)

func RegisterFlightRoutes(mux *http.ServeMux, h *handler.FlightHandler) {
	mux.HandleFunc("/flights/search", h.Search)
}
