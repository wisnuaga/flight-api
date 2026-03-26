package entity

import "strings"

var airportCityMap = map[string]string{
	"CGK": "Jakarta",
	"DPS": "Denpasar",
	"SUB": "Surabaya",
	"JKT": "Jakarta",
	// Add more airport codes as needed for different providers
}

type Location struct {
	Airport string
	City    string
}

// GetCity returns city name from airport code
// Uses the airport mapping, fallback to code if not found
func (l *Location) GetCity() string {
	if city, ok := airportCityMap[strings.ToUpper(l.Airport)]; ok {
		return city
	}
	return l.Airport
}

// Normalize normalizes the location by setting city from airport code
func (l *Location) Normalize() {
	l.Airport = strings.ToUpper(l.Airport)
	l.City = l.GetCity()
}
