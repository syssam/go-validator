package validator

import (
	"reflect"
	"testing"
)

// TestNoDuplicateFields is a regression test for a field-cache bug: typefields
// (ported from encoding/json) incremented count per field and kept the
// duplicate-add without the annihilation pass, so every field after the first
// was duplicated.
func TestNoDuplicateFields(t *testing.T) {
	type Flat struct {
		A string `valid:"required"`
		B string `valid:"required"`
		C string `valid:"required"`
	}
	fields := cachedTypefields(reflect.TypeOf(Flat{}))
	if len(fields) != 3 {
		names := make([]string, len(fields))
		for i, f := range fields {
			names[i] = f.name
		}
		t.Errorf("expected 3 distinct fields, got %d: %v", len(fields), names)
	}
}

// TestNoDuplicateErrors verifies the duplicate-field bug no longer doubles the
// reported errors.
func TestNoDuplicateErrors(t *testing.T) {
	type Flat struct {
		A string `valid:"required"`
		B string `valid:"required"`
		C string `valid:"required"`
	}
	err := ValidateStruct(&Flat{})
	if err == nil {
		t.Fatal("expected errors")
	}
	if n := len(err.(Errors)); n != 3 {
		t.Errorf("expected exactly 3 errors (one per field), got %d", n)
	}
}

type embBase struct {
	ID string `valid:"required"`
}
type embMid struct {
	Code string `valid:"required"`
}

// TestEmbeddedStructFieldsValidated is a regression test: promoted fields from
// embedded (anonymous) structs were never collected — line 130 returned early
// for the untagged embedded field — so their rules silently never fired.
func TestEmbeddedStructFieldsValidated(t *testing.T) {
	type User struct {
		embBase        // value-embedded
		*embMid        // pointer-embedded
		Name    string `valid:"required"`
	}

	// All promoted required fields empty -> all must be reported.
	err := ValidateStruct(&User{embMid: &embMid{}})
	if err == nil {
		t.Fatal("expected errors for empty promoted fields")
	}
	got := map[string]bool{}
	for _, e := range err.(Errors) {
		if fe, ok := e.(*FieldError); ok {
			got[fe.Name] = true
		}
	}
	for _, want := range []string{"ID", "Code", "Name"} {
		if !got[want] {
			t.Errorf("expected promoted field %q to be validated; got %v", want, got)
		}
	}

	// nil pointer-embedded struct must not panic; its fields are simply absent.
	err = ValidateStruct(&User{})
	if err == nil {
		t.Fatal("expected errors")
	}
	for _, e := range err.(Errors) {
		if fe, ok := e.(*FieldError); ok && fe.Name == "Code" {
			t.Error("Code should be skipped when its embedded *struct is nil")
		}
	}

	// Fully valid -> no error.
	if err := ValidateStruct(&User{embBase: embBase{ID: "x"}, embMid: &embMid{Code: "c"}, Name: "n"}); err != nil {
		t.Errorf("unexpected error for valid embedded struct: %v", err)
	}
}

type nestedAddr struct {
	Street string `valid:"required"`
}

// TestUntaggedNestedStructsValidated verifies untagged nested struct and
// pointer-to-struct fields are validated recursively (previously they were
// silently skipped unless the field itself carried a valid tag).
func TestUntaggedNestedStructsValidated(t *testing.T) {
	type ByValue struct{ Addr nestedAddr }
	if err := ValidateStruct(&ByValue{Addr: nestedAddr{Street: ""}}); err == nil {
		t.Error("untagged named nested struct: inner required field not validated")
	}

	type ByPtr struct{ Addr *nestedAddr }
	if err := ValidateStruct(&ByPtr{Addr: &nestedAddr{Street: ""}}); err == nil {
		t.Error("untagged *struct: inner required field not validated")
	}

	// nil pointer-to-struct must be skipped without a panic or error.
	if err := ValidateStruct(&ByPtr{}); err != nil {
		t.Errorf("nil *struct should be skipped, got %v", err)
	}

	// Valid inner -> no error.
	if err := ValidateStruct(&ByValue{Addr: nestedAddr{Street: "x"}}); err != nil {
		t.Errorf("valid nested struct should pass, got %v", err)
	}
}

type linkNode struct {
	Name string    `valid:"required"`
	Next *linkNode `valid:"required"`
}

// TestDeepNestingIsLinear guards against the O(2^depth) regression: validating a
// deep (but bounded) chain must complete quickly. Before the field-cache fix
// this was exponential — a depth of ~20 already took seconds. The suite's
// -timeout fails the test if the exponential behavior returns.
func TestDeepNestingIsLinear(t *testing.T) {
	head := &linkNode{Name: "n"}
	cur := head
	for i := 0; i < 500; i++ {
		cur.Next = &linkNode{Name: "n"}
		cur = cur.Next
	}
	// The tail's Next is nil (required) -> a validation error, not a hang.
	if err := ValidateStruct(head); err == nil {
		t.Error("expected a validation error for the nil tail Next")
	}
}

// TestCyclicReferenceTerminates is a regression test for an infinite-recursion
// DoS: a cyclic object graph must terminate (via the depth guard) instead of
// hanging the goroutine forever.
func TestCyclicReferenceTerminates(t *testing.T) {
	a := &linkNode{Name: "a"}
	b := &linkNode{Name: "b"}
	a.Next = b
	b.Next = a // cycle

	// If the depth guard regresses, this hangs and the suite -timeout fails it.
	err := ValidateStruct(a)
	if err == nil {
		t.Error("expected an error (max depth) for a cyclic reference, got nil")
	}
}
