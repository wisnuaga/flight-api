package command_test

import (
	"testing"
	"time"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/wisnuaga/flight-api/internal/command"
	"github.com/wisnuaga/flight-api/internal/domain/entity"
)

func ptrDecimal(d float64) *decimal.Decimal { a := decimal.NewFromFloat(d); return &a }
func ptrInt(i int) *int                     { return &i }
func ptrString(s string) *string            { return &s }

func TestFlightFilterCommand_Execute(t *testing.T) {
	cmd := command.NewFlightFilterCommand()
	baseTime := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)

	flights := []*entity.Flight{
		{
			ID: "1", Price: decimal.NewFromInt(1000), Stops: 0, Provider: "Garuda", CabinClass: "economy",
			DepartureTime: baseTime, ArrivalTime: baseTime.Add(2 * time.Hour), Duration: 120 * time.Minute,
		},
		{
			ID: "2", Price: decimal.NewFromInt(2000), Stops: 1, Provider: "LionAir", CabinClass: "business",
			DepartureTime: baseTime.Add(1 * time.Hour), ArrivalTime: baseTime.Add(4 * time.Hour), Duration: 180 * time.Minute,
		},
		{
			ID: "3", Price: decimal.NewFromInt(1500), Stops: 0, Provider: "AirAsia", CabinClass: "economy",
			DepartureTime: baseTime.Add(2 * time.Hour), ArrivalTime: baseTime.Add(5 * time.Hour), Duration: 180 * time.Minute,
		},
	}

	t.Run("filter by exact cabin class", func(t *testing.T) {
		f := &entity.SearchFilter{CabinClass: ptrString("business")}
		// Filtering in-place, map copies to prevent altering global mock
		copied := append([]*entity.Flight(nil), flights...)
		res := cmd.Execute(copied, f)
		assert.Len(t, res, 1)
		assert.Equal(t, "2", res[0].ID)
	})

	t.Run("filter by max price", func(t *testing.T) {
		f := &entity.SearchFilter{MaxPrice: ptrDecimal(1600)}
		copied := append([]*entity.Flight(nil), flights...)
		res := cmd.Execute(copied, f)
		assert.Len(t, res, 2)
		assert.Equal(t, "1", res[0].ID)
		assert.Equal(t, "3", res[1].ID)
	})

	t.Run("filter by multiple providers", func(t *testing.T) {
		f := &entity.SearchFilter{AirlineCodes: []string{"Garuda", "LionAir"}}
		copied := append([]*entity.Flight(nil), flights...)
		res := cmd.Execute(copied, f)
		assert.Len(t, res, 2)
		assert.Equal(t, "1", res[0].ID)
		assert.Equal(t, "2", res[1].ID)
	})

	t.Run("filter by max stops", func(t *testing.T) {
		f := &entity.SearchFilter{MaxStops: ptrInt(0)}
		copied := append([]*entity.Flight(nil), flights...)
		res := cmd.Execute(copied, f)
		assert.Len(t, res, 2)
		assert.Equal(t, "1", res[0].ID)
		assert.Equal(t, "3", res[1].ID)
	})

	t.Run("no filter supplied", func(t *testing.T) {
		copied := append([]*entity.Flight(nil), flights...)
		res := cmd.Execute(copied, nil)
		assert.Len(t, res, 3)
	})
}
