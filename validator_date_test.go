package validator

import (
	"testing"
	"time"
)

// =============================================================================
// Tests for parseRelativeDate
// =============================================================================

func TestParseRelativeDate(t *testing.T) {
	now := time.Now()
	todayMidnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	tests := []struct {
		name     string
		input    string
		wantOk   bool
		validate func(t time.Time) bool
	}{
		{
			name:   "today",
			input:  "today",
			wantOk: true,
			validate: func(t time.Time) bool {
				return t.Equal(todayMidnight)
			},
		},
		{
			name:   "tomorrow",
			input:  "tomorrow",
			wantOk: true,
			validate: func(t time.Time) bool {
				return t.Equal(todayMidnight.AddDate(0, 0, 1))
			},
		},
		{
			name:   "yesterday",
			input:  "yesterday",
			wantOk: true,
			validate: func(t time.Time) bool {
				return t.Equal(todayMidnight.AddDate(0, 0, -1))
			},
		},
		{
			name:   "today+7d",
			input:  "today+7d",
			wantOk: true,
			validate: func(t time.Time) bool {
				return t.Equal(todayMidnight.AddDate(0, 0, 7))
			},
		},
		{
			name:   "today-7d",
			input:  "today-7d",
			wantOk: true,
			validate: func(t time.Time) bool {
				return t.Equal(todayMidnight.AddDate(0, 0, -7))
			},
		},
		{
			name:   "today+1m",
			input:  "today+1m",
			wantOk: true,
			validate: func(t time.Time) bool {
				return t.Equal(todayMidnight.AddDate(0, 1, 0))
			},
		},
		{
			name:   "today-18y",
			input:  "today-18y",
			wantOk: true,
			validate: func(t time.Time) bool {
				return t.Equal(todayMidnight.AddDate(-18, 0, 0))
			},
		},
		{
			name:   "tomorrow+30d",
			input:  "tomorrow+30d",
			wantOk: true,
			validate: func(t time.Time) bool {
				return t.Equal(todayMidnight.AddDate(0, 0, 31))
			},
		},
		{
			name:   "invalid",
			input:  "invalid",
			wantOk: false,
		},
		{
			name:   "empty",
			input:  "",
			wantOk: false,
		},
		{
			name:   "static date not parsed as relative",
			input:  "2024-01-01",
			wantOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := parseRelativeDate(tt.input)
			if ok != tt.wantOk {
				t.Errorf("parseRelativeDate(%q) ok = %v, want %v", tt.input, ok, tt.wantOk)
				return
			}
			if ok && tt.validate != nil && !tt.validate(got) {
				t.Errorf("parseRelativeDate(%q) = %v, validation failed", tt.input, got)
			}
		})
	}
}

// =============================================================================
// Tests for parseDateParam (combined relative + static)
// =============================================================================

