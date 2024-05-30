package utils

import "time"

// ParseTimestamp parses a timestamp string in RFC3339 format. Returns nil if the input string is empty, otherwise returns
// the parsed time or an error.
func ParseTimestamp(timestamp string) (time.Time, error) {
	if timestamp == "" {
		return time.Time{}, nil
	}
	parsedTime, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return time.Time{}, err
	}

	return parsedTime, nil
}
