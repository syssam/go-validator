package validator

import (
	"testing"
)

// Consolidated numeric validation tests
// Contains tests for integer, float, and uint validation functions

// Float validation tests
func TestValidateDigitsBetweenFloat64(t *testing.T) {
	tests := []struct {
		value    float64
		left     float64
		right    float64
		expected bool
	}{
		{5.0, 1.0, 10.0, true},
		{0.0, 1.0, 10.0, false},
		{15.0, 1.0, 10.0, false},
		{1.0, 1.0, 10.0, true},
		{10.0, 1.0, 10.0, true},
		{5.5, 5.5, 5.5, true},
		{5.0, 10.0, 1.0, true}, // Test swapping when left > right
		{0.5, 10.0, 1.0, false},
		{15.0, 10.0, 1.0, false},
	}

	for _, test := range tests {
		result := ValidateDigitsBetweenFloat64(test.value, test.left, test.right)
		if result != test.expected {
			t.Errorf("ValidateDigitsBetweenFloat64(%f, %f, %f) = %t; expected %t",
				test.value, test.left, test.right, result, test.expected)
		}
	}
}

func TestValidateLtFloat64(t *testing.T) {
	tests := []struct {
		value    float64
		param    float64
		expected bool
	}{
		{5.0, 10.0, true},
		{10.0, 5.0, false},
		{5.5, 5.5, false},
		{-5.0, 0.0, true},
		{0.0, -5.0, false},
	}

	for _, test := range tests {
		result := ValidateLtFloat64(test.value, test.param)
		if result != test.expected {
			t.Errorf("ValidateLtFloat64(%f, %f) = %t; expected %t",
				test.value, test.param, result, test.expected)
		}
	}
}

func TestValidateLteFloat64(t *testing.T) {
	tests := []struct {
		value    float64
		param    float64
		expected bool
	}{
		{5.0, 10.0, true},
		{10.0, 5.0, false},
		{5.5, 5.5, true},
		{-5.0, 0.0, true},
		{0.0, -5.0, false},
	}

	for _, test := range tests {
		result := ValidateLteFloat64(test.value, test.param)
		if result != test.expected {
			t.Errorf("ValidateLteFloat64(%f, %f) = %t; expected %t",
				test.value, test.param, result, test.expected)
		}
	}
}

func TestValidateGteFloat64(t *testing.T) {
	tests := []struct {
		value    float64
		param    float64
		expected bool
	}{
		{10.0, 5.0, true},
		{5.0, 10.0, false},
		{5.5, 5.5, true},
		{0.0, -5.0, true},
		{-5.0, 0.0, false},
	}

	for _, test := range tests {
		result := ValidateGteFloat64(test.value, test.param)
		if result != test.expected {
			t.Errorf("ValidateGteFloat64(%f, %f) = %t; expected %t",
				test.value, test.param, result, test.expected)
		}
	}
}

func TestValidateGtFloat64(t *testing.T) {
	tests := []struct {
		value    float64
		param    float64
		expected bool
	}{
		{10.0, 5.0, true},
		{5.0, 10.0, false},
		{5.5, 5.5, false},
		{0.0, -5.0, true},
		{-5.0, 0.0, false},
	}

	for _, test := range tests {
		result := ValidateGtFloat64(test.value, test.param)
		if result != test.expected {
			t.Errorf("ValidateGtFloat64(%f, %f) = %t; expected %t",
				test.value, test.param, result, test.expected)
		}
	}
}

func TestCompareFloat64(t *testing.T) {
	tests := []struct {
		name        string
		first       float64
		second      float64
		operator    string
		expected    bool
		expectError bool
	}{
		{"Less than true", 5.5, 10.5, "<", true, false},
		{"Less than false", 10.5, 5.5, "<", false, false},
		{"Greater than true", 10.5, 5.5, ">", true, false},
		{"Greater than false", 5.5, 10.5, ">", false, false},
		{"Less than or equal true (less)", 5.5, 10.5, "<=", true, false},
		{"Less than or equal true (equal)", 5.5, 5.5, "<=", true, false},
		{"Less than or equal false", 10.5, 5.5, "<=", false, false},
		{"Greater than or equal true (greater)", 10.5, 5.5, ">=", true, false},
		{"Greater than or equal true (equal)", 5.5, 5.5, ">=", true, false},
		{"Greater than or equal false", 5.5, 10.5, ">=", false, false},
		{"Equal true", 5.5, 5.5, "==", true, false},
		{"Equal false", 5.5, 10.5, "==", false, false},
		{"Invalid operator", 5.5, 10.5, "!=", false, true},
		{"Invalid operator symbol", 5.5, 10.5, "invalid", false, true},
		{"Negative numbers", -10.5, -5.5, "<", true, false},
		{"Zero comparisons", 0.0, 1.5, "<", true, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := compareFloat64(test.first, test.second, test.operator)
			if test.expectError && err == nil {
				t.Errorf("Expected error for %s", test.name)
			}
			if !test.expectError && err != nil {
				t.Errorf("Unexpected error for %s: %v", test.name, err)
			}
			if !test.expectError && result != test.expected {
				t.Errorf("Expected %t for %s, got %t", test.expected, test.name, result)
			}
		})
	}
}

