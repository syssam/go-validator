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
