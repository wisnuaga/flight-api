package router

import (
	"github.com/gin-gonic/gin"
	"github.com/wisnuaga/flight-api/internal/delivery/http/handler"
)

func registerFlightRoutes(r *gin.Engine, h *handler.FlightHandler) {
	flights := r.Group("/flights")
	flights.POST("/search", h.Search)
}
