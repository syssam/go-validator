package validator

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"
	"unicode/utf8"
)

const tagName string = "valid"

// Validator contruct
type Validator struct {
	Translator    *Translator
	Attributes    map[string]string
	CustomMessage map[string]string
}

var loadValidatorOnce *Validator
var once sync.Once

// New returns a new instance of 'valid' with sane defaults.
func New() *Validator {
	once.Do(func() {
		loadValidatorOnce = &Validator{}
	})
	return loadValidatorOnce
}

// newValidator returns a new instance of 'valid' with sane defaults.
func newValidator() *Validator {
	once.Do(func() {
		loadValidatorOnce = &Validator{}
	})
	return loadValidatorOnce
}

// Between check The field under validation must have a size between the given min and max. Strings, numerics, arrays, and files are evaluated in the same fashion as the size rule.
func Between(v reflect.Value, params ...string) bool {
	if len(params) != 2 {
		return false
	}

	switch v.Kind() {
	case reflect.String:
		min, _ := ToInt(params[0])
		max, _ := ToInt(params[1])
		return BetweenString(v.String(), min, max)
	case reflect.Slice, reflect.Map, reflect.Array:
		min, _ := ToInt(params[0])
		max, _ := ToInt(params[1])
		return DigitsBetweenInt64(int64(v.Len()), min, max)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		min, _ := ToInt(params[0])
		max, _ := ToInt(params[1])
		return DigitsBetweenInt64(v.Int(), min, max)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		min, _ := ToUint(params[0])
		max, _ := ToUint(params[1])
		return DigitsBetweenUint64(v.Uint(), min, max)
	case reflect.Float32, reflect.Float64:
		min, _ := ToFloat(params[0])
		max, _ := ToFloat(params[1])
		return DigitsBetweenFloat64(v.Float(), min, max)
	}

	panic(fmt.Sprintf("validator: Between unsupport Type %T", v.Interface()))
}

// DigitsBetween check The field under validation must have a length between the given min and max.
func DigitsBetween(v reflect.Value, params ...string) bool {
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

		return BetweenString(value, min, max)
	}

	panic(fmt.Sprintf("validator: DigitsBetween unsupport Type %T", v.Interface()))
}

// Size The field under validation must have a size matching the given value.
// For string data, value corresponds to the number of characters.
// For numeric data, value corresponds to a given integer value.
// For an array | map | slice, size corresponds to the count of the array | map | slice.
func Size(v reflect.Value, param ...string) bool {
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
		return compareUnit64(v.Uint(), p, "==")
	case reflect.Float32, reflect.Float64:
		p, _ := ToFloat(param[0])
		return compareFloat64(v.Float(), p, "==")
	}

	panic(fmt.Sprintf("validator: Size unsupport Type %T", v.Interface()))
}

// Max is the validation function for validating if the current field's value is less than or equal to the param's value.
func Max(v reflect.Value, param ...string) bool {
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
		return compareUnit64(v.Uint(), p, "<=")
	case reflect.Float32, reflect.Float64:
		p, _ := ToFloat(param[0])
		return compareFloat64(v.Float(), p, "<=")
	}

	panic(fmt.Sprintf("validator: Max unsupport Type %T", v.Interface()))
}

// Min is the validation function for validating if the current field's value is greater than or equal to the param's value.
func Min(v reflect.Value, param ...string) bool {
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
		return compareUnit64(v.Uint(), p, ">=")
	case reflect.Float32, reflect.Float64:
		p, _ := ToFloat(param[0])
		return compareFloat64(v.Float(), p, ">=")
	}

	panic(fmt.Sprintf("validator: Min unsupport Type %T", v.Interface()))
}

// Same is the validation function for validating if the current field's value euqal the param's value.
func Same(v reflect.Value, anotherField reflect.Value) bool {
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
		return compareUnit64(v.Uint(), anotherField.Uint(), "==")
	case reflect.Float32, reflect.Float64:
		return compareFloat64(v.Float(), anotherField.Float(), "==")
	}

	panic(fmt.Sprintf("validator: Lt unsupport Type %T", v.Interface()))
}

