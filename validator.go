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

// Validator contruct
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
func validateBetween(v reflect.Value, params []string) bool {
	if len(params) != 2 {
		return false
	}

	switch v.Kind() {
	case reflect.String:
		min, _ := ToInt(params[0])
		max, _ := ToInt(params[1])
		return ValidateBetweenString(v.String(), min, max)
	case reflect.Slice, reflect.Map, reflect.Array:
		min, _ := ToInt(params[0])
		max, _ := ToInt(params[1])
		return ValidateDigitsBetweenInt64(int64(v.Len()), min, max)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		min, _ := ToInt(params[0])
		max, _ := ToInt(params[1])
		return ValidateDigitsBetweenInt64(v.Int(), min, max)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		min, _ := ToUint(params[0])
		max, _ := ToUint(params[1])
		return ValidateDigitsBetweenUint64(v.Uint(), min, max)
	case reflect.Float32, reflect.Float64:
		min, _ := ToFloat(params[0])
		max, _ := ToFloat(params[1])
		return ValidateDigitsBetweenFloat64(v.Float(), min, max)
	}

	panic(fmt.Sprintf("validator: Between unsupport Type %T", v.Interface()))
}

// ValidateBetween check The field under validation must have a size between the given min and max. Strings, numerics, arrays, and files are evaluated in the same fashion as the size rule.
func ValidateBetween(i interface{}, params []string) bool {
	v := reflect.ValueOf(i)
	return validateBetween(v, params)
}

// validateDigitsBetween check The field under validation must have a length between the given min and max.
func validateDigitsBetween(v reflect.Value, params []string) bool {
	if len(params) != 2 {
		return false
	}

	switch v.Kind() {
	case reflect.String, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		min, _ := ToInt(params[0])
		max, _ := ToInt(params[1])
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
			return false
		}

		return ValidateBetweenString(value, min, max)
	}

	panic(fmt.Sprintf("validator: DigitsBetween unsupport Type %T", v.Interface()))
}

// ValidateDigitsBetween check The field under validation must have a length between the given min and max.
func ValidateDigitsBetween(i interface{}, params []string) bool {
	v := reflect.ValueOf(i)
	return validateDigitsBetween(v, params)
}

// validateSize The field under validation must have a size matching the given value.
// For string data, value corresponds to the number of characters.
// For numeric data, value corresponds to a given integer value.
// For an array | map | slice, size corresponds to the count of the array | map | slice.
func validateSize(v reflect.Value, param []string) bool {
	switch v.Kind() {
	case reflect.String:
		p, _ := ToInt(param[0])
		return compareString(v.String(), p, "==")
	case reflect.Slice, reflect.Map, reflect.Array:
		p, _ := ToInt(param[0])
		return compareInt64(int64(v.Len()), p, "==")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p, _ := ToInt(param[0])
		return compareInt64(v.Int(), p, "==")
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p, _ := ToUint(param[0])
		return compareUint64(v.Uint(), p, "==")
	case reflect.Float32, reflect.Float64:
		p, _ := ToFloat(param[0])
		return compareFloat64(v.Float(), p, "==")
	}

	panic(fmt.Sprintf("validator: Size unsupport Type %T", v.Interface()))
}

// ValidateSize The field under validation must have a size matching the given value.
// For string data, value corresponds to the number of characters.
// For numeric data, value corresponds to a given integer value.
// For an array | map | slice, size corresponds to the count of the array | map | slice.
func ValidateSize(i interface{}, params []string) bool {
	v := reflect.ValueOf(i)
	return validateSize(v, params)
}

// validateMax is the validation function for validating if the current field's value is less than or equal to the param's value.
func validateMax(v reflect.Value, param []string) bool {
	switch v.Kind() {
	case reflect.String:
		p, _ := ToInt(param[0])
		return compareString(v.String(), p, "<=")
	case reflect.Slice, reflect.Map, reflect.Array:
		p, _ := ToInt(param[0])
		return compareInt64(int64(v.Len()), p, "<=")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p, _ := ToInt(param[0])
		return compareInt64(v.Int(), p, "<=")
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p, _ := ToUint(param[0])
		return compareUint64(v.Uint(), p, "<=")
	case reflect.Float32, reflect.Float64:
		p, _ := ToFloat(param[0])
		return compareFloat64(v.Float(), p, "<=")
	}

	panic(fmt.Sprintf("validator: Max unsupport Type %T", v.Interface()))
}

