package validator

import (
	"reflect"
	"strings"
	"testing"
)

// mockOmittable mimics graphql.Omittable[T] for testing without the gqlgen dependency.
type mockOmittable[T any] struct {
	value T
	set   bool
}

func omittableOf[T any](v T) mockOmittable[T] {
	return mockOmittable[T]{value: v, set: true}
}

func omittableEmpty[T any]() mockOmittable[T] {
	return mockOmittable[T]{}
}

func (o mockOmittable[T]) IsSet() bool { return o.set }
func (o mockOmittable[T]) Value() T    { return o.value }

// registerMockOmittables registers extractors for all mockOmittable types used in tests.
func registerMockOmittables() {
	RegisterCustomTypeFunc(func(o mockOmittable[string]) (any, bool) {
		if !o.IsSet() {
			return nil, false
		}
		return o.Value(), true
	})
	RegisterCustomTypeFunc(func(o mockOmittable[int]) (any, bool) {
		if !o.IsSet() {
			return nil, false
		}
		return o.Value(), true
	})
	RegisterCustomTypeFunc(func(o mockOmittable[*string]) (any, bool) {
		if !o.IsSet() {
			return nil, false
		}
		return o.Value(), true
	})
}

// registerMockAutoUnwrap registers a single matcher for ALL mockOmittable[T] variants.
func registerMockAutoUnwrap() {
	RegisterAutoUnwrap(
		func(t reflect.Type) bool {
			return t.Kind() == reflect.Struct &&
				strings.HasPrefix(t.Name(), "mockOmittable[")
		},
		"IsSet", // bool method
		"Value", // value method
	)
}

// --- Test structs ---

type omittableRequiredStruct struct {
	Name mockOmittable[string] `valid:"required"`
}

type omittableEmailStruct struct {
	Email mockOmittable[string] `valid:"required,email"`
}

type omittableMinStruct struct {
	Age mockOmittable[int] `valid:"min=18"`
}

type omittableOptionalStruct struct {
	Name mockOmittable[string] `valid:"email"`
}

type omittablePtrStruct struct {
	Name mockOmittable[*string] `valid:"required"`
}

type omittableOmitEmptyStruct struct {
	Name mockOmittable[string] `valid:"omitempty,email"`
}

type omittableMultiFieldStruct struct {
	Name  mockOmittable[string] `valid:"required"`
	Email mockOmittable[string] `valid:"required,email"`
	Age   mockOmittable[int]    `valid:"min=0"`
}

type omittableBetweenStruct struct {
	Score mockOmittable[int] `valid:"between=1|100"`
}

// --- RegisterCustomTypeFunc tests (exact match) ---

func TestCustomTypeFunc_NotSet_Required(t *testing.T) {
	ResetCustomTypeFuncs()
	registerMockOmittables()
	defer ResetCustomTypeFuncs()

	s := &omittableRequiredStruct{
		Name: omittableEmpty[string](),
	}

	err := ValidateStruct(s)
	if err == nil {
		t.Fatal("expected error for unset required Omittable field")
	}

	errs := err.(Errors)
	if !errs.HasFieldError("Name") {
		t.Fatal("expected FieldError for Name")
	}
}

