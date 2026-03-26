package batikair

type BatikResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Results []struct {
		FlightNumber      string `json:"flightNumber"`
		AirlineName       string `json:"airlineName"`
		AirlineIATA       string `json:"airlineIATA"`
		Origin            string `json:"origin"`
		Destination       string `json:"destination"`
		DepartureDateTime string `json:"departureDateTime"`
		ArrivalDateTime   string `json:"arrivalDateTime"`
		Fare              struct {
			BasePrice    float64 `json:"basePrice"`
			Taxes        float64 `json:"taxes"`
			TotalPrice   float64 `json:"totalPrice"`
			CurrencyCode string  `json:"currencyCode"`
			Class        string  `json:"class"`
		} `json:"fare"`
		SeatsAvailable int `json:"seatsAvailable"`
	} `json:"results"`
}
