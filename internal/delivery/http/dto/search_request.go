package dto

type SearchRequest struct {
	Origin        string  `json:"origin"`
	Destination   string  `json:"destination"`
	DepartureDate string  `json:"departure_date"`
	ReturnDate    *string `json:"return_date"`
	Passengers    int     `json:"passengers"`
	CabinClass    string  `json:"cabin_class"`

	// Filters
	MinPrice       *float64 `json:"min_price"`
	MaxPrice       *float64 `json:"max_price"`
	MaxStops       *int     `json:"max_stops"`
	DepartureStart *string  `json:"departure_start"`
	DepartureEnd   *string  `json:"departure_end"`
	ArrivalStart   *string  `json:"arrival_start"`
	ArrivalEnd     *string  `json:"arrival_end"`
	AirlineCodes   []string `json:"airline_codes"`
	MaxDuration    *int     `json:"max_duration"` // Minutes
}
