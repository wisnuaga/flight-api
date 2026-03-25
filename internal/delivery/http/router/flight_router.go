package router

import (
	"net/http"

	"github.com/wisnuaga/flight-api/internal/delivery/http/handler"
)

func registerFlightRoutes(mux *http.ServeMux, h *handler.FlightHandler) {
	mux.HandleFunc("/flights/search", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		h.Search(w, r)
	})
}
