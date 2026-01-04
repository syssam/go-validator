package validator

import (
	"fmt"
	"reflect"

	"github.com/shopspring/decimal"
)

// decimalType is the reflect.Type for decimal.Decimal (cached for performance).
var decimalType = reflect.TypeOf(decimal.Decimal{})

// asDecimal checks if v is a decimal.Decimal and returns it.
// Returns (value, true) if v is a decimal, (zero, false) otherwise.
func asDecimal(v reflect.Value) (decimal.Decimal, bool) {
	if v.Kind() == reflect.Struct && v.Type() == decimalType {
		return v.Interface().(decimal.Decimal), true
	}
	return decimal.Decimal{}, false
}

// parseDecimal converts a string to decimal.Decimal.
func parseDecimal(s string) (decimal.Decimal, error) {
	return decimal.NewFromString(s)
}

// parseDecimalParams parses one or two string parameters into decimals.
func parseDecimalParams(params []string) (decimal.Decimal, decimal.Decimal, error) {
	if len(params) == 0 {
		return decimal.Decimal{}, decimal.Decimal{}, fmt.Errorf("no parameters provided")
	}

	first, err := parseDecimal(params[0])
	if err != nil {
		return decimal.Decimal{}, decimal.Decimal{}, fmt.Errorf("invalid first parameter: %w", err)
	}

	if len(params) == 1 {
		return first, decimal.Decimal{}, nil
	}

	second, err := parseDecimal(params[1])
	if err != nil {
		return decimal.Decimal{}, decimal.Decimal{}, fmt.Errorf("invalid second parameter: %w", err)
	}

	return first, second, nil
}

// compareDecimal compares two decimals using the given operator.
func compareDecimal(first, second decimal.Decimal, operator string) (bool, error) {
	switch operator {
	case "<":
		return first.LessThan(second), nil
	case ">":
		return first.GreaterThan(second), nil
	case "<=":
		return first.LessThanOrEqual(second), nil
	case ">=":
		return first.GreaterThanOrEqual(second), nil
	case "==":
		return first.Equal(second), nil
	default:
		return false, fmt.Errorf("validator: compareDecimal unsupported operator %s", operator)
	}
}

// Decimal validation functions - single responsibility, no duplication.

// IsDecimalZero checks if decimal is zero.
func IsDecimalZero(value decimal.Decimal) bool {
	return value.IsZero()
}

// IsDecimalPositive checks if value > 0.
func IsDecimalPositive(value decimal.Decimal) bool {
	return value.IsPositive()
}

// IsDecimalNegative checks if value < 0.
func IsDecimalNegative(value decimal.Decimal) bool {
	return value.IsNegative()
}

// IsDecimalEqual checks if two decimals are equal.
func IsDecimalEqual(a, b decimal.Decimal) bool {
	return a.Equal(b)
}

// IsDecimalGt checks if a > b.
func IsDecimalGt(a, b decimal.Decimal) bool {
	return a.GreaterThan(b)
}

// IsDecimalGte checks if a >= b.
func IsDecimalGte(a, b decimal.Decimal) bool {
	return a.GreaterThanOrEqual(b)
}

// IsDecimalLt checks if a < b.
func IsDecimalLt(a, b decimal.Decimal) bool {
	return a.LessThan(b)
}

// IsDecimalLte checks if a <= b.
func IsDecimalLte(a, b decimal.Decimal) bool {
	return a.LessThanOrEqual(b)
}

// IsDecimalBetween checks if min <= value <= max.
func IsDecimalBetween(value, min, max decimal.Decimal) bool {
	return value.GreaterThanOrEqual(min) && value.LessThanOrEqual(max)
}

// Aliases for consistency with other validators.
var (
	IsDecimalRequired      = func(v decimal.Decimal) bool { return !v.IsZero() }
	IsDecimalMin           = IsDecimalGte
	IsDecimalMax           = IsDecimalLte
	ValidateDecimalMin     = IsDecimalGte
	ValidateDecimalMax     = IsDecimalLte
	ValidateDecimalGt      = IsDecimalGt
	ValidateDecimalGte     = IsDecimalGte
	ValidateDecimalLt      = IsDecimalLt
	ValidateDecimalLte     = IsDecimalLte
	ValidateDecimalZero    = IsDecimalZero
	ValidateDecimalBetween = IsDecimalBetween
)

// Deprecated aliases for backward compatibility.
var (
	ValidateDecimalPositive = IsDecimalPositive
	ValidateDecimalNegative = IsDecimalNegative
)
