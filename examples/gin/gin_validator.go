package main

import (
	"sync"

	"github.com/gin-gonic/gin/binding"
	validator "github.com/syssam/go-validator"
)

// DefaultValidator struct
type DefaultValidator struct {
	once     sync.Once
	validate *validator.Validator
}

var _ binding.StructValidator = &DefaultValidator{}

// ValidateStruct return error
func (v *DefaultValidator) ValidateStruct(obj interface{}) error {
	v.lazyinit()
	if err := v.validate.Struct(obj); err != nil {
		return error(err)
	}
	return nil
}

// Engine return v.validate
func (v *DefaultValidator) Engine() interface{} {
	v.lazyinit()
	return v.validate
}

func (v *DefaultValidator) lazyinit() {
	v.once.Do(func() {
		v.validate = validator.New()
	})
}
