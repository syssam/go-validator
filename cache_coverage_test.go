package validator

import (
	"reflect"
	"testing"
)

// Test shouldSkipField function edge cases
func TestShouldSkipField(t *testing.T) {
	type TestStruct struct {
		PublicField  string
		privateField string
	}

	structType := reflect.TypeOf(TestStruct{})

	// Test public field - should not be skipped
	publicField := structType.Field(0)
	if shouldSkipField(publicField) {
		t.Error("Public field should not be skipped")
	}

	// Test private field - should be skipped
	privateField := structType.Field(1)
	if !shouldSkipField(privateField) {
		t.Error("Private field should be skipped")
	}
}

// Test isValidAttribute function (currently unused)
func TestIsValidAttribute(t *testing.T) {
	f := &field{}

	// Test function exists and handles empty string
	result := f.isValidAttribute("")
	if result {
		t.Error("Expected empty string to be invalid")
	}

	// The function rejects most strings due to its character restrictions
	result = f.isValidAttribute("validattribute") // no special chars
	if !result {
		t.Error("Expected simple string without special chars to be valid")
	}
}

// Test parseMessageParameterIntoSlice edge cases
func TestParseMessageParameterIntoSlice(t *testing.T) {
	f := &field{}

	tests := []struct {
		name        string
		rule        string
		params      []string
		expectError bool
		expectNil   bool
	}{
		{
			"requiredUnless with valid params",
			"requiredUnless",
			[]string{"field", "value1", "value2"},
			false,
			false,
		},
		{
			"requiredUnless with insufficient params",
			"requiredUnless",
			[]string{"field"},
			true,
			false,
		},
		{
			"between with valid params",
			"between",
			[]string{"1", "10"},
			false,
			false,
		},
		{
			"between with invalid param count",
			"between",
			[]string{"1"},
			true,
			false,
		},
		{
			"digitsBetween with valid params",
			"digitsBetween",
			[]string{"1", "10"},
			false,
			false,
		},
		{
			"gt with valid param",
			"gt",
			[]string{"5"},
			false,
			false,
		},
		{
			"gt with invalid param count",
			"gt",
			[]string{},
			true,
			false,
		},
		{
			"gte with valid param",
			"gte",
			[]string{"5"},
			false,
			false,
		},
		{
			"lt with valid param",
			"lt",
			[]string{"5"},
			false,
			false,
		},
		{
			"lte with valid param",
			"lte",
			[]string{"5"},
			false,
			false,
		},
		{
			"max with valid param",
			"max",
			[]string{"10"},
			false,
			false,
		},
		{
			"min with valid param",
			"min",
			[]string{"1"},
			false,
			false,
		},
		{
			"size with valid param",
			"size",
			[]string{"5"},
			false,
			false,
		},
		{
			"unknown rule",
			"unknown",
			[]string{},
			false,
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := f.parseMessageParameterIntoSlice(test.rule, test.params...)

			if test.expectError && err == nil {
				t.Errorf("Expected error for %s", test.name)
			}
			if !test.expectError && err != nil {
				t.Errorf("Unexpected error for %s: %v", test.name, err)
			}
			if test.expectNil && result != nil {
				t.Errorf("Expected nil result for %s", test.name)
			}
			if !test.expectNil && !test.expectError && result == nil {
				t.Errorf("Expected non-nil result for %s", test.name)
			}
		})
	}
}

