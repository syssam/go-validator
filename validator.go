package validator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

// regexCache stores compiled regular expressions to avoid recompilation
var regexCache sync.Map // map[string]*regexp.Regexp

// getCompiledRegex returns a cached compiled regex or compiles and caches it.
// Uses LoadOrStore pattern to handle concurrent access correctly.
func getCompiledRegex(pattern string) (*regexp.Regexp, error) {
	// Fast path: check if already cached
	if cached, ok := regexCache.Load(pattern); ok {
		return cached.(*regexp.Regexp), nil
	}

	// Slow path: compile and store
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	// Use LoadOrStore to handle race condition - if another goroutine
	// stored the same pattern concurrently, use that one instead
	actual, _ := regexCache.LoadOrStore(pattern, re)
	return actual.(*regexp.Regexp), nil
}

// defaultBufferCap is the default capacity for pooled byte buffers
const defaultBufferCap = 128

// maxBufferCap prevents unbounded buffer growth in the pool
const maxBufferCap = 1024

// byteBufferPool provides reusable byte buffers to reduce allocations
var byteBufferPool = sync.Pool{
	New: func() any {
		buf := make([]byte, 0, defaultBufferCap)
		return &buf
	},
}

// getBuffer gets a buffer from the pool
func getBuffer() *[]byte {
	return byteBufferPool.Get().(*[]byte)
}

// putBuffer returns a buffer to the pool.
// Buffers that grew too large are discarded to prevent memory bloat.
func putBuffer(buf *[]byte) {
	// Discard buffers that grew too large to prevent memory bloat
	if cap(*buf) > maxBufferCap {
		return
	}
	*buf = (*buf)[:0]
	byteBufferPool.Put(buf)
}

// buildFieldName efficiently builds a field name string
func buildFieldName(namespace, fieldName []byte) string {
	if len(namespace) == 0 {
		return string(fieldName)
	}
	buf := getBuffer()
	*buf = append(*buf, namespace...)
	*buf = append(*buf, fieldName...)
	result := string(*buf)
	putBuffer(buf)
	return result
}

const tagName string = "valid"

// Validator construct
type Validator struct {
	Attributes    map[string]string
	CustomMessage map[string]string
	Translator    *Translator
}

// Default returns a instance of Validator
var Default = New()

// New returns a new instance of Validator
func New() *Validator {
	return &Validator{}
}

// validateBetween check The field under validation must have a size between the given min and max. Strings, numerics, arrays, and files are evaluated in the same fashion as the size rule.
//
//nolint:gocyclo,gocritic // Complex validation logic with parameter names
func validateBetween(v reflect.Value, params []string) (bool, error) {
	if len(params) != 2 {
		return false, fmt.Errorf("validator: Between params length must be 2")
	}

	// Check for decimal.Decimal type first
	if d, ok := asDecimal(v); ok {
		minVal, maxVal, err := parseDecimalParams(params)
		if err != nil {
			return false, fmt.Errorf("validator: Between decimal: %w", err)
		}
		return IsDecimalBetween(d, minVal, maxVal), nil
	}

	var valid bool
	var err error

	switch v.Kind() {
	case reflect.String:
		minVal, err := ToInt(params[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Between rule on string field, min value: %w", err)
		}
		maxVal, err := ToInt(params[1])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Between rule on string field, max value: %w", err)
		}
		valid = IsStringBetween(v.String(), minVal, maxVal)
	case reflect.Slice, reflect.Map, reflect.Array:
		minVal, err := ToInt(params[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Between rule on collection field, min value: %w", err)
		}
		maxVal, err := ToInt(params[1])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Between rule on collection field, max value: %w", err)
		}
		valid = IsInt64Between(int64(v.Len()), minVal, maxVal)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		minVal, err := ToInt(params[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Between rule on numeric field, min value: %w", err)
		}
		maxVal, err := ToInt(params[1])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Between rule on numeric field, max value: %w", err)
		}
		valid = IsInt64Between(v.Int(), minVal, maxVal)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		minVal, err := ToUint(params[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Between rule on numeric field, min value: %w", err)
		}
		maxVal, err := ToUint(params[1])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Between rule on numeric field, max value: %w", err)
		}
		valid = IsUint64Between(v.Uint(), minVal, maxVal)
	case reflect.Float32, reflect.Float64:
		minVal, err := ToFloat(params[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Between rule on numeric field, min value: %w", err)
		}
		maxVal, err := ToFloat(params[1])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Between rule on numeric field, max value: %w", err)
		}
		valid = IsFloat64Between(v.Float(), minVal, maxVal)
	default:
		return false, fmt.Errorf("validator: Between unsupported type %T", v.Interface())
	}

	return valid, err
}

// createFieldError creates a FieldError struct with common fields populated
func (v *Validator) createFieldError(name, structName, tagName, messageName string, messageParameters MessageParameters, attribute, defaultAttribute, value string, funcError error) *FieldError {
	return &FieldError{
		Name:              name,
		StructName:        structName,
		Tag:               tagName,
		MessageName:       messageName,
		MessageParameters: messageParameters,
		Attribute:         attribute,
		DefaultAttribute:  defaultAttribute,
		Value:             value,
		FuncError:         funcError,
	}
}

// validateWithRuleMap validates a value using RuleMap and returns formatted error if validation fails
func (v *Validator) validateWithRuleMap(tag *ValidTag, value reflect.Value, f *field, name, structName string, o reflect.Value) error {
	if validfunc, ok := RuleMap[tag.name]; ok {
		isValid, funcError := validfunc(value)
		if !isValid {
			return v.formatsMessages(v.createFieldError(
				name, structName, tag.name, tag.messageName,
				parseValidatorMessageParameters(tag, o),
				f.attribute, f.defaultAttribute,
				ToString(value.Interface()), funcError,
			))
		}
	}
	return nil
}

// validateWithParamRuleMap validates a value using ParamRuleMap and returns formatted error if validation fails
func (v *Validator) validateWithParamRuleMap(tag *ValidTag, value reflect.Value, f *field, name, structName string, o reflect.Value) error {
	if validfunc, ok := ParamRuleMap[tag.name]; ok {
		isValid, funcError := validfunc(value, tag.params)
		if !isValid {
			return v.formatsMessages(v.createFieldError(
				name, structName, tag.name, tag.messageName,
				parseValidatorMessageParameters(tag, o),
				f.attribute, f.defaultAttribute,
				ToString(value.Interface()), funcError,
			))
		}
	}
	return nil
}

// validateWithStringRulesMap validates a string value using StringRulesMap and returns formatted error if validation fails
func (v *Validator) validateWithStringRulesMap(tag *ValidTag, value reflect.Value, f *field, name, structName string, o reflect.Value) error {
	if validfunc, ok := StringRulesMap[tag.name]; ok {
		isValid := validfunc(value.String())
		if !isValid {
			return v.formatsMessages(v.createFieldError(
				name, structName, tag.name, tag.messageName,
				parseValidatorMessageParameters(tag, o),
				f.attribute, f.defaultAttribute,
				ToString(value.Interface()), nil,
			))
		}
	}
	return nil
}

// validateWithStringParamRulesMap validates a string value using StringParamRulesMap and returns formatted error if validation fails
func (v *Validator) validateWithStringParamRulesMap(tag *ValidTag, value reflect.Value, f *field, name, structName string, o reflect.Value) error {
	if validfunc, ok := StringParamRulesMap[tag.name]; ok {
		isValid := validfunc(value.String(), tag.params)
		if !isValid {
			return v.formatsMessages(v.createFieldError(
				name, structName, tag.name, tag.messageName,
				parseValidatorMessageParameters(tag, o),
				f.attribute, f.defaultAttribute,
				ToString(value.Interface()), nil,
			))
		}
	}
	return nil
}

// validateCommonRules applies common validation rules (RuleMap, ParamRuleMap, dependent rules)
func (v *Validator) validateCommonRules(tags otherValidTags, value reflect.Value, f *field, name, structName string, o reflect.Value) error {
	for _, tag := range tags {
		handled, err := v.checkDependentRulesWithStatus(tag, f, value, o, name, structName)
		if err != nil {
			return err
		}

		// Skip ParamRuleMap for comparison rules if they were handled by field comparison
		skipParamRule := handled && (tag.name == "gt" || tag.name == "gte" || tag.name == "lt" || tag.name == "lte")

		if err := v.validateWithRuleMap(tag, value, f, name, structName, o); err != nil {
			return err
		}

		if !skipParamRule {
			if err := v.validateWithParamRuleMap(tag, value, f, name, structName, o); err != nil {
				return err
			}
		}

		if value.Kind() == reflect.String {
			if err := v.validateWithStringRulesMap(tag, value, f, name, structName, o); err != nil {
				return err
			}
			if err := v.validateWithStringParamRulesMap(tag, value, f, name, structName, o); err != nil {
				return err
			}
		}
	}
	return nil
}

