package validator

import (
	"encoding/json"
	"errors"
	"testing"
)

func TestErrorsError(t *testing.T) {
	// Test empty errors
	var errs Errors
	if errs.Error() != "" {
		t.Errorf("Expected empty string for empty errors, got %s", errs.Error())
	}

	// Test single error
	errs = Errors{errors.New("single error")}
	if errs.Error() != "single error" {
		t.Errorf("Expected 'single error', got %s", errs.Error())
	}

	// Test multiple errors
	errs = Errors{
		errors.New("first error"),
		errors.New("second error"),
		errors.New("third error"),
	}
	expected := "first error\nsecond error\nthird error"
	if errs.Error() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, errs.Error())
	}
}

func TestErrorsErrors(t *testing.T) {
	err1 := errors.New("error 1")
	err2 := errors.New("error 2")
	errs := Errors{err1, err2}

	result := errs.Errors()
	if len(result) != 2 {
		t.Errorf("Expected 2 errors, got %d", len(result))
	}
	if result[0] != err1 || result[1] != err2 {
		t.Error("Errors() did not return the original errors")
	}
}

func TestErrorsFieldErrors(t *testing.T) {
	fieldErr := &FieldError{Name: "email", Message: "invalid email"}
	genericErr := errors.New("generic error")
	errs := Errors{fieldErr, genericErr}

	fieldErrors := errs.FieldErrors()
	if len(fieldErrors) != 2 {
		t.Errorf("Expected 2 field errors, got %d", len(fieldErrors))
	}

	if fieldErrors[0].Name != "email" {
		t.Error("First field error should be the original FieldError")
	}

	if fieldErrors[1].Message != "generic error" {
		t.Error("Second field error should be converted from generic error")
	}
}

func TestErrorsHasFieldError(t *testing.T) {
	fieldErr := &FieldError{Name: "email", Message: "invalid email"}
	errs := Errors{fieldErr}

	if !errs.HasFieldError("email") {
		t.Error("Expected HasFieldError to return true for existing field")
	}

	if errs.HasFieldError("name") {
		t.Error("Expected HasFieldError to return false for non-existing field")
	}
}

func TestErrorsGetFieldError(t *testing.T) {
	fieldErr := &FieldError{Name: "email", Message: "invalid email"}
	errs := Errors{fieldErr}

	result := errs.GetFieldError("email")
	if result == nil {
		t.Error("Expected GetFieldError to return the field error")
	} else if result.Name != "email" {
		t.Error("Expected field error name to be 'email'")
	}

	result = errs.GetFieldError("name")
	if result != nil {
		t.Error("Expected GetFieldError to return nil for non-existing field")
	}
}

func TestErrorsGroupByField(t *testing.T) {
	fieldErr1 := &FieldError{Name: "email", Message: "required"}
	fieldErr2 := &FieldError{Name: "email", Message: "invalid format"}
	fieldErr3 := &FieldError{Name: "name", Message: "required"}
	errs := Errors{fieldErr1, fieldErr2, fieldErr3}

	groups := errs.GroupByField()
	if len(groups) != 2 {
		t.Errorf("Expected 2 groups, got %d", len(groups))
	}

	if len(groups["email"]) != 2 {
		t.Errorf("Expected 2 errors for email field, got %d", len(groups["email"]))
	}

	if len(groups["name"]) != 1 {
		t.Errorf("Expected 1 error for name field, got %d", len(groups["name"]))
	}
}

func TestErrorsMarshalJSON(t *testing.T) {
	// Test empty errors
	var errs Errors
	data, err := json.Marshal(errs)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if string(data) != "[]" {
		t.Errorf("Expected '[]', got %s", string(data))
	}

	// Test with field errors
	fieldErr := &FieldError{Name: "email", Message: "invalid email"}
	errs = Errors{fieldErr}

	data, err = json.Marshal(errs)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	var responses []ErrorResponse
	err = json.Unmarshal(data, &responses)
	if err != nil {
		t.Errorf("Failed to unmarshal JSON: %v", err)
	}

	if len(responses) != 1 {
		t.Errorf("Expected 1 response, got %d", len(responses))
	}

	if responses[0].Message != "invalid email" || responses[0].Parameter != "email" {
		t.Error("JSON output does not match expected format")
	}
}

func TestFieldErrorError(t *testing.T) {
	// Test with custom message
	fe := &FieldError{Name: "email", Message: "Custom error message"}
	if fe.Error() != "Custom error message" {
		t.Errorf("Expected 'Custom error message', got %s", fe.Error())
	}

	// Test with function error
	fe = &FieldError{Name: "email", FuncError: errors.New("function error")}
	expected := "validation failed for field 'email': function error"
	if fe.Error() != expected {
		t.Errorf("Expected '%s', got %s", expected, fe.Error())
	}

	// Test with no message or function error
	fe = &FieldError{Name: "email"}
	expected = "validation failed for field 'email'"
	if fe.Error() != expected {
		t.Errorf("Expected '%s', got %s", expected, fe.Error())
	}
}

func TestFieldErrorUnwrap(t *testing.T) {
	originalErr := errors.New("original error")
	fe := &FieldError{Name: "email", FuncError: originalErr}

	if fe.Unwrap() != originalErr {
		t.Error("Unwrap should return the original function error")
	}

	fe = &FieldError{Name: "email"}
	if fe.Unwrap() != nil {
		t.Error("Unwrap should return nil when no function error")
	}
}

func TestFieldErrorHasFuncError(t *testing.T) {
	fe := &FieldError{Name: "email", FuncError: errors.New("error")}
	if !fe.HasFuncError() {
		t.Error("Expected HasFuncError to return true")
	}

	fe = &FieldError{Name: "email"}
	if fe.HasFuncError() {
		t.Error("Expected HasFuncError to return false")
	}
}

func TestFieldErrorSetMessage(t *testing.T) {
	fe := &FieldError{Name: "email"}
	fe.SetMessage("New message")

	if fe.Message != "New message" {
		t.Errorf("Expected 'New message', got %s", fe.Message)
	}
}
