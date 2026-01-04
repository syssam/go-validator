package validator

import "testing"

func TestIsCountryAlpha2(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"US", true},
		{"GB", true},
		{"CN", true},
		{"TW", true},
		{"JP", true},
		{"DE", true},
		{"FR", true},
		{"us", true}, // lowercase should work
		{"gb", true},
		{"XX", false},
		{"USA", false}, // 3-letter code
		{"", false},
		{"A", false},
		{"ABC", false},
	}
	for _, test := range tests {
		result := IsCountryAlpha2(test.input)
		if result != test.expected {
			t.Errorf("IsCountryAlpha2(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestIsCountryAlpha3(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"USA", true},
		{"GBR", true},
		{"CHN", true},
		{"TWN", true},
		{"JPN", true},
		{"DEU", true},
		{"FRA", true},
		{"usa", true}, // lowercase should work
		{"gbr", true},
		{"XXX", false},
		{"US", false}, // 2-letter code
		{"", false},
		{"AB", false},
		{"ABCD", false},
	}
	for _, test := range tests {
		result := IsCountryAlpha3(test.input)
		if result != test.expected {
			t.Errorf("IsCountryAlpha3(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestIsCountryNumeric(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"840", true}, // USA
		{"826", true}, // GBR
		{"156", true}, // CHN
		{"158", true}, // TWN
		{"392", true}, // JPN
		{"276", true}, // DEU
		{"250", true}, // FRA
		{"000", false},
		{"999", false},
		{"", false},
	}
	for _, test := range tests {
		result := IsCountryNumeric(test.input)
		if result != test.expected {
			t.Errorf("IsCountryNumeric(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestIsCountryCode(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"US", true},
		{"USA", true},
		{"840", true},
		{"GB", true},
		{"GBR", true},
		{"826", true},
		{"XX", false},
		{"XXX", false},
		{"000", false},
		{"", false},
	}
	for _, test := range tests {
		result := IsCountryCode(test.input)
		if result != test.expected {
			t.Errorf("IsCountryCode(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

// Benchmark tests
func BenchmarkIsCountryAlpha2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsCountryAlpha2("US")
	}
}

func BenchmarkIsCountryAlpha3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsCountryAlpha3("USA")
	}
}

func BenchmarkIsCountryNumeric(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsCountryNumeric("840")
	}
}
