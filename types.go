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
type ValidateFunc func(v reflect.Value) (bool, error)

// ParamValidateFunc is
type ParamValidateFunc func(v reflect.Value, params []string) (bool, error)

// StringValidateFunc is
type StringValidateFunc func(str string) bool

// StringParamValidateFunc is
type StringParamValidateFunc func(str string, params []string) bool

// CustomTypeValidateFunc is a wrapper for validator functions that returns bool.
// first parameter is field value
// second parameter is struct field
// third parameter is validTag message, pass the variable to the message
type CustomTypeValidateFunc func(v reflect.Value, o reflect.Value, validTag *ValidTag) bool

type customTypeRuleMap struct {
	validateFunc map[string]CustomTypeValidateFunc
	sync.RWMutex
}

// CustomTypeRuleMap is a map of functions that can be used as tags for ValidateStruct function.
var CustomTypeRuleMap = &customTypeRuleMap{validateFunc: make(map[string]CustomTypeValidateFunc)}

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
//
// Thread-safety: RuleMap, ParamRuleMap, StringRulesMap, StringParamRulesMap and
// Mimes are plain maps read on the validation hot path. Register custom rules
// during init/startup, before any concurrent ValidateStruct calls. For
// runtime-safe registration, use CustomTypeRuleMap instead.
var RuleMap = map[string]ValidateFunc{
	"distinct":   isDistinct,
	"accepted":   func(v reflect.Value) (bool, error) { return isAccepted(v), nil },
	"declined":   func(v reflect.Value) (bool, error) { return isDeclined(v), nil },
	"prohibited": func(v reflect.Value) (bool, error) { return isProhibited(v), nil },
	"missing":    func(v reflect.Value) (bool, error) { return isMissing(v), nil },
	"present":    func(v reflect.Value) (bool, error) { return isPresent(v), nil },
	"list":       isList,
	"date":       isDate,
	"json":       isJSON,
	"boolean":    isBoolean,
	"string":     isString,
	"array":      isArray,
	"filled":     isFilled,
}

// ParamRuleMap is a map of functions, that can be used as tags for ValidateStruct function.
var ParamRuleMap = map[string]ParamValidateFunc{
	"between":           isBetween,
	"digitsBetween":     isDigitsBetween,
	"min":               isMin,
	"max":               isMax,
	"size":              isSize,
	"gt":                isGtParam,
	"gte":               isGteParam,
	"lt":                isLtParam,
	"lte":               isLteParam,
	"multipleOf":        isMultipleOf,
	"maxDigits":         isMaxDigits,
	"minDigits":         isMinDigits,
	"decimal":           isDecimalPrecision,
	"contains":          isContains,
	"doesntContain":     isDoesntContain,
	"requiredArrayKeys": isRequiredArrayKeys,
	"dateFormat":        isDateFormat,
	"after":             isAfter,
	"afterOrEqual":      isAfterOrEqual,
	"before":            isBefore,
	"beforeOrEqual":     isBeforeOrEqual,
	"in":                isIn,
	"notIn":             isNotIn,
	"regex":             isRegex,
	"notRegex":          isNotRegex,
}

// StringRulesMap is a map of functions, that can be used as tags for ValidateStruct function when reflect type is string.
var StringRulesMap = map[string]StringValidateFunc{
	"numeric":          IsNumeric,
	"int":              IsInt,
	"integer":          IsInt,
	"float":            IsFloat,
	"email":            IsEmail,
	"alpha":            IsAlpha,
	"alphaNum":         IsAlphaNum,
	"alphaDash":        IsAlphaDash,
	"alphaUnicode":     IsAlphaUnicode,
	"alphaNumUnicode":  IsAlphaNumUnicode,
	"alphaDashUnicode": IsAlphaDashUnicode,
	"ip":               IsIP,
	"ipv4":             IsIPv4,
	"ipv6":             IsIPv6,
	"uuid3":            IsUUID3,
	"uuid4":            IsUUID4,
	"uuid5":            IsUUID5,
	"uuid":             IsUUID,
	"ulid":             IsULID,
	"url":              IsURL,
	"hexColor":         IsHexColor,
	"timezone":         IsTimezone,
	"ascii":            IsASCII,
	"lowercase":        IsLowercase,
	"uppercase":        IsUppercase,
	"macAddress":       IsMACAddress,
	// Country codes (ISO 3166-1)
	"country":         IsCountryCode,
	"country.alpha2":  IsCountryAlpha2,
	"country.alpha3":  IsCountryAlpha3,
	"country.numeric": IsCountryNumeric,
	// Currency codes (ISO 4217 + Crypto)
	"currency":        IsCurrencyAll,
	"currency.all":    IsCurrencyAll,
	"currency.fiat":   IsCurrencyFiat,
	"currency.crypto": IsCurrencyCrypto,
	// Language codes (ISO 639)
	"language":        IsLanguageCode,
	"language.alpha2": IsLanguageAlpha2,
	"language.alpha3": IsLanguageAlpha3,
	// Phone validation
	"phone.e164": IsPhoneE164,
	"phone":      IsPhoneValid,
}

// StringParamRulesMap is a map of functions with params, that can be used as tags for ValidateStruct function when reflect type is string.
var StringParamRulesMap = map[string]StringParamValidateFunc{
	"startsWith":      IsStartsWith,
	"endsWith":        IsEndsWith,
	"doesntStartWith": IsDoesntStartWith,
	"doesntEndWith":   IsDoesntEndWith,
}
