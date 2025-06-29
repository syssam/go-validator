package validator

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
)

// Errors is an array of multiple errors and conforms to the error interface.
type Errors []error

// Error implements the error interface
func (es Errors) Error() string {
	if len(es) == 0 {
		return ""
	}
	if len(es) == 1 {
		return es[0].Error()
	}

	var builder strings.Builder
	builder.Grow(len(es) * 50) // Pre-allocate estimated capacity

	for i, e := range es {
		if i > 0 {
			builder.WriteByte('\n')
		}
		builder.WriteString(e.Error())
	}
	return builder.String()
}

// Errors returns itself for compatibility
func (es Errors) Errors() []error {
	return es
}

// FieldErrors returns all FieldError instances
func (es Errors) FieldErrors() []*FieldError {
	fieldErrors := make([]*FieldError, 0, len(es))
	for _, e := range es {
		if fieldErr, ok := e.(*FieldError); ok {
			fieldErrors = append(fieldErrors, fieldErr)
		} else {
			// Convert generic error to FieldError
			fieldErrors = append(fieldErrors, &FieldError{
				Message: e.Error(),
			})
		}
	}
	return fieldErrors
}

// HasFieldError checks if there's an error for the specified field
func (es Errors) HasFieldError(fieldName string) bool {
	for _, e := range es {
		if fieldErr, ok := e.(*FieldError); ok && fieldErr.Name == fieldName {
			return true
		}
	}
	return false
}

// GetFieldError returns the first error for the specified field
func (es Errors) GetFieldError(fieldName string) *FieldError {
	for _, e := range es {
		if fieldErr, ok := e.(*FieldError); ok && fieldErr.Name == fieldName {
			return fieldErr
		}
	}
	return nil
}

// GroupByField groups errors by field name
func (es Errors) GroupByField() map[string][]*FieldError {
	groups := make(map[string][]*FieldError)
	for _, e := range es {
		if fieldErr, ok := e.(*FieldError); ok {
			groups[fieldErr.Name] = append(groups[fieldErr.Name], fieldErr)
		}
	}
	return groups
}

type ErrorResponse struct {
	Message   string `json:"message"`
	Parameter string `json:"parameter"`
}

var errorResponsePool = sync.Pool{
	New: func() interface{} {
		slice := make([]ErrorResponse, 0, 10)
		return &slice
	},
}

// MarshalJSON output Json format.
func (es Errors) MarshalJSON() ([]byte, error) {
	if len(es) == 0 {
		return []byte("[]"), nil
	}

	responsesPtr := errorResponsePool.Get().(*[]ErrorResponse)
	responses := (*responsesPtr)[:0]

	defer errorResponsePool.Put(responsesPtr)

	for _, e := range es {
		if fieldErr, ok := e.(*FieldError); ok {
			responses = append(responses, ErrorResponse{
				Message:   fieldErr.Message,
				Parameter: fieldErr.Name,
			})
		}
	}

	*responsesPtr = responses
	return json.Marshal(responses)
}

// FieldError encapsulates name, message, and value etc.
type FieldError struct {
	Name              string            `json:"name"`
	StructName        string            `json:"struct_name,omitempty"`
	Tag               string            `json:"tag"`
	MessageName       string            `json:"message_name,omitempty"`
	MessageParameters MessageParameters `json:"message_parameters,omitempty"`
	Attribute         string            `json:"attribute,omitempty"`
	DefaultAttribute  string            `json:"default_attribute,omitempty"`
	Value             string            `json:"value,omitempty"`
	Message           string            `json:"message"`
	FuncError         error             `json:"func_error,omitempty"`
}

// Unwrap implements the errors.Unwrap interface for error chain support
func (fe *FieldError) Unwrap() error {
	return fe.FuncError
}

// Error returns the error message with optional function error details
func (fe *FieldError) Error() string {
	if fe.Message != "" {
		return fe.Message
	}
	if fe.FuncError != nil {
		return fmt.Sprintf("validation failed for field '%s': %v", fe.Name, fe.FuncError)
	}
	return fmt.Sprintf("validation failed for field '%s'", fe.Name)
}

// HasFuncError checks if there's an underlying function error
func (fe *FieldError) HasFuncError() bool {
	return fe.FuncError != nil
}

// SetMessage sets the user-friendly message while preserving function error
func (fe *FieldError) SetMessage(msg string) {
	fe.Message = msg
}
