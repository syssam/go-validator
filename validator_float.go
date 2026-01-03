package validator

import "fmt"

// IsFloat64Between returns true if value lies between min and max (inclusive)
func IsFloat64Between(value, min, max float64) bool {
	if min > max {
		min, max = max, min
	}
	return value >= min && value <= max
}

// IsFloat64Gt returns true if value is greater than threshold
func IsFloat64Gt(value, threshold float64) bool {
	return value > threshold
}

// IsFloat64Gte returns true if value is greater than or equal to threshold
func IsFloat64Gte(value, threshold float64) bool {
	return value >= threshold
}

// IsFloat64Lt returns true if value is less than threshold
func IsFloat64Lt(value, threshold float64) bool {
	return value < threshold
}

// IsFloat64Lte returns true if value is less than or equal to threshold
func IsFloat64Lte(value, threshold float64) bool {
	return value <= threshold
}

// Aliases for consistency with other validators.
var (
	IsFloat64Min = IsFloat64Gte
	IsFloat64Max = IsFloat64Lte
)

// Deprecated: Use Is* functions instead.
var (
	ValidateDigitsBetweenFloat64 = IsFloat64Between
	ValidateMaxFloat64           = IsFloat64Lte
	ValidateMinFloat64           = IsFloat64Gte
	ValidateLtFloat64            = IsFloat64Lt
	ValidateLteFloat64           = IsFloat64Lte
	ValidateGteFloat64           = IsFloat64Gte
	ValidateGtFloat64            = IsFloat64Gt
)

// compareFloat64 determines if a comparison passes between the given values.
func compareFloat64(first, second float64, operator string) (bool, error) {
	switch operator {
	case "<":
		return first < second, nil
	case ">":
		return first > second, nil
	case "<=":
		return first <= second, nil
	case ">=":
		return first >= second, nil
	case "==":
		return first == second, nil
	default:
		return false, fmt.Errorf("validator: compareFloat64 unsupported operator %s", operator)
	}
}
