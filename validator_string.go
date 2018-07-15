package validator

import (
	"unicode/utf8"
)

// MaxString is the validation function for validating if the current field's value is less than or equal to the param's value.
func MaxString(v string, param int64) bool {
	return IsLteString(v, param)
}

// MinString is the validation function for validating if the current field's value is greater than or equal to the param's value.
func MinString(v string, param int64) bool {
	return IsGteString(v, param)
}

// IsLtString is the validation function for validating if the current field's value is less than the param's value.
func IsLtString(v string, param int64) bool {
	return int64(utf8.RuneCountInString(v)) < param
}

// IsLteString is the validation function for validating if the current field's value is less than or equal to the param's value.
func IsLteString(v string, param int64) bool {
	return int64(utf8.RuneCountInString(v)) <= param
}

// IsGteString is the validation function for validating if the current field's value is greater than or equal to the param's value.
func IsGteString(v string, param int64) bool {
	return int64(utf8.RuneCountInString(v)) >= param
}

// IsGtString is the validation function for validating if the current field's value is greater than to the param's value.
func IsGtString(v string, param int64) bool {
	return int64(utf8.RuneCountInString(v)) > param
}

// BetweenString is
func BetweenString(v string, left int64, right int64) bool {
	return DigitsBetweenInt64(int64(utf8.RuneCountInString(v)), left, right)
}
