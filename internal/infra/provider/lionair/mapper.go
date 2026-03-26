package lionair

import (
	"time"

	"github.com/wisnuaga/flight-api/internal/domain/entity"
)

func mapToDomain(resp LionResponse, req *entity.SearchRequest) []*entity.Flight {
	var flights []*entity.Flight
	for _, f := range resp.Data.AvailableFlights {
		depLoc, _ := time.LoadLocation(f.Schedule.DepartureTimezone)
		dep, err := time.ParseInLocation("2006-01-02T15:04:05", f.Schedule.Departure, depLoc)
		if err != nil {
			continue
		}

		arrLoc, _ := time.LoadLocation(f.Schedule.ArrivalTimezone)
		arr, err := time.ParseInLocation("2006-01-02T15:04:05", f.Schedule.Arrival, arrLoc)
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
			ID:             f.ID,
			Provider:       "Lion Air",
			FlightNumber:   f.ID,
			Origin:         f.Route.From.Code,
			Destination:    f.Route.To.Code,
			DepartureTime:  dep,
			ArrivalTime:    arr,
			Price:          f.Pricing.Total,
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
