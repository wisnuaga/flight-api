package domain

// NormalizeFlight normalizes the flight data by handling missing optional fields
// and ensuring proper duration calculation and timezone consistency.
func NormalizeFlight(f Flight) Flight {
	if f.CabinClass == "" {
		f.CabinClass = "economy"
	}

	if f.AvailableSeats == 0 {
		f.AvailableSeats = 1 // Default to 1 seat available minimum if returned but empty
	}

	// Always convert to UTC internally for consistent timezone handling
	if !f.DepartureTime.IsZero() {
		f.DepartureTime = f.DepartureTime.UTC()
	}
	if !f.ArrivalTime.IsZero() {
		f.ArrivalTime = f.ArrivalTime.UTC()
	}

	// Ensure duration is strictly set once during normalization
	if !f.ArrivalTime.IsZero() && !f.DepartureTime.IsZero() {
		f.Duration = f.ArrivalTime.Sub(f.DepartureTime)
	}

	// Apply basic normalizations specific to the entity instance methods
	f.Normalize()

	// Stops is an int and its zero value (0) correctly represents no-stops,
	// so it functions naturally without explicit "if missing -> 0" checks.

	return f
}
