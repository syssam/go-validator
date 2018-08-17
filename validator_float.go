package validator

import "fmt"

// ValidateDigitsBetweenFloat64 returns true if value lies between left and right border
func ValidateDigitsBetweenFloat64(value, left, right float64) bool {
	if left > right {
		left, right = right, left
	}
	return value >= left && value <= right
}

// ValidateMaxFloat64 is the validation function for validating if the current field's value is less than or equal to the param's value.
func ValidateMaxFloat64(v, param float64) bool {
	return ValidateLteFloat64(v, param)
}

// ValidateMinFloat64 is the validation function for validating if the current field's value is greater than or equal to the param's value.
func ValidateMinFloat64(v, param float64) bool {
	return ValidateGteFloat64(v, param)
}

// ValidateLtFloat64 is the validation function for validating if the current field's value is less than the param's value.
func ValidateLtFloat64(v, param float64) bool {
	return v < param
}

// ValidateLteFloat64 is the validation function for validating if the current field's value is less than or equal to the param's value.
func ValidateLteFloat64(v, param float64) bool {
	return v <= param
}

// ValidateGteFloat64 is the validation function for validating if the current field's value is greater than or equal to the param's value.
func ValidateGteFloat64(v, param float64) bool {
	return v >= param
}

// ValidateGtFloat64 is the validation function for validating if the current field's value is greater than to the param's value.
func ValidateGtFloat64(v, param float64) bool {
	return v > param
}

//  compareFloat64 determine if a comparison passes between the given values.
func compareFloat64(first float64, second float64, operator string) bool {
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
		panic(fmt.Sprintf("validator: compareFloat64 unsupport operator %s", operator))
	}
}
