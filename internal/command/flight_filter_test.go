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
func ptrTime(t time.Time) *time.Time        { return &t }

func TestFlightFilterCommand_Execute(t *testing.T) {
	cmd := command.NewFlightFilterCommand()
	baseTime := time.Date(2025, 1, 1, 10, 0, 0, 0, time.UTC)

	flights := []*entity.Flight{
		{
			ID: "1", Price: decimal.NewFromInt(1000), Stops: 0, Provider: "Garuda", CabinClass: "economy",
			Origin:      entity.Location{Airport: "CGK", Time: baseTime},
			Destination: entity.Location{Airport: "DPS", Time: baseTime.Add(2 * time.Hour)},
			Duration:    120 * time.Minute,
		},
		{
			ID: "2", Price: decimal.NewFromInt(2000), Stops: 1, Provider: "LionAir", CabinClass: "business",
			Origin:      entity.Location{Airport: "CGK", Time: baseTime.Add(1 * time.Hour)},
			Destination: entity.Location{Airport: "DPS", Time: baseTime.Add(4 * time.Hour)},
			Duration:    180 * time.Minute,
		},
		{
			ID: "3", Price: decimal.NewFromInt(1500), Stops: 0, Provider: "AirAsia", CabinClass: "economy",
			Origin:      entity.Location{Airport: "CGK", Time: baseTime.Add(2 * time.Hour)},
			Destination: entity.Location{Airport: "DPS", Time: baseTime.Add(5 * time.Hour)},
			Duration:    180 * time.Minute,
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

	t.Run("filter by departure start (UTC) — excludes earlier flights", func(t *testing.T) {
		// Only flights departing at baseTime+1h or later should pass
		f := &entity.SearchFilter{DepartureStart: ptrTime(baseTime.Add(1 * time.Hour))}
		copied := append([]*entity.Flight(nil), flights...)
		res := cmd.Execute(copied, f)
		assert.Len(t, res, 2)
		assert.Equal(t, "2", res[0].ID) // dep = baseTime+1h
		assert.Equal(t, "3", res[1].ID) // dep = baseTime+2h
	})

	t.Run("filter by arrival end (UTC) — excludes later flights", func(t *testing.T) {
		// Only flights arriving at baseTime+3h or earlier should pass
		f := &entity.SearchFilter{ArrivalEnd: ptrTime(baseTime.Add(3 * time.Hour))}
		copied := append([]*entity.Flight(nil), flights...)
		res := cmd.Execute(copied, f)
		assert.Len(t, res, 1)
		assert.Equal(t, "1", res[0].ID) // arr = baseTime+2h
	})

	t.Run("zero departure time rejected when departure start is set", func(t *testing.T) {
		// A flight with a zero Origin.Time must never pass a DepartureStart filter
		zeroTimeFlight := &entity.Flight{
			ID:          "zero",
			Price:       decimal.NewFromInt(999),
			Origin:      entity.Location{Airport: "CGK"}, // Time is zero value
			Destination: entity.Location{Airport: "DPS", Time: baseTime.Add(2 * time.Hour)},
			Duration:    120 * time.Minute,
		}
		f := &entity.SearchFilter{DepartureStart: ptrTime(baseTime)}
		copied := []*entity.Flight{zeroTimeFlight}
		res := cmd.Execute(copied, f)
		assert.Len(t, res, 0, "flight with zero departure time should be rejected")
	})
}
