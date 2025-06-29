package validator

import "fmt"

// ValidateDigitsBetweenUint64 returns true if value lies between left and right border
func ValidateDigitsBetweenUint64(value, left, right uint64) bool {
	if left > right {
		left, right = right, left
	}
	return value >= left && value <= right
}

// compareUint64 determine if a comparison passes between the given values.
func compareUint64(first uint64, second uint64, operator string) (bool, error) {
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
