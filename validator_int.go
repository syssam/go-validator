package validator

import "fmt"

// IsInt64Between returns true if value lies between min and max (inclusive)
func IsInt64Between(value, min, max int64) bool {
	if min > max {
		min, max = max, min
	}
	return value >= min && value <= max
}

// IsInt64Gt returns true if value is greater than threshold
func IsInt64Gt(value, threshold int64) bool {
	return value > threshold
}

// IsInt64Gte returns true if value is greater than or equal to threshold
func IsInt64Gte(value, threshold int64) bool {
	return value >= threshold
}

// IsInt64Lt returns true if value is less than threshold
func IsInt64Lt(value, threshold int64) bool {
	return value < threshold
}

// IsInt64Lte returns true if value is less than or equal to threshold
func IsInt64Lte(value, threshold int64) bool {
	return value <= threshold
}

// Aliases for consistency with other validators.
var (
	IsInt64Min = IsInt64Gte
	IsInt64Max = IsInt64Lte
)

// Deprecated: Use Is* functions instead.
var (
	ValidateDigitsBetweenInt64 = IsInt64Between
)

// compareInt64 determine if a comparison passes between the given values.
func compareInt64(first, second int64, operator string) (bool, error) {
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
		return false, fmt.Errorf("validator: compareInt64 unsupported operator %s", operator)
	}
}
