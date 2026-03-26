package garuda

type GarudaSearchResponse struct {
	Status  string         `json:"status"`
	Flights []GarudaFlight `json:"flights"`
}

type GarudaFlight struct {
	FlightID        string          `json:"flight_id"`
	Airline         string          `json:"airline"`
	AirlineCode     string          `json:"airline_code"`
	Departure       *GarudaAirport  `json:"departure"`
	Arrival         *GarudaAirport  `json:"arrival"`
	DurationMinutes int             `json:"duration_minutes"`
	Stops           int             `json:"stops"`
	Aircraft        string          `json:"aircraft"`
	Price           *GarudaPrice    `json:"price"`
	AvailableSeats  int             `json:"available_seats"`
	FareClass       string          `json:"fare_class"`
	Baggage         *GarudaBaggage  `json:"baggage"`
	Amenities       []string        `json:"amenities"`
	Segments        []GarudaSegment `json:"segments"`
}

type GarudaAirport struct {
	Airport  string `json:"airport"`
	City     string `json:"city"`
	Time     string `json:"time"`
	Terminal string `json:"terminal"`
}

type GarudaPrice struct {
	Amount   int64  `json:"amount"`
	Currency string `json:"currency"`
}

type GarudaBaggage struct {
	CarryOn int `json:"carry_on"`
	Checked int `json:"checked"`
}

type GarudaSegment struct {
	FlightNumber    string              `json:"flight_number"`
	Departure       *GarudaSegmentPoint `json:"departure"`
	Arrival         *GarudaSegmentPoint `json:"arrival"`
	DurationMinutes int                 `json:"duration_minutes"`
	LayoverMinutes  int                 `json:"layover_minutes,omitempty"`
}

type GarudaSegmentPoint struct {
	Airport string `json:"airport"`
	Time    string `json:"time"`
}
