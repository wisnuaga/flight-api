package airasia

import (
	"github.com/wisnuaga/flight-api/internal/domain/entity"
)

func mapToDomain(resp AirAsiaResponse, req *entity.SearchRequest) []*entity.Flight {
	var flights []*entity.Flight
	for _, f := range resp.Flights {
		dep, err := parseTime(f.DepartTime)
		if err != nil {
			continue
		}
		arr, err := parseTime(f.ArriveTime)
		if err != nil {
			continue
		}

		if req.Origin != "" && f.FromAirport != req.Origin {
			continue
		}
		if req.Destination != "" && f.ToAirport != req.Destination {
			continue
		}

		flight := entity.Flight{
			ID:             f.FlightCode,
			Provider:       "AirAsia",
			FlightNumber:   f.FlightCode,
			Origin:         f.FromAirport,
			Destination:    f.ToAirport,
			DepartureTime:  dep,
			ArrivalTime:    arr,
			Price:          f.PriceIDR,
			Currency:       "IDR",
			CabinClass:     f.CabinClass,
			AvailableSeats: f.Seats,
		}

		flight = entity.NormalizeFlight(flight)
		if !entity.IsValidFlight(flight) {
			continue
		}
		flights = append(flights, &flight)
	}

	return flights
}