func TestParseDateParam(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"today", "today", false},
		{"tomorrow", "tomorrow", false},
		{"yesterday", "yesterday", false},
		{"today+7d", "today+7d", false},
		{"today-18y", "today-18y", false},
		{"static date YYYY-MM-DD", "2024-01-15", false},
		{"static date with time", "2024-01-15 10:30:00", false},
		{"static date YYYY/MM/DD", "2024/01/15", false},
		{"RFC3339", "2024-01-15T10:30:00Z", false},
		{"invalid", "not-a-date", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseDateParam(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDateParam(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

// =============================================================================
// Tests for struct tag validation with relative dates
// =============================================================================

func TestValidateStructWithRelativeDates(t *testing.T) {
	type BookingForm struct {
		CheckIn  time.Time `valid:"after=today"`
		CheckOut time.Time `valid:"after=tomorrow"`
	}

	type BirthDateForm struct {
		BirthDate time.Time `valid:"before=today-18y"`
	}

	type EventForm struct {
		EventDate time.Time `valid:"afterOrEqual=today,beforeOrEqual=today+90d"`
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	t.Run("BookingForm valid", func(t *testing.T) {
		form := BookingForm{
			CheckIn:  today.AddDate(0, 0, 1),  // tomorrow
			CheckOut: today.AddDate(0, 0, 3),  // 3 days from now
		}
		err := ValidateStruct(&form)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
	})

	t.Run("BookingForm invalid - CheckIn not after today", func(t *testing.T) {
		form := BookingForm{
			CheckIn:  today.AddDate(0, 0, -1), // yesterday
			CheckOut: today.AddDate(0, 0, 3),
		}
		err := ValidateStruct(&form)
		if err == nil {
			t.Error("Expected error for CheckIn before today")
		}
	})

	t.Run("BirthDateForm valid - over 18", func(t *testing.T) {
		form := BirthDateForm{
			BirthDate: today.AddDate(-20, 0, 0), // 20 years ago
		}
		err := ValidateStruct(&form)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
	})

	t.Run("BirthDateForm invalid - under 18", func(t *testing.T) {
		form := BirthDateForm{
			BirthDate: today.AddDate(-16, 0, 0), // 16 years ago
		}
		err := ValidateStruct(&form)
		if err == nil {
			t.Error("Expected error for BirthDate under 18")
		}
	})

	t.Run("EventForm valid - within 90 days", func(t *testing.T) {
		form := EventForm{
			EventDate: today.AddDate(0, 0, 30), // 30 days from now
		}
		err := ValidateStruct(&form)
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
	})

	t.Run("EventForm invalid - too far in future", func(t *testing.T) {
		form := EventForm{
			EventDate: today.AddDate(0, 0, 100), // 100 days from now
		}
		err := ValidateStruct(&form)
		if err == nil {
			t.Error("Expected error for EventDate too far in future")
		}
	})
}

// =============================================================================
// Tests for type-safe validation functions (Google-style)
// =============================================================================

func TestTypeSafeDateValidation(t *testing.T) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	t.Run("DateAfter", func(t *testing.T) {
		if !DateAfter(today.AddDate(0, 0, 1), today) {
			t.Error("tomorrow should be after today")
		}
		if DateAfter(today, today.AddDate(0, 0, 1)) {
			t.Error("today should not be after tomorrow")
		}
	})

	t.Run("DateBefore", func(t *testing.T) {
		if !DateBefore(today, today.AddDate(0, 0, 1)) {
			t.Error("today should be before tomorrow")
		}
		if DateBefore(today.AddDate(0, 0, 1), today) {
			t.Error("tomorrow should not be before today")
		}
	})

	t.Run("DateBetween", func(t *testing.T) {
		start := today.AddDate(0, 0, -7)
		end := today.AddDate(0, 0, 7)
		if !DateBetween(today, start, end) {
			t.Error("today should be between last week and next week")
		}
		if DateBetween(today.AddDate(0, 0, -10), start, end) {
			t.Error("10 days ago should not be between last week and next week")
		}
	})

	t.Run("IsToday", func(t *testing.T) {
		if !IsToday(now) {
			t.Error("now should be today")
		}
		if IsToday(today.AddDate(0, 0, -1)) {
			t.Error("yesterday should not be today")
		}
	})

	t.Run("IsFuture", func(t *testing.T) {
		if !IsFuture(now.Add(time.Hour)) {
			t.Error("1 hour from now should be in the future")
		}
		if IsFuture(now.Add(-time.Hour)) {
			t.Error("1 hour ago should not be in the future")
		}
	})

	t.Run("IsPast", func(t *testing.T) {
		if !IsPast(now.Add(-time.Hour)) {
			t.Error("1 hour ago should be in the past")
		}
		if IsPast(now.Add(time.Hour)) {
			t.Error("1 hour from now should not be in the past")
		}
	})
}

// =============================================================================
// Tests for DateRule fluent builder
// =============================================================================

func TestDateRuleBuilder(t *testing.T) {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	t.Run("After", func(t *testing.T) {
		rule := Date().After(today)
		if !rule.Validate(today.AddDate(0, 0, 1)) {
			t.Error("tomorrow should pass After(today)")
		}
		if rule.Validate(today.AddDate(0, 0, -1)) {
			t.Error("yesterday should fail After(today)")
		}
	})

	t.Run("Before", func(t *testing.T) {
		rule := Date().Before(today)
		if !rule.Validate(today.AddDate(0, 0, -1)) {
			t.Error("yesterday should pass Before(today)")
		}
		if rule.Validate(today.AddDate(0, 0, 1)) {
			t.Error("tomorrow should fail Before(today)")
		}
	})

	t.Run("Chained rules", func(t *testing.T) {
		// Must be after today and before 30 days from now
		rule := Date().After(today).Before(today.AddDate(0, 0, 30))

		if !rule.Validate(today.AddDate(0, 0, 15)) {
			t.Error("15 days from now should pass")
		}
		if rule.Validate(today.AddDate(0, 0, -1)) {
			t.Error("yesterday should fail")
		}
		if rule.Validate(today.AddDate(0, 0, 45)) {
			t.Error("45 days from now should fail")
		}
	})

	t.Run("Between", func(t *testing.T) {
		start := today.AddDate(0, 0, -7)
		end := today.AddDate(0, 0, 7)
		rule := Date().Between(start, end)

		if !rule.Validate(today) {
			t.Error("today should be between last week and next week")
		}
		if rule.Validate(today.AddDate(0, 0, -10)) {
			t.Error("10 days ago should fail")
		}
	})

	t.Run("ValidateWithError", func(t *testing.T) {
		rule := Date("birth_date").Before(today.AddDate(-18, 0, 0))

		err := rule.ValidateWithError(today.AddDate(-20, 0, 0))
		if err != nil {
			t.Errorf("20 years ago should pass, got: %v", err)
		}

		err = rule.ValidateWithError(today.AddDate(-16, 0, 0))
		if err == nil {
			t.Error("16 years ago should fail")
		}
	})
}

// =============================================================================
// Tests for TimeHelper
// =============================================================================

func TestTimeHelper(t *testing.T) {
	today := Today()

	t.Run("SubDays", func(t *testing.T) {
		result := T(today).SubDays(7)
		expected := today.AddDate(0, 0, -7)
		if !result.Equal(expected) {
			t.Errorf("SubDays(7) = %v, want %v", result, expected)
		}
	})

	t.Run("AddDays", func(t *testing.T) {
		result := T(today).AddDays(7)
		expected := today.AddDate(0, 0, 7)
		if !result.Equal(expected) {
			t.Errorf("AddDays(7) = %v, want %v", result, expected)
		}
	})

	t.Run("SubMonths", func(t *testing.T) {
		result := T(today).SubMonths(1)
		expected := today.AddDate(0, -1, 0)
		if !result.Equal(expected) {
			t.Errorf("SubMonths(1) = %v, want %v", result, expected)
		}
	})

	t.Run("SubYears", func(t *testing.T) {
		result := T(today).SubYears(18)
		expected := today.AddDate(-18, 0, 0)
		if !result.Equal(expected) {
			t.Errorf("SubYears(18) = %v, want %v", result, expected)
		}
	})
}

// =============================================================================
// Tests for helper functions
// =============================================================================

func TestDateHelpers(t *testing.T) {
	now := time.Now()
	todayExpected := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	t.Run("Today", func(t *testing.T) {
		result := Today()
		if !result.Equal(todayExpected) {
			t.Errorf("Today() = %v, want %v", result, todayExpected)
		}
	})

	t.Run("Tomorrow", func(t *testing.T) {
		result := Tomorrow()
		expected := todayExpected.AddDate(0, 0, 1)
		if !result.Equal(expected) {
			t.Errorf("Tomorrow() = %v, want %v", result, expected)
		}
	})

	t.Run("Yesterday", func(t *testing.T) {
		result := Yesterday()
		expected := todayExpected.AddDate(0, 0, -1)
		if !result.Equal(expected) {
			t.Errorf("Yesterday() = %v, want %v", result, expected)
		}
	})
}
