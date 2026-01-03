package validator

import (
	"fmt"
	"net"
	"net/url"
	"strings"
	"time"
	"unicode/utf8"
)

// IsStringBetween checks if string length is between left and right
func IsStringBetween(v string, left, right int64) bool {
	return IsInt64Between(int64(utf8.RuneCountInString(v)), left, right)
}

// InString check if string str is a member of the set of strings params
func InString(str string, params []string) bool {
	for _, param := range params {
		if str == param {
			return true
		}
	}

	return false
}

// compareString determine if a comparison passes between the given values.
func compareString(first string, second int64, operator string) (bool, error) {
	switch operator {
	case "<":
		return int64(utf8.RuneCountInString(first)) < second, nil
	case ">":
		return int64(utf8.RuneCountInString(first)) > second, nil
	case "<=":
		return int64(utf8.RuneCountInString(first)) <= second, nil
	case ">=":
		return int64(utf8.RuneCountInString(first)) >= second, nil
	case "==":
		return int64(utf8.RuneCountInString(first)) == second, nil
	default:
		return false, fmt.Errorf("validator: compareString unsupported operator %s", operator)
	}
}

// IsNumeric check if the string must be numeric. Empty string is valid.
func IsNumeric(str string) bool {
	if IsNull(str) {
		return true
	}
	return rxNumeric.MatchString(str)
}

// IsInt check if the string must be an integer. Empty string is valid.
func IsInt(str string) bool {
	if IsNull(str) {
		return true
	}
	return rxInt.MatchString(str)
}

// IsFloat check if the string must be an float. Empty string is valid.
func IsFloat(str string) bool {
	if IsNull(str) {
		return true
	}
	return rxFloat.MatchString(str)
}

// IsNull check if the string is null.
func IsNull(str string) bool {
	return str == ""
}

// IsEmptyString check if the string is empty.
func IsEmptyString(str string) bool {
	return strings.TrimSpace(str) == ""
}

// IsEmail check if the string is an email.
func IsEmail(str string) bool {
	return rxEmail.MatchString(str)
}

// IsAlpha check if the string may be only contains letters (a-zA-Z). Empty string is valid.
func IsAlpha(str string) bool {
	if IsNull(str) {
		return true
	}
	return rxAlpha.MatchString(str)
}

// IsAlphaNum check if the string may be only contains letters and numbers. Empty string is valid.
func IsAlphaNum(str string) bool {
	if IsNull(str) {
		return true
	}
	return rxAlphaNum.MatchString(str)
}

// IsAlphaDash check if the string may be only contains letters, numbers, dashes and underscores. Empty string is valid.
func IsAlphaDash(str string) bool {
	if IsNull(str) {
		return true
	}
	return rxAlphaDash.MatchString(str)
}

// IsAlphaUnicode check if the string may be only contains letters (a-zA-Z). Empty string is valid.
func IsAlphaUnicode(str string) bool {
	if IsNull(str) {
		return true
	}
	return rxAlphaUnicode.MatchString(str)
}

// IsAlphaNumUnicode check if the string may be only contains letters and numbers. Empty string is valid.
func IsAlphaNumUnicode(str string) bool {
	if IsNull(str) {
		return true
	}
	return rxAlphaNumUnicode.MatchString(str)
}

// IsAlphaDashUnicode check if the string may be only contains letters, numbers, dashes and underscores. Empty string is valid.
func IsAlphaDashUnicode(str string) bool {
	if IsNull(str) {
		return true
	}
	return rxAlphaDashUnicode.MatchString(str)
}

// IsIP check if the string is an ip address.
func IsIP(v string) bool {
	ip := net.ParseIP(v)
	return ip != nil
}

// IsIPv4 check if the string is an ipv4 address.
func IsIPv4(v string) bool {
	ip := net.ParseIP(v)
	return ip != nil && ip.To4() != nil
}

// IsIPv6 check if the string is an ipv6 address.
func IsIPv6(v string) bool {
	ip := net.ParseIP(v)
	return ip != nil && ip.To4() == nil
}

// IsUUID3 check if the string is an uuid3.
func IsUUID3(str string) bool {
	if IsNull(str) {
		return true
	}
	return rxUUID3.MatchString(str)
}

// IsUUID4 check if the string is an uuid4.
func IsUUID4(str string) bool {
	if IsNull(str) {
		return true
	}
	return rxUUID4.MatchString(str)
}

// IsUUID5 check if the string is an uuid5.
func IsUUID5(str string) bool {
	if IsNull(str) {
		return true
	}
	return rxUUID5.MatchString(str)
}

// IsUUID check if the string is an uuid.
func IsUUID(str string) bool {
	if IsNull(str) {
		return true
	}
	return rxUUID.MatchString(str)
}