// ValidatMax is the validation function for validating if the current field's value is less than or equal to the param's value.
func ValidatMax(i interface{}, params []string) bool {
	v := reflect.ValueOf(i)
	return validateMax(v, params)
}

// validateMin is the validation function for validating if the current field's value is greater than or equal to the param's value.
func validateMin(v reflect.Value, param []string) bool {
	switch v.Kind() {
	case reflect.String:
		p, _ := ToInt(param[0])
		return compareString(v.String(), p, ">=")
	case reflect.Slice, reflect.Map, reflect.Array:
		p, _ := ToInt(param[0])
		return compareInt64(int64(v.Len()), p, ">=")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p, _ := ToInt(param[0])
		return compareInt64(v.Int(), p, ">=")
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p, _ := ToUint(param[0])
		return compareUint64(v.Uint(), p, ">=")
	case reflect.Float32, reflect.Float64:
		p, _ := ToFloat(param[0])
		return compareFloat64(v.Float(), p, ">=")
	}

	panic(fmt.Sprintf("validator: Min unsupport Type %T", v.Interface()))
}

// ValidateMin is the validation function for validating if the current field's value is greater than or equal to the param's value.
func ValidateMin(i interface{}, params []string) bool {
	v := reflect.ValueOf(i)
	return validateMin(v, params)
}

// validateSame is the validation function for validating if the current field's value euqal the param's value.
func validateSame(v reflect.Value, anotherField reflect.Value) bool {
	if v.Kind() != anotherField.Kind() {
		panic(fmt.Sprintf("validator: Same The two fields must be of the same type %T, %T", v.Interface(), anotherField.Interface()))
	}

	switch v.Kind() {
	case reflect.String:
		return v.String() == anotherField.String()
	case reflect.Slice, reflect.Map, reflect.Array:
		return compareInt64(int64(v.Len()), int64(anotherField.Len()), "==")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return compareInt64(v.Int(), anotherField.Int(), "==")
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return compareUint64(v.Uint(), anotherField.Uint(), "==")
	case reflect.Float32, reflect.Float64:
		return compareFloat64(v.Float(), anotherField.Float(), "==")
	}

	panic(fmt.Sprintf("validator: Same unsupport Type %T", v.Interface()))
}

// ValidateSame is the validation function for validating if the current field's value is greater than or equal to the param's value.
func ValidateSame(i interface{}, a interface{}) bool {
	v := reflect.ValueOf(i)
	anotherField := reflect.ValueOf(a)
	return validateSame(v, anotherField)
}

// validateLt is the validation function for validating if the current field's value is less than the param's value.
func validateLt(v reflect.Value, anotherField reflect.Value) bool {
	if v.Kind() != anotherField.Kind() {
		panic(fmt.Sprintf("validator: Lt The two fields must be of the same type %T, %T", v.Interface(), anotherField.Interface()))
	}

	switch v.Kind() {
	case reflect.String:
		return compareString(v.String(), int64(utf8.RuneCountInString(anotherField.String())), "<")
	case reflect.Slice, reflect.Map, reflect.Array:
		return compareInt64(int64(v.Len()), int64(anotherField.Len()), "<")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return compareInt64(v.Int(), anotherField.Int(), "<")
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return compareUint64(v.Uint(), anotherField.Uint(), "<")
	case reflect.Float32, reflect.Float64:
		return compareFloat64(v.Float(), anotherField.Float(), "<")
	}

	panic(fmt.Sprintf("validator: Lt unsupport Type %T", v.Interface()))
}

