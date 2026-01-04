package validator

import "testing"

func TestValidateHexColor(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		// Valid formats
		{"empty string", "", true},
		{"3 char with hash", "#fff", true},
		{"3 char uppercase with hash", "#FFF", true},
		{"3 char mixed case with hash", "#fFf", true},
		{"6 char with hash", "#ffffff", true},
		{"6 char uppercase with hash", "#FFFFFF", true},
		{"6 char mixed case with hash", "#ffFFff", true},
		{"3 char without hash", "fff", true},
		{"6 char without hash", "ffffff", true},
		{"valid color #000", "#000", true},
		{"valid color #abc123", "#abc123", true},

		// Invalid formats
		{"4 char", "#ffff", false},
		{"5 char", "#fffff", false},
		{"7 char", "#fffffff", false},
		{"invalid char g", "#ggg", false},
		{"invalid char z", "#zzz", false},
		{"too short", "#ff", false},
		{"too long", "#fffffffff", false},
		{"double hash", "##fff", false},
		{"spaces", "# fff", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateHexColor(tt.input); got != tt.want {
				t.Errorf("ValidateHexColor(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestValidateTimezone(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		// Valid timezones
		{"empty string", "", true},
		{"UTC", "UTC", true},
		{"Local", "Local", true},
		{"America/New_York", "America/New_York", true},
		{"America/Los_Angeles", "America/Los_Angeles", true},
		{"Europe/London", "Europe/London", true},
		{"Europe/Paris", "Europe/Paris", true},
		{"Asia/Tokyo", "Asia/Tokyo", true},
		{"Asia/Shanghai", "Asia/Shanghai", true},
		{"Australia/Sydney", "Australia/Sydney", true},
		{"Pacific/Auckland", "Pacific/Auckland", true},

		// Invalid timezones
		{"invalid", "Invalid/Timezone", false},
		{"random string", "not_a_timezone", false},
		{"partial", "America", false},
		{"typo", "America/NewYork", false},
		{"spaces", "America/New York", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateTimezone(tt.input); got != tt.want {
				t.Errorf("ValidateTimezone(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestValidateStructHexColor(t *testing.T) {
	type Theme struct {
		PrimaryColor   string `valid:"required,hexColor"`
		SecondaryColor string `valid:"hexColor"`
	}

	tests := []struct {
		name      string
		theme     Theme
		wantError bool
	}{
		{
			name:      "valid colors",
			theme:     Theme{PrimaryColor: "#ff5733", SecondaryColor: "#333"},
			wantError: false,
		},
		{
			name:      "valid without secondary",
			theme:     Theme{PrimaryColor: "#ffffff"},
			wantError: false,
		},
		{
			name:      "invalid primary color",
			theme:     Theme{PrimaryColor: "not-a-color"},
			wantError: true,
		},
		{
			name:      "invalid secondary color",
			theme:     Theme{PrimaryColor: "#fff", SecondaryColor: "#gggggg"},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.theme)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateStruct() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateStructTimezone(t *testing.T) {
	type UserSettings struct {
		Timezone string `valid:"required,timezone"`
	}

	tests := []struct {
		name      string
		settings  UserSettings
		wantError bool
	}{
		{
			name:      "valid timezone UTC",
			settings:  UserSettings{Timezone: "UTC"},
			wantError: false,
		},
		{
			name:      "valid timezone America/New_York",
			settings:  UserSettings{Timezone: "America/New_York"},
			wantError: false,
		},
		{
			name:      "valid timezone Asia/Tokyo",
			settings:  UserSettings{Timezone: "Asia/Tokyo"},
			wantError: false,
		},
		{
			name:      "invalid timezone",
			settings:  UserSettings{Timezone: "Invalid/Zone"},
			wantError: true,
		},
		{
			name:      "empty timezone fails required",
			settings:  UserSettings{Timezone: ""},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.settings)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateStruct() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateASCII(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"empty string", "", true},
		{"ascii only", "hello world", true},
		{"ascii with numbers", "hello123", true},
		{"ascii with symbols", "hello@world.com", true},
		{"unicode", "héllo", false},
		{"chinese", "你好", false},
		{"emoji", "hello 😀", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateASCII(tt.input); got != tt.want {
				t.Errorf("ValidateASCII(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestValidateLowercase(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"empty string", "", true},
		{"lowercase", "hello", true},
		{"lowercase with numbers", "hello123", true},
		{"uppercase", "HELLO", false},
		{"mixed case", "Hello", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateLowercase(tt.input); got != tt.want {
				t.Errorf("ValidateLowercase(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestValidateUppercase(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"empty string", "", true},
		{"uppercase", "HELLO", true},
		{"uppercase with numbers", "HELLO123", true},
		{"lowercase", "hello", false},
		{"mixed case", "Hello", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateUppercase(tt.input); got != tt.want {
				t.Errorf("ValidateUppercase(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestValidateStartsWith(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		params []string
		want   bool
	}{
		{"empty string", "", []string{"hello"}, true},
		{"starts with", "hello world", []string{"hello"}, true},
		{"starts with one of many", "hello world", []string{"hi", "hello"}, true},
		{"does not start with", "hello world", []string{"world"}, false},
		{"case sensitive", "Hello world", []string{"hello"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateStartsWith(tt.input, tt.params); got != tt.want {
				t.Errorf("ValidateStartsWith(%q, %v) = %v, want %v", tt.input, tt.params, got, tt.want)
			}
		})
	}
}

func TestValidateEndsWith(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		params []string
		want   bool
	}{
		{"empty string", "", []string{"world"}, true},
		{"ends with", "hello world", []string{"world"}, true},
		{"ends with one of many", "hello world", []string{"hello", "world"}, true},
		{"does not end with", "hello world", []string{"hello"}, false},
		{"case sensitive", "hello World", []string{"world"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateEndsWith(tt.input, tt.params); got != tt.want {
				t.Errorf("ValidateEndsWith(%q, %v) = %v, want %v", tt.input, tt.params, got, tt.want)
			}
		})
	}
}

func TestValidateDoesntStartWith(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		params []string
		want   bool
	}{
		{"empty string", "", []string{"hello"}, true},
		{"does not start with", "hello world", []string{"world"}, true},
		{"starts with", "hello world", []string{"hello"}, false},
		{"starts with one of many", "hello world", []string{"hi", "hello"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateDoesntStartWith(tt.input, tt.params); got != tt.want {
				t.Errorf("ValidateDoesntStartWith(%q, %v) = %v, want %v", tt.input, tt.params, got, tt.want)
			}
		})
	}
}

func TestValidateDoesntEndWith(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		params []string
		want   bool
	}{
		{"empty string", "", []string{"world"}, true},
		{"does not end with", "hello world", []string{"hello"}, true},
		{"ends with", "hello world", []string{"world"}, false},
		{"ends with one of many", "hello world", []string{"hello", "world"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateDoesntEndWith(tt.input, tt.params); got != tt.want {
				t.Errorf("ValidateDoesntEndWith(%q, %v) = %v, want %v", tt.input, tt.params, got, tt.want)
			}
		})
	}
}

func TestValidateMACAddress(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"empty string", "", true},
		{"colon format", "00:1A:2B:3C:4D:5E", true},
		{"hyphen format", "00-1A-2B-3C-4D-5E", true},
		{"cisco format", "001A.2B3C.4D5E", true},
		{"lowercase", "00:1a:2b:3c:4d:5e", true},
		{"invalid", "00:1A:2B:3C:4D", false},
		{"invalid chars", "00:1G:2B:3C:4D:5E", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateMACAddress(tt.input); got != tt.want {
				t.Errorf("ValidateMACAddress(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestValidateULID(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{"empty string", "", true},
		{"valid ulid", "01ARZ3NDEKTSV4RRFFQ69G5FAV", true},
		{"lowercase invalid", "01arz3ndektsv4rrffq69g5fav", false},
		{"too short", "01ARZ3NDEKTSV4RRFFQ69G5FA", false},
		{"too long", "01ARZ3NDEKTSV4RRFFQ69G5FAVX", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidateULID(tt.input); got != tt.want {
				t.Errorf("ValidateULID(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestValidateStructASCII(t *testing.T) {
	type TestASCII struct {
		Name string `valid:"ascii"`
	}

	tests := []struct {
		name      string
		data      TestASCII
		wantError bool
	}{
		{"ascii valid", TestASCII{Name: "hello"}, false},
		{"ascii invalid", TestASCII{Name: "héllo"}, true},
		{"empty valid", TestASCII{Name: ""}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.data)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateStruct() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateStructCase(t *testing.T) {
	type TestCase struct {
		Lower string `valid:"lowercase"`
		Upper string `valid:"uppercase"`
	}

	tests := []struct {
		name      string
		data      TestCase
		wantError bool
	}{
		{"both valid", TestCase{Lower: "hello", Upper: "WORLD"}, false},
		{"lower invalid", TestCase{Lower: "Hello", Upper: "WORLD"}, true},
		{"upper invalid", TestCase{Lower: "hello", Upper: "World"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.data)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateStruct() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateStructStartsEndsWith(t *testing.T) {
	type TestPrefix struct {
		URL string `valid:"startsWith=http|https"`
	}

	type TestSuffix struct {
		Email string `valid:"endsWith=.com|.org"`
	}

	t.Run("startsWith valid", func(t *testing.T) {
		err := ValidateStruct(TestPrefix{URL: "https://example.com"})
		if err != nil {
			t.Errorf("ValidateStruct() error = %v", err)
		}
	})

	t.Run("startsWith invalid", func(t *testing.T) {
		err := ValidateStruct(TestPrefix{URL: "ftp://example.com"})
		if err == nil {
			t.Errorf("ValidateStruct() expected error")
		}
	})

	t.Run("endsWith valid", func(t *testing.T) {
		err := ValidateStruct(TestSuffix{Email: "test@example.com"})
		if err != nil {
			t.Errorf("ValidateStruct() error = %v", err)
		}
	})

	t.Run("endsWith invalid", func(t *testing.T) {
		err := ValidateStruct(TestSuffix{Email: "test@example.net"})
		if err == nil {
			t.Errorf("ValidateStruct() expected error")
		}
	})
}

func TestValidateStructMAC(t *testing.T) {
	type TestMAC struct {
		Address string `valid:"macAddress"`
	}

	tests := []struct {
		name      string
		data      TestMAC
		wantError bool
	}{
		{"valid mac", TestMAC{Address: "00:1A:2B:3C:4D:5E"}, false},
		{"invalid mac", TestMAC{Address: "invalid"}, true},
		{"empty valid", TestMAC{Address: ""}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.data)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateStruct() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateAccepted(t *testing.T) {
	type TestAccepted struct {
		Terms string `valid:"accepted"`
	}

	tests := []struct {
		name      string
		data      TestAccepted
		wantError bool
	}{
		{"accepted yes", TestAccepted{Terms: "yes"}, false},
		{"accepted on", TestAccepted{Terms: "on"}, false},
		{"accepted 1", TestAccepted{Terms: "1"}, false},
		{"accepted true", TestAccepted{Terms: "true"}, false},
		{"not accepted no", TestAccepted{Terms: "no"}, true},
		{"not accepted empty", TestAccepted{Terms: ""}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.data)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateStruct() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateDeclined(t *testing.T) {
	type TestDeclined struct {
		OptOut string `valid:"declined"`
	}

	tests := []struct {
		name      string
		data      TestDeclined
		wantError bool
	}{
		{"declined no", TestDeclined{OptOut: "no"}, false},
		{"declined off", TestDeclined{OptOut: "off"}, false},
		{"declined 0", TestDeclined{OptOut: "0"}, false},
		{"declined false", TestDeclined{OptOut: "false"}, false},
		{"not declined yes", TestDeclined{OptOut: "yes"}, true},
		{"not declined empty", TestDeclined{OptOut: ""}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.data)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateStruct() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateProhibited(t *testing.T) {
	type TestProhibited struct {
		SecretField string `valid:"prohibited"`
	}

	tests := []struct {
		name      string
		data      TestProhibited
		wantError bool
	}{
		{"empty allowed", TestProhibited{SecretField: ""}, false},
		{"has value prohibited", TestProhibited{SecretField: "secret"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.data)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateStruct() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateMissing(t *testing.T) {
	type TestMissing struct {
		Field string `valid:"missing"`
	}

	tests := []struct {
		name      string
		data      TestMissing
		wantError bool
	}{
		{"empty is missing", TestMissing{Field: ""}, false},
		{"has value not missing", TestMissing{Field: "value"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.data)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateStruct() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidatePresent(t *testing.T) {
	type TestPresent struct {
		Field string `valid:"present"`
	}

	tests := []struct {
		name      string
		data      TestPresent
		wantError bool
	}{
		{"empty is present", TestPresent{Field: ""}, false},
		{"with value is present", TestPresent{Field: "value"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.data)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateStruct() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateMultipleOf(t *testing.T) {
	type TestMultiple struct {
		Number int `valid:"multipleOf=5"`
	}

	tests := []struct {
		name      string
		data      TestMultiple
		wantError bool
	}{
		{"10 is multiple of 5", TestMultiple{Number: 10}, false},
		{"15 is multiple of 5", TestMultiple{Number: 15}, false},
		{"0 is multiple of 5", TestMultiple{Number: 0}, false},
		{"7 is not multiple of 5", TestMultiple{Number: 7}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.data)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateStruct() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateMaxDigits(t *testing.T) {
	type TestMaxDigits struct {
		Code int `valid:"maxDigits=4"`
	}

	tests := []struct {
		name      string
		data      TestMaxDigits
		wantError bool
	}{
		{"1234 has 4 digits", TestMaxDigits{Code: 1234}, false},
		{"123 has 3 digits", TestMaxDigits{Code: 123}, false},
		{"12345 has 5 digits", TestMaxDigits{Code: 12345}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.data)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateStruct() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateMinDigits(t *testing.T) {
	type TestMinDigits struct {
		Code int `valid:"minDigits=4"`
	}

	tests := []struct {
		name      string
		data      TestMinDigits
		wantError bool
	}{
		{"1234 has 4 digits", TestMinDigits{Code: 1234}, false},
		{"12345 has 5 digits", TestMinDigits{Code: 12345}, false},
		{"123 has 3 digits", TestMinDigits{Code: 123}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.data)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateStruct() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateContains(t *testing.T) {
	type TestContains struct {
		Text string `valid:"contains=hello"`
	}

	tests := []struct {
		name      string
		data      TestContains
		wantError bool
	}{
		{"contains hello", TestContains{Text: "say hello world"}, false},
		{"does not contain", TestContains{Text: "say goodbye"}, true},
		{"empty string", TestContains{Text: ""}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.data)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateStruct() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateDoesntContain(t *testing.T) {
	type TestDoesntContain struct {
		Text string `valid:"doesntContain=spam"`
	}

	tests := []struct {
		name      string
		data      TestDoesntContain
		wantError bool
	}{
		{"doesnt contain spam", TestDoesntContain{Text: "hello world"}, false},
		{"contains spam", TestDoesntContain{Text: "this is spam"}, true},
		{"empty is ok", TestDoesntContain{Text: ""}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.data)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateStruct() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateDecimalPrecision(t *testing.T) {
	type TestDecimal struct {
		Price float64 `valid:"decimal=2"`
	}

	tests := []struct {
		name      string
		data      TestDecimal
		wantError bool
	}{
		{"10.99 has 2 decimals", TestDecimal{Price: 10.99}, false},
		{"10.12 has 2 decimals", TestDecimal{Price: 10.12}, false},
		{"10 has 0 decimals", TestDecimal{Price: 10}, true},
		{"10.1 has 1 decimal", TestDecimal{Price: 10.1}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.data)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateStruct() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateList(t *testing.T) {
	type TestList struct {
		Items []string `valid:"list"`
	}

	tests := []struct {
		name      string
		data      TestList
		wantError bool
	}{
		{"valid list", TestList{Items: []string{"a", "b", "c"}}, false},
		{"empty list", TestList{Items: []string{}}, false},
		{"nil list", TestList{Items: nil}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.data)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateStruct() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateRequiredArrayKeys(t *testing.T) {
	type TestArrayKeys struct {
		Data map[string]string `valid:"requiredArrayKeys=name|email"`
	}

	tests := []struct {
		name      string
		data      TestArrayKeys
		wantError bool
	}{
		{"has all keys", TestArrayKeys{Data: map[string]string{"name": "John", "email": "john@example.com"}}, false},
		{"missing email", TestArrayKeys{Data: map[string]string{"name": "John"}}, true},
		{"missing name", TestArrayKeys{Data: map[string]string{"email": "john@example.com"}}, true},
		{"empty map", TestArrayKeys{Data: map[string]string{}}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.data)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateStruct() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateRequiredIfAccepted(t *testing.T) {
	type TestRequiredIfAccepted struct {
		Newsletter string `valid:""`
		Email      string `valid:"requiredIfAccepted=Newsletter"`
	}

	tests := []struct {
		name      string
		data      TestRequiredIfAccepted
		wantError bool
	}{
		{"newsletter yes, email provided", TestRequiredIfAccepted{Newsletter: "yes", Email: "test@example.com"}, false},
		{"newsletter yes, email missing", TestRequiredIfAccepted{Newsletter: "yes", Email: ""}, true},
		{"newsletter no, email not required", TestRequiredIfAccepted{Newsletter: "no", Email: ""}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.data)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateStruct() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateRequiredIfDeclined(t *testing.T) {
	type TestRequiredIfDeclined struct {
		AutoPay string `valid:""`
		Reason  string `valid:"requiredIfDeclined=AutoPay"`
	}

	tests := []struct {
		name      string
		data      TestRequiredIfDeclined
		wantError bool
	}{
		{"autopay no, reason provided", TestRequiredIfDeclined{AutoPay: "no", Reason: "prefer manual"}, false},
		{"autopay no, reason missing", TestRequiredIfDeclined{AutoPay: "no", Reason: ""}, true},
		{"autopay yes, reason not required", TestRequiredIfDeclined{AutoPay: "yes", Reason: ""}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.data)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateStruct() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateProhibitedIfAccepted(t *testing.T) {
	type TestProhibitedIfAccepted struct {
		AutoFill string `valid:""`
		Manual   string `valid:"prohibitedIfAccepted=AutoFill"`
	}

	tests := []struct {
		name      string
		data      TestProhibitedIfAccepted
		wantError bool
	}{
		{"autofill yes, manual empty", TestProhibitedIfAccepted{AutoFill: "yes", Manual: ""}, false},
		{"autofill yes, manual provided", TestProhibitedIfAccepted{AutoFill: "yes", Manual: "data"}, true},
		{"autofill no, manual allowed", TestProhibitedIfAccepted{AutoFill: "no", Manual: "data"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.data)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateStruct() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateProhibitedIfDeclined(t *testing.T) {
	type TestProhibitedIfDeclined struct {
		Subscribe   string `valid:""`
		Preferences string `valid:"prohibitedIfDeclined=Subscribe"`
	}

	tests := []struct {
		name      string
		data      TestProhibitedIfDeclined
		wantError bool
	}{
		{"subscribe no, prefs empty", TestProhibitedIfDeclined{Subscribe: "no", Preferences: ""}, false},
		{"subscribe no, prefs provided", TestProhibitedIfDeclined{Subscribe: "no", Preferences: "daily"}, true},
		{"subscribe yes, prefs allowed", TestProhibitedIfDeclined{Subscribe: "yes", Preferences: "daily"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.data)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateStruct() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateAcceptedIf(t *testing.T) {
	type TestAcceptedIf struct {
		Country string `valid:""`
		Terms   string `valid:"acceptedIf=Country|US"`
	}

	tests := []struct {
		name      string
		data      TestAcceptedIf
		wantError bool
	}{
		{"US country, terms accepted", TestAcceptedIf{Country: "US", Terms: "yes"}, false},
		{"US country, terms not accepted", TestAcceptedIf{Country: "US", Terms: "no"}, true},
		{"other country, terms not required", TestAcceptedIf{Country: "UK", Terms: "no"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.data)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateStruct() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateDeclinedIf(t *testing.T) {
	type TestDeclinedIf struct {
		Type   string `valid:""`
		OptOut string `valid:"declinedIf=Type|premium"`
	}

	tests := []struct {
		name      string
		data      TestDeclinedIf
		wantError bool
	}{
		{"premium type, optout declined", TestDeclinedIf{Type: "premium", OptOut: "no"}, false},
		{"premium type, optout accepted", TestDeclinedIf{Type: "premium", OptOut: "yes"}, true},
		{"basic type, no restriction", TestDeclinedIf{Type: "basic", OptOut: "yes"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.data)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateStruct() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateMissingIf(t *testing.T) {
	type TestMissingIf struct {
		Type   string `valid:""`
		Secret string `valid:"missingIf=Type|public"`
	}

	tests := []struct {
		name      string
		data      TestMissingIf
		wantError bool
	}{
		{"public type, secret missing", TestMissingIf{Type: "public", Secret: ""}, false},
		{"public type, secret present", TestMissingIf{Type: "public", Secret: "value"}, true},
		{"private type, secret allowed", TestMissingIf{Type: "private", Secret: "value"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.data)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateStruct() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateMissingUnless(t *testing.T) {
	type TestMissingUnless struct {
		Role  string `valid:""`
		Admin string `valid:"missingUnless=Role|admin"`
	}

	tests := []struct {
		name      string
		data      TestMissingUnless
		wantError bool
	}{
		{"admin role, admin field allowed", TestMissingUnless{Role: "admin", Admin: "true"}, false},
		{"user role, admin field missing", TestMissingUnless{Role: "user", Admin: ""}, false},
		{"user role, admin field present", TestMissingUnless{Role: "user", Admin: "true"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.data)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateStruct() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateMissingWith(t *testing.T) {
	type TestMissingWith struct {
		Email string `valid:""`
		Phone string `valid:"missingWith=Email"`
	}

	tests := []struct {
		name      string
		data      TestMissingWith
		wantError bool
	}{
		{"email present, phone missing", TestMissingWith{Email: "test@example.com", Phone: ""}, false},
		{"email present, phone also present", TestMissingWith{Email: "test@example.com", Phone: "123"}, true},
		{"email missing, phone allowed", TestMissingWith{Email: "", Phone: "123"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.data)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateStruct() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateMissingWithAll(t *testing.T) {
	type TestMissingWithAll struct {
		Email string `valid:""`
		Phone string `valid:""`
		Fax   string `valid:"missingWithAll=Email|Phone"`
	}

	tests := []struct {
		name      string
		data      TestMissingWithAll
		wantError bool
	}{
		{"all present, fax must be missing", TestMissingWithAll{Email: "a@b.com", Phone: "123", Fax: ""}, false},
		{"all present, fax also present", TestMissingWithAll{Email: "a@b.com", Phone: "123", Fax: "456"}, true},
		{"one missing, fax allowed", TestMissingWithAll{Email: "a@b.com", Phone: "", Fax: "456"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.data)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateStruct() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidatePresentWith(t *testing.T) {
	type TestPresentWith struct {
		Password string `valid:""`
		Confirm  string `valid:"presentWith=Password"`
	}

	tests := []struct {
		name      string
		data      TestPresentWith
		wantError bool
	}{
		{"password present, confirm present", TestPresentWith{Password: "pass", Confirm: "pass"}, false},
		{"password present, confirm missing", TestPresentWith{Password: "pass", Confirm: ""}, true},
		{"password missing, confirm not required", TestPresentWith{Password: "", Confirm: ""}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.data)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateStruct() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidatePresentWithAll(t *testing.T) {
	type TestPresentWithAll struct {
		First    string `valid:""`
		Last     string `valid:""`
		FullName string `valid:"presentWithAll=First|Last"`
	}

	tests := []struct {
		name      string
		data      TestPresentWithAll
		wantError bool
	}{
		{"all present", TestPresentWithAll{First: "John", Last: "Doe", FullName: "John Doe"}, false},
		{"both present, full missing", TestPresentWithAll{First: "John", Last: "Doe", FullName: ""}, true},
		{"one missing, full not required", TestPresentWithAll{First: "John", Last: "", FullName: ""}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.data)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateStruct() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateProhibitedIf(t *testing.T) {
	type TestProhibitedIf struct {
		Type  string `valid:""`
		Extra string `valid:"prohibitedIf=Type|basic"`
	}

	tests := []struct {
		name      string
		data      TestProhibitedIf
		wantError bool
	}{
		{"basic type, extra empty", TestProhibitedIf{Type: "basic", Extra: ""}, false},
		{"basic type, extra present", TestProhibitedIf{Type: "basic", Extra: "data"}, true},
		{"premium type, extra allowed", TestProhibitedIf{Type: "premium", Extra: "data"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.data)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateStruct() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateProhibitedUnless(t *testing.T) {
	type TestProhibitedUnless struct {
		Role  string `valid:""`
		Admin string `valid:"prohibitedUnless=Role|admin"`
	}

	tests := []struct {
		name      string
		data      TestProhibitedUnless
		wantError bool
	}{
		{"admin role, admin allowed", TestProhibitedUnless{Role: "admin", Admin: "true"}, false},
		{"user role, admin empty", TestProhibitedUnless{Role: "user", Admin: ""}, false},
		{"user role, admin present", TestProhibitedUnless{Role: "user", Admin: "true"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.data)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateStruct() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidateProhibits(t *testing.T) {
	type TestProhibits struct {
		Email string `valid:"prohibits=Phone"`
		Phone string `valid:""`
	}

	tests := []struct {
		name      string
		data      TestProhibits
		wantError bool
	}{
		{"email present, phone empty", TestProhibits{Email: "test@example.com", Phone: ""}, false},
		{"email present, phone present", TestProhibits{Email: "test@example.com", Phone: "123"}, true},
		{"email empty, phone allowed", TestProhibits{Email: "", Phone: "123"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.data)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateStruct() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidatePresentIf(t *testing.T) {
	type TestPresentIf struct {
		Type    string `valid:""`
		Details string `valid:"presentIf=Type|custom"`
	}

	tests := []struct {
		name      string
		data      TestPresentIf
		wantError bool
	}{
		{"custom type, details present", TestPresentIf{Type: "custom", Details: "info"}, false},
		{"custom type, details empty", TestPresentIf{Type: "custom", Details: ""}, true},
		{"standard type, details not required", TestPresentIf{Type: "standard", Details: ""}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.data)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateStruct() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestValidatePresentUnless(t *testing.T) {
	type TestPresentUnless struct {
		Status string `valid:""`
		Reason string `valid:"presentUnless=Status|approved"`
	}

	tests := []struct {
		name      string
		data      TestPresentUnless
		wantError bool
	}{
		{"approved status, reason optional", TestPresentUnless{Status: "approved", Reason: ""}, false},
		{"rejected status, reason required", TestPresentUnless{Status: "rejected", Reason: "issue"}, false},
		{"rejected status, reason missing", TestPresentUnless{Status: "rejected", Reason: ""}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateStruct(tt.data)
			if (err != nil) != tt.wantError {
				t.Errorf("ValidateStruct() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}
