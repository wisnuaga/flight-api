package batikair

import (
	"github.com/shopspring/decimal"
	"github.com/wisnuaga/flight-api/internal/domain/entity"
	"github.com/wisnuaga/flight-api/internal/util"
)

func mapToDomain(resp BatikResponse, req *entity.SearchRequest) []*entity.Flight {
	var flights []*entity.Flight
	for _, f := range resp.Results {
		depTimeUTC, err := util.ParseTimeWithOptionalTZ(f.DepartureDateTime, "")
		if err != nil {
			continue
		}

		arrTimeUTC, err := util.ParseTimeWithOptionalTZ(f.ArrivalDateTime, "")
		if err != nil {
			continue
		}

		if req.Origin != "" && f.Origin != req.Origin {
			continue
		}
		if req.Destination != "" && f.Destination != req.Destination {
			continue
		}

		// Map to domain Flight with timezone-aware locations
		flight := entity.Flight{
			ID:           f.FlightNumber,
			Provider:     "Batik Air",
			FlightNumber: f.FlightNumber,
			Origin: entity.Location{
				Airport:  f.Origin,
				Time:     depTimeUTC,
				Timezone: depTimeUTC.Location(),
			},
			Destination: entity.Location{
				Airport:  f.Destination,
				Time:     arrTimeUTC,
				Timezone: arrTimeUTC.Location(),
			},
			Price:          decimal.NewFromFloat(f.Fare.TotalPrice),
			Currency:       f.Fare.CurrencyCode,
			CabinClass:     f.Fare.Class,
			AvailableSeats: f.SeatsAvailable,
		}

		// Normalize: ensure UTC times, set defaults, compute duration
		flight = entity.NormalizeFlight(flight)

		// Validate flight data
		if !entity.IsValidFlight(flight) {
			continue
		}

		flights = append(flights, &flight)
	}

	return flights
}
