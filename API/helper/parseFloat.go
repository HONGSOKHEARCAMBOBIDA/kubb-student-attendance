package helper

import (
	"fmt"
	"strconv"
	"time"
)

func ParseFloat(s string) float64 {
	if s == "" {
		return 0
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}
	return v
}

func FormatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', 2, 64)
}

func ParseLeaveDate(value string) (time.Time, error) {
	if t, err := time.Parse("2006-01-02", value); err == nil {
		return t, nil
	}

	if t, err := time.Parse(time.RFC3339, value); err == nil {
		return t, nil
	}

	return time.Time{}, fmt.Errorf("invalid date format: %s", value)
}
