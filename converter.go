package validator

import (
	"fmt"
	"math"
	"strconv"
)

// ToString convert the input to a string with optimized fast paths.
func ToString(obj interface{}) string {
	// Fast path for common types to avoid fmt.Sprintf overhead
	switch v := obj.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	default:
		// Fallback to fmt.Sprintf for complex types
		return fmt.Sprintf("%v", obj)
	}
}

// ToFloat convert the input string to a float, or 0.0 if the input is not a float.
func ToFloat(str string) (float64, error) {
	res, err := strconv.ParseFloat(str, 64)
	if err != nil {
		res = 0.0
	}
	return res, err
}

// ToBool convert the input string to a bool if the input is not a string.
func ToBool(str string) bool {
	if str == "true" || str == "1" {
		return true
	}

	return false
}

// ToInt convert the input string or any int type to an integer type 64, or 0 if the input is not an integer.
func ToInt(value interface{}) (res int64, err error) {
	switch v := value.(type) {
	case int:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case uint:
		if v > math.MaxInt64 {
			return 0, fmt.Errorf("value %v exceeds int64 maximum", v)
		}
		return int64(v), nil
	case uint8:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case uint64:
		if v > math.MaxInt64 {
			return 0, fmt.Errorf("value %v exceeds int64 maximum", v)
		}
		return int64(v), nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	default:
		return 0, fmt.Errorf("unable to convert %T to int64", value)
	}
}

// ToUint convert the input string if the input is not an unit.
func ToUint(param string) (res uint64, err error) {
	i, err := strconv.ParseUint(param, 0, 64)
	return i, err
}
