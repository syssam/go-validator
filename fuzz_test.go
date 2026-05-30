package validator

import (
	"strings"
	"testing"
)

// adversarial seeds shared across fuzz targets
func addSeeds(f *testing.F) {
	for _, s := range []string{
		"", " ", "test@example.com", "http://example.com/a?b=c",
		"123e4567-e89b-12d3-a456-426614174000",
		strings.Repeat("a", 200000),          // very long
		strings.Repeat("a@", 50000),          // long with structure
		"\x00\xff\xfe",                       // null + invalid UTF-8
		"\u202e\u0301\u00e8",                 // RTL override + combining mark + multibyte
		"+14155552671", "++++", "0000000000", // phone-ish
		"2024-13-45T99:99:99Z", // bogus date
	} {
		f.Add(s)
	}
}

// FuzzStringValidators ensures the regex/string validators never panic.
func FuzzStringValidators(f *testing.F) {
	addSeeds(f)
	f.Fuzz(func(t *testing.T, s string) {
		_ = IsEmail(s)
		_ = IsURL(s)
		_ = IsUUID(s)
		_ = IsAlpha(s)
		_ = IsASCII(s)
	})
}

// FuzzPhone exercises the third-party phonenumbers library (most likely to panic).
func FuzzPhone(f *testing.F) {
	addSeeds(f)
	f.Fuzz(func(t *testing.T, s string) {
		_ = IsPhoneE164(s)
		_ = IsPhoneValid(s)
	})
}

// FuzzConverters covers the string->number/bool conversion helpers.
func FuzzConverters(f *testing.F) {
	addSeeds(f)
	f.Fuzz(func(t *testing.T, s string) {
		_, _ = ToFloat(s)
		_ = ToBool(s)
		_, _ = ToInt(s)
		_ = ToString(s)
	})
}

// FuzzDate covers the date parser.
func FuzzDate(f *testing.F) {
	addSeeds(f)
	f.Fuzz(func(t *testing.T, s string) {
		_, _ = IsDate(s)
		_, _ = IsDateFormat(s, []string{"2006-01-02"})
		_, _ = IsAfter(s, []string{"2020-01-01"})
	})
}

// FuzzValidateStruct drives the full reflection path with adversarial field values.
func FuzzValidateStruct(f *testing.F) {
	f.Add("a@b.com", "http://x", "Name")
	f.Add("", "", "")
	f.Add(strings.Repeat("x", 100000), "\x00", "\xff\xfe")
	f.Fuzz(func(t *testing.T, email, url, name string) {
		type S struct {
			Email string `valid:"email"`
			URL   string `valid:"url"`
			Name  string `valid:"required,alphaNum"`
		}
		_ = ValidateStruct(&S{Email: email, URL: url, Name: name})
	})
}
