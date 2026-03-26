package util

import (
	"fmt"
	"time"
)

var timeFormats = []string{
	time.RFC3339,
	"2006-01-02T15:04:05-0700",
	"2006-01-02 15:04",
	"2006-01-02 15:04:05",
	"2006-01-02T15:04:05",
	"2006-01-02",
}

// ParseTime tries to parse the given time string using multiple formats.
// It returns a properly parsed time.Time with timezone awareness, defaulting to UTC
// if the format doesn't inherently contain timezone information.
func ParseTime(value string) (time.Time, error) {
	for _, format := range timeFormats {
		t, err := time.Parse(format, value)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unable to parse time value: %s", value)
}

// ParseTimeWithOptionalTZ parses timestamps like:
// 1. "2006-01-02T15:04:05" + separate timezone string
// 2. "2006-01-02T15:04:05-0700" with offset embedded
func ParseTimeWithOptionalTZ(value string, tz string) (time.Time, error) {
	if tz != "" {
		loc, err := time.LoadLocation(tz)
		if err != nil {
			return time.Time{}, fmt.Errorf("invalid timezone: %s", tz)
		}
		t, err := time.ParseInLocation("2006-01-02T15:04:05", value, loc)
		if err != nil {
			return time.Time{}, err
		}
		return t, nil
	}

	// fallback: parse using embedded timezone or UTC
	return ParseTime(value) // reuse existing util.ParseTime
}
