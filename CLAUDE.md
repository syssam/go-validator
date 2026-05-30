# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go validation library (`go-validator`) that provides comprehensive validation for structs, strings, and other data types. It's optimized for Go 1.19+ and emphasizes proper error handling without panics. The library is designed for production use with web frameworks like Gin, Echo, and Iris.

## Key Architecture

### Core Components

- **Validator (`validator.go`)**: Main validation engine with struct tag-based validation
- **Error Handling (`error.go`)**: Comprehensive error types including `Errors`, `FieldError`, and `UnsupportedTypeError`
- **Types (`types.go`)**: Custom validation function types and thread-safe custom rule mapping
- **Cache (`cache.go`)**: Field metadata structures and validation tag parsing
- **Translator (`translator.go`)**: Internationalization support with language-specific error messages
- **Validation Rules**: Split across multiple files:
  - `validator_string.go` - String validation rules (email, alpha, numeric, etc.)
  - `validator_int.go` - Integer validation rules
  - `validator_float.go` - Float validation rules
  - `validator_unit.go` - Unit-specific validations

### Error Handling Architecture

The library uses a sophisticated error handling system:
- `Errors` type: Collection of multiple validation errors
- `FieldError` type: Detailed field-specific errors with optional `FuncError` for internal errors
- Error chaining support with `Unwrap()` method
- Utility methods: `HasFieldError()`, `GetFieldError()`, `GroupByField()`

### Validation Flow

1. **Entry Point**: `ValidateStruct(s interface{})` delegates to the default `Validator` instance
2. **Field Caching**: `cachedTypefields()` caches field metadata in thread-safe `sync.Map` to reduce reflection overhead
3. **Tag Parsing**: Parses `valid:"rule1,rule2=param"` struct tags into `ValidTag` structs
4. **Rule Dispatch**: Routes validation to the appropriate rule map based on rule type:
   - `RuleMap` - No-parameter rules (e.g., `distinct`)
   - `ParamRuleMap` - Parameter rules (e.g., `between=1|10`, `min=5`)
   - `StringRulesMap` - String-specific pattern rules (e.g., `email`, `alpha`, `uuid`)
   - `CustomTypeRuleMap` - User-defined custom validators (thread-safe)
5. **Error Collection**: All validation errors are accumulated in an `Errors` slice

## Development Commands

### Testing
```bash
go test                    # Run all tests
go test -v                 # Run tests with verbose output
go test -run TestName      # Run specific test
go test -bench=.           # Run benchmarks
go test -bench=. -benchmem # Run benchmarks with memory stats
go test -cover             # Run tests with coverage
go test -race              # Run tests with race detection
```

### Building & Validation
```bash
go build                   # Build the package
go mod tidy                # Clean up dependencies
go vet                     # Run go vet for static analysis
go fmt                     # Format code
gofmt -s -w .              # Format and simplify code
```

### Performance Testing
```bash
go test -bench=BenchmarkErrorHandling           # Test error handling performance
go test -bench=BenchmarkGo119Performance        # Test Go 1.19 optimizations
go test -bench=BenchmarkMemoryAllocation        # Test memory allocation patterns
go test -bench=BenchmarkStringBuilderOptimization # Test string building performance
go test -bench=BenchmarkFuncErrorHandling       # Test FuncError functionality performance
go test -bench=BenchmarkErrorUnwrapping         # Test error unwrapping performance
```

## Usage Patterns

### Basic Struct Validation
```go
type User struct {
    Name  string `valid:"required"`
    Email string `valid:"required,email"`
    Age   int    `valid:"min=18"`
}

err := validator.ValidateStruct(user)
if err != nil {
    errors := err.(validator.Errors)
    // Use utility methods for error handling
    if errors.HasFieldError("Email") {
        fieldError := errors.GetFieldError("Email")
        fmt.Println("Email error:", fieldError.Message)
    }
    
    // Group errors by field for organized display
    groupedErrors := errors.GroupByField()
    for field, errs := range groupedErrors {
        fmt.Printf("Field %s has %d errors\n", field, len(errs))
    }
}
```

### Nested, Embedded, and Recursive Structs

