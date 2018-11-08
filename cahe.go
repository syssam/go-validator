package validator

import (
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"unicode"
)

// A field represents a single field found in a struct.
type field struct {
	name            string
	nameBytes       []byte // []byte(name)
	structName      string
	structNameBytes []byte // []byte(structName)
	attribute       string
	tag             bool
	index           []int
	requiredTags    requiredTags
	validTags       otherValidTags
	typ             reflect.Type
	omitEmpty       bool
}

// A ValidTag represents parse validTag into field struct.
type ValidTag struct {
	name              string
	params            []string
	messageName       string
	messageParameters messageParameters
}

// A otherValidTags represents parse validTag into field struct when validTag is not required...
type otherValidTags []*ValidTag

// A requiredTags represents parse validTag into field struct when validTag is required...
type requiredTags []*ValidTag

var fieldCache sync.Map // map[reflect.Type][]field

// cachedTypefields is like typefields but uses a cache to avoid repeated work.
func cachedTypefields(t reflect.Type) []field {
	if f, ok := fieldCache.Load(t); ok {
		return f.([]field)
	}
	f, _ := fieldCache.LoadOrStore(t, typefields(t))
	return f.([]field)
}

// typefields returns a list of fields that Validator should recognize for the given type.
// The algorithm is breadth-first search over the set of structs to include - the top struct
// and then any reachable anonymous structs.
func typefields(t reflect.Type) []field {
	current := []field{}
	next := []field{{typ: t}}

	// Count of queued names for current level and the next.
	count := map[reflect.Type]int{}
	nextCount := map[reflect.Type]int{}

	// Types already visited at an earlier level.
	visited := map[reflect.Type]bool{}

	var fields []field

	for len(next) > 0 {
		current, next = next, current[:0]
		count, nextCount = nextCount, map[reflect.Type]int{}

		for _, f := range current {
			if visited[f.typ] {
				continue
			}
			visited[f.typ] = true
			for i := 0; i < f.typ.NumField(); i++ {
				sf := f.typ.Field(i)
				isUnexported := sf.PkgPath != ""
				if sf.Anonymous {
					t := sf.Type
					if t.Kind() == reflect.Ptr {
						t = t.Elem()
					}
					if isUnexported && t.Kind() != reflect.Struct {
						// Ignore embedded fields of unexported non-struct types.
						continue
					}
					// Do not ignore embedded fields of unexported struct types
					// since they may have exported fields.
				} else if isUnexported {
					// Ignore unexported non-embedded fields.
					continue
				}
				validTag := sf.Tag.Get(tagName)
				name := sf.Tag.Get("json")
				if !f.isvalidTag(name) {
					name = ""
				}
				if validTag == "-" || validTag == "" {
					continue
				}

				index := make([]int, len(f.index)+1)
				copy(index, f.index)
				index[len(f.index)] = i

				ft := sf.Type
				if ft.Name() == "" && ft.Kind() == reflect.Ptr {
					// Follow pointer.
					ft = ft.Elem()
				}

				// Record found field and index sequence.
				if name != "" || !sf.Anonymous || ft.Kind() != reflect.Struct {
					tagged := name != ""
					if name == "" {
						name = sf.Name
					}

					requiredTags, otherValidTags := f.parseTagIntoSlice(validTag, ft)

					fields = append(fields, field{
						name:            name,
						nameBytes:       []byte(name),
						structName:      t.Name() + "." + sf.Name,
						structNameBytes: []byte(t.Name() + "." + sf.Name),
						attribute:       sf.Name,
						tag:             tagged,
						index:           index,
						requiredTags:    requiredTags,
						validTags:       otherValidTags,
						typ:             ft,
						omitEmpty:       strings.Contains(validTag, "omitempty"),
					})

					if count[f.typ] > 1 {
						// If there were multiple instances, add a second,
						// so that the annihilation code will see a duplicate.
						// It only cares about the distinction between 1 or 2,
						// so don't bother generating any more copies.
						fields = append(fields, fields[len(fields)-1])
					}

					continue
				}

				// Record new anonymous struct to explore in next round.
				nextCount[ft]++
				if nextCount[ft] == 1 {
					requiredTags, otherValidTags := f.parseTagIntoSlice(validTag, ft)

					next = append(next, field{
						name:            sf.Name,
						nameBytes:       []byte(sf.Name),
						structName:      t.Name() + "." + sf.Name,
						structNameBytes: []byte(t.Name() + "." + sf.Name),
						attribute:       sf.Name,
						index:           index,
						requiredTags:    requiredTags,
						validTags:       otherValidTags,
						typ:             ft,
						omitEmpty:       strings.Contains(validTag, "omitempty"),
					})
				}
			}
		}
	}

	return fields
}

func (f *field) parseTagIntoSlice(tag string, ft reflect.Type) (requiredTags, otherValidTags) {
	options := strings.Split(tag, ",")
	var otherValidTags otherValidTags
	var requiredTags requiredTags

	for _, option := range options {
		option = strings.TrimSpace(option)

		if !f.isvalidTag(option) {
			continue
		}

		tag := strings.Split(option, "=")
		var params []string

		if len(tag) == 2 {
			params = strings.Split(tag[1], "|")
		}

		switch tag[0] {
		case "required", "requiredIf", "requiredUnless", "requiredWith", "requiredWithAll", "requiredWithout", "requiredWithoutAll":
			requiredTags = append(requiredTags, &ValidTag{
				name:              tag[0],
				params:            params,
				messageName:       f.parseMessageName(tag[0], ft),
				messageParameters: f.parseMessageParameterIntoSlice(tag[0], params...),
			})
			continue
		}

		otherValidTags = append(otherValidTags, &ValidTag{
			name:              tag[0],
			params:            params,
			messageName:       f.parseMessageName(tag[0], ft),
			messageParameters: f.parseMessageParameterIntoSlice(tag[0], params...),
		})
	}
	return requiredTags, otherValidTags
}

func (f *field) isvalidTag(s string) bool {
	if s == "" {
		return false
	}
	for _, c := range s {
		switch {
		case strings.ContainsRune("\\'\"!#$%&()*+-./:<=>?@[]^_{|}~ ", c):
			// Backslash and quote chars are reserved, but
			// otherwise any punctuation chars are allowed
			// in a tag name.
		default:
			if !unicode.IsLetter(c) && !unicode.IsDigit(c) {
				return false
			}
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

type messageParameters []messageParameter

func (f *field) parseMessageParameterIntoSlice(rule string, params ...string) messageParameters {
	var messageParameters messageParameters
	switch rule {
	case "requiredUnless":
		if len(params) < 2 {
			panic(fmt.Sprintf("validator: " + rule + " format is not valid"))
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
			panic(fmt.Sprintf("validator: " + rule + " format is not valid"))
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
			panic(fmt.Sprintf("validator: " + rule + " format is not valid"))
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
			panic(fmt.Sprintf("validator: " + rule + " format is not valid"))
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
			panic(fmt.Sprintf("validator: " + rule + " format is not valid"))
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
			panic(fmt.Sprintf("validator: " + rule + " format is not valid"))
		}
		messageParameters = append(
			messageParameters,
			messageParameter{
				Key:   "Size",
				Value: params[0],
			},
		)
	}

	if messageParameters != nil && len(messageParameters) > 0 {
		return messageParameters
	}

	return nil
}
