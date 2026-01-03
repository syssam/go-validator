package validator

import (
	"errors"
	"regexp"
	"strconv"
	"time"
)

// =============================================================================
// Relative Date Parser - Supports: today, tomorrow, yesterday, now, today+7d, today-1m, today-18y
// =============================================================================

var dateModifierRegex = regexp.MustCompile(`^(today|tomorrow|yesterday|now)([+-]\d+[dmyDMY])?$`)

// parseRelativeDate parses keywords like "today", "tomorrow-7d", "today+1m", "today-18y"
// Returns the parsed time and true if successful, or zero time and false if not a relative date
func parseRelativeDate(input string) (time.Time, bool) {
	matches := dateModifierRegex.FindStringSubmatch(input)
	if matches == nil {
		return time.Time{}, false
	}

	base := matches[1]
	modifier := matches[2]

	// Get base time
	var t time.Time
	now := time.Now()

	switch base {
	case "today":
		t = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	case "tomorrow":
		t = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, 1)
	case "yesterday":
		t = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, -1)
	case "now":
		t = now
	default:
		return time.Time{}, false
	}

	// Apply modifier if present (e.g., +7d, -1m, -18y)
	if modifier != "" {
		sign := 1
		if modifier[0] == '-' {
			sign = -1
		}

		numStr := modifier[1 : len(modifier)-1]
		unit := modifier[len(modifier)-1]

		num, err := strconv.Atoi(numStr)
		if err != nil {
			return time.Time{}, false
		}
		num *= sign

		switch unit {
		case 'd', 'D':
			t = t.AddDate(0, 0, num)
		case 'm', 'M':
			t = t.AddDate(0, num, 0)
		case 'y', 'Y':
			t = t.AddDate(num, 0, 0)
		}
	}

	return t, true
}

// parseDateParam parses a date parameter, trying relative date first, then static date formats
func parseDateParam(param string) (time.Time, error) {
	// Try relative date first
	if t, ok := parseRelativeDate(param); ok {
		return t, nil
	}

	// Try standard date formats using local timezone (consistent with parseDate in validator.go)
	formats := []string{
		"2006-01-02",
		"2006-01-02 15:04:05",
		"2006/01/02",
		"02-01-2006",
		"01/02/2006",
	}

	for _, format := range formats {
		// Use ParseInLocation for formats without timezone to use local time
		if t, err := time.ParseInLocation(format, param, time.Local); err == nil {
			return t, nil
		}
	}

	// Try RFC3339 which has timezone info
	if t, err := time.Parse(time.RFC3339, param); err == nil {
		return t, nil
	}

	return time.Time{}, errors.New("invalid date format: " + param)
}

// =============================================================================
// Time Helper Functions - For relative date calculations
// =============================================================================

// Today returns today's date at midnight (00:00:00) in local timezone
func Today() time.Time {
	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
}

// TodayUTC returns today's date at midnight (00:00:00) in UTC
func TodayUTC() time.Time {
	now := time.Now().UTC()
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
}

// Tomorrow returns tomorrow's date at midnight
func Tomorrow() time.Time {
	return Today().AddDate(0, 0, 1)
}

// Yesterday returns yesterday's date at midnight
func Yesterday() time.Time {
	return Today().AddDate(0, 0, -1)
}

// Now returns the current time
func Now() time.Time {
	return time.Now()
}

// =============================================================================
// Type-Safe Date Validation Functions (Google-style)
// =============================================================================

// DateValidator represents a date validation function
type DateValidator func(t time.Time) bool

// DateRule provides fluent date validation building
type DateRule struct {
	validators []DateValidator
	fieldName  string
}

// Date creates a new DateRule for fluent validation
func Date(fieldName ...string) *DateRule {
	name := ""
	if len(fieldName) > 0 {
		name = fieldName[0]
	}
	return &DateRule{fieldName: name}
}

// Before adds a "before" validation
func (r *DateRule) Before(t time.Time) *DateRule {
	r.validators = append(r.validators, func(v time.Time) bool {
		return v.Before(t)
	})
	return r
}

// BeforeOrEqual adds a "before or equal" validation
func (r *DateRule) BeforeOrEqual(t time.Time) *DateRule {
	r.validators = append(r.validators, func(v time.Time) bool {
		return v.Before(t) || v.Equal(t)
	})
	return r
}

