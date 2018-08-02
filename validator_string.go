package validator

import (
	"fmt"
	"net"
	"unicode/utf8"
)

// BetweenString is
func BetweenString(v string, left int64, right int64) bool {
	return DigitsBetweenInt64(int64(utf8.RuneCountInString(v)), left, right)
}

// InString check if string str is a member of the set of strings params
func InString(str string, params ...string) bool {
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
	return rxAlphaNum.MatchString(str)
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

// IsEmail check if the string is an email.
func IsEmail(str string) bool {
	// TODO uppercase letters are not supported
	return rxEmail.MatchString(str)
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

// IsIP check if the string is an ip address.
func IsIP(v string) bool {
	ip := net.ParseIP(v)
	return ip != nil
}
