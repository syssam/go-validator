package validator

import "regexp"

// Basic regular expressions for validating strings
const (
	Email            string = "^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$"
	Alpha            string = "^[a-zA-Z]+$"
	AlphaNum         string = "^[a-zA-Z0-9]+$"
	AlphaDash        string = "^[a-zA-Z_-]+$"
	AlphaUnicode     string = "^[\\p{L}]+$"
	AlphaNumUnicode  string = "^[\\p{L}\\p{M}\\p{N}]+$"
	AlphaDashUnicode string = "^[\\p{L}\\p{M}\\p{N}_-]+$"
	Numeric          string = "^[0-9]+$"
	Int              string = "^(?:[-+]?(?:0|[1-9][0-9]*))$"
	Float            string = "^(?:[-+]?(?:[0-9]+))?(?:\\.[0-9]*)?(?:[eE][\\+\\-]?(?:[0-9]+))?$"
	UUID3            string = "^[0-9a-f]{8}-[0-9a-f]{4}-3[0-9a-f]{3}-[0-9a-f]{4}-[0-9a-f]{12}$"
	UUID4            string = "^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
	UUID5            string = "^[0-9a-f]{8}-[0-9a-f]{4}-5[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$"
	UUID             string = "^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$"
)

// Used by IsFilePath func
const (
	// Unknown is unresolved OS type
	Unknown = iota
	// Win is Windows type
	Win
	// Unix is *nix OS types
	Unix
)

var (
	rxEmail            = regexp.MustCompile(Email)
	rxAlpha            = regexp.MustCompile(Alpha)
	rxAlphaNum         = regexp.MustCompile(AlphaNum)
	rxAlphaDash        = regexp.MustCompile(AlphaDash)
	rxAlphaUnicode     = regexp.MustCompile(AlphaUnicode)
	rxAlphaNumUnicode  = regexp.MustCompile(AlphaNumUnicode)
	rxAlphaDashUnicode = regexp.MustCompile(AlphaDashUnicode)
	rxNumeric          = regexp.MustCompile(Numeric)
	rxInt              = regexp.MustCompile(Int)
	rxFloat            = regexp.MustCompile(Float)
	rxUUID3            = regexp.MustCompile(UUID3)
	rxUUID4            = regexp.MustCompile(UUID4)
	rxUUID5            = regexp.MustCompile(UUID5)
	rxUUID             = regexp.MustCompile(UUID)
)
