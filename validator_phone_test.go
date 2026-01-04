package validator

import "testing"

func TestIsPhoneE164(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"+14155551234", true},
		{"+442071234567", true},
		{"+8613812345678", true},
		{"+81312345678", true},
		{"+886912345678", true},
		{"+821012345678", true},
		{"+6591234567", true},  // Singapore 8-digit mobile
		{"+33123456789", true},
		{"+4930123456", true},
		{"14155551234", false},         // missing +
		{"+1415555", false},            // too short
		{"+1234567890123456789", false}, // too long
		{"", false},
		{"+", false},
		{"+abc", false},
	}
	for _, test := range tests {
		result := IsPhoneE164(test.input)
		if result != test.expected {
			t.Errorf("IsPhoneE164(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestIsPhone(t *testing.T) {
	tests := []struct {
		input    string
		region   string
		expected bool
	}{
		{"(415) 555-1234", "US", true},
		{"415-555-1234", "US", true},
		{"4155551234", "US", true},
		{"+14155551234", "US", true},
		{"020 7123 4567", "GB", true},
		{"+442071234567", "GB", true},
		{"13812345678", "CN", true},
		{"+8613812345678", "CN", true},
		{"03-1234-5678", "JP", true},
		{"+81312345678", "JP", true},
		{"0912345678", "TW", true},
		{"+886912345678", "TW", true},
		{"010-1234-5678", "KR", true},
		{"+821012345678", "KR", true},
		{"invalid", "US", false},
		{"", "US", false},
		{"123", "US", false},
	}
	for _, test := range tests {
		result := IsPhone(test.input, test.region)
		if result != test.expected {
			t.Errorf("IsPhone(%q, %q) = %v, expected %v", test.input, test.region, result, test.expected)
		}
	}
}

func TestIsPhoneValid(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"+14155551234", true},
		{"+442071234567", true},
		{"+8613812345678", true},
		{"14155551234", false}, // missing +
		{"4155551234", false},  // no country code
		{"", false},
	}
	for _, test := range tests {
		result := IsPhoneValid(test.input)
		if result != test.expected {
			t.Errorf("IsPhoneValid(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestIsPhoneMobile(t *testing.T) {
	tests := []struct {
		input    string
		region   string
		expected bool
	}{
		{"4155551234", "US", true},      // US mobile
		{"13812345678", "CN", true},     // CN mobile
		{"0912345678", "TW", true},      // TW mobile
		{"09012345678", "JP", true},     // JP mobile
		{"01012345678", "KR", true},     // KR mobile
		{"07123456789", "GB", true},     // UK mobile
		{"invalid", "US", false},
		{"", "US", false},
	}
	for _, test := range tests {
		result := IsPhoneMobile(test.input, test.region)
		if result != test.expected {
			t.Errorf("IsPhoneMobile(%q, %q) = %v, expected %v", test.input, test.region, result, test.expected)
		}
	}
}

func TestFormatPhoneE164(t *testing.T) {
	tests := []struct {
		input    string
		region   string
		expected string
	}{
		{"4155551234", "US", "+14155551234"},
		{"(415) 555-1234", "US", "+14155551234"},
		{"020 7123 4567", "GB", "+442071234567"},
		{"13812345678", "CN", "+8613812345678"},
	}
	for _, test := range tests {
		result, err := FormatPhoneE164(test.input, test.region)
		if err != nil {
			t.Errorf("FormatPhoneE164(%q, %q) error: %v", test.input, test.region, err)
			continue
		}
		if result != test.expected {
			t.Errorf("FormatPhoneE164(%q, %q) = %v, expected %v", test.input, test.region, result, test.expected)
		}
	}
}

func TestFormatPhoneNational(t *testing.T) {
	tests := []struct {
		input    string
		region   string
		expected string
	}{
		{"+14155551234", "US", "(415) 555-1234"},
		{"+442071234567", "GB", "020 7123 4567"},
	}
	for _, test := range tests {
		result, err := FormatPhoneNational(test.input, test.region)
		if err != nil {
			t.Errorf("FormatPhoneNational(%q, %q) error: %v", test.input, test.region, err)
			continue
		}
		if result != test.expected {
			t.Errorf("FormatPhoneNational(%q, %q) = %v, expected %v", test.input, test.region, result, test.expected)
		}
	}
}

func TestFormatPhoneInternational(t *testing.T) {
	tests := []struct {
		input    string
		region   string
		expected string
	}{
		{"+14155551234", "US", "+1 415-555-1234"},
		{"+442071234567", "GB", "+44 20 7123 4567"},
	}
	for _, test := range tests {
		result, err := FormatPhoneInternational(test.input, test.region)
		if err != nil {
			t.Errorf("FormatPhoneInternational(%q, %q) error: %v", test.input, test.region, err)
			continue
		}
		if result != test.expected {
			t.Errorf("FormatPhoneInternational(%q, %q) = %v, expected %v", test.input, test.region, result, test.expected)
		}
	}
}

// Benchmark tests
func BenchmarkIsPhoneE164(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsPhoneE164("+14155551234")
	}
}

func BenchmarkIsPhone(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsPhone("4155551234", "US")
	}
}

func BenchmarkFormatPhoneE164(b *testing.B) {
	for i := 0; i < b.N; i++ {
		FormatPhoneE164("4155551234", "US")
	}
}
