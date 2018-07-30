package validator

import "fmt"

// DigitsBetweenUint64 returns true if value lies between left and right border
func DigitsBetweenUint64(value, left, right uint64) bool {
	if left > right {
		left, right = right, left
	}
	return value >= left && value <= right
}

// MaxUnit64 is the validation function for validating if the current field's value is less than or equal to the param's value.
func MaxUnit64(v, param uint64) bool {
	return LteUnit64(v, param)
}

// MinUnit64 is the validation function for validating if the current field's value is greater than or equal to the param's value.
func MinUnit64(v, param uint64) bool {
	return GteUnit64(v, param)
}

// LtUnit64 is the validation function for validating if the current field's value is less than the param's value.
func LtUnit64(v, param uint64) bool {
	return v < param
}

// LteUnit64 is the validation function for validating if the current field's value is less than or equal to the param's value.
func LteUnit64(v, param uint64) bool {
	return v <= param
}

// GteUnit64 is the validation function for validating if the current field's value is greater than or equal to the param's value.
func GteUnit64(v, param uint64) bool {
	return v >= param
}

// GtUnit64 is the validation function for validating if the current field's value is greater than to the param's value.
func GtUnit64(v, param uint64) bool {
	return v > param
}

//  compareUnit64 determine if a comparison passes between the given values.
func compareUnit64(first uint64, second uint64, operator string) bool {
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
		panic(fmt.Sprintf("validator: compareUnit64 unsupport operator %s", operator))
	}
}
