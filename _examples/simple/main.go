package main

import (
	"fmt"

	validator "github.com/syssam/go-validator"
)

// User contains user information with comprehensive validation
type User struct {
	FirstName string     `valid:"required"`
	LastName  string     `valid:"required"`
	Age       uint8      `valid:"between=0|120"`
	Email     string     `valid:"required,email"`
	Addresses []*Address `valid:"required"`
}

// Address houses a user's address information
type Address struct {
	Street string `valid:"required"`
	City   string `valid:"required"`
	Planet string `valid:"required"`
	Phone  string `valid:"required"`
}

func main() {
	fmt.Println("🚀 Go Validator - Simple Example with Error Handling")
	fmt.Println("==================================================")

	address := &Address{
		Street: "Eavesdown Docks",
		Planet: "Persephone", // Fixed spelling
		Phone:  "none",
		// Missing City to demonstrate error handling
	}

	user := &User{
		FirstName: "Badger",
		LastName:  "Smith",
		Age:       135, // Invalid age to demonstrate error handling
		Email:     "Badger.Smith@gmail.com",
		Addresses: []*Address{address},
	}

	fmt.Println("📋 Validating user data...")
	err := validator.ValidateStruct(user)
	if err != nil {
		fmt.Println("❌ Validation failed with the following errors:")

		// Cast to Errors type for advanced error handling
		if validationErrors, ok := err.(validator.Errors); ok {
			// Group errors by field for better organization
			groupedErrors := validationErrors.GroupByField()

			for fieldName, fieldErrors := range groupedErrors {
				fmt.Printf("\n🔸 Field '%s':\n", fieldName)
				for _, fieldError := range fieldErrors {
					fmt.Println(fieldError)
				}
			}
		}
		return
	}
}
