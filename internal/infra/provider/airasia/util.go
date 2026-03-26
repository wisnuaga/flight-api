package airasia

import (
	"regexp"
)

func getFlightCodePrefix(code string) string {
	re := regexp.MustCompile(`^[A-Za-z]+`)
	return re.FindString(code)
}
