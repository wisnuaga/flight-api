package airasia

import (
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"github.com/wisnuaga/flight-api/internal/domain/entity"
	"github.com/wisnuaga/flight-api/internal/util"
)

func mapToDomain(resp SearchResponse, req *entity.SearchRequest) []*entity.Flight {
	var flights []*entity.Flight
	for _, f := range resp.Flights {
		// AirAsia embeds the UTC offset in the time string (e.g. "2025-12-15T04:45:00+07:00").
		// Extract the instant (UTC) and fixed-offset location from the string itself.
		depTimeUTC, depTz, err := util.ParseTimeFromString(f.DepartTime)
		if err != nil {
			continue
		}

		arrTimeUTC, arrTz, err := util.ParseTimeFromString(f.ArriveTime)
		if err != nil {
			continue
		}

		if req.Origin != "" && f.FromAirport != req.Origin {
			continue
		}
		if req.Destination != "" && f.ToAirport != req.Destination {
			continue
		}

		layovers := []*entity.Layover{}
		for _, stop := range f.Stops {
			layovers = append(layovers, &entity.Layover{
				Airport:  stop.Airport,
				Duration: time.Duration(stop.WaitTimeMinutes) * time.Minute,
			})
		}

		flight := entity.Flight{
			ID:           fmt.Sprintf("%s_%s", f.FlightCode, util.NormalizeAirlineName(f.Airline)),
			Provider:     f.Airline,
			FlightNumber: f.FlightCode,
			AirlineCode:  util.GetFlightCodePrefix(f.FlightCode),
			Origin: entity.Location{
				Airport:  f.FromAirport,
				Time:     depTimeUTC, // UTC for internal filtering/sorting
				Timezone: depTz,      // Fixed-offset location extracted from the time string
			},
			Destination: entity.Location{
				Airport:  f.ToAirport,
				Time:     arrTimeUTC, // UTC for internal filtering/sorting
				Timezone: arrTz,      // Fixed-offset location extracted from the time string
			},
			Price:          decimal.NewFromFloat(f.PriceIDR),
			Currency:       "IDR",
			CabinClass:     f.CabinClass,
			AvailableSeats: f.Seats,
			Stops:          len(layovers),
			Layovers:       layovers,
		}

		flight = entity.NormalizeFlight(flight)

		if !entity.IsValidFlight(flight) {
			continue
		}

		flights = append(flights, &flight)
	}

	return flights
}
