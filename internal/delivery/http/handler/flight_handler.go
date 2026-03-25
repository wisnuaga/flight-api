package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/wisnuaga/flight-api/internal/delivery/http/dto"
)

type FlightHandler struct{}

func NewFlightHandler() *FlightHandler {
	return &FlightHandler{}
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

	now := time.Now()

	resp := dto.SearchResponse{
		SearchCriteria: dto.SearchCriteria{
			Origin:        domainReq.Origin,
			Destination:   domainReq.Destination,
			DepartureDate: "2025-12-15",
			Passengers:    1,
			CabinClass:    "economy",
		},
		Metadata: dto.Metadata{
			TotalResults:       15,
			ProvidersQueried:   4,
			ProvidersSucceeded: 4,
			ProvidersFailed:    0,
			SearchTimeMs:       285,
			CacheHit:           false,
		},
		Flights: []dto.Flight{
			{
				ID:       "QZ7250_AirAsia",
				Provider: "AirAsia",
				Airline: dto.Airline{
					Name: "AirAsia",
					Code: "QZ",
				},
				FlightNumber: "QZ7250",
				Departure: dto.Location{
					Airport:   "CGK",
					City:      "Jakarta",
					Datetime:  "2025-12-15T15:15:00+07:00",
					Timestamp: now.Unix(),
				},
				Arrival: dto.Location{
					Airport:   "DPS",
					City:      "Denpasar",
					Datetime:  "2025-12-15T20:35:00+08:00",
					Timestamp: now.Add(4*time.Hour + 20*time.Minute).Unix(),
				},
				Duration: dto.Duration{
					TotalMinutes: 260,
					Formatted:    "4h 20m",
				},
				Stops: 1,
				Price: dto.Price{
					Amount:   485000,
					Currency: "IDR",
				},
				AvailableSeats: 88,
				CabinClass:     "economy",
				Aircraft:       nil,
				Amenities:      []string{},
				Baggage: dto.Baggage{
					CarryOn: "Cabin baggage only",
					Checked: "Additional fee",
				},
			},
		},
	}

	c.JSON(http.StatusOK, resp)
}
