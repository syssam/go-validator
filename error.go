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

// MarshalJSON output Json format.
func (es Errors) MarshalJSON() ([]byte, error) {
	var buff bytes.Buffer
	first := true
	buff.WriteByte('[')
	for _, e := range es {
		if first {
			first = false
		} else {
			buff.WriteByte(',')
		}
		buff.WriteByte('{')
		buff.WriteString(`"message":`)
		buff.WriteByte('"')
		buff.WriteString(e.Error())
		buff.WriteByte('"')
		buff.WriteByte(',')
		buff.WriteString(`"parameter":`)
		buff.WriteByte('"')
		buff.WriteString(e.(*FieldError).name)
		buff.WriteByte('"')
		buff.WriteByte('}')
	}
	buff.WriteByte(']')
	return buff.Bytes(), nil
}

// FieldError encapsulates name, message, and value etc.
type FieldError struct {
	name              string
	structName        string
	tag               string // Tag indicates the name of the validator that failed
	messageName       string
	messageParameters messageParameters
	attribute         string
	value             string
	err               error
}

func (fe *FieldError) Error() string {
	return fe.err.Error()
}
