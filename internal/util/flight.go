package util

import (
	"regexp"
	"strings"
)

func GetFlightCodePrefix(code string) string {
	re := regexp.MustCompile(`^[A-Za-z]+`)
	return re.FindString(code)
}

func NormalizeAirlineName(name string) string {
	return strings.ReplaceAll(name, " ", "")
}