// Lt is the validation function for validating if the current field's value is less than the param's value.
func Lt(v reflect.Value, anotherField reflect.Value) bool {
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
		return compareUnit64(v.Uint(), anotherField.Uint(), "<")
	case reflect.Float32, reflect.Float64:
		return compareFloat64(v.Float(), anotherField.Float(), "<")
	}

	panic(fmt.Sprintf("validator: Lt unsupport Type %T", v.Interface()))
}

// Lte is the validation function for validating if the current field's value is less than or equal to the param's value.
func Lte(v reflect.Value, anotherField reflect.Value) bool {
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
		return compareUnit64(v.Uint(), anotherField.Uint(), "<=")
	case reflect.Float32, reflect.Float64:
		return compareFloat64(v.Float(), anotherField.Float(), "<=")
	}

	panic(fmt.Sprintf("validator: Lte unsupport Type %T", v.Interface()))
}

// Gt is the validation function for validating if the current field's value is greater than to the param's value.
func Gt(v reflect.Value, anotherField reflect.Value) bool {
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
		return compareUnit64(v.Uint(), anotherField.Uint(), ">")
	case reflect.Float32, reflect.Float64:
		return compareFloat64(v.Float(), anotherField.Float(), ">")
	}

	panic(fmt.Sprintf("validator: Gt unsupport Type %T", v.Interface()))
}

// Gte is the validation function for validating if the current field's value is greater than or equal to the param's value.
func Gte(v reflect.Value, anotherField reflect.Value) bool {
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
		return compareUnit64(v.Uint(), anotherField.Uint(), ">=")
	case reflect.Float32, reflect.Float64:
		return compareFloat64(v.Float(), anotherField.Float(), ">=")
	}

	panic(fmt.Sprintf("validator: Gte unsupport Type %T", v.Interface()))
}

func validateStruct(s interface{}, jsonNamespace []byte, structNamespace []byte) error {
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
		err := newTypeValidator(valuefield, &f, val, jsonNamespace, structNamespace)
		if err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		err = errs
	}

	return err
}

// ValidateStruct use tags for fields.
// result will be equal to `false` if there are any errors.
func ValidateStruct(s interface{}) error {
	newValidator()
	return validateStruct(s, nil, nil)
}

