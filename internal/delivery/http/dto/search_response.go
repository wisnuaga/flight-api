package dto

type SearchResponse struct {
	SearchCriteria SearchCriteria `json:"search_criteria"`
	Metadata       Metadata       `json:"metadata"`
	Flights        []Flight       `json:"flights"`
}

// RoundTripResponse represents a round-trip search response
type RoundTripResponse struct {
	SearchCriteria SearchCriteria       `json:"search_criteria"`
	Metadata       Metadata             `json:"metadata"`
	Itineraries    []RoundTripItinerary `json:"round_trip_itineraries"`
}

// RoundTripItinerary represents a single round-trip option (outbound + return)
type RoundTripItinerary struct {
	OutboundFlight Flight   `json:"outbound_flight"`
	ReturnFlight   Flight   `json:"return_flight"`
	TotalPrice     Price    `json:"total_price"`
	TotalDuration  Duration `json:"total_duration"`
}

type SearchCriteria struct {
	Origin        string  `json:"origin"`
	Destination   string  `json:"destination"`
	DepartureDate string  `json:"departure_date"`
	ReturnDate    *string `json:"return_date,omitempty"`
	Passengers    int     `json:"passengers"`
	CabinClass    string  `json:"cabin_class"`
}

type Metadata struct {
	TotalResults       int  `json:"total_results"`
	ProvidersQueried   int  `json:"providers_queried"`
	ProvidersSucceeded int  `json:"providers_succeeded"`
	ProvidersFailed    int  `json:"providers_failed"`
	SearchTimeMs       int  `json:"search_time_ms"`
	CacheHit           bool `json:"cache_hit"`
}

type Flight struct {
	ID             string   `json:"id"`
	Provider       string   `json:"provider"`
	Airline        Airline  `json:"airline"`
	FlightNumber   string   `json:"flight_number"`
	Departure      Location `json:"departure"`
	Arrival        Location `json:"arrival"`
	Duration       Duration `json:"duration"`
	Stops          int      `json:"stops"`
	Price          Price    `json:"price"`
	AvailableSeats int      `json:"available_seats"`
	CabinClass     string   `json:"cabin_class"`
	Aircraft       *string  `json:"aircraft"`
	Amenities      []string `json:"amenities"`
	Baggage        Baggage  `json:"baggage"`
}

type Airline struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type Location struct {
	Airport   string `json:"airport"`
	City      string `json:"city"`
	Datetime  string `json:"datetime"`
	Timestamp int64  `json:"timestamp"`
}

type Duration struct {
	TotalMinutes int    `json:"total_minutes"`
	Formatted    string `json:"formatted"`
}

type Price struct {
	Amount    int64  `json:"amount"`
	Currency  string `json:"currency"`
	Formatted string `json:"formatted"`
}

type Baggage struct {
	CarryOn string `json:"carry_on"`
	Checked string `json:"checked"`
}
