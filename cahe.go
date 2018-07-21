package validator

import (
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
}

// A validTag represents parse validTag into field struct.
type validTag struct {
	name             string
	params           []string
	messageName      string
	messageParameter messageParameterMap
}

// A otherValidTags represents parse validTag into field struct when validTag is not required...
type otherValidTags []*validTag

// A requiredTags represents parse validTag into field struct when validTag is required...
type requiredTags []*validTag

var fieldCache sync.Map // map[reflect.Type][]field

// cachedTypefields is like typefields but uses a cache to avoid repeated work.
func cachedTypefields(t reflect.Type) []field {
	if f, ok := fieldCache.Load(t); ok {
		return f.([]field)
	}
	f, _ := fieldCache.LoadOrStore(t, typefields(t))
	return f.([]field)
}

// typefields returns a list of fields that JSON should recognize for the given type.
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
				if !isvalidTag(name) {
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

					requiredTags, otherValidTags := parseTagIntoArray(validTag, ft)

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
					requiredTags, otherValidTags := parseTagIntoArray(validTag, ft)

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
					})
				}
			}
		}
	}

	return fields
}

func parseTagIntoArray(tag string, ft reflect.Type) (requiredTags, otherValidTags) {
	options := strings.Split(tag, ",")
	var otherValidTags otherValidTags
	var requiredTags requiredTags

	for _, option := range options {
		option = strings.TrimSpace(option)

		if !isvalidTag(option) {
			continue
		}

		tag := strings.Split(option, "=")
		var params []string

		if len(tag) == 2 {
			params = strings.Split(tag[1], "|")
		}

		switch tag[0] {
		case "required", "requiredIf", "requiredUnless", "requiredWith", "requiredWithAll", "requiredWithout", "requiredWithoutAll":
			requiredTags = append(requiredTags, &validTag{
				name:             tag[0],
				params:           params,
				messageName:      parseMessageName(tag[0], ft),
				messageParameter: parseMessageParameterIntoMap(tag[0], params...),
			})
			continue
		}

		otherValidTags = append(otherValidTags, &validTag{
			name:             tag[0],
			params:           params,
			messageName:      parseMessageName(tag[0], ft),
			messageParameter: parseMessageParameterIntoMap(tag[0], params...),
		})
	}
	return requiredTags, otherValidTags
}

func isvalidTag(s string) bool {
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

func parseMessageName(rule string, ft reflect.Type) string {
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

type messageParameterMap map[string]string

func parseMessageParameterIntoMap(rule string, params ...string) messageParameterMap {
	switch rule {
	case "requiredUnless":
		if len(params) != 1 {
			panic(fmt.Sprintf("validator: " + rule + " format is not valid"))
		}
		return messageParameterMap{
			"values": params[0],
		}
	case "between", "digitsBetween":
		if len(params) != 2 {
			panic(fmt.Sprintf("validator: " + rule + " format is not valid"))
		}
		return messageParameterMap{
			"min": params[0],
			"max": params[1],
		}
	case "gt", "gte", "lt", "lte":
		if len(params) != 1 {
			panic(fmt.Sprintf("validator: " + rule + " format is not valid"))
		}
		return messageParameterMap{
			"value": params[0],
		}
	case "max":
		if len(params) != 1 {
			panic(fmt.Sprintf("validator: " + rule + " format is not valid"))
		}
		return messageParameterMap{
			"max": params[0],
		}
	case "min":
		if len(params) != 1 {
			panic(fmt.Sprintf("validator: " + rule + " format is not valid"))
		}
		return messageParameterMap{
			"min": params[0],
		}
	}

	return nil
}
