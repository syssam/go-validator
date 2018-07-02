package validator

import (
	"bytes"
	"strings"
)

// Errors is an array of multiple errors and conforms to the error interface.
type Errors []error

// Errors returns itself.
func (es Errors) Errors() []error {
	return es
}

func (es Errors) Error() string {
	buff := bytes.NewBufferString("")

	for _, e := range es {
		buff.WriteString(e.Error())
		buff.WriteString("\n")
	}

	return strings.TrimSpace(buff.String())
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
	return e.Name + ": " + e.Err.Error()
}
