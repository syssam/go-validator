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
	messageParameters MessageParameters
	attribute         string
	value             string
	err               error
}

// Name returns the FieldError.name
func (fe *FieldError) Name() string {
	return fe.name
}

// StructName returns the FieldError.structName
func (fe *FieldError) StructName() string {
	return fe.structName
}

// Tag returns the FieldError.tag
func (fe *FieldError) Tag() string {
	return fe.tag
}

// MessageName returns the FieldError.messageName
func (fe *FieldError) MessageName() string {
	return fe.messageName
}

// MessageParameters returns the FieldError.messageParameters
func (fe *FieldError) MessageParameters() MessageParameters {
	return fe.messageParameters
}

// Attribute returns the FieldError.attribute
func (fe *FieldError) Attribute() string {
	return fe.attribute
}

// Value returns the FieldError.value, which is validate value
func (fe *FieldError) Value() string {
	return fe.value
}

// Error returns the FieldError.err as string
func (fe *FieldError) Error() string {
	return fe.err.Error()
}
