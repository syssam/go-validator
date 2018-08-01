package validator

import (
	"fmt"
	"unicode/utf8"
)

// BetweenString is
func BetweenString(v string, left int64, right int64) bool {
	return DigitsBetweenInt64(int64(utf8.RuneCountInString(v)), left, right)
}

// InString check if string str is a member of the set of strings params
func InString(str string, params ...string) bool {
	for _, param := range params {
		if str == param {
			return true
		}
	}

	return false
}

//  compareString determine if a comparison passes between the given values.
func compareString(first string, second int64, operator string) bool {
	switch operator {
	case "<":
		return int64(utf8.RuneCountInString(first)) < second
	case ">":
		return int64(utf8.RuneCountInString(first)) > second
	case "<=":
		return int64(utf8.RuneCountInString(first)) <= second
	case ">=":
		return int64(utf8.RuneCountInString(first)) >= second
	case "==":
		return int64(utf8.RuneCountInString(first)) == second
	default:
		panic(fmt.Sprintf("validator: compareString unsupport operator %s", operator))
	}
}
