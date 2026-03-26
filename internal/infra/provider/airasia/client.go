package airasia

import (
	"context"
	"errors"
	"math/rand"
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
	return "AirAsia"
}

type AirAsiaResponse struct {
	Status  string `json:"status"`
	Flights []struct {
		FlightCode  string  `json:"flight_code"`
		Airline     string  `json:"airline"`
		FromAirport string  `json:"from_airport"`
		ToAirport   string  `json:"to_airport"`
		DepartTime  string  `json:"depart_time"`
		ArriveTime  string  `json:"arrive_time"`
		PriceIDR    float64 `json:"price_idr"`
		Seats       int     `json:"seats"`
		CabinClass  string  `json:"cabin_class"`
	} `json:"flights"`
}

func (c *Client) Search(ctx context.Context, req *entity.SearchRequest) ([]*entity.Flight, error) {
	if rand.Float32() < 0.10 {
		return nil, errors.New("upstream gateway timeout simulated")
	}

	delayMs := rand.Intn(101) + 50
	select {
	case <-time.After(time.Duration(delayMs) * time.Millisecond):
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	mockResp, err := util.LoadMock[AirAsiaResponse](c.mockPath)
	if err != nil {
		return nil, err
	}

	var flights []*entity.Flight
	for _, f := range mockResp.Flights {
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

	return flights, nil
}
