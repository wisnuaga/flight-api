package garuda

type SearchResponse struct {
	Status  string         `json:"status"`
	Flights []SearchFlight `json:"flights"`
}

type SearchFlight struct {
	FlightID        string          `json:"flight_id"`
	Airline         string          `json:"airline"`
	AirlineCode     string          `json:"airline_code"`
	Departure       *Searchirport   `json:"departure"`
	Arrival         *Searchirport   `json:"arrival"`
	DurationMinutes int             `json:"duration_minutes"`
	Stops           int             `json:"stops"`
	Aircraft        string          `json:"aircraft"`
	Price           *SearchPrice    `json:"price"`
	AvailableSeats  int             `json:"available_seats"`
	FareClass       string          `json:"fare_class"`
	Baggage         *SearchBaggage  `json:"baggage"`
	Amenities       []string        `json:"amenities"`
	Segments        []SearchSegment `json:"segments"`
}

type Searchirport struct {
	Airport  string `json:"airport"`
	City     string `json:"city"`
	Time     string `json:"time"`
	Terminal string `json:"terminal"`
}

type SearchPrice struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

type SearchBaggage struct {
	CarryOn int `json:"carry_on"`
	Checked int `json:"checked"`
}

type SearchSegment struct {
	FlightNumber    string              `json:"flight_number"`
	Departure       *SearchSegmentPoint `json:"departure"`
	Arrival         *SearchSegmentPoint `json:"arrival"`
	DurationMinutes int                 `json:"duration_minutes"`
	LayoverMinutes  *int                `json:"layover_minutes,omitempty"`
}

type SearchSegmentPoint struct {
	Airport string `json:"airport"`
	Time    string `json:"time"`
}
