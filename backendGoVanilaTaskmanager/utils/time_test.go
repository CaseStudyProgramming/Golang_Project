package utils

import (
	"testing"
	"time"
)

func TestCurrentEpochMillis(t *testing.T) {
	// Test that CurrentEpochMillis returns a reasonable value
	epoch := CurrentEpochMillis()

	// Check that it's positive and not too far from current time
	if epoch <= 0 {
		t.Errorf("CurrentEpochMillis returned non-positive value: %d", epoch)
	}

	// Check that it's within a reasonable range (within last minute)
	now := time.Now().UnixMilli()
	if epoch < now-60000 || epoch > now+60000 {
		t.Errorf("CurrentEpochMillis returned value outside reasonable range: %d, expected around %d", epoch, now)
	}
}

func TestTimeToEpochMillis(t *testing.T) {
	// Test conversion of a known time
	testTime := time.Date(2024, 9, 3, 12, 0, 0, 0, time.UTC)
	epoch := TimeToEpochMillis(testTime)

	expected := testTime.UnixMilli()
	if epoch != expected {
		t.Errorf("TimeToEpochMillis(%v) = %d, expected %d", testTime, epoch, expected)
	}
}

func TestEpochMillisToTime(t *testing.T) {
	// Test conversion of known epoch milliseconds
	epoch := int64(1725364800000) // 2024-09-03 12:00:00 UTC
	result := EpochMillisToTime(epoch)

	expected := time.UnixMilli(epoch)
	if !result.Equal(expected) {
		t.Errorf("EpochMillisToTime(%d) = %v, expected %v", epoch, result, expected)
	}
}

func TestEpochMillisPtr(t *testing.T) {
	// Test conversion to pointer
	epoch := int64(1725364800000) // 2024-09-03 12:00:00 UTC
	ptr := EpochMillisPtr(epoch)

	if ptr == nil {
		t.Fatal("EpochMillisPtr returned nil")
	}
	if *ptr != epoch {
		t.Errorf("EpochMillisPtr(%d) = %d, expected %d", epoch, *ptr, epoch)
	}
}

func TestTimePtrToEpochMillis(t *testing.T) {
	// Test conversion from time pointer
	testTime := time.Date(2024, 9, 3, 12, 0, 0, 0, time.UTC)
	ptr := TimePtrToEpochMillis(&testTime)

	if ptr == nil {
		t.Fatal("TimePtrToEpochMillis returned nil")
	}
	expected := testTime.UnixMilli()
	if *ptr != expected {
		t.Errorf("TimePtrToEpochMillis(%v) = %d, expected %d", testTime, *ptr, expected)
	}

	// Test with nil pointer
	nilPtr := TimePtrToEpochMillis(nil)
	if nilPtr != nil {
		t.Errorf("TimePtrToEpochMillis(nil) should return nil, got %v", nilPtr)
	}
}

func TestEpochMillisPtrToTime(t *testing.T) {
	// Test conversion from epoch pointer
	epoch := int64(1725357600000)
	ptr := &epoch
	timePtr := EpochMillisPtrToTime(ptr)

	if timePtr == nil {
		t.Fatal("EpochMillisPtrToTime returned nil")
	}
	expected := time.UnixMilli(epoch)
	if !timePtr.Equal(expected) {
		t.Errorf("EpochMillisPtrToTime(%d) = %v, expected %v", epoch, *timePtr, expected)
	}

	// Test with nil pointer
	nilTimePtr := EpochMillisPtrToTime(nil)
	if nilTimePtr != nil {
		t.Errorf("EpochMillisPtrToTime(nil) should return nil, got %v", nilTimePtr)
	}
}

func TestConvertToTimezone(t *testing.T) {
	// Test conversion to different timezones
	epoch := int64(1725364800000) // 2024-09-03 12:00:00 UTC

	// Test UTC
	utcTime, err := ConvertToTimezone(epoch, "UTC")
	if err != nil {
		t.Errorf("ConvertToTimezone failed for UTC: %v", err)
	}
	expectedUTC := time.UnixMilli(epoch).In(time.UTC)
	if !utcTime.Equal(expectedUTC) {
		t.Errorf("ConvertToTimezone UTC: got %v, expected %v", utcTime, expectedUTC)
	}

	// Test America/New_York
	nyTime, err := ConvertToTimezone(epoch, "America/New_York")
	if err != nil {
		t.Errorf("ConvertToTimezone failed for America/New_York: %v", err)
	}
	// New York is UTC-4 or UTC-5 depending on DST
	// Just verify it's different from UTC
	loc, _ := time.LoadLocation("America/New_York")
	expectedNY := time.UnixMilli(epoch).In(loc)
	if !nyTime.Equal(expectedNY) {
		t.Errorf("ConvertToTimezone America/New_York: got %v, expected %v", nyTime, expectedNY)
	}

	// Test invalid timezone
	_, err = ConvertToTimezone(epoch, "Invalid/Timezone")
	if err == nil {
		t.Error("ConvertToTimezone should return error for invalid timezone")
	}
}

