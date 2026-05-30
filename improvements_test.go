package validator

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"testing"
)

// TestRequiredIfDoesNotMutateCachedTag is a regression test for a data race +
// unbounded-growth bug: isRequiredIf/checkRequiredIfCondition appended the
// "Value" message parameter to the SHARED cached *ValidTag on every failing
// validation, so parameters accumulated across calls and concurrent validation
// raced on the slice.
func TestRequiredIfDoesNotMutateCachedTag(t *testing.T) {
	type Form struct {
		Status string
		Reason string `valid:"requiredIf=Status|active"`
	}

	var first int
	for i := 1; i <= 3; i++ {
		err := ValidateStruct(&Form{Status: "active", Reason: ""})
		if err == nil {
			t.Fatal("expected error: Reason required when Status is active")
		}
		fe := err.(Errors)[0].(*FieldError)
		if i == 1 {
			first = len(fe.MessageParameters)
		} else if len(fe.MessageParameters) != first {
			t.Fatalf("message parameters grew across calls (%d -> %d): cached tag is being mutated",
				first, len(fe.MessageParameters))
		}
		// The runtime value must still render in the message.
		if !strings.Contains(fe.Message, "active") {
			t.Errorf("call %d: expected message to contain 'active', got %q", i, fe.Message)
		}
	}
}

// TestRequiredIfConcurrentNoRace runs the same requiredIf validation from many
// goroutines; with the cached-tag mutation present this trips the race detector.
func TestRequiredIfConcurrentNoRace(t *testing.T) {
	type Form struct {
		Status string
		Reason string `valid:"requiredIf=Status|active"`
	}
	var wg sync.WaitGroup
	for i := 0; i < 64; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = ValidateStruct(&Form{Status: "active", Reason: ""})
		}()
	}
	wg.Wait()
}

// TestTransTranslatesAllErrorsAfterCustomMessage is a regression test for the
// break-vs-continue bug in Trans: once one field matched a custom message, the
// old code broke out of the loop and left every later error untranslated.
func TestTransTranslatesAllErrorsAfterCustomMessage(t *testing.T) {
	translator := NewTranslator()
	translator.SetMessage("en", Translate{
		"required": "The {{.Attribute}} field is required",
	})
	translator.customMessage["en"] = Translate{
		"email.required": "Email is mandatory",
	}

	errors := Errors{
		&FieldError{Name: "email", MessageName: "required", Attribute: "email"},
		&FieldError{Name: "name", MessageName: "required", Attribute: "name"},
	}

	translated := translator.Trans(errors, "en")

	if got := translated[0].(*FieldError).Message; got != "Email is mandatory" {
		t.Errorf("errors[0]: expected custom message, got %q", got)
	}
	// This is the assertion that fails with the old `break`: the second error
	// must still be translated even though the first used a custom message.
	if got := translated[1].(*FieldError).Message; got != "The name field is required" {
		t.Errorf("errors[1]: expected standard translation, got %q", got)
	}
}

// TestTransSkipsNonFieldErrorButContinues is a regression test for the second
// break: a non-FieldError in the slice must be skipped, not halt translation of
// the remaining FieldErrors.
func TestTransSkipsNonFieldErrorButContinues(t *testing.T) {
	translator := NewTranslator()
	translator.SetMessage("en", Translate{
		"required": "The {{.Attribute}} field is required",
	})

	errors := Errors{
		&struct{ error }{fmt.Errorf("generic error")},
		&FieldError{Name: "name", MessageName: "required", Attribute: "name"},
	}

	translated := translator.Trans(errors, "en")

	if got := translated[1].(*FieldError).Message; got != "The name field is required" {
		t.Errorf("errors[1]: expected translation after a non-FieldError, got %q", got)
	}
}

// TestAppendNamespaceNoAliasing verifies appendNamespace never aliases its base
// slice's backing array, so sibling fields cannot corrupt each other's paths.
func TestAppendNamespaceNoAliasing(t *testing.T) {
	base := make([]byte, 0, 64) // spare capacity is what triggers the old aliasing bug
	base = append(base, "root."...)

	a := appendNamespace(base, []byte("alpha"))
	b := appendNamespace(base, []byte("beta"))

	if string(a) != "root.alpha." {
		t.Errorf("a = %q, want %q", string(a), "root.alpha.")
	}
	if string(b) != "root.beta." {
		t.Errorf("b = %q, want %q", string(b), "root.beta.")
	}
	// With the old append(append(base,...),'.') pattern, building b overwrote a's
	// bytes through the shared backing array. Confirm a is intact.
	if string(a) != "root.alpha." {
		t.Errorf("a was corrupted by building b: got %q", string(a))
	}
}

// TestNestedSliceStructPathBuiltCorrectly exercises the namespace fix end to end:
// a failing element inside a slice of structs must report a correctly composed,
// indexed, uncorrupted path (parent + index + field).
func TestNestedSliceStructPathBuiltCorrectly(t *testing.T) {
	type Item struct {
		Name string `valid:"required"`
	}
	type Order struct {
		Items []Item `valid:"required"`
	}

	err := ValidateStruct(&Order{Items: []Item{{Name: ""}}})
	if err == nil {
		t.Fatal("expected a validation error for the empty item name")
	}

	fe, ok := err.(Errors)[0].(*FieldError)
	if !ok {
		t.Fatalf("expected a *FieldError, got %T", err.(Errors)[0])
	}
	if fe.Name != "Items.0.Name" {
		t.Errorf("expected path %q, got %q", "Items.0.Name", fe.Name)
	}
}

// TestClearSyncMap verifies the in-place clear helper empties a sync.Map without
// reassigning it (the reassignment was a data race with concurrent Load).
func TestClearSyncMap(t *testing.T) {
	var m sync.Map
	m.Store("a", 1)
	m.Store("b", 2)

	clearSyncMap(&m)

	m.Range(func(key, _ any) bool {
		t.Errorf("expected empty map, found key %v", key)
		return true
	})
}

// TestRegisterAutoUnwrapInvalidatesCache ensures registering a new matcher clears
// the resolution cache so previously-cached "no match" results are reconsidered.
func TestRegisterAutoUnwrapInvalidatesCache(t *testing.T) {
	ResetCustomTypeFuncs()
	defer ResetCustomTypeFuncs()

	autoUnwrapCache.Store("sentinel", customTypeExtractFunc(nil))

	RegisterAutoUnwrap(func(reflect.Type) bool { return false }, "IsSet", "Value")

	if _, ok := autoUnwrapCache.Load("sentinel"); ok {
		t.Error("expected autoUnwrapCache to be cleared after RegisterAutoUnwrap")
	}
}
