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

// Reusable zero-size map value for set-style reflect maps (e.g. isDistinct).
var (
	emptyStructType  = reflect.TypeOf(struct{}{})
	emptyStructValue = reflect.ValueOf(struct{}{})
)

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

// appendNamespace returns a fresh namespace of the form base + part + ".".
// It always allocates a new backing array so sibling fields can never corrupt
// each other's paths through append-aliasing of a shared base slice.
func appendNamespace(base, part []byte) []byte {
	out := make([]byte, 0, len(base)+len(part)+1)
	out = append(out, base...)
	out = append(out, part...)
	out = append(out, '.')
	return out
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

// =============================================================================
// Helper Functions - Reduce code duplication and improve performance
// =============================================================================

// deref dereferences interface and pointer values to their underlying value.
// This consolidates the repeated pattern found 30+ times in the codebase.
func deref(v reflect.Value) reflect.Value {
	for v.Kind() == reflect.Interface || v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return v
		}
		v = v.Elem()
	}
	return v
}

// comparisonOp represents a comparison operator
type comparisonOp string

const (
	opGt  comparisonOp = ">"
	opGte comparisonOp = ">="
	opLt  comparisonOp = "<"
	opLte comparisonOp = "<="
	opEq  comparisonOp = "=="
)

// compareValue is a unified comparison function that handles all numeric types.
// This replaces the duplicated switch statements in isMin, isMax, isGt, isGte, isLt, isLte.
func compareValue(v reflect.Value, params []string, op comparisonOp, ruleName string) (bool, error) {
	if len(params) == 0 {
		return false, fmt.Errorf("validator: %s rule requires at least one parameter", ruleName)
	}

	// Check for decimal.Decimal type first
	if d, ok := asDecimal(v); ok {
		p, _, err := parseDecimalParams(params)
		if err != nil {
			return false, fmt.Errorf("validator: %s decimal: %w", ruleName, err)
		}
		switch op {
		case opGt:
			return IsDecimalGt(d, p), nil
		case opGte:
			return IsDecimalGte(d, p), nil
		case opLt:
			return IsDecimalLt(d, p), nil
		case opLte:
			return IsDecimalLte(d, p), nil
		case opEq:
			return IsDecimalEqual(d, p), nil
		}
	}

	switch v.Kind() {
	case reflect.String:
		p, err := ToInt(params[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for %s rule on string field: %w", ruleName, err)
		}
		return compareString(v.String(), p, string(op))

	case reflect.Slice, reflect.Map, reflect.Array:
		p, err := ToInt(params[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for %s rule on collection field: %w", ruleName, err)
		}
		return compareInt64(int64(v.Len()), p, string(op))

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p, err := ToInt(params[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for %s rule on int field: %w", ruleName, err)
		}
		return compareInt64(v.Int(), p, string(op))

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p, err := ToUint(params[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for %s rule on uint field: %w", ruleName, err)
		}
		return compareUint64(v.Uint(), p, string(op))

	case reflect.Float32, reflect.Float64:
		p, err := ToFloat(params[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for %s rule on float field: %w", ruleName, err)
		}
		return compareFloat64(v.Float(), p, string(op))

	default:
		return false, fmt.Errorf("validator: %s rule is not supported for type %s", ruleName, v.Kind())
	}
}

// compareFields compares two reflect.Value fields with the given operator.
// This consolidates isSame, isLt, isLte, isGt, isGte field comparison functions.
func compareFields(v, anotherField reflect.Value, op comparisonOp, ruleName string) (bool, error) {
	if !v.IsValid() || !anotherField.IsValid() {
		return false, fmt.Errorf("validator: %s invalid reflection values", ruleName)
	}
	if v.Kind() != anotherField.Kind() {
		return false, fmt.Errorf("validator: %s The two fields must be of the same type %T, %T", ruleName, v.Interface(), anotherField.Interface())
	}

	// Check for decimal.Decimal type first
	if d1, ok := asDecimal(v); ok {
		if d2, ok := asDecimal(anotherField); ok {
			switch op {
			case opGt:
				return IsDecimalGt(d1, d2), nil
			case opGte:
				return IsDecimalGte(d1, d2), nil
			case opLt:
				return IsDecimalLt(d1, d2), nil
			case opLte:
				return IsDecimalLte(d1, d2), nil
			case opEq:
				return IsDecimalEqual(d1, d2), nil
			}
		}
	}

	switch v.Kind() {
	case reflect.String:
		if op == opEq {
			return v.String() == anotherField.String(), nil
		}
		return compareString(v.String(), int64(utf8.RuneCountInString(anotherField.String())), string(op))
	case reflect.Slice, reflect.Map, reflect.Array:
		return compareInt64(int64(v.Len()), int64(anotherField.Len()), string(op))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return compareInt64(v.Int(), anotherField.Int(), string(op))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return compareUint64(v.Uint(), anotherField.Uint(), string(op))
	case reflect.Float32, reflect.Float64:
		return compareFloat64(v.Float(), anotherField.Float(), string(op))
	default:
		return false, fmt.Errorf("validator: %s unsupported type %T", ruleName, v.Interface())
	}
}

const tagName string = "valid"

// Validator construct
type Validator struct {
	Attributes    map[string]string
	CustomMessage map[string]string
	Translator    *Translator
	// FailFast stops validation at the first field that fails and returns
	// immediately, instead of collecting every error. The default (false)
	// preserves the collect-all behavior. Set it once at setup, before
	// concurrent ValidateStruct calls.
	FailFast bool
}

// Default returns a instance of Validator
var Default = New()

// New returns a new instance of Validator
func New() *Validator {
	return &Validator{}
}

// isBetween check The field under validation must have a size between the given min and max. Strings, numerics, arrays, and files are evaluated in the same fashion as the size rule.
//
//nolint:gocyclo,gocritic // Complex validation logic with parameter names
func isBetween(v reflect.Value, params []string) (bool, error) {
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

// isWithRuleMap validates a value using RuleMap and returns formatted error if validation fails
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

// isWithParamRuleMap validates a value using ParamRuleMap and returns formatted error if validation fails
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

// isWithStringRulesMap validates a string value using StringRulesMap and returns formatted error if validation fails
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

// isWithStringParamRulesMap validates a string value using StringParamRulesMap and returns formatted error if validation fails
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

// isFieldComparisonRule reports whether a rule name is a comparison rule that may
// have already been satisfied by dependent field comparison (gt/gte/lt/lte).
func isFieldComparisonRule(name string) bool {
	return name == "gt" || name == "gte" || name == "lt" || name == "lte"
}

// applyRuleMaps runs RuleMap then ParamRuleMap for a single tag. ParamRuleMap is
// skipped for comparison rules already handled by dependent field comparison.
func (v *Validator) applyRuleMaps(tag *ValidTag, value reflect.Value, f *field, name, structName string, o reflect.Value, handled bool) error {
	if err := v.validateWithRuleMap(tag, value, f, name, structName, o); err != nil {
		return err
	}
	if handled && isFieldComparisonRule(tag.name) {
		return nil
	}
	return v.validateWithParamRuleMap(tag, value, f, name, structName, o)
}

// validateCollectionRules applies dependent rules and RuleMap/ParamRuleMap to a
// map or slice value (without string-specific rules).
func (v *Validator) validateCollectionRules(f *field, value reflect.Value, name, structName string, o reflect.Value) error {
	for _, tag := range f.validTags {
		handled, err := v.checkDependentRulesWithStatus(tag, f, value, o, name, structName)
		if err != nil {
			return err
		}
		if err := v.applyRuleMaps(tag, value, f, name, structName, o, handled); err != nil {
			return err
		}
	}
	return nil
}

// isCommonRules applies common validation rules (RuleMap, ParamRuleMap, dependent rules)
func (v *Validator) validateCommonRules(tags otherValidTags, value reflect.Value, f *field, name, structName string, o reflect.Value) error {
	for _, tag := range tags {
		handled, err := v.checkDependentRulesWithStatus(tag, f, value, o, name, structName)
		if err != nil {
			return err
		}

		if err := v.applyRuleMaps(tag, value, f, name, structName, o, handled); err != nil {
			return err
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
// checkRequiredIfCondition reports whether the field is required (invalid) and,
// when it is, the matched value that triggered the requirement. It must not
// mutate any cached tag — the matched value is returned to the caller so the
// "Value" message parameter can be built per call.
func checkRequiredIfCondition(v reflect.Value, values, params []string) (valid bool, matchedValue string, err error) {
	for _, value := range values {
		if InString(value, params) && Empty(v) {
			return false, value, nil
		}
	}
	return true, "", nil
}

// isCustomTypeRules validates using CustomTypeRuleMap
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

// isMapFields validates map structure and each element
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
			key := []byte(k.String())
			newJSONNamespace := appendNamespace(appendNamespace(jsonNamespace, f.nameBytes), key)
			newstructNamespace := appendNamespace(appendNamespace(structNamespace, f.structNameBytes), key)
			err = v.ValidateStruct(item.Interface(), newJSONNamespace, newstructNamespace)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// isSliceFields validates slice/array structure and each element
func (v *Validator) validateSliceFields(value reflect.Value, f *field, jsonNamespace, structNamespace []byte) error {
	for i := 0; i < value.Len(); i++ {
		var err error
		item := value.Index(i)
		if item.Kind() == reflect.Interface {
			item = item.Elem()
		}

		if item.Kind() == reflect.Struct || item.Kind() == reflect.Ptr {
			index := []byte(strconv.Itoa(i))
			newJSONNamespace := appendNamespace(appendNamespace(jsonNamespace, f.nameBytes), index)
			newStructNamespace := appendNamespace(appendNamespace(structNamespace, f.structNameBytes), index)
			err = v.ValidateStruct(value.Index(i).Interface(), newJSONNamespace, newStructNamespace)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// IsBetween check The field under validation must have a size between the given min and max. Strings, numerics, arrays, and files are evaluated in the same fashion as the size rule.
func IsBetween(i interface{}, params []string) (bool, error) {
	v := reflect.ValueOf(i)
	return isBetween(v, params)
}

// isDigitsBetween check The field under validation must have a length between the given min and max.
func isDigitsBetween(v reflect.Value, params []string) (bool, error) {
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

// IsDigitsBetween check The field under validation must have a length between the given min and max.
func IsDigitsBetween(i interface{}, params []string) (bool, error) {
	v := reflect.ValueOf(i)
	return isDigitsBetween(v, params)
}

// isSize The field under validation must have a size matching the given value.
// For string data, value corresponds to the number of characters.
// For numeric data, value corresponds to a given integer value.
// For an array | map | slice, size corresponds to the count of the array | map | slice.
func isSize(v reflect.Value, param []string) (bool, error) {
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

// IsSize The field under validation must have a size matching the given value.
// For string data, value corresponds to the number of characters.
// For numeric data, value corresponds to a given integer value.
// For an array | map | slice, size corresponds to the count of the array | map | slice.
func IsSize(i interface{}, params []string) (bool, error) {
	v := reflect.ValueOf(i)
	return isSize(v, params)
}

// isMax is the validation function for validating if the current field's value is less than or equal to the param's value.
func isMax(v reflect.Value, params []string) (bool, error) {
	return compareValue(v, params, opLte, "Max")
}

// IsMax is the validation function for validating if the current field's value is less than or equal to the param's value.
func IsMax(i interface{}, params []string) (bool, error) {
	v := reflect.ValueOf(i)
	return isMax(v, params)
}

// isMin is the validation function for validating if the current field's value is greater than or equal to the param's value.
func isMin(v reflect.Value, params []string) (bool, error) {
	return compareValue(v, params, opGte, "Min")
}

// IsMin is the validation function for validating if the current field's value is greater than or equal to the param's value.
func IsMin(i interface{}, params []string) (bool, error) {
	v := reflect.ValueOf(i)
	return isMin(v, params)
}

// isGtParam is the validation function for validating if the current field's value is greater than the param's value.
func isGtParam(v reflect.Value, params []string) (bool, error) {
	return compareValue(v, params, opGt, "Gt")
}

// IsGtParam is the validation function for validating if the current field's value is greater than the param's value.
func IsGtParam(i interface{}, params []string) (bool, error) {
	v := reflect.ValueOf(i)
	return isGtParam(v, params)
}

// isGteParam is the validation function for validating if the current field's value is greater than or equal to the param's value.
func isGteParam(v reflect.Value, params []string) (bool, error) {
	return compareValue(v, params, opGte, "Gte")
}

// IsGteParam is the validation function for validating if the current field's value is greater than or equal to the param's value.
func IsGteParam(i interface{}, params []string) (bool, error) {
	v := reflect.ValueOf(i)
	return isGteParam(v, params)
}

// isLtParam is the validation function for validating if the current field's value is less than the param's value.
func isLtParam(v reflect.Value, params []string) (bool, error) {
	return compareValue(v, params, opLt, "Lt")
}

// IsLtParam is the validation function for validating if the current field's value is less than the param's value.
func IsLtParam(i interface{}, params []string) (bool, error) {
	v := reflect.ValueOf(i)
	return isLtParam(v, params)
}

// isLteParam is the validation function for validating if the current field's value is less than or equal to the param's value.
func isLteParam(v reflect.Value, params []string) (bool, error) {
	return compareValue(v, params, opLte, "Lte")
}

// IsLteParam is the validation function for validating if the current field's value is less than or equal to the param's value.
func IsLteParam(i interface{}, params []string) (bool, error) {
	v := reflect.ValueOf(i)
	return isLteParam(v, params)
}

// isSame is the validation function for validating if the current field's value equal the param's value.
func isSame(v, anotherField reflect.Value) (bool, error) {
	return compareFields(v, anotherField, opEq, "Same")
}

// IsSame is the validation function for validating if the current field's value is greater than or equal to the param's value.
func IsSame(i, a interface{}) (bool, error) {
	v := reflect.ValueOf(i)
	anotherField := reflect.ValueOf(a)
	return isSame(v, anotherField)
}

// isLt is the validation function for validating if the current field's value is less than the param's value.
func isLt(v, anotherField reflect.Value) (bool, error) {
	return compareFields(v, anotherField, opLt, "Lt")
}

// IsLt is the validation function for validating if the current field's value is less than the param's value.
func IsLt(i, a interface{}) (bool, error) {
	v := reflect.ValueOf(i)
	anotherField := reflect.ValueOf(a)
	return isLt(v, anotherField)
}

// isLte is the validation function for validating if the current field's value is less than or equal to the param's value.
func isLte(v, anotherField reflect.Value) (bool, error) {
	return compareFields(v, anotherField, opLte, "Lte")
}

// IsLte is the validation function for validating if the current field's value is less than or equal to the param's value.
func IsLte(i, a interface{}) (bool, error) {
	v := reflect.ValueOf(i)
	anotherField := reflect.ValueOf(a)
	return isLte(v, anotherField)
}

// isGt is the validation function for validating if the current field's value is greater than to the param's value.
func isGt(v, anotherField reflect.Value) (bool, error) {
	return compareFields(v, anotherField, opGt, "Gt")
}

// IsGt is the validation function for validating if the current field's value is greater than to the param's value.
func IsGt(i, a interface{}) (bool, error) {
	v := reflect.ValueOf(i)
	anotherField := reflect.ValueOf(a)
	return isGt(v, anotherField)
}

// isGte is the validation function for validating if the current field's value is greater than or equal to the param's value.
func isGte(v, anotherField reflect.Value) (bool, error) {
	return compareFields(v, anotherField, opGte, "Gte")
}

// IsGte is the validation function for validating if the current field's value is greater than to the param's value.
func IsGte(i, a interface{}) (bool, error) {
	v := reflect.ValueOf(i)
	anotherField := reflect.ValueOf(a)
	return isGte(v, anotherField)
}

// isDistinct is the validation function for validating an attribute is unique among other values.
func isDistinct(v reflect.Value) (bool, error) {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64:
		return true, nil
	case reflect.Slice, reflect.Array:
		// Only keys matter for uniqueness; use a zero-size value type so each
		// entry stores nothing instead of a full copy of the collection.
		seen := reflect.MakeMapWithSize(reflect.MapOf(v.Type().Elem(), emptyStructType), v.Len())
		for i := 0; i < v.Len(); i++ {
			seen.SetMapIndex(v.Index(i), emptyStructValue)
		}
		return v.Len() == seen.Len(), nil
	case reflect.Map:
		seen := reflect.MakeMapWithSize(reflect.MapOf(v.Type().Elem(), emptyStructType), v.Len())
		for _, k := range v.MapKeys() {
			seen.SetMapIndex(v.MapIndex(k), emptyStructValue)
		}
		return v.Len() == seen.Len(), nil
	}

	return false, fmt.Errorf("validator: Distinct unsupported type %T", v.Interface())
}

// IsDistinct is the validation function for validating an attribute is unique among other values.
func IsDistinct(i interface{}) bool {
	v := reflect.ValueOf(i)
	valid, _ := isDistinct(v)
	return valid
}

// isMultipleOf is the validation function for validating if the current field's value is a multiple of the param's value.
func isMultipleOf(v reflect.Value, params []string) (bool, error) {
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

// IsMultipleOf is the validation function for validating if a value is a multiple of another.
func IsMultipleOf(i interface{}, params []string) (bool, error) {
	v := reflect.ValueOf(i)
	return isMultipleOf(v, params)
}

// isMaxDigits is the validation function for validating the maximum number of digits.
func isMaxDigits(v reflect.Value, params []string) (bool, error) {
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

// IsMaxDigits is the validation function for validating the maximum number of digits.
func IsMaxDigits(i interface{}, params []string) (bool, error) {
	v := reflect.ValueOf(i)
	return isMaxDigits(v, params)
}

// isMinDigits is the validation function for validating the minimum number of digits.
func isMinDigits(v reflect.Value, params []string) (bool, error) {
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

// IsMinDigits is the validation function for validating the minimum number of digits.
func IsMinDigits(i interface{}, params []string) (bool, error) {
	v := reflect.ValueOf(i)
	return isMinDigits(v, params)
}

// isDecimalPrecision is the validation function for validating decimal places.
// params[0] is min decimal places, params[1] (optional) is max decimal places.
func isDecimalPrecision(v reflect.Value, params []string) (bool, error) {
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

// IsDecimalPrecision is the validation function for validating decimal places.
func IsDecimalPrecision(i interface{}, params []string) (bool, error) {
	v := reflect.ValueOf(i)
	return isDecimalPrecision(v, params)
}

// isContains is the validation function for validating a string/array/slice contains specified values.
func isContains(v reflect.Value, params []string) (bool, error) {
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

// IsContains is the validation function for validating an array/slice contains specified values.
func IsContains(i interface{}, params []string) (bool, error) {
	v := reflect.ValueOf(i)
	return isContains(v, params)
}

// isDoesntContain is the validation function for validating a string/array/slice does not contain specified values.
func isDoesntContain(v reflect.Value, params []string) (bool, error) {
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

// IsDoesntContain is the validation function for validating an array/slice does not contain specified values.
func IsDoesntContain(i interface{}, params []string) (bool, error) {
	v := reflect.ValueOf(i)
	return isDoesntContain(v, params)
}

// IsMimeTypes is the validation function for the file must match one of the given MIME types.
func IsMimeTypes(data []byte, mimeTypes []string) bool {
	mimeType := http.DetectContentType(data)
	for _, value := range mimeTypes {
		if mimeType == value {
			return true
		}
	}
	return false
}

// IsMimes is the validation function for the file must have a MIME type corresponding to one of the listed extensions.
func IsMimes(data []byte, mimes []string) (bool, error) {
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

// IsImage is the validation function for the The file under validation must be an image (jpeg, png, bmp, gif, or svg)
func IsImage(data []byte) bool {
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

	// errs is left nil: the common case is zero errors, and appending from nil
	// avoids allocating a backing array that a successful validation never uses.
	var errs Errors
	fields := cachedTypefields(val.Type())

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
			if v.FailFast {
				break
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

	// Check for custom type functions or auto-detect types with IsSet()/Value() methods
	// (e.g., graphql.Omittable[T], sql.NullString). Results are cached per type.
	if value.Kind() == reflect.Struct {
		if extract := resolveCustomTypeFunc(value.Type()); extract != nil {
			inner, shouldValidate := extract(value)
			if !shouldValidate || !inner.IsValid() {
				// Value was not set — only check required rules
				if err := v.checkRequired(value, f, o, name, structName); err != nil {
					return err
				}
				return nil
			}
			value = inner
		}
	}

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
		if err := v.validateCollectionRules(f, value, name, structName, o); err != nil {
			return err
		}
		return v.validateMapFields(value, f, jsonNamespace, structNamespace)
	case reflect.Slice, reflect.Array:
		// Validate slice/array-specific rules (without string-specific rules)
		if err := v.validateCollectionRules(f, value, name, structName, o); err != nil {
			return err
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
		jsonNamespace = appendNamespace(jsonNamespace, f.nameBytes)
		structNamespace = appendNamespace(structNamespace, f.structNameBytes)
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
	case reflect.Chan, reflect.Func:
		return v.IsNil()
	case reflect.Struct:
		// Use IsZero() for structs - much faster than DeepEqual
		// Available since Go 1.13
		return v.IsZero()
	case reflect.Complex64, reflect.Complex128:
		return v.Complex() == 0
	case reflect.Invalid:
		return true
	}

	// For any other type (should be rare), use IsZero
	return v.IsZero()
}

// Error returns string equivalent for reflect.Type
func (e *UnsupportedTypeError) Error() string {
	return "validator: unsupported type: " + e.Type.String()
}

func (sv stringValues) Len() int           { return len(sv) }
func (sv stringValues) Swap(i, j int)      { sv[i], sv[j] = sv[j], sv[i] }
func (sv stringValues) Less(i, j int) bool { return sv.get(i) < sv.get(j) }
func (sv stringValues) get(i int) string   { return sv[i].String() }

// isRequired check value required when anotherField str is a member of the set of strings params
func isRequired(v reflect.Value) bool {
	return !Empty(v)
}

// IsRequired check value required when anotherField str is a member of the set of strings params
func IsRequired(i interface{}) bool {
	v := reflect.ValueOf(i)
	return isRequired(v)
}

// acceptedValues are the values that are considered "accepted"
var acceptedValues = []string{"yes", "on", "1", "true"}

// declinedValues are the values that are considered "declined"
var declinedValues = []string{"no", "off", "0", "false"}

// isAccepted checks if the field is accepted (yes, on, 1, "1", true, "true")
func isAccepted(v reflect.Value) bool {
	if v.Kind() == reflect.Bool {
		return v.Bool()
	}
	value := strings.ToLower(ToString(v.Interface()))
	return InString(value, acceptedValues)
}

// IsAccepted checks if the field is accepted
func IsAccepted(i interface{}) bool {
	v := reflect.ValueOf(i)
	return isAccepted(v)
}

// isDeclined checks if the field is declined (no, off, 0, "0", false, "false")
func isDeclined(v reflect.Value) bool {
	if v.Kind() == reflect.Bool {
		return !v.Bool()
	}
	value := strings.ToLower(ToString(v.Interface()))
	return InString(value, declinedValues)
}

// IsDeclined checks if the field is declined
func IsDeclined(i interface{}) bool {
	v := reflect.ValueOf(i)
	return isDeclined(v)
}

// isProhibited checks if the field is empty (prohibited means it must be empty)
func isProhibited(v reflect.Value) bool {
	return Empty(v)
}

// IsProhibited checks if the field is prohibited (must be empty)
func IsProhibited(i interface{}) bool {
	v := reflect.ValueOf(i)
	return isProhibited(v)
}

// isAcceptedIf checks if the field is accepted when another field equals a specific value
func isAcceptedIf(v, anotherField reflect.Value, params []string) (bool, error) {
	if anotherField.Kind() == reflect.Interface || anotherField.Kind() == reflect.Ptr {
		anotherField = anotherField.Elem()
	}

	if !anotherField.IsValid() {
		return true, nil
	}

	value := ToString(anotherField.Interface())
	if InString(value, params) {
		return isAccepted(v), nil
	}
	return true, nil
}

// isDeclinedIf checks if the field is declined when another field equals a specific value
func isDeclinedIf(v, anotherField reflect.Value, params []string) (bool, error) {
	if anotherField.Kind() == reflect.Interface || anotherField.Kind() == reflect.Ptr {
		anotherField = anotherField.Elem()
	}

	if !anotherField.IsValid() {
		return true, nil
	}

	value := ToString(anotherField.Interface())
	if InString(value, params) {
		return isDeclined(v), nil
	}
	return true, nil
}

// isProhibitedIf checks if the field is empty when another field equals a specific value
func isProhibitedIf(v, anotherField reflect.Value, params []string) (bool, error) {
	if anotherField.Kind() == reflect.Interface || anotherField.Kind() == reflect.Ptr {
		anotherField = anotherField.Elem()
	}

	if !anotherField.IsValid() {
		return true, nil
	}

	value := ToString(anotherField.Interface())
	if InString(value, params) {
		return isProhibited(v), nil
	}
	return true, nil
}

// isProhibitedUnless checks if the field is empty unless another field equals a specific value
func isProhibitedUnless(v, anotherField reflect.Value, params []string) (bool, error) {
	if anotherField.Kind() == reflect.Interface || anotherField.Kind() == reflect.Ptr {
		anotherField = anotherField.Elem()
	}

	if !anotherField.IsValid() {
		return true, nil
	}

	value := ToString(anotherField.Interface())
	if !InString(value, params) {
		return isProhibited(v), nil
	}
	return true, nil
}

// isMissing checks if the field is not present (must be absent from the data)
func isMissing(v reflect.Value) bool {
	return !v.IsValid() || Empty(v)
}

// IsMissing checks if the field is missing
func IsMissing(i interface{}) bool {
	v := reflect.ValueOf(i)
	return isMissing(v)
}

// isMissingIf checks if the field is missing when another field equals a specific value
func isMissingIf(v, anotherField reflect.Value, params []string) (bool, error) {
	if anotherField.Kind() == reflect.Interface || anotherField.Kind() == reflect.Ptr {
		anotherField = anotherField.Elem()
	}

	if !anotherField.IsValid() {
		return true, nil
	}

	value := ToString(anotherField.Interface())
	if InString(value, params) {
		return isMissing(v), nil
	}
	return true, nil
}

// isMissingUnless checks if the field is missing unless another field equals a specific value
func isMissingUnless(v, anotherField reflect.Value, params []string) (bool, error) {
	if anotherField.Kind() == reflect.Interface || anotherField.Kind() == reflect.Ptr {
		anotherField = anotherField.Elem()
	}

	if !anotherField.IsValid() {
		return true, nil
	}

	value := ToString(anotherField.Interface())
	if !InString(value, params) {
		return isMissing(v), nil
	}
	return true, nil
}

// isRequiredIfAccepted checks if the field is required when another field is accepted
func isRequiredIfAccepted(v, anotherField reflect.Value) (bool, error) {
	if isAccepted(anotherField) {
		return isRequired(v), nil
	}
	return true, nil
}

// isRequiredIfDeclined checks if the field is required when another field is declined
func isRequiredIfDeclined(v, anotherField reflect.Value) (bool, error) {
	if isDeclined(anotherField) {
		return isRequired(v), nil
	}
	return true, nil
}

// isProhibitedIfAccepted checks if the field is prohibited when another field is accepted
func isProhibitedIfAccepted(v, anotherField reflect.Value) (bool, error) {
	if isAccepted(anotherField) {
		return isProhibited(v), nil
	}
	return true, nil
}

// isProhibitedIfDeclined checks if the field is prohibited when another field is declined
func isProhibitedIfDeclined(v, anotherField reflect.Value) (bool, error) {
	if isDeclined(anotherField) {
		return isProhibited(v), nil
	}
	return true, nil
}

// isProhibits checks if when this field has a value, the other specified fields must be empty
func isProhibits(v reflect.Value, otherFields []string, obj reflect.Value) bool {
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

// isMissingWith checks if the field is missing when any of the specified fields are present
func isMissingWith(v reflect.Value, otherFields []string, obj reflect.Value) bool {
	for _, fieldName := range otherFields {
		otherField, err := findField(fieldName, obj)
		if err != nil {
			continue
		}
		if !Empty(otherField) {
			return isMissing(v)
		}
	}
	return true
}

// isMissingWithAll checks if the field is missing when all of the specified fields are present
func isMissingWithAll(v reflect.Value, otherFields []string, obj reflect.Value) bool {
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
		return isMissing(v)
	}
	return true
}

// isPresent checks if the field is present (for struct fields, always true since they exist)
func isPresent(v reflect.Value) bool {
	// In a struct context, all fields are always "present" - they exist in the struct
	// This is different from "required" which checks for non-empty values
	return v.IsValid()
}

// IsPresent checks if the field is present
func IsPresent(i interface{}) bool {
	v := reflect.ValueOf(i)
	return isPresent(v)
}

// isPresentIf checks if the field is present (non-empty) when another field equals a specific value
func isPresentIf(v, anotherField reflect.Value, params []string) (bool, error) {
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

// isPresentUnless checks if the field is present (non-empty) unless another field equals a specific value
func isPresentUnless(v, anotherField reflect.Value, params []string) (bool, error) {
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

// isPresentWith checks if the field is present (non-empty) when any of the specified fields are present
func isPresentWith(v reflect.Value, otherFields []string, obj reflect.Value) bool {
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

// isPresentWithAll checks if the field is present (non-empty) when all of the specified fields are present
func isPresentWithAll(v reflect.Value, otherFields []string, obj reflect.Value) bool {
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

// isRequiredArrayKeys checks if the array has all the specified keys
func isRequiredArrayKeys(v reflect.Value, keys []string) (bool, error) {
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

// IsRequiredArrayKeys checks if the map has all the specified keys
func IsRequiredArrayKeys(i interface{}, keys []string) (bool, error) {
	v := reflect.ValueOf(i)
	return isRequiredArrayKeys(v, keys)
}

// isList checks if the slice/array is a valid list (sequential integer keys starting from 0)
// In Go, slices are always lists, so this just checks if it's a slice/array
func isList(v reflect.Value) (bool, error) {
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		return true, nil
	default:
		return false, nil
	}
}

// IsList checks if the value is a list (slice or array)
func IsList(i interface{}) (bool, error) {
	v := reflect.ValueOf(i)
	return isList(v)
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
	v = deref(v)
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

// isDate checks if the value is a valid date
func isDate(v reflect.Value) (bool, error) {
	v = deref(v)
	if !v.IsValid() || Empty(v) {
		return true, nil
	}

	_, ok := getTimeValue(v)
	return ok, nil
}

// IsDate checks if the value is a valid date
func IsDate(i interface{}) (bool, error) {
	v := reflect.ValueOf(i)
	return isDate(v)
}

// isDateFormat checks if the value matches the given date format
func isDateFormat(v reflect.Value, params []string) (bool, error) {
	if len(params) == 0 {
		return false, fmt.Errorf("validator: dateFormat rule requires a format parameter")
	}

	v = deref(v)
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

// IsDateFormat checks if the value matches the given date format
func IsDateFormat(i interface{}, params []string) (bool, error) {
	v := reflect.ValueOf(i)
	return isDateFormat(v, params)
}

// isAfter checks if the date is after the given date
func isAfter(v reflect.Value, params []string) (bool, error) {
	if len(params) == 0 {
		return false, fmt.Errorf("validator: after rule requires a date parameter")
	}

	v = deref(v)
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

// IsAfter checks if the date is after the given date
func IsAfter(i interface{}, params []string) (bool, error) {
	v := reflect.ValueOf(i)
	return isAfter(v, params)
}

// isAfterOrEqual checks if the date is after or equal to the given date
func isAfterOrEqual(v reflect.Value, params []string) (bool, error) {
	if len(params) == 0 {
		return false, fmt.Errorf("validator: afterOrEqual rule requires a date parameter")
	}

	v = deref(v)
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

// IsAfterOrEqual checks if the date is after or equal to the given date
func IsAfterOrEqual(i interface{}, params []string) (bool, error) {
	v := reflect.ValueOf(i)
	return isAfterOrEqual(v, params)
}

// isBefore checks if the date is before the given date
func isBefore(v reflect.Value, params []string) (bool, error) {
	if len(params) == 0 {
		return false, fmt.Errorf("validator: before rule requires a date parameter")
	}

	v = deref(v)
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

// IsBefore checks if the date is before the given date
func IsBefore(i interface{}, params []string) (bool, error) {
	v := reflect.ValueOf(i)
	return isBefore(v, params)
}

// isBeforeOrEqual checks if the date is before or equal to the given date
func isBeforeOrEqual(v reflect.Value, params []string) (bool, error) {
	if len(params) == 0 {
		return false, fmt.Errorf("validator: beforeOrEqual rule requires a date parameter")
	}

	v = deref(v)
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

// IsBeforeOrEqual checks if the date is before or equal to the given date
func IsBeforeOrEqual(i interface{}, params []string) (bool, error) {
	v := reflect.ValueOf(i)
	return isBeforeOrEqual(v, params)
}

// isIn checks if the value is in the given list
func isIn(v reflect.Value, params []string) (bool, error) {
	if len(params) == 0 {
		return false, fmt.Errorf("validator: in rule requires at least one parameter")
	}

	v = deref(v)
	if !v.IsValid() || Empty(v) {
		return true, nil
	}

	value := ToString(v.Interface())
	return InString(value, params), nil
}

// IsIn checks if the value is in the given list
func IsIn(i interface{}, params []string) (bool, error) {
	v := reflect.ValueOf(i)
	return isIn(v, params)
}

// isNotIn checks if the value is not in the given list
func isNotIn(v reflect.Value, params []string) (bool, error) {
	if len(params) == 0 {
		return false, fmt.Errorf("validator: notIn rule requires at least one parameter")
	}

	v = deref(v)
	if !v.IsValid() || Empty(v) {
		return true, nil
	}

	value := ToString(v.Interface())
	return !InString(value, params), nil
}

// IsNotIn checks if the value is not in the given list
func IsNotIn(i interface{}, params []string) (bool, error) {
	v := reflect.ValueOf(i)
	return isNotIn(v, params)
}

// isDifferent checks if the value is different from another field
func isDifferent(v, anotherField reflect.Value) (bool, error) {
	v = deref(v)
	if anotherField.Kind() == reflect.Interface || anotherField.Kind() == reflect.Ptr {
		anotherField = anotherField.Elem()
	}

	if !v.IsValid() || !anotherField.IsValid() {
		return true, nil
	}

	return ToString(v.Interface()) != ToString(anotherField.Interface()), nil
}

// IsDifferent checks if the value is different from another value
func IsDifferent(a, b interface{}) (bool, error) {
	return isDifferent(reflect.ValueOf(a), reflect.ValueOf(b))
}

// isConfirmed checks if the field has a matching confirmation field
// This is handled specially in the validation loop since it needs to find {field}_confirmation
func isConfirmed(v, confirmationField reflect.Value) (bool, error) {
	v = deref(v)
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

// isJSON checks if the value is a valid JSON string
func isJSON(v reflect.Value) (bool, error) {
	v = deref(v)
	if !v.IsValid() || Empty(v) {
		return true, nil
	}

	if v.Kind() != reflect.String {
		return false, nil
	}

	var js interface{}
	return json.Unmarshal([]byte(v.String()), &js) == nil, nil
}

// IsJSON checks if the value is a valid JSON string
func IsJSON(i interface{}) (bool, error) {
	v := reflect.ValueOf(i)
	return isJSON(v)
}

// isRegex checks if the value matches the given regex pattern
func isRegex(v reflect.Value, params []string) (bool, error) {
	if len(params) == 0 {
		return false, fmt.Errorf("validator: regex rule requires a pattern parameter")
	}

	v = deref(v)
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

// IsRegex checks if the value matches the given regex pattern
func IsRegex(i interface{}, params []string) (bool, error) {
	v := reflect.ValueOf(i)
	return isRegex(v, params)
}

// isNotRegex checks if the value does not match the given regex pattern
func isNotRegex(v reflect.Value, params []string) (bool, error) {
	if len(params) == 0 {
		return false, fmt.Errorf("validator: notRegex rule requires a pattern parameter")
	}

	v = deref(v)
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

// IsNotRegex checks if the value does not match the given regex pattern
func IsNotRegex(i interface{}, params []string) (bool, error) {
	v := reflect.ValueOf(i)
	return isNotRegex(v, params)
}

// isBoolean checks if the value is boolean-like (true/false, 1/0, "yes"/"no", "on"/"off")
func isBoolean(v reflect.Value) (bool, error) {
	v = deref(v)
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

// IsBoolean checks if the value is boolean-like
func IsBoolean(i interface{}) (bool, error) {
	v := reflect.ValueOf(i)
	return isBoolean(v)
}

// isString checks if the value is a string
func isString(v reflect.Value) (bool, error) {
	v = deref(v)
	if !v.IsValid() {
		return true, nil
	}

	return v.Kind() == reflect.String, nil
}

// IsString checks if the value is a string
func IsString(i interface{}) (bool, error) {
	v := reflect.ValueOf(i)
	return isString(v)
}

// isArray checks if the value is an array or slice
func isArray(v reflect.Value) (bool, error) {
	v = deref(v)
	if !v.IsValid() {
		return true, nil
	}

	return v.Kind() == reflect.Slice || v.Kind() == reflect.Array, nil
}

// IsArray checks if the value is an array or slice
func IsArray(i interface{}) (bool, error) {
	v := reflect.ValueOf(i)
	return isArray(v)
}

// isFilled checks if the field is not empty when it is present
func isFilled(v reflect.Value) (bool, error) {
	v = deref(v)
	// If the field is not valid (doesn't exist), it passes
	if !v.IsValid() {
		return true, nil
	}
	// If the field exists, it must not be empty
	return !Empty(v), nil
}

// IsFilled checks if the field is not empty when present
func IsFilled(i interface{}) (bool, error) {
	v := reflect.ValueOf(i)
	return isFilled(v)
}

// isRequiredIf check value required when anotherField str is a member of the set of strings params
func isRequiredIf(v, anotherField reflect.Value, params []string) (valid bool, matchedValue string, err error) {
	if anotherField.Kind() == reflect.Interface || anotherField.Kind() == reflect.Ptr {
		anotherField = anotherField.Elem()
	}

	if !anotherField.IsValid() {
		return true, "", nil
	}

	switch anotherField.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.String:
		value := ToString(anotherField)
		if InString(value, params) && Empty(v) {
			return false, value, nil
		}
	case reflect.Map, reflect.Slice, reflect.Array:
		values, err := extractValuesFromCollection(anotherField)
		if err != nil {
			return false, "", err
		}
		return checkRequiredIfCondition(v, values, params)
	default:
		return false, "", fmt.Errorf("validator: RequiredIf unsupported type %T", anotherField.Interface())
	}

	return true, "", nil
}

// isRequiredUnless check value required when anotherField str is a member of the set of strings params
func isRequiredUnless(v, anotherField reflect.Value, params []string) (bool, error) {
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
		var requiredIfValue string
		switch tag.name {
		case "required":
			isError = !isRequired(value)
		case "requiredIf":
			if len(tag.params) == 0 {
				continue
			}
			anotherField, err := findField(tag.params[0], o)
			if err == nil && len(tag.params) >= 2 {
				isValid, requiredIfValue, funcError = isRequiredIf(value, anotherField, tag.params[1:])
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
				isValid, funcError = isRequiredUnless(value, anotherField, tag.params[1:])
				if !isValid {
					isError = true
				}
			}
		case "requiredWith":
			if !isRequiredWith(tag.params, value, o) {
				isError = true
			}
		case "requiredWithAll":
			if !isRequiredWithAll(tag.params, value, o) {
				isError = true
			}
		case "requiredWithout":
			if !isRequiredWithout(tag.params, value, o) {
				isError = true
			}
		case "requiredWithoutAll":
			if !isRequiredWithoutAll(tag.params, value, o) {
				isError = true
			}
		case "requiredIfAccepted":
			if len(tag.params) == 0 {
				continue
			}
			anotherField, err := findField(tag.params[0], o)
			if err == nil {
				isValid, funcError = isRequiredIfAccepted(value, anotherField)
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
				isValid, funcError = isRequiredIfDeclined(value, anotherField)
				if !isValid {
					isError = true
				}
			}
		}

		if isError {
			messageParameters := parseValidatorMessageParameters(tag, o)
			if tag.name == "requiredIf" {
				// "Value" is the runtime value of the referenced field that
				// triggered the requirement; built per call, never cached.
				messageParameters = append(messageParameters, messageParameter{
					Key:   "Value",
					Value: requiredIfValue,
				})
			}
			return v.formatsMessages(&FieldError{
				Name:              name,
				StructName:        structName,
				Tag:               tag.name,
				MessageName:       tag.messageName,
				MessageParameters: messageParameters,
				Attribute:         f.attribute,
				DefaultAttribute:  f.defaultAttribute,
				Value:             ToString(value.Interface()),
				FuncError:         funcError,
			})
		}
	}

	return nil
}

// isRequiredWith The field under validation must be present and not empty only if any of the other specified fields are present.
func isRequiredWith(otherFields []string, currentField, obj reflect.Value) bool {
	if !allFailingRequired(otherFields, obj) {
		return isRequired(currentField)
	}
	return true
}

// isRequiredWithAll The field under validation must be present and not empty only if all of the other specified fields are present.
func isRequiredWithAll(otherFields []string, currentField, obj reflect.Value) bool {
	if !anyFailingRequired(otherFields, obj) {
		return isRequired(currentField)
	}
	return true
}

// RequiredWithout The field under validation must be present and not empty only when any of the other specified fields are not present.
func isRequiredWithout(otherFields []string, currentField, obj reflect.Value) bool {
	if anyFailingRequired(otherFields, obj) {
		return isRequired(currentField)
	}
	return true
}

// isRequiredWithoutAll The field under validation must be present and not empty only when all of the other specified fields are not present.
func isRequiredWithoutAll(otherFields []string, currentField, obj reflect.Value) bool {
	if allFailingRequired(otherFields, obj) {
		return isRequired(currentField)
	}
	return true
}

func parseValidatorMessageParameters(validTag *ValidTag, o reflect.Value) MessageParameters {
	// Copy the cached base params; appending below must never mutate or alias
	// the shared (cached) tag, since this runs on concurrent validations.
	messageParameters := append(MessageParameters(nil), validTag.messageParameters...)
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
		other := getDisplayableAttribute(validTag.params[0])
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
	message = strings.ReplaceAll(message, "{{.Attribute}}", attribute)
	for _, parameter := range messageParameters {
		message = strings.ReplaceAll(message, "{{."+parameter.Key+"}}", parameter.Value)
	}
	return message
}

func getDisplayableAttribute(attribute string) string {
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
			isValid, funcError = isGt(value, anotherField)
		} else {
			return false, nil // Let ParamRuleMap handle it
		}
	case "gte":
		if anotherField.IsValid() {
			isValid, funcError = isGte(value, anotherField)
		} else {
			return false, nil
		}
	case "lt":
		if anotherField.IsValid() {
			isValid, funcError = isLt(value, anotherField)
		} else {
			return false, nil
		}
	case "lte":
		if anotherField.IsValid() {
			isValid, funcError = isLte(value, anotherField)
		} else {
			return false, nil
		}
	case "same":
		isValid, funcError = isSame(value, anotherField)
	case "different":
		isValid, funcError = isDifferent(value, anotherField)
	case "confirmed":
		isValid, funcError = isConfirmed(value, anotherField)
	case "acceptedIf":
		if len(validTag.params) >= 2 {
			isValid, funcError = isAcceptedIf(value, anotherField, validTag.params[1:])
		}
	case "declinedIf":
		if len(validTag.params) >= 2 {
			isValid, funcError = isDeclinedIf(value, anotherField, validTag.params[1:])
		}
	case "prohibitedIf":
		if len(validTag.params) >= 2 {
			isValid, funcError = isProhibitedIf(value, anotherField, validTag.params[1:])
		}
	case "prohibitedUnless":
		if len(validTag.params) >= 2 {
			isValid, funcError = isProhibitedUnless(value, anotherField, validTag.params[1:])
		}
	case "missingIf":
		if len(validTag.params) >= 2 {
			isValid, funcError = isMissingIf(value, anotherField, validTag.params[1:])
		}
	case "missingUnless":
		if len(validTag.params) >= 2 {
			isValid, funcError = isMissingUnless(value, anotherField, validTag.params[1:])
		}
	case "presentIf":
		if len(validTag.params) >= 2 {
			isValid, funcError = isPresentIf(value, anotherField, validTag.params[1:])
		}
	case "presentUnless":
		if len(validTag.params) >= 2 {
			isValid, funcError = isPresentUnless(value, anotherField, validTag.params[1:])
		}
	case "prohibitedIfAccepted":
		isValid, funcError = isProhibitedIfAccepted(value, anotherField)
	case "prohibitedIfDeclined":
		isValid, funcError = isProhibitedIfDeclined(value, anotherField)
	case "prohibits":
		isValid = isProhibits(value, validTag.params, o)
	case "missingWith":
		isValid = isMissingWith(value, validTag.params, o)
	case "missingWithAll":
		isValid = isMissingWithAll(value, validTag.params, o)
	case "presentWith":
		isValid = isPresentWith(value, validTag.params, o)
	case "presentWithAll":
		isValid = isPresentWithAll(value, validTag.params, o)
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

// =============================================================================
// Deprecated: Use Is* functions instead.
// =============================================================================

var (
	ValidateBetween           = IsBetween
	ValidateDigitsBetween     = IsDigitsBetween
	ValidateSize              = IsSize
	ValidateMax               = IsMax
	ValidateMin               = IsMin
	ValidateGtParam           = IsGtParam
	ValidateGteParam          = IsGteParam
	ValidateLtParam           = IsLtParam
	ValidateLteParam          = IsLteParam
	ValidateSame              = IsSame
	ValidateLt                = IsLt
	ValidateLte               = IsLte
	ValidateGt                = IsGt
	ValidateGte               = IsGte
	ValidateDistinct          = IsDistinct
	ValidateMultipleOf        = IsMultipleOf
	ValidateMaxDigits         = IsMaxDigits
	ValidateMinDigits         = IsMinDigits
	ValidateDecimalPrecision  = IsDecimalPrecision
	ValidateContains          = IsContains
	ValidateDoesntContain     = IsDoesntContain
	ValidateMimeTypes         = IsMimeTypes
	ValidateMimes             = IsMimes
	ValidateImage             = IsImage
	ValidateRequired          = IsRequired
	ValidateAccepted          = IsAccepted
	ValidateDeclined          = IsDeclined
	ValidateProhibited        = IsProhibited
	ValidateMissing           = IsMissing
	ValidatePresent           = IsPresent
	ValidateRequiredArrayKeys = IsRequiredArrayKeys
	ValidateList              = IsList
	ValidateDate              = IsDate
	ValidateDateFormat        = IsDateFormat
	ValidateAfter             = IsAfter
	ValidateAfterOrEqual      = IsAfterOrEqual
	ValidateBefore            = IsBefore
	ValidateBeforeOrEqual     = IsBeforeOrEqual
	ValidateIn                = IsIn
	ValidateNotIn             = IsNotIn
	ValidateDifferent         = IsDifferent
	ValidateJSON              = IsJSON
	ValidateRegex             = IsRegex
	ValidateNotRegex          = IsNotRegex
	ValidateBoolean           = IsBoolean
	ValidateString            = IsString
	ValidateArray             = IsArray
	ValidateFilled            = IsFilled
)
