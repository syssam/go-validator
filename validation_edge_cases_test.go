package validator

import (
	"reflect"
	"testing"
)

// Test validateSame edge cases
func TestValidateSameEdgeCases(t *testing.T) {
	type SameTest struct {
		Password        string `valid:"required"`
		ConfirmPassword string `valid:"same=Password"`
	}

	tests := []struct {
		name     string
		data     SameTest
		expected bool
	}{
		{"Passwords match - valid", SameTest{Password: "secret123", ConfirmPassword: "secret123"}, true},
		{"Passwords don't match - invalid", SameTest{Password: "secret123", ConfirmPassword: "different"}, false},
		{"Empty passwords match - valid", SameTest{Password: "", ConfirmPassword: ""}, true},
		{"One empty, one filled - invalid", SameTest{Password: "secret", ConfirmPassword: ""}, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateStruct(test.data)
			actual := err == nil
			if actual != test.expected {
				t.Errorf("Expected %t for %s, got %t. Error: %v", test.expected, test.name, actual, err)
			}
		})
	}
}

// Test validateLt edge cases
func TestValidateLtEdgeCases(t *testing.T) {
	type LtTest struct {
		Value string `valid:"lt=10"`
	}

	tests := []struct {
		name     string
		data     LtTest
		expected bool
	}{
		{"Short string - valid", LtTest{Value: "test"}, true},
		{"Long string - invalid", LtTest{Value: "this is a very long string"}, false},
		{"Empty string - valid", LtTest{Value: ""}, true},
		{"Exactly 9 chars - valid", LtTest{Value: "123456789"}, true},
		{"Exactly 10 chars - invalid", LtTest{Value: "1234567890"}, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateStruct(test.data)
			actual := err == nil
			if actual != test.expected {
				t.Errorf("Expected %t for %s, got %t. Error: %v", test.expected, test.name, actual, err)
			}
		})
	}
}

// Test validateLte edge cases
func TestValidateLteEdgeCases(t *testing.T) {
	type LteTest struct {
		Value string `valid:"lte=5"`
	}

	tests := []struct {
		name     string
		data     LteTest
		expected bool
	}{
		{"Short string - valid", LteTest{Value: "test"}, true},
		{"Exactly 5 chars - valid", LteTest{Value: "hello"}, true},
		{"6 chars - invalid", LteTest{Value: "hello!"}, false},
		{"Empty string - valid", LteTest{Value: ""}, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateStruct(test.data)
			actual := err == nil
			if actual != test.expected {
				t.Errorf("Expected %t for %s, got %t. Error: %v", test.expected, test.name, actual, err)
			}
		})
	}
}

// Test validateGt edge cases
func TestValidateGtEdgeCases(t *testing.T) {
	type GtTest struct {
		Value string `valid:"gt=3"`
	}

	tests := []struct {
		name     string
		data     GtTest
		expected bool
	}{
		{"4 chars - valid", GtTest{Value: "test"}, true},
		{"Exactly 3 chars - invalid", GtTest{Value: "abc"}, false},
		{"2 chars - invalid", GtTest{Value: "ab"}, false},
		{"Empty string - invalid", GtTest{Value: ""}, false},
		{"10 chars - valid", GtTest{Value: "verylongst"}, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateStruct(test.data)
			actual := err == nil
			if actual != test.expected {
				t.Errorf("Expected %t for %s, got %t. Error: %v", test.expected, test.name, actual, err)
			}
		})
	}
}

// Test validateGte edge cases
func TestValidateGteEdgeCases(t *testing.T) {
	type GteTest struct {
		Value string `valid:"gte=4"`
	}

	tests := []struct {
		name     string
		data     GteTest
		expected bool
	}{
		{"Exactly 4 chars - valid", GteTest{Value: "test"}, true},
		{"5 chars - valid", GteTest{Value: "tests"}, true},
		{"3 chars - invalid", GteTest{Value: "abc"}, false},
		{"Empty string - invalid", GteTest{Value: ""}, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateStruct(test.data)
			actual := err == nil
			if actual != test.expected {
				t.Errorf("Expected %t for %s, got %t. Error: %v", test.expected, test.name, actual, err)
			}
		})
	}
}

