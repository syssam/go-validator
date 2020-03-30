package validator

import (
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
			t.Errorf("Expected validateateStruct(%q) Case %d to be %v, got %v", test.param, i, test.expected, actual)
			if err != nil {
				t.Errorf("Got Error on validateateStruct(%q): %s", test.param, err)
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
		if v.Kind() != reflect.String {
			return true
		}

		return false
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

	ValidateStruct(u)
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
	if err != nil {
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

	ValidateStruct(u)
}
