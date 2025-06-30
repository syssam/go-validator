package validator

import (
	"bytes"
	"errors"
	"reflect"
	"strings"
	"sync"
)

// A field represents a single field found in a struct.
type field struct {
	name             string
	nameBytes        []byte // []byte(name)
	structName       string
	structNameBytes  []byte // []byte(structName)
	attribute        string
	defaultAttribute string
	tag              bool
	index            []int
	requiredTags     requiredTags
	validTags        otherValidTags
	typ              reflect.Type
	omitEmpty        bool
}

// A ValidTag represents parse validTag into field struct.
type ValidTag struct {
	name              string
	params            []string
	messageName       string
	messageParameters MessageParameters
}

// A otherValidTags represents parse validTag into field struct when validTag is not required...
type otherValidTags []*ValidTag

// A requiredTags represents parse validTag into field struct when validTag is required...
type requiredTags []*ValidTag

var fieldCache sync.Map // map[reflect.Type][]field

// cachedTypefields is like typefields but uses a cache to avoid repeated work.
func cachedTypefields(t reflect.Type) []field {
	if f, ok := fieldCache.Load(t); ok {
		if fields, ok := f.([]field); ok {
			return fields
		}
	}
	f, _ := fieldCache.LoadOrStore(t, typefields(t))
	if fields, ok := f.([]field); ok {
		return fields
	}
	return []field{}
}

// shouldSkipField determines if a field should be skipped based on export status
func shouldSkipField(sf reflect.StructField) bool {
	isUnexported := sf.PkgPath != ""
	if sf.Anonymous {
		t := sf.Type
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		if isUnexported && t.Kind() != reflect.Struct {
			// Ignore embedded fields of unexported non-struct types.
			return true
		}
		// Do not ignore embedded fields of unexported struct types
		// since they may have exported fields.
	} else if isUnexported {
		// Ignore unexported non-embedded fields.
		return true
	}
	return false
}

// getFieldName extracts the field name from json tag or struct field name
func getFieldName(sf reflect.StructField, f *field) string {
	name := sf.Tag.Get("json")
	if !f.isvalidTag(name) {
		name = ""
	}
	if name == "" {
		name = sf.Name
	}
	return name
}

// createFieldFromStructField creates a field struct from reflect.StructField
func createFieldFromStructField(sf reflect.StructField, f *field, t, ft reflect.Type, index []int, validTag string) field {
	name := getFieldName(sf, f)
	tagged := sf.Tag.Get("json") != "" && f.isvalidTag(sf.Tag.Get("json"))
	requiredTags, otherValidTags, defaultAttribute := f.parseTagIntoSlice(validTag, ft)

	return field{
		name:             name,
		nameBytes:        []byte(name),
		structName:       t.Name() + "." + sf.Name,
		structNameBytes:  []byte(t.Name() + "." + sf.Name),
		attribute:        sf.Name,
		defaultAttribute: defaultAttribute,
		tag:              tagged,
		index:            index,
		requiredTags:     requiredTags,
		validTags:        otherValidTags,
		typ:              ft,
		omitEmpty:        strings.Contains(validTag, "omitempty"),
	}
}

// processStructField processes a single struct field and updates fields/next accordingly
func processStructField(sf reflect.StructField, f *field, t reflect.Type, i int, count, nextCount map[reflect.Type]int, fields, next *[]field) {
	if shouldSkipField(sf) {
		return
	}

	validTag := sf.Tag.Get(tagName)
	if validTag == "-" {
		return
	}

	index := make([]int, len(f.index)+1)
	copy(index, f.index)
	index[len(f.index)] = i

	ft := sf.Type
	if validTag == "" && ft.Kind() != reflect.Slice && ft.Kind() != reflect.Array {
		return
	}

	if ft.Name() == "" && ft.Kind() == reflect.Ptr {
		// Follow pointer.
		ft = ft.Elem()
	}

	name := getFieldName(sf, f)

	// Record found field and index sequence.
	if name != sf.Name || !sf.Anonymous || ft.Kind() != reflect.Struct {
		count[f.typ]++
		newField := createFieldFromStructField(sf, f, t, ft, index, validTag)
		*fields = append(*fields, newField)

		if count[f.typ] > 1 {
			// If there were multiple instances, add a second,
			// so that the annihilation code will see a duplicate.
			// It only cares about the distinction between 1 or 2,
			// so don't bother generating any more copies.
			*fields = append(*fields, (*fields)[len(*fields)-1])
		}
		return
	}

	// Record new anonymous struct to explore in next round.
	nextCount[ft]++
	if nextCount[ft] == 1 {
		newField := createFieldFromStructField(sf, f, t, ft, index, validTag)
		*next = append(*next, newField)
	}
}

