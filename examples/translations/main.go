package main

import (
	"fmt"

	validator "github.com/syssam/go-validator"
	lang_en "github.com/syssam/go-validator/examples/translations/lang/en"
	lang_zhCN "github.com/syssam/go-validator/examples/translations/lang/zh_CN"
	validator_en "github.com/syssam/go-validator/lang/en"
	validator_zhCN "github.com/syssam/go-validator/lang/zh_CN"
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

// use a single instance of Validate, it caches struct info
var validatorInstance *validator.Validator

func main() {

	validatorInstance = validator.New()
	validatorInstance.Translator = validator.NewTranslator()
	validatorInstance.Translator.SetMessage("en", validator_en.MessageMap)
	validatorInstance.Translator.SetMessage("zh_CN", validator_zhCN.MessageMap)
	validatorInstance.Translator.SetAttributes("en", lang_en.AttributeMap)
	validatorInstance.Translator.SetAttributes("zh_CN", lang_zhCN.AttributeMap)
	validatorInstance.Translator.SetLocale("zh_CN")

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
