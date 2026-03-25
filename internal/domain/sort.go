package domain

type SortField string

const (
	SortByPrice     SortField = "price"
	SortByDuration  SortField = "duration"
	SortByDeparture SortField = "departure_time"
	SortByArrival   SortField = "arrival_time"
)

type SortOrder string

const (
	SortAsc  SortOrder = "asc"
	SortDesc SortOrder = "desc"
)

type SearchSort struct {
	Field SortField
	Order SortOrder
}
