package validator

import (
	"reflect"
	"testing"
)

// Test direct validation functions to improve coverage
func TestCoverageImprovementDirectFunctions(t *testing.T) {
	// Test findField function
	type TestStruct struct {
		Name string
		Age  int
	}
	testStruct := TestStruct{Name: "test", Age: 25}
	v := reflect.ValueOf(testStruct)

	// Test successful field finding
	field, err := findField("Name", v)
	if err != nil {
		t.Errorf("findField should not error for valid field: %v", err)
	}
	if field.String() != "test" {
		t.Errorf("Expected 'test', got %v", field.String())
	}

	// Test non-existent field (returns zero Value, not error)
	field, err = findField("NonExistent", v)
	if err != nil {
		t.Errorf("findField should not error for non-existent field: %v", err)
	}
	if field.IsValid() {
		t.Error("findField should return invalid Value for non-existent field")
	}

	// Test with non-struct
	nonStruct := reflect.ValueOf("not a struct")
	_, err = findField("field", nonStruct)
	if err == nil {
		t.Error("findField should error for non-struct value")
	}
}

// Test validateRequiredWith function directly
func TestValidateRequiredWithDirect(t *testing.T) {
	type TestStruct struct {
		Field1 string
		Field2 string
	}

	// Case 1: Field2 present, Field1 present - should be valid
	testStruct := TestStruct{Field1: "value1", Field2: "value2"}
	objValue := reflect.ValueOf(testStruct)
	field1Value := reflect.ValueOf("value1")

	result := validateRequiredWith([]string{"Field2"}, field1Value, objValue)
	if !result {
		t.Error("validateRequiredWith should return true when both fields are present")
	}

	// Case 2: Field2 present, Field1 empty - should be invalid
	field1EmptyValue := reflect.ValueOf("")
	result = validateRequiredWith([]string{"Field2"}, field1EmptyValue, objValue)
	if result {
		t.Error("validateRequiredWith should return false when Field2 is present but Field1 is empty")
	}

	// Case 3: Field2 empty, Field1 empty - should be valid
	testStructEmpty := TestStruct{Field1: "", Field2: ""}
	objValueEmpty := reflect.ValueOf(testStructEmpty)
	result = validateRequiredWith([]string{"Field2"}, field1EmptyValue, objValueEmpty)
	if !result {
		t.Error("validateRequiredWith should return true when both fields are empty")
	}
}

// Test validateRequiredWithAll function directly
func TestValidateRequiredWithAllDirect(t *testing.T) {
	type TestStruct struct {
		Field1 string
		Field2 string
		Field3 string
	}

	// Case 1: All fields present - Field1 should be required and present
	testStruct := TestStruct{Field1: "value1", Field2: "value2", Field3: "value3"}
	objValue := reflect.ValueOf(testStruct)
	field1Value := reflect.ValueOf("value1")

	result := validateRequiredWithAll([]string{"Field2", "Field3"}, field1Value, objValue)
	if !result {
		t.Error("validateRequiredWithAll should return true when all fields are present")
	}

	// Case 2: Field2 and Field3 present, Field1 empty - should be invalid
	field1EmptyValue := reflect.ValueOf("")
	result = validateRequiredWithAll([]string{"Field2", "Field3"}, field1EmptyValue, objValue)
	if result {
		t.Error("validateRequiredWithAll should return false when Field2 and Field3 are present but Field1 is empty")
	}

	// Case 3: Only Field2 present, Field1 empty - should be valid (Field1 not required)
	testStructPartial := TestStruct{Field1: "", Field2: "value2", Field3: ""}
	objValuePartial := reflect.ValueOf(testStructPartial)
	result = validateRequiredWithAll([]string{"Field2", "Field3"}, field1EmptyValue, objValuePartial)
	if !result {
		t.Error("validateRequiredWithAll should return true when not all required fields are present")
	}
}

