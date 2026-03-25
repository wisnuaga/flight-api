package garuda

import (
	"context"
	"time"

	"github.com/wisnuaga/flight-api/internal/domain"
	"github.com/wisnuaga/flight-api/utils"
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

func (c *Client) Search(ctx context.Context, req *domain.SearchRequest) ([]*domain.Flight, error) {
	select {
	case <-time.After(80 * time.Millisecond):
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	mockResp, err := utils.LoadMock[GarudaSearchResponse](c.mockPath)
	if err != nil {
		return nil, err
	}

	return mapToDomain(mockResp), nil
}
