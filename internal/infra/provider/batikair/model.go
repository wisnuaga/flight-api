package batikair

type SearchResponse struct {
	Code    int             `json:"code"`
	Message string          `json:"message"`
	Results []*SearchResult `json:"results"`
}

type SearchResult struct {
	FlightNumber      string                    `json:"flightNumber"`
	AirlineName       string                    `json:"airlineName"`
	AirlineIATA       string                    `json:"airlineIATA"`
	Origin            string                    `json:"origin"`
	Destination       string                    `json:"destination"`
	DepartureDateTime string                    `json:"departureDateTime"`
	ArrivalDateTime   string                    `json:"arrivalDateTime"`
	Fare              *SearchResultFare         `json:"fare"`
	SeatsAvailable    int                       `json:"seatsAvailable"`
	Connections       []*SearchResultConnection `json:"connections"`
}

type SearchResultFare struct {
	BasePrice    float64 `json:"basePrice"`
	Taxes        float64 `json:"taxes"`
	TotalPrice   float64 `json:"totalPrice"`
	CurrencyCode string  `json:"currencyCode"`
	Class        string  `json:"class"`
}

type SearchResultConnection struct {
	StopAirport  string `json:"stopAirport"`
	StopDuration string `json:"stopDuration"`
}