// typefields returns a list of fields that Validator should recognize for the given type.
// The algorithm is breadth-first search over the set of structs to include - the top struct
// and then any reachable anonymous structs.
func typefields(t reflect.Type) []field {
	current := make([]field, 0, t.NumField())
	next := []field{{typ: t}}

	// Count of queued names for current level and the next.
	nextCount := map[reflect.Type]int{}

	// Types already visited at an earlier level.
	visited := map[reflect.Type]bool{}

	var fields []field

	for len(next) > 0 {
		current, next = next, current[:0]
		count, nextCount := nextCount, map[reflect.Type]int{}

		for _, f := range current {
			if visited[f.typ] {
				continue
			}
			visited[f.typ] = true
			for i := 0; i < f.typ.NumField(); i++ {
				sf := f.typ.Field(i)
				processStructField(sf, &f, t, i, count, nextCount, &fields, &next)
			}
		}
	}

	return fields
}

func (f *field) parseTagIntoSlice(tag string, ft reflect.Type) (requiredTags, otherValidTags, string) {
	options := strings.Split(tag, ",")
	var otherValidTags otherValidTags
	var requiredTags requiredTags
	defaultAttribute := ""

	for _, option := range options {
		option = strings.TrimSpace(option)

		tag := strings.Split(option, "=")
		var params []string

		if len(tag) == 2 {
			params = strings.Split(tag[1], "|")
		}

		switch tag[0] {
		case "attribute":
			if len(tag) == 2 {
				defaultAttribute = tag[1]
			}
			continue
		case "required", "requiredIf", "requiredUnless", "requiredWith", "requiredWithAll", "requiredWithout", "requiredWithoutAll":
			messageParameters, _ := f.parseMessageParameterIntoSlice(tag[0], params...)
			requiredTags = append(requiredTags, &ValidTag{
				name:              tag[0],
				params:            params,
				messageName:       f.parseMessageName(tag[0], ft),
				messageParameters: messageParameters,
			})
			continue
		}

		messageParameters, _ := f.parseMessageParameterIntoSlice(tag[0], params...)
		otherValidTags = append(otherValidTags, &ValidTag{
			name:              tag[0],
			params:            params,
			messageName:       f.parseMessageName(tag[0], ft),
			messageParameters: messageParameters,
		})
	}

	return requiredTags, otherValidTags, defaultAttribute
}

func (f *field) isvalidTag(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if strings.ContainsRune("\\'\"!#$%&()*+-./:<=>?@[]^_{|}~ ", c) {
			// Backslash and quote chars are reserved, but
			// otherwise anything goes.
			return false
		}
	}
	return true
}

//nolint:unused // Kept for potential future use
func (f *field) isValidAttribute(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		if strings.ContainsRune("\\'\"!#$%&()*+-./:<=>?@[]^_{|}~ ", c) {
			// Backslash and quote chars are reserved, but
			// otherwise anything goes.
			return false
		}
	}
	return true
}

func (f *field) parseMessageName(rule string, ft reflect.Type) string {
	messageName := rule

	switch rule {
	case "between", "gt", "gte", "lt", "lte", "min", "max", "size":
		switch ft.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16,
			reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16,
			reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64:
			return messageName + ".numeric"
		case reflect.String:
			return messageName + ".string"
		case reflect.Array, reflect.Slice, reflect.Map:
			return messageName + ".array"
		case reflect.Struct, reflect.Ptr:
			return messageName
		default:
			return messageName
		}
	default:
		return messageName
	}
}

type messageParameter struct {
	Key   string
	Value string
}

// A MessageParameters represents store message parameter into field struct.
type MessageParameters []messageParameter

func (f *field) parseMessageParameterIntoSlice(rule string, params ...string) (MessageParameters, error) {
	var messageParameters MessageParameters

	switch rule {
	case "requiredUnless":
		if len(params) < 2 {
			return nil, errors.New("validator: " + rule + " format is not valid")
		}

		first := true
		var buff bytes.Buffer
		for _, v := range params[1:] {
			if first {
				first = false
			} else {
				buff.WriteByte(' ')
				buff.WriteByte(',')
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
	case "between", "digitsBetween":
		if len(params) != 2 {
			return nil, errors.New("validator: " + rule + " format is not valid")
		}

		messageParameters = append(
			messageParameters,
			messageParameter{
				Key:   "Min",
				Value: params[0],
			}, messageParameter{
				Key:   "Max",
				Value: params[1],
			},
		)
	case "gt", "gte", "lt", "lte":
		if len(params) != 1 {
			return nil, errors.New("validator: " + rule + " format is not valid")
		}

		messageParameters = append(
			messageParameters,
			messageParameter{
				Key:   "Value",
				Value: params[0],
			},
		)
	case "max":
		if len(params) != 1 {
			return nil, errors.New("validator: " + rule + " format is not valid")
		}

		messageParameters = append(
			messageParameters,
			messageParameter{
				Key:   "Max",
				Value: params[0],
			},
		)
	case "min":
		if len(params) != 1 {
			return nil, errors.New("validator: " + rule + " format is not valid")
		}

		messageParameters = append(
			messageParameters,
			messageParameter{
				Key:   "Min",
				Value: params[0],
			},
		)
	case "size":
		if len(params) != 1 {
			return nil, errors.New("validator: " + rule + " format is not valid")
		}
		messageParameters = append(
			messageParameters,
			messageParameter{
				Key:   "Size",
				Value: params[0],
			},
		)
	}

	if len(messageParameters) > 0 {
		return messageParameters, nil
	}

	return nil, nil
}
