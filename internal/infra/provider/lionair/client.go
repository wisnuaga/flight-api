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

	flights := mapToDomain(mockResp, req)

	return flights, nil
}