func TestConvertFromTimezone(t *testing.T) {
	// Test conversion from timezone to UTC epoch
	testTime := time.Date(2024, 9, 3, 12, 0, 0, 0, time.UTC)

	// Test UTC conversion
	epoch, err := ConvertFromTimezone(testTime, "UTC")
	if err != nil {
		t.Errorf("ConvertFromTimezone failed for UTC: %v", err)
	}
	expected := testTime.UnixMilli()
	if epoch != expected {
		t.Errorf("ConvertFromTimezone UTC: got %d, expected %d", epoch, expected)
	}

	// Test with different timezone
	loc, _ := time.LoadLocation("America/New_York")
	nyTime := time.Date(2024, 9, 3, 12, 0, 0, 0, loc)
	epoch, err = ConvertFromTimezone(nyTime, "America/New_York")
	if err != nil {
		t.Errorf("ConvertFromTimezone failed for America/New_York: %v", err)
	}
	// The epoch should be the same regardless of timezone representation
	expected = nyTime.In(time.UTC).UnixMilli()
	if epoch != expected {
		t.Errorf("ConvertFromTimezone America/New_York: got %d, expected %d", epoch, expected)
	}

	// Test invalid timezone
	_, err = ConvertFromTimezone(testTime, "Invalid/Timezone")
	if err == nil {
		t.Error("ConvertFromTimezone should return error for invalid timezone")
	}
}

func TestFormatEpochMillis(t *testing.T) {
	// Test formatting with custom format
	// Use a known epoch: 2024-09-03 12:00:00 UTC
	epoch := int64(1725364800000) // 2024-09-03 12:00:00 UTC

	formatted, err := FormatEpochMillis(epoch, "UTC", "2006-01-02 15:04:05")
	if err != nil {
		t.Errorf("FormatEpochMillis failed: %v", err)
	}
	expected := "2024-09-03 12:00:00"
	if formatted != expected {
		t.Errorf("FormatEpochMillis: got %s, expected %s", formatted, expected)
	}

	// Test with default format
	formatted, err = FormatEpochMillis(epoch, "UTC", "")
	if err != nil {
		t.Errorf("FormatEpochMillis with default format failed: %v", err)
	}
	if formatted != expected {
		t.Errorf("FormatEpochMillis with default format: got %s, expected %s", formatted, expected)
	}

	// Test with different timezone
	formatted, err = FormatEpochMillis(epoch, "America/New_York", "2006-01-02 15:04:05")
	if err != nil {
		t.Errorf("FormatEpochMillis with New York timezone failed: %v", err)
	}
	// Should be different from UTC (New York is UTC-4 or UTC-5 depending on DST)
	utcFormatted, _ := FormatEpochMillis(epoch, "UTC", "2006-01-02 15:04:05")
	if formatted == utcFormatted {
		t.Errorf("FormatEpochMillis should produce different results for different timezones. UTC: %s, NY: %s", utcFormatted, formatted)
	}

	// Test invalid timezone
	_, err = FormatEpochMillis(epoch, "Invalid/Timezone", "2006-01-02 15:04:05")
	if err == nil {
		t.Error("FormatEpochMillis should return error for invalid timezone")
	}
}

func TestGetAvailableTimezones(t *testing.T) {
	// Test that we get a list of timezones
	timezones := GetAvailableTimezones()

	if len(timezones) == 0 {
		t.Error("GetAvailableTimezones returned empty list")
	}

	// Check that UTC is in the list
	foundUTC := false
	for _, tz := range timezones {
		if tz == "UTC" {
			foundUTC = true
			break
		}
	}
	if !foundUTC {
		t.Error("GetAvailableTimezones should include UTC")
	}

	// Check that all timezones are valid
	for _, tz := range timezones {
		if !IsValidTimezone(tz) {
			t.Errorf("GetAvailableTimezones returned invalid timezone: %s", tz)
		}
	}
}

func TestIsValidTimezone(t *testing.T) {
	// Test valid timezones
	validTimezones := []string{
		"UTC",
		"America/New_York",
		"Europe/London",
		"Asia/Tokyo",
	}

	for _, tz := range validTimezones {
		if !IsValidTimezone(tz) {
			t.Errorf("IsValidTimezone should return true for valid timezone: %s", tz)
		}
	}

	// Test invalid timezones
	invalidTimezones := []string{
		"Invalid/Timezone",
		"NOT/A/TIMEZONE",
	}

	for _, tz := range invalidTimezones {
		if IsValidTimezone(tz) {
			t.Errorf("IsValidTimezone should return false for invalid timezone: %s", tz)
		}
	}
}