// ValidateLt is the validation function for validating if the current field's value is less than the param's value.
func ValidateLt(i interface{}, a interface{}) bool {
	v := reflect.ValueOf(i)
	anotherField := reflect.ValueOf(a)
	return validateLt(v, anotherField)
}

// validateLte is the validation function for validating if the current field's value is less than or equal to the param's value.
func validateLte(v reflect.Value, anotherField reflect.Value) bool {
	if v.Kind() != anotherField.Kind() {
		panic(fmt.Sprintf("validator: Lte The two fields must be of the same type %T, %T", v.Interface(), anotherField.Interface()))
	}

	switch v.Kind() {
	case reflect.String:
		return compareString(v.String(), int64(utf8.RuneCountInString(anotherField.String())), "<=")
	case reflect.Slice, reflect.Map, reflect.Array:
		return compareInt64(int64(v.Len()), int64(anotherField.Len()), "<=")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return compareInt64(v.Int(), anotherField.Int(), "<=")
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return compareUint64(v.Uint(), anotherField.Uint(), "<=")
	case reflect.Float32, reflect.Float64:
		return compareFloat64(v.Float(), anotherField.Float(), "<=")
	}

	panic(fmt.Sprintf("validator: Lte unsupport Type %T", v.Interface()))
}

// ValidateLte is the validation function for validating if the current field's value is less than or equal to the param's value.
func ValidateLte(i interface{}, a interface{}) bool {
	v := reflect.ValueOf(i)
	anotherField := reflect.ValueOf(a)
	return validateLte(v, anotherField)
}

// validateGt is the validation function for validating if the current field's value is greater than to the param's value.
func validateGt(v reflect.Value, anotherField reflect.Value) bool {
	if v.Kind() != anotherField.Kind() {
		panic(fmt.Sprintf("validator: Gt The two fields must be of the same type %T, %T", v.Interface(), anotherField.Interface()))
	}

	switch v.Kind() {
	case reflect.String:
		return compareString(v.String(), int64(utf8.RuneCountInString(anotherField.String())), ">")
	case reflect.Slice, reflect.Map, reflect.Array:
		return compareInt64(int64(v.Len()), int64(anotherField.Len()), ">")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return compareInt64(v.Int(), anotherField.Int(), ">")
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return compareUint64(v.Uint(), anotherField.Uint(), ">")
	case reflect.Float32, reflect.Float64:
		return compareFloat64(v.Float(), anotherField.Float(), ">")
	}

	panic(fmt.Sprintf("validator: Gt unsupport Type %T", v.Interface()))
}

// ValidateGt is the validation function for validating if the current field's value is greater than to the param's value.
func ValidateGt(i interface{}, a interface{}) bool {
	v := reflect.ValueOf(i)
	anotherField := reflect.ValueOf(a)
	return validateGt(v, anotherField)
}

// validateGte is the validation function for validating if the current field's value is greater than or equal to the param's value.
func validateGte(v reflect.Value, anotherField reflect.Value) bool {
	if v.Kind() != anotherField.Kind() {
		panic(fmt.Sprintf("validator: Gte The two fields must be of the same type %T, %T", v.Interface(), anotherField.Interface()))
	}

	switch v.Kind() {
	case reflect.String:
		return compareString(v.String(), int64(utf8.RuneCountInString(anotherField.String())), ">=")
	case reflect.Slice, reflect.Map, reflect.Array:
		return compareInt64(int64(v.Len()), int64(anotherField.Len()), ">=")
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return compareInt64(v.Int(), anotherField.Int(), ">=")
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return compareUint64(v.Uint(), anotherField.Uint(), ">=")
	case reflect.Float32, reflect.Float64:
		return compareFloat64(v.Float(), anotherField.Float(), ">=")
	}

	panic(fmt.Sprintf("validator: Gte unsupport Type %T", v.Interface()))
}

// ValidateGte is the validation function for validating if the current field's value is greater than to the param's value.
func ValidateGte(i interface{}, a interface{}) bool {
	v := reflect.ValueOf(i)
	anotherField := reflect.ValueOf(a)
	return validateGte(v, anotherField)
}

