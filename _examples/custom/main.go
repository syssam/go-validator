package main

import (
	"reflect"

	validator "github.com/syssam/go-validator"
)

func CustomValidator(v reflect.Value, o reflect.Value, validTag *validator.ValidTag) bool {
	return false
}

func main() {
	validator.CustomTypeRuleMap.Set("customValidator", CustomValidator)
	validator.CustomTypeRuleMap.Set("customValidator2", func(v reflect.Value, o reflect.Value, validTag *validator.ValidTag) bool {
		return false
	})
}
