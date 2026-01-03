package validator

import "fmt"

// IsUint64Between returns true if value lies between min and max (inclusive)
func IsUint64Between(value, min, max uint64) bool {
	if min > max {
		min, max = max, min
	}
	return value >= min && value <= max
}

// IsUint64Gt returns true if value is greater than threshold
func IsUint64Gt(value, threshold uint64) bool {
	return value > threshold
}

// IsUint64Gte returns true if value is greater than or equal to threshold
func IsUint64Gte(value, threshold uint64) bool {
	return value >= threshold
}

// IsUint64Lt returns true if value is less than threshold
func IsUint64Lt(value, threshold uint64) bool {
	return value < threshold
}

// IsUint64Lte returns true if value is less than or equal to threshold
func IsUint64Lte(value, threshold uint64) bool {
	return value <= threshold
}

// Aliases for consistency with other validators.
var (
	IsUint64Min = IsUint64Gte
	IsUint64Max = IsUint64Lte
)

// Deprecated: Use Is* functions instead.
var (
	ValidateDigitsBetweenUint64 = IsUint64Between
)

// compareUint64 determine if a comparison passes between the given values.
func compareUint64(first, second uint64, operator string) (bool, error) {
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
		return false, fmt.Errorf("validator: compareUint64 unsupported operator %s", operator)
	}
}
