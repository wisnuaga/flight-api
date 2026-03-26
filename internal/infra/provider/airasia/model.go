package airasia

type AirAsiaResponse struct {
	Status  string `json:"status"`
	Flights []struct {
		FlightCode  string  `json:"flight_code"`
		Airline     string  `json:"airline"`
		FromAirport string  `json:"from_airport"`
		ToAirport   string  `json:"to_airport"`
		DepartTime  string  `json:"depart_time"`
		ArriveTime  string  `json:"arrive_time"`
		PriceIDR    float64 `json:"price_idr"`
		Seats       int     `json:"seats"`
		CabinClass  string  `json:"cabin_class"`
	} `json:"flights"`
}
