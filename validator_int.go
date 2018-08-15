package validator

import "fmt"

// DigitsBetweenInt64 returns true if value lies between left and right border
func DigitsBetweenInt64(value, left, right int64) bool {
	if left > right {
		left, right = right, left
	}
	return value >= left && value <= right
}

//  compareInt64 determine if a comparison passes between the given values.
func compareInt64(first int64, second int64, operator string) bool {
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
		panic(fmt.Sprintf("validator: compareInt64 unsupport operator %s", operator))
	}
}

func DistinctInt(v []int) bool {
	return inArrayInt(v, v)
}

func DistinctInt8(v []int8) bool {
	return inArrayInt8(v, v)
}

func DistinctInt16(v []int16) bool {
	return inArrayInt16(v, v)
}

func DistinctInt32(v []int32) bool {
	return inArrayInt32(v, v)
}

func DistinctInt64(v []int64) bool {
	return inArrayInt64(v, v)
}

func inArrayInt(needle []int, haystack []int) bool {
	for _, n := range needle {
		for _, s := range haystack {
			if n == s {
				return true
			}
		}
	}

	return false
}

func inArrayInt8(needle []int8, haystack []int8) bool {
	for _, n := range needle {
		for _, s := range haystack {
			if n == s {
				return true
			}
		}
	}

	return false
}

func inArrayInt16(needle []int16, haystack []int16) bool {
	for _, n := range needle {
		for _, s := range haystack {
			if n == s {
				return true
			}
		}
	}

	return false
}

func inArrayInt32(needle []int32, haystack []int32) bool {
	for _, n := range needle {
		for _, s := range haystack {
			if n == s {
				return true
			}
		}
	}

	return false
}

func inArrayInt64(needle []int64, haystack []int64) bool {
	for _, n := range needle {
		for _, s := range haystack {
			if n == s {
				return true
			}
		}
	}

	return false
}