// Test validateRequiredWithout function directly
func TestValidateRequiredWithoutDirect(t *testing.T) {
	type TestStruct struct {
		Field1 string
		Field2 string
	}

	// Case 1: Field2 absent, Field1 present - should be valid
	testStruct := TestStruct{Field1: "value1", Field2: ""}
	objValue := reflect.ValueOf(testStruct)
	field1Value := reflect.ValueOf("value1")

	result := validateRequiredWithout([]string{"Field2"}, field1Value, objValue)
	if !result {
		t.Error("validateRequiredWithout should return true when Field2 is absent and Field1 is present")
	}

	// Case 2: Field2 absent, Field1 absent - should be invalid
	field1EmptyValue := reflect.ValueOf("")
	result = validateRequiredWithout([]string{"Field2"}, field1EmptyValue, objValue)
	if result {
		t.Error("validateRequiredWithout should return false when Field2 is absent and Field1 is also absent")
	}

	// Case 3: Field2 present, Field1 absent - should be valid (Field1 not required)
	testStructWithField2 := TestStruct{Field1: "", Field2: "value2"}
	objValueWithField2 := reflect.ValueOf(testStructWithField2)
	result = validateRequiredWithout([]string{"Field2"}, field1EmptyValue, objValueWithField2)
	if !result {
		t.Error("validateRequiredWithout should return true when Field2 is present (Field1 not required)")
	}
}

// Test validateRequiredWithoutAll function directly
func TestValidateRequiredWithoutAllDirect(t *testing.T) {
	type TestStruct struct {
		Field1 string
		Field2 string
		Field3 string
	}

	// Case 1: All fields absent, Field1 present - should be valid
	testStruct := TestStruct{Field1: "value1", Field2: "", Field3: ""}
	objValue := reflect.ValueOf(testStruct)
	field1Value := reflect.ValueOf("value1")

	result := validateRequiredWithoutAll([]string{"Field2", "Field3"}, field1Value, objValue)
	if !result {
		t.Error("validateRequiredWithoutAll should return true when all other fields are absent and Field1 is present")
	}

	// Case 2: All fields absent, Field1 absent - should be invalid
	field1EmptyValue := reflect.ValueOf("")
	result = validateRequiredWithoutAll([]string{"Field2", "Field3"}, field1EmptyValue, objValue)
	if result {
		t.Error("validateRequiredWithoutAll should return false when all fields including Field1 are absent")
	}

	// Case 3: Some field present, Field1 absent - should be valid (Field1 not required)
	testStructPartial := TestStruct{Field1: "", Field2: "value2", Field3: ""}
	objValuePartial := reflect.ValueOf(testStructPartial)
	result = validateRequiredWithoutAll([]string{"Field2", "Field3"}, field1EmptyValue, objValuePartial)
	if !result {
		t.Error("validateRequiredWithoutAll should return true when some other fields are present (Field1 not required)")
	}
}

// Test parameter-based comparison validators directly
func TestParameterValidatorsDirect(t *testing.T) {
	// Test validateGtParam
	stringValue := reflect.ValueOf("hello")
	result, err := validateGtParam(stringValue, []string{"3"})
	if err != nil || !result {
		t.Errorf("validateGtParam should return true for 'hello' > 3 characters: result=%v, err=%v", result, err)
	}

	shortString := reflect.ValueOf("hi")
	result, err = validateGtParam(shortString, []string{"3"})
	if err != nil || result {
		t.Errorf("validateGtParam should return false for 'hi' > 3 characters: result=%v, err=%v", result, err)
	}

	// Test validateGteParam
	exactString := reflect.ValueOf("test")
	result, err = validateGteParam(exactString, []string{"4"})
	if err != nil || !result {
		t.Errorf("validateGteParam should return true for 'test' >= 4 characters: result=%v, err=%v", result, err)
	}

	// Test validateLtParam
	result, err = validateLtParam(shortString, []string{"3"})
	if err != nil || !result {
		t.Errorf("validateLtParam should return true for 'hi' < 3 characters: result=%v, err=%v", result, err)
	}

	// Test validateLteParam
	result, err = validateLteParam(exactString, []string{"4"})
	if err != nil || !result {
		t.Errorf("validateLteParam should return true for 'test' <= 4 characters: result=%v, err=%v", result, err)
	}

	// Test with integer values
	intValue := reflect.ValueOf(5)
	result, err = validateGtParam(intValue, []string{"3"})
	if err != nil || !result {
		t.Errorf("validateGtParam should return true for 5 > 3: result=%v, err=%v", result, err)
	}

	// Test error cases
	_, err = validateGtParam(stringValue, []string{})
	if err == nil {
		t.Error("validateGtParam should return error for empty params")
	}

	_, err = validateGtParam(stringValue, []string{"invalid"})
	if err == nil {
		t.Error("validateGtParam should return error for invalid numeric param")
	}
}
