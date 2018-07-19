package validator

import (
	"reflect"
	"sync"
)

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

type CustomTypeValidateFunc func(v reflect.Value, f *field, o reflect.Value) bool

type customTypeRuleMap struct {
	validateFunc map[string]CustomTypeValidateFunc
	sync.RWMutex
}

var CustomTypeTagMap = &customTypeRuleMap{validateFunc: make(map[string]CustomTypeValidateFunc)}

func (tm *customTypeRuleMap) Get(name string) (CustomTypeValidateFunc, bool) {
	tm.RLock()
	defer tm.RUnlock()
	v, ok := tm.validateFunc[name]
	return v, ok
}

func (tm *customTypeRuleMap) Set(name string, ctv CustomTypeValidateFunc) {
	tm.Lock()
	defer tm.Unlock()
	tm.validateFunc[name] = ctv
}

// RuleMap is a map of functions, that can be used as tags for ValidateStruct function.
var RuleMap = map[string]ParamValidateFunc{}

// ParamRuleMap is a map of functions, that can be used as tags for ValidateStruct function.
var ParamRuleMap = map[string]ParamValidateFunc{
	"between": Between,
}

// ParamsRuleMap is a map of functions, that can be used as tags for ValidateStruct function.
var ParamsRuleMap = map[string]ParamValidateFunc{
	"between": Between,
	"gt":      Gt,
	"gte":     Gte,
	"lt":      Lt,
	"lte":     Lte,
	"min":     Min,
	"max":     Max,
}

// StringRulesMap is a map of functions, that can be used as tags for ValidateStruct function when refelect type is string.
var StringRulesMap = map[string]StringValidateFunc{
	"email": IsEmail,
}
