package validator

import "fmt"

// ValidateDigitsBetweenInt64 returns true if value lies between left and right border
func ValidateDigitsBetweenInt64(value, left, right int64) bool {
	if left > right {
		left, right = right, left
	}
	return value >= left && value <= right
}

//  compareInt64 determine if a comparison passes between the given values.
func compareInt64(first int64, second int64, operator string) bool {
	switch operator {
	case "<":
		return first < second
	case ">":
		return first > second
	case "<=":
		return first <= second
	case ">=":
		return first >= second
	case "==":
		return first == second
	default:
		panic(fmt.Sprintf("validator: compareInt64 unsupport operator %s", operator))
	}
}
