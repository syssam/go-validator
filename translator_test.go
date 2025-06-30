package validator

import (
	"fmt"
	"testing"
)

func TestNewTranslator(t *testing.T) {
	translator := NewTranslator()
	if translator == nil {
		t.Error("Expected NewTranslator to return a non-nil translator")
		return
	}
	if translator.messages == nil {
		t.Error("Expected translator.messages to be initialized")
	}
	if translator.attributes == nil {
		t.Error("Expected translator.attributes to be initialized")
	}
	if translator.customMessage == nil {
		t.Error("Expected translator.customMessage to be initialized")
	}
}

func TestSetMessage(t *testing.T) {
	translator := NewTranslator()
	messages := Translate{"required": "Field is required"}
	translator.SetMessage("en", messages)

	if translator.messages["en"] == nil {
		t.Error("Expected messages to be set for 'en' language")
	}
	if translator.messages["en"]["required"] != "Field is required" {
		t.Error("Expected message to be set correctly")
	}
}

func TestLoadMessage(t *testing.T) {
	translator := NewTranslator()
	messages := Translate{"required": "Field is required"}
	translator.SetMessage("en", messages)

	loaded := translator.LoadMessage("en")
	if loaded["required"] != "Field is required" {
		t.Error("Expected loaded message to match set message")
	}

	// Test loading non-existent language
	empty := translator.LoadMessage("fr")
	if empty != nil {
		t.Error("Expected nil for non-existent language")
	}
}

func TestSetAttributes(t *testing.T) {
	translator := NewTranslator()
	attributes := Translate{"User.Name": "Full Name"}
	translator.SetAttributes("en", attributes)

	if translator.attributes["en"] == nil {
		t.Error("Expected attributes to be set for 'en' language")
	}
	if translator.attributes["en"]["User.Name"] != "Full Name" {
		t.Error("Expected attribute to be set correctly")
	}
}

func TestTrans(t *testing.T) {
	translator := NewTranslator()

	// Set up messages
	messages := Translate{
		"required": "The {{.Attribute}} field is required",
		"email":    "The {{.Attribute}} field must be a valid email",
	}
	translator.SetMessage("en", messages)

	// Set up attributes
	attributes := Translate{"User.Email": "Email Address"}
	translator.SetAttributes("en", attributes)

	// Create test field error
	fieldError := &FieldError{
		Name:        "email",
		StructName:  "User.Email",
		MessageName: "required",
		Attribute:   "email",
	}

	errors := Errors{fieldError}
	translatedErrors := translator.Trans(errors, "en")

	if len(translatedErrors) != 1 {
		t.Error("Expected one translated error")
	}

	translated := translatedErrors[0].(*FieldError)
	if translated.Message != "The Email Address field is required" {
		t.Errorf("Expected 'The Email Address field is required', got '%s'", translated.Message)
	}
}

func TestTransWithCustomMessage(t *testing.T) {
	translator := NewTranslator()

	// Set up custom message
	customMessage := Translate{"email.required": "Email is mandatory"}
	translator.customMessage["en"] = customMessage

	// Create test field error
	fieldError := &FieldError{
		Name:        "email",
		MessageName: "required",
		Attribute:   "email",
	}

	errors := Errors{fieldError}
	translatedErrors := translator.Trans(errors, "en")

	translated := translatedErrors[0].(*FieldError)
	if translated.Message != "Email is mandatory" {
		t.Errorf("Expected 'Email is mandatory', got '%s'", translated.Message)
	}
}

func TestTransWithMessageParameters(t *testing.T) {
	translator := NewTranslator()

	// Set up messages with parameters
	messages := Translate{
		"between": "The {{.Attribute}} field must be between {{.Min}} and {{.Max}}",
	}
	translator.SetMessage("en", messages)

	// Create test field error with parameters
	fieldError := &FieldError{
		Name:        "age",
		MessageName: "between",
		Attribute:   "age",
		MessageParameters: MessageParameters{
			{Key: "Min", Value: "18"},
			{Key: "Max", Value: "65"},
		},
	}

	errors := Errors{fieldError}
	translatedErrors := translator.Trans(errors, "en")

	translated := translatedErrors[0].(*FieldError)
	expected := "The age field must be between 18 and 65"
	if translated.Message != expected {
		t.Errorf("Expected '%s', got '%s'", expected, translated.Message)
	}
}

func TestTransWithDefaultAttribute(t *testing.T) {
	translator := NewTranslator()

	// Set up messages
	messages := Translate{"required": "The {{.Attribute}} field is required"}
	translator.SetMessage("en", messages)

	// Create test field error with default attribute
	fieldError := &FieldError{
		Name:             "user_name",
		MessageName:      "required",
		Attribute:        "user_name",
		DefaultAttribute: "User Name",
	}

	errors := Errors{fieldError}
	translatedErrors := translator.Trans(errors, "en")

	translated := translatedErrors[0].(*FieldError)
	if translated.Message != "The User Name field is required" {
		t.Errorf("Expected 'The User Name field is required', got '%s'", translated.Message)
	}
}

func TestTransWithNonFieldError(t *testing.T) {
	translator := NewTranslator()

	// Create error that's not a FieldError
	genericError := &struct{ error }{fmt.Errorf("generic error")}

	errors := Errors{genericError}
	translatedErrors := translator.Trans(errors, "en")

	// Should not panic and should return original errors
	if len(translatedErrors) != 1 {
		t.Error("Expected one error")
	}
}
