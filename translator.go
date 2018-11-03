package validator

import (
	"fmt"
	"strings"
	"sync"
)

// Translate translate type
type Translate map[string]string

// Translator type
type Translator struct {
	locale        string
	customMessage map[string]Translate
	messages      map[string]Translate
	attributes    map[string]Translate
}

var loadTranslatorOnce *Translator
var translatorOnce sync.Once

// NewTranslator returns a new instance of 'translator' with sane defaults.
func NewTranslator() *Translator {
	translatorOnce.Do(func() {
		loadTranslatorOnce = &Translator{
			locale: "en",
		}

		if loadTranslatorOnce.messages == nil {
			loadTranslatorOnce.messages = make(map[string]Translate)
		}

		if loadTranslatorOnce.attributes == nil {
			loadTranslatorOnce.attributes = make(map[string]Translate)
		}

		if loadTranslatorOnce.customMessage == nil {
			loadTranslatorOnce.customMessage = make(map[string]Translate)
		}
	})
	return loadTranslatorOnce
}

// SetLocale set locale
func (t *Translator) SetLocale(locale string) {
	t.locale = locale
}

// GetLocale get locale
func (t *Translator) GetLocale() string {
	return t.locale
}

// SetMessage set Message
func (t *Translator) SetMessage(langCode string, messages Translate) {
	t.messages[langCode] = messages
}

// LoadMessage load message
func (t *Translator) LoadMessage() Translate {
	return t.messages[t.locale]
}

// SetAttributes set attributes
func (t *Translator) SetAttributes(langCode string, messages Translate) {
	t.attributes[langCode] = messages
}

// Trans trans
func (t *Translator) Trans(structName string, messageName string, attribute string) string {
	locale := t.GetLocale()

	if m, ok := t.customMessage[locale][structName+"."+messageName]; ok {
		return m
	}

	message, ok := t.messages[locale][messageName]
	if ok {
		if customAttribute, ok := t.attributes[locale][structName]; ok {
			attribute = customAttribute
		}

		return strings.Replace(message, ":attribute", attribute, -1)
	}

	panic(fmt.Sprintf("validator: Trans undefined message %s on locale %s", messageName, locale))
}
