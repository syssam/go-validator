package main

import (
	"fmt"

	validator "github.com/syssam/go-validator"
	validator_zh_CN "github.com/syssam/go-validator/lang/zh_CN"
)

// User contains user information
type User struct {
	FirstName string     `valid:"required,attribute=名字"`
	LastName  string     `valid:"required,attribute=姓氏"`
	Age       uint8      `valid:"between=0|30,attribute=年龄"`
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
	translator.SetMessage("zh_CN", validator_zh_CN.MessageMap)
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