// extractValuesFromCollection extracts string values from map or slice/array
func extractValuesFromCollection(field reflect.Value) ([]string, error) {
	var values []string

	switch field.Kind() {
	case reflect.Map:
		var sv stringValues
		sv = field.MapKeys()
		sort.Sort(sv)
		for _, k := range sv {
			mapValue := field.MapIndex(k)
			if mapValue.Kind() == reflect.Interface || mapValue.Kind() == reflect.Ptr {
				mapValue = mapValue.Elem()
			}
			if mapValue.Kind() != reflect.Struct {
				values = append(values, ToString(mapValue.Interface()))
			} else {
				return nil, fmt.Errorf("validator: RequiredIf unsupported type %T", mapValue.Interface())
			}
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < field.Len(); i++ {
			sliceValue := field.Index(i)
			if sliceValue.Kind() == reflect.Interface || sliceValue.Kind() == reflect.Ptr {
				sliceValue = sliceValue.Elem()
			}
			if sliceValue.Kind() != reflect.Struct {
				values = append(values, ToString(sliceValue.Interface()))
			} else {
				return nil, fmt.Errorf("validator: RequiredIf unsupported type %T", sliceValue.Interface())
			}
		}
	}

	return values, nil
}

// checkRequiredIfCondition checks if the required condition is met and updates tag parameters
func checkRequiredIfCondition(v reflect.Value, values, params []string, tag *ValidTag) (bool, error) {
	for _, value := range values {
		if InString(value, params) && Empty(v) {
			if tag != nil {
				tag.messageParameters = append(
					tag.messageParameters,
					messageParameter{
						Key:   "Value",
						Value: value,
					},
				)
			}
			return false, nil
		}
	}
	return true, nil
}

// validateCustomTypeRules validates using CustomTypeRuleMap
func (v *Validator) validateCustomTypeRules(tags otherValidTags, value reflect.Value, f *field, name, structName string, o reflect.Value) error {
	for _, tag := range tags {
		if validatefunc, ok := CustomTypeRuleMap.Get(tag.name); ok {
			if result := validatefunc(value, o, tag); !result {
				return v.formatsMessages(v.createFieldError(
					name, structName, tag.name, tag.messageName,
					parseValidatorMessageParameters(tag, o),
					f.attribute, f.defaultAttribute,
					ToString(value.Interface()), nil,
				))
			}
		}
	}
	return nil
}

// validateMapFields validates map structure and each element
func (v *Validator) validateMapFields(value reflect.Value, f *field, jsonNamespace, structNamespace []byte) error {
	if value.Type().Key().Kind() != reflect.String {
		return &UnsupportedTypeError{value.Type()}
	}

	sv := stringValues(value.MapKeys())
	sort.Sort(sv)
	for _, k := range sv {
		var err error
		item := value.MapIndex(k)
		if value.Kind() == reflect.Interface {
			item = item.Elem()
		}

		if item.Kind() == reflect.Struct || item.Kind() == reflect.Ptr {
			newJSONNamespace := append(append(jsonNamespace, f.nameBytes...), '.')
			newJSONNamespace = append(append(newJSONNamespace, []byte(k.String())...), '.')
			newstructNamespace := append(append(structNamespace, f.structNameBytes...), '.')
			newstructNamespace = append(append(newstructNamespace, []byte(k.String())...), '.')
			err = v.ValidateStruct(item.Interface(), newJSONNamespace, newstructNamespace)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// validateSliceFields validates slice/array structure and each element
func (v *Validator) validateSliceFields(value reflect.Value, f *field, jsonNamespace, structNamespace []byte) error {
	for i := 0; i < value.Len(); i++ {
		var err error
		item := value.Index(i)
		if item.Kind() == reflect.Interface {
			item = item.Elem()
		}

		if item.Kind() == reflect.Struct || item.Kind() == reflect.Ptr {
			newJSONNamespace := append(append(jsonNamespace, f.nameBytes...), '.')
			newJSONNamespace = append(append(newJSONNamespace, []byte(strconv.Itoa(i))...), '.')
			newStructNamespace := append(append(structNamespace, f.structNameBytes...), '.')
			newStructNamespace = append(append(newStructNamespace, []byte(strconv.Itoa(i))...), '.')
			err = v.ValidateStruct(value.Index(i).Interface(), newJSONNamespace, newStructNamespace)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// ValidateBetween check The field under validation must have a size between the given min and max. Strings, numerics, arrays, and files are evaluated in the same fashion as the size rule.
func ValidateBetween(i interface{}, params []string) (bool, error) {
	v := reflect.ValueOf(i)
	return validateBetween(v, params)
}

// validateDigitsBetween check The field under validation must have a length between the given min and max.
func validateDigitsBetween(v reflect.Value, params []string) (bool, error) {
	if len(params) != 2 {
		return false, fmt.Errorf("validator: DigitsBetween params length must be 2")
	}

	switch v.Kind() {
	case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		min, err := ToInt(params[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for DigitsBetween rule on string field, min value: %w", err)
		}
		max, err := ToInt(params[1])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for DigitsBetween rule on string field, max value: %w", err)
		}
		var value string
		switch v.Kind() {
		case reflect.String:
			value = v.String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			value = ToString(v.Int())
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			value = ToString(v.Uint())
		}

		if value == "" || !IsNumeric(value) {
			return false, fmt.Errorf("validator: DigitsBetween value is not numeric")
		}

		return IsStringBetween(value, min, max), nil
	}

	return false, fmt.Errorf("validator: DigitsBetween unsupported type %T", v.Interface())
}

// ValidateDigitsBetween check The field under validation must have a length between the given min and max.
func ValidateDigitsBetween(i interface{}, params []string) (bool, error) {
	v := reflect.ValueOf(i)
	return validateDigitsBetween(v, params)
}

// validateSize The field under validation must have a size matching the given value.
// For string data, value corresponds to the number of characters.
// For numeric data, value corresponds to a given integer value.
// For an array | map | slice, size corresponds to the count of the array | map | slice.
func validateSize(v reflect.Value, param []string) (bool, error) {
	valid := false
	var err error
	switch v.Kind() {
	case reflect.String:
		p, err := ToInt(param[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Size rule on string field, value: %w", err)
		}
		valid, err = compareString(v.String(), p, "==")
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Size rule on string field, value: %w", err)
		}
	case reflect.Slice, reflect.Map, reflect.Array:
		p, err := ToInt(param[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Size rule on collection field, value: %w", err)
		}
		valid, err = compareInt64(int64(v.Len()), p, "==")
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Size rule on collection field, value: %w", err)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p, err := ToInt(param[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Size rule on numeric field, value: %w", err)
		}
		valid, err = compareInt64(v.Int(), p, "==")
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Size rule on numeric field, value: %w", err)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p, err := ToUint(param[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Size rule on numeric field, value: %w", err)
		}
		valid, err = compareUint64(v.Uint(), p, "==")
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Size rule on numeric field, value: %w", err)
		}
	case reflect.Float32, reflect.Float64:
		p, err := ToFloat(param[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Size rule on numeric field, value: %w", err)
		}
		valid, err = compareFloat64(v.Float(), p, "==")
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Size rule on numeric field, value: %w", err)
		}
	default:
		return false, fmt.Errorf("validator: Size unsupported type %T", v.Interface())
	}

	return valid, err
}

// ValidateSize The field under validation must have a size matching the given value.
// For string data, value corresponds to the number of characters.
// For numeric data, value corresponds to a given integer value.
// For an array | map | slice, size corresponds to the count of the array | map | slice.
func ValidateSize(i interface{}, params []string) (bool, error) {
	v := reflect.ValueOf(i)
	return validateSize(v, params)
}

// validateMax is the validation function for validating if the current field's value is less than or equal to the param's value.
//
//nolint:gocyclo,gocritic // Complex validation logic
func validateMax(v reflect.Value, param []string) (bool, error) {
	// Check for decimal.Decimal type first
	if d, ok := asDecimal(v); ok {
		p, _, err := parseDecimalParams(param)
		if err != nil {
			return false, fmt.Errorf("validator: Max decimal: %w", err)
		}
		return IsDecimalLte(d, p), nil
	}

	var valid bool
	var err error

	switch v.Kind() {
	case reflect.String:
		p, err := ToInt(param[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Max rule on string field, value: %w", err)
		}
		valid, err = compareString(v.String(), p, "<=")
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Max rule on string field, value: %w", err)
		}
	case reflect.Slice, reflect.Map, reflect.Array:
		p, err := ToInt(param[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Max rule on collection field, value: %w", err)
		}
		valid, err = compareInt64(int64(v.Len()), p, "<=")
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Max rule on collection field, value: %w", err)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p, err := ToInt(param[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Max rule on numeric field, value: %w", err)
		}
		valid, err = compareInt64(v.Int(), p, "<=")
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Max rule on numeric field, value: %w", err)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p, err := ToUint(param[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Max rule on numeric field, value: %w", err)
		}
		valid, err = compareUint64(v.Uint(), p, "<=")
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Max rule on numeric field, value: %w", err)
		}
	case reflect.Float32, reflect.Float64:
		p, err := ToFloat(param[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Max rule on numeric field, value: %w", err)
		}
		valid, err = compareFloat64(v.Float(), p, "<=")
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Max rule on numeric field, value: %w", err)
		}
	default:
		return false, fmt.Errorf("validator: Max unsupported type %T", v.Interface())
	}

	return valid, err
}

// ValidateMax is the validation function for validating if the current field's value is less than or equal to the param's value.
func ValidateMax(i interface{}, params []string) (bool, error) {
	v := reflect.ValueOf(i)
	return validateMax(v, params)
}

// validateMin is the validation function for validating if the current field's value is greater than or equal to the param's value.
func validateMin(v reflect.Value, param []string) (bool, error) {
	// Check for decimal.Decimal type first
	if d, ok := asDecimal(v); ok {
		p, _, err := parseDecimalParams(param)
		if err != nil {
			return false, fmt.Errorf("validator: Min decimal: %w", err)
		}
		return IsDecimalGte(d, p), nil
	}

	var valid bool
	var err error

	switch v.Kind() {
	case reflect.String:
		p, err := ToInt(param[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Min rule on string field, value: %w", err)
		}
		valid, err = compareString(v.String(), p, ">=")
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Min rule on string field, value: %w", err)
		}
	case reflect.Slice, reflect.Map, reflect.Array:
		p, err := ToInt(param[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Min rule on collection field, value: %w", err)
		}
		valid, err = compareInt64(int64(v.Len()), p, ">=")
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Min rule on collection field, value: %w", err)
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p, err := ToInt(param[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Min rule on numeric field, value: %w", err)
		}
		valid, err = compareInt64(v.Int(), p, ">=")
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Min rule on numeric field, value: %w", err)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p, err := ToUint(param[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Min rule on numeric field, value: %w", err)
		}
		valid, err = compareUint64(v.Uint(), p, ">=")
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Min rule on numeric field, value: %w", err)
		}
	case reflect.Float32, reflect.Float64:
		p, err := ToFloat(param[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Min rule on numeric field, value: %w", err)
		}
		valid, err = compareFloat64(v.Float(), p, ">=")
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Min rule on numeric field, value: %w", err)
		}
	default:
		return false, fmt.Errorf("validator: Min unsupported type %T", v.Interface())
	}

	return valid, err
}

// ValidateMin is the validation function for validating if the current field's value is greater than or equal to the param's value.
func ValidateMin(i interface{}, params []string) (bool, error) {
	v := reflect.ValueOf(i)
	return validateMin(v, params)
}

