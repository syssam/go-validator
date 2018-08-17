package main

import (
	"sync"

	validator "github.com/syssam/go-validator"
)

// DefaultValidator struct
type DefaultValidator struct {
	once     sync.Once
	validate *validator.Validator
}

// Validate return error
func (v *DefaultValidator) Validate(obj interface{}) error {
	v.lazyinit()
	if err := validator.ValidateStruct(obj); err != nil {
		return error(err)
	}
	return nil
}

func (v *DefaultValidator) lazyinit() {
	v.once.Do(func() {
		v.validate = validator.New()
	})
}
