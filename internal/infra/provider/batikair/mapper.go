package batikair

import (
	"fmt"

	"github.com/shopspring/decimal"
	"github.com/wisnuaga/flight-api/internal/domain/entity"
	"github.com/wisnuaga/flight-api/internal/util"
)

func mapToDomain(resp BatikResponse, req *entity.SearchRequest) []*entity.Flight {
	var flights []*entity.Flight
	for _, f := range resp.Results {
		// Batik Air embeds the UTC offset in the time string (e.g. "2025-12-15T07:15:00+0700").
		// Extract the instant (UTC) and fixed-offset location from the string itself.
		depTimeUTC, depTz, err := util.ParseTimeFromString(f.DepartureDateTime)
		if err != nil {
			continue
		}

		arrTimeUTC, arrTz, err := util.ParseTimeFromString(f.ArrivalDateTime)
		if err != nil {
			continue
		}

		if req.Origin != "" && f.Origin != req.Origin {
			continue
		}
		if req.Destination != "" && f.Destination != req.Destination {
			continue
		}

		flight := entity.Flight{
			ID:           fmt.Sprintf("%s_%s", f.FlightNumber, util.NormalizeAirlineName(f.AirlineName)),
			Provider:     f.AirlineName,
			FlightNumber: f.FlightNumber,
			AirlineCode:  f.AirlineIATA,
			Origin: entity.Location{
				Airport:  f.Origin,
				Time:     depTimeUTC, // UTC for internal filtering/sorting
				Timezone: depTz,      // Fixed-offset location extracted from the time string
			},
			Destination: entity.Location{
				Airport:  f.Destination,
				Time:     arrTimeUTC, // UTC for internal filtering/sorting
				Timezone: arrTz,      // Fixed-offset location extracted from the time string
			},
			Price:          decimal.NewFromFloat(f.Fare.TotalPrice),
			Currency:       f.Fare.CurrencyCode,
			CabinClass:     f.Fare.Class,
			AvailableSeats: f.SeatsAvailable,
		}

		flight = entity.NormalizeFlight(flight)

		if !entity.IsValidFlight(flight) {
			continue
		}

		flights = append(flights, &flight)
	}

	return flights
}
