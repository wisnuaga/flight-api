package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wisnuaga/flight-api/internal/delivery/http/dto"
	"github.com/wisnuaga/flight-api/internal/usecase"
)

type FlightHandlerUsecases struct {
	FlightUsecase usecase.FlightUsecase
}

type FlightHandler struct {
	Usecases *FlightHandlerUsecases
}

func NewFlightHandler(usecases *FlightHandlerUsecases) *FlightHandler {
	return &FlightHandler{Usecases: usecases}
}

func (h *FlightHandler) Search(c *gin.Context) {
	// Read req body
	var req dto.SearchRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	// Convert to domain
	domainReq, err := req.ToDomain()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid date format",
		})
		return
	}

	flights, err := h.Usecases.FlightUsecase.Search(c, &domainReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, flights)
}
