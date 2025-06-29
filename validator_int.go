package validator

import "fmt"

// ValidateDigitsBetweenInt64 returns true if value lies between left and right border
func ValidateDigitsBetweenInt64(value, left, right int64) bool {
	if left > right {
		left, right = right, left
	}
	return value >= left && value <= right
}

// compareInt64 determine if a comparison passes between the given values.
func compareInt64(first int64, second int64, operator string) (bool, error) {
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
