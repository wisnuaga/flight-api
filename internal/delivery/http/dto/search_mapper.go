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

func convertAirlineNames(airlineStrs []string) []entity.AirlineName {
	if len(airlineStrs) == 0 {
		return nil
	}
	airlines := make([]entity.AirlineName, 0, len(airlineStrs))
	for _, s := range airlineStrs {
		if airline := entity.AirlineNameFromString(s); airline != "" {
			airlines = append(airlines, airline)
		}
	}
	return airlines
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

	// Parse optional return date for round-trip searches
	var returnDate *time.Time
	if r.ReturnDate != nil && *r.ReturnDate != "" {
		returnTime, err := time.Parse("2006-01-02", *r.ReturnDate)
		if err == nil {
			returnDate = &returnTime
		}
	}

	return entity.SearchRequest{
		Origin:        r.Origin,
		Destination:   r.Destination,
		DepartureDate: departureTime,
		ReturnDate:    returnDate,
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
			Airlines:       convertAirlineNames(r.Airlines),
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
				Provider: f.Airline.String(),
				Airline: Airline{
					Name: f.Airline.String(),
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
					Formatted: util.FormatPriceDecimal(f.Price, f.Currency),
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
		City:      loc.City,
		Datetime:  datetime,
		Timestamp: timestamp,
	}
}

// ToRoundTripResponse converts round-trip itineraries to a response DTO
func ToRoundTripResponse(req *SearchRequest, itineraries []*entity.RoundTripItinerary) RoundTripResponse {
	return ToRoundTripResponseWithMeta(req, itineraries, nil)
}

// ToRoundTripResponseWithMeta converts round-trip itineraries with metadata to a response DTO
func ToRoundTripResponseWithMeta(req *SearchRequest, itineraries []*entity.RoundTripItinerary, result *entity.SearchResult) RoundTripResponse {
	var roundTrips []RoundTripItinerary

	for _, itinerary := range itineraries {
		// Format outbound flight
		outboundDurationMins := int(itinerary.OutboundFlight.TotalTripDuration().Minutes())
		outboundHours := outboundDurationMins / 60
		outboundMins := outboundDurationMins % 60

		outboundFlight := Flight{
			ID:       itinerary.OutboundFlight.ID,
			Provider: itinerary.OutboundFlight.Airline.String(),
			Airline: Airline{
				Name: itinerary.OutboundFlight.Airline.String(),
				Code: itinerary.OutboundFlight.AirlineCode,
			},
			FlightNumber: itinerary.OutboundFlight.FlightNumber,
			Departure:    ConvertLocationToResponse(itinerary.OutboundFlight.Origin),
			Arrival:      ConvertLocationToResponse(itinerary.OutboundFlight.Destination),
			Duration: Duration{
				TotalMinutes: outboundDurationMins,
				Formatted:    fmt.Sprintf("%dh %dm", outboundHours, outboundMins),
			},
			Stops: itinerary.OutboundFlight.Stops,
			Price: Price{
				Amount:    itinerary.OutboundFlight.Price.IntPart(),
				Currency:  itinerary.OutboundFlight.Currency,
				Formatted: util.FormatPriceDecimal(itinerary.OutboundFlight.Price, itinerary.OutboundFlight.Currency),
			},
			AvailableSeats: itinerary.OutboundFlight.AvailableSeats,
			CabinClass:     itinerary.OutboundFlight.CabinClass,
		}

		// Format return flight
		returnDurationMins := int(itinerary.ReturnFlight.TotalTripDuration().Minutes())
		returnHours := returnDurationMins / 60
		returnMins := returnDurationMins % 60

		returnFlight := Flight{
			ID:       itinerary.ReturnFlight.ID,
			Provider: itinerary.ReturnFlight.Airline.String(),
			Airline: Airline{
				Name: itinerary.ReturnFlight.Airline.String(),
				Code: itinerary.ReturnFlight.AirlineCode,
			},
			FlightNumber: itinerary.ReturnFlight.FlightNumber,
			Departure:    ConvertLocationToResponse(itinerary.ReturnFlight.Origin),
			Arrival:      ConvertLocationToResponse(itinerary.ReturnFlight.Destination),
			Duration: Duration{
				TotalMinutes: returnDurationMins,
				Formatted:    fmt.Sprintf("%dh %dm", returnHours, returnMins),
			},
			Stops: itinerary.ReturnFlight.Stops,
			Price: Price{
				Amount:    itinerary.ReturnFlight.Price.IntPart(),
				Currency:  itinerary.ReturnFlight.Currency,
				Formatted: util.FormatPriceDecimal(itinerary.ReturnFlight.Price, itinerary.ReturnFlight.Currency),
			},
			AvailableSeats: itinerary.ReturnFlight.AvailableSeats,
			CabinClass:     itinerary.ReturnFlight.CabinClass,
		}

		// Format total duration
		totalDurationMins := int(itinerary.TotalDuration.Minutes())
		totalHours := totalDurationMins / 60
		totalMins := totalDurationMins % 60

		roundTrips = append(roundTrips, RoundTripItinerary{
			OutboundFlight: outboundFlight,
			ReturnFlight:   returnFlight,
			TotalPrice: Price{
				Amount:    itinerary.TotalPrice.IntPart(),
				Currency:  itinerary.OutboundFlight.Currency,
				Formatted: util.FormatPriceDecimal(itinerary.TotalPrice, itinerary.OutboundFlight.Currency),
			},
			TotalDuration: Duration{
				TotalMinutes: totalDurationMins,
				Formatted:    fmt.Sprintf("%dh %dm", totalHours, totalMins),
			},
		})
	}

	meta := Metadata{}
	if result != nil && result.Meta != nil {
		meta = Metadata{
			TotalResults:       result.Meta.TotalFlights,
			ProvidersQueried:   result.Meta.Providers,
			ProvidersSucceeded: result.Meta.SuccessCount,
			ProvidersFailed:    result.Meta.FailedCount,
			SearchTimeMs:       result.Meta.SearchTimeMs,
			CacheHit:           result.Meta.CacheHit,
		}
	}

	return RoundTripResponse{
		SearchCriteria: SearchCriteria{
			Origin:        req.Origin,
			Destination:   req.Destination,
			DepartureDate: req.DepartureDate,
			ReturnDate:    req.ReturnDate,
			Passengers:    req.Passengers,
			CabinClass:    req.CabinClass,
		},
		Metadata:    meta,
		Itineraries: roundTrips,
	}
}
