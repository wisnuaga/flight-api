package lionair

import (
	"context"
	"time"

	"github.com/wisnuaga/flight-api/internal/domain/entity"
	"github.com/wisnuaga/flight-api/internal/util"
)

type Client struct {
	mockPath string
}

func NewClient(mockPath string) *Client {
	return &Client{mockPath: mockPath}
}

func (c *Client) Name() string {
	return "Lion Air"
}

type LionResponse struct {
	Success bool `json:"success"`
	Data    struct {
		AvailableFlights []struct {
			ID      string `json:"id"`
			Carrier struct {
				Name string `json:"name"`
				IATA string `json:"iata"`
			} `json:"carrier"`
			Route struct {
				From struct {
					Code string `json:"code"`
				} `json:"from"`
				To struct {
					Code string `json:"code"`
				} `json:"to"`
			} `json:"route"`
			Schedule struct {
				Departure         string `json:"departure"`
				DepartureTimezone string `json:"departure_timezone"`
				Arrival           string `json:"arrival"`
				ArrivalTimezone   string `json:"arrival_timezone"`
			} `json:"schedule"`
			Pricing struct {
				Total    float64 `json:"total"`
				Currency string  `json:"currency"`
				FareType string  `json:"fare_type"`
			} `json:"pricing"`
			SeatsLeft int `json:"seats_left"`
		} `json:"available_flights"`
	} `json:"data"`
}

func (c *Client) Search(ctx context.Context, req *entity.SearchRequest) ([]*entity.Flight, error) {
	select {
	case <-time.After(150 * time.Millisecond):
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	mockResp, err := util.LoadMock[LionResponse](c.mockPath)
	if err != nil {
		return nil, err
	}

	var flights []*entity.Flight
	for _, f := range mockResp.Data.AvailableFlights {
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

	return flights, nil
}
