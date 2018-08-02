package validator

import (
	"fmt"
	"reflect"
	"strconv"
)

// ToString convert the input to a string.
func ToString(obj interface{}) string {
	res := fmt.Sprintf("%v", obj)
	return string(res)
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
	val := reflect.ValueOf(value)

	switch value.(type) {
	case int, int8, int16, int32, int64:
		res = val.Int()
	case uint, uint8, uint16, uint32, uint64:
		res = int64(val.Uint())
	case string:
		if IsInt(val.String()) {
			res, err = strconv.ParseInt(val.String(), 0, 64)
			if err != nil {
				res = 0
			}
		} else {
			err = fmt.Errorf("math: square root of negative number %g", value)
			res = 0
		}
	default:

		err = fmt.Errorf("math: square root of negative number %g", value)
		res = 0
	}

	return
}

// ToUint convert the input string if the input is not an unit.
func ToUint(param string) (res uint64, err error) {
	i, err := strconv.ParseUint(param, 0, 64)
	return i, err
}
