package lionair

import (
	"context"
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
	return "Lion Air"
}

func (c *Client) Search(ctx context.Context, req *entity.SearchRequest) ([]*entity.Flight, error) {
	// Random delay between 100ms and 200ms
	delay := time.Duration(100+rand.Intn(101)) * time.Millisecond // rand.Intn(101) -> 0..101

	select {
	case <-time.After(delay):
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	mockResp, err := util.LoadMock[LionResponse](c.mockPath)
	if err != nil {
		return nil, err
	}

	flights := mapToDomain(mockResp, req)

	return flights, nil
}
