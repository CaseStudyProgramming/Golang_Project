package utils

import (
	"time"
)

// CurrentEpochMillis returns the current time in epoch milliseconds
func CurrentEpochMillis() int64 {
	return time.Now().UnixMilli()
}

// TimeToEpochMillis converts a time.Time to epoch milliseconds
func TimeToEpochMillis(t time.Time) int64 {
	return t.UnixMilli()
}

// EpochMillisToTime converts epoch milliseconds to time.Time
func EpochMillisToTime(epochMillis int64) time.Time {
	return time.UnixMilli(epochMillis)
}

// EpochMillisPtr converts epoch milliseconds to a pointer
func EpochMillisPtr(epochMillis int64) *int64 {
	return &epochMillis
}

// TimePtrToEpochMillis converts a time.Time pointer to epoch milliseconds pointer
func TimePtrToEpochMillis(t *time.Time) *int64 {
	if t == nil {
		return nil
	}
	epochMillis := t.UnixMilli()
	return &epochMillis
}

// EpochMillisPtrToTime converts epoch milliseconds pointer to time.Time pointer
func EpochMillisPtrToTime(epochMillis *int64) *time.Time {
	if epochMillis == nil {
		return nil
	}
	t := time.UnixMilli(*epochMillis)
	return &t
}

// ConvertToTimezone converts epoch milliseconds to a specific timezone
func ConvertToTimezone(epochMillis int64, timezone string) (time.Time, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return time.Time{}, err
	}
	return time.UnixMilli(epochMillis).In(loc), nil
}

// ConvertFromTimezone converts a time from a specific timezone to UTC epoch milliseconds
func ConvertFromTimezone(t time.Time, timezone string) (int64, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return 0, err
	}
	// Parse the time as if it's in the specified timezone
	timeInLoc := time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond(), loc)
	return timeInLoc.UnixMilli(), nil
}

// FormatEpochMillis formats epoch milliseconds to a human-readable string in the specified timezone
func FormatEpochMillis(epochMillis int64, timezone string, format string) (string, error) {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return "", err
	}
	if format == "" {
		format = "2006-01-02 15:04:05"
	}
	return time.UnixMilli(epochMillis).In(loc).Format(format), nil
}

// GetAvailableTimezones returns a list of common timezones
func GetAvailableTimezones() []string {
	return []string{
		"UTC",
		"America/New_York",
		"America/Chicago",
		"America/Denver",
		"America/Los_Angeles",
		"Europe/London",
		"Europe/Paris",
		"Europe/Berlin",
		"Asia/Tokyo",
		"Asia/Shanghai",
		"Asia/Singapore",
		"Australia/Sydney",
		"Pacific/Auckland",
	}
}

// IsValidTimezone checks if a timezone string is valid
func IsValidTimezone(timezone string) bool {
	_, err := time.LoadLocation(timezone)
	return err == nil
}
