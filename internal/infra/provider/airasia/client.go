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

	flights := mapToDomain(mockResp, req)

	return flights, nil
}
