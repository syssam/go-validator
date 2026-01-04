package validator

import (
	"github.com/nyaruka/phonenumbers"
)

// IsPhoneE164 validates phone number in E.164 format (+[country code][number])
func IsPhoneE164(str string) bool {
	if str == "" {
		return false
	}
	// E.164 must start with +
	if str[0] != '+' {
		return false
	}
	num, err := phonenumbers.Parse(str, "")
	if err != nil {
		return false
	}
	return phonenumbers.IsValidNumber(num)
}

// IsPhone validates phone number for a specific region (ISO 3166-1 alpha-2)
func IsPhone(str string, region string) bool {
	if str == "" {
		return false
	}
	num, err := phonenumbers.Parse(str, region)
	if err != nil {
		return false
	}
	return phonenumbers.IsValidNumberForRegion(num, region)
}

// IsPhoneValid validates phone number (any format, any region)
func IsPhoneValid(str string) bool {
	if str == "" {
		return false
	}
	// Try E.164 first
	if str[0] == '+' {
		num, err := phonenumbers.Parse(str, "")
		if err != nil {
			return false
		}
		return phonenumbers.IsValidNumber(num)
	}
	// Without country code, we can't validate properly
	return false
}

// IsPhoneMobile validates if phone number is a mobile number
func IsPhoneMobile(str string, region string) bool {
	if str == "" {
		return false
	}
	num, err := phonenumbers.Parse(str, region)
	if err != nil {
		return false
	}
	if !phonenumbers.IsValidNumber(num) {
		return false
	}
	numType := phonenumbers.GetNumberType(num)
	return numType == phonenumbers.MOBILE || numType == phonenumbers.FIXED_LINE_OR_MOBILE
}

// FormatPhoneE164 formats phone number to E.164 format
func FormatPhoneE164(str string, region string) (string, error) {
	num, err := phonenumbers.Parse(str, region)
	if err != nil {
		return "", err
	}
	return phonenumbers.Format(num, phonenumbers.E164), nil
}

// FormatPhoneNational formats phone number to national format
func FormatPhoneNational(str string, region string) (string, error) {
	num, err := phonenumbers.Parse(str, region)
	if err != nil {
		return "", err
	}
	return phonenumbers.Format(num, phonenumbers.NATIONAL), nil
}

// FormatPhoneInternational formats phone number to international format
func FormatPhoneInternational(str string, region string) (string, error) {
	num, err := phonenumbers.Parse(str, region)
	if err != nil {
		return "", err
	}
	return phonenumbers.Format(num, phonenumbers.INTERNATIONAL), nil
}
