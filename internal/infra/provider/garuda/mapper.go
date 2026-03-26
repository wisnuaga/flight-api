package garuda

import (
	"github.com/shopspring/decimal"

	"github.com/wisnuaga/flight-api/internal/domain/entity"
	"github.com/wisnuaga/flight-api/internal/util"
)

func mapToDomain(resp GarudaSearchResponse) []*entity.Flight {
	var flights []*entity.Flight
	for _, f := range resp.Flights {
		// Guard against nil pointers from the provider
		if f.Departure == nil || f.Arrival == nil {
			continue
		}

		// Garuda embeds the UTC offset directly in the time string (e.g. "2025-12-15T06:00:00+07:00").
		// Extract the instant (UTC) and the fixed-offset location from the string itself.
		depTimeUTC, depTz, err := util.ParseTimeFromString(f.Departure.Time)
		if err != nil {
			continue
		}
		arrTimeUTC, arrTz, err := util.ParseTimeFromString(f.Arrival.Time)
		if err != nil {
			continue
		}

		var price float64
		var currency string
		if f.Price != nil {
			price = float64(f.Price.Amount)
			currency = f.Price.Currency
		}

		flight := entity.Flight{
			ID:           f.FlightID,
			Provider:     "Garuda",
			FlightNumber: f.FlightID,
			AirlineCode:  f.AirlineCode,
			Origin: entity.Location{
				Airport:  f.Departure.Airport,
				City:     f.Departure.City,
				Time:     depTimeUTC, // UTC for internal filtering/sorting
				Timezone: depTz,      // Fixed-offset location extracted from the time string
			},
			Destination: entity.Location{
				Airport:  f.Arrival.Airport,
				City:     f.Arrival.City,
				Time:     arrTimeUTC, // UTC for internal filtering/sorting
				Timezone: arrTz,      // Fixed-offset location extracted from the time string
			},
			Price:          decimal.NewFromFloat(price),
			Currency:       currency,
			CabinClass:     f.FareClass,
			AvailableSeats: f.AvailableSeats,
		}

		flight = entity.NormalizeFlight(flight)

		if !entity.IsValidFlight(flight) {
			continue
		}

		flights = append(flights, &flight)
	}

	return flights
}
