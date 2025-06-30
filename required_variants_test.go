package validator

import (
	"testing"
)

// Test struct for RequiredWith validation
type RequiredWithTest struct {
	Field1 string `valid:"requiredWith=Field2"`
	Field2 string
}

// Test RequiredWith validation
func TestRequiredWithValidation(t *testing.T) {
	tests := []struct {
		name     string
		data     RequiredWithTest
		expected bool
	}{
		{"Both fields present - valid", RequiredWithTest{Field1: "value1", Field2: "value2"}, true},
		{"Field2 present, Field1 missing - invalid", RequiredWithTest{Field2: "value2"}, false},
		{"Field2 absent, Field1 absent - valid", RequiredWithTest{}, true},
		{"Field2 absent, Field1 present - valid", RequiredWithTest{Field1: "value1"}, true},
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

// Test struct for RequiredWithAll validation
type RequiredWithAllTest struct {
	Field1 string `valid:"requiredWithAll=Field2,Field3"`
	Field2 string
	Field3 string
}

// Test RequiredWithAll validation
func TestRequiredWithAllValidation(t *testing.T) {
	tests := []struct {
		name     string
		data     RequiredWithAllTest
		expected bool
	}{
		{"All fields present - valid", RequiredWithAllTest{Field1: "v1", Field2: "v2", Field3: "v3"}, true},
		{"Field2 and Field3 present, Field1 missing - invalid", RequiredWithAllTest{Field2: "v2", Field3: "v3"}, false},
		{"Only Field2 present - valid", RequiredWithAllTest{Field2: "v2"}, true},
		{"Only Field3 present - valid", RequiredWithAllTest{Field3: "v3"}, true},
		{"No fields present - valid", RequiredWithAllTest{}, true},
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

// Test struct for RequiredWithout validation
type RequiredWithoutTest struct {
	Field1 string `valid:"requiredWithout=Field2"`
	Field2 string
}

// Test RequiredWithout validation
func TestRequiredWithoutValidation(t *testing.T) {
	tests := []struct {
		name     string
		data     RequiredWithoutTest
		expected bool
	}{
		{"Field2 absent, Field1 present - valid", RequiredWithoutTest{Field1: "value1"}, true},
		{"Field2 absent, Field1 absent - invalid", RequiredWithoutTest{}, false},
		{"Field2 present, Field1 present - valid", RequiredWithoutTest{Field1: "v1", Field2: "v2"}, true},
		{"Field2 present, Field1 absent - valid", RequiredWithoutTest{Field2: "v2"}, true},
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

// Test struct for RequiredWithoutAll validation
type RequiredWithoutAllTest struct {
	Field1 string `valid:"requiredWithoutAll=Field2,Field3"`
	Field2 string
	Field3 string
}

// Test RequiredWithoutAll validation
func TestRequiredWithoutAllValidation(t *testing.T) {
	tests := []struct {
		name     string
		data     RequiredWithoutAllTest
		expected bool
	}{
		{"All fields absent, Field1 present - valid", RequiredWithoutAllTest{Field1: "v1"}, true},
		{"All fields absent, Field1 absent - invalid", RequiredWithoutAllTest{}, false},
		{"Field2 present - valid", RequiredWithoutAllTest{Field2: "v2"}, true},
		{"Field3 present - valid", RequiredWithoutAllTest{Field3: "v3"}, true},
		{"Both Field2 and Field3 present - valid", RequiredWithoutAllTest{Field2: "v2", Field3: "v3"}, true},
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

// Test struct for RequiredUnless validation
type RequiredUnlessTest struct {
	Field1 string `valid:"requiredUnless=Field2=exempt"`
	Field2 string
}

// Test RequiredUnless validation
func TestRequiredUnlessValidation(t *testing.T) {
	tests := []struct {
		name     string
		data     RequiredUnlessTest
		expected bool
	}{
		{"Field2 is exempt, Field1 absent - valid", RequiredUnlessTest{Field2: "exempt"}, true},
		{"Field2 is exempt, Field1 present - valid", RequiredUnlessTest{Field1: "v1", Field2: "exempt"}, true},
		{"Field2 not exempt, Field1 present - valid", RequiredUnlessTest{Field1: "v1", Field2: "other"}, true},
		{"Field2 not exempt, Field1 absent - invalid", RequiredUnlessTest{Field2: "other"}, false},
		{"Field2 absent, Field1 present - valid", RequiredUnlessTest{Field1: "v1"}, true},
		{"Field2 absent, Field1 absent - invalid", RequiredUnlessTest{}, false},
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

// Test struct for RequiredIf with multiple conditions
type RequiredIfAdvancedTest struct {
	Field1 string `valid:"requiredIf=Field2=active,Field3=enabled"`
	Field2 string
	Field3 string
}

// Test RequiredIf with multiple conditions
func TestRequiredIfAdvancedValidation(t *testing.T) {
	tests := []struct {
		name     string
		data     RequiredIfAdvancedTest
		expected bool
	}{
		{"Condition met, field present - valid", RequiredIfAdvancedTest{Field1: "v1", Field2: "active"}, true},
		{"Condition met, field absent - invalid", RequiredIfAdvancedTest{Field2: "active"}, false},
		{"Condition not met, field absent - valid", RequiredIfAdvancedTest{Field2: "inactive"}, true},
		{"Condition not met, field present - valid", RequiredIfAdvancedTest{Field1: "v1", Field2: "inactive"}, true},
		{"Multiple conditions - one met", RequiredIfAdvancedTest{Field1: "v1", Field2: "active", Field3: "disabled"}, true},
		{"Multiple conditions - both met", RequiredIfAdvancedTest{Field1: "v1", Field2: "active", Field3: "enabled"}, true},
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
