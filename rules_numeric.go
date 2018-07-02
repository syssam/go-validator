package validator

import (
	"fmt"
	"reflect"
)

// DigitsBetween returns true if value lies between left and right border, generic type to handle int, float32 or float64, all types must the same type
func DigitsBetween(v reflect.Value, params ...string) bool {
	if len(params) != 2 {
		return false
	}

	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		min, _ := ToInt(params[0])
		max, _ := ToInt(params[1])
		return DigitsBetweenInt64(v.Int(), min, max)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		min, _ := ToUint(params[0])
		max, _ := ToUint(params[1])
		return DigitsBetweenUint64(v.Uint(), min, max)
	case reflect.Float32, reflect.Float64:
		min, _ := ToFloat(params[0])
		max, _ := ToFloat(params[1])
		return DigitsBetweenFloat64(v.Float(), min, max)
	}

	panic(fmt.Sprintf("validator: DigitsBetween unsupport Type %T", v.Interface()))
}
