package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/wisnuaga/flight-api/internal/delivery/http/dto"
	"github.com/wisnuaga/flight-api/internal/port"
)

type FlightHandlerUsecases struct {
	FlightUsecase port.FlightUsecase
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

	// Convert to domain entity
	domainReq, err := req.ToDomain()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid date format",
		})
		return
	}

	// The Search method handles both one-way and round-trip searches
	// based on whether ReturnDate is provided
	result, err := h.Usecases.FlightUsecase.Search(c, &domainReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Check if this is a round-trip search based on the ReturnDate
	if domainReq.ReturnDate != nil {
		// Round-trip search - include round-trip itineraries in response
		c.JSON(http.StatusOK, dto.ToRoundTripResponseWithMeta(&req, result.RoundTripItineraries, result))
	} else {
		// One-way search
		c.JSON(http.StatusOK, dto.ToSearchResponse(&req, result))
	}
}
