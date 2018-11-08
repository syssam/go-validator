package main

import (
	"fmt"

	validator "github.com/syssam/go-validator"
	lang_en "github.com/syssam/go-validator/_examples/translations/lang/en"
	lang_zh_CN "github.com/syssam/go-validator/_examples/translations/lang/zh_CN"
	validator_en "github.com/syssam/go-validator/lang/en"
	validator_zh_CN "github.com/syssam/go-validator/lang/zh_CN"
)

// User contains user information
type User struct {
	FirstName string     `valid:"required"`
	LastName  string     `valid:"required"`
	Age       uint8      `valid:"between=0|30"`
	Email     string     `valid:"required,email"`
	Addresses []*Address `valid:"required"` // a person can have a home and cottage...
}

// Address houses a users address information
type Address struct {
	Street string `valid:"required"`
	City   string `valid:"required"`
	Planet string `valid:"required"`
	Phone  string `valid:"required"`
}

func main() {
	translator := validator.NewTranslator()
	translator.SetMessage("en", validator_en.MessageMap)
	translator.SetMessage("zh_CN", validator_zh_CN.MessageMap)
	translator.SetAttributes("en", lang_en.AttributeMap)
	translator.SetAttributes("zh_CN", lang_zh_CN.AttributeMap)
	validator.Default.Translator = translator

	address := &Address{
		Street: "Eavesdown Docks",
		Planet: "Persphone",
		Phone:  "none",
	}

	user := &User{
		FirstName: "Badger",
		LastName:  "Smith",
		Age:       135,
		Email:     "Badger.Smith@gmail.com",
		Addresses: []*Address{address},
	}

	err := validator.ValidateStruct(user)
	if err != nil {
		errs := validator.Default.Translator.Trans(err.(validator.Errors), "zh_CN")
		for _, err := range errs {
			fmt.Println(err)
		}
		return
	}
}
