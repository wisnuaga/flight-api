package garuda

import (
	"github.com/shopspring/decimal"
	"github.com/wisnuaga/flight-api/internal/domain/entity"
	"github.com/wisnuaga/flight-api/internal/util"
)

func mapToDomain(resp GarudaSearchResponse) []*entity.Flight {
	var flights []*entity.Flight

	for _, f := range resp.Flights {
		dep, err := util.ParseTime(f.Departure.Time)
		if err != nil {
			continue
		}

		arr, err := util.ParseTime(f.Arrival.Time)
		if err != nil {
			continue
		}

		var price float64
		var currency string
		if f.Price != nil {
			price = float64(f.Price.Amount)
			currency = f.Price.Currency
		}

		// Initial Raw Flight Mapping
		flight := entity.Flight{
			ID:             f.FlightID,
			Provider:       "Garuda",
			FlightNumber:   f.AirlineCode + f.FlightID[len(f.AirlineCode):],
			Origin:         f.Departure.Airport,
			Destination:    f.Arrival.Airport,
			DepartureTime:  dep,
			ArrivalTime:    arr,
			Price:          decimal.NewFromFloat(price),
			Currency:       currency,
			CabinClass:     f.FareClass,
			AvailableSeats: f.AvailableSeats,
		}

		// Let domain rules normalize basic values and enforce duration calculation
		flight = entity.NormalizeFlight(flight)

		// Hard drop invalid/malformed response payload flight items
		if !entity.IsValidFlight(flight) {
			continue
		}

		flights = append(flights, &flight)
	}

	return flights
}
