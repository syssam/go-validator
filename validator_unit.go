package validator

// DigitsBetweenUint64 returns true if value lies between left and right border
func DigitsBetweenUint64(value, left, right uint64) bool {
	if left > right {
		left, right = right, left
	}
	return value >= left && value <= right
}

// MaxUnit64 is the validation function for validating if the current field's value is less than or equal to the param's value.
func MaxUnit64(v, param uint64) bool {
	return IsLteUnit64(v, param)
}

// MinUnit64 is the validation function for validating if the current field's value is greater than or equal to the param's value.
func MinUnit64(v, param uint64) bool {
	return IsGteUnit64(v, param)
}

// IsLtUnit64 is the validation function for validating if the current field's value is less than the param's value.
func IsLtUnit64(v, param uint64) bool {
	return v < param
}

// IsLteUnit64 is the validation function for validating if the current field's value is less than or equal to the param's value.
func IsLteUnit64(v, param uint64) bool {
	return v <= param
}

// IsGteUnit64 is the validation function for validating if the current field's value is greater than or equal to the param's value.
func IsGteUnit64(v, param uint64) bool {
	return v >= param
}

// IsGtUnit64 is the validation function for validating if the current field's value is greater than to the param's value.
func IsGtUnit64(v, param uint64) bool {
	return v > param
}
