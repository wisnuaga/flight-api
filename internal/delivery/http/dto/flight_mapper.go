package dto

import (
	"fmt"
	"time"

	"github.com/wisnuaga/flight-api/internal/domain"
)

func (r *SearchRequest) ToDomain() (domain.SearchRequest, error) {
	departureTime, err := time.Parse("2006-01-02", r.DepartureDate)
	if err != nil {
		return domain.SearchRequest{}, err
	}

	return domain.SearchRequest{
		Origin:        r.Origin,
		Destination:   r.Destination,
		DepartureDate: departureTime,
		Passengers:    r.Passengers,
		CabinClass:    r.CabinClass,
	}, nil
}

func ToSearchResponse(req *SearchRequest, result *domain.SearchResult) SearchResponse {
	var flights []Flight
	if result.Flights != nil {
		for _, f := range result.Flights {
			durationMins := int(f.Duration.Minutes())
			hours := durationMins / 60
			mins := durationMins % 60

			flights = append(flights, Flight{
				ID:       f.ID,
				Provider: f.Provider,
				Airline: Airline{
					Name: f.Provider,
					Code: f.Provider,
				},
				FlightNumber: f.FlightNumber,
				Departure: Location{
					Airport:   f.Origin,
					Datetime:  f.DepartureTime.Format(time.RFC3339),
					Timestamp: f.DepartureTime.Unix(),
				},
				Arrival: Location{
					Airport:   f.Destination,
					Datetime:  f.ArrivalTime.Format(time.RFC3339),
					Timestamp: f.ArrivalTime.Unix(),
				},
				Duration: Duration{
					TotalMinutes: durationMins,
					Formatted:    fmt.Sprintf("%dh %dm", hours, mins),
				},
				Stops: 0,
				Price: Price{
					Amount:   int64(f.Price),
					Currency: f.Currency,
				},
				AvailableSeats: f.AvailableSeats,
				CabinClass:     f.CabinClass,
				Amenities:      make([]string, 0),
			})
		}
	}

	if flights == nil {
		flights = make([]Flight, 0)
	}

	meta := Metadata{}
	if result.Meta != nil {
		meta = Metadata{
			TotalResults:       result.Meta.TotalFlights,
			ProvidersQueried:   result.Meta.Providers,
			ProvidersSucceeded: result.Meta.SuccessCount,
			ProvidersFailed:    result.Meta.FailedCount,
		}
	}

	return SearchResponse{
		SearchCriteria: SearchCriteria{
			Origin:        req.Origin,
			Destination:   req.Destination,
			DepartureDate: req.DepartureDate,
			Passengers:    req.Passengers,
			CabinClass:    req.CabinClass,
		},
		Metadata: meta,
		Flights:  flights,
	}
}
