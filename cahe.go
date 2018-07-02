package validator

import (
	"reflect"
	"strings"
	"sync"
	"unicode"
)

// A field represents a single field found in a struct.
type field struct {
	name       string
	structName string
	attribute  string
	tag        bool
	index      []int
	validTags  []*validTag
	typ        reflect.Type
}

// A validTag represents parse validTag into field struct.
type validTag struct {
	name   string
	params []string
}

var fieldCache sync.Map // map[reflect.Type][]field

// cachedTypeFields is like typeFields but uses a cache to avoid repeated work.
func cachedTypeFields(t reflect.Type) []field {
	if f, ok := fieldCache.Load(t); ok {
		return f.([]field)
	}
	f, _ := fieldCache.LoadOrStore(t, typeFields(t))
	return f.([]field)
}

// typeFields returns a list of fields that JSON should recognize for the given type.
// The algorithm is breadth-first search over the set of structs to include - the top struct
// and then any reachable anonymous structs.
func typeFields(t reflect.Type) []field {
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
				if !isValidTag(name) {
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

					fields = append(fields, field{
						name:       name,
						structName: t.Name() + "." + sf.Name,
						attribute:  sf.Name,
						tag:        tagged,
						index:      index,
						validTags:  parseTagIntoArray(validTag),
						typ:        ft,
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
					next = append(next, field{
						name:       sf.Name,
						structName: t.Name() + "." + sf.Name,
						attribute:  sf.Name,
						index:      index,
						validTags:  parseTagIntoArray(validTag),
						typ:        ft,
					})
				}
			}
		}
	}

	return fields
}

func parseTagIntoArray(tag string) []*validTag {
	options := strings.Split(tag, ",")
	var rules []*validTag
	for _, option := range options {
		option = strings.TrimSpace(option)

		if !isValidTag(option) {
			continue
		}

		tag := strings.Split(option, "=")
		var params []string

		if len(tag) == 2 {
			params = strings.Split(tag[1], "|")
		}

		rules = append(rules, &validTag{
			name:   tag[0],
			params: params,
		})
	}
	return rules
}

func isValidTag(s string) bool {
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
