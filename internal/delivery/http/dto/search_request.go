package dto

type SearchRequest struct {
	Origin        string  `json:"origin"`
	Destination   string  `json:"destination"`
	DepartureDate string  `json:"departure_date"`
	ReturnDate    *string `json:"return_date"`
	Passengers    int     `json:"passengers"`
	CabinClass    string  `json:"cabin_class"`
}
