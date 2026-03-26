package lionair

type LionResponse struct {
	Success bool `json:"success"`
	Data    struct {
		AvailableFlights []struct {
			ID      string `json:"id"`
			Carrier struct {
				Name string `json:"name"`
				IATA string `json:"iata"`
			} `json:"carrier"`
			Route struct {
				From struct {
					Code string `json:"code"`
				} `json:"from"`
				To struct {
					Code string `json:"code"`
				} `json:"to"`
			} `json:"route"`
			Schedule struct {
				Departure         string `json:"departure"`
				DepartureTimezone string `json:"departure_timezone"`
				Arrival           string `json:"arrival"`
				ArrivalTimezone   string `json:"arrival_timezone"`
			} `json:"schedule"`
			Pricing struct {
				Total    float64 `json:"total"`
				Currency string  `json:"currency"`
				FareType string  `json:"fare_type"`
			} `json:"pricing"`
			SeatsLeft int `json:"seats_left"`
		} `json:"available_flights"`
	} `json:"data"`
}
