package dto

import "github.com/shopspring/decimal"

type SearchRequest struct {
	Origin        string  `json:"origin"`
	Destination   string  `json:"destination"`
	DepartureDate string  `json:"departure_date"`
	ReturnDate    *string `json:"return_date"`
	Passengers    int     `json:"passengers"`
	CabinClass    string  `json:"cabin_class"`

	// Filters
	MinPrice       *decimal.Decimal `json:"min_price"`
	MaxPrice       *decimal.Decimal `json:"max_price"`
	MaxStops       *int             `json:"max_stops"`
	DepartureStart *string          `json:"departure_start"`
	DepartureEnd   *string          `json:"departure_end"`
	ArrivalStart   *string          `json:"arrival_start"`
	ArrivalEnd     *string          `json:"arrival_end"`
	Airlines       []string         `json:"airlines"`     // Airline names: "Garuda Indonesia", "Lion Air", "AirAsia", "Batik Air"
	MaxDuration    *int             `json:"max_duration"` // Minutes

	// Sort
	SortBy    string `json:"sort_by"`    // price | duration | departure_time | arrival_time | best_value  (default: price)
	SortOrder string `json:"sort_order"` // asc | desc  (default: asc)
}
