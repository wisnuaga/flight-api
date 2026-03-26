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
// if the format does not inherently contain timezone information.
func ParseTime(value string) (time.Time, error) {
	for _, format := range timeFormats {
		t, err := time.Parse(format, value)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unable to parse time value: %s", value)
}

// ParseTimeWithOptionalTZ parses timestamps and converts them to UTC.
// Use when the provider supplies a separate IANA timezone name alongside a naive time string.
func ParseTimeWithOptionalTZ(value string, tz string) (time.Time, error) {
	if tz != "" {
		loc, err := time.LoadLocation(tz)
		if err != nil {
			return ParseTime(value)
		}
		t, err := time.ParseInLocation("2006-01-02T15:04:05", value, loc)
		if err != nil {
			for _, format := range timeFormats {
				t2, err2 := time.ParseInLocation(format, value, loc)
				if err2 == nil {
					return t2.UTC(), nil
				}
			}
			return time.Time{}, err
		}
		return t.UTC(), nil
	}
	return ParseTime(value)
}

// ParseTimeWithTZInfo parses a time string and also returns the original timezone location.
// Use this when the provider supplies a separate IANA timezone name (e.g. Lion Air).
//
// Returns: (UTC time, original timezone location, error)
func ParseTimeWithTZInfo(value string, tzStr string) (time.Time, *time.Location, error) {
	var originalTz *time.Location = time.UTC

	if tzStr != "" {
		loc, err := time.LoadLocation(tzStr)
		if err != nil {
			parsedTime, err := ParseTime(value)
			return parsedTime, time.UTC, err
		}
		originalTz = loc

		t, err := time.ParseInLocation("2006-01-02T15:04:05", value, loc)
		if err != nil {
			for _, format := range timeFormats {
				t2, err2 := time.ParseInLocation(format, value, loc)
				if err2 == nil {
					return t2.UTC(), originalTz, nil
				}
			}
			return time.Time{}, originalTz, err
		}
		return t.UTC(), originalTz, nil
	}

	parsedTime, err := ParseTime(value)
	return parsedTime, time.UTC, err
}

// ParseTimeFromString parses a time string that contains an embedded UTC offset
// (e.g. RFC3339 "2006-01-02T15:04:05+07:00" or compact "2006-01-02T15:04:05+0700").
// It returns the instant in UTC and the fixed-offset *time.Location extracted from
// the string itself — no separate IANA timezone name is needed.
//
// Use this for providers (Garuda, AirAsia, Batik Air) that embed the offset in the
// time value but do NOT supply a named timezone field.
//
// Returns: (UTC time, fixed-offset location from string, error)
func ParseTimeFromString(value string) (time.Time, *time.Location, error) {
	t, err := ParseTime(value)
	if err != nil {
		return time.Time{}, time.UTC, err
	}
	// t.Location() holds the fixed-offset location parsed from the embedded offset.
	// If the format had no offset (naive time), Location() returns time.UTC.
	return t.UTC(), t.Location(), nil
}
