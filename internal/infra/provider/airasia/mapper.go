package airasia

import (
	"github.com/shopspring/decimal"
	"github.com/wisnuaga/flight-api/internal/domain/entity"
	"github.com/wisnuaga/flight-api/internal/util"
)

func mapToDomain(resp AirAsiaResponse, req *entity.SearchRequest) []*entity.Flight {
	var flights []*entity.Flight
	for _, f := range resp.Flights {
		depTimeUTC, err := util.ParseTimeWithOptionalTZ(f.DepartTime, "")
		if err != nil {
			continue
		}

		arrTimeUTC, err := util.ParseTimeWithOptionalTZ(f.ArriveTime, "")
		if err != nil {
			continue
		}

		if req.Origin != "" && f.FromAirport != req.Origin {
			continue
		}
		if req.Destination != "" && f.ToAirport != req.Destination {
			continue
		}

		// Map to domain Flight with timezone-aware locations
		flight := entity.Flight{
			ID:           f.FlightCode,
			Provider:     "AirAsia",
			FlightNumber: f.FlightCode,
			Origin: entity.Location{
				Airport:  f.FromAirport,
				Time:     depTimeUTC,
				Timezone: depTimeUTC.Location(),
			},
			Destination: entity.Location{
				Airport:  f.ToAirport,
				Time:     arrTimeUTC,
				Timezone: arrTimeUTC.Location(),
			},
			Price:          decimal.NewFromFloat(f.PriceIDR),
			Currency:       "IDR",
			CabinClass:     f.CabinClass,
			AvailableSeats: f.Seats,
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
