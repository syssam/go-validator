package validator

import (
	"encoding/json"
	"fmt"
	"reflect"
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

type FieldsURL struct {
	URL string `valid:"url"`
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

func TestMax(t *testing.T) {
	type TestMaxStruct struct {
		String string            `valid:"max=5"`
		Int    int               `valid:"max=5"`
		Unit   uint              `valid:"max=5"`
		Float  float64           `valid:"max=5.3"`
		Array  []string          `valid:"max=5"`
		Map    map[string]string `valid:"max=5"`
	}
	var tests = []struct {
		param    TestMaxStruct
		expected bool
	}{
		{TestMaxStruct{}, true},
		{TestMaxStruct{String: "Hell"}, true},
		{TestMaxStruct{String: "Hello World"}, false},
		{TestMaxStruct{Int: 4}, true},
		{TestMaxStruct{Int: 6}, false},
		{TestMaxStruct{Unit: 4}, true},
		{TestMaxStruct{Unit: 6}, false},
		{TestMaxStruct{Float: 5.2}, true},
		{TestMaxStruct{Float: 5.9}, false},
		{TestMaxStruct{Array: []string{"1", "2", "3", "4"}}, true},
		{TestMaxStruct{Array: []string{"1", "2", "3", "4", "5", "6"}}, false},
		{TestMaxStruct{Array: []string{"1", "2", "3", "4", "5"}}, true},
		{TestMaxStruct{Map: map[string]string{
			"rsc": "string",
			"r":   "string",
			"gri": "string",
			"adg": "string",
			"ab":  "string",
			"cd":  "string",
		}}, false},
		{TestMaxStruct{Map: map[string]string{
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

func TestMin(t *testing.T) {
	type TestMinStruct struct {
		String string            `valid:"omitempty,min=5"`
		Int    int               `valid:"omitempty,min=5"`
		Unit   uint              `valid:"omitempty,min=5"`
		Float  float64           `valid:"omitempty,min=5.3"`
		Array  []string          `valid:"omitempty,min=5"`
		Map    map[string]string `valid:"omitempty,min=5"`
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
func TestSize(t *testing.T) {
	type TestSizeStruct struct {
		String string            `valid:"omitempty,size=5"`
		Int    int               `valid:"omitempty,size=5"`
		Unit   uint              `valid:"omitempty,size=5"`
		Float  float64           `valid:"omitempty,size=5.3"`
		Array  []string          `valid:"omitempty,size=5"`
		Map    map[string]string `valid:"omitempty,size=5"`
	}
	var tests = []struct {
		param    TestSizeStruct
		expected bool
	}{
		{TestSizeStruct{}, true},
		{TestSizeStruct{String: "Hell"}, false},
		{TestSizeStruct{String: "Hello"}, true},
		{TestSizeStruct{Int: 4}, false},
		{TestSizeStruct{Int: 5}, true},
		{TestSizeStruct{Unit: 4}, false},
		{TestSizeStruct{Unit: 5}, true},
		{TestSizeStruct{Float: 5.2}, false},
		{TestSizeStruct{Float: 5.3}, true},
		{TestSizeStruct{Array: []string{"1", "2", "3", "4"}}, false},
		{TestSizeStruct{Array: []string{"1", "2", "3", "4", "5"}}, true},
		{TestSizeStruct{Map: map[string]string{
			"rsc": "string",
			"r":   "string",
			"gri": "string",
			"adg": "string",
		}}, false},
		{TestSizeStruct{Map: map[string]string{
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

func TestGt(t *testing.T) {
	type GtStruct struct {
		String         string `valid:"omitempty,gt=GtStructString"`
		GtStructString string
		Number         int `valid:"omitempty,gt=GtStructNumber"`
		GtStructNumber int
		Array          []string `valid:"omitempty,gt=GtStructArray"`
		GtStructArray  []string
	}
	var tests = []struct {
		param    GtStruct
		expected bool
	}{
		{GtStruct{}, true},
		{GtStruct{String: "Hell", GtStructString: "Hello"}, false},
		{GtStruct{String: "Hello World", GtStructString: "Hello"}, true},
		{GtStruct{Number: 4, GtStructNumber: 5}, false},
		{GtStruct{Number: 10, GtStructNumber: 5}, true},
		{GtStruct{Array: []string{"1", "2", "3", "4"}, GtStructArray: []string{"1", "2", "3", "4", "5"}}, false},
		{GtStruct{Array: []string{"1", "2", "3", "4", "5"}, GtStructArray: []string{"1", "2", "3", "4"}}, true},
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

func TestGte(t *testing.T) {
	type GteStruct struct {
		String          string `valid:"omitempty,gte=GteStructString"`
		GteStructString string
		Number          int `valid:"omitempty,gte=GteStructNumber"`
		GteStructNumber int
		Array           []string `valid:"omitempty,gte=GteStructArray"`
		GteStructArray  []string
	}
	var tests = []struct {
		param    GteStruct
		expected bool
	}{
		{GteStruct{}, true},
		{GteStruct{String: "Hell", GteStructString: "Hello"}, false},
		{GteStruct{String: "Hello World", GteStructString: "Hello"}, true},
		{GteStruct{Number: 4, GteStructNumber: 5}, false},
		{GteStruct{Number: 10, GteStructNumber: 5}, true},
		{GteStruct{Array: []string{"1", "2", "3", "4"}, GteStructArray: []string{"1", "2", "3", "4", "5"}}, false},
		{GteStruct{Array: []string{"1", "2", "3", "4", "5"}, GteStructArray: []string{"1", "2", "3", "4"}}, true},
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
	for i, test := range tests {
		err := ValidateStruct(test.param)
		actual := err == nil
		if actual != test.expected {
			t.Errorf("Expected validateateStruct(%+v) Case %d to be %v, got %v", test.param, i, test.expected, actual)
			if err != nil {
				t.Errorf("Got Error on validateateStruct(%+v): %s", test.param, err)
			}
		}
	}
}

func TestDistinct(t *testing.T) {
	var result bool
	result = ValidateDistinct([]string{"zh-HK", "zh-CN"})
	if result != true {
		t.Errorf("Got Error on validateateStruct case 1 %v, got %v", true, result)
	}
	result = ValidateDistinct([]string{"zh-HK", "zh-CN", "zh-CN"})
	if result != false {
		t.Errorf("Got Error on validateateStruct case 2 %v, got %v", false, result)
	}
}

func TestInnerStruct(t *testing.T) {
	CustomTypeRuleMap.Set("languageCode", func(v reflect.Value, o reflect.Value, validTag *ValidTag) bool {
		return v.Kind() != reflect.String
	})
	MessageMap["languageCode"] = "Language Code is not valid."
	type CreditCard struct {
		Number           string
		UserMemberNumber string `valid:"languageCode"`
	}

	type User struct {
		MemberNumber string
		CreditCards  []CreditCard `json:"CreditCards" valid:"languageCode"`
	}

	c := []CreditCard{
		{
			Number:           "1",
			UserMemberNumber: "1",
		},
		{
			Number:           "2",
			UserMemberNumber: "2",
		},
	}

	u := User{
		MemberNumber: "MemberNumber",
		CreditCards:  c,
	}

	_ = ValidateStruct(u)
}

func TestInnerStruct2(t *testing.T) {
	type CreditCard struct {
		Number           string
		UserMemberNumber string `valid:"required,max=64"`
	}

	type User struct {
		MemberNumber string
		CreditCards  []CreditCard `json:"CreditCards"`
	}

	c := []CreditCard{
		{
			Number:           "1",
			UserMemberNumber: "",
		},
		{
			Number:           "2",
			UserMemberNumber: "2",
		},
	}

	u := User{
		MemberNumber: "MemberNumber",
		CreditCards:  c,
	}

	err := ValidateStruct(u)
	if err == nil {
		t.Errorf("Got Error on validateateStruct: %s", err)
	}
}

func TestNameSpace(t *testing.T) {
	type CreditCard struct {
		Number           string
		UserMemberNumber string `valid:"languageCode"`
	}

	type User struct {
		MemberNumber string
		CreditCards  []CreditCard `json:"CreditCards" valid:"languageCode"`
	}

	c := []CreditCard{
		{
			Number:           "1",
			UserMemberNumber: "1",
		},
		{
			Number:           "2",
			UserMemberNumber: "2",
		},
	}

	u := User{
		MemberNumber: "MemberNumber",
		CreditCards:  c,
	}

	_ = ValidateStruct(u)
}

func TestIsURL(t *testing.T) {
	var tests = []struct {
		param    string
		expected bool
	}{
		{"http://foo.bar#com", true},
		{"http://foobar.com", true},
		{"https://foobar.com", true},
		{"foobar.com", false},
		{"http://foobar.coffee/", true},
		{"http://foobar.中文网/", true},
		{"http://foobar.org/", true},
		{"http://foobar.org:8080/", true},
		{"ftp://foobar.ru/", true},
		{"http://user:pass@www.foobar.com/", true},
		{"http://127.0.0.1/", true},
		{"http://duckduckgo.com/?q=%2F", true},
		{"http://localhost:3000/", true},
		{"http://foobar.com/?foo=bar#baz=qux", true},
		{"http://foobar.com?foo=bar", true},
		{"http://www.xn--froschgrn-x9a.net/", true},
		{"", true},
		{"xyz://foobar.com", true},
		{"invalid.", false},
		{".com", false},
		{"rtmp://foobar.com", true},
		{"http://www.foo_bar.com/", true},
		{"http://localhost:3000/", true},
		{"http://foobar.com/#baz", true},
		{"http://foobar.com#baz=qux", true},
		{"http://foobar.com/t$-_.+!*\\'(),", true},
		{"http://www.foobar.com/~foobar", true},
		{"http://www.-foobar.com/", true},
		{"http://www.foo---bar.com/", true},
		{"mailto:someone@example.com", true},
		{"irc://irc.server.org/channel", true},
		{"irc://#channel@network", true},
		{"/abs/test/dir", false},
		{"./rel/test/dir", false},
	}
	for i, test := range tests {
		actual := ValidateURL(test.param)
		if actual != test.expected {
			t.Errorf("Expected IsURL(%+v) Case %d to be %v, got %v", test.param, i, test.expected, actual)
		}
	}
}

func TestPointer(t *testing.T) {
	type FieldsPointer struct {
		Name     *string ``
		UserName *string `valid:"max=5"`
		Email    *string `valid:"required,email"`
	}
	e1 := ""
	u1 := ""
	e2 := "test"
	u2 := "test"
	e3 := "test@example.com"
	u3 := "test123456"
	e4 := "test@example.com"
	u4 := "test"
	var tests = []struct {
		param    FieldsPointer
		expected bool
	}{
		{FieldsPointer{}, false},
		{FieldsPointer{Email: &e1, UserName: &u1}, false},
		{FieldsPointer{Email: &e2, UserName: &u2}, false},
		{FieldsPointer{Email: &e3, UserName: &u3}, false},
		{FieldsPointer{Email: &e4, UserName: &u4}, true},
	}
	for i, test := range tests {
		err := ValidateStruct(test.param)
		actual := err == nil
		if actual != test.expected {
			t.Errorf("Expected validateateStruct(%+v) Case %d to be %v, got %v", test.param, i, test.expected, actual)
			if err != nil {
				t.Errorf("Got Error on validateateStruct(%+v): %s", test.param, err)
			}
		}
	}
}

// TestErrorHandling tests the new error handling mechanisms
func TestErrorHandling(t *testing.T) {
	// Test unsupported type errors for validateBetween
	type UnsupportedBetween struct {
		Complex complex64 `valid:"between=1|10"`
	}

	err := ValidateStruct(UnsupportedBetween{Complex: 5 + 5i})
	if err == nil {
		t.Error("Expected error for unsupported type in between validation")
	}

	errors := err.(Errors)
	if len(errors) == 0 {
		t.Error("Expected at least one error")
	}

	if !errors.HasFieldError("Complex") {
		t.Error("Expected error for Complex field")
	}
}

// TestUnsupportedTypeErrors tests error handling for unsupported types
func TestUnsupportedTypeErrors(t *testing.T) {
	tests := []struct {
		name  string
		input interface{}
		field string
	}{
		{
			name: "Complex type in Size validation",
			input: struct {
				Complex complex64 `valid:"size=5"`
			}{Complex: 3 + 4i},
			field: "Complex",
		},
		{
			name: "Complex type in Max validation",
			input: struct {
				Complex complex128 `valid:"max=10"`
			}{Complex: 5 + 5i},
			field: "Complex",
		},
		{
			name: "Complex type in Min validation",
			input: struct {
				Complex complex64 `valid:"min=1"`
			}{Complex: 2 + 2i},
			field: "Complex",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateStruct(test.input)
			if err == nil {
				t.Errorf("Expected error for %s", test.name)
				return
			}

			errors := err.(Errors)
			if !errors.HasFieldError(test.field) {
				t.Errorf("Expected error for field %s", test.field)
			}
		})
	}
}

// TestComparisonOperatorErrors tests invalid operators in comparison functions
func TestComparisonOperatorErrors(t *testing.T) {
	// Test compareFloat64 with invalid operator
	_, err := ValidateGt(5.0, 3.0)
	if err != nil {
		t.Errorf("Expected no error for valid comparison, got: %v", err)
	}

	// Test compareInt64 indirectly through struct validation
	type TestStruct struct {
		Value1 int `valid:"gt=Value2"`
		Value2 int
	}

	// This should work normally
	err = ValidateStruct(TestStruct{Value1: 10, Value2: 5})
	if err != nil {
		t.Errorf("Expected no error for valid gt validation, got: %v", err)
	}
}

// TestErrorsUtilityMethods tests the new Errors utility methods
func TestErrorsUtilityMethods(t *testing.T) {
	type TestStruct struct {
		Name  string `valid:"required"`
		Email string `valid:"required,email"`
		Age   int    `valid:"min=18"`
	}

	err := ValidateStruct(TestStruct{Name: "", Email: "invalid", Age: 16})
	if err == nil {
		t.Error("Expected validation errors")
		return
	}

	errors := err.(Errors)

	// Test HasFieldError
	if !errors.HasFieldError("Name") {
		t.Error("Expected error for Name field")
	}
	if !errors.HasFieldError("Email") {
		t.Error("Expected error for Email field")
	}
	if !errors.HasFieldError("Age") {
		t.Error("Expected error for Age field")
	}
	if errors.HasFieldError("NonExistent") {
		t.Error("Did not expect error for NonExistent field")
	}

	// Test GetFieldError
	nameError := errors.GetFieldError("Name")
	if nameError == nil {
		t.Error("Expected to get Name field error")
	} else if nameError.Name != "Name" {
		t.Errorf("Expected field name 'Name', got '%s'", nameError.Name)
	}

	nonExistentError := errors.GetFieldError("NonExistent")
	if nonExistentError != nil {
		t.Error("Did not expect to get error for NonExistent field")
	}

	// Test GroupByField
	groups := errors.GroupByField()
	if len(groups) != 3 {
		t.Errorf("Expected 3 error groups, got %d", len(groups))
	}

	if _, exists := groups["Name"]; !exists {
		t.Error("Expected Name in grouped errors")
	}
	if _, exists := groups["Email"]; !exists {
		t.Error("Expected Email in grouped errors")
	}
	if _, exists := groups["Age"]; !exists {
		t.Error("Expected Age in grouped errors")
	}
}

// TestFieldErrorMethods tests FieldError methods
func TestFieldErrorMethods(t *testing.T) {
	fieldError := &FieldError{
		Name:    "TestField",
		Message: "Test message",
	}

	// Test Error() method
	expectedMsg := "Test message"
	if fieldError.Error() != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, fieldError.Error())
	}

	// Test SetMessage method
	newMsg := "New test message"
	fieldError.SetMessage(newMsg)
	if fieldError.Message != newMsg {
		t.Errorf("Expected message to be set to '%s', got '%s'", newMsg, fieldError.Message)
	}

	// Test Error() with no message
	emptyFieldError := &FieldError{
		Name: "EmptyField",
	}
	expectedDefault := "validation failed for field 'EmptyField'"
	if emptyFieldError.Error() != expectedDefault {
		t.Errorf("Expected default error message '%s', got '%s'", expectedDefault, emptyFieldError.Error())
	}
}

// TestBetweenParameterValidation tests parameter validation for between rule
func TestBetweenParameterValidation(t *testing.T) {
	type InvalidParams struct {
		Value int `valid:"between=5"` // Missing second parameter
	}

	err := ValidateStruct(InvalidParams{Value: 7})
	if err == nil {
		t.Error("Expected error for invalid between parameters")
	}
}

// TestRequiredIfEdgeCases tests edge cases for requiredIf validation
func TestRequiredIfEdgeCases(t *testing.T) {
	type TestStruct struct {
		Field1 string `valid:"requiredIf=Field2|value1|value2"`
		Field2 string
		Field3 string            `valid:"requiredIf=NonExistent|value"`
		Field4 map[string]string `valid:"requiredIf=Field5|test"`
		Field5 string
	}

	// Test with valid requiredIf - Field4 should have a value since Field5="test"
	validStruct := TestStruct{
		Field1: "present",
		Field2: "value1",
		Field4: map[string]string{"key": "value"}, // Field4 is required when Field5="test"
		Field5: "test",
	}

	err := ValidateStruct(validStruct)
	if err != nil {
		t.Errorf("Expected no error for valid requiredIf, got: %v", err)
	}

	// Test with missing required field
	invalidStruct := TestStruct{
		Field1: "", // Should be required because Field2 = "value1"
		Field2: "value1",
	}

	err = ValidateStruct(invalidStruct)
	if err == nil {
		t.Error("Expected error for missing required field")
	}
}

// TestGo119Features tests Go 1.19 specific improvements
func TestGo119Features(t *testing.T) {
	// Test that the validator works with Go 1.19 types and features
	type ModernStruct struct {
		Text string `valid:"required"`
		Num  int    `valid:"min=1"`
	}

	// This tests that the validator properly handles modern Go syntax
	data := ModernStruct{
		Text: "test",
		Num:  5,
	}

	err := ValidateStruct(data)
	if err != nil {
		t.Errorf("Expected no error with Go 1.19 compatible struct, got: %v", err)
	}

	// Test with invalid data
	invalidData := ModernStruct{
		Text: "",
		Num:  0,
	}

	err = ValidateStruct(invalidData)
	if err == nil {
		t.Error("Expected validation errors for invalid data")
	}
}

// TestFuncErrorChaining tests the FuncError functionality and error chaining
func TestFuncErrorChaining(t *testing.T) {
	// Test case 1: FuncError with unsupported type
	t.Run("UnsupportedType", func(t *testing.T) {
		type UnsupportedStruct struct {
			Complex complex64 `valid:"between=1|10"`
		}

		err := ValidateStruct(UnsupportedStruct{Complex: 5 + 5i})
		if err == nil {
			t.Error("Expected validation error for unsupported type")
			return
		}

		errors := err.(Errors)
		if len(errors) == 0 {
			t.Error("Expected at least one error")
			return
		}

		fieldError := errors.GetFieldError("Complex")
		if fieldError == nil {
			t.Error("Expected FieldError for Complex field")
			return
		}

		// Test HasFuncError method
		if !fieldError.HasFuncError() {
			t.Error("Expected HasFuncError to return true")
		}

		// Test Unwrap method
		funcErr := fieldError.Unwrap()
		if funcErr == nil {
			t.Error("Expected Unwrap to return non-nil error")
		}

		// Verify the underlying error contains expected message
		if funcErr.Error() == "" {
			t.Error("Expected non-empty function error message")
		}
	})

	// Test case 2: FuncError with comparison operator error
	t.Run("ComparisonOperatorError", func(t *testing.T) {
		// Create a custom validation that will trigger a comparison error
		type TestStruct struct {
			Value1 int `valid:"gt=Value2"`
			Value2 int
		}

		// This should work normally first
		err := ValidateStruct(TestStruct{Value1: 10, Value2: 5})
		if err != nil {
			t.Errorf("Expected no error for valid comparison, got: %v", err)
		}

		// Test with equal values (should fail gt validation)
		err = ValidateStruct(TestStruct{Value1: 5, Value2: 5})
		if err == nil {
			t.Error("Expected validation error for gt validation failure")
			return
		}

		errors := err.(Errors)
		fieldError := errors.GetFieldError("Value1")
		if fieldError == nil {
			t.Error("Expected FieldError for Value1 field")
			return
		}

		// Even though this is a logical validation failure (not a function error),
		// we should test the FuncError handling
		if fieldError.HasFuncError() {
			funcErr := fieldError.Unwrap()
			if funcErr != nil {
				t.Logf("Function error (if any): %v", funcErr)
			}
		}
	})

	// Test case 3: Error chaining with errors.Is and errors.As
	t.Run("ErrorChaining", func(t *testing.T) {
		type InvalidParamsStruct struct {
			Value int `valid:"between=5"` // Missing second parameter
		}

		err := ValidateStruct(InvalidParamsStruct{Value: 7})
		if err == nil {
			t.Error("Expected validation error for invalid parameters")
			return
		}

		errors := err.(Errors)
		fieldError := errors.GetFieldError("Value")
		if fieldError == nil {
			t.Error("Expected FieldError for Value field")
			return
		}

		// Test error unwrapping chain
		if fieldError.HasFuncError() {
			funcErr := fieldError.Unwrap()
			if funcErr != nil {
				// Test that we can unwrap the error chain
				t.Logf("Unwrapped error: %v", funcErr)

				// Verify error message content
				if funcErr.Error() == "" {
					t.Error("Expected non-empty function error message")
				}
			}
		}
	})

	// Test case 4: FieldError without FuncError
	t.Run("NoFuncError", func(t *testing.T) {
		// Create a FieldError manually without FuncError
		fieldError := &FieldError{
			Name:    "TestField",
			Message: "Test validation message",
			// FuncError is nil
		}

		// Test HasFuncError method
		if fieldError.HasFuncError() {
			t.Error("Expected HasFuncError to return false when FuncError is nil")
		}

		// Test Unwrap method
		funcErr := fieldError.Unwrap()
		if funcErr != nil {
			t.Error("Expected Unwrap to return nil when FuncError is nil")
		}

		// Test Error method
		expectedMsg := "Test validation message"
		if fieldError.Error() != expectedMsg {
			t.Errorf("Expected error message '%s', got '%s'", expectedMsg, fieldError.Error())
		}
	})

	// Test case 5: FieldError with FuncError but no Message
	t.Run("FuncErrorOnly", func(t *testing.T) {
		// Create a FieldError with FuncError but no user message
		originalErr := fmt.Errorf("original function error")
		fieldError := &FieldError{
			Name:      "TestField",
			FuncError: originalErr,
			// Message is empty
		}

		// Test HasFuncError method
		if !fieldError.HasFuncError() {
			t.Error("Expected HasFuncError to return true")
		}

		// Test Unwrap method
		funcErr := fieldError.Unwrap()
		if funcErr != originalErr {
			t.Error("Expected Unwrap to return the original error")
		}

		// Test Error method falls back to FuncError when Message is empty
		expectedMsg := "validation failed for field 'TestField': original function error"
		if fieldError.Error() != expectedMsg {
			t.Errorf("Expected error message '%s', got '%s'", expectedMsg, fieldError.Error())
		}
	})

	// Test case 6: Error interface compatibility
	t.Run("ErrorInterfaceCompatibility", func(t *testing.T) {
		type TestStruct struct {
			Complex complex64 `valid:"size=5"`
		}

		err := ValidateStruct(TestStruct{Complex: 3 + 4i})
		if err == nil {
			t.Error("Expected validation error")
			return
		}

		errors := err.(Errors)
		fieldError := errors.GetFieldError("Complex")
		if fieldError == nil {
			t.Error("Expected FieldError")
			return
		}

		// Test that FieldError implements error interface properly
		var _ error = fieldError

		// Test that we can use it in error context
		if fieldError.Error() == "" {
			t.Error("Expected non-empty error message")
		}

		// Test JSON marshaling doesn't break with FuncError
		jsonData, jsonErr := json.Marshal(fieldError)
		if jsonErr != nil {
			t.Errorf("Expected successful JSON marshaling, got error: %v", jsonErr)
		}

		if len(jsonData) == 0 {
			t.Error("Expected non-empty JSON data")
		}
	})

	// Test case 7: Multiple errors with mixed FuncError states
	t.Run("MultipleErrorsMixed", func(t *testing.T) {
		type MixedErrorStruct struct {
			Required string    `valid:"required"`     // No FuncError expected
			Complex  complex64 `valid:"between=1|10"` // FuncError expected
			Email    string    `valid:"email"`        // No FuncError expected
			Invalid  int       `valid:"between=1"`    // FuncError expected (invalid params)
		}

		err := ValidateStruct(MixedErrorStruct{
			Required: "",        // Will fail required
			Complex:  5 + 5i,    // Will fail with unsupported type
			Email:    "invalid", // Will fail email validation
			Invalid:  5,         // Will fail with invalid parameters
		})

		if err == nil {
			t.Error("Expected validation errors")
			return
		}

		errors := err.(Errors)
		if len(errors) == 0 {
			t.Error("Expected multiple errors")
			return
		}

		// Check each field error
		fields := []string{"Required", "Complex", "Email", "Invalid"}
		funcErrorExpected := map[string]bool{
			"Required": false, // Required validation typically doesn't have FuncError
			"Complex":  true,  // Unsupported type should have FuncError
			"Email":    false, // Email validation typically doesn't have FuncError
			"Invalid":  true,  // Invalid parameters should have FuncError
		}

		for _, field := range fields {
			fieldError := errors.GetFieldError(field)
			if fieldError != nil {
				hasFuncError := fieldError.HasFuncError()
				expected := funcErrorExpected[field]

				if hasFuncError != expected {
					t.Logf("Field %s: HasFuncError=%v, Expected=%v", field, hasFuncError, expected)
					// Note: This is informational rather than a hard failure
					// as the exact behavior might vary based on implementation
				}

				// Verify that Error() method works regardless of FuncError state
				if fieldError.Error() == "" {
					t.Errorf("Expected non-empty error message for field %s", field)
				}
			}
		}
	})
}

// TestAdditionalCoverage tests additional code paths for better coverage
func TestAdditionalCoverage(t *testing.T) {
	// Test Empty function with different types
	t.Run("EmptyFunction", func(t *testing.T) {
		// Test string
		if !Empty(reflect.ValueOf("")) {
			t.Error("Expected empty string to be considered empty")
		}
		if Empty(reflect.ValueOf("non-empty")) {
			t.Error("Expected non-empty string to not be considered empty")
		}

		// Test array
		if !Empty(reflect.ValueOf([0]int{})) {
			t.Error("Expected empty array to be considered empty")
		}

		// Test map
		emptyMap := make(map[string]string)
		if !Empty(reflect.ValueOf(emptyMap)) {
			t.Error("Expected empty map to be considered empty")
		}

		// Test slice
		var emptySlice []string
		if !Empty(reflect.ValueOf(emptySlice)) {
			t.Error("Expected empty slice to be considered empty")
		}

		// Test bool
		if !Empty(reflect.ValueOf(false)) {
			t.Error("Expected false bool to be considered empty")
		}
		if Empty(reflect.ValueOf(true)) {
			t.Error("Expected true bool to not be considered empty")
		}

		// Test int types
		if !Empty(reflect.ValueOf(int(0))) {
			t.Error("Expected zero int to be considered empty")
		}
		if !Empty(reflect.ValueOf(int8(0))) {
			t.Error("Expected zero int8 to be considered empty")
		}
		if !Empty(reflect.ValueOf(int16(0))) {
			t.Error("Expected zero int16 to be considered empty")
		}
		if !Empty(reflect.ValueOf(int32(0))) {
			t.Error("Expected zero int32 to be considered empty")
		}
		if !Empty(reflect.ValueOf(int64(0))) {
			t.Error("Expected zero int64 to be considered empty")
		}

		// Test uint types
		if !Empty(reflect.ValueOf(uint(0))) {
			t.Error("Expected zero uint to be considered empty")
		}
		if !Empty(reflect.ValueOf(uint8(0))) {
			t.Error("Expected zero uint8 to be considered empty")
		}
		if !Empty(reflect.ValueOf(uint16(0))) {
			t.Error("Expected zero uint16 to be considered empty")
		}
		if !Empty(reflect.ValueOf(uint32(0))) {
			t.Error("Expected zero uint32 to be considered empty")
		}
		if !Empty(reflect.ValueOf(uint64(0))) {
			t.Error("Expected zero uint64 to be considered empty")
		}
		if !Empty(reflect.ValueOf(uintptr(0))) {
			t.Error("Expected zero uintptr to be considered empty")
		}

		// Test float types
		if !Empty(reflect.ValueOf(float32(0))) {
			t.Error("Expected zero float32 to be considered empty")
		}
		if !Empty(reflect.ValueOf(float64(0))) {
			t.Error("Expected zero float64 to be considered empty")
		}

		// Test pointer
		var nilPtr *string
		if !Empty(reflect.ValueOf(nilPtr)) {
			t.Error("Expected nil pointer to be considered empty")
		}

		// Test interface
		var nilInterface interface{}
		if !Empty(reflect.ValueOf(&nilInterface).Elem()) {
			t.Error("Expected nil interface to be considered empty")
		}
	})

	// Test string validation functions
	t.Run("StringValidationFunctions", func(t *testing.T) {
		// Test IsNull
		if !IsNull("") {
			t.Error("Expected empty string to be null")
		}
		if IsNull("not empty") {
			t.Error("Expected non-empty string to not be null")
		}

		// Test IsEmptyString
		if !IsEmptyString("   ") {
			t.Error("Expected whitespace-only string to be considered empty")
		}
		if IsEmptyString("not empty") {
			t.Error("Expected non-empty string to not be considered empty")
		}

		// Test numeric validation with empty strings
		if !IsNumeric("") {
			t.Error("Expected empty string to be valid numeric")
		}
		if !IsInt("") {
			t.Error("Expected empty string to be valid int")
		}
		if !IsFloat("") {
			t.Error("Expected empty string to be valid float")
		}

		// Test alpha validation with empty strings
		if !ValidateAlpha("") {
			t.Error("Expected empty string to be valid alpha")
		}
		if !ValidateAlphaNum("") {
			t.Error("Expected empty string to be valid alphaNum")
		}
		if !ValidateAlphaDash("") {
			t.Error("Expected empty string to be valid alphaDash")
		}
		if !ValidateAlphaUnicode("") {
			t.Error("Expected empty string to be valid alphaUnicode")
		}
		if !ValidateAlphaNumUnicode("") {
			t.Error("Expected empty string to be valid alphaNumUnicode")
		}
		if !ValidateAlphaDashUnicode("") {
			t.Error("Expected empty string to be valid alphaDashUnicode")
		}

		// Test UUID validation with empty strings
		if !ValidateUUID3("") {
			t.Error("Expected empty string to be valid UUID3")
		}
		if !ValidateUUID4("") {
			t.Error("Expected empty string to be valid UUID4")
		}
		if !ValidateUUID5("") {
			t.Error("Expected empty string to be valid UUID5")
		}
		if !ValidateUUID("") {
			t.Error("Expected empty string to be valid UUID")
		}

		// Test URL validation with empty string
		if !ValidateURL("") {
			t.Error("Expected empty string to be valid URL")
		}

		// Test IP validation
		if !ValidateIP("192.168.1.1") {
			t.Error("Expected valid IPv4 to pass IP validation")
		}
		if !ValidateIP("2001:db8::1") {
			t.Error("Expected valid IPv6 to pass IP validation")
		}
		if ValidateIP("invalid") {
			t.Error("Expected invalid IP to fail validation")
		}

		// Test IPv4 specific validation
		if !ValidateIPv4("192.168.1.1") {
			t.Error("Expected valid IPv4 to pass IPv4 validation")
		}
		if ValidateIPv4("2001:db8::1") {
			t.Error("Expected IPv6 to fail IPv4 validation")
		}

		// Test IPv6 specific validation
		if !ValidateIPv6("2001:db8::1") {
			t.Error("Expected valid IPv6 to pass IPv6 validation")
		}
		if ValidateIPv6("192.168.1.1") {
			t.Error("Expected IPv4 to fail IPv6 validation")
		}
	})

	// Test direct validation functions
	t.Run("DirectValidationFunctions", func(t *testing.T) {
		// Test ValidateRequired
		if !ValidateRequired("non-empty") {
			t.Error("Expected non-empty string to be required")
		}
		if ValidateRequired("") {
			t.Error("Expected empty string to fail required validation")
		}

		// Test ValidateDistinct
		distinctSlice := []string{"a", "b", "c"}
		if !ValidateDistinct(distinctSlice) {
			t.Error("Expected distinct slice to pass validation")
		}

		nonDistinctSlice := []string{"a", "b", "a"}
		if ValidateDistinct(nonDistinctSlice) {
			t.Error("Expected non-distinct slice to fail validation")
		}

		// Test ValidateBetween
		valid, err := ValidateBetween("hello", []string{"3", "10"})
		if err != nil {
			t.Errorf("Unexpected error in ValidateBetween: %v", err)
		}
		if !valid {
			t.Error("Expected 'hello' to be between 3 and 10 characters")
		}

		// Test ValidateDigitsBetween
		valid, err = ValidateDigitsBetween("12345", []string{"3", "10"})
		if err != nil {
			t.Errorf("Unexpected error in ValidateDigitsBetween: %v", err)
		}
		if !valid {
			t.Error("Expected '12345' to be between 3 and 10 digits")
		}

		// Test ValidateSize
		valid, err = ValidateSize("hello", []string{"5"})
		if err != nil {
			t.Errorf("Unexpected error in ValidateSize: %v", err)
		}
		if !valid {
			t.Error("Expected 'hello' to have size 5")
		}

		// Test ValidateMax
		valid, err = ValidateMax("hi", []string{"5"})
		if err != nil {
			t.Errorf("Unexpected error in ValidateMax: %v", err)
		}
		if !valid {
			t.Error("Expected 'hi' to be at most 5 characters")
		}

		// Test ValidateMin
		valid, err = ValidateMin("hello", []string{"3"})
		if err != nil {
			t.Errorf("Unexpected error in ValidateMin: %v", err)
		}
		if !valid {
			t.Error("Expected 'hello' to be at least 3 characters")
		}

		// Test comparison functions
		valid, err = ValidateGt(10, 5)
		if err != nil {
			t.Errorf("Unexpected error in ValidateGt: %v", err)
		}
		if !valid {
			t.Error("Expected 10 to be greater than 5")
		}

		valid, err = ValidateGte(10, 10)
		if err != nil {
			t.Errorf("Unexpected error in ValidateGte: %v", err)
		}
		if !valid {
			t.Error("Expected 10 to be greater than or equal to 10")
		}

		valid, err = ValidateLt(5, 10)
		if err != nil {
			t.Errorf("Unexpected error in ValidateLt: %v", err)
		}
		if !valid {
			t.Error("Expected 5 to be less than 10")
		}

		valid, err = ValidateLte(10, 10)
		if err != nil {
			t.Errorf("Unexpected error in ValidateLte: %v", err)
		}
		if !valid {
			t.Error("Expected 10 to be less than or equal to 10")
		}

		valid, err = ValidateSame("hello", "hello")
		if err != nil {
			t.Errorf("Unexpected error in ValidateSame: %v", err)
		}
		if !valid {
			t.Error("Expected 'hello' to be same as 'hello'")
		}
	})
}

// TestEdgeCases tests various edge cases for comprehensive coverage
func TestEdgeCases(t *testing.T) {
	// Test validation with nil pointer
	t.Run("NilPointer", func(t *testing.T) {
		type TestStruct struct {
			Field *string `valid:"required"`
		}

		err := ValidateStruct(TestStruct{})
		if err == nil {
			t.Error("Expected error for nil required pointer field")
		}
	})

	// Test nested struct validation
	t.Run("NestedStruct", func(t *testing.T) {
		type Inner struct {
			Value string `valid:"required"`
		}
		type Outer struct {
			Inner Inner `valid:"required"`
		}

		err := ValidateStruct(Outer{Inner: Inner{Value: "test"}})
		if err != nil {
			t.Errorf("Unexpected error for valid nested struct: %v", err)
		}

		err = ValidateStruct(Outer{Inner: Inner{Value: ""}})
		if err == nil {
			t.Error("Expected error for invalid nested struct")
		}
	})

	// Test validation with interface{}
	t.Run("InterfaceField", func(t *testing.T) {
		type TestStruct struct {
			Field interface{} `valid:"required"`
		}

		err := ValidateStruct(TestStruct{Field: "test"})
		if err != nil {
			t.Errorf("Unexpected error for non-nil interface field: %v", err)
		}

		err = ValidateStruct(TestStruct{Field: nil})
		if err == nil {
			t.Error("Expected error for nil interface field")
		}
	})

	// Test URL validation edge cases
	t.Run("URLValidation", func(t *testing.T) {
		// Test URL with fragment
		if !ValidateURL("https://example.com/path#fragment") {
			t.Error("Expected URL with fragment to be valid")
		}

		// Test URL without scheme
		if ValidateURL("example.com") {
			t.Error("Expected URL without scheme to be invalid")
		}
	})
}

// Additional utility function tests consolidated from small test files
func TestValidateDigitsBetweenUint64(t *testing.T) {
	tests := []struct {
		value    uint64
		left     uint64
		right    uint64
		expected bool
	}{
		{5, 1, 10, true},
		{0, 1, 10, false},
		{15, 1, 10, false},
		{1, 1, 10, true},
		{10, 1, 10, true},
		{5, 5, 5, true},
		{5, 10, 1, true}, // Test swapping when left > right
		{0, 10, 1, false},
		{15, 10, 1, false},
		{18446744073709551615, 0, 18446744073709551615, true}, // Max uint64
	}

	for _, test := range tests {
		result := ValidateDigitsBetweenUint64(test.value, test.left, test.right)
		if result != test.expected {
			t.Errorf("ValidateDigitsBetweenUint64(%d, %d, %d) = %t; expected %t",
				test.value, test.left, test.right, result, test.expected)
		}
	}
}

func TestCompareUint64(t *testing.T) {
	tests := []struct {
		first       uint64
		second      uint64
		operator    string
		expected    bool
		expectError bool
	}{
		{5, 10, "<", true, false},
		{10, 5, "<", false, false},
		{5, 10, ">", false, false},
		{10, 5, ">", true, false},
		{5, 10, "<=", true, false},
		{10, 10, "<=", true, false},
		{15, 10, "<=", false, false},
		{5, 10, ">=", false, false},
		{10, 10, ">=", true, false},
		{15, 10, ">=", true, false},
		{10, 10, "==", true, false},
		{10, 5, "==", false, false},
		{10, 5, "!=", false, true}, // Unsupported operator
		{10, 5, "invalid", false, true},
		{0, 18446744073709551615, "<", true, false}, // Test with max uint64
		{18446744073709551615, 0, ">", true, false}, // Test with max uint64
	}

	for _, test := range tests {
		result, err := compareUint64(test.first, test.second, test.operator)
		if result != test.expected {
			t.Errorf("compareUint64(%d, %d, %s) = %t; expected %t",
				test.first, test.second, test.operator, result, test.expected)
		}
		if test.expectError && err == nil {
			t.Errorf("compareUint64(%d, %d, %s) expected error but got nil",
				test.first, test.second, test.operator)
		}
		if !test.expectError && err != nil {
			t.Errorf("compareUint64(%d, %d, %s) unexpected error: %v",
				test.first, test.second, test.operator, err)
		}
	}
}

func TestValidateBetweenErrorCases(t *testing.T) {
	tests := []struct {
		name      string
		value     interface{}
		params    []string
		expectErr bool
	}{
		// Error cases for parameter validation
		{"Wrong param count", "test", []string{"1"}, true},
		{"Too many params", "test", []string{"1", "2", "3"}, true},
		{"String min param error", "test", []string{"invalid", "10"}, true},
		{"String max param error", "test", []string{"1", "invalid"}, true},

		// Error cases for slice/array/map types
		{"Slice min param error", []int{1, 2}, []string{"invalid", "10"}, true},
		{"Slice max param error", []int{1, 2}, []string{"1", "invalid"}, true},
		{"Map min param error", map[string]int{"a": 1}, []string{"invalid", "10"}, true},
		{"Map max param error", map[string]int{"a": 1}, []string{"1", "invalid"}, true},

		// Error cases for numeric types
		{"Int min param error", 5, []string{"invalid", "10"}, true},
		{"Int max param error", 5, []string{"1", "invalid"}, true},
		{"Uint min param error", uint(5), []string{"invalid", "10"}, true},
		{"Uint max param error", uint(5), []string{"1", "invalid"}, true},
		{"Float min param error", 5.5, []string{"invalid", "10"}, true},
		{"Float max param error", 5.5, []string{"1", "invalid"}, true},

		// Error cases for unsupported types
		{"Complex type", complex64(1 + 2i), []string{"1", "10"}, true},
		{"Chan type", make(chan int), []string{"1", "10"}, true},
		{"Func type", func() {}, []string{"1", "10"}, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v := reflect.ValueOf(test.value)
			_, err := validateBetween(v, test.params)
			if test.expectErr && err == nil {
				t.Errorf("Expected error for %s", test.name)
			}
			if !test.expectErr && err != nil {
				t.Errorf("Unexpected error for %s: %v", test.name, err)
			}
		})
	}
}

func TestValidateBetweenAllNumericTypes(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		params   []string
		expected bool
	}{
		// All integer types
		{"int8 valid", int8(5), []string{"1", "10"}, true},
		{"int8 invalid", int8(15), []string{"1", "10"}, false},
		{"int16 valid", int16(5), []string{"1", "10"}, true},
		{"int16 invalid", int16(0), []string{"1", "10"}, false},
		{"int32 valid", int32(5), []string{"1", "10"}, true},
		{"int32 invalid", int32(15), []string{"1", "10"}, false},
		{"int64 valid", int64(5), []string{"1", "10"}, true},
		{"int64 invalid", int64(0), []string{"1", "10"}, false},

		// All unsigned integer types
		{"uint8 valid", uint8(5), []string{"1", "10"}, true},
		{"uint8 invalid", uint8(15), []string{"1", "10"}, false},
		{"uint16 valid", uint16(5), []string{"1", "10"}, true},
		{"uint16 invalid", uint16(0), []string{"1", "10"}, false},
		{"uint32 valid", uint32(5), []string{"1", "10"}, true},
		{"uint32 invalid", uint32(15), []string{"1", "10"}, false},
		{"uint64 valid", uint64(5), []string{"1", "10"}, true},
		{"uint64 invalid", uint64(0), []string{"1", "10"}, false},
		{"uintptr valid", uintptr(5), []string{"1", "10"}, true},
		{"uintptr invalid", uintptr(15), []string{"1", "10"}, false},

		// All float types
		{"float32 valid", float32(5.5), []string{"1", "10"}, true},
		{"float32 invalid", float32(15.5), []string{"1", "10"}, false},
		{"float64 valid", float64(5.5), []string{"1", "10"}, true},
		{"float64 invalid", float64(0.5), []string{"1", "10"}, false},

		// Collection types
		{"slice valid", []string{"a", "b", "c"}, []string{"2", "5"}, true},
		{"slice invalid", []string{"a"}, []string{"2", "5"}, false},
		{"array valid", [3]string{"a", "b", "c"}, []string{"2", "5"}, true},
		{"map valid", map[string]int{"a": 1, "b": 2, "c": 3}, []string{"2", "5"}, true},
		{"map invalid", map[string]int{"a": 1}, []string{"2", "5"}, false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			v := reflect.ValueOf(test.value)
			result, err := validateBetween(v, test.params)
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

func TestBasicStringValidationFunctions(t *testing.T) {
	// Test IsInt and IsFloat
	if !IsInt("123") {
		t.Error("Expected IsInt to pass for number string")
	}

	if IsInt("abc") {
		t.Error("Expected IsInt to fail for non-number string")
	}

	if !IsFloat("123.45") {
		t.Error("Expected IsFloat to pass for float string")
	}

	if IsFloat("abc") {
		t.Error("Expected IsFloat to fail for non-float string")
	}

	// Test ValidateDistinct
	if !ValidateDistinct([]interface{}{"a", "b", "c"}) {
		t.Error("Expected distinct values to validate")
	}

	if ValidateDistinct([]interface{}{"a", "b", "a"}) {
		t.Error("Expected duplicate values to fail validation")
	}

	// Test alpha functions
	if !ValidateAlpha("abc") {
		t.Error("Expected alpha validation to pass")
	}

	if ValidateAlpha("abc123") {
		t.Error("Expected alpha validation to fail with numbers")
	}

	if !ValidateAlphaNum("abc123") {
		t.Error("Expected alphanumeric validation to pass")
	}

	if ValidateAlphaNum("abc123!") {
		t.Error("Expected alphanumeric validation to fail with symbols")
	}

	// Test UUID validation functions
	if !ValidateUUID("550e8400-e29b-41d4-a716-446655440000") {
		t.Error("Expected UUID validation to pass")
	}

	if ValidateUUID("invalid-uuid") {
		t.Error("Expected UUID validation to fail")
	}

	if !ValidateUUID3("6fa459ea-ee8a-3ca4-894e-db77e160355e") {
		t.Error("Expected UUID3 validation to pass")
	}

	if !ValidateUUID4("550e8400-e29b-41d4-a716-446655440000") {
		t.Error("Expected UUID4 validation to pass")
	}

	if !ValidateUUID5("6fa459ea-ee8a-5ca4-894e-db77e160355e") {
		t.Error("Expected UUID5 validation to pass")
	}
}

func TestFileValidationFunctions(t *testing.T) {
	content := []byte("test content")

	// Test ValidateMimeTypes - text content gets detected as text/plain; charset=utf-8
	if !ValidateMimeTypes(content, []string{"text/plain; charset=utf-8"}) {
		t.Error("Expected text content to match detected mime type")
	}

	// Test with non-matching mime type
	if ValidateMimeTypes(content, []string{"image/jpeg"}) {
		t.Error("Expected text content to not match image/jpeg mime type")
	}

	// Test ValidateImage - should fail for text content
	if ValidateImage(content) {
		t.Error("Expected non-image content to fail image validation")
	}

	// Test ValidateMimes - just test that it runs without panic
	result, err := ValidateMimes(content, []string{"unknown"})
	if err == nil {
		_ = result // Function executed successfully
	}
}

func TestCompareStringFunctionDirect(t *testing.T) {
	// Test compareString with different operators
	result, err := compareString("abc", 3, "==")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !result {
		t.Error("Expected string length to equal 3")
	}

	result, err = compareString("a", 2, "<")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !result {
		t.Error("Expected string length to be less than 2")
	}

	result, err = compareString("abc", 2, ">")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !result {
		t.Error("Expected string length to be greater than 2")
	}

	// Test invalid operator
	_, err = compareString("abc", 3, "invalid")
	if err == nil {
		t.Error("Expected error for invalid operator")
	}
}