- **Nested structs** (named or pointer) are validated **recursively and
  automatically** — a field no longer needs its own `valid` tag for its inner
  fields to be checked. `time.Time`/`decimal.Decimal` are treated as scalar
  types, not recursed into.
- **Embedded (anonymous) structs** have their promoted fields validated as if
  declared on the parent. A nil embedded pointer is skipped (no panic).
- **Cyclic references** (a node pointing back to an ancestor) terminate via a
  depth guard (`maxValidationDepth`) that returns an error instead of looping
  forever. Validation is linear in the number of fields/elements.

### Fail-Fast vs Collect-All

By default the validator collects every error. Set `FailFast` on a `Validator`
instance to stop at the first field that fails and return immediately (similar to
Joi's `abortEarly` / Laravel's `stopOnFirstFailure`):

```go
v := validator.New()
v.FailFast = true // stop at the first failing field
err := v.ValidateStruct(user, nil, nil)
```

`FailFast` defaults to `false` (collect-all), preserving backward-compatible
behavior. Configure it once at setup, before concurrent validation.

### Custom Validation Rules
```go
validator.CustomTypeRuleMap.Set("customRule", func(v reflect.Value, o reflect.Value, validTag *validator.ValidTag) bool {
    // Custom validation logic
    return true
})

// Set custom error message
validator.MessageMap["customRule"] = "Custom validation message with {Value}"
```

### Error Handling with FuncError
```go
// The library supports error chaining for internal function errors
if fieldError.HasFuncError() {
    internalErr := fieldError.Unwrap()
    log.Printf("Internal validation error: %v", internalErr)
}
```

### Available Validation Rules

The authoritative rule registry lives in `types.go` (`RuleMap`, `ParamRuleMap`, `StringRulesMap`, `StringParamRulesMap`). When you add or rename a rule, update this list to match.

**Core / Presence Rules**: `required`, `filled`, `present`, `missing`, `prohibited`, `accepted`, `declined`

**Conditional Rules**: `requiredIf`, `requiredUnless`, `requiredWith`, `requiredWithAll`, `requiredWithout`, `requiredWithoutAll` *(implemented in validator logic)*

**Comparison Rules**: `gt`, `gte`, `lt`, `lte`, `same`, `distinct`, `multipleOf`

**Size Rules**: `min`, `max`, `between`, `size`, `digitsBetween`, `minDigits`, `maxDigits`, `decimal`

**Type Rules**: `int`, `integer`, `float`, `numeric`, `string`, `boolean`, `array`, `list`, `json`, `date`

**String Rules**: `email`, `alpha`, `alphaNum`, `alphaDash`, `alphaUnicode`, `alphaNumUnicode`, `alphaDashUnicode`, `ascii`, `lowercase`, `uppercase`, `hexColor`, `startsWith`, `endsWith`, `doesntStartWith`, `doesntEndWith`, `contains`, `doesntContain`

**Network Rules**: `ip`, `ipv4`, `ipv6`, `url`, `macAddress`, `uuid`, `uuid3`, `uuid4`, `uuid5`, `ulid`

**Pattern Rules**: `regex`, `notRegex`, `in`, `notIn`

**Date Rules**: `dateFormat`, `after`, `afterOrEqual`, `before`, `beforeOrEqual`, `timezone`

**Locale Rules**: `country`, `currency`, `language`, `phone` *(see `validator_country.go`, `validator_currency.go`, `validator_language.go`, `validator_phone.go`; locale variants use a dot separator, e.g. `country.alpha2`)*

**Note**: All regex patterns are defined in `patterns.go`. Some conditional rules are implemented in the main validation logic but may require specific struct field relationships to function.

## Framework Integration Examples

The `_examples/` directory contains production-ready integration examples:

### Web Framework Examples
- **`simple/`**: Basic validation with comprehensive error handling
- **`gin/`**: Gin framework integration with JSON API error responses
- **`echo/`**: Echo framework integration with custom validators
- **`iris/`**: Iris framework integration
- **`translation/`**: Simple internationalization example
- **`translations/`**: Advanced multi-language error messages
- **`custom/`**: Custom validation rule implementation

