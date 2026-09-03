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
