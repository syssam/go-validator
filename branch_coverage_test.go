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
	if gotNil.Kind() != reflect.Pointer || !gotNil.IsNil() {
		t.Errorf("deref(nil *string) should return the nil pointer value, got %s", gotNil.Kind())
	}

	// interface wrapping a pointer.
	var iface any = p
	gotIface := deref(reflect.ValueOf(&iface).Elem())
	if gotIface.Kind() != reflect.String || gotIface.String() != "value" {
		t.Errorf("deref(interface{*string}) = %v (%s), want \"value\"", gotIface, gotIface.Kind())
	}
}

// TestDateParsesConsistently is a regression test: a date string must parse to
// the same date whether it appears as a field value (parseDate) or a rule
// parameter (parseDateParam). They previously used divergent format lists, so
// "03/04/2006" became Apr 3 in one path and Mar 4 in the other.
func TestDateParsesConsistently(t *testing.T) {
	for _, s := range []string{"03/04/2006", "2024-01-15", "2024/01/15", "01/02/2006"} {
		d1, err1 := parseDate(s)
		d2, err2 := parseDateParam(s)
		if err1 != nil || err2 != nil {
			t.Fatalf("%s: parse error (%v / %v)", s, err1, err2)
		}
		if !d1.Equal(d2) {
			t.Errorf("%s parses inconsistently: value=%s param=%s",
				s, d1.Format("2006-01-02"), d2.Format("2006-01-02"))
		}
	}
}

// TestOmitemptyNotMatchedAsSubstring is a regression test: a parameter value
// that merely contains "omitempty" must not flip the field into omitempty mode.
func TestOmitemptyNotMatchedAsSubstring(t *testing.T) {
	type S struct {
		// "omitempty" appears only inside a rule parameter, not as an option.
		Name string `valid:"required,contains=omitempty"`
	}
	// Name is empty: required must fire. With the old strings.Contains check the
	// field was wrongly treated as omitempty and skipped entirely.
	if err := ValidateStruct(&S{Name: ""}); err == nil {
		t.Error("empty required field was skipped — omitempty matched as a substring")
	}
}

// TestInterfaceValuedMapOfStructsValidated is a regression test: entries of a
// map[string]any holding structs must be validated (they were silently
// skipped because the deref check tested the map's kind, not the entry's).
func TestInterfaceValuedMapOfStructsValidated(t *testing.T) {
	type Inner struct {
		Name string `valid:"required"`
	}
	type Outer struct {
		Items map[string]any `valid:"required"`
	}
	err := ValidateStruct(&Outer{Items: map[string]any{"a": Inner{Name: ""}}})
	if err == nil {
		t.Error("expected the interface-valued map entry's struct to be validated")
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
