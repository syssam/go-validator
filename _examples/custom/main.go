package main

import (
	"fmt"
	"reflect"
	"regexp"

	validator "github.com/syssam/go-validator"
)

// User contains user information
type User struct {
	UserName string `valid:"customValidator"`
	Password string `valid:"customValidator2"`
}

func CustomValidator(v reflect.Value, o reflect.Value, validTag *validator.ValidTag) bool {
	return false
}

func main() {
	validator.MessageMap["customValidator"] = "customValidator is not valid."
	validator.MessageMap["customValidator2"] = "Beginning with a letter, allowing 5-16 bytes, allowing alphanumeric underlining."
	validator.CustomTypeRuleMap.Set("customValidator", CustomValidator)
	validator.CustomTypeRuleMap.Set("customValidator2", func(v reflect.Value, o reflect.Value, validTag *validator.ValidTag) bool {
		switch v.Kind() {
		case reflect.String:
			return regexp.MustCompile("^[a-zA-Z]\\w{5,17}$").MatchString(v.String())
		}
		return false
	})

	user := &User{
		UserName: "Tester",
		Password: "12345678",
	}

	err := validator.ValidateStruct(user)
	if err != nil {
		for _, err := range err.(validator.Errors) {
			fmt.Println(err)
		}
		return
	}
}
