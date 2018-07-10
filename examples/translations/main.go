package main

import (
	"fmt"

	validator "github.com/syssam/go-validator"
	lang_en "github.com/syssam/go-validator/examples/translations/lang/en"
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
	validator = validator.New()
	validator.Translator = validator.NewTranslator()
	validator.Translator.SetMessage("en", validator_en.MessageMap)
	validator.Translator.SetMessage("zh_CN", validator_zh_CN.MessageMap)
	validator.Translator.SetAttributes("en", lang_en.AttributeMap)
	validator.Translator.SetAttributes("zh_CN", validator_zh_CN.AttributeMap)
	validator.Translator.SetLocale("zh_CN")

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

	err := validatorInstance.Struct(user)
	if err != nil {
		for _, err := range err.(validator.Errors) {
			fmt.Println(err)
		}
		return
	}
}
