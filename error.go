package validator

import (
	"bytes"
)

// Errors is an array of multiple errors and conforms to the error interface.
type Errors []error

// Errors returns itself.
func (es Errors) Errors() []error {
	return es
}

func (es Errors) Error() string {
	var buff bytes.Buffer
	first := true
	for _, e := range es {
		if first {
			first = false
		} else {
			buff.WriteByte('\n')
		}
		buff.WriteString(e.Error())
	}
	return buff.String()
}

// JSON output json format.
func (es Errors) JSON() string {
	var buff bytes.Buffer
	first := true
	buff.WriteByte('{')
	for _, e := range es {
		if first {
			first = false
		} else {
			buff.WriteByte(',')
		}
		buff.WriteString(e.(Error).Name)
		buff.WriteByte(':')
		buff.WriteString(e.Error())
	}
	buff.WriteByte('}')
	return buff.String()
}

// Error encapsulates a name, an error and whether there's a custom error message or not.
type Error struct {
	Name       string
	StructName string
	Err        error

	// Tag indicates the name of the validator that failed
	Tag string
}

func (e Error) Error() string {
	return e.Err.Error()
}