// IsURL check if the string is an URL.
func IsURL(str string) bool {
	var i int

	if IsNull(str) {
		return true
	}

	if i = strings.Index(str, "#"); i > -1 {
		str = str[:i]
	}

	url, err := url.ParseRequestURI(str)
	if err != nil || url.Scheme == "" {
		return false
	}

	return true
}

// IsHexColor check if the string is a valid hex color code.
// Supports 3 and 6 character formats with or without # prefix.
// Examples: #fff, #FFF, #ffffff, #FFFFFF, fff, ffffff
func IsHexColor(str string) bool {
	if IsNull(str) {
		return true
	}
	return rxHexColor.MatchString(str)
}

// IsTimezone check if the string is a valid timezone identifier.
// Uses Go's time.LoadLocation to validate.
// Examples: UTC, America/New_York, Europe/London, Asia/Tokyo
func IsTimezone(str string) bool {
	if IsNull(str) {
		return true
	}
	_, err := time.LoadLocation(str)
	return err == nil
}

// IsASCII check if the string contains only ASCII characters.
// Empty string is valid.
func IsASCII(str string) bool {
	if IsNull(str) {
		return true
	}
	return rxASCII.MatchString(str)
}

// IsLowercase check if the string is all lowercase.
// Empty string is valid.
func IsLowercase(str string) bool {
	if IsNull(str) {
		return true
	}
	return str == strings.ToLower(str)
}

// IsUppercase check if the string is all uppercase.
// Empty string is valid.
func IsUppercase(str string) bool {
	if IsNull(str) {
		return true
	}
	return str == strings.ToUpper(str)
}

// IsStartsWith check if the string starts with any of the given prefixes.
// Empty string is valid.
func IsStartsWith(str string, params []string) bool {
	if IsNull(str) {
		return true
	}
	for _, prefix := range params {
		if strings.HasPrefix(str, prefix) {
			return true
		}
	}
	return false
}

// IsEndsWith check if the string ends with any of the given suffixes.
// Empty string is valid.
func IsEndsWith(str string, params []string) bool {
	if IsNull(str) {
		return true
	}
	for _, suffix := range params {
		if strings.HasSuffix(str, suffix) {
			return true
		}
	}
	return false
}

// IsDoesntStartWith check if the string does not start with any of the given prefixes.
// Empty string is valid.
func IsDoesntStartWith(str string, params []string) bool {
	if IsNull(str) {
		return true
	}
	for _, prefix := range params {
		if strings.HasPrefix(str, prefix) {
			return false
		}
	}
	return true
}

// IsDoesntEndWith check if the string does not end with any of the given suffixes.
// Empty string is valid.
func IsDoesntEndWith(str string, params []string) bool {
	if IsNull(str) {
		return true
	}
	for _, suffix := range params {
		if strings.HasSuffix(str, suffix) {
			return false
		}
	}
	return true
}

// IsMACAddress check if the string is a valid MAC address.
// Supports formats: XX:XX:XX:XX:XX:XX, XX-XX-XX-XX-XX-XX, XXXX.XXXX.XXXX
// Empty string is valid.
func IsMACAddress(str string) bool {
	if IsNull(str) {
		return true
	}
	return rxMACAddress.MatchString(str)
}

// IsULID check if the string is a valid ULID.
// ULID format: 01ARZ3NDEKTSV4RRFFQ69G5FAV (26 characters, Crockford's base32)
// Empty string is valid.
func IsULID(str string) bool {
	if IsNull(str) {
		return true
	}
	return rxULID.MatchString(str)
}

// Deprecated: Use Is* functions instead.
var (
	ValidateBetweenString    = IsStringBetween
	ValidateEmail            = IsEmail
	ValidateAlpha            = IsAlpha
	ValidateAlphaNum         = IsAlphaNum
	ValidateAlphaDash        = IsAlphaDash
	ValidateAlphaUnicode     = IsAlphaUnicode
	ValidateAlphaNumUnicode  = IsAlphaNumUnicode
	ValidateAlphaDashUnicode = IsAlphaDashUnicode
	ValidateIP               = IsIP
	ValidateIPv4             = IsIPv4
	ValidateIPv6             = IsIPv6
	ValidateUUID             = IsUUID
	ValidateUUID3            = IsUUID3
	ValidateUUID4            = IsUUID4
	ValidateUUID5            = IsUUID5
	ValidateURL              = IsURL
	ValidateHexColor         = IsHexColor
	ValidateTimezone         = IsTimezone
	ValidateASCII            = IsASCII
	ValidateLowercase        = IsLowercase
	ValidateUppercase        = IsUppercase
	ValidateStartsWith       = IsStartsWith
	ValidateEndsWith         = IsEndsWith
	ValidateDoesntStartWith  = IsDoesntStartWith
	ValidateDoesntEndWith    = IsDoesntEndWith
	ValidateMACAddress       = IsMACAddress
	ValidateULID             = IsULID
)
