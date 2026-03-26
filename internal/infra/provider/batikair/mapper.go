package batikair

import (
	"github.com/shopspring/decimal"
	"github.com/wisnuaga/flight-api/internal/domain/entity"
	"github.com/wisnuaga/flight-api/internal/util"
)

func mapToDomain(resp BatikResponse, req *entity.SearchRequest) []*entity.Flight {
	var flights []*entity.Flight
	for _, f := range resp.Results {
		dep, err := util.ParseTime(f.DepartureDateTime)
		if err != nil {
			continue
		}
		arr, err := util.ParseTime(f.ArrivalDateTime)
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
			ID:           f.FlightNumber,
			Provider:     "Batik Air",
			FlightNumber: f.FlightNumber,
			Origin: entity.Location{
				Airport: f.Origin,
			},
			Destination: entity.Location{
				Airport: f.Destination,
			},
			DepartureTime:  dep,
			ArrivalTime:    arr,
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
