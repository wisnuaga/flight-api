package garuda

import (
	"time"

	"github.com/wisnuaga/flight-api/internal/domain"
)

func mapToDomain(resp GarudaSearchResponse) []*domain.Flight {
	flights := make([]*domain.Flight, 0, len(resp.Flights))

	for _, f := range resp.Flights {
		dep, _ := time.Parse(time.RFC3339, f.Departure.Time)
		arr, _ := time.Parse(time.RFC3339, f.Arrival.Time)

		var price float64
		var currency string
		if f.Price != nil {
			price = float64(f.Price.Amount)
			currency = f.Price.Currency
		}

		flight := &domain.Flight{
			ID:             f.FlightID,
			Provider:       "Garuda",
			FlightNumber:   f.AirlineCode + f.FlightID[len(f.AirlineCode):],
			Origin:         f.Departure.Airport,
			Destination:    f.Arrival.Airport,
			DepartureTime:  dep,
			ArrivalTime:    arr,
			Duration:       time.Duration(f.DurationMinutes) * time.Minute,
			Price:          price,
			Currency:       currency,
			CabinClass:     f.FareClass,
			AvailableSeats: f.AvailableSeats,
		}
		flight.Normalize()
		flights = append(flights, flight)
	}

	return flights
}
