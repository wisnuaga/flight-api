package lionair

import (
	"github.com/shopspring/decimal"

	"github.com/wisnuaga/flight-api/internal/domain/entity"
	"github.com/wisnuaga/flight-api/internal/util"
)

func mapToDomain(resp LionResponse, req *entity.SearchRequest) []*entity.Flight {
	var flights []*entity.Flight
	for _, f := range resp.Data.AvailableFlights {
		// Parse departure time with timezone info (Lion Air provides timezone data)
		depTimeUTC, err := util.ParseTimeWithOptionalTZ(f.Schedule.Departure, f.Schedule.DepartureTimezone)
		if err != nil {
			continue
		}

		// Parse arrival time with timezone info
		arrTimeUTC, err := util.ParseTimeWithOptionalTZ(f.Schedule.Arrival, f.Schedule.ArrivalTimezone)
		if err != nil {
			continue
		}

		if req.Origin != "" && f.Route.From.Code != req.Origin {
			continue
		}
		if req.Destination != "" && f.Route.To.Code != req.Destination {
			continue
		}

		// Map to domain Flight with timezone-aware locations
		flight := entity.Flight{
			ID:           f.ID,
			Provider:     "Lion Air",
			FlightNumber: f.ID,
			Origin: entity.Location{
				Airport:  f.Route.From.Code,
				Time:     depTimeUTC,
				Timezone: depTimeUTC.Location(),
			},
			Destination: entity.Location{
				Airport:  f.Route.To.Code,
				Time:     arrTimeUTC,
				Timezone: arrTimeUTC.Location(),
			},
			Price:          decimal.NewFromFloat(f.Pricing.Total),
			Currency:       f.Pricing.Currency,
			CabinClass:     f.Pricing.FareType,
			AvailableSeats: f.SeatsLeft,
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
