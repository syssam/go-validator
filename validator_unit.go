package validator

import "fmt"

// DigitsBetweenUint64 returns true if value lies between left and right border
func DigitsBetweenUint64(value, left, right uint64) bool {
	if left > right {
		left, right = right, left
	}
	return value >= left && value <= right
}

//  compareUint64 determine if a comparison passes between the given values.
func compareUint64(first uint64, second uint64, operator string) bool {
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
		panic(fmt.Sprintf("validator: compareUint64 unsupport operator %s", operator))
	}
}

// DistinctUint is the validation function for validating an attribute is unique among other values.
func DistinctUint(v []uint) bool {
	return inArrayUint(v, v)
}

// DistinctUint8 is the validation function for validating an attribute is unique among other values.
func DistinctUint8(v []uint8) bool {
	return inArrayUint8(v, v)
}

// DistinctUint16 is the validation function for validating an attribute is unique among other values.
func DistinctUint16(v []uint16) bool {
	return inArrayUint16(v, v)
}

// DistinctUint32 is the validation function for validating an attribute is unique among other values.
func DistinctUint32(v []uint32) bool {
	return inArrayUint32(v, v)
}

// DistinctUint64 is the validation function for validating an attribute is unique among other values.
func DistinctUint64(v []uint64) bool {
	return inArrayUint64(v, v)
}

func inArrayUint(needle []uint, haystack []uint) bool {
	for _, n := range needle {
		for _, s := range haystack {
			if n == s {
				return true
			}
		}
	}

	return false
}

func inArrayUint8(needle []uint8, haystack []uint8) bool {
	for _, n := range needle {
		for _, s := range haystack {
			if n == s {
				return true
			}
		}
	}

	return false
}

func inArrayUint16(needle []uint16, haystack []uint16) bool {
	for _, n := range needle {
		for _, s := range haystack {
			if n == s {
				return true
			}
		}
	}

	return false
}

func inArrayUint32(needle []uint32, haystack []uint32) bool {
	for _, n := range needle {
		for _, s := range haystack {
			if n == s {
				return true
			}
		}
	}

	return false
}

func inArrayUint64(needle []uint64, haystack []uint64) bool {
	for _, n := range needle {
		for _, s := range haystack {
			if n == s {
				return true
			}
		}
	}

	return false
}