### Key Integration Patterns
```go
// Gin integration with localized errors
func ValidateJSON(c *gin.Context, obj interface{}) error {
    if err := c.ShouldBindJSON(obj); err != nil {
        return err
    }
    
    // Apply custom validation
    if err := validator.ValidateStruct(obj); err != nil {
        // Convert to API-friendly error format
        return formatValidationErrors(err.(validator.Errors))
    }
    return nil
}
```

## Testing Strategy

### Test Organization
- **`validator_test.go`**: Core validation logic tests (1100+ lines)
- **`benchmarks_test.go`**: Performance benchmarks for all major operations
- **Error handling tests**: Edge cases, Go 1.19 features, FuncError chaining
- **Framework examples**: Real-world integration testing

### Test Categories
- **Unit tests**: Individual validation rules and error handling (`validator_test.go` - 1,177 lines)
- **Integration tests**: Struct validation with complex nested types
- **Performance tests**: Memory allocation and execution time benchmarks (`benchmarks_test.go` - 10 benchmark functions)
- **Error chain tests**: FuncError unwrapping and error interface compatibility
- **Framework integration**: Real-world examples in `_examples/` directory

## Language Support & Internationalization

### Built-in Languages
- **English (`lang/en/`)**: Default language with comprehensive error messages
- **Chinese Simplified (`lang/zh_CN/`)**: Complete translation set
- **Chinese Traditional (`lang/zh_HK/`)**: Regional variant support

### Custom Language Implementation
```go
// Register a translator and load per-language message sets.
translator := validator.NewTranslator()
translator.SetMessage("fr", validator.Translate{
    "required": "{{.Attribute}} est requis",
})

// Translate a collected Errors set into the target language.
if err := validator.ValidateStruct(obj); err != nil {
    errs := translator.Trans(err.(validator.Errors), "fr")
    // errs now carry localized messages
}
```

> **Thread-safety note:** the package-level configuration maps — `MessageMap`,
> `RuleMap`, `ParamRuleMap`, `StringRulesMap`, and `Mimes` — are plain maps with
> no internal locking. Configure them **once at init, before** any concurrent
> `ValidateStruct` calls. Only `CustomTypeRuleMap` is safe to mutate at runtime.
> `*Translator` instances are read-only after setup and safe for concurrent use.

## Performance Characteristics

### Go 1.19+ Optimizations
- **String operations**: `strings.Builder` with pre-allocation (`builder.Grow()`) for efficient error message construction
- **Object pooling**: `sync.Pool` for ErrorResponse slices in JSON marshaling to reduce GC pressure
- **Slice pre-allocation**: Strategic capacity pre-sizing to minimize reallocations
- **Runtime benefits**: Leverages Go 1.19+ improved garbage collector and memory management

### Benchmark Results Focus Areas
- **Error handling**: Reduced allocations through efficient string building and object reuse
- **JSON marshaling**: Object pooling reduces allocation overhead for API responses
- **String concatenation**: Improved error message building with pre-allocated builders
- **Memory profiling**: Allocation patterns monitored via `go test -benchmem`

## Development Best Practices

### Code Organization
- Validation rules are logically separated by type (`validator_string.go`, `validator_int.go`, etc.)
- Error types support both simple usage and advanced error chaining
- Custom validation functions should be thread-safe and stateless
- Use the existing patterns for parameter validation and error creation

### Performance Guidelines
- Take advantage of object pooling patterns shown in `error.go` for high-frequency operations
- Use benchmark tests when adding new validation rules to measure allocation impact
- Pre-allocate slices with known capacity to reduce reallocations
- Test with `-race` flag for concurrent usage validation
- Consider using `strings.Builder` with `Grow()` for string concatenation in custom validators

## Adding New Validation Rules

1. **No-parameter rules** (e.g., `distinct`): Add to `RuleMap` in `types.go`
2. **Parameter rules** (e.g., `min=value`): Add to `ParamRuleMap` in `types.go`
3. **String pattern rules** (e.g., `email`): Add regex to `patterns.go`, function to `validator_string.go`, register in `StringRulesMap`
4. **Error messages**: Add default message to `MessageMap` in `message.go`
5. **Tests**: Add test cases in `validator_test.go`
6. **Benchmarks**: Add performance tests in `benchmarks_test.go` if the rule involves complex logic