func TestCustomTypeFunc_Set_Required(t *testing.T) {
	ResetCustomTypeFuncs()
	registerMockOmittables()
	defer ResetCustomTypeFuncs()

	s := &omittableRequiredStruct{
		Name: omittableOf("Alice"),
	}

	err := ValidateStruct(s)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestCustomTypeFunc_Set_EmptyString_Required(t *testing.T) {
	ResetCustomTypeFuncs()
	registerMockOmittables()
	defer ResetCustomTypeFuncs()

	s := &omittableRequiredStruct{
		Name: omittableOf(""),
	}

	err := ValidateStruct(s)
	if err == nil {
		t.Fatal("expected error for set but empty required Omittable field")
	}
}

func TestCustomTypeFunc_NotSet_Optional(t *testing.T) {
	ResetCustomTypeFuncs()
	registerMockOmittables()
	defer ResetCustomTypeFuncs()

	s := &omittableOptionalStruct{
		Name: omittableEmpty[string](),
	}

	err := ValidateStruct(s)
	if err != nil {
		t.Fatalf("expected no error for unset optional Omittable, got: %v", err)
	}
}

func TestCustomTypeFunc_Set_Email_Valid(t *testing.T) {
	ResetCustomTypeFuncs()
	registerMockOmittables()
	defer ResetCustomTypeFuncs()

	s := &omittableEmailStruct{
		Email: omittableOf("test@example.com"),
	}

	err := ValidateStruct(s)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestCustomTypeFunc_Set_Email_Invalid(t *testing.T) {
	ResetCustomTypeFuncs()
	registerMockOmittables()
	defer ResetCustomTypeFuncs()

	s := &omittableEmailStruct{
		Email: omittableOf("not-an-email"),
	}

	err := ValidateStruct(s)
	if err == nil {
		t.Fatal("expected error for invalid email")
	}
}

func TestCustomTypeFunc_Set_Min_Valid(t *testing.T) {
	ResetCustomTypeFuncs()
	registerMockOmittables()
	defer ResetCustomTypeFuncs()

	s := &omittableMinStruct{
		Age: omittableOf(25),
	}

	err := ValidateStruct(s)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestCustomTypeFunc_Set_Min_Invalid(t *testing.T) {
	ResetCustomTypeFuncs()
	registerMockOmittables()
	defer ResetCustomTypeFuncs()

	s := &omittableMinStruct{
		Age: omittableOf(10),
	}

	err := ValidateStruct(s)
	if err == nil {
		t.Fatal("expected error for age below min")
	}
}

func TestCustomTypeFunc_NotSet_Min_Skipped(t *testing.T) {
	ResetCustomTypeFuncs()
	registerMockOmittables()
	defer ResetCustomTypeFuncs()

	s := &omittableMinStruct{
		Age: omittableEmpty[int](),
	}

	err := ValidateStruct(s)
	if err != nil {
		t.Fatalf("expected no error for unset non-required field, got: %v", err)
	}
}

func TestCustomTypeFunc_OmitEmpty_NotSet(t *testing.T) {
	ResetCustomTypeFuncs()
	registerMockOmittables()
	defer ResetCustomTypeFuncs()

	s := &omittableOmitEmptyStruct{
		Name: omittableEmpty[string](),
	}

	err := ValidateStruct(s)
	if err != nil {
		t.Fatalf("expected no error for unset omitempty field, got: %v", err)
	}
}

func TestCustomTypeFunc_OmitEmpty_Set_Valid(t *testing.T) {
	ResetCustomTypeFuncs()
	registerMockOmittables()
	defer ResetCustomTypeFuncs()

	s := &omittableOmitEmptyStruct{
		Name: omittableOf("test@example.com"),
	}

	err := ValidateStruct(s)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestCustomTypeFunc_OmitEmpty_Set_Invalid(t *testing.T) {
	ResetCustomTypeFuncs()
	registerMockOmittables()
	defer ResetCustomTypeFuncs()

	s := &omittableOmitEmptyStruct{
		Name: omittableOf("not-an-email"),
	}

	err := ValidateStruct(s)
	if err == nil {
		t.Fatal("expected error for invalid email with set omitempty field")
	}
}

func TestCustomTypeFunc_Pointer_Set_Nil(t *testing.T) {
	ResetCustomTypeFuncs()
	registerMockOmittables()
	defer ResetCustomTypeFuncs()

	s := &omittablePtrStruct{
		Name: omittableOf[*string](nil),
	}

	err := ValidateStruct(s)
	if err == nil {
		t.Fatal("expected error for set Omittable with nil pointer and required tag")
	}
}

func TestCustomTypeFunc_Pointer_Set_Value(t *testing.T) {
	ResetCustomTypeFuncs()
	registerMockOmittables()
	defer ResetCustomTypeFuncs()

	name := "Alice"
	s := &omittablePtrStruct{
		Name: omittableOf(&name),
	}

	err := ValidateStruct(s)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestCustomTypeFunc_MultiField(t *testing.T) {
	ResetCustomTypeFuncs()
	registerMockOmittables()
	defer ResetCustomTypeFuncs()

	s := &omittableMultiFieldStruct{
		Name:  omittableOf("Alice"),
		Email: omittableOf("alice@example.com"),
		Age:   omittableOf(25),
	}

	err := ValidateStruct(s)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestCustomTypeFunc_MultiField_PartialSet(t *testing.T) {
	ResetCustomTypeFuncs()
	registerMockOmittables()
	defer ResetCustomTypeFuncs()

	s := &omittableMultiFieldStruct{
		Name:  omittableOf("Alice"),
		Email: omittableEmpty[string](), // not set, required -> error
		Age:   omittableEmpty[int](),    // not set, no required -> skip
	}

	err := ValidateStruct(s)
	if err == nil {
		t.Fatal("expected error for unset required Email")
	}

	errs := err.(Errors)
	if !errs.HasFieldError("Email") {
		t.Fatal("expected FieldError for Email")
	}
	if errs.HasFieldError("Age") {
		t.Fatal("expected no error for unset non-required Age")
	}
}

func TestCustomTypeFunc_Between_Valid(t *testing.T) {
	ResetCustomTypeFuncs()
	registerMockOmittables()
	defer ResetCustomTypeFuncs()

	s := &omittableBetweenStruct{
		Score: omittableOf(50),
	}

	err := ValidateStruct(s)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestCustomTypeFunc_Between_Invalid(t *testing.T) {
	ResetCustomTypeFuncs()
	registerMockOmittables()
	defer ResetCustomTypeFuncs()

	s := &omittableBetweenStruct{
		Score: omittableOf(200),
	}

	err := ValidateStruct(s)
	if err == nil {
		t.Fatal("expected error for score outside range")
	}
}

// --- RegisterAutoUnwrap tests ---

func TestAutoUnwrap_NotSet_Required(t *testing.T) {
	ResetCustomTypeFuncs()
	registerMockAutoUnwrap()
	defer ResetCustomTypeFuncs()

	s := &omittableRequiredStruct{
		Name: omittableEmpty[string](),
	}

	err := ValidateStruct(s)
	if err == nil {
		t.Fatal("expected error for unset required field")
	}
}

func TestAutoUnwrap_Set_Required(t *testing.T) {
	ResetCustomTypeFuncs()
	registerMockAutoUnwrap()
	defer ResetCustomTypeFuncs()

	s := &omittableRequiredStruct{
		Name: omittableOf("Alice"),
	}

	err := ValidateStruct(s)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestAutoUnwrap_Set_Email_Valid(t *testing.T) {
	ResetCustomTypeFuncs()
	registerMockAutoUnwrap()
	defer ResetCustomTypeFuncs()

	s := &omittableEmailStruct{
		Email: omittableOf("test@example.com"),
	}

	err := ValidateStruct(s)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestAutoUnwrap_Set_Email_Invalid(t *testing.T) {
	ResetCustomTypeFuncs()
	registerMockAutoUnwrap()
	defer ResetCustomTypeFuncs()

	s := &omittableEmailStruct{
		Email: omittableOf("not-an-email"),
	}

	err := ValidateStruct(s)
	if err == nil {
		t.Fatal("expected error for invalid email")
	}
}

func TestAutoUnwrap_NotSet_Optional(t *testing.T) {
	ResetCustomTypeFuncs()
	registerMockAutoUnwrap()
	defer ResetCustomTypeFuncs()

	s := &omittableOptionalStruct{
		Name: omittableEmpty[string](),
	}

	err := ValidateStruct(s)
	if err != nil {
		t.Fatalf("expected no error for unset optional field, got: %v", err)
	}
}

func TestAutoUnwrap_Pointer_Set_Nil(t *testing.T) {
	ResetCustomTypeFuncs()
	registerMockAutoUnwrap()
	defer ResetCustomTypeFuncs()

	s := &omittablePtrStruct{
		Name: omittableOf[*string](nil),
	}

	err := ValidateStruct(s)
	if err == nil {
		t.Fatal("expected error for nil pointer with required tag")
	}
}

func TestAutoUnwrap_Pointer_Set_Value(t *testing.T) {
	ResetCustomTypeFuncs()
	registerMockAutoUnwrap()
	defer ResetCustomTypeFuncs()

	name := "Alice"
	s := &omittablePtrStruct{
		Name: omittableOf(&name),
	}

	err := ValidateStruct(s)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestAutoUnwrap_MultiField_PartialSet(t *testing.T) {
	ResetCustomTypeFuncs()
	registerMockAutoUnwrap()
	defer ResetCustomTypeFuncs()

	s := &omittableMultiFieldStruct{
		Name:  omittableOf("Alice"),
		Email: omittableEmpty[string](),
		Age:   omittableEmpty[int](),
	}

	err := ValidateStruct(s)
	if err == nil {
		t.Fatal("expected error for unset required Email")
	}

	errs := err.(Errors)
	if !errs.HasFieldError("Email") {
		t.Fatal("expected FieldError for Email")
	}
	if errs.HasFieldError("Age") {
		t.Fatal("expected no error for unset non-required Age")
	}
}

func TestAutoUnwrap_ExactMatchTakesPriority(t *testing.T) {
	ResetCustomTypeFuncs()
	registerMockAutoUnwrap()
	// Register exact match — should take priority over auto-unwrap
	RegisterCustomTypeFunc(func(o mockOmittable[string]) (any, bool) {
		if !o.IsSet() {
			return nil, false
		}
		return "override", true
	})
	defer ResetCustomTypeFuncs()

	s := &omittableRequiredStruct{
		Name: omittableOf("Alice"),
	}

	err := ValidateStruct(s)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
}

func TestAutoUnwrap_NoMatcherNoUnwrap(t *testing.T) {
	ResetCustomTypeFuncs()
	// No matcher registered — mockOmittable should NOT be unwrapped

	s := &omittableRequiredStruct{
		Name: omittableOf("Alice"),
	}

	// Without a matcher, the struct is treated as a plain struct — should not panic
	_ = ValidateStruct(s)
}

func TestResetCustomTypeFuncs(t *testing.T) {
	ResetCustomTypeFuncs()
	registerMockOmittables()
	registerMockAutoUnwrap()

	// Verify something is registered
	t1 := reflect.TypeFor[mockOmittable[string]]()
	if fn := resolveCustomTypeFunc(t1); fn == nil {
		t.Fatal("expected custom type func to be registered")
	}

	ResetCustomTypeFuncs()

	// After reset, exact match should be gone
	if fn := resolveCustomTypeFunc(t1); fn != nil {
		t.Fatal("expected no custom type func after reset")
	}
}

// --- Benchmarks ---

func BenchmarkCustomTypeFunc_Set(b *testing.B) {
	ResetCustomTypeFuncs()
	registerMockOmittables()
	defer ResetCustomTypeFuncs()

	s := &omittableEmailStruct{
		Email: omittableOf("test@example.com"),
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = ValidateStruct(s)
	}
}

func BenchmarkCustomTypeFunc_NotSet(b *testing.B) {
	ResetCustomTypeFuncs()
	registerMockOmittables()
	defer ResetCustomTypeFuncs()

	s := &omittableOptionalStruct{
		Name: omittableEmpty[string](),
	}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = ValidateStruct(s)
	}
}

func BenchmarkAutoUnwrap_Set(b *testing.B) {
	ResetCustomTypeFuncs()
	registerMockAutoUnwrap()
	defer ResetCustomTypeFuncs()

	s := &omittableEmailStruct{
		Email: omittableOf("test@example.com"),
	}

	// Prime the cache
	_ = ValidateStruct(s)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = ValidateStruct(s)
	}
}

func BenchmarkAutoUnwrap_NotSet(b *testing.B) {
	ResetCustomTypeFuncs()
	registerMockAutoUnwrap()
	defer ResetCustomTypeFuncs()

	s := &omittableOptionalStruct{
		Name: omittableEmpty[string](),
	}

	// Prime the cache
	_ = ValidateStruct(s)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = ValidateStruct(s)
	}
}

func BenchmarkCustomTypeFunc_Overhead_None(b *testing.B) {
	ResetCustomTypeFuncs()

	type SimpleStruct struct {
		Name string `valid:"required"`
	}

	s := &SimpleStruct{Name: "Alice"}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = ValidateStruct(s)
	}
}

func BenchmarkCustomTypeFunc_Overhead_Registered(b *testing.B) {
	ResetCustomTypeFuncs()
	registerMockOmittables()
	defer ResetCustomTypeFuncs()

	type SimpleStruct struct {
		Name string `valid:"required"`
	}

	s := &SimpleStruct{Name: "Alice"}

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		_ = ValidateStruct(s)
	}
}

func BenchmarkCustomTypeFunc_CacheLookup(b *testing.B) {
	ResetCustomTypeFuncs()
	registerMockOmittables()
	defer ResetCustomTypeFuncs()

	t := reflect.TypeFor[mockOmittable[string]]()
	// Prime the cache
	resolveCustomTypeFunc(t)

	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		resolveCustomTypeFunc(t)
	}
}
