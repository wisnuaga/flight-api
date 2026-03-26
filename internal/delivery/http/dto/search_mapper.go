package dto

import (
	"fmt"
	"time"

	"github.com/wisnuaga/flight-api/internal/domain/entity"
	"github.com/wisnuaga/flight-api/internal/util"
)

func parseOptionalTime(t *string) *time.Time {
	if t == nil {
		return nil
	}
	parsed, err := util.ParseTime(*t)
	if err != nil {
		return nil
	}
	// Always normalise to UTC so filter comparisons are consistent with
	// the UTC times stored on Origin.Time / Destination.Time.
	utc := parsed.UTC()
	return &utc
}

func (r *SearchRequest) ToDomain() (entity.SearchRequest, error) {
	departureTime, err := time.Parse("2006-01-02", r.DepartureDate)
	if err != nil {
		return entity.SearchRequest{}, err
	}

	var cabinClass *string
	if r.CabinClass != "" {
		cabinClass = &r.CabinClass
	}

	var maxDuration *time.Duration
	if r.MaxDuration != nil {
		d := time.Duration(*r.MaxDuration) * time.Minute
		maxDuration = &d
	}

	return entity.SearchRequest{
		Origin:        r.Origin,
		Destination:   r.Destination,
		DepartureDate: departureTime,
		Passengers:    r.Passengers,
		Filter: entity.SearchFilter{
			MinPrice:       r.MinPrice,
			MaxPrice:       r.MaxPrice,
			MaxStops:       r.MaxStops,
			DepartureStart: parseOptionalTime(r.DepartureStart),
			DepartureEnd:   parseOptionalTime(r.DepartureEnd),
			ArrivalStart:   parseOptionalTime(r.ArrivalStart),
			ArrivalEnd:     parseOptionalTime(r.ArrivalEnd),
			MaxDuration:    maxDuration,
			AirlineCodes:   r.AirlineCodes,
			CabinClass:     cabinClass,
		},
		Sort: entity.SearchSort{
			Field: entity.SortField(r.SortBy),
			Order: entity.SortOrder(r.SortOrder),
		},
	}, nil
}

func ToSearchResponse(req *SearchRequest, result *entity.SearchResult) SearchResponse {
	var flights []Flight
	if result.Flights != nil {
		for _, f := range result.Flights {
			durationMins := int(f.TotalTripDuration().Minutes())
			hours := durationMins / 60
			mins := durationMins % 60

			// Convert locations to response format with timezone awareness
			departure := ConvertLocationToResponse(f.Origin)
			arrival := ConvertLocationToResponse(f.Destination)

			flights = append(flights, Flight{
				ID:       f.ID,
				Provider: f.Provider,
				Airline: Airline{
					Name: f.Provider,
					Code: f.AirlineCode,
				},
				FlightNumber: f.FlightNumber,
				Departure:    departure,
				Arrival:      arrival,
				Duration: Duration{
					TotalMinutes: durationMins,
					Formatted:    fmt.Sprintf("%dh %dm", hours, mins),
				},
				Stops: f.Stops,
				Price: Price{
					Amount:    f.Price.IntPart(),
					Currency:  f.Currency,
					Formatted: fmt.Sprintf("%s %d", f.Currency, f.Price.IntPart()),
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
			SearchTimeMs:       result.Meta.SearchTimeMs,
			CacheHit:           result.Meta.CacheHit,
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

func ConvertLocationToResponse(loc entity.Location) Location {
	// Default to UTC if timezone is not set
	tz := loc.Timezone
	if tz == nil {
		tz = time.UTC
	}

	// Convert UTC time to the original timezone for display
	localTime := loc.Time.In(tz)

	// Format datetime with timezone offset (RFC3339 format)
	datetime := localTime.Format(time.RFC3339)

	// Unix timestamp (seconds since epoch) in UTC
	timestamp := loc.Time.Unix()

	return Location{
		Airport:   loc.Airport,
		City:      loc.City, // populated from entity after Normalize() sets it from airport code map
		Datetime:  datetime,
		Timestamp: timestamp,
	}
}
