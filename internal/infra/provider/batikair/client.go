package batikair

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
	return "Batik Air"
}

func (c *Client) Search(ctx context.Context, req *entity.SearchRequest) ([]*entity.Flight, error) {
	// Random delay between 200ms and 400ms
	delay := time.Duration(200+rand.Intn(201)) * time.Millisecond // rand.Intn(201) -> 0..200

	select {
	case <-time.After(delay):
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	mockResp, err := util.LoadMock[SearchResponse](c.mockPath)
	if err != nil {
		return nil, err
	}

	flights := mapToDomain(mockResp, req)

	return flights, nil
}