func TestConvertToTimezoneAware(t *testing.T) {
	// Test normal conversion
	epoch := int64(1725364800000) // 2024-09-03 12:00:00 UTC

	result, err := ConvertToTimezoneAware(epoch, "UTC")
	if err != nil {
		t.Errorf("ConvertToTimezoneAware failed: %v", err)
	}

	if result == nil {
		t.Fatal("ConvertToTimezoneAware returned nil for valid epoch")
	}

	if result.EpochMillis != epoch {
		t.Errorf("ConvertToTimezoneAware EpochMillis: got %d, expected %d", result.EpochMillis, epoch)
	}

	if result.Timezone != "UTC" {
		t.Errorf("ConvertToTimezoneAware Timezone: got %s, expected UTC", result.Timezone)
	}

	if result.ISO8601 == "" {
		t.Error("ConvertToTimezoneAware ISO8601 should not be empty")
	}

	if result.Formatted == "" {
		t.Error("ConvertToTimezoneAware Formatted should not be empty")
	}

	// Test with zero epoch (should return nil)
	result, err = ConvertToTimezoneAware(0, "UTC")
	if err != nil {
		t.Errorf("ConvertToTimezoneAware with zero epoch failed: %v", err)
	}
	if result != nil {
		t.Error("ConvertToTimezoneAware should return nil for zero epoch")
	}

	// Test with invalid timezone
	_, err = ConvertToTimezoneAware(epoch, "Invalid/Timezone")
	if err == nil {
		t.Error("ConvertToTimezoneAware should return error for invalid timezone")
	}

	// Test with different timezone
	result, err = ConvertToTimezoneAware(epoch, "America/New_York")
	if err != nil {
		t.Errorf("ConvertToTimezoneAware with New York timezone failed: %v", err)
	}
	if result.Timezone != "America/New_York" {
		t.Errorf("ConvertToTimezoneAware Timezone: got %s, expected America/New_York", result.Timezone)
	}
	// The formatted time should be different from UTC
	utcResult, _ := ConvertToTimezoneAware(epoch, "UTC")
	if result.Formatted == utcResult.Formatted {
		t.Error("ConvertToTimezoneAware should produce different formatted times for different timezones")
	}
}

func TestConvertTimezoneAwarePtr(t *testing.T) {
	// Test normal conversion
	epoch := int64(1725357600000)

	result, err := ConvertTimezoneAwarePtr(&epoch, "UTC")
	if err != nil {
		t.Errorf("ConvertTimezoneAwarePtr failed: %v", err)
	}

	if result == nil {
		t.Fatal("ConvertTimezoneAwarePtr returned nil for valid epoch pointer")
	}

	if result.EpochMillis != epoch {
		t.Errorf("ConvertTimezoneAwarePtr EpochMillis: got %d, expected %d", result.EpochMillis, epoch)
	}

	// Test with nil pointer
	result, err = ConvertTimezoneAwarePtr(nil, "UTC")
	if err != nil {
		t.Errorf("ConvertTimezoneAwarePtr with nil pointer failed: %v", err)
	}
	if result != nil {
		t.Error("ConvertTimezoneAwarePtr should return nil for nil pointer")
	}

	// Test with zero value pointer
	zeroEpoch := int64(0)
	result, err = ConvertTimezoneAwarePtr(&zeroEpoch, "UTC")
	if err != nil {
		t.Errorf("ConvertTimezoneAwarePtr with zero epoch pointer failed: %v", err)
	}
	if result != nil {
		t.Error("ConvertTimezoneAwarePtr should return nil for zero epoch pointer")
	}
}

func TestRoundTripConversion(t *testing.T) {
	// Test that converting to timezone and back preserves the epoch
	originalEpoch := int64(1725364800000) // 2024-09-03 12:00:00 UTC

	// Convert to New York timezone
	nyTime, err := ConvertToTimezone(originalEpoch, "America/New_York")
	if err != nil {
		t.Errorf("ConvertToTimezone failed: %v", err)
	}

	// Convert back to epoch
	roundTripEpoch, err := ConvertFromTimezone(nyTime, "America/New_York")
	if err != nil {
		t.Errorf("ConvertFromTimezone failed: %v", err)
	}

	if roundTripEpoch != originalEpoch {
		t.Errorf("Round trip conversion failed: got %d, expected %d", roundTripEpoch, originalEpoch)
	}
}

func TestTimezoneAwareTimestampFormatting(t *testing.T) {
	// Test that the TimezoneAwareTimestamp produces consistent formatting
	epoch := int64(1725364800000) // 2024-09-03 12:00:00 UTC

	result, err := ConvertToTimezoneAware(epoch, "UTC")
	if err != nil {
		t.Errorf("ConvertToTimezoneAware failed: %v", err)
	}

	// Verify ISO8601 format can be parsed back
	_, err = time.Parse(time.RFC3339, result.ISO8601)
	if err != nil {
		t.Errorf("TimezoneAwareTimestamp ISO8601 format is invalid: %v", err)
	}

	// Verify custom format can be parsed back
	_, err = time.Parse("2006-01-02 15:04:05", result.Formatted)
	if err != nil {
		t.Errorf("TimezoneAwareTimestamp Formatted format is invalid: %v", err)
	}
}
