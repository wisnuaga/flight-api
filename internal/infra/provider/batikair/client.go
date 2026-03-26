package batikair

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
	return "Batik Air"
}

type BatikResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Results []struct {
		FlightNumber      string `json:"flightNumber"`
		AirlineName       string `json:"airlineName"`
		AirlineIATA       string `json:"airlineIATA"`
		Origin            string `json:"origin"`
		Destination       string `json:"destination"`
		DepartureDateTime string `json:"departureDateTime"`
		ArrivalDateTime   string `json:"arrivalDateTime"`
		Fare              struct {
			BasePrice    float64 `json:"basePrice"`
			Taxes        float64 `json:"taxes"`
			TotalPrice   float64 `json:"totalPrice"`
			CurrencyCode string  `json:"currencyCode"`
			Class        string  `json:"class"`
		} `json:"fare"`
		SeatsAvailable int `json:"seatsAvailable"`
	} `json:"results"`
}

func (c *Client) Search(ctx context.Context, req *entity.SearchRequest) ([]*entity.Flight, error) {
	select {
	case <-time.After(300 * time.Millisecond):
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	mockResp, err := util.LoadMock[BatikResponse](c.mockPath)
	if err != nil {
		return nil, err
	}

	var flights []*entity.Flight
	for _, f := range mockResp.Results {
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
			ID:             f.FlightNumber,
			Provider:       "Batik Air",
			FlightNumber:   f.FlightNumber,
			Origin:         f.Origin,
			Destination:    f.Destination,
			DepartureTime:  dep,
			ArrivalTime:    arr,
			Price:          f.Fare.TotalPrice,
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

	return flights, nil
}
