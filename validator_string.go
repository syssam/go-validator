package validator

import (
	"fmt"
	"net"
	"net/url"
	"strings"
	"unicode/utf8"
)

// ValidateBetweenString is
func ValidateBetweenString(v string, left int64, right int64) bool {
	return ValidateDigitsBetweenInt64(int64(utf8.RuneCountInString(v)), left, right)
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

//  compareString determine if a comparison passes between the given values.
func compareString(first string, second int64, operator string) bool {
	switch operator {
	case "<":
		return int64(utf8.RuneCountInString(first)) < second
	case ">":
		return int64(utf8.RuneCountInString(first)) > second
	case "<=":
		return int64(utf8.RuneCountInString(first)) <= second
	case ">=":
		return int64(utf8.RuneCountInString(first)) >= second
	case "==":
		return int64(utf8.RuneCountInString(first)) == second
	default:
		panic(fmt.Sprintf("validator: compareString unsupport operator %s", operator))
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
	return len(str) == 0
}

// ValidateEmail check if the string is an email.
func ValidateEmail(str string) bool {
	// TODO uppercase letters are not supported
	return rxEmail.MatchString(str)
}

// ValidateAlpha check if the string may be only contains letters (a-zA-Z). Empty string is valid.
func ValidateAlpha(str string) bool {
	if IsNull(str) {
		return true
	}
	return rxAlpha.MatchString(str)
}

// ValidateAlphaNum check if the string may be only contains letters and numbers. Empty string is valid.
func ValidateAlphaNum(str string) bool {
	if IsNull(str) {
		return true
	}
	return rxAlphaNum.MatchString(str)
}

// ValidateAlphaDash check if the string may be only contains letters, numbers, dashes and underscores. Empty string is valid.
func ValidateAlphaDash(str string) bool {
	if IsNull(str) {
		return true
	}
	return rxAlphaDash.MatchString(str)
}

// ValidateAlphaUnicode check if the string may be only contains letters (a-zA-Z). Empty string is valid.
func ValidateAlphaUnicode(str string) bool {
	if IsNull(str) {
		return true
	}
	return rxAlphaUnicode.MatchString(str)
}

// ValidateAlphaNumUnicode check if the string may be only contains letters and numbers. Empty string is valid.
func ValidateAlphaNumUnicode(str string) bool {
	if IsNull(str) {
		return true
	}
	return rxAlphaNumUnicode.MatchString(str)
}

// ValidateAlphaDashUnicode check if the string may be only contains letters, numbers, dashes and underscores. Empty string is valid.
func ValidateAlphaDashUnicode(str string) bool {
	if IsNull(str) {
		return true
	}
	return rxAlphaDashUnicode.MatchString(str)
}

// ValidateIP check if the string is an ip address.
func ValidateIP(v string) bool {
	ip := net.ParseIP(v)
	return ip != nil
}

// ValidateIPv4 check if the string is an ipv4 address.
func ValidateIPv4(v string) bool {
	ip := net.ParseIP(v)
	return ip != nil && ip.To4() != nil
}

// ValidateIPv6 check if the string is an ipv6 address.
func ValidateIPv6(v string) bool {
	ip := net.ParseIP(v)
	return ip != nil && ip.To4() == nil
}

// ValidateUUID3 check if the string is an uuid3.
func ValidateUUID3(str string) bool {
	if IsNull(str) {
		return true
	}
	return rxUUID3.MatchString(str)
}

// ValidateUUID4 check if the string is an uuid4.
func ValidateUUID4(str string) bool {
	if IsNull(str) {
		return true
	}
	return rxUUID4.MatchString(str)
}

// ValidateUUID5 check if the string is an uuid5.
func ValidateUUID5(str string) bool {
	if IsNull(str) {
		return true
	}
	return rxUUID5.MatchString(str)
}

// ValidateUUID check if the string is an uuid.
func ValidateUUID(str string) bool {
	if IsNull(str) {
		return true
	}
	return rxUUID.MatchString(str)
}

// ValidateURL check if the string is an URL.
func ValidateURL(str string) bool {
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