// Integer validation tests
func TestValidateDigitsBetweenInt64EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		value    int64
		left     int64
		right    int64
		expected bool
	}{
		{"Value between bounds", 5, 1, 10, true},
		{"Value at left bound", 1, 1, 10, true},
		{"Value at right bound", 10, 1, 10, true},
		{"Value below bounds", 0, 1, 10, false},
		{"Value above bounds", 15, 1, 10, false},
		{"Swapped bounds - value valid", 5, 10, 1, true},
		{"Swapped bounds - value invalid", 15, 10, 1, false},
		{"Negative values", -5, -10, -1, true},
		{"Negative values invalid", -15, -10, -1, false},
		{"Equal bounds", 5, 5, 5, true},
		{"Equal bounds invalid", 4, 5, 5, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := ValidateDigitsBetweenInt64(test.value, test.left, test.right)
			if result != test.expected {
				t.Errorf("ValidateDigitsBetweenInt64(%d, %d, %d) = %t; expected %t",
					test.value, test.left, test.right, result, test.expected)
			}
		})
	}
}

func TestCompareInt64AllOperators(t *testing.T) {
	tests := []struct {
		name        string
		first       int64
		second      int64
		operator    string
		expected    bool
		expectError bool
	}{
		{"Less than true", 5, 10, "<", true, false},
		{"Less than false", 10, 5, "<", false, false},
		{"Greater than true", 10, 5, ">", true, false},
		{"Greater than false", 5, 10, ">", false, false},
		{"Less than or equal true (less)", 5, 10, "<=", true, false},
		{"Less than or equal true (equal)", 5, 5, "<=", true, false},
		{"Less than or equal false", 10, 5, "<=", false, false},
		{"Greater than or equal true (greater)", 10, 5, ">=", true, false},
		{"Greater than or equal true (equal)", 5, 5, ">=", true, false},
		{"Greater than or equal false", 5, 10, ">=", false, false},
		{"Equal true", 5, 5, "==", true, false},
		{"Equal false", 5, 10, "==", false, false},
		{"Invalid operator", 5, 10, "!=", false, true},
		{"Invalid operator symbol", 5, 10, "invalid", false, true},
		{"Negative numbers less than", -10, -5, "<", true, false},
		{"Negative numbers greater than", -5, -10, ">", true, false},
		{"Zero comparisons", 0, 1, "<", true, false},
		{"Zero equal", 0, 0, "==", true, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := compareInt64(test.first, test.second, test.operator)
			if test.expectError && err == nil {
				t.Errorf("Expected error for %s", test.name)
			}
			if !test.expectError && err != nil {
				t.Errorf("Unexpected error for %s: %v", test.name, err)
			}
			if !test.expectError && result != test.expected {
				t.Errorf("Expected %t for %s, got %t", test.expected, test.name, result)
			}
		})
	}
}

// Additional comprehensive validation tests using the generic ValidateDigitsBetween function
func TestValidateDigitsBetweenGeneric(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		params   []string
		expected bool
		hasError bool
	}{
		// Supported integer types (checks number of digits, not value)
		{"int 1 digit valid", 5, []string{"1", "10"}, true, false},
		{"int 2 digits valid", 10, []string{"1", "10"}, true, false},
		{"int 1 digit boundary", 0, []string{"1", "10"}, true, false},
		{"int 3 digits valid", 123, []string{"1", "10"}, true, false},
		{"int64 1 digit", int64(5), []string{"1", "10"}, true, false},
		{"int64 2 digits", int64(12), []string{"1", "10"}, true, false},
		{"int64 too many digits", int64(12345678901), []string{"1", "10"}, false, false},

		// String digit validation (checks string length, must be numeric)
		{"string 1 digit", "5", []string{"1", "10"}, true, false},
		{"string 2 digits", "12", []string{"1", "10"}, true, false},
		{"string too long", "12345678901", []string{"1", "10"}, false, false},
		{"string non-numeric", "ab", []string{"1", "10"}, false, true},

		// Unsupported types should error
		{"uint unsupported", uint(5), []string{"1", "10"}, false, true},
		{"uint64 unsupported", uint64(5), []string{"1", "10"}, false, true},
		{"float32 unsupported", float32(5.5), []string{"1", "10"}, false, true},
		{"float64 unsupported", float64(5.5), []string{"1", "10"}, false, true},

		// Error cases
		{"wrong param count", 5, []string{"1"}, false, true},
		{"invalid param", 5, []string{"invalid", "10"}, false, true},
		{"complex type", complex64(1 + 2i), []string{"1", "10"}, false, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := ValidateDigitsBetween(test.value, test.params)
			if test.hasError {
				if err == nil {
					t.Errorf("Expected error for %s", test.name)
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error for %s: %v", test.name, err)
				return
			}
			if result != test.expected {
				t.Errorf("Expected %t for %s, got %t", test.expected, test.name, result)
			}
		})
	}
}
