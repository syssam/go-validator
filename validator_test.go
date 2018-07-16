package validator

import (
	"testing"
)

type FieldsRequired struct {
	Name  string ``
	Email string `valid:"required"`
}

type MultipleFieldsRequired struct {
	Url   string `valid:"required"`
	Email string `valid:"required"`
}

type FieldsEmail struct {
	Email string `valid:"email"`
}

func TestFieldsRequired(t *testing.T) {
	var tests = []struct {
		param    FieldsRequired
		expected bool
	}{
		{FieldsRequired{}, false},
		{FieldsRequired{Name: "", Email: ""}, false},
		{FieldsRequired{Name: "TEST"}, false},
		{FieldsRequired{Name: "TEST", Email: ""}, false},
		{FieldsRequired{Email: "test@example.com"}, true},
		{FieldsRequired{Name: "", Email: "test@example.com"}, true},
		{FieldsRequired{Name: "TEST", Email: "test@example.com"}, true},
	}
	for _, test := range tests {
		err := ValidateStruct(test.param)
		actual := err == nil
		if actual != test.expected {
			t.Errorf("Expected validateateStruct(%q) to be %v, got %v", test.param, test.expected, actual)
			if err == nil {
				t.Errorf("Got Error on validateateStruct(%q): %s", test.param, err)
			}
		}
	}
}

func TestMultipleFieldsRequired(t *testing.T) {
	var tests = []struct {
		param    MultipleFieldsRequired
		expected bool
	}{
		{MultipleFieldsRequired{}, false},
		{MultipleFieldsRequired{Url: "", Email: ""}, false},
		{MultipleFieldsRequired{Url: "TEST"}, false},
		{MultipleFieldsRequired{Url: "TEST", Email: ""}, false},
		{MultipleFieldsRequired{Email: "test@example.com"}, false},
		{MultipleFieldsRequired{Url: "", Email: "test@example.com"}, false},
		{MultipleFieldsRequired{Url: "TEST", Email: "test@example.com"}, true},
	}
	for _, test := range tests {
		err := ValidateStruct(&test.param)
		actual := err == nil
		if actual != test.expected {
			t.Errorf("Expected validateateStruct(%q) to be %v, got %v", test.param, test.expected, actual)
			if err != nil {
				t.Errorf("Got Error on validateateStruct(%q): %s", test.param, err)
			}
		}
	}
}

func TestRequiredIf(t *testing.T) {
	type RequiredIf struct {
		First string
		Last  string `valid:"requiredIf=First|taylor"`
	}

	var tests = []struct {
		param    RequiredIf
		expected bool
	}{
		{RequiredIf{}, true},
		{RequiredIf{First: "", Last: ""}, true},
		{RequiredIf{First: "TEST"}, true},
		{RequiredIf{First: "taylor", Last: ""}, false},
		{RequiredIf{First: "taylor", Last: "otwell"}, true},
	}
	for _, test := range tests {
		err := ValidateStruct(&test.param)
		actual := err == nil
		if actual != test.expected {
			t.Errorf("Expected validateateStruct(%q) to be %v, got %v", test.param, test.expected, actual)
			if err != nil {
				t.Errorf("Got Error on validateateStruct(%q): %s", test.param, err)
			}
		}
	}
}

func TestMultipleRequiredIf(t *testing.T) {
	type RequiredIf struct {
		First string
		Last  string `valid:"requiredIf=First|taylor|dayle"`
	}

	var tests = []struct {
		param    RequiredIf
		expected bool
	}{
		{RequiredIf{}, true},
		{RequiredIf{First: "", Last: ""}, true},
		{RequiredIf{First: "TEST"}, true},
		{RequiredIf{First: "taylor", Last: ""}, false},
		{RequiredIf{First: "dayle", Last: ""}, false},
		{RequiredIf{First: "taylor", Last: "otwell"}, true},
		{RequiredIf{First: "dayle", Last: "otwell"}, true},
	}
	for _, test := range tests {
		err := ValidateStruct(&test.param)
		actual := err == nil
		if actual != test.expected {
			t.Errorf("Expected validateateStruct(%q) to be %v, got %v", test.param, test.expected, actual)
			if err != nil {
				t.Errorf("Got Error on validateateStruct(%q): %s", test.param, err)
			}
		}
	}
}

func TestFieldsEmail(t *testing.T) {
	var tests = []struct {
		param    FieldsEmail
		expected bool
	}{
		{FieldsEmail{}, false},
		{FieldsEmail{Email: ""}, false},
		{FieldsEmail{Email: "aaa"}, false},
		{FieldsEmail{Email: "test@example.com"}, true},
	}
	for _, test := range tests {
		err := ValidateStruct(test.param)
		actual := err == nil
		if actual != test.expected {
			t.Errorf("Expected validateateStruct(%q) to be %v, got %v", test.param, test.expected, actual)
			if err != nil {
				t.Errorf("Got Error on validateateStruct(%q): %s", test.param, err)
			}
		}
	}
}
