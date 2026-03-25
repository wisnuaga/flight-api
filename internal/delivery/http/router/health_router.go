package router

import (
	"github.com/gin-gonic/gin"
	"github.com/wisnuaga/flight-api/internal/delivery/http/handler"
)

func registerHealthRoutes(r *gin.Engine, h *handler.HealthHandler) {
	// health endpoint
	r.GET("/healthz", h.Check)
}
