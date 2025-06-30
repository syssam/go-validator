package validator

import (
	"fmt"
	"testing"
)

func BenchmarkFieldsRequired(t *testing.B) {
	model := FieldsRequired{Name: "TEST", Email: "test@example.com"}
	expected := true
	for i := 0; i < t.N; i++ {
		err := ValidateStruct(&model)
		actual := err == nil
		if actual != expected {
			t.Errorf("Expected validateateStruct(%q) to be %v, got %v", model, expected, actual)
			if err != nil {
				t.Errorf("Got Error on validateateStruct(%q): %s", model, err)
			}
		}
	}
}

// BenchmarkErrorHandling benchmarks the new error handling
func BenchmarkErrorHandling(t *testing.B) {
	type TestStruct struct {
		Name  string `valid:"required"`
		Email string `valid:"required,email"`
		Age   int    `valid:"min=18"`
	}

	invalidData := TestStruct{Name: "", Email: "invalid", Age: 16}

	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		_ = ValidateStruct(invalidData)
	}
}

// BenchmarkErrorsUtilityMethods benchmarks the Errors utility methods
func BenchmarkErrorsUtilityMethods(t *testing.B) {
	type TestStruct struct {
		Name  string `valid:"required"`
		Email string `valid:"required,email"`
		Age   int    `valid:"min=18"`
	}

	invalidData := TestStruct{Name: "", Email: "invalid", Age: 16}
	err := ValidateStruct(invalidData)
	//nolint:errcheck // Safe type assertion in benchmark context
	errors := err.(Errors)

	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		_ = errors.HasFieldError("Name")
		_ = errors.GetFieldError("Email")
		_ = errors.GroupByField()
	}
}

// BenchmarkComparisonFunctions benchmarks the comparison functions
func BenchmarkComparisonFunctions(t *testing.B) {
	type TestStruct struct {
		Value1 int `valid:"gt=Value2"`
		Value2 int
	}

	data := TestStruct{Value1: 10, Value2: 5}

	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		_ = ValidateStruct(data)
	}
}

// BenchmarkGo119Performance benchmarks performance with Go 1.19 optimizations
func BenchmarkGo119Performance(t *testing.B) {
	type LargeStruct struct {
		Field1  string `valid:"required"`
		Field2  string `valid:"email"`
		Field3  int    `valid:"min=1"`
		Field4  int    `valid:"max=100"`
		Field5  string `valid:"between=5|20"`
		Field6  string `valid:"alpha"`
		Field7  string `valid:"numeric"`
		Field8  string `valid:"url"`
		Field9  string `valid:"ip"`
		Field10 string `valid:"uuid"`
	}

	data := LargeStruct{
		Field1:  "test",
		Field2:  "test@example.com",
		Field3:  5,
		Field4:  50,
		Field5:  "medium length text",
		Field6:  "alphabet",
		Field7:  "123456",
		Field8:  "https://example.com",
		Field9:  "192.168.1.1",
		Field10: "550e8400-e29b-41d4-a716-446655440000",
	}

	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		_ = ValidateStruct(data)
	}
}

// BenchmarkMemoryAllocation benchmarks memory allocation patterns
func BenchmarkMemoryAllocation(t *testing.B) {
	type SimpleStruct struct {
		Name  string `valid:"required"`
		Email string `valid:"email"`
		Age   int    `valid:"min=18,max=120"`
	}

	data := SimpleStruct{
		Name:  "John Doe",
		Email: "john.doe@example.com",
		Age:   25,
	}

	t.ResetTimer()
	t.ReportAllocs()
	for i := 0; i < t.N; i++ {
		_ = ValidateStruct(data)
	}
}

// BenchmarkStringBuilderOptimization benchmarks the strings.Builder optimization in error handling
func BenchmarkStringBuilderOptimization(t *testing.B) {
	type MultiErrorStruct struct {
		Field1 string `valid:"required"`
		Field2 string `valid:"required"`
		Field3 string `valid:"required"`
		Field4 string `valid:"required"`
		Field5 string `valid:"required"`
	}

	// Create data that will generate multiple errors
	invalidData := MultiErrorStruct{} // All fields empty, will generate 5 errors

	t.ResetTimer()
	t.ReportAllocs()
	for i := 0; i < t.N; i++ {
		err := ValidateStruct(invalidData)
		if err != nil {
			// Force error string generation to test strings.Builder optimization
			_ = err.Error()
		}
	}
}

