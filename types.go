package validator

import "reflect"

// UnsupportedTypeError is a wrapper for reflect.Type
type UnsupportedTypeError struct {
	Type reflect.Type
}

// stringValues is a slice of reflect.Value holding *reflect.StringValue.
// It implements the methods to sort by string.
type stringValues []reflect.Value

// ValidateFunc is
type ValidateFunc func(v reflect.Value) bool

// ParamValidateFunc is
type ParamValidateFunc func(v reflect.Value, params ...string) bool

// StringValidateFunc is
type StringValidateFunc func(str string) bool

// StringParamValidateFunc is
type StringParamValidateFunc func(str string, params ...string) bool

// RuleMap is a map of functions, that can be used as tags for ValidateStruct function.
var RuleMap = map[string]ParamValidateFunc{}

// ParamRuleMap is a map of functions, that can be used as tags for ValidateStruct function.
var ParamRuleMap = map[string]ParamValidateFunc{
	"between": Between,
}

// StringRulesMap is a map of functions, that can be used as tags for ValidateStruct function when refelect type is string.
var StringRulesMap = map[string]StringValidateFunc{
	"email": IsEmail,
}
