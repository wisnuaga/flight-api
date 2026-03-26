package entity

// AirlineName represents the valid airline names (providers) in the system.
type AirlineName string

const (
	AirlineGaruda   AirlineName = "Garuda Indonesia"
	AirlineLionAir  AirlineName = "Lion Air"
	AirlineAirAsia  AirlineName = "AirAsia"
	AirlineBatikAir AirlineName = "Batik Air"
)

// String returns the string representation of the airline name.
func (a AirlineName) String() string {
	return string(a)
}

// IsValid checks if the airline name is valid.
func (a AirlineName) IsValid() bool {
	switch a {
	case AirlineGaruda, AirlineLionAir, AirlineAirAsia, AirlineBatikAir:
		return true
	default:
		return false
	}
}

// AirlineNameFromString converts a string to AirlineName.
// Returns the matched AirlineName or empty string if not found.
func AirlineNameFromString(s string) AirlineName {
	switch s {
	case "Garuda Indonesia", "Garuda", "GA":
		return AirlineGaruda
	case "Lion Air", "LionAir", "JT":
		return AirlineLionAir
	case "AirAsia", "QZ":
		return AirlineAirAsia
	case "Batik Air", "BatikAir", "ID":
		return AirlineBatikAir
	default:
		return ""
	}
}

// ValidAirlineNames returns all valid airline names.
func ValidAirlineNames() []AirlineName {
	return []AirlineName{
		AirlineGaruda,
		AirlineLionAir,
		AirlineAirAsia,
		AirlineBatikAir,
	}
}