// Test validateSize edge cases
func TestValidateSizeEdgeCases(t *testing.T) {
	type SizeTest struct {
		Value string `valid:"size=5"`
	}

	tests := []struct {
		name     string
		data     SizeTest
		expected bool
	}{
		{"Exactly 5 chars - valid", SizeTest{Value: "hello"}, true},
		{"4 chars - invalid", SizeTest{Value: "test"}, false},
		{"6 chars - invalid", SizeTest{Value: "hellos"}, false},
		{"Empty string - invalid", SizeTest{Value: ""}, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateStruct(test.data)
			actual := err == nil
			if actual != test.expected {
				t.Errorf("Expected %t for %s, got %t. Error: %v", test.expected, test.name, actual, err)
			}
		})
	}
}

// Test validateMax edge cases
func TestValidateMaxEdgeCases(t *testing.T) {
	type MaxTest struct {
		Value string `valid:"max=8"`
	}

	tests := []struct {
		name     string
		data     MaxTest
		expected bool
	}{
		{"Under max - valid", MaxTest{Value: "short"}, true},
		{"Exactly max - valid", MaxTest{Value: "exactly8"}, true},
		{"Over max - invalid", MaxTest{Value: "toolongstring"}, false},
		{"Empty string - valid", MaxTest{Value: ""}, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateStruct(test.data)
			actual := err == nil
			if actual != test.expected {
				t.Errorf("Expected %t for %s, got %t. Error: %v", test.expected, test.name, actual, err)
			}
		})
	}
}

// Test validateMin edge cases
func TestValidateMinEdgeCases(t *testing.T) {
	type MinTest struct {
		Value string `valid:"min=3"`
	}

	tests := []struct {
		name     string
		data     MinTest
		expected bool
	}{
		{"Above min - valid", MinTest{Value: "test"}, true},
		{"Exactly min - valid", MinTest{Value: "abc"}, true},
		{"Below min - invalid", MinTest{Value: "ab"}, false},
		{"Empty string - invalid", MinTest{Value: ""}, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateStruct(test.data)
			actual := err == nil
			if actual != test.expected {
				t.Errorf("Expected %t for %s, got %t. Error: %v", test.expected, test.name, actual, err)
			}
		})
	}
}

// Test validateDistinct with different scenarios
func TestValidateDistinctAdvanced(t *testing.T) {
	type DistinctTest struct {
		Items []string `valid:"distinct"`
	}

	tests := []struct {
		name     string
		data     DistinctTest
		expected bool
	}{
		{"All unique - valid", DistinctTest{Items: []string{"a", "b", "c", "d"}}, true},
		{"Has duplicates - invalid", DistinctTest{Items: []string{"a", "b", "c", "a"}}, false},
		{"Empty slice - valid", DistinctTest{Items: []string{}}, true},
		{"Single item - valid", DistinctTest{Items: []string{"a"}}, true},
		{"Two identical - invalid", DistinctTest{Items: []string{"a", "a"}}, false},
		{"Case sensitive - valid", DistinctTest{Items: []string{"A", "a", "B", "b"}}, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateStruct(test.data)
			actual := err == nil
			if actual != test.expected {
				t.Errorf("Expected %t for %s, got %t. Error: %v", test.expected, test.name, actual, err)
			}
		})
	}
}

// Test direct validation functions to increase coverage
func TestDirectValidationFunctions(t *testing.T) {
	// Test validateDistinct function directly
	result, err := validateDistinct(reflect.ValueOf([]string{"a", "b", "c"}))
	if err != nil || !result {
		t.Error("Expected distinct values to validate")
	}

	result, err = validateDistinct(reflect.ValueOf([]string{"a", "b", "a"}))
	if err != nil || result {
		t.Error("Expected duplicate values to fail")
	}

	// Test empty slice
	result, err = validateDistinct(reflect.ValueOf([]string{}))
	if err != nil || !result {
		t.Error("Expected empty slice to validate")
	}

	// Test single item
	result, err = validateDistinct(reflect.ValueOf([]string{"single"}))
	if err != nil || !result {
		t.Error("Expected single item to validate")
	}
}
