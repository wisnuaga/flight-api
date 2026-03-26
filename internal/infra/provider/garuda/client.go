package garuda

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
	return "Garuda"
}

func (c *Client) Search(ctx context.Context, req *entity.SearchRequest) ([]*entity.Flight, error) {
	// Random delay between 50ms and 100ms
	delay := time.Duration(50+rand.Intn(51)) * time.Millisecond // rand.Intn(51) -> 0..50

	select {
	case <-time.After(delay):
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	mockResp, err := util.LoadMock[SearchResponse](c.mockPath)
	if err != nil {
		return nil, err
	}

	return mapToDomain(mockResp), nil
}
