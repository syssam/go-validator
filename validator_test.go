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
	type SubTest struct {
		Test string
	}

	type RequiredIf struct {
		String           string `valid:"requiredIf=RequiredIfString|otwell"`
		RequiredIfString string
		Bool             string `valid:"requiredIf=RequiredIfBool|true"`
		RequiredIfBool   bool
		Number           string `valid:"requiredIf=RequiredIfNumber|888"`
		RequiredIfNumber int
		Array            string `valid:"requiredIf=RequiredIfArray|888"`
		RequiredIfArray  []string
		SubTest          string `valid:"requiredIf=RequiredIfSub.Test|otwell"`
		RequiredIfSub    *SubTest
	}
	var tests = []struct {
		param    RequiredIf
		expected bool
	}{
		{RequiredIf{}, true},
		{RequiredIf{String: "", RequiredIfString: ""}, true},
		{RequiredIf{String: "String"}, true},
		{RequiredIf{RequiredIfString: "otwell"}, false},
		{RequiredIf{String: "", RequiredIfString: "otwell"}, false},
		{RequiredIf{String: "String", RequiredIfString: "otwell"}, true},
		{RequiredIf{Bool: "Bool"}, true},
		{RequiredIf{Bool: "", RequiredIfBool: false}, true},
		{RequiredIf{Bool: "", RequiredIfBool: true}, false},
		{RequiredIf{Number: "Number"}, true},
		{RequiredIf{Number: "", RequiredIfNumber: 555}, true},
		{RequiredIf{Number: "", RequiredIfNumber: 888}, false},
		{RequiredIf{SubTest: "SubTest"}, true},
		{RequiredIf{SubTest: "", RequiredIfSub: &SubTest{Test: "aaa"}}, true},
		{RequiredIf{SubTest: "", RequiredIfSub: &SubTest{Test: "otwell"}}, false},
	}
	for i, test := range tests {
		err := ValidateStruct(&test.param)
		actual := err == nil
		if actual != test.expected {
			t.Errorf("Expected validateateStruct(%T) Case %d to be %v, got %v", test.param, i, test.expected, actual)
			if err != nil {
				t.Errorf("Got Error on validateateStruct(%T): %s", test.param, err)
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

func TestRequiredUnless(t *testing.T) {
	type SubTest struct {
		Test string
	}

	type RequiredUnless struct {
		String               string `valid:"RequiredUnless=RequiredUnlessString|otwell"`
		RequiredUnlessString string
		Bool                 string `valid:"RequiredUnless=RequiredUnlessBool|true"`
		RequiredUnlessBool   bool
		Number               string `valid:"RequiredUnless=RequiredUnlessNumber|888"`
		RequiredUnlessNumber int
		Array                string `valid:"RequiredUnless=RequiredUnlessArray|888"`
		RequiredUnlessArray  []string
		SubTest              string `valid:"RequiredUnless=RequiredUnlessSub.Test|otwell"`
		RequiredUnlessSub    *SubTest
	}
	var tests = []struct {
		param    RequiredUnless
		expected bool
	}{
		{RequiredUnless{}, true},
		{RequiredUnless{String: "", RequiredUnlessString: ""}, true},
		{RequiredUnless{String: "String"}, true},
		{RequiredUnless{RequiredUnlessString: "otwell"}, false},
		{RequiredUnless{String: "", RequiredUnlessString: "otwell"}, false},
		{RequiredUnless{String: "String", RequiredUnlessString: "otwell"}, true},
		{RequiredUnless{Bool: "Bool"}, true},
		{RequiredUnless{Bool: "", RequiredUnlessBool: false}, true},
		{RequiredUnless{Bool: "", RequiredUnlessBool: true}, false},
		{RequiredUnless{Number: "Number"}, true},
		{RequiredUnless{Number: "", RequiredUnlessNumber: 555}, true},
		{RequiredUnless{Number: "", RequiredUnlessNumber: 888}, false},
		{RequiredUnless{SubTest: "SubTest"}, true},
		{RequiredUnless{SubTest: "", RequiredUnlessSub: &SubTest{Test: "aaa"}}, true},
		{RequiredUnless{SubTest: "", RequiredUnlessSub: &SubTest{Test: "otwell"}}, false},
	}
	for i, test := range tests {
		err := ValidateStruct(&test.param)
		actual := err == nil
		if actual != test.expected {
			t.Errorf("Expected validateateStruct(%T) Case %d to be %v, got %v", test.param, i, test.expected, actual)
			if err != nil {
				t.Errorf("Got Error on validateateStruct(%T): %s", test.param, err)
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