// validateDistinct is the validation function for validating an attribute is unique among other values.
func validateDistinct(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64:
		return true
	case reflect.Slice, reflect.Array:
		m := reflect.MakeMap(reflect.MapOf(v.Type().Elem(), v.Type()))

		for i := 0; i < v.Len(); i++ {
			m.SetMapIndex(v.Index(i), v)
		}
		return v.Len() == m.Len()
	case reflect.Map:
		m := reflect.MakeMap(reflect.MapOf(v.Type().Elem(), v.Type()))

		for _, k := range v.MapKeys() {
			m.SetMapIndex(v.MapIndex(k), v)
		}
		return v.Len() == m.Len()
	}

	panic(fmt.Sprintf("validator: Distinct unsupport Type %T", v.Interface()))
}

// ValidateDistinct is the validation function for validating an attribute is unique among other values.
func ValidateDistinct(i interface{}) bool {
	v := reflect.ValueOf(i)
	return validateDistinct(v)
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
func ValidateMimes(data []byte, mimes []string) bool {
	mimeTypes := make([]string, len(mimes))
	for i, mime := range mimes {
		if val, ok := Mimes[mime]; ok {
			mimeTypes[i] = Mimes[val]
		} else {
			panic(fmt.Sprintf("validator: Mimes unsupport Type %s", mime))
		}
	}

	return ValidateMimeTypes(data, mimeTypes)
}

// ValidateImage is the validation function for the The file under validation must be an image (jpeg, png, bmp, gif, or svg)
func ValidateImage(data []byte) bool {
	return ValidateMimes(data, []string{"jpeg", "png", "gif", "bmp", "svg"})
}

// ValidateStruct use tags for fields.
// result will be equal to `false` if there are any errors.
func ValidateStruct(s interface{}) error {
	return Default.ValidateStruct(s, nil, nil)
}

// ValidateStruct use tags for fields.
// result will be equal to `false` if there are any errors.
func (v *Validator) ValidateStruct(s interface{}, jsonNamespace []byte, structNamespace []byte) error {
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

	for _, f := range fields {
		valuefield := val.Field(f.index[0])
		err := v.newTypeValidator(valuefield, &f, val, jsonNamespace, structNamespace)
		if err != nil {
			if errors, ok := err.(Errors); ok {
				for _, fieldError := range errors {
					errs = append(errs, fieldError)
				}
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

func (v *Validator) newTypeValidator(value reflect.Value, f *field, o reflect.Value, jsonNamespace []byte, structNamespace []byte) (resultErr error) {
	if !value.IsValid() || f.omitEmpty && Empty(value) {
		return nil
	}

	name := string(append(jsonNamespace, f.nameBytes...))
	structName := string(append(structNamespace, f.structName...))

	if err := v.checkRequired(value, f, o, name, structName); err != nil {
		return err
	}

	for _, tag := range f.validTags {
		if validatefunc, ok := CustomTypeRuleMap.Get(tag.name); ok {
			if result := validatefunc(value, o, tag); !result {
				return v.formatsMessages(&FieldError{
					name:              name,
					structName:        structName,
					tag:               tag.name,
					messageName:       tag.messageName,
					messageParameters: parseValidatorMessageParameters(tag, o),
					attribute:         f.attribute,
					defaultAttribute:  f.defaultAttribute,
					value:             ToString(value.Interface()),
				})
			}
		}
	}

	switch value.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.String:

		for _, tag := range f.validTags {

			if err := v.checkDependentRules(tag, f, value, o, name, structName); err != nil {
				return err
			}

			if validfunc, ok := RuleMap[tag.name]; ok {
				isValid := validfunc(value)
				if !isValid {
					return v.formatsMessages(&FieldError{
						name:              name,
						structName:        structName,
						tag:               tag.name,
						messageName:       tag.messageName,
						messageParameters: parseValidatorMessageParameters(tag, o),
						attribute:         f.attribute,
						defaultAttribute:  f.defaultAttribute,
						value:             ToString(value.Interface()),
					})
				}
			}

			if validfunc, ok := ParamRuleMap[tag.name]; ok {
				isValid := validfunc(value, tag.params)
				if !isValid {
					return v.formatsMessages(&FieldError{
						name:              name,
						structName:        structName,
						tag:               tag.name,
						messageName:       tag.messageName,
						messageParameters: parseValidatorMessageParameters(tag, o),
						attribute:         f.attribute,
						defaultAttribute:  f.defaultAttribute,
						value:             ToString(value.Interface()),
					})
				}
			}

			switch value.Kind() {
			case reflect.String:
				if validfunc, ok := StringRulesMap[tag.name]; ok {
					isValid := validfunc(value.String())
					if !isValid {
						return v.formatsMessages(&FieldError{
							name:              name,
							structName:        structName,
							tag:               tag.name,
							messageName:       tag.messageName,
							messageParameters: parseValidatorMessageParameters(tag, o),
							attribute:         f.attribute,
							defaultAttribute:  f.defaultAttribute,
							value:             ToString(value.Interface()),
						})
					}
				}
			}
		}
		return nil
	case reflect.Map:
		if value.Type().Key().Kind() != reflect.String {
			return &UnsupportedTypeError{value.Type()}
		}

		for _, tag := range f.validTags {

			if err := v.checkDependentRules(tag, f, value, o, name, structName); err != nil {
				return err
			}

			if validfunc, ok := RuleMap[tag.name]; ok {
				isValid := validfunc(value)
				if !isValid {
					return v.formatsMessages(&FieldError{
						name:              name,
						structName:        structName,
						tag:               tag.name,
						messageName:       tag.messageName,
						messageParameters: parseValidatorMessageParameters(tag, o),
						attribute:         f.attribute,
						defaultAttribute:  f.defaultAttribute,
						value:             ToString(value.Interface()),
					})
				}
			}

			if validfunc, ok := ParamRuleMap[tag.name]; ok {
				isValid := validfunc(value, tag.params)
				if !isValid {
					return v.formatsMessages(&FieldError{
						name:              name,
						structName:        structName,
						tag:               tag.name,
						messageName:       tag.messageName,
						messageParameters: parseValidatorMessageParameters(tag, o),
						attribute:         f.attribute,
						defaultAttribute:  f.defaultAttribute,
						value:             ToString(value.Interface()),
					})
				}
			}
		}

		var sv stringValues
		sv = value.MapKeys()
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
	case reflect.Slice, reflect.Array:
		for _, tag := range f.validTags {
			if err := v.checkDependentRules(tag, f, value, o, name, structName); err != nil {
				return err
			}

			if validfunc, ok := RuleMap[tag.name]; ok {
				isValid := validfunc(value)
				if !isValid {
					return v.formatsMessages(&FieldError{
						name:              name,
						structName:        structName,
						tag:               tag.name,
						messageName:       tag.messageName,
						messageParameters: parseValidatorMessageParameters(tag, o),
						attribute:         f.attribute,
						defaultAttribute:  f.defaultAttribute,
						value:             ToString(value.Interface()),
					})
				}
			}

			if validfunc, ok := ParamRuleMap[tag.name]; ok {
				isValid := validfunc(value, tag.params)
				if !isValid {
					return v.formatsMessages(&FieldError{
						name:              name,
						structName:        structName,
						tag:               tag.name,
						messageName:       tag.messageName,
						messageParameters: parseValidatorMessageParameters(tag, o),
						attribute:         f.attribute,
						defaultAttribute:  f.defaultAttribute,
						value:             ToString(value.Interface()),
					})
				}
			}
		}

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
	case reflect.Interface:
		if value.IsNil() {
			return nil
		}
		return v.ValidateStruct(value.Interface(), jsonNamespace, structNamespace)
	case reflect.Ptr:
		if value.IsNil() {
			return nil
		}
		jsonNamespace = append(append(jsonNamespace, f.nameBytes...), '.')
		structNamespace = append(append(structNamespace, f.structNameBytes...), '.')
		return v.ValidateStruct(value.Interface(), jsonNamespace, structNamespace)
	case reflect.Struct:
		jsonNamespace = append(append(jsonNamespace, f.nameBytes...), '.')
		structNamespace = append(append(structNamespace, f.structNameBytes...), '.')
		return v.ValidateStruct(value.Interface(), jsonNamespace, structNamespace)
	default:
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
func validateRequiredIf(v reflect.Value, anotherField reflect.Value, params []string, tag *ValidTag) bool {
	if anotherField.Kind() == reflect.Interface || anotherField.Kind() == reflect.Ptr {
		anotherField = anotherField.Elem()
	}

	if !anotherField.IsValid() {
		return true
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
			return false
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
				panic(fmt.Sprintf("validator: RequiredIf unsupport Type %T", value.Interface()))
			}
		}

		for _, value := range values {
			if InString(value, params) && Empty(v) {
				tag.messageParameters = append(
					tag.messageParameters,
					messageParameter{
						Key:   "Value",
						Value: value,
					},
				)
				return false
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
				panic(fmt.Sprintf("validator: RequiredIf unsupport Type %T", value.Interface()))
			}
		}

		for _, value := range values {
			if InString(value, params) {
				if Empty(v) {
					tag.messageParameters = append(
						tag.messageParameters,
						messageParameter{
							Key:   "Value",
							Value: value,
						},
					)
					return false
				}
			}
		}
	default:
		panic(fmt.Sprintf("validator: RequiredIf unsupport Type %T", anotherField.Interface()))
	}

	return true
}

// validateRequiredUnless check value required when anotherField str is a member of the set of strings params
func validateRequiredUnless(v reflect.Value, anotherField reflect.Value, params []string) bool {
	if anotherField.Kind() == reflect.Interface || anotherField.Kind() == reflect.Ptr {
		anotherField = anotherField.Elem()
	}

	if !anotherField.IsValid() {
		return true
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
				return false
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
				panic(fmt.Sprintf("validator: requiredUnless unsupport Type %T", value.Interface()))
			}
		}

		for _, value := range values {
			if !InString(value, params) {
				if Empty(v) {
					return false
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
				panic(fmt.Sprintf("validator: requiredUnless unsupport Type %T", value.Interface()))
			}
		}

		for _, value := range values {
			if !InString(value, params) {
				if Empty(v) {
					return false
				}
			}
		}
	default:
		panic(fmt.Sprintf("validator: requiredUnless unsupport Type %T", anotherField.Interface()))
	}

	return true
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

func (v *Validator) checkRequired(value reflect.Value, f *field, o reflect.Value, name string, structName string) *FieldError {
	for _, tag := range f.requiredTags {
		isError := false
		switch tag.name {
		case "required":
			isError = !validateRequired(value)
		case "requiredIf":
			anotherField, err := findField(tag.params[0], o)
			if err == nil && len(tag.params) >= 2 && !validateRequiredIf(value, anotherField, tag.params[1:], tag) {
				isError = true
			}
		case "requiredUnless":
			anotherField, err := findField(tag.params[0], o)
			if err == nil && len(tag.params) >= 2 && !validateRequiredUnless(value, anotherField, tag.params[1:]) {
				isError = true
			}
		case "requiredWith":
			if !validateRequiredWith(tag.params, value) {
				isError = true
			}
		case "requiredWithAll":
			if !validateRequiredWithAll(tag.params, value) {
				isError = true
			}
		case "requiredWithout":
			if !validateRequiredWithout(tag.params, value) {
				isError = true
			}
		case "requiredWithoutAll":
			if !validateRequiredWithoutAll(tag.params, value) {
				isError = true
			}
		}

		if isError {
			return v.formatsMessages(&FieldError{
				name:              name,
				structName:        structName,
				tag:               tag.name,
				messageName:       tag.messageName,
				messageParameters: parseValidatorMessageParameters(tag, o),
				attribute:         f.attribute,
				defaultAttribute:  f.defaultAttribute,
				value:             ToString(value.Interface()),
			})
		}
	}

	return nil
}

// validateRequiredWith The field under validation must be present and not empty only if any of the other specified fields are present.
func validateRequiredWith(otherFields []string, v reflect.Value) bool {
	if !allFailingRequired(otherFields, v) {
		return validateRequired(v)
	}
	return true
}

// validateRequiredWithAll The field under validation must be present and not empty only if all of the other specified fields are present.
func validateRequiredWithAll(otherFields []string, v reflect.Value) bool {
	if !anyFailingRequired(otherFields, v) {
		return validateRequired(v)
	}
	return true
}

// RequiredWithout The field under validation must be present and not empty only when any of the other specified fields are not present.
func validateRequiredWithout(otherFields []string, v reflect.Value) bool {
	if anyFailingRequired(otherFields, v) {
		return validateRequired(v)
	}
	return true
}

// validateRequiredWithoutAll The field under validation must be present and not empty only when all of the other specified fields are not present.
func validateRequiredWithoutAll(otherFields []string, v reflect.Value) bool {
	if allFailingRequired(otherFields, v) {
		return validateRequired(v)
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

	if message, ok = v.CustomMessage[fieldError.structName+"."+fieldError.messageName]; ok {
		fieldError.err = fmt.Errorf(message)
		return fieldError
	}

	message, ok = MessageMap[fieldError.messageName]
	if ok {
		attribute := fieldError.attribute
		if customAttribute, ok := v.Attributes[fieldError.structName]; ok {
			attribute = customAttribute
		} else if fieldError.defaultAttribute != "" {
			attribute = fieldError.defaultAttribute
		}
		message = replaceAttributes(message, attribute, fieldError.messageParameters)

		fieldError.err = fmt.Errorf(message)
		return fieldError
	}

	fieldError.err = fmt.Errorf("validator: undefined message : %s", fieldError.messageName)
	return fieldError
}

func replaceAttributes(message string, attribute string, messageParameters MessageParameters) string {
	message = strings.Replace(message, "{{.Attribute}}", attribute, -1)
	for _, parameter := range messageParameters {
		message = strings.Replace(message, "{{."+parameter.Key+"}}", parameter.Value, -1)
	}
	return message
}

func getDisplayableAttribute(o reflect.Value, attribute string) string {
	attributes := strings.Split(attribute, ".")
	if len(attributes) > 0 {
		attribute = o.Type().Name() + attributes[0]
	} else {
		attribute = strings.Join(attributes[len(attributes)-2:], ".")
	}

	return attributes[len(attributes)-1]
}

func findField(fieldName string, v reflect.Value) (reflect.Value, error) {
	fields := strings.Split(fieldName, ".")
	current := v.FieldByName(fields[0])
	i := 1
	if len(fields) > i {
		for true {
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

func (v *Validator) checkDependentRules(validTag *ValidTag, f *field, value reflect.Value, o reflect.Value, name string, structName string) error {
	isValid := true
	var anotherField reflect.Value
	var err error
	switch validTag.name {
	case "gt", "gte", "lt", "lte", "same":
		anotherField, err = findField(validTag.params[0], o)
		if err != nil {
			return nil
		}
	}

	switch validTag.name {
	case "gt":
		isValid = validateGt(value, anotherField)
	case "gte":
		isValid = validateGte(value, anotherField)
	case "lt":
		isValid = validateLt(value, anotherField)
	case "lte":
		isValid = validateLte(value, anotherField)
	case "same":
		isValid = validateSame(value, anotherField)
	}

	if !isValid {
		return v.formatsMessages(&FieldError{
			name:              name,
			structName:        structName,
			tag:               validTag.name,
			messageName:       validTag.messageName,
			messageParameters: parseValidatorMessageParameters(validTag, o),
			attribute:         f.attribute,
			defaultAttribute:  f.defaultAttribute,
			value:             ToString(value.Interface()),
		})
	}

	return nil
}
