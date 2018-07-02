package validator

// DigitsBetweenInt64 returns true if value lies between left and right border
func DigitsBetweenInt64(value, left, right int64) bool {
	if left > right {
		left, right = right, left
	}
	return value >= left && value <= right
}

// MaxInt64 is the validation function for validating if the current field's value is less than or equal to the param's value.
func MaxInt64(v, param int64) bool {
	return IsLteInt64(v, param)
}

// MinInt64 is the validation function for validating if the current field's value is greater than or equal to the param's value.
func MinInt64(v, param int64) bool {
	return IsGteInt64(v, param)
}

// IsLtInt64 is the validation function for validating if the current field's value is less than the param's value.
func IsLtInt64(v, param int64) bool {
	return v < param
}

// IsLteInt64 is the validation function for validating if the current field's value is less than or equal to the param's value.
func IsLteInt64(v, param int64) bool {
	return v <= param
}

// IsGteInt64 is the validation function for validating if the current field's value is greater than or equal to the param's value.
func IsGteInt64(v, param int64) bool {
	return v >= param
}

// IsGtInt64 is the validation function for validating if the current field's value is greater than to the param's value.
func IsGtInt64(v, param int64) bool {
	return v > param
}
