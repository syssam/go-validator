package validator

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"sync"
)

const tagName string = "valid"

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

// IsRequiredIf check value required when anotherfield str is a member of the set of strings params
func IsRequiredIf(v reflect.Value, anotherfield reflect.Value, params []string, tag *validTag) bool {
	if !anotherfield.IsValid() {
		return false
	}

	switch anotherfield.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.String:

		if IsIn(ToString(anotherfield), params...) {
			if isEmptyValue(v) {
				if tag != nil {
					if tag.messageParameter == nil {
						tag.messageParameter = make(messageParameterMap)
					}
					tag.messageParameter["value"] = anotherfield.String()
				}
				return true
			}
		}

		return false
	}

	return false
}

// IsIn check if string str is a member of the set of strings params
func IsIn(str string, params ...string) bool {
	for _, param := range params {
		if str == param {
			return true
		}
	}

	return false
}

// IsEmail check if the string is an email.
func IsEmail(str string) bool {
	// TODO uppercase letters are not supported
	return rxEmail.MatchString(str)
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

// IsNumeric check if the string contains only numbers. Empty string is valid.
func IsNumeric(str string) bool {
	if IsNull(str) {
		return true
	}
	return rxNumeric.MatchString(str)
}

// IsInt check if the string is an integer. Empty string is valid.
func IsInt(str string) bool {
	if IsNull(str) {
		return true
	}
	return rxInt.MatchString(str)
}

// IsFloat check if the string is a float.
func IsFloat(str string) bool {
	return str != "" && rxFloat.MatchString(str)
}

// IsNull check if the string is null.
func IsNull(str string) bool {
	return len(str) == 0
}

// Max is the validation function for validating if the current field's value is less than or equal to the param's value.
func Max(v reflect.Value, param ...string) bool {
	return IsLte(v, param[0])
}

// Min is the validation function for validating if the current field's value is greater than or equal to the param's value.
func Min(v reflect.Value, param ...string) bool {
	return IsGte(v, param[0])
}

// IsLt is the validation function for validating if the current field's value is less than the param's value.
func IsLt(v reflect.Value, param ...string) bool {
	switch v.Kind() {
	case reflect.String:
		p, _ := ToInt(param[0])
		return IsLtString(v.String(), p)
	case reflect.Slice, reflect.Map, reflect.Array:
		p, _ := ToInt(param[0])
		return int64(v.Len()) < p
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p, _ := ToInt(param[0])
		return IsLtInt64(v.Int(), p)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p, _ := ToUint(param[0])
		return IsLtUnit64(v.Uint(), p)
	case reflect.Float32, reflect.Float64:
		p, _ := ToFloat(param[0])
		return IsLtFloat64(v.Float(), p)
	}

	panic(fmt.Sprintf("validator: IsLt unsupport Type %T", v.Interface()))
}

// IsLte is the validation function for validating if the current field's value is less than or equal to the param's value.
func IsLte(v reflect.Value, param ...string) bool {
	switch v.Kind() {
	case reflect.String:
		p, _ := ToInt(param[0])
		return IsLteString(v.String(), p)
	case reflect.Slice, reflect.Map, reflect.Array:
		p, _ := ToInt(param[0])
		return int64(v.Len()) <= p
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p, _ := ToInt(param[0])
		return IsLteInt64(v.Int(), p)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p, _ := ToUint(param[0])
		return IsLteUnit64(v.Uint(), p)
	case reflect.Float32, reflect.Float64:
		p, _ := ToFloat(param[0])
		return IsLteFloat64(v.Float(), p)
	}

	panic(fmt.Sprintf("validator: IsLte unsupport Type %T", v.Interface()))
}

// IsGte is the validation function for validating if the current field's value is greater than or equal to the param's value.
func IsGte(v reflect.Value, param ...string) bool {
	switch v.Kind() {
	case reflect.String:
		p, _ := ToInt(param[0])
		return IsGteString(v.String(), p)
	case reflect.Slice, reflect.Map, reflect.Array:
		p, _ := ToInt(param[0])
		return int64(v.Len()) >= p
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p, _ := ToInt(param[0])
		return IsGteInt64(v.Int(), p)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p, _ := ToUint(param[0])
		return IsGteUnit64(v.Uint(), p)
	case reflect.Float32, reflect.Float64:
		p, _ := ToFloat(param[0])
		return IsGteFloat64(v.Float(), p)
	}

	panic(fmt.Sprintf("validator: IsGte unsupport Type %T", v.Interface()))
}

// IsGt is the validation function for validating if the current field's value is greater than to the param's value.
func IsGt(v reflect.Value, param ...string) bool {
	switch v.Kind() {
	case reflect.String:
		p, _ := ToInt(param[0])
		return IsGtString(v.String(), p)
	case reflect.Slice, reflect.Map, reflect.Array:
		p, _ := ToInt(param[0])
		return int64(v.Len()) > p
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		p, _ := ToInt(param[0])
		return IsGtInt64(v.Int(), p)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		p, _ := ToUint(param[0])
		return IsGtUnit64(v.Uint(), p)
	case reflect.Float32, reflect.Float64:
		p, _ := ToFloat(param[0])
		return IsGtFloat64(v.Float(), p)
	}

	panic(fmt.Sprintf("validator: IsGt unsupport Type %T", v.Interface()))
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
	if !v.IsValid() {
		return nil
	}

	name := string(append(jsonNamespace, f.nameBytes...))
	structName := string(append(structNamespace, f.structName...))

	if err := checkRequired(v, f, o, name, structName); err != nil {
		return err
	}

	switch v.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.String:

		for _, tag := range f.validTags {

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
			if v.MapIndex(k).Kind() != reflect.Struct {
				err = newTypeValidator(v.MapIndex(k), f, o, jsonNamespace, structNamespace)
				if err != nil {
					return err
				}
			} else {
				jsonNamespace = append(append(jsonNamespace, f.nameBytes...), '.')
				structNamespace = append(append(structNamespace, f.structNameBytes...), '.')
				err = validateStruct(v.MapIndex(k).Interface(), jsonNamespace, structNamespace)
				if err != nil {
					return err
				}
			}
		}
		return nil
	case reflect.Slice, reflect.Array:
		for _, tag := range f.validTags {
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
			if v.Index(i).Kind() != reflect.Struct {
				err = newTypeValidator(v.Index(i), f, o, jsonNamespace, structNamespace)
				if err != nil {
					return err
				}
			} else {
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

func isEmptyValue(v reflect.Value) bool {
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

// IsRequiredUnless check value required when anotherfield str is a member of the set of strings params
func IsRequiredUnless(v reflect.Value, anotherfield reflect.Value, params ...string) bool {
	switch anotherfield.Kind() {
	case reflect.Bool,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr,
		reflect.Float32, reflect.Float64,
		reflect.String:

		if !IsIn(anotherfield.String(), params...) {
			return isEmptyValue(v)
		}

		return false
	}

	return false
}

func allFailingRequired(parameters []string, v reflect.Value) bool {
	for _, p := range parameters {
		if isEmptyValue(v.FieldByName(p)) {
			return false
		}
	}
	return false
}

func anyFailingRequired(parameters []string, v reflect.Value) bool {
	for _, p := range parameters {
		if isEmptyValue(v.FieldByName(p)) {
			return true
		}
	}
	return false
}

func checkRequired(v reflect.Value, f *field, o reflect.Value, name string, structName string) error {
	for _, tag := range f.requiredTags {
		result := false
		switch tag.name {
		case "required":
			result = isEmptyValue(v)
		case "requiredIf":
			anotherfield := findField(tag.params[0], o)
			if len(tag.params) >= 2 && IsRequiredIf(v, anotherfield, tag.params[1:], tag) {
				result = true
			}
		case "requiredUnless":
			if len(tag.params) >= 2 && IsRequiredUnless(v, o.FieldByName(tag.params[0]), tag.params[1:]...) {
				result = true
			}
		case "requiredWith":
			if IsRequiredWith(tag.params, v) {
				result = true
			}
		case "requiredWithAll":
			if IsRequiredWithAll(tag.params, v) {
				result = true
			}
		case "requiredWithout":
			if IsRequiredWithout(tag.params, v) {
				result = true
			}
		case "requiredWithoutAll":
			if IsRequiredWithoutAll(tag.params, v) {
				result = true
			}
		}

		if result {
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

// IsRequiredWith The field under validation must be present and not empty only if any of the other specified fields are present.
func IsRequiredWith(otherFields []string, v reflect.Value) bool {
	if !allFailingRequired(otherFields, v) {
		return isEmptyValue(v)
	}
	return false
}

// IsRequiredWithAll The field under validation must be present and not empty only if all of the other specified fields are present.
func IsRequiredWithAll(otherFields []string, v reflect.Value) bool {
	if !anyFailingRequired(otherFields, v) {
		return isEmptyValue(v)
	}
	return false
}

// IsRequiredWithout The field under validation must be present and not empty only when any of the other specified fields are not present.
func IsRequiredWithout(otherFields []string, v reflect.Value) bool {
	if anyFailingRequired(otherFields, v) {
		return isEmptyValue(v)
	}
	return false
}

// IsRequiredWithoutAll The field under validation must be present and not empty only when all of the other specified fields are not present.
func IsRequiredWithoutAll(otherFields []string, v reflect.Value) bool {
	if allFailingRequired(otherFields, v) {
		return isEmptyValue(v)
	}
	return false
}

func formatsMessages(validTag *validTag, v reflect.Value, f *field, o reflect.Value) error {
	validator := newValidator()
	if validator.Translator != nil {
		message := validator.Translator.Trans(f.structName, validTag.messageName, f.attribute)
		message = replaceAttributes(message, "", validTag.messageParameter)
		if shouldReplaceRequiredWith(validTag.name) {
			message = replaceRequiredWith(message, validTag.params, validator)
		}
		if shouldReplaceRequiredIf(validTag.name) {
			message = replaceRequiredIf(message, o.Type().Name()+"."+validTag.params[0], validator)
		}
		return fmt.Errorf(message)
	}

	if m, ok := validator.CustomMessage[f.structName+"."+validTag.messageName]; ok {
		return fmt.Errorf(m)
	}

	message, ok := MessageMap[validTag.messageName]
	if ok {
		attribute := f.attribute
		if customAttribute, ok := validator.Attributes[f.structName]; ok {
			attribute = customAttribute
		}

		message = replaceAttributes(message, attribute, validTag.messageParameter)

		if shouldReplaceRequiredWith(validTag.name) {
			message = replaceRequiredWith(message, validTag.params, validator)
		}

		if shouldReplaceRequiredIf(validTag.name) {
			message = replaceRequiredIf(message, o.Type().Name()+"."+validTag.params[0], validator)
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

func replaceRequiredIf(message string, attribute string, validator *Validator) string {
	if validator.Translator != nil {
		if customAttribute, ok := validator.Translator.attributes[validator.Translator.locale][attribute]; ok {
			return strings.Replace(message, ":other", customAttribute, -1)
		}
	}

	if customAttribute, ok := validator.Attributes[attribute]; ok {
		return strings.Replace(message, ":other", customAttribute, -1)
	}

	return strings.Replace(message, ":other", attribute, -1)
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

func findField(fieldName string, v reflect.Value) reflect.Value {
	fields := strings.Split(fieldName, ".")
	current := v.FieldByName(fields[0])

	i := 2
	if len(fields) > i {
		for true {
			if current.Kind() == reflect.Interface || current.Kind() == reflect.Ptr {
				current = current.Elem()
			}
			name := fields[i-1]
			current = current.FieldByName(name)
			if i == len(fields) {
				break
			}
			i++
		}
	}

	return current
}
