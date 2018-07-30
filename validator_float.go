package validator

import "fmt"

// DigitsBetweenFloat64 returns true if value lies between left and right border
func DigitsBetweenFloat64(value, left, right float64) bool {
	if left > right {
		left, right = right, left
	}
	return value >= left && value <= right
}

// MaxFloat64 is the validation function for validating if the current field's value is less than or equal to the param's value.
func MaxFloat64(v, param float64) bool {
	return LteFloat64(v, param)
}

// MinFloat64 is the validation function for validating if the current field's value is greater than or equal to the param's value.
func MinFloat64(v, param float64) bool {
	return GteFloat64(v, param)
}

// LtFloat64 is the validation function for validating if the current field's value is less than the param's value.
func LtFloat64(v, param float64) bool {
	return v < param
}

// LteFloat64 is the validation function for validating if the current field's value is less than or equal to the param's value.
func LteFloat64(v, param float64) bool {
	return v <= param
}

// GteFloat64 is the validation function for validating if the current field's value is greater than or equal to the param's value.
func GteFloat64(v, param float64) bool {
	return v >= param
}

// GtFloat64 is the validation function for validating if the current field's value is greater than to the param's value.
func GtFloat64(v, param float64) bool {
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
