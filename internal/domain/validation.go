package domain

// IsValidFlight validates the flight properties ensuring it meets business rules.
func IsValidFlight(f Flight) bool {
	if f.Origin == "" || f.Destination == "" {
		return false
	}

	if f.Price <= 0 {
		return false
	}

	if f.DepartureTime.IsZero() || f.ArrivalTime.IsZero() {
		return false
	}

	if !f.ArrivalTime.After(f.DepartureTime) {
		return false
	}

	if f.Duration <= 0 {
		return false
	}

	return true
}
