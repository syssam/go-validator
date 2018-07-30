package validator

import (
	"testing"
)

type FieldsRequired struct {
	Name  string ``
	Email string `valid:"required"`
}

type MultipleFieldsRequired struct {
	URL   string `valid:"required"`
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
		{MultipleFieldsRequired{URL: "", Email: ""}, false},
		{MultipleFieldsRequired{URL: "TEST"}, false},
		{MultipleFieldsRequired{URL: "TEST", Email: ""}, false},
		{MultipleFieldsRequired{Email: "test@example.com"}, false},
		{MultipleFieldsRequired{URL: "", Email: "test@example.com"}, false},
		{MultipleFieldsRequired{URL: "TEST", Email: "test@example.com"}, true},
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

	type RequiredUnlessString struct {
		First string `valid:"requiredUnless=Last|otwell"`
		Last  string
	}

	type RequiredUnlessBool struct {
		First string `valid:"requiredUnless=Last|true"`
		Last  bool
	}

	type RequiredUnlessNumber struct {
		First string `valid:"requiredUnless=Last|888"`
		Last  int
	}

	type RequiredUnlessSub struct {
		First string `valid:"requiredUnless=Last.Test|otwell"`
		Last  *SubTest
	}

	var tests = []struct {
		param    interface{}
		expected bool
	}{
		{RequiredUnlessString{}, false},
		{RequiredUnlessString{Last: "aa"}, false},
		{RequiredUnlessString{Last: "otwell"}, true},
		{RequiredUnlessBool{Last: false}, false},
		{RequiredUnlessBool{Last: true}, true},
		{RequiredUnlessNumber{Last: 555}, false},
		{RequiredUnlessNumber{Last: 888}, true},
		{RequiredUnlessSub{Last: &SubTest{Test: "test"}}, false},
		{RequiredUnlessSub{Last: &SubTest{Test: "otwell"}}, true},
	}
	for i, test := range tests {
		var err error
		switch test.param.(type) {
		case RequiredUnlessString:
			err = ValidateStruct(test.param.(RequiredUnlessString))
		case RequiredUnlessBool:
			err = ValidateStruct(test.param.(RequiredUnlessBool))
		case RequiredUnlessNumber:
			err = ValidateStruct(test.param.(RequiredUnlessNumber))
		case RequiredUnlessSub:
			err = ValidateStruct(test.param.(RequiredUnlessSub))
		}
		actual := err == nil
		if actual != test.expected {
			t.Errorf("Expected validateateStruct(%T) Case %d to be %v, got %v", test.param, i, test.expected, actual)
			if err != nil {
				t.Errorf("Got Error on validateateStruct(%T): %s", test.param, err)
			}
		}
	}
}

func TestMin(t *testing.T) {
	type TestMinStruct struct {
		String string            `valid:"min=5"`
		Int    int               `valid:"min=5"`
		Unit   uint              `valid:"min=5"`
		Float  float64           `valid:"min=5.3"`
		Array  []string          `valid:"min=5"`
		Map    map[string]string `valid:"min=5"`
	}
	var tests = []struct {
		param    TestMinStruct
		expected bool
	}{
		{TestMinStruct{}, true},
		{TestMinStruct{String: "Hell"}, false},
		{TestMinStruct{String: "Hello"}, true},
		{TestMinStruct{Int: 4}, false},
		{TestMinStruct{Int: 5}, true},
		{TestMinStruct{Unit: 4}, false},
		{TestMinStruct{Unit: 5}, true},
		{TestMinStruct{Float: 5.2}, false},
		{TestMinStruct{Float: 5.9}, true},
		{TestMinStruct{Array: []string{"1", "2", "3", "4"}}, false},
		{TestMinStruct{Array: []string{"1", "2", "3", "4", "5"}}, true},
		{TestMinStruct{Map: map[string]string{
			"rsc": "string",
			"r":   "string",
			"gri": "string",
			"adg": "string",
		}}, false},
		{TestMinStruct{Map: map[string]string{
			"rsc": "string",
			"r":   "string",
			"gri": "string",
			"adg": "string",
			"tt":  "string",
		}}, true},
	}

	for i, test := range tests {
		err := ValidateStruct(test.param)
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
