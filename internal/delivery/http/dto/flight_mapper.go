package dto

import (
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