// Test parseMessageName edge cases
func TestParseMessageName(t *testing.T) {
	f := &field{}

	tests := []struct {
		rule      string
		fieldType reflect.Type
		expected  string
	}{
		{"between", reflect.TypeOf(0), "between.numeric"},
		{"between", reflect.TypeOf(""), "between.string"},
		{"between", reflect.TypeOf([]int{}), "between.array"},
		{"between", reflect.TypeOf(map[string]int{}), "between.array"},
		{"between", reflect.TypeOf(struct{}{}), "between"},
		{"gt", reflect.TypeOf(int8(0)), "gt.numeric"},
		{"gt", reflect.TypeOf(int16(0)), "gt.numeric"},
		{"gt", reflect.TypeOf(int32(0)), "gt.numeric"},
		{"gt", reflect.TypeOf(int64(0)), "gt.numeric"},
		{"gt", reflect.TypeOf(uint(0)), "gt.numeric"},
		{"gt", reflect.TypeOf(uint8(0)), "gt.numeric"},
		{"gt", reflect.TypeOf(uint16(0)), "gt.numeric"},
		{"gt", reflect.TypeOf(uint32(0)), "gt.numeric"},
		{"gt", reflect.TypeOf(uint64(0)), "gt.numeric"},
		{"gt", reflect.TypeOf(float32(0)), "gt.numeric"},
		{"gt", reflect.TypeOf(float64(0)), "gt.numeric"},
		{"gte", reflect.TypeOf(""), "gte.string"},
		{"lt", reflect.TypeOf([]int{}), "lt.array"},
		{"lte", reflect.TypeOf(&struct{}{}), "lte"},
		{"min", reflect.TypeOf(0), "min.numeric"},
		{"max", reflect.TypeOf(""), "max.string"},
		{"size", reflect.TypeOf([]int{}), "size.array"},
		{"unknown", reflect.TypeOf(0), "unknown"},
	}

	for _, test := range tests {
		t.Run(test.rule+"_"+test.fieldType.String(), func(t *testing.T) {
			result := f.parseMessageName(test.rule, test.fieldType)
			if result != test.expected {
				t.Errorf("parseMessageName(%s, %s) = %s; expected %s",
					test.rule, test.fieldType, result, test.expected)
			}
		})
	}
}

// Test cachedTypefields with complex nested types
func TestCachedTypefieldsComplex(t *testing.T) {
	type NestedStruct struct {
		Field1 string `valid:"required"`
		Field2 int    `valid:"min=1"`
	}

	type ComplexStruct struct {
		Simple  string       `valid:"email"`
		Nested  NestedStruct `valid:"required"`
		Slice   []string     `valid:"required"`
		Map     map[string]int
		Pointer *string `valid:"required"`
		ignored string  `valid:"-"`
	}

	// First call should populate cache
	fields1 := cachedTypefields(reflect.TypeOf(ComplexStruct{}))

	// Second call should use cache
	fields2 := cachedTypefields(reflect.TypeOf(ComplexStruct{}))

	// Should return same results
	if len(fields1) != len(fields2) {
		t.Errorf("Cache returned different number of fields: %d vs %d", len(fields1), len(fields2))
	}

	// Verify specific fields are present and have correct properties
	fieldNames := make(map[string]bool)
	for _, field := range fields1 {
		fieldNames[field.name] = true
	}

	expectedFields := []string{"Simple", "Nested", "Slice", "Pointer"}
	for _, expected := range expectedFields {
		if !fieldNames[expected] {
			t.Errorf("Expected field %s not found in cached results", expected)
		}
	}

	// Verify ignored field is not present
	if fieldNames["ignored"] {
		t.Error("Field marked with '-' should be ignored")
	}
}

// Test processStructField with various field configurations
func TestProcessStructField(t *testing.T) {
	type TestStruct struct {
		ValidField   string `valid:"required"`
		IgnoredField string `valid:"-"`
		EmptyTag     string
		unexported   string `valid:"required"`
	}

	structType := reflect.TypeOf(TestStruct{})
	f := &field{typ: structType}
	count := make(map[reflect.Type]int)
	nextCount := make(map[reflect.Type]int)
	var fields []field
	var next []field

	// Test each field
	for i := 0; i < structType.NumField(); i++ {
		sf := structType.Field(i)
		processStructField(sf, f, structType, i, count, nextCount, &fields, &next)
	}

	// Should have processed valid fields but skipped others
	fieldNames := make(map[string]bool)
	for _, field := range fields {
		fieldNames[field.name] = true
	}

	if !fieldNames["ValidField"] {
		t.Error("ValidField should be processed")
	}
	if fieldNames["IgnoredField"] {
		t.Error("IgnoredField should be ignored due to '-' tag")
	}
	if fieldNames["EmptyTag"] {
		t.Error("EmptyTag should be ignored due to empty validation tag")
	}
	if fieldNames["unexported"] {
		t.Error("unexported field should be skipped")
	}
}