// validateGtParam is the validation function for validating if the current field's value is greater than the param's value.
func validateGtParam(v reflect.Value, params []string) (bool, error) {
	if len(params) == 0 {
		return false, fmt.Errorf("validator: Gt rule requires at least one parameter")
	}

	// Check for decimal.Decimal type first
	if d, ok := asDecimal(v); ok {
		p, _, err := parseDecimalParams(params)
		if err != nil {
			return false, fmt.Errorf("validator: Gt decimal: %w", err)
		}
		return IsDecimalGt(d, p), nil
	}

	var valid bool
	var err error

	switch v.Kind() {
	case reflect.String:
		p, err := ToInt(params[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Gt rule on string field: %w", err)
		}
		valid, err = compareString(v.String(), p, ">")
		if err != nil {
			return false, fmt.Errorf("validator: comparison error for Gt rule on string field: %w", err)
		}
	case reflect.Slice, reflect.Map, reflect.Array:
		p, err := ToInt(params[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Gt rule on collection field: %w", err)
		}
		valid, err = compareInt64(int64(v.Len()), p, ">")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p, err := ToInt(params[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Gt rule on int field: %w", err)
		}
		valid, err = compareInt64(v.Int(), p, ">")
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p, err := ToUint(params[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Gt rule on uint field: %w", err)
		}
		valid, err = compareUint64(v.Uint(), p, ">")
	case reflect.Float32, reflect.Float64:
		p, err := ToFloat(params[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Gt rule on float field: %w", err)
		}
		valid, err = compareFloat64(v.Float(), p, ">")
	default:
		return false, fmt.Errorf("validator: Gt rule is not supported for type %s", v.Kind())
	}

	return valid, err
}

// ValidateGtParam is the validation function for validating if the current field's value is greater than the param's value.
func ValidateGtParam(i interface{}, params []string) (bool, error) {
	v := reflect.ValueOf(i)
	return validateGtParam(v, params)
}

// validateGteParam is the validation function for validating if the current field's value is greater than or equal to the param's value.
func validateGteParam(v reflect.Value, params []string) (bool, error) {
	if len(params) == 0 {
		return false, fmt.Errorf("validator: Gte rule requires at least one parameter")
	}

	// Check for decimal.Decimal type first
	if d, ok := asDecimal(v); ok {
		p, _, err := parseDecimalParams(params)
		if err != nil {
			return false, fmt.Errorf("validator: Gte decimal: %w", err)
		}
		return IsDecimalGte(d, p), nil
	}

	var valid bool
	var err error

	switch v.Kind() {
	case reflect.String:
		p, err := ToInt(params[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Gte rule on string field: %w", err)
		}
		valid, err = compareString(v.String(), p, ">=")
	case reflect.Slice, reflect.Map, reflect.Array:
		p, err := ToInt(params[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Gte rule on collection field: %w", err)
		}
		valid, err = compareInt64(int64(v.Len()), p, ">=")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p, err := ToInt(params[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Gte rule on int field: %w", err)
		}
		valid, err = compareInt64(v.Int(), p, ">=")
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p, err := ToUint(params[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Gte rule on uint field: %w", err)
		}
		valid, err = compareUint64(v.Uint(), p, ">=")
	case reflect.Float32, reflect.Float64:
		p, err := ToFloat(params[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Gte rule on float field: %w", err)
		}
		valid, err = compareFloat64(v.Float(), p, ">=")
	default:
		return false, fmt.Errorf("validator: Gte rule is not supported for type %s", v.Kind())
	}

	return valid, err
}

// ValidateGteParam is the validation function for validating if the current field's value is greater than or equal to the param's value.
func ValidateGteParam(i interface{}, params []string) (bool, error) {
	v := reflect.ValueOf(i)
	return validateGteParam(v, params)
}

// validateLtParam is the validation function for validating if the current field's value is less than the param's value.
func validateLtParam(v reflect.Value, params []string) (bool, error) {
	if len(params) == 0 {
		return false, fmt.Errorf("validator: Lt rule requires at least one parameter")
	}

	// Check for decimal.Decimal type first
	if d, ok := asDecimal(v); ok {
		p, _, err := parseDecimalParams(params)
		if err != nil {
			return false, fmt.Errorf("validator: Lt decimal: %w", err)
		}
		return IsDecimalLt(d, p), nil
	}

	var valid bool
	var err error

	switch v.Kind() {
	case reflect.String:
		p, err := ToInt(params[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Lt rule on string field: %w", err)
		}
		valid, err = compareString(v.String(), p, "<")
	case reflect.Slice, reflect.Map, reflect.Array:
		p, err := ToInt(params[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Lt rule on collection field: %w", err)
		}
		valid, err = compareInt64(int64(v.Len()), p, "<")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p, err := ToInt(params[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Lt rule on int field: %w", err)
		}
		valid, err = compareInt64(v.Int(), p, "<")
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p, err := ToUint(params[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Lt rule on uint field: %w", err)
		}
		valid, err = compareUint64(v.Uint(), p, "<")
	case reflect.Float32, reflect.Float64:
		p, err := ToFloat(params[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Lt rule on float field: %w", err)
		}
		valid, err = compareFloat64(v.Float(), p, "<")
	default:
		return false, fmt.Errorf("validator: Lt rule is not supported for type %s", v.Kind())
	}

	return valid, err
}

// ValidateLtParam is the validation function for validating if the current field's value is less than the param's value.
func ValidateLtParam(i interface{}, params []string) (bool, error) {
	v := reflect.ValueOf(i)
	return validateLtParam(v, params)
}

// validateLteParam is the validation function for validating if the current field's value is less than or equal to the param's value.
func validateLteParam(v reflect.Value, params []string) (bool, error) {
	if len(params) == 0 {
		return false, fmt.Errorf("validator: Lte rule requires at least one parameter")
	}

	// Check for decimal.Decimal type first
	if d, ok := asDecimal(v); ok {
		p, _, err := parseDecimalParams(params)
		if err != nil {
			return false, fmt.Errorf("validator: Lte decimal: %w", err)
		}
		return IsDecimalLte(d, p), nil
	}

	var valid bool
	var err error

	switch v.Kind() {
	case reflect.String:
		p, err := ToInt(params[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Lte rule on string field: %w", err)
		}
		valid, err = compareString(v.String(), p, "<=")
	case reflect.Slice, reflect.Map, reflect.Array:
		p, err := ToInt(params[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Lte rule on collection field: %w", err)
		}
		valid, err = compareInt64(int64(v.Len()), p, "<=")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p, err := ToInt(params[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Lte rule on int field: %w", err)
		}
		valid, err = compareInt64(v.Int(), p, "<=")
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p, err := ToUint(params[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Lte rule on uint field: %w", err)
		}
		valid, err = compareUint64(v.Uint(), p, "<=")
	case reflect.Float32, reflect.Float64:
		p, err := ToFloat(params[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Lte rule on float field: %w", err)
		}
		valid, err = compareFloat64(v.Float(), p, "<=")
	default:
		return false, fmt.Errorf("validator: Lte rule is not supported for type %s", v.Kind())
	}

	return valid, err
}

// ValidateLteParam is the validation function for validating if the current field's value is less than or equal to the param's value.
func ValidateLteParam(i interface{}, params []string) (bool, error) {
	v := reflect.ValueOf(i)
	return validateLteParam(v, params)
}

// validateSame is the validation function for validating if the current field's value equal the param's value.
func validateSame(v, anotherField reflect.Value) (bool, error) {
	if !v.IsValid() || !anotherField.IsValid() {
		return false, fmt.Errorf("validator: Same invalid reflection values")
	}
	if v.Kind() != anotherField.Kind() {
		return false, fmt.Errorf("validator: Same The two fields must be of the same type %T, %T", v.Interface(), anotherField.Interface())
	}

	var valid bool
	var err error

	// Check for decimal.Decimal type first
	if d1, ok := asDecimal(v); ok {
		if d2, ok := asDecimal(anotherField); ok {
			return IsDecimalEqual(d1, d2), nil
		}
	}

	switch v.Kind() {
	case reflect.String:
		valid, err = v.String() == anotherField.String(), nil
	case reflect.Slice, reflect.Map, reflect.Array:
		valid, err = compareInt64(int64(v.Len()), int64(anotherField.Len()), "==")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		valid, err = compareInt64(v.Int(), anotherField.Int(), "==")
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		valid, err = compareUint64(v.Uint(), anotherField.Uint(), "==")
	case reflect.Float32, reflect.Float64:
		valid, err = compareFloat64(v.Float(), anotherField.Float(), "==")
	default:
		return false, fmt.Errorf("validator: Same unsupported type %T", v.Interface())
	}

	return valid, err
}

// ValidateSame is the validation function for validating if the current field's value is greater than or equal to the param's value.
func ValidateSame(i, a interface{}) (bool, error) {
	v := reflect.ValueOf(i)
	anotherField := reflect.ValueOf(a)
	return validateSame(v, anotherField)
}

// validateLt is the validation function for validating if the current field's value is less than the param's value.
func validateLt(v, anotherField reflect.Value) (bool, error) {
	if !v.IsValid() || !anotherField.IsValid() {
		return false, fmt.Errorf("validator: Lt invalid reflection values")
	}
	if v.Kind() != anotherField.Kind() {
		return false, fmt.Errorf("validator: Lt The two fields must be of the same type %T, %T", v.Interface(), anotherField.Interface())
	}

	// Check for decimal.Decimal type first
	if d1, ok := asDecimal(v); ok {
		if d2, ok := asDecimal(anotherField); ok {
			return IsDecimalLt(d1, d2), nil
		}
	}

	var valid bool
	var err error

	switch v.Kind() {
	case reflect.String:
		valid, err = compareString(v.String(), int64(utf8.RuneCountInString(anotherField.String())), "<")
	case reflect.Slice, reflect.Map, reflect.Array:
		valid, err = compareInt64(int64(v.Len()), int64(anotherField.Len()), "<")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		valid, err = compareInt64(v.Int(), anotherField.Int(), "<")
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		valid, err = compareUint64(v.Uint(), anotherField.Uint(), "<")
	case reflect.Float32, reflect.Float64:
		valid, err = compareFloat64(v.Float(), anotherField.Float(), "<")
	default:
		return false, fmt.Errorf("validator: Lt unsupported type %T", v.Interface())
	}

	return valid, err
}

// ValidateLt is the validation function for validating if the current field's value is less than the param's value.
func ValidateLt(i, a interface{}) (bool, error) {
	v := reflect.ValueOf(i)
	anotherField := reflect.ValueOf(a)
	return validateLt(v, anotherField)
}

// validateLte is the validation function for validating if the current field's value is less than or equal to the param's value.
func validateLte(v, anotherField reflect.Value) (bool, error) {
	if !v.IsValid() || !anotherField.IsValid() {
		return false, fmt.Errorf("validator: Lte invalid reflection values")
	}
	if v.Kind() != anotherField.Kind() {
		return false, fmt.Errorf("validator: Lte The two fields must be of the same type %T, %T", v.Interface(), anotherField.Interface())
	}

	// Check for decimal.Decimal type first
	if d1, ok := asDecimal(v); ok {
		if d2, ok := asDecimal(anotherField); ok {
			return IsDecimalLte(d1, d2), nil
		}
	}

	var valid bool
	var err error

	switch v.Kind() {
	case reflect.String:
		valid, err = compareString(v.String(), int64(utf8.RuneCountInString(anotherField.String())), "<=")
	case reflect.Slice, reflect.Map, reflect.Array:
		valid, err = compareInt64(int64(v.Len()), int64(anotherField.Len()), "<=")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		valid, err = compareInt64(v.Int(), anotherField.Int(), "<=")
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		valid, err = compareUint64(v.Uint(), anotherField.Uint(), "<=")
	case reflect.Float32, reflect.Float64:
		valid, err = compareFloat64(v.Float(), anotherField.Float(), "<=")
	default:
		return false, fmt.Errorf("validator: Lte unsupported type %T", v.Interface())
	}

	return valid, err
}

// ValidateLte is the validation function for validating if the current field's value is less than or equal to the param's value.
func ValidateLte(i, a interface{}) (bool, error) {
	v := reflect.ValueOf(i)
	anotherField := reflect.ValueOf(a)
	return validateLte(v, anotherField)
}

// validateGt is the validation function for validating if the current field's value is greater than to the param's value.
func validateGt(v, anotherField reflect.Value) (bool, error) {
	if !v.IsValid() || !anotherField.IsValid() {
		return false, fmt.Errorf("validator: Gt invalid reflection values")
	}
	if v.Kind() != anotherField.Kind() {
		return false, fmt.Errorf("validator: Gt The two fields must be of the same type %T, %T", v.Interface(), anotherField.Interface())
	}

	// Check for decimal.Decimal type first
	if d1, ok := asDecimal(v); ok {
		if d2, ok := asDecimal(anotherField); ok {
			return IsDecimalGt(d1, d2), nil
		}
	}

	var valid bool
	var err error

	switch v.Kind() {
	case reflect.String:
		valid, err = compareString(v.String(), int64(utf8.RuneCountInString(anotherField.String())), ">")
	case reflect.Slice, reflect.Map, reflect.Array:
		valid, err = compareInt64(int64(v.Len()), int64(anotherField.Len()), ">")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		valid, err = compareInt64(v.Int(), anotherField.Int(), ">")
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		valid, err = compareUint64(v.Uint(), anotherField.Uint(), ">")
	case reflect.Float32, reflect.Float64:
		valid, err = compareFloat64(v.Float(), anotherField.Float(), ">")
	default:
		return false, fmt.Errorf("validator: Gt unsupported type %T", v.Interface())
	}

	return valid, err
}

// ValidateGt is the validation function for validating if the current field's value is greater than to the param's value.
func ValidateGt(i, a interface{}) (bool, error) {
	v := reflect.ValueOf(i)
	anotherField := reflect.ValueOf(a)
	return validateGt(v, anotherField)
}

// validateGte is the validation function for validating if the current field's value is greater than or equal to the param's value.
func validateGte(v, anotherField reflect.Value) (bool, error) {
	if !v.IsValid() || !anotherField.IsValid() {
		return false, fmt.Errorf("validator: Gte invalid reflection values")
	}
	if v.Kind() != anotherField.Kind() {
		return false, fmt.Errorf("validator: Gte The two fields must be of the same type %T, %T", v.Interface(), anotherField.Interface())
	}

	// Check for decimal.Decimal type first
	if d1, ok := asDecimal(v); ok {
		if d2, ok := asDecimal(anotherField); ok {
			return IsDecimalGte(d1, d2), nil
		}
	}

	var valid bool
	var err error

	switch v.Kind() {
	case reflect.String:
		valid, err = compareString(v.String(), int64(utf8.RuneCountInString(anotherField.String())), ">=")
	case reflect.Slice, reflect.Map, reflect.Array:
		valid, err = compareInt64(int64(v.Len()), int64(anotherField.Len()), ">=")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		valid, err = compareInt64(v.Int(), anotherField.Int(), ">=")
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		valid, err = compareUint64(v.Uint(), anotherField.Uint(), ">=")
	case reflect.Float32, reflect.Float64:
		valid, err = compareFloat64(v.Float(), anotherField.Float(), ">=")
	default:
		return false, fmt.Errorf("validator: Gte unsupported type %T", v.Interface())
	}

	return valid, err
}

// ValidateGte is the validation function for validating if the current field's value is greater than to the param's value.
func ValidateGte(i, a interface{}) (bool, error) {
	v := reflect.ValueOf(i)
	anotherField := reflect.ValueOf(a)
	return validateGte(v, anotherField)
}

// validateDistinct is the validation function for validating an attribute is unique among other values.
func validateDistinct(v reflect.Value) (bool, error) {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64:
		return true, nil
	case reflect.Slice, reflect.Array:
		m := reflect.MakeMap(reflect.MapOf(v.Type().Elem(), v.Type()))

		for i := 0; i < v.Len(); i++ {
			m.SetMapIndex(v.Index(i), v)
		}
		return v.Len() == m.Len(), nil
	case reflect.Map:
		m := reflect.MakeMap(reflect.MapOf(v.Type().Elem(), v.Type()))

		for _, k := range v.MapKeys() {
			m.SetMapIndex(v.MapIndex(k), v)
		}
		return v.Len() == m.Len(), nil
	}

	return false, fmt.Errorf("validator: Distinct unsupported type %T", v.Interface())
}

// ValidateDistinct is the validation function for validating an attribute is unique among other values.
func ValidateDistinct(i interface{}) bool {
	v := reflect.ValueOf(i)
	valid, _ := validateDistinct(v)
	return valid
}

// validateMultipleOf is the validation function for validating if the current field's value is a multiple of the param's value.
func validateMultipleOf(v reflect.Value, params []string) (bool, error) {
	if len(params) == 0 {
		return false, fmt.Errorf("validator: multipleOf rule requires at least one parameter")
	}

	// Check for decimal.Decimal type first
	if d, ok := asDecimal(v); ok {
		p, _, err := parseDecimalParams(params)
		if err != nil {
			return false, fmt.Errorf("validator: multipleOf decimal: %w", err)
		}
		if p.IsZero() {
			return false, fmt.Errorf("validator: multipleOf cannot divide by zero")
		}
		return d.Mod(p).IsZero(), nil
	}

	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p, err := ToInt(params[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for multipleOf rule: %w", err)
		}
		if p == 0 {
			return false, fmt.Errorf("validator: multipleOf cannot divide by zero")
		}
		return v.Int()%p == 0, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p, err := ToUint(params[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for multipleOf rule: %w", err)
		}
		if p == 0 {
			return false, fmt.Errorf("validator: multipleOf cannot divide by zero")
		}
		return v.Uint()%p == 0, nil
	case reflect.Float32, reflect.Float64:
		p, err := ToFloat(params[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for multipleOf rule: %w", err)
		}
		if p == 0 {
			return false, fmt.Errorf("validator: multipleOf cannot divide by zero")
		}
		remainder := v.Float() / p
		return remainder == float64(int64(remainder)), nil
	default:
		return false, fmt.Errorf("validator: multipleOf rule is not supported for type %s", v.Kind())
	}
}

// ValidateMultipleOf is the validation function for validating if a value is a multiple of another.
func ValidateMultipleOf(i interface{}, params []string) (bool, error) {
	v := reflect.ValueOf(i)
	return validateMultipleOf(v, params)
}

// validateMaxDigits is the validation function for validating the maximum number of digits.
func validateMaxDigits(v reflect.Value, params []string) (bool, error) {
	if len(params) == 0 {
		return false, fmt.Errorf("validator: maxDigits rule requires at least one parameter")
	}

	maxDigits, err := ToInt(params[0])
	if err != nil {
		return false, fmt.Errorf("validator: invalid parameter for maxDigits rule: %w", err)
	}

	var numStr string
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		numStr = strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		numStr = strconv.FormatUint(v.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		numStr = strconv.FormatFloat(v.Float(), 'f', -1, 64)
	case reflect.String:
		numStr = v.String()
	default:
		return false, fmt.Errorf("validator: maxDigits rule is not supported for type %s", v.Kind())
	}

	// Remove negative sign and decimal point for counting
	numStr = strings.TrimPrefix(numStr, "-")
	numStr = strings.ReplaceAll(numStr, ".", "")

	return int64(len(numStr)) <= maxDigits, nil
}

// ValidateMaxDigits is the validation function for validating the maximum number of digits.
func ValidateMaxDigits(i interface{}, params []string) (bool, error) {
	v := reflect.ValueOf(i)
	return validateMaxDigits(v, params)
}

// validateMinDigits is the validation function for validating the minimum number of digits.
func validateMinDigits(v reflect.Value, params []string) (bool, error) {
	if len(params) == 0 {
		return false, fmt.Errorf("validator: minDigits rule requires at least one parameter")
	}

	minDigits, err := ToInt(params[0])
	if err != nil {
		return false, fmt.Errorf("validator: invalid parameter for minDigits rule: %w", err)
	}

	var numStr string
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		numStr = strconv.FormatInt(v.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		numStr = strconv.FormatUint(v.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		numStr = strconv.FormatFloat(v.Float(), 'f', -1, 64)
	case reflect.String:
		numStr = v.String()
	default:
		return false, fmt.Errorf("validator: minDigits rule is not supported for type %s", v.Kind())
	}

	// Remove negative sign and decimal point for counting
	numStr = strings.TrimPrefix(numStr, "-")
	numStr = strings.ReplaceAll(numStr, ".", "")

	return int64(len(numStr)) >= minDigits, nil
}

// ValidateMinDigits is the validation function for validating the minimum number of digits.
func ValidateMinDigits(i interface{}, params []string) (bool, error) {
	v := reflect.ValueOf(i)
	return validateMinDigits(v, params)
}

// validateDecimalPrecision is the validation function for validating decimal places.
// params[0] is min decimal places, params[1] (optional) is max decimal places.
func validateDecimalPrecision(v reflect.Value, params []string) (bool, error) {
	if len(params) == 0 {
		return false, fmt.Errorf("validator: decimal rule requires at least one parameter")
	}

	minPlaces, err := ToInt(params[0])
	if err != nil {
		return false, fmt.Errorf("validator: invalid min parameter for decimal rule: %w", err)
	}

	maxPlaces := minPlaces
	if len(params) > 1 {
		maxPlaces, err = ToInt(params[1])
		if err != nil {
			return false, fmt.Errorf("validator: invalid max parameter for decimal rule: %w", err)
		}
	}

	var numStr string
	switch v.Kind() {
	case reflect.Float32, reflect.Float64:
		numStr = strconv.FormatFloat(v.Float(), 'f', -1, 64)
	case reflect.String:
		numStr = v.String()
	default:
		return false, fmt.Errorf("validator: decimal rule is not supported for type %s", v.Kind())
	}

	// Find decimal point
	parts := strings.Split(numStr, ".")
	if len(parts) == 1 {
		// No decimal point, 0 decimal places
		return minPlaces <= 0, nil
	}

	decimalPlaces := int64(len(parts[1]))
	return decimalPlaces >= minPlaces && decimalPlaces <= maxPlaces, nil
}

// ValidateDecimalPrecision is the validation function for validating decimal places.
func ValidateDecimalPrecision(i interface{}, params []string) (bool, error) {
	v := reflect.ValueOf(i)
	return validateDecimalPrecision(v, params)
}

// validateContains is the validation function for validating a string/array/slice contains specified values.
func validateContains(v reflect.Value, params []string) (bool, error) {
	if len(params) == 0 {
		return false, fmt.Errorf("validator: contains rule requires at least one parameter")
	}

	switch v.Kind() {
	case reflect.String:
		str := v.String()
		for _, param := range params {
			if !strings.Contains(str, param) {
				return false, nil
			}
		}
		return true, nil
	case reflect.Slice, reflect.Array:
		for _, param := range params {
			found := false
			for i := 0; i < v.Len(); i++ {
				elem := v.Index(i)
				if elem.Kind() == reflect.Interface || elem.Kind() == reflect.Ptr {
					elem = elem.Elem()
				}
				if ToString(elem.Interface()) == param {
					found = true
					break
				}
			}
			if !found {
				return false, nil
			}
		}
		return true, nil
	default:
		return false, fmt.Errorf("validator: contains rule is not supported for type %s", v.Kind())
	}
}

// ValidateContains is the validation function for validating an array/slice contains specified values.
func ValidateContains(i interface{}, params []string) (bool, error) {
	v := reflect.ValueOf(i)
	return validateContains(v, params)
}

// validateDoesntContain is the validation function for validating a string/array/slice does not contain specified values.
func validateDoesntContain(v reflect.Value, params []string) (bool, error) {
	if len(params) == 0 {
		return false, fmt.Errorf("validator: doesntContain rule requires at least one parameter")
	}

	switch v.Kind() {
	case reflect.String:
		str := v.String()
		for _, param := range params {
			if strings.Contains(str, param) {
				return false, nil
			}
		}
		return true, nil
	case reflect.Slice, reflect.Array:
		for _, param := range params {
			for i := 0; i < v.Len(); i++ {
				elem := v.Index(i)
				if elem.Kind() == reflect.Interface || elem.Kind() == reflect.Ptr {
					elem = elem.Elem()
				}
				if ToString(elem.Interface()) == param {
					return false, nil
				}
			}
		}
		return true, nil
	default:
		return false, fmt.Errorf("validator: doesntContain rule is not supported for type %s", v.Kind())
	}
}

// ValidateDoesntContain is the validation function for validating an array/slice does not contain specified values.
func ValidateDoesntContain(i interface{}, params []string) (bool, error) {
	v := reflect.ValueOf(i)
	return validateDoesntContain(v, params)
}

// ValidateMimeTypes is the validation function for the file must match one of the given MIME types.
func ValidateMimeTypes(data []byte, mimeTypes []string) bool {
	mimeType := http.DetectContentType(data)
	for _, value := range mimeTypes {
		if mimeType == value {
			return true
		}
	}
	return false
}

// ValidateMimes is the validation function for the file must have a MIME type corresponding to one of the listed extensions.
func ValidateMimes(data []byte, mimes []string) (bool, error) {
	mimeTypes := make([]string, len(mimes))
	for i, mime := range mimes {
		if val, ok := Mimes[mime]; ok {
			mimeTypes[i] = val
		} else {
			return false, fmt.Errorf("validator: Mimes unsupported type %s", mime)
		}
	}

	return ValidateMimeTypes(data, mimeTypes), nil
}

// ValidateImage is the validation function for the The file under validation must be an image (jpeg, png, bmp, gif, or svg)
func ValidateImage(data []byte) bool {
	v, err := ValidateMimes(data, []string{"jpeg", "png", "gif", "bmp", "svg"})
	if err != nil {
		return false
	}
	return v
}

// ValidateStruct use tags for fields.
// result will be equal to `false` if there are any errors.
func ValidateStruct(s interface{}) error {
	return Default.ValidateStruct(s, nil, nil)
}

// ValidateStruct use tags for fields.
// result will be equal to `false` if there are any errors.
func (v *Validator) ValidateStruct(s interface{}, jsonNamespace, structNamespace []byte) error {
	if s == nil {
		return nil
	}

	var err error

	val := reflect.ValueOf(s)
	if val.Kind() == reflect.Interface || val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	// we only accept structs
	if val.Kind() != reflect.Struct {
		return fmt.Errorf("function only accepts structs; got %s", val.Kind())
	}

	var errs Errors
	fields := cachedTypefields(val.Type())

	// Pre-allocate slice capacity to reduce allocations
	if len(fields) > 0 {
		errs = make(Errors, 0, len(fields)/2) // Assume ~50% will have validation errors
	}

	//nolint:gocritic // Field struct copying is acceptable for validation library performance
	for _, f := range fields {
		valuefield := val.Field(f.index[0])
		err := v.newTypeValidator(valuefield, &f, val, jsonNamespace, structNamespace)
		if err != nil {
			if errors, ok := err.(Errors); ok {
				errs = append(errs, errors...)
			} else {
				errs = append(errs, err)
			}
		}
	}

	if len(errs) > 0 {
		err = errs
	}

	return err
}

func (v *Validator) newTypeValidator(value reflect.Value, f *field, o reflect.Value, jsonNamespace, structNamespace []byte) (resultErr error) {
	if !value.IsValid() || (f.omitEmpty && Empty(value)) {
		return nil
	}

	name := buildFieldName(jsonNamespace, f.nameBytes)
	structName := buildFieldName(structNamespace, f.structNameBytes)

	// Handle pointer and interface dereferencing
	if value.Kind() == reflect.Interface || value.Kind() == reflect.Ptr {
		if err := v.checkRequired(value, f, o, name, structName); err != nil {
			return err
		}
		if value.IsNil() {
			return nil
		}
		value = value.Elem()
	} else if err := v.checkRequired(value, f, o, name, structName); err != nil {
		return err
	}

	// If nullable and empty, skip remaining validations (required already checked above)
	if f.nullable && Empty(value) {
		return nil
	}

	// Validate custom type rules
	if err := v.validateCustomTypeRules(f.validTags, value, f, name, structName, o); err != nil {
		return err
	}

	switch value.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.String:

		if err := v.validateCommonRules(f.validTags, value, f, name, structName, o); err != nil {
			return err
		}
		return nil
	case reflect.Map:
		// Validate map-specific rules (without string-specific rules)
		for _, tag := range f.validTags {
			handled, err := v.checkDependentRulesWithStatus(tag, f, value, o, name, structName)
			if err != nil {
				return err
			}

			// Skip ParamRuleMap for comparison rules if they were handled by field comparison
			skipParamRule := handled && (tag.name == "gt" || tag.name == "gte" || tag.name == "lt" || tag.name == "lte")

			if err := v.validateWithRuleMap(tag, value, f, name, structName, o); err != nil {
				return err
			}

			if !skipParamRule {
				if err := v.validateWithParamRuleMap(tag, value, f, name, structName, o); err != nil {
					return err
				}
			}
		}
		return v.validateMapFields(value, f, jsonNamespace, structNamespace)
	case reflect.Slice, reflect.Array:
		// Validate slice/array-specific rules (without string-specific rules)
		for _, tag := range f.validTags {
			handled, err := v.checkDependentRulesWithStatus(tag, f, value, o, name, structName)
			if err != nil {
				return err
			}

			// Skip ParamRuleMap for comparison rules if they were handled by field comparison
			skipParamRule := handled && (tag.name == "gt" || tag.name == "gte" || tag.name == "lt" || tag.name == "lte")

			if err := v.validateWithRuleMap(tag, value, f, name, structName, o); err != nil {
				return err
			}

			if !skipParamRule {
				if err := v.validateWithParamRuleMap(tag, value, f, name, structName, o); err != nil {
					return err
				}
			}
		}
		return v.validateSliceFields(value, f, jsonNamespace, structNamespace)
	case reflect.Struct:
		// Check for decimal.Decimal type - validate it like a numeric type
		if _, ok := asDecimal(value); ok {
			if err := v.validateCommonRules(f.validTags, value, f, name, structName, o); err != nil {
				return err
			}
			return nil
		}
		// Check for time.Time type - validate with date validators
		if _, ok := getTimeValue(value); ok {
			if err := v.validateCommonRules(f.validTags, value, f, name, structName, o); err != nil {
				return err
			}
			return nil
		}
		// Regular struct - recursively validate
		jsonNamespace = append(append(jsonNamespace, f.nameBytes...), '.')
		structNamespace = append(append(structNamespace, f.structNameBytes...), '.')
		return v.ValidateStruct(value.Interface(), jsonNamespace, structNamespace)
	default:
		// For unsupported types with validation tags, return a FieldError with FuncError
		if len(f.validTags) > 0 {
			unsupportedErr := &UnsupportedTypeError{value.Type()}
			// Return the first validation tag's error with the unsupported type as FuncError
			return v.formatsMessages(&FieldError{
				Name:              name,
				StructName:        structName,
				Tag:               f.validTags[0].name,
				MessageName:       f.validTags[0].messageName,
				MessageParameters: parseValidatorMessageParameters(f.validTags[0], o),
				Attribute:         f.attribute,
				DefaultAttribute:  f.defaultAttribute,
				Value:             ToString(value.Interface()),
				FuncError:         unsupportedErr,
			})
		}
		return &UnsupportedTypeError{value.Type()}
	}
}

// Empty determine whether a variable is empty
func Empty(v reflect.Value) bool {
	// Check for decimal.Decimal type first
	if d, ok := asDecimal(v); ok {
		return d.IsZero()
	}

	switch v.Kind() {
	case reflect.String, reflect.Array:
		return v.Len() == 0
	case reflect.Map, reflect.Slice:
		return v.Len() == 0 || v.IsNil()
	case reflect.Bool:
		return !v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}

	return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}

// Error returns string equivalent for reflect.Type
func (e *UnsupportedTypeError) Error() string {
	return "validator: unsupported type: " + e.Type.String()
}

func (sv stringValues) Len() int           { return len(sv) }
func (sv stringValues) Swap(i, j int)      { sv[i], sv[j] = sv[j], sv[i] }
func (sv stringValues) Less(i, j int) bool { return sv.get(i) < sv.get(j) }
func (sv stringValues) get(i int) string   { return sv[i].String() }

// validateRequired check value required when anotherField str is a member of the set of strings params
func validateRequired(v reflect.Value) bool {
	return !Empty(v)
}

// ValidateRequired check value required when anotherField str is a member of the set of strings params
func ValidateRequired(i interface{}) bool {
	v := reflect.ValueOf(i)
	return validateRequired(v)
}

// acceptedValues are the values that are considered "accepted"
var acceptedValues = []string{"yes", "on", "1", "true"}

// declinedValues are the values that are considered "declined"
var declinedValues = []string{"no", "off", "0", "false"}

// validateAccepted checks if the field is accepted (yes, on, 1, "1", true, "true")
func validateAccepted(v reflect.Value) bool {
	if v.Kind() == reflect.Bool {
		return v.Bool()
	}
	value := strings.ToLower(ToString(v.Interface()))
	return InString(value, acceptedValues)
}

// ValidateAccepted checks if the field is accepted
func ValidateAccepted(i interface{}) bool {
	v := reflect.ValueOf(i)
	return validateAccepted(v)
}

// validateDeclined checks if the field is declined (no, off, 0, "0", false, "false")
func validateDeclined(v reflect.Value) bool {
	if v.Kind() == reflect.Bool {
		return !v.Bool()
	}
	value := strings.ToLower(ToString(v.Interface()))
	return InString(value, declinedValues)
}

// ValidateDeclined checks if the field is declined
func ValidateDeclined(i interface{}) bool {
	v := reflect.ValueOf(i)
	return validateDeclined(v)
}

// validateProhibited checks if the field is empty (prohibited means it must be empty)
func validateProhibited(v reflect.Value) bool {
	return Empty(v)
}

// ValidateProhibited checks if the field is prohibited (must be empty)
func ValidateProhibited(i interface{}) bool {
	v := reflect.ValueOf(i)
	return validateProhibited(v)
}

// validateAcceptedIf checks if the field is accepted when another field equals a specific value
func validateAcceptedIf(v, anotherField reflect.Value, params []string) (bool, error) {
	if anotherField.Kind() == reflect.Interface || anotherField.Kind() == reflect.Ptr {
		anotherField = anotherField.Elem()
	}

	if !anotherField.IsValid() {
		return true, nil
	}

	value := ToString(anotherField.Interface())
	if InString(value, params) {
		return validateAccepted(v), nil
	}
	return true, nil
}

// validateDeclinedIf checks if the field is declined when another field equals a specific value
func validateDeclinedIf(v, anotherField reflect.Value, params []string) (bool, error) {
	if anotherField.Kind() == reflect.Interface || anotherField.Kind() == reflect.Ptr {
		anotherField = anotherField.Elem()
	}

	if !anotherField.IsValid() {
		return true, nil
	}

	value := ToString(anotherField.Interface())
	if InString(value, params) {
		return validateDeclined(v), nil
	}
	return true, nil
}

// validateProhibitedIf checks if the field is empty when another field equals a specific value
func validateProhibitedIf(v, anotherField reflect.Value, params []string) (bool, error) {
	if anotherField.Kind() == reflect.Interface || anotherField.Kind() == reflect.Ptr {
		anotherField = anotherField.Elem()
	}

	if !anotherField.IsValid() {
		return true, nil
	}

	value := ToString(anotherField.Interface())
	if InString(value, params) {
		return validateProhibited(v), nil
	}
	return true, nil
}

// validateProhibitedUnless checks if the field is empty unless another field equals a specific value
func validateProhibitedUnless(v, anotherField reflect.Value, params []string) (bool, error) {
	if anotherField.Kind() == reflect.Interface || anotherField.Kind() == reflect.Ptr {
		anotherField = anotherField.Elem()
	}

	if !anotherField.IsValid() {
		return true, nil
	}

	value := ToString(anotherField.Interface())
	if !InString(value, params) {
		return validateProhibited(v), nil
	}
	return true, nil
}

// validateMissing checks if the field is not present (must be absent from the data)
func validateMissing(v reflect.Value) bool {
	return !v.IsValid() || Empty(v)
}

// ValidateMissing checks if the field is missing
func ValidateMissing(i interface{}) bool {
	v := reflect.ValueOf(i)
	return validateMissing(v)
}

// validateMissingIf checks if the field is missing when another field equals a specific value
func validateMissingIf(v, anotherField reflect.Value, params []string) (bool, error) {
	if anotherField.Kind() == reflect.Interface || anotherField.Kind() == reflect.Ptr {
		anotherField = anotherField.Elem()
	}

	if !anotherField.IsValid() {
		return true, nil
	}

	value := ToString(anotherField.Interface())
	if InString(value, params) {
		return validateMissing(v), nil
	}
	return true, nil
}

// validateMissingUnless checks if the field is missing unless another field equals a specific value
func validateMissingUnless(v, anotherField reflect.Value, params []string) (bool, error) {
	if anotherField.Kind() == reflect.Interface || anotherField.Kind() == reflect.Ptr {
		anotherField = anotherField.Elem()
	}

	if !anotherField.IsValid() {
		return true, nil
	}

	value := ToString(anotherField.Interface())
	if !InString(value, params) {
		return validateMissing(v), nil
	}
	return true, nil
}

// validateRequiredIfAccepted checks if the field is required when another field is accepted
func validateRequiredIfAccepted(v, anotherField reflect.Value) (bool, error) {
	if validateAccepted(anotherField) {
		return validateRequired(v), nil
	}
	return true, nil
}

// validateRequiredIfDeclined checks if the field is required when another field is declined
func validateRequiredIfDeclined(v, anotherField reflect.Value) (bool, error) {
	if validateDeclined(anotherField) {
		return validateRequired(v), nil
	}
	return true, nil
}

// validateProhibitedIfAccepted checks if the field is prohibited when another field is accepted
func validateProhibitedIfAccepted(v, anotherField reflect.Value) (bool, error) {
	if validateAccepted(anotherField) {
		return validateProhibited(v), nil
	}
	return true, nil
}

// validateProhibitedIfDeclined checks if the field is prohibited when another field is declined
func validateProhibitedIfDeclined(v, anotherField reflect.Value) (bool, error) {
	if validateDeclined(anotherField) {
		return validateProhibited(v), nil
	}
	return true, nil
}

// validateProhibits checks if when this field has a value, the other specified fields must be empty
func validateProhibits(v reflect.Value, otherFields []string, obj reflect.Value) bool {
	if Empty(v) {
		return true
	}
	for _, fieldName := range otherFields {
		otherField, err := findField(fieldName, obj)
		if err != nil {
			continue
		}
		if !Empty(otherField) {
			return false
		}
	}
	return true
}

// validateMissingWith checks if the field is missing when any of the specified fields are present
func validateMissingWith(v reflect.Value, otherFields []string, obj reflect.Value) bool {
	for _, fieldName := range otherFields {
		otherField, err := findField(fieldName, obj)
		if err != nil {
			continue
		}
		if !Empty(otherField) {
			return validateMissing(v)
		}
	}
	return true
}

// validateMissingWithAll checks if the field is missing when all of the specified fields are present
func validateMissingWithAll(v reflect.Value, otherFields []string, obj reflect.Value) bool {
	allPresent := true
	for _, fieldName := range otherFields {
		otherField, err := findField(fieldName, obj)
		if err != nil {
			allPresent = false
			break
		}
		if Empty(otherField) {
			allPresent = false
			break
		}
	}
	if allPresent {
		return validateMissing(v)
	}
	return true
}

// validatePresent checks if the field is present (for struct fields, always true since they exist)
func validatePresent(v reflect.Value) bool {
	// In a struct context, all fields are always "present" - they exist in the struct
	// This is different from "required" which checks for non-empty values
	return v.IsValid()
}

// ValidatePresent checks if the field is present
func ValidatePresent(i interface{}) bool {
	v := reflect.ValueOf(i)
	return validatePresent(v)
}

// validatePresentIf checks if the field is present (non-empty) when another field equals a specific value
func validatePresentIf(v, anotherField reflect.Value, params []string) (bool, error) {
	if anotherField.Kind() == reflect.Interface || anotherField.Kind() == reflect.Ptr {
		anotherField = anotherField.Elem()
	}

	if !anotherField.IsValid() {
		return true, nil
	}

	value := ToString(anotherField.Interface())
	if InString(value, params) {
		return !Empty(v), nil
	}
	return true, nil
}

// validatePresentUnless checks if the field is present (non-empty) unless another field equals a specific value
func validatePresentUnless(v, anotherField reflect.Value, params []string) (bool, error) {
	if anotherField.Kind() == reflect.Interface || anotherField.Kind() == reflect.Ptr {
		anotherField = anotherField.Elem()
	}

	if !anotherField.IsValid() {
		return true, nil
	}

	value := ToString(anotherField.Interface())
	if !InString(value, params) {
		return !Empty(v), nil
	}
	return true, nil
}

// validatePresentWith checks if the field is present (non-empty) when any of the specified fields are present
func validatePresentWith(v reflect.Value, otherFields []string, obj reflect.Value) bool {
	for _, fieldName := range otherFields {
		otherField, err := findField(fieldName, obj)
		if err != nil {
			continue
		}
		if !Empty(otherField) {
			return !Empty(v)
		}
	}
	return true
}

// validatePresentWithAll checks if the field is present (non-empty) when all of the specified fields are present
func validatePresentWithAll(v reflect.Value, otherFields []string, obj reflect.Value) bool {
	allPresent := true
	for _, fieldName := range otherFields {
		otherField, err := findField(fieldName, obj)
		if err != nil {
			allPresent = false
			break
		}
		if Empty(otherField) {
			allPresent = false
			break
		}
	}
	if allPresent {
		return !Empty(v)
	}
	return true
}

// validateRequiredArrayKeys checks if the array has all the specified keys
func validateRequiredArrayKeys(v reflect.Value, keys []string) (bool, error) {
	if v.Kind() != reflect.Map {
		return false, fmt.Errorf("validator: requiredArrayKeys only supports map types")
	}

	for _, key := range keys {
		found := false
		for _, k := range v.MapKeys() {
			if ToString(k.Interface()) == key {
				found = true
				break
			}
		}
		if !found {
			return false, nil
		}
	}
	return true, nil
}

// ValidateRequiredArrayKeys checks if the map has all the specified keys
func ValidateRequiredArrayKeys(i interface{}, keys []string) (bool, error) {
	v := reflect.ValueOf(i)
	return validateRequiredArrayKeys(v, keys)
}

// validateList checks if the slice/array is a valid list (sequential integer keys starting from 0)
// In Go, slices are always lists, so this just checks if it's a slice/array
func validateList(v reflect.Value) (bool, error) {
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		return true, nil
	default:
		return false, nil
	}
}

// ValidateList checks if the value is a list (slice or array)
func ValidateList(i interface{}) (bool, error) {
	v := reflect.ValueOf(i)
	return validateList(v)
}

// Common date formats to try when parsing
var dateFormats = []string{
	time.RFC3339,
	"2006-01-02T15:04:05",
	"2006-01-02 15:04:05",
	"2006-01-02",
	"02/01/2006",
	"01/02/2006",
	"2006/01/02",
	time.RFC1123,
	time.RFC822,
}

// parseDate attempts to parse a date string using common formats
// It uses local timezone for formats without timezone info
func parseDate(dateStr string) (time.Time, error) {
	for _, format := range dateFormats {
		// Use ParseInLocation for formats without timezone to use local time
		if t, err := time.ParseInLocation(format, dateStr, time.Local); err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unable to parse date: %s", dateStr)
}

// getTimeValue extracts a time.Time from a reflect.Value
func getTimeValue(v reflect.Value) (time.Time, bool) {
	if v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if !v.IsValid() {
		return time.Time{}, false
	}

	// Check if it's a time.Time
	if t, ok := v.Interface().(time.Time); ok {
		return t, true
	}

	// Try to parse as string
	if v.Kind() == reflect.String {
		str := v.String()
		if str == "" {
			return time.Time{}, false
		}
		if t, err := parseDate(str); err == nil {
			return t, true
		}
	}

	return time.Time{}, false
}

// validateDate checks if the value is a valid date
func validateDate(v reflect.Value) (bool, error) {
	if v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if !v.IsValid() || Empty(v) {
		return true, nil
	}

	_, ok := getTimeValue(v)
	return ok, nil
}

// ValidateDate checks if the value is a valid date
func ValidateDate(i interface{}) (bool, error) {
	v := reflect.ValueOf(i)
	return validateDate(v)
}

// validateDateFormat checks if the value matches the given date format
func validateDateFormat(v reflect.Value, params []string) (bool, error) {
	if len(params) == 0 {
		return false, fmt.Errorf("validator: dateFormat rule requires a format parameter")
	}

	if v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if !v.IsValid() || Empty(v) {
		return true, nil
	}

	format := params[0]
	var dateStr string

	if v.Kind() == reflect.String {
		dateStr = v.String()
	} else if t, ok := v.Interface().(time.Time); ok {
		// For time.Time, always valid
		_ = t
		return true, nil
	} else {
		return false, nil
	}

	_, err := time.Parse(format, dateStr)
	return err == nil, nil
}

// ValidateDateFormat checks if the value matches the given date format
func ValidateDateFormat(i interface{}, params []string) (bool, error) {
	v := reflect.ValueOf(i)
	return validateDateFormat(v, params)
}

// validateAfter checks if the date is after the given date
func validateAfter(v reflect.Value, params []string) (bool, error) {
	if len(params) == 0 {
		return false, fmt.Errorf("validator: after rule requires a date parameter")
	}

	if v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if !v.IsValid() || Empty(v) {
		return true, nil
	}

	currentTime, ok := getTimeValue(v)
	if !ok {
		return false, nil
	}

	// Use parseDateParam which supports: today, tomorrow, yesterday, now, today+7d, today-1m, today-18y
	compareTime, err := parseDateParam(params[0])
	if err != nil {
		return false, fmt.Errorf("validator: invalid date parameter for after rule: %s", params[0])
	}

	return currentTime.After(compareTime), nil
}

// ValidateAfter checks if the date is after the given date
func ValidateAfter(i interface{}, params []string) (bool, error) {
	v := reflect.ValueOf(i)
	return validateAfter(v, params)
}

// validateAfterOrEqual checks if the date is after or equal to the given date
func validateAfterOrEqual(v reflect.Value, params []string) (bool, error) {
	if len(params) == 0 {
		return false, fmt.Errorf("validator: afterOrEqual rule requires a date parameter")
	}

	if v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if !v.IsValid() || Empty(v) {
		return true, nil
	}

	currentTime, ok := getTimeValue(v)
	if !ok {
		return false, nil
	}

	// Use parseDateParam which supports: today, tomorrow, yesterday, now, today+7d, today-1m, today-18y
	compareTime, err := parseDateParam(params[0])
	if err != nil {
		return false, fmt.Errorf("validator: invalid date parameter for afterOrEqual rule: %s", params[0])
	}

	return currentTime.After(compareTime) || currentTime.Equal(compareTime), nil
}

// ValidateAfterOrEqual checks if the date is after or equal to the given date
func ValidateAfterOrEqual(i interface{}, params []string) (bool, error) {
	v := reflect.ValueOf(i)
	return validateAfterOrEqual(v, params)
}

// validateBefore checks if the date is before the given date
func validateBefore(v reflect.Value, params []string) (bool, error) {
	if len(params) == 0 {
		return false, fmt.Errorf("validator: before rule requires a date parameter")
	}

	if v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if !v.IsValid() || Empty(v) {
		return true, nil
	}

	currentTime, ok := getTimeValue(v)
	if !ok {
		return false, nil
	}

	// Use parseDateParam which supports: today, tomorrow, yesterday, now, today+7d, today-1m, today-18y
	compareTime, err := parseDateParam(params[0])
	if err != nil {
		return false, fmt.Errorf("validator: invalid date parameter for before rule: %s", params[0])
	}

	return currentTime.Before(compareTime), nil
}

// ValidateBefore checks if the date is before the given date
func ValidateBefore(i interface{}, params []string) (bool, error) {
	v := reflect.ValueOf(i)
	return validateBefore(v, params)
}

// validateBeforeOrEqual checks if the date is before or equal to the given date
func validateBeforeOrEqual(v reflect.Value, params []string) (bool, error) {
	if len(params) == 0 {
		return false, fmt.Errorf("validator: beforeOrEqual rule requires a date parameter")
	}

	if v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if !v.IsValid() || Empty(v) {
		return true, nil
	}

	currentTime, ok := getTimeValue(v)
	if !ok {
		return false, nil
	}

	// Use parseDateParam which supports: today, tomorrow, yesterday, now, today+7d, today-1m, today-18y
	compareTime, err := parseDateParam(params[0])
	if err != nil {
		return false, fmt.Errorf("validator: invalid date parameter for beforeOrEqual rule: %s", params[0])
	}

	return currentTime.Before(compareTime) || currentTime.Equal(compareTime), nil
}

// ValidateBeforeOrEqual checks if the date is before or equal to the given date
func ValidateBeforeOrEqual(i interface{}, params []string) (bool, error) {
	v := reflect.ValueOf(i)
	return validateBeforeOrEqual(v, params)
}

// validateIn checks if the value is in the given list
func validateIn(v reflect.Value, params []string) (bool, error) {
	if len(params) == 0 {
		return false, fmt.Errorf("validator: in rule requires at least one parameter")
	}

	if v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if !v.IsValid() || Empty(v) {
		return true, nil
	}

	value := ToString(v.Interface())
	return InString(value, params), nil
}

// ValidateIn checks if the value is in the given list
func ValidateIn(i interface{}, params []string) (bool, error) {
	v := reflect.ValueOf(i)
	return validateIn(v, params)
}

// validateNotIn checks if the value is not in the given list
func validateNotIn(v reflect.Value, params []string) (bool, error) {
	if len(params) == 0 {
		return false, fmt.Errorf("validator: notIn rule requires at least one parameter")
	}

	if v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if !v.IsValid() || Empty(v) {
		return true, nil
	}

	value := ToString(v.Interface())
	return !InString(value, params), nil
}

// ValidateNotIn checks if the value is not in the given list
func ValidateNotIn(i interface{}, params []string) (bool, error) {
	v := reflect.ValueOf(i)
	return validateNotIn(v, params)
}

// validateDifferent checks if the value is different from another field
func validateDifferent(v, anotherField reflect.Value) (bool, error) {
	if v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if anotherField.Kind() == reflect.Interface || anotherField.Kind() == reflect.Ptr {
		anotherField = anotherField.Elem()
	}

	if !v.IsValid() || !anotherField.IsValid() {
		return true, nil
	}

	return ToString(v.Interface()) != ToString(anotherField.Interface()), nil
}

// ValidateDifferent checks if the value is different from another value
func ValidateDifferent(a, b interface{}) (bool, error) {
	return validateDifferent(reflect.ValueOf(a), reflect.ValueOf(b))
}

// validateConfirmed checks if the field has a matching confirmation field
// This is handled specially in the validation loop since it needs to find {field}_confirmation
func validateConfirmed(v, confirmationField reflect.Value) (bool, error) {
	if v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if confirmationField.Kind() == reflect.Interface || confirmationField.Kind() == reflect.Ptr {
		confirmationField = confirmationField.Elem()
	}

	if !v.IsValid() || Empty(v) {
		return true, nil
	}

	if !confirmationField.IsValid() {
		return false, nil
	}

	return ToString(v.Interface()) == ToString(confirmationField.Interface()), nil
}

// validateJSON checks if the value is a valid JSON string
func validateJSON(v reflect.Value) (bool, error) {
	if v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if !v.IsValid() || Empty(v) {
		return true, nil
	}

	if v.Kind() != reflect.String {
		return false, nil
	}

	var js interface{}
	return json.Unmarshal([]byte(v.String()), &js) == nil, nil
}

// ValidateJSON checks if the value is a valid JSON string
func ValidateJSON(i interface{}) (bool, error) {
	v := reflect.ValueOf(i)
	return validateJSON(v)
}

// validateRegex checks if the value matches the given regex pattern
func validateRegex(v reflect.Value, params []string) (bool, error) {
	if len(params) == 0 {
		return false, fmt.Errorf("validator: regex rule requires a pattern parameter")
	}

	if v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if !v.IsValid() || Empty(v) {
		return true, nil
	}

	if v.Kind() != reflect.String {
		return false, nil
	}

	pattern := params[0]
	re, err := getCompiledRegex(pattern)
	if err != nil {
		return false, fmt.Errorf("validator: invalid regex pattern: %w", err)
	}

	return re.MatchString(v.String()), nil
}

// ValidateRegex checks if the value matches the given regex pattern
func ValidateRegex(i interface{}, params []string) (bool, error) {
	v := reflect.ValueOf(i)
	return validateRegex(v, params)
}

// validateNotRegex checks if the value does not match the given regex pattern
func validateNotRegex(v reflect.Value, params []string) (bool, error) {
	if len(params) == 0 {
		return false, fmt.Errorf("validator: notRegex rule requires a pattern parameter")
	}

	if v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if !v.IsValid() || Empty(v) {
		return true, nil
	}

	if v.Kind() != reflect.String {
		return true, nil
	}

	pattern := params[0]
	re, err := getCompiledRegex(pattern)
	if err != nil {
		return false, fmt.Errorf("validator: invalid regex pattern: %w", err)
	}

	return !re.MatchString(v.String()), nil
}

// ValidateNotRegex checks if the value does not match the given regex pattern
func ValidateNotRegex(i interface{}, params []string) (bool, error) {
	v := reflect.ValueOf(i)
	return validateNotRegex(v, params)
}

// validateBoolean checks if the value is boolean-like (true/false, 1/0, "yes"/"no", "on"/"off")
func validateBoolean(v reflect.Value) (bool, error) {
	if v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if !v.IsValid() || Empty(v) {
		return true, nil
	}

	switch v.Kind() {
	case reflect.Bool:
		return true, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		val := v.Int()
		return val == 0 || val == 1, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val := v.Uint()
		return val == 0 || val == 1, nil
	case reflect.String:
		str := strings.ToLower(v.String())
		return str == "true" || str == "false" || str == "1" || str == "0" ||
			str == "yes" || str == "no" || str == "on" || str == "off", nil
	}

	return false, nil
}

// ValidateBoolean checks if the value is boolean-like
func ValidateBoolean(i interface{}) (bool, error) {
	v := reflect.ValueOf(i)
	return validateBoolean(v)
}

// validateString checks if the value is a string
func validateString(v reflect.Value) (bool, error) {
	if v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if !v.IsValid() {
		return true, nil
	}

	return v.Kind() == reflect.String, nil
}

// ValidateString checks if the value is a string
func ValidateString(i interface{}) (bool, error) {
	v := reflect.ValueOf(i)
	return validateString(v)
}

// validateArray checks if the value is an array or slice
func validateArray(v reflect.Value) (bool, error) {
	if v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	if !v.IsValid() {
		return true, nil
	}

	return v.Kind() == reflect.Slice || v.Kind() == reflect.Array, nil
}

// ValidateArray checks if the value is an array or slice
func ValidateArray(i interface{}) (bool, error) {
	v := reflect.ValueOf(i)
	return validateArray(v)
}

// validateFilled checks if the field is not empty when it is present
func validateFilled(v reflect.Value) (bool, error) {
	if v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	// If the field is not valid (doesn't exist), it passes
	if !v.IsValid() {
		return true, nil
	}
	// If the field exists, it must not be empty
	return !Empty(v), nil
}

// ValidateFilled checks if the field is not empty when present
func ValidateFilled(i interface{}) (bool, error) {
	v := reflect.ValueOf(i)
	return validateFilled(v)
}

// validateRequiredIf check value required when anotherField str is a member of the set of strings params
func validateRequiredIf(v, anotherField reflect.Value, params []string, tag *ValidTag) (bool, error) {
	if anotherField.Kind() == reflect.Interface || anotherField.Kind() == reflect.Ptr {
		anotherField = anotherField.Elem()
	}

	if !anotherField.IsValid() {
		return true, nil
	}

	switch anotherField.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.String:
		value := ToString(anotherField)
		if InString(value, params) && Empty(v) && tag != nil {
			tag.messageParameters = append(
				tag.messageParameters,
				messageParameter{
					Key:   "Value",
					Value: value,
				},
			)
			return false, nil
		}
	case reflect.Map, reflect.Slice, reflect.Array:
		values, err := extractValuesFromCollection(anotherField)
		if err != nil {
			return false, err
		}
		return checkRequiredIfCondition(v, values, params, tag)
	default:
		return false, fmt.Errorf("validator: RequiredIf unsupported type %T", anotherField.Interface())
	}

	return true, nil
}

// validateRequiredUnless check value required when anotherField str is a member of the set of strings params
func validateRequiredUnless(v, anotherField reflect.Value, params []string) (bool, error) {
	if anotherField.Kind() == reflect.Interface || anotherField.Kind() == reflect.Ptr {
		anotherField = anotherField.Elem()
	}

	if !anotherField.IsValid() {
		return true, nil
	}

	switch anotherField.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.String:
		value := ToString(anotherField)
		if !InString(value, params) {
			if Empty(v) {
				return false, nil
			}
		}
	case reflect.Map, reflect.Slice, reflect.Array:
		values, err := extractValuesFromCollection(anotherField)
		if err != nil {
			return false, err
		}
		for _, value := range values {
			if !InString(value, params) {
				if Empty(v) {
					return false, nil
				}
			}
		}
	default:
		return false, fmt.Errorf("validator: requiredUnless unsupported type %T", anotherField.Interface())
	}

	return true, nil
}

// allFailingRequired validate that an attribute exists when all other attributes do not.
func allFailingRequired(parameters []string, v reflect.Value) bool {
	for _, p := range parameters {
		anotherField, err := findField(p, v)
		if err != nil {
			continue
		}
		if !Empty(anotherField) {
			return false
		}
	}
	return true
}

// anyFailingRequired determine if any of the given attributes fail the required test.
func anyFailingRequired(parameters []string, v reflect.Value) bool {
	for _, p := range parameters {
		anotherField, err := findField(p, v)
		if err != nil {
			return true
		}
		if Empty(anotherField) {
			return true
		}
	}
	return false
}

func (v *Validator) checkRequired(value reflect.Value, f *field, o reflect.Value, name, structName string) *FieldError {
	for _, tag := range f.requiredTags {
		var funcError error
		isError := false
		var isValid bool
		switch tag.name {
		case "required":
			isError = !validateRequired(value)
		case "requiredIf":
			if len(tag.params) == 0 {
				continue
			}
			anotherField, err := findField(tag.params[0], o)
			if err == nil && len(tag.params) >= 2 {
				isValid, funcError = validateRequiredIf(value, anotherField, tag.params[1:], tag)
				if !isValid {
					isError = true
				}
			}
		case "requiredUnless":
			if len(tag.params) == 0 {
				continue
			}
			anotherField, err := findField(tag.params[0], o)
			if err == nil && len(tag.params) >= 2 {
				isValid, funcError = validateRequiredUnless(value, anotherField, tag.params[1:])
				if !isValid {
					isError = true
				}
			}
		case "requiredWith":
			if !validateRequiredWith(tag.params, value, o) {
				isError = true
			}
		case "requiredWithAll":
			if !validateRequiredWithAll(tag.params, value, o) {
				isError = true
			}
		case "requiredWithout":
			if !validateRequiredWithout(tag.params, value, o) {
				isError = true
			}
		case "requiredWithoutAll":
			if !validateRequiredWithoutAll(tag.params, value, o) {
				isError = true
			}
		case "requiredIfAccepted":
			if len(tag.params) == 0 {
				continue
			}
			anotherField, err := findField(tag.params[0], o)
			if err == nil {
				isValid, funcError = validateRequiredIfAccepted(value, anotherField)
				if !isValid {
					isError = true
				}
			}
		case "requiredIfDeclined":
			if len(tag.params) == 0 {
				continue
			}
			anotherField, err := findField(tag.params[0], o)
			if err == nil {
				isValid, funcError = validateRequiredIfDeclined(value, anotherField)
				if !isValid {
					isError = true
				}
			}
		}

		if isError {
			return v.formatsMessages(&FieldError{
				Name:              name,
				StructName:        structName,
				Tag:               tag.name,
				MessageName:       tag.messageName,
				MessageParameters: parseValidatorMessageParameters(tag, o),
				Attribute:         f.attribute,
				DefaultAttribute:  f.defaultAttribute,
				Value:             ToString(value.Interface()),
				FuncError:         funcError,
			})
		}
	}

	return nil
}

// validateRequiredWith The field under validation must be present and not empty only if any of the other specified fields are present.
func validateRequiredWith(otherFields []string, currentField, obj reflect.Value) bool {
	if !allFailingRequired(otherFields, obj) {
		return validateRequired(currentField)
	}
	return true
}

// validateRequiredWithAll The field under validation must be present and not empty only if all of the other specified fields are present.
func validateRequiredWithAll(otherFields []string, currentField, obj reflect.Value) bool {
	if !anyFailingRequired(otherFields, obj) {
		return validateRequired(currentField)
	}
	return true
}

// RequiredWithout The field under validation must be present and not empty only when any of the other specified fields are not present.
func validateRequiredWithout(otherFields []string, currentField, obj reflect.Value) bool {
	if anyFailingRequired(otherFields, obj) {
		return validateRequired(currentField)
	}
	return true
}

// validateRequiredWithoutAll The field under validation must be present and not empty only when all of the other specified fields are not present.
func validateRequiredWithoutAll(otherFields []string, currentField, obj reflect.Value) bool {
	if allFailingRequired(otherFields, obj) {
		return validateRequired(currentField)
	}
	return true
}

func parseValidatorMessageParameters(validTag *ValidTag, o reflect.Value) MessageParameters {
	messageParameters := validTag.messageParameters
	switch validTag.name {
	case "requiredWith", "requiredWithAll", "requiredWithout", "requiredWithoutAll":
		first := true
		var buff bytes.Buffer
		for _, v := range validTag.params {
			if first {
				first = false
			} else {
				buff.WriteByte(' ')
				buff.WriteByte('/')
				buff.WriteByte(' ')
			}

			buff.WriteString(v)
		}
		messageParameters = append(
			messageParameters,
			messageParameter{
				Key:   "Values",
				Value: buff.String(),
			},
		)
	case "requiredIf", "requiredUnless", "same":
		other := getDisplayableAttribute(o, validTag.params[0])
		messageParameters = append(
			messageParameters,
			messageParameter{
				Key:   "Other",
				Value: other,
			},
		)
	}

	return messageParameters
}

func (v *Validator) formatsMessages(fieldError *FieldError) *FieldError {
	var message string
	var ok bool

	if message, ok = v.CustomMessage[fieldError.StructName+"."+fieldError.MessageName]; ok {
		fieldError.SetMessage(message)
		return fieldError
	}

	message, ok = MessageMap[fieldError.MessageName]
	if ok {
		attribute := fieldError.Attribute
		if customAttribute, ok := v.Attributes[fieldError.StructName]; ok {
			attribute = customAttribute
		} else if fieldError.DefaultAttribute != "" {
			attribute = fieldError.DefaultAttribute
		}
		message = replaceAttributes(message, attribute, fieldError.MessageParameters)

		fieldError.SetMessage(message)
		return fieldError
	}

	fieldError.SetMessage(fmt.Sprintf("validator: undefined message : %s", fieldError.MessageName))
	return fieldError
}

func replaceAttributes(message, attribute string, messageParameters MessageParameters) string {
	message = strings.Replace(message, "{{.Attribute}}", attribute, -1)
	for _, parameter := range messageParameters {
		message = strings.Replace(message, "{{."+parameter.Key+"}}", parameter.Value, -1)
	}
	return message
}

func getDisplayableAttribute(o reflect.Value, attribute string) string {
	attributes := strings.Split(attribute, ".")
	return attributes[len(attributes)-1]
}

func findField(fieldName string, v reflect.Value) (reflect.Value, error) {
	if v.Kind() != reflect.Struct {
		return reflect.Value{}, fmt.Errorf("findField: value is not a struct, got %s", v.Kind())
	}
	fields := strings.Split(fieldName, ".")
	current := v.FieldByName(fields[0])
	i := 1
	if len(fields) > i {
		for {
			if current.Kind() == reflect.Interface || current.Kind() == reflect.Ptr {
				current = current.Elem()
			}

			if !current.IsValid() {
				return current, fmt.Errorf("validator: findField Struct is nil")
			}

			name := fields[i]
			current = current.FieldByName(name)
			if i == len(fields)-1 {
				break
			}
			i++
		}
	}

	return current, nil
}

func (v *Validator) checkDependentRulesWithStatus(validTag *ValidTag, f *field, value, o reflect.Value, name, structName string) (bool, error) {
	isValid := true
	var funcError error
	var anotherField reflect.Value
	var err error
	var handled bool

	switch validTag.name {
	case "gt", "gte", "lt", "lte":
		// Check if the parameter is numeric (parameter comparison) or a field name (field comparison)
		if len(validTag.params) > 0 {
			if _, err := ToFloat(validTag.params[0]); err == nil {
				// It's a numeric parameter, skip field lookup and let ParamRuleMap handle it
				return false, nil
			}
		}
		// It's a field name, proceed with field comparison
		anotherField, err = findField(validTag.params[0], o)
		if err != nil {
			return false, nil
		}
		handled = true
	case "same", "different":
		if len(validTag.params) == 0 {
			return false, nil
		}
		anotherField, err = findField(validTag.params[0], o)
		if err != nil {
			return false, nil
		}
		handled = true
	case "confirmed":
		// Look for {fieldname}_confirmation field
		confirmationFieldName := f.attribute + "_confirmation"
		anotherField, err = findField(confirmationFieldName, o)
		if err != nil || !anotherField.IsValid() {
			// Also try with CamelCase naming: {fieldname}Confirmation
			confirmationFieldName = f.attribute + "Confirmation"
			anotherField, err = findField(confirmationFieldName, o)
			if err != nil || !anotherField.IsValid() {
				return false, nil
			}
		}
		handled = true
	case "acceptedIf", "declinedIf", "prohibitedIf", "prohibitedUnless", "missingIf", "missingUnless", "presentIf", "presentUnless", "prohibitedIfAccepted", "prohibitedIfDeclined":
		if len(validTag.params) == 0 {
			return false, nil
		}
		anotherField, err = findField(validTag.params[0], o)
		if err != nil {
			return false, nil
		}
		handled = true
	case "prohibits", "missingWith", "missingWithAll", "presentWith", "presentWithAll":
		// These validators need object access but don't use a single anotherField
		handled = true
	}

	switch validTag.name {
	case "gt":
		// Only handle field comparison, parameter comparison is handled by ParamRuleMap
		if anotherField.IsValid() {
			isValid, funcError = validateGt(value, anotherField)
		} else {
			return false, nil // Let ParamRuleMap handle it
		}
	case "gte":
		if anotherField.IsValid() {
			isValid, funcError = validateGte(value, anotherField)
		} else {
			return false, nil
		}
	case "lt":
		if anotherField.IsValid() {
			isValid, funcError = validateLt(value, anotherField)
		} else {
			return false, nil
		}
	case "lte":
		if anotherField.IsValid() {
			isValid, funcError = validateLte(value, anotherField)
		} else {
			return false, nil
		}
	case "same":
		isValid, funcError = validateSame(value, anotherField)
	case "different":
		isValid, funcError = validateDifferent(value, anotherField)
	case "confirmed":
		isValid, funcError = validateConfirmed(value, anotherField)
	case "acceptedIf":
		if len(validTag.params) >= 2 {
			isValid, funcError = validateAcceptedIf(value, anotherField, validTag.params[1:])
		}
	case "declinedIf":
		if len(validTag.params) >= 2 {
			isValid, funcError = validateDeclinedIf(value, anotherField, validTag.params[1:])
		}
	case "prohibitedIf":
		if len(validTag.params) >= 2 {
			isValid, funcError = validateProhibitedIf(value, anotherField, validTag.params[1:])
		}
	case "prohibitedUnless":
		if len(validTag.params) >= 2 {
			isValid, funcError = validateProhibitedUnless(value, anotherField, validTag.params[1:])
		}
	case "missingIf":
		if len(validTag.params) >= 2 {
			isValid, funcError = validateMissingIf(value, anotherField, validTag.params[1:])
		}
	case "missingUnless":
		if len(validTag.params) >= 2 {
			isValid, funcError = validateMissingUnless(value, anotherField, validTag.params[1:])
		}
	case "presentIf":
		if len(validTag.params) >= 2 {
			isValid, funcError = validatePresentIf(value, anotherField, validTag.params[1:])
		}
	case "presentUnless":
		if len(validTag.params) >= 2 {
			isValid, funcError = validatePresentUnless(value, anotherField, validTag.params[1:])
		}
	case "prohibitedIfAccepted":
		isValid, funcError = validateProhibitedIfAccepted(value, anotherField)
	case "prohibitedIfDeclined":
		isValid, funcError = validateProhibitedIfDeclined(value, anotherField)
	case "prohibits":
		isValid = validateProhibits(value, validTag.params, o)
	case "missingWith":
		isValid = validateMissingWith(value, validTag.params, o)
	case "missingWithAll":
		isValid = validateMissingWithAll(value, validTag.params, o)
	case "presentWith":
		isValid = validatePresentWith(value, validTag.params, o)
	case "presentWithAll":
		isValid = validatePresentWithAll(value, validTag.params, o)
	}

	if !isValid {
		return handled, v.formatsMessages(&FieldError{
			Name:              name,
			StructName:        structName,
			Tag:               validTag.name,
			MessageName:       validTag.messageName,
			MessageParameters: parseValidatorMessageParameters(validTag, o),
			Attribute:         f.attribute,
			DefaultAttribute:  f.defaultAttribute,
			Value:             ToString(value.Interface()),
			FuncError:         funcError,
		})
	}

	return handled, nil
}

func (v *Validator) checkDependentRules(validTag *ValidTag, f *field, value, o reflect.Value, name, structName string) error {
	_, err := v.checkDependentRulesWithStatus(validTag, f, value, o, name, structName)
	return err
}
