package validator

import (
	"bytes"
	"fmt"
	"net/http"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

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
		valid = ValidateBetweenString(v.String(), minVal, maxVal)
	case reflect.Slice, reflect.Map, reflect.Array:
		minVal, err := ToInt(params[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Between rule on collection field, min value: %w", err)
		}
		maxVal, err := ToInt(params[1])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Between rule on collection field, max value: %w", err)
		}
		valid = ValidateDigitsBetweenInt64(int64(v.Len()), minVal, maxVal)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		minVal, err := ToInt(params[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Between rule on numeric field, min value: %w", err)
		}
		maxVal, err := ToInt(params[1])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Between rule on numeric field, max value: %w", err)
		}
		valid = ValidateDigitsBetweenInt64(v.Int(), minVal, maxVal)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		minVal, err := ToUint(params[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Between rule on numeric field, min value: %w", err)
		}
		maxVal, err := ToUint(params[1])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Between rule on numeric field, max value: %w", err)
		}
		valid = ValidateDigitsBetweenUint64(v.Uint(), minVal, maxVal)
	case reflect.Float32, reflect.Float64:
		minVal, err := ToFloat(params[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Between rule on numeric field, min value: %w", err)
		}
		maxVal, err := ToFloat(params[1])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Between rule on numeric field, max value: %w", err)
		}
		valid = ValidateDigitsBetweenFloat64(v.Float(), minVal, maxVal)
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

		return ValidateBetweenString(value, min, max), nil
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

	var valid bool
	var err error

	switch v.Kind() {
	case reflect.String:
		p, err := ToInt(params[0])
		if err != nil {
			return false, fmt.Errorf("validator: invalid parameter for Gt rule on string field: %w", err)
		}
		valid, err = compareString(v.String(), p, ">")
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
			mimeTypes[i] = Mimes[val]
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

	name := string(append(jsonNamespace, f.nameBytes...))
	structName := string(append(structNamespace, f.structName...))

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
	case reflect.Map:
		values := []string{}
		var sv stringValues
		sv = anotherField.MapKeys()
		sort.Sort(sv)
		for _, k := range sv {
			value := v.MapIndex(k)
			if value.Kind() == reflect.Interface || value.Kind() == reflect.Ptr {
				value = value.Elem()
			}

			if value.Kind() != reflect.Struct {
				values = append(values, ToString(value.Interface()))
			} else {
				return false, fmt.Errorf("validator: requiredUnless unsupported type %T", value.Interface())
			}
		}

		for _, value := range values {
			if !InString(value, params) {
				if Empty(v) {
					return false, nil
				}
			}
		}
	case reflect.Slice, reflect.Array:
		values := []string{}
		for i := 0; i < v.Len(); i++ {
			value := v.Index(i)
			if value.Kind() == reflect.Interface || value.Kind() == reflect.Ptr {
				value = value.Elem()
			}

			if value.Kind() != reflect.Struct {
				values = append(values, ToString(value.Interface()))
			} else {
				return false, fmt.Errorf("validator: requiredUnless unsupported type %T", value.Interface())
			}
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
func validateRequiredWith(otherFields []string, currentField reflect.Value, obj reflect.Value) bool {
	if !allFailingRequired(otherFields, obj) {
		return validateRequired(currentField)
	}
	return true
}

// validateRequiredWithAll The field under validation must be present and not empty only if all of the other specified fields are present.
func validateRequiredWithAll(otherFields []string, currentField reflect.Value, obj reflect.Value) bool {
	if !anyFailingRequired(otherFields, obj) {
		return validateRequired(currentField)
	}
	return true
}

// RequiredWithout The field under validation must be present and not empty only when any of the other specified fields are not present.
func validateRequiredWithout(otherFields []string, currentField reflect.Value, obj reflect.Value) bool {
	if anyFailingRequired(otherFields, obj) {
		return validateRequired(currentField)
	}
	return true
}

// validateRequiredWithoutAll The field under validation must be present and not empty only when all of the other specified fields are not present.
func validateRequiredWithoutAll(otherFields []string, currentField reflect.Value, obj reflect.Value) bool {
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
	case "same":
		anotherField, err = findField(validTag.params[0], o)
		if err != nil {
			return false, nil
		}
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
