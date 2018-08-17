package validator

import "fmt"

// ValidateDigitsBetweenUint64 returns true if value lies between left and right border
func ValidateDigitsBetweenUint64(value, left, right uint64) bool {
	if left > right {
		left, right = right, left
	}
	return value >= left && value <= right
}

//  compareUint64 determine if a comparison passes between the given values.
func compareUint64(first uint64, second uint64, operator string) bool {
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
		panic(fmt.Sprintf("validator: compareUint64 unsupport operator %s", operator))
	}
}
