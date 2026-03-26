package lionair

import (
	"github.com/shopspring/decimal"

	"github.com/wisnuaga/flight-api/internal/domain/entity"
	"github.com/wisnuaga/flight-api/internal/util"
)

func mapToDomain(resp LionResponse, req *entity.SearchRequest) []*entity.Flight {
	var flights []*entity.Flight
	for _, f := range resp.Data.AvailableFlights {
		// Lion Air is the only provider that supplies separate IANA timezone names
		// (departure_timezone / arrival_timezone fields). Use ParseTimeWithTZInfo to
		// parse the naive datetime string in that named timezone and preserve it.
		depTimeUTC, depTz, err := util.ParseTimeWithTZInfo(f.Schedule.Departure, f.Schedule.DepartureTimezone)
		if err != nil {
			continue
		}

		arrTimeUTC, arrTz, err := util.ParseTimeWithTZInfo(f.Schedule.Arrival, f.Schedule.ArrivalTimezone)
		if err != nil {
			continue
		}

		if req.Origin != "" && f.Route.From.Code != req.Origin {
			continue
		}
		if req.Destination != "" && f.Route.To.Code != req.Destination {
			continue
		}

		flight := entity.Flight{
			ID:           f.ID,
			Provider:     "Lion Air",
			FlightNumber: f.ID,
			AirlineCode:  f.Carrier.IATA,
			Origin: entity.Location{
				Airport:  f.Route.From.Code,
				Time:     depTimeUTC, // UTC for internal filtering/sorting
				Timezone: depTz,      // Named IANA location from the provider
			},
			Destination: entity.Location{
				Airport:  f.Route.To.Code,
				Time:     arrTimeUTC, // UTC for internal filtering/sorting
				Timezone: arrTz,      // Named IANA location from the provider
			},
			Price:          decimal.NewFromFloat(f.Pricing.Total),
			Currency:       f.Pricing.Currency,
			CabinClass:     f.Pricing.FareType,
			AvailableSeats: f.SeatsLeft,
		}

		flight = entity.NormalizeFlight(flight)

		if !entity.IsValidFlight(flight) {
			continue
		}

		flights = append(flights, &flight)
	}

	return flights
}