// BenchmarkErrorGrouping benchmarks the error grouping functionality
func BenchmarkErrorGrouping(t *testing.B) {
	type ComplexStruct struct {
		Name    string `valid:"required"`
		Email   string `valid:"required,email"`
		Age     int    `valid:"min=18"`
		Phone   string `valid:"required"`
		Address string `valid:"required,min=10"`
	}

	// Create invalid data to generate multiple errors
	invalidData := ComplexStruct{
		Name:    "",
		Email:   "invalid-email",
		Age:     15,
		Phone:   "",
		Address: "short",
	}

	err := ValidateStruct(invalidData)
	if err == nil {
		t.Fatal("Expected validation errors")
	}

	//nolint:errcheck // Safe type assertion in benchmark context
	errors := err.(Errors)

	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		_ = errors.GroupByField()
	}
}

// BenchmarkFuncErrorHandling benchmarks the FuncError functionality
func BenchmarkFuncErrorHandling(t *testing.B) {
	type ErrorStruct struct {
		Complex complex64 `valid:"between=1|10"` // Will trigger FuncError
		Valid   string    `valid:"required"`     // Normal validation
	}

	data := ErrorStruct{
		Complex: 5 + 5i, // Unsupported type
		Valid:   "test",
	}

	t.ResetTimer()
	t.ReportAllocs()
	for i := 0; i < t.N; i++ {
		err := ValidateStruct(data)
		if err != nil {
			//nolint:errcheck // Safe type assertion in benchmark context
			errors := err.(Errors)
			// Benchmark accessing FuncError methods
			for _, e := range errors {
				if fieldErr, ok := e.(*FieldError); ok {
					_ = fieldErr.HasFuncError()
					_ = fieldErr.Unwrap()
				}
			}
		}
	}
}

// BenchmarkErrorUnwrapping benchmarks error unwrapping performance
func BenchmarkErrorUnwrapping(t *testing.B) {
	// Create a FieldError with FuncError
	originalErr := fmt.Errorf("test function error")
	fieldError := &FieldError{
		Name:      "TestField",
		Message:   "Test message",
		FuncError: originalErr,
	}

	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		_ = fieldError.HasFuncError()
		_ = fieldError.Unwrap()
		_ = fieldError.Error()
	}
}

// BenchmarkPerformanceOptimizations benchmarks the new performance optimizations
func BenchmarkPerformanceOptimizations(t *testing.B) {
	type OptimizedStruct struct {
		Name   string `valid:"required"`
		Email  string `valid:"email"`
		Age    int    `valid:"min=18,max=120"`
		Score  int64  `valid:"between=0|100"`
		Active bool   `valid:"required"`
	}

	data := OptimizedStruct{
		Name:   "John Doe",
		Email:  "john.doe@example.com",
		Age:    25,
		Score:  85,
		Active: true,
	}

	t.ResetTimer()
	t.ReportAllocs()
	for i := 0; i < t.N; i++ {
		_ = ValidateStruct(data)
	}
}

// BenchmarkToStringOptimization benchmarks the optimized ToString function
func BenchmarkToStringOptimization(t *testing.B) {
	testValues := []interface{}{
		"string value",
		42,
		int64(123456789),
		uint64(987654321),
		3.14159,
		true,
	}

	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		for _, v := range testValues {
			_ = ToString(v)
		}
	}
}

// BenchmarkObjectPooling benchmarks object pooling performance
func BenchmarkObjectPooling(t *testing.B) {
	t.Run("AllocationPattern", func(t *testing.B) {
		for i := 0; i < t.N; i++ {
			fe := &FieldError{}
			fe.Name = "test"
			fe.Message = "test message"
			// Simulate usage
			_ = fe.Error()
		}
	})
}

// BenchmarkComplexValidation tests performance on complex nested structures
func BenchmarkComplexValidation(t *testing.B) {
	type Address struct {
		Street  string `valid:"required"`
		City    string `valid:"required"`
		Zip     string `valid:"numeric,size=5"`
		Country string `valid:"required,alpha"`
	}

	type User struct {
		Name      string    `valid:"required,alpha"`
		Email     string    `valid:"required,email"`
		Age       int       `valid:"min=18,max=120"`
		Score     float64   `valid:"between=0|100"`
		Active    bool      `valid:"required"`
		Addresses []Address `valid:"required"`
	}

	data := User{
		Name:   "JohnDoe",
		Email:  "john.doe@example.com",
		Age:    30,
		Score:  85.5,
		Active: true,
		Addresses: []Address{
			{
				Street:  "Main St",
				City:    "CityName",
				Zip:     "12345",
				Country: "USA",
			},
		},
	}

	t.ResetTimer()
	t.ReportAllocs()
	for i := 0; i < t.N; i++ {
		_ = ValidateStruct(data)
	}
}
