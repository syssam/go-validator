package validator

import (
	"unicode/utf8"
)

// MaxString is the validation function for validating if the current field's value is less than or equal to the param's value.
func MaxString(v string, param int64) bool {
	return LteString(v, param)
}

// MinString is the validation function for validating if the current field's value is greater than or equal to the param's value.
func MinString(v string, param int64) bool {
	return GteString(v, param)
}

// LtString is the validation function for validating if the current field's value is less than the param's value.
func LtString(v string, param int64) bool {
	return int64(utf8.RuneCountInString(v)) < param
}

// LteString is the validation function for validating if the current field's value is less than or equal to the param's value.
func LteString(v string, param int64) bool {
	return int64(utf8.RuneCountInString(v)) <= param
}

// GteString is the validation function for validating if the current field's value is greater than or equal to the param's value.
func GteString(v string, param int64) bool {
	return int64(utf8.RuneCountInString(v)) >= param
}

// GtString is the validation function for validating if the current field's value is greater than to the param's value.
func GtString(v string, param int64) bool {
	return int64(utf8.RuneCountInString(v)) > param
}

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