func newTypeValidator(v reflect.Value, f *field, o reflect.Value, jsonNamespace []byte, structNamespace []byte) (resultErr error) {
	if !v.IsValid() || f.omitEmpty && Empty(v) {
		return nil

	}
	name := string(append(jsonNamespace, f.nameBytes...))
	structName := string(append(structNamespace, f.structName...))

	if err := checkRequired(v, f, o, name, structName); err != nil {
		return err
	}

	for _, tag := range f.validTags {
		if validatefunc, ok := CustomTypeRuleMap.Get(tag.name); ok {
			if result := validatefunc(v, o, tag); !result {
				return &Error{
					Name:       name,
					StructName: structName,
					Err:        formatsMessages(tag, v, f, o),
					Tag:        tag.name,
				}
			}
		}
	}

	switch v.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.String:

		for _, tag := range f.validTags {

			if err := checkDependentRules(tag, f, v, o, name, structName); err != nil {
				return err
			}

			if validfunc, ok := RuleMap[tag.name]; ok {
				isValid := validfunc(v)
				if !isValid {
					return &Error{
						Name:       name,
						StructName: structName,
						Err:        formatsMessages(tag, v, f, o),
						Tag:        tag.name,
					}
				}
			}

			if validfunc, ok := ParamRuleMap[tag.name]; ok {
				isValid := validfunc(v, tag.params...)
				if !isValid {
					return &Error{
						Name:       name,
						StructName: structName,
						Err:        formatsMessages(tag, v, f, o),
						Tag:        tag.name,
					}
				}
			}

			switch v.Kind() {
			case reflect.String:
				if validfunc, ok := StringRulesMap[tag.name]; ok {
					isValid := validfunc(v.String())
					if !isValid {
						return &Error{
							Name:       name,
							StructName: structName,
							Err:        formatsMessages(tag, v, f, o),
							Tag:        tag.name,
						}
					}
				}
			}
		}
		return nil
	case reflect.Map:
		if v.Type().Key().Kind() != reflect.String {
			return &UnsupportedTypeError{v.Type()}
		}

		for _, tag := range f.validTags {

			if err := checkDependentRules(tag, f, v, o, name, structName); err != nil {
				return err
			}

			if validfunc, ok := ParamRuleMap[tag.name]; ok {
				isValid := validfunc(v, tag.params...)
				if !isValid {
					return &Error{
						Name:       name,
						StructName: structName,
						Err:        formatsMessages(tag, v, f, o),
						Tag:        tag.name,
					}
				}
			}
		}

		var sv stringValues
		sv = v.MapKeys()
		sort.Sort(sv)
		for _, k := range sv {
			var err error
			value := v.MapIndex(k)
			if value.Kind() == reflect.Interface {
				value = value.Elem()
			}

			if value.Kind() == reflect.Struct || value.Kind() == reflect.Ptr {
				jsonNamespace = append(append(jsonNamespace, f.nameBytes...), '.')
				structNamespace = append(append(structNamespace, f.structNameBytes...), '.')
				err = validateStruct(value.Interface(), jsonNamespace, structNamespace)
				if err != nil {
					return err
				}
			}
		}
		return nil
	case reflect.Slice, reflect.Array:
		for _, tag := range f.validTags {
			if err := checkDependentRules(tag, f, v, o, name, structName); err != nil {
				return err
			}
			if validfunc, ok := ParamRuleMap[tag.name]; ok {
				isValid := validfunc(v, tag.params...)
				if !isValid {
					return &Error{
						Name:       name,
						StructName: structName,
						Err:        formatsMessages(tag, v, f, o),
						Tag:        tag.name,
					}
				}
			}
		}
		for i := 0; i < v.Len(); i++ {
			var err error
			value := v.Index(i)
			if value.Kind() == reflect.Interface {
				value = value.Elem()
			}
			if value.Kind() == reflect.Struct || value.Kind() == reflect.Ptr {
				jsonNamespace = append(append(jsonNamespace, f.nameBytes...), '.')
				structNamespace = append(append(structNamespace, f.structNameBytes...), '.')
				err = validateStruct(v.Index(i).Interface(), jsonNamespace, structNamespace)
				if err != nil {
					return err
				}
			}
		}
		return nil
	case reflect.Interface:
		// If the value is an interface then encode its element
		if v.IsNil() {
			return nil
		}
		return validateStruct(v.Interface(), jsonNamespace, structNamespace)
	case reflect.Ptr:
		// If the value is a pointer then check its element
		if v.IsNil() {
			return nil
		}

		jsonNamespace = append(append(jsonNamespace, f.nameBytes...), '.')
		structNamespace = append(append(structNamespace, f.structNameBytes...), '.')
		return validateStruct(v.Interface(), jsonNamespace, structNamespace)
	case reflect.Struct:
		jsonNamespace = append(append(jsonNamespace, f.nameBytes...), '.')
		structNamespace = append(append(structNamespace, f.structNameBytes...), '.')
		return validateStruct(v.Interface(), jsonNamespace, structNamespace)
	default:
		return &UnsupportedTypeError{v.Type()}
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

// Required check value required when anotherField str is a member of the set of strings params
func Required(v reflect.Value) bool {
	return !Empty(v)
}

// RequiredIf check value required when anotherField str is a member of the set of strings params
func RequiredIf(v reflect.Value, anotherField reflect.Value, params []string, tag *ValidTag) bool {
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
		if InString(value, params...) {
			if Empty(v) {
				if tag != nil {
					if tag.messageParameter == nil {
						tag.messageParameter = make(messageParameterMap)
					}
					tag.messageParameter["value"] = value
				}
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
				panic(fmt.Sprintf("validator: RequiredIf unsupport Type %T", value.Interface()))
			}
		}

		for _, value := range values {
			if InString(value, params...) {
				if Empty(v) {
					if tag != nil {
						if tag.messageParameter == nil {
							tag.messageParameter = make(messageParameterMap)
						}
						tag.messageParameter["value"] = value
					}
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
				panic(fmt.Sprintf("validator: RequiredIf unsupport Type %T", value.Interface()))
			}
		}

		for _, value := range values {
			if InString(value, params...) {
				if Empty(v) {
					if tag != nil {
						if tag.messageParameter == nil {
							tag.messageParameter = make(messageParameterMap)
						}
						tag.messageParameter["value"] = value
					}
					return false
				}
			}
		}
	default:
		panic(fmt.Sprintf("validator: RequiredIf unsupport Type %T", anotherField.Interface()))
	}

	return true
}

// RequiredUnless check value required when anotherField str is a member of the set of strings params
func RequiredUnless(v reflect.Value, anotherField reflect.Value, params ...string) bool {
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
		if !InString(value, params...) {
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
			if !InString(value, params...) {
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
			if !InString(value, params...) {
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

func checkRequired(v reflect.Value, f *field, o reflect.Value, name string, structName string) error {
	for _, tag := range f.requiredTags {
		isError := false
		switch tag.name {
		case "required":
			isError = !Required(v)
		case "requiredIf":
			anotherField, err := findField(tag.params[0], o)
			if err == nil && len(tag.params) >= 2 && !RequiredIf(v, anotherField, tag.params[1:], tag) {
				isError = true
			}
		case "requiredUnless":
			anotherField, err := findField(tag.params[0], o)
			if err == nil && len(tag.params) >= 2 && !RequiredUnless(v, anotherField, tag.params[1:]...) {
				isError = true
			}
		case "requiredWith":
			if !RequiredWith(tag.params, v) {
				isError = true
			}
		case "requiredWithAll":
			if !RequiredWithAll(tag.params, v) {
				isError = true
			}
		case "requiredWithout":
			if !RequiredWithout(tag.params, v) {
				isError = true
			}
		case "requiredWithoutAll":
			if !RequiredWithoutAll(tag.params, v) {
				isError = true
			}
		}

		if isError {
			return &Error{
				Name:       name,
				StructName: structName,
				Err:        formatsMessages(tag, v, f, o),
				Tag:        tag.name,
			}
		}
	}

	return nil
}

// RequiredWith The field under validation must be present and not empty only if any of the other specified fields are present.
func RequiredWith(otherFields []string, v reflect.Value) bool {
	if !allFailingRequired(otherFields, v) {
		return Required(v)
	}
	return true
}

// RequiredWithAll The field under validation must be present and not empty only if all of the other specified fields are present.
func RequiredWithAll(otherFields []string, v reflect.Value) bool {
	if !anyFailingRequired(otherFields, v) {
		return Required(v)
	}
	return true
}

// RequiredWithout The field under validation must be present and not empty only when any of the other specified fields are not present.
func RequiredWithout(otherFields []string, v reflect.Value) bool {
	if anyFailingRequired(otherFields, v) {
		return Required(v)
	}
	return true
}

// RequiredWithoutAll The field under validation must be present and not empty only when all of the other specified fields are not present.
func RequiredWithoutAll(otherFields []string, v reflect.Value) bool {
	if allFailingRequired(otherFields, v) {
		return Required(v)
	}
	return true
}

func formatsMessages(validTag *ValidTag, v reflect.Value, f *field, o reflect.Value) error {
	validator := newValidator()

	var message string
	var ok bool

	if message, ok = validator.CustomMessage[f.structName+"."+validTag.messageName]; ok {
		return fmt.Errorf(message)
	}

	if validator.Translator != nil {
		message = validator.Translator.Trans(f.structName, validTag.messageName, f.attribute)
		message = replaceAttributes(message, "", validTag.messageParameter)
	} else {
		message, ok = MessageMap[validTag.messageName]
		if ok {
			attribute := f.attribute
			if customAttribute, ok := validator.Attributes[f.structName]; ok {
				attribute = customAttribute
			}
			message = replaceAttributes(message, attribute, validTag.messageParameter)
		}
	}

	if message != "" {
		if shouldReplaceRequiredWith(validTag.name) {
			message = replaceRequiredWith(message, validTag.params, validator)
		}

		if shouldReplaceRequiredIf(validTag.name) {
			message = replaceRequiredIf(message, o, validTag.params[0], validator)
		}

		if validTag.name == "same" {
			message = replaceSame(message, o, validTag.params[0], validator)
		}

		return fmt.Errorf(message)
	}

	return fmt.Errorf("validator: undefined message : %s", validTag.messageName)
}

func replaceAttributes(message string, attribute string, messageParameter messageParameterMap) string {
	message = strings.Replace(message, ":attribute", attribute, -1)
	for key, value := range messageParameter {
		message = strings.Replace(message, ":"+key, value, -1)
	}
	return message
}

func replaceRequiredWith(message string, attributes []string, validator *Validator) string {
	first := true
	var buff bytes.Buffer
	for _, v := range attributes {
		if first {
			first = false
		} else {
			buff.WriteByte(' ')
			buff.WriteByte('/')
			buff.WriteByte(' ')
		}

		if validator.Translator != nil {
			if customAttribute, ok := validator.Translator.attributes[validator.Translator.locale][v]; ok {
				buff.WriteString(customAttribute)
				continue
			}
		}

		if customAttribute, ok := validator.Attributes[v]; ok {
			buff.WriteString(customAttribute)
			continue
		}

		buff.WriteString(v)
	}

	return strings.Replace(message, ":values", buff.String(), -1)
}

func shouldReplaceRequiredWith(tag string) bool {
	switch tag {
	case "requiredWith", "requiredWithAll", "requiredWithout", "requiredWithoutAll":
		return true
	default:
		return false
	}
}

func getDisplayableAttribute(o reflect.Value, attribute string, validator *Validator) string {
	attributes := strings.Split(attribute, ".")
	if len(attributes) > 0 {
		attribute = o.Type().Name() + attributes[0]
	} else {
		attribute = strings.Join(attributes[len(attributes)-2:], ".")
	}

	if validator.Translator != nil {
		if customAttribute, ok := validator.Translator.attributes[validator.Translator.locale][attribute]; ok {
			return customAttribute
		}
	}

	if customAttribute, ok := validator.Attributes[attribute]; ok {
		return customAttribute
	}

	return attributes[len(attributes)-1]
}

func replaceSame(message string, o reflect.Value, attribute string, validator *Validator) string {
	other := getDisplayableAttribute(o, attribute, validator)
	return strings.Replace(message, ":other", other, -1)
}

func replaceRequiredIf(message string, o reflect.Value, attribute string, validator *Validator) string {
	other := getDisplayableAttribute(o, attribute, validator)
	return strings.Replace(message, ":other", other, -1)
}

func shouldReplaceRequiredIf(tag string) bool {
	switch tag {
	case "requiredIf", "requiredUnless":
		return true
	default:
		return false
	}
}

func inArray(needle []string, haystack []string) bool {
	for _, n := range needle {
		for _, s := range haystack {
			if n == s {
				return true
			}
		}
	}

	return false
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

func checkDependentRules(validTag *ValidTag, f *field, v reflect.Value, o reflect.Value, name string, structName string) error {
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
		isValid = Gt(v, anotherField)
	case "gte":
		isValid = Gte(v, anotherField)
	case "lt":
		isValid = Lt(v, anotherField)
	case "lte":
		isValid = Lte(v, anotherField)
	case "same":
		isValid = Same(v, anotherField)
	}

	if !isValid {
		return &Error{
			Name:       name,
			StructName: structName,
			Err:        formatsMessages(validTag, v, f, o),
			Tag:        validTag.name,
		}
	}

	return nil
}