// After adds an "after" validation
func (r *DateRule) After(t time.Time) *DateRule {
	r.validators = append(r.validators, func(v time.Time) bool {
		return v.After(t)
	})
	return r
}

// AfterOrEqual adds an "after or equal" validation
func (r *DateRule) AfterOrEqual(t time.Time) *DateRule {
	r.validators = append(r.validators, func(v time.Time) bool {
		return v.After(t) || v.Equal(t)
	})
	return r
}

// Between adds a "between" validation (exclusive)
func (r *DateRule) Between(start, end time.Time) *DateRule {
	r.validators = append(r.validators, func(v time.Time) bool {
		return v.After(start) && v.Before(end)
	})
	return r
}

// BetweenOrEqual adds a "between" validation (inclusive)
func (r *DateRule) BetweenOrEqual(start, end time.Time) *DateRule {
	r.validators = append(r.validators, func(v time.Time) bool {
		return (v.After(start) || v.Equal(start)) && (v.Before(end) || v.Equal(end))
	})
	return r
}

// Validate checks if the given time passes all validations
func (r *DateRule) Validate(t time.Time) bool {
	for _, v := range r.validators {
		if !v(t) {
			return false
		}
	}
	return true
}

// ValidateWithError checks if the given time passes all validations and returns an error if not
func (r *DateRule) ValidateWithError(t time.Time) error {
	if !r.Validate(t) {
		if r.fieldName != "" {
			return errors.New(r.fieldName + ": date validation failed")
		}
		return errors.New("date validation failed")
	}
	return nil
}

// =============================================================================
// Standalone Type-Safe Validation Functions
// =============================================================================

// DateAfter checks if the date is after the given time
func DateAfter(value, after time.Time) bool {
	return value.After(after)
}

// DateAfterOrEqual checks if the date is after or equal to the given time
func DateAfterOrEqual(value, after time.Time) bool {
	return value.After(after) || value.Equal(after)
}

// DateBefore checks if the date is before the given time
func DateBefore(value, before time.Time) bool {
	return value.Before(before)
}

// DateBeforeOrEqual checks if the date is before or equal to the given time
func DateBeforeOrEqual(value, before time.Time) bool {
	return value.Before(before) || value.Equal(before)
}

// DateBetween checks if the date is between start and end (exclusive)
func DateBetween(value, start, end time.Time) bool {
	return value.After(start) && value.Before(end)
}

// DateBetweenOrEqual checks if the date is between start and end (inclusive)
func DateBetweenOrEqual(value, start, end time.Time) bool {
	return (value.After(start) || value.Equal(start)) && (value.Before(end) || value.Equal(end))
}

// IsToday checks if the date is today
func IsToday(value time.Time) bool {
	today := Today()
	return value.Year() == today.Year() &&
		value.Month() == today.Month() &&
		value.Day() == today.Day()
}

// IsFuture checks if the date is in the future
func IsFuture(value time.Time) bool {
	return value.After(time.Now())
}

// IsPast checks if the date is in the past
func IsPast(value time.Time) bool {
	return value.Before(time.Now())
}

// =============================================================================
// Time Extension Methods - For fluent date arithmetic
// =============================================================================

// TimeHelper wraps time.Time with helper methods
type TimeHelper struct {
	time.Time
}

// T wraps a time.Time for fluent operations
func T(t time.Time) TimeHelper {
	return TimeHelper{t}
}

// SubDays subtracts days from the time
func (t TimeHelper) SubDays(days int) time.Time {
	return t.AddDate(0, 0, -days)
}

// AddDays adds days to the time
func (t TimeHelper) AddDays(days int) time.Time {
	return t.AddDate(0, 0, days)
}

// SubMonths subtracts months from the time
func (t TimeHelper) SubMonths(months int) time.Time {
	return t.AddDate(0, -months, 0)
}

// AddMonths adds months to the time
func (t TimeHelper) AddMonths(months int) time.Time {
	return t.AddDate(0, months, 0)
}

// SubYears subtracts years from the time
func (t TimeHelper) SubYears(years int) time.Time {
	return t.AddDate(-years, 0, 0)
}

// AddYears adds years to the time
func (t TimeHelper) AddYears(years int) time.Time {
	return t.AddDate(years, 0, 0)
}
