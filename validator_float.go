package validator

// DigitsBetweenFloat64 returns true if value lies between left and right border
func DigitsBetweenFloat64(value, left, right float64) bool {
	if left > right {
		left, right = right, left
	}
	return value >= left && value <= right
}

// MaxFloat64 is the validation function for validating if the current field's value is less than or equal to the param's value.
func MaxFloat64(v, param float64) bool {
	return IsLteFloat64(v, param)
}

// MinFloat64 is the validation function for validating if the current field's value is greater than or equal to the param's value.
func MinFloat64(v, param float64) bool {
	return IsGteFloat64(v, param)
}

// IsLtFloat64 is the validation function for validating if the current field's value is less than the param's value.
func IsLtFloat64(v, param float64) bool {
	return v < param
}

// IsLteFloat64 is the validation function for validating if the current field's value is less than or equal to the param's value.
func IsLteFloat64(v, param float64) bool {
	return v <= param
}

// IsGteFloat64 is the validation function for validating if the current field's value is greater than or equal to the param's value.
func IsGteFloat64(v, param float64) bool {
	return v >= param
}

// IsGtFloat64 is the validation function for validating if the current field's value is greater than to the param's value.
func IsGtFloat64(v, param float64) bool {
	return v > param
}
