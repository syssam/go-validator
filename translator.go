package validator

import (
	"strings"
)

// Translate translate type
type Translate map[string]string

// Translator type
type Translator struct {
	customMessage map[string]Translate
	messages      map[string]Translate
	attributes    map[string]Translate
}

// NewTranslator returns a new instance of 'translator' with sane defaults.
func NewTranslator() *Translator {
	translator := &Translator{
		messages:      make(map[string]Translate),
		attributes:    make(map[string]Translate),
		customMessage: make(map[string]Translate),
	}
	return translator
}

// SetMessage set Message
func (t *Translator) SetMessage(langCode string, messages Translate) {
	t.messages[langCode] = messages
}

// LoadMessage load message
func (t *Translator) LoadMessage(langCode string) Translate {
	return t.messages[langCode]
}

// SetAttributes set attributes
func (t *Translator) SetAttributes(langCode string, messages Translate) {
	t.attributes[langCode] = messages
}

// Trans translate errors
func (t *Translator) Trans(errors Errors, language string) Errors {
	for i := 0; i < len(errors); i++ {
		fieldError, ok := errors[i].(*FieldError)
		if !ok {
			break
		}

		if m, ok := t.customMessage[language][fieldError.Name+"."+fieldError.MessageName]; ok {
			errors[i].(*FieldError).SetMessage(m)
			break
		}

		message, ok := t.messages[language][fieldError.MessageName]
		if ok {
			attribute := fieldError.Attribute
			if customAttribute, ok := t.attributes[language][fieldError.StructName]; ok {
				attribute = customAttribute
			} else if fieldError.DefaultAttribute != "" {
				attribute = fieldError.DefaultAttribute
			}

			message = strings.Replace(message, "{{.Attribute}}", attribute, -1)

			for _, parameter := range fieldError.MessageParameters {
				message = strings.Replace(message, "{{."+parameter.Key+"}}", parameter.Value, -1)
			}

			errors[i].(*FieldError).SetMessage(message)
		}
	}

	return errors
}
