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
	return LteInt64(v, param)
}

// MinInt64 is the validation function for validating if the current field's value is greater than or equal to the param's value.
func MinInt64(v, param int64) bool {
	return GteInt64(v, param)
}

// LtInt64 is the validation function for validating if the current field's value is less than the param's value.
func LtInt64(v, param int64) bool {
	return v < param
}

// LteInt64 is the validation function for validating if the current field's value is less than or equal to the param's value.
func LteInt64(v, param int64) bool {
	return v <= param
}

// GteInt64 is the validation function for validating if the current field's value is greater than or equal to the param's value.
func GteInt64(v, param int64) bool {
	return v >= param
}

// GtInt64 is the validation function for validating if the current field's value is greater than to the param's value.
func GtInt64(v, param int64) bool {
	return v > param
}
