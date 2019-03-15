package main

import (
	"fmt"
	"reflect"

	validator "github.com/syssam/go-validator"
)

// User contains user information
type User struct {
	FirstName string `valid:"customValidator"`
	LastName  string `valid:"customValidator2"`
}

func CustomValidator(v reflect.Value, o reflect.Value, validTag *validator.ValidTag) bool {
	return false
}

func main() {
	validator.MessageMap["customValidator"] = "customValidator is not valid."
	validator.MessageMap["customValidator2"] = "customValidator2 is not valid."
	validator.CustomTypeRuleMap.Set("customValidator", CustomValidator)
	validator.CustomTypeRuleMap.Set("customValidator2", func(v reflect.Value, o reflect.Value, validTag *validator.ValidTag) bool {
		return false
	})

	user := &User{
		FirstName: "Badger",
		LastName:  "Smith",
	}

	err := validator.ValidateStruct(user)
	if err != nil {
		for _, err := range err.(validator.Errors) {
			fmt.Println(err)
		}
		return
	}
}
