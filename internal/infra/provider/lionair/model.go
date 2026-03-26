package lionair

type SearchResponse struct {
	Success bool `json:"success"`
	Data    struct {
		AvailableFlights []*SearchFlight `json:"available_flights"`
	} `json:"data"`
}

type SearchFlight struct {
	ID        string                 `json:"id"`
	Carrier   *SearchFlightCarrier   `json:"carrier"`
	Route     *SearchFlightRoute     `json:"route"`
	Schedule  *SearchFlightSchedule  `json:"schedule"`
	Pricing   *SearchFlightPricing   `json:"pricing"`
	SeatsLeft int                    `json:"seats_left"`
	Layovers  []*SearchFlightLayover `json:"layovers"`
}

type SearchFlightCarrier struct {
	Name string `json:"name"`
	IATA string `json:"iata"`
}

type SearchFlightRoute struct {
	From *SearchFlightRouteLocation `json:"from"`
	To   *SearchFlightRouteLocation `json:"to"`
}

type SearchFlightSchedule struct {
	Departure         string `json:"departure"`
	DepartureTimezone string `json:"departure_timezone"`
	Arrival           string `json:"arrival"`
	ArrivalTimezone   string `json:"arrival_timezone"`
}

type SearchFlightRouteLocation struct {
	Code string `json:"code"`
	Name string `json:"name"`
	City string `json:"city"`
}

type SearchFlightPricing struct {
	Total    float64 `json:"total"`
	Currency string  `json:"currency"`
	FareType string  `json:"fare_type"`
}

type SearchFlightLayover struct {
	Airport         string `json:"airport"`
	DurationMinutes int    `json:"duration_minutes"`
}
