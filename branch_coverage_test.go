package validator

import (
	"reflect"
	"testing"

	"github.com/shopspring/decimal"
)

func TestIsMultipleOfTypeBranches(t *testing.T) {
	// uint branch
	b, err := IsMultipleOf(uint(10), []string{"5"})
	okErr(t, "IsMultipleOf(uint ok)", b, err, true)

	// float branch
	b, err = IsMultipleOf(4.5, []string{"1.5"})
	okErr(t, "IsMultipleOf(float ok)", b, err, true)
	b, err = IsMultipleOf(4.0, []string{"1.5"})
	okErr(t, "IsMultipleOf(float no)", b, err, false)

	// decimal branch
	b, err = IsMultipleOf(decimal.NewFromInt(9), []string{"3"})
	okErr(t, "IsMultipleOf(decimal ok)", b, err, true)

	// error branches
	if _, err = IsMultipleOf(5, []string{"0"}); err == nil {
		t.Error("IsMultipleOf divide-by-zero should error")
	}
	if _, err = IsMultipleOf(5, nil); err == nil {
		t.Error("IsMultipleOf missing param should error")
	}
	if _, err = IsMultipleOf("notnumeric", []string{"2"}); err == nil {
		t.Error("IsMultipleOf unsupported type should error")
	}
}

func TestIsDistinctTypeBranches(t *testing.T) {
	okBool(t, "IsDistinct(unique slice)", IsDistinct([]int{1, 2, 3}), true)
	okBool(t, "IsDistinct(dup slice)", IsDistinct([]int{1, 1, 2}), false)
	okBool(t, "IsDistinct(map)", IsDistinct(map[string]int{"a": 1, "b": 2}), true)
	okBool(t, "IsDistinct(numeric)", IsDistinct(5), true)
	okBool(t, "IsDistinct(unsupported)", IsDistinct("string"), false)
}

func TestIsDoesntContainSliceBranch(t *testing.T) {
	b, err := IsDoesntContain([]string{"a", "b"}, []string{"c"})
	okErr(t, "IsDoesntContain(slice ok)", b, err, true)
	b, err = IsDoesntContain([]string{"a", "b"}, []string{"a"})
	okErr(t, "IsDoesntContain(slice hit)", b, err, false)
}

func TestMissingParamErrorPaths(t *testing.T) {
	cases := []struct {
		name string
		fn   func() (bool, error)
	}{
		{"IsMaxDigits", func() (bool, error) { return IsMaxDigits(1, nil) }},
		{"IsMinDigits", func() (bool, error) { return IsMinDigits(1, nil) }},
		{"IsContains", func() (bool, error) { return IsContains("x", nil) }},
		{"IsDecimalPrecision", func() (bool, error) { return IsDecimalPrecision(1.0, nil) }},
		{"IsMultipleOf", func() (bool, error) { return IsMultipleOf(1, nil) }},
	}
	for _, c := range cases {
		if _, err := c.fn(); err == nil {
			t.Errorf("%s with no params should return an error", c.name)
		}
	}
}

func TestDerefChains(t *testing.T) {
	s := "value"
	p := &s
	pp := &p

	got := deref(reflect.ValueOf(pp))
	if got.Kind() != reflect.String || got.String() != "value" {
		t.Errorf("deref(**string) = %v (%s), want \"value\"", got, got.Kind())
	}

	// nil pointer: deref must stop and return the (nil) pointer value, not panic.
	var np *string
	gotNil := deref(reflect.ValueOf(np))
	if gotNil.Kind() != reflect.Ptr || !gotNil.IsNil() {
		t.Errorf("deref(nil *string) should return the nil pointer value, got %s", gotNil.Kind())
	}

	// interface wrapping a pointer.
	var iface interface{} = p
	gotIface := deref(reflect.ValueOf(&iface).Elem())
	if gotIface.Kind() != reflect.String || gotIface.String() != "value" {
		t.Errorf("deref(interface{*string}) = %v (%s), want \"value\"", gotIface, gotIface.Kind())
	}
}

// TestFailFast verifies the opt-in FailFast option: collect-all by default,
// stop at the first failing field when enabled.
func TestFailFast(t *testing.T) {
	type Form struct {
		A string `valid:"required"`
		B string `valid:"required"`
		C string `valid:"required"`
	}

	// Default (collect-all): all failing fields are reported.
	collect := New()
	err := collect.ValidateStruct(&Form{}, nil, nil)
	if err == nil {
		t.Fatal("expected errors for empty required fields")
	}
	if n := len(err.(Errors)); n < 3 {
		t.Errorf("collect-all: expected >=3 errors, got %d", n)
	}

	// FailFast: stop at the first failing field.
	ff := New()
	ff.FailFast = true
	err = ff.ValidateStruct(&Form{}, nil, nil)
	if err == nil {
		t.Fatal("expected an error in fail-fast mode")
	}
	if n := len(err.(Errors)); n != 1 {
		t.Errorf("fail-fast: expected exactly 1 error, got %d", n)
	}
	if fe, ok := err.(Errors)[0].(*FieldError); ok && fe.Name != "A" {
		t.Errorf("fail-fast: expected first field 'A', got %q", fe.Name)
	}
}

// TestRequiredIfWithSliceOtherField exercises extractValuesFromCollection's slice
// path: requiredIf referencing a field whose value is a slice.
func TestRequiredIfWithSliceOtherField(t *testing.T) {
	type Form struct {
		Tags []string
		Role string `valid:"requiredIf=Tags|admin"`
	}

	// Tags contains "admin" -> Role becomes required, so empty Role fails.
	if err := ValidateStruct(&Form{Tags: []string{"admin"}, Role: ""}); err == nil {
		t.Error("expected error: Role required when Tags contains admin")
	}

	// Tags contains "admin" and Role provided -> passes.
	if err := ValidateStruct(&Form{Tags: []string{"admin"}, Role: "x"}); err != nil {
		t.Errorf("unexpected error when Role provided: %v", err)
	}

	// Tags has no matching value -> Role not required, empty passes.
	if err := ValidateStruct(&Form{Tags: []string{"user"}, Role: ""}); err != nil {
		t.Errorf("unexpected error when Tags has no admin: %v", err)
	}
}
