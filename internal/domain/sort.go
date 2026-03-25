package domain

type SortField string

const (
	SortByPrice     SortField = "price"
	SortByDuration  SortField = "duration"
	SortByDeparture SortField = "departure_time"
	SortByArrival   SortField = "arrival_time"
	SortByBestValue SortField = "best_value"
)

type SortOrder string

const (
	SortAsc  SortOrder = "asc"
	SortDesc SortOrder = "desc"
)

type SearchSort struct {
	Field SortField
	Order SortOrder

	// Weights for Best Value ranking (0 defaults to 1.0 internally)
	PriceWeight    float64
	DurationWeight float64
}
