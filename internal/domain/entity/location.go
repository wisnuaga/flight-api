package entity

import (
	"strings"
	"time"
)

var airportCityMap = map[string]string{
	"CGK": "Jakarta",
	"DPS": "Denpasar",
	"SUB": "Surabaya",
	"JKT": "Jakarta",
	// Add more airport codes as needed for different providers
}

// Location represents a flight location (airport) with timezone-aware scheduling
type Location struct {
	Airport  string
	City     string
	Time     time.Time      // Stored in UTC internally for filtering/sorting consistency
	Timezone *time.Location // Original timezone from provider for output formatting
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
