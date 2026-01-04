package validator

import "testing"

func TestIsLanguageAlpha2(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"en", true},
		{"zh", true},
		{"ja", true},
		{"ko", true},
		{"fr", true},
		{"de", true},
		{"es", true},
		{"pt", true},
		{"ru", true},
		{"ar", true},
		{"EN", true}, // uppercase should work
		{"ZH", true},
		{"xx", false},
		{"eng", false}, // 3-letter code
		{"", false},
		{"e", false},
		{"abc", false},
	}
	for _, test := range tests {
		result := IsLanguageAlpha2(test.input)
		if result != test.expected {
			t.Errorf("IsLanguageAlpha2(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestIsLanguageAlpha3(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"eng", true},
		{"zho", true},
		{"jpn", true},
		{"kor", true},
		{"fra", true},
		{"deu", true},
		{"spa", true},
		{"por", true},
		{"rus", true},
		{"ara", true},
		{"ENG", true}, // uppercase should work
		{"ZHO", true},
		{"xxx", false},
		{"en", false}, // 2-letter code
		{"", false},
		{"ab", false},
		{"abcd", false},
	}
	for _, test := range tests {
		result := IsLanguageAlpha3(test.input)
		if result != test.expected {
			t.Errorf("IsLanguageAlpha3(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

func TestIsLanguageCode(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"en", true},
		{"eng", true},
		{"zh", true},
		{"zho", true},
		{"ja", true},
		{"jpn", true},
		{"xx", false},
		{"xxx", false},
		{"", false},
	}
	for _, test := range tests {
		result := IsLanguageCode(test.input)
		if result != test.expected {
			t.Errorf("IsLanguageCode(%q) = %v, expected %v", test.input, result, test.expected)
		}
	}
}

// Benchmark tests
func BenchmarkIsLanguageAlpha2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsLanguageAlpha2("en")
	}
}

func BenchmarkIsLanguageAlpha3(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsLanguageAlpha3("eng")
	}
}
