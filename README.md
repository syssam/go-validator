<h1>go-validator</h1>
<p>
  <a href='https://github.com/syssam/go-validator/actions/workflows/ci.yml'><img src='https://github.com/syssam/go-validator/actions/workflows/ci.yml/badge.svg' alt='CI Status' /></a>
  <a href="https://goreportcard.com/report/github.com/syssam/go-validator" rel="nofollow"><img src="https://goreportcard.com/badge/github.com/syssam/go-validator" alt="Go Report Card" data-canonical-src="https://goreportcard.com/badge/github.com/syssam/go-validator" style="max-width:100%;"></a>
  <a href="https://godoc.org/github.com/syssam/go-validator" rel="nofollow"><img src="https://godoc.org/github.com/syssam/go-validator?status.svg" alt="GoDoc" data-canonical-src="https://godoc.org/github.com/syssam/go-validator?status.svg" style="max-width:100%;"></a>
  <a href="https://codecov.io/gh/syssam/go-validator"><img src="https://codecov.io/gh/syssam/go-validator/graph/badge.svg" alt="Code Coverage"/></a>
</p>
<p>A package of validators and sanitizers for strings, structs and collections.</p>
<p>features:</p>
<ul>
  <li>Customizable Attributes.</li>
  <li>Customizable error messages.</li>
  <li>Support i18n messages</li>
</ul>
<h2>Installation</h2>
<p>Make sure that Go is installed on your computer. Type the following command in your terminal:</p>
<p>go get github.com/syssam/go-validator</p>
<h2>Usage and documentation</h2>
<h5>Examples:</h5>
<ul>
  <li><a href="https://github.com/syssam/go-validator/tree/master/_examples/simple">Simple</a></li>
  <li><a href="https://github.com/syssam/go-validator/tree/master/_examples/translations">Translations</a></li>
  <li><a href="https://github.com/syssam/go-validator/tree/master/_examples/translation">Simple Translation</a></li>
  <li><a href="https://github.com/syssam/go-validator/tree/master/_examples/gin">Gin</a></li>
  <li><a href="https://github.com/syssam/go-validator/tree/master/_examples/echo">Echo</a></li>
  <li><a href="https://github.com/syssam/go-validator/tree/master/_examples/iris">Iris</a></li>
  <li><a href="https://github.com/syssam/go-validator/tree/master/_examples/custom">Custom Validation Rules</a></li>
  <li><a href="https://github.com/syssam/go-validator/tree/master/_examples/customtype">Custom Type Functions (Omittable, sql.Null)</a></li>
</ul>
<h2>Available Validation Rules</h2>
<ul>
    <li><a>omitempty</a></li>
    <li><a>nullable</a></li>
    <li><a>required</a></li>
    <li><a>requiredIf</a></li>
    <li><a>requiredUnless</a></li>
    <li><a>requiredWith</a></li>
    <li><a>requiredWithAll</a></li>
    <li><a>requiredWithout</a></li>
    <li><a>requiredWithoutAll</a></li>
    <li><a>between</a></li>
    <li><a>digitsBetween</a></li>
    <li><a>size</a></li>
    <li><a>max</a></li>
    <li><a>min</a></li>
    <li><a>same</a></li>
    <li><a>gt</a></li>
    <li><a>gte</a></li>
    <li><a>lt</a></li>
    <li><a>lte</a></li>
    <li><a>distinct</a></li>
    <li><a>email</a></li>
    <li><a>alpha</a></li>
    <li><a>alphaNum</a></li>
    <li><a>alphaDash</a></li>
    <li><a>alphaUnicode</a></li>
    <li><a>alphaNumUnicode</a></li>
    <li><a>alphaDashUnicode</a></li>
    <li><a>numeric</a></li>
    <li><a>int</a></li>
    <li><a>integer</a></li>
    <li><a>float</a></li>
    <li><a>null</a></li>
    <li><a>ip</a></li>
    <li><a>ipv4</a></li>
    <li><a>ipv6</a></li>
    <li><a>uuid3</a></li>
    <li><a>uuid4</a></li>
    <li><a>uuid5</a></li>
    <li><a>uuid</a></li>
    <li><a>url</a></li>
    <li><a>date</a></li>
    <li><a>dateFormat</a></li>
    <li><a>after</a></li>
    <li><a>afterOrEqual</a></li>
    <li><a>before</a></li>
    <li><a>beforeOrEqual</a></li>
    <li><a>regex</a></li>
    <li><a>notRegex</a></li>
    <li><a>in</a></li>
    <li><a>notIn</a></li>
    <li><a>boolean</a></li>
    <li><a>accepted</a></li>
    <li><a>declined</a></li>
    <li><a>country</a></li>
    <li><a>country.alpha2</a></li>
    <li><a>country.alpha3</a></li>
    <li><a>country.numeric</a></li>
    <li><a>currency</a></li>
    <li><a>currency.all</a></li>
    <li><a>currency.fiat</a></li>
    <li><a>currency.crypto</a></li>
    <li><a>language</a></li>
    <li><a>language.alpha2</a></li>
    <li><a>language.alpha3</a></li>
    <li><a>phone</a></li>
    <li><a>phone.e164</a></li>
</ul>
<h4 id="rule-omitempty">omitempty</h4>
<p>The "omitempty" option specifies that the field should be omitted from the encoding if the field has an empty value, defined as false, 0, a nil pointer, a nil interface value, and any empty array, slice, map, or string.</p>
<h4 id="rule-required">required</h4>
<p>The field under validation must be present in the input data and not empty. A field is considered "empty" if one of the following conditions are true:</p>
<div class="content-list">
  <ul>
    <li>The value is <code class="language-php"><span class="token keyword">nil</span></code>.</li>
    <li>The value is an empty string.</li>
    <li>The value is an empty array | map</li>
  </ul>
</div>
<h4 id="rule-requiredIf">requiredIf=anotherfield|value|...</h4>
<p>The field under validation must be present and not empty if the anotherfield field is equal to any value.</p>
<h4 id="rule-requiredIf">requiredUnless=anotherfield|value|...</h4>
<p>The field under validation must be present and not empty unless the anotherfield field is equal to any value.</p>
<h4 id="rule-requiredIf">requiredWith=anotherfield|anotherfield|...</h4>
<p>The field under validation must be present and not empty only if any of the other specified fields are present.</p>
<h4 id="rule-requiredIf">requiredWithAll=anotherfield|anotherfield|...</h4>
<p>The field under validation must be present and not empty only if all of the other specified fields are present.</p>
<h4 id="rule-requiredIf">requiredWithout=anotherfield|anotherfield|...</h4>
<p>The field under validation must be present and not empty only when any of the other specified fields are not present.</p>
<h4 id="rule-requiredIf">requiredWithoutAll=anotherfield|anotherfield|...</h4>
<p>The field under validation must be present and not empty only when all of the other specified fields are not present.</p>
<h4 id="rule-between">between=min|max</h4>
<p>The field under validation must have a size between the given min and max. String, Number, Array, Map are evaluated in the same fashion as the size rule.</p>
<h4 id="rule-between">digitsBetween=min|max</h4>
<p>The field under validation must have a length between the given min and max.</p>
<h4 id="rule-max">size=value</h4>
<p>The field under validation must have a size matching the given value. For string data, value corresponds to the number of characters. For numeric data, value corresponds to a given integer value. For an array | map | slice, size corresponds to the count of the array | map | slice.</p>
<h4 id="rule-max">max=value</h4>
<p>The field under validation must be less than or equal to a maximum value. String, Number, Array, Map are evaluated in the same fashion as the size rule.</p>
<h4 id="rule-min">min=value</h4>
<p>The field under validation must be greater than or equal to a minimum value. String, Number, Array, Map are evaluated in the same fashion as the size rule.</p>
<h4 id="rule-same">same=anotherfield</h4>
<p>The given field must match the field under validation.</p>
<h4 id="rule-gt">gt=anotherfield</h4>
<p>The field under validation must be greater than the given field. The two fields must be of the same type. String, Number, Array, Map are evaluated using the same conventions as the size rule.</p>
<h4 id="rule-gte">gte=anotherfield</h4>
<p>The field under validation must be greater than or equal to the given field. The two fields must be of the same type. String, Number, Array, Map are evaluated using the same conventions as the size rule.</p>
<h4 id="rule-lt">lt=anotherfield</h4>
<p>The field under validation must be less than the given field. The two fields must be of the same type. String, Number, Array, Map are evaluated using the same conventions as the size rule.</p>
<h4 id="rule-lte">lte=anotherfield</h4>
<p>The field under validation must be less than or equal to the given field. The two fields must be of the same type. String, Number, Array, Map are evaluated using the same conventions as the size rule.</p>
<h4 id="rule-distinct">distinct</h4>
<p>The field under validation must not have any duplicate values.</p>
<h4 id="rule-email">email</h4>
<p>The field under validation must be formatted as an e-mail address.</p>
<h4 id="rule-alpha">alpha</h4>
<p>The field under validation may be only contains letters. Empty string is valid.</p>
<h4 id="rule-alphaNum">alphaNum</h4>
<p>The field under validation may be only contains letters and numbers. Empty string is valid.</p>
<h4 id="rule-alphaDash">alphaDash</h4>
<p>The field under validation may be only contains letters, numbers, dashes and underscores. Empty string is valid.</p>
<h4 id="rule-alpha">alphaUnicode</h4>
<p>The field under validation may be only contains letters. Empty string is valid.</p>
<h4 id="rule-alphaNum">alphaNumUnicode</h4>
<p>The field under validation may be only contains letters and numbers. Empty string is valid.</p>
<h4 id="rule-alphaDash">alphaDashUnicode</h4>
<p>The field under validation may be only contains letters, numbers, dashes and underscores. Empty string is valid.</p>
<h4 id="rule-numeric">numeric</h4>
<p>The field under validation must be numbers. Empty string is valid.</p>
<h4 id="rule-int">int</h4>
<p>The field under validation must be int. Empty string is valid.</p>
<h4 id="rule-float">float</h4>
<p>The field under validation must be float. Empty string is valid.</p>
<h4 id="rule-ip">ip</h4>
<p>The field under validation must be an IP address.</p>
<h4 id="rule-ipv4">ipv4</h4>
<p>The field under validation must be an IPv4 address.</p>
<h4 id="rule-ipv6">ipv6</h4>
<p>The field under validation must be an IPv6 address.</p>
<h4 id="rule-ipv6">uuid3</h4>
<p>The field under validation must be an uuid3.</p>
<h4 id="rule-ipv6">uuid4</h4>
<p>The field under validation must be an uuid4.</p>
<h4 id="rule-ipv6">uuid5</h4>
<p>The field under validation must be an uuid5.</p>
<h4 id="rule-ipv6">uuid</h4>
<p>The field under validation must be an uuid.</p>
<h4 id="rule-url">url</h4>
<p>The field under validation must be a valid URL.</p>
<h4 id="rule-date">date</h4>
<p>The field under validation must be a valid date.</p>
<h4 id="rule-dateFormat">dateFormat=format</h4>
<p>The field under validation must match the given format. Example: <code>dateFormat=2006-01-02</code></p>
<h4 id="rule-after">after=date</h4>
<p>The field under validation must be a date after the given date. Supports relative dates.</p>
<pre>
type BookingForm struct {
    CheckIn time.Time `valid:"after=today"`
    Event   time.Time `valid:"after=tomorrow"`
    Future  time.Time `valid:"after=today+7d"`
}
</pre>
<p><strong>Relative date keywords:</strong></p>
<ul>
  <li><code>today</code> - Today at midnight</li>
  <li><code>tomorrow</code> - Tomorrow at midnight</li>
  <li><code>yesterday</code> - Yesterday at midnight</li>
  <li><code>now</code> - Current time</li>
  <li><code>today+7d</code> - 7 days from today</li>
  <li><code>today-1m</code> - 1 month ago</li>
  <li><code>today-18y</code> - 18 years ago</li>
</ul>
<h4 id="rule-afterOrEqual">afterOrEqual=date</h4>
<p>The field under validation must be a date after or equal to the given date. Supports relative dates.</p>
<h4 id="rule-before">before=date</h4>
<p>The field under validation must be a date before the given date. Supports relative dates.</p>
<pre>
type UserForm struct {
    BirthDate time.Time `valid:"before=today-18y"` // Must be 18 years or older
}
</pre>
<h4 id="rule-beforeOrEqual">beforeOrEqual=date</h4>
<p>The field under validation must be a date before or equal to the given date. Supports relative dates.</p>
<h4 id="rule-regex">regex=pattern</h4>
<p>The field under validation must match the given regular expression.</p>
<pre>
type Form struct {
    Code string `valid:"regex=^[A-Z]{3}-[0-9]{4}$"`
}
</pre>
<h4 id="rule-notRegex">notRegex=pattern</h4>
<p>The field under validation must not match the given regular expression.</p>
<h4 id="rule-in">in=value1|value2|...</h4>
<p>The field under validation must be included in the given list of values.</p>
<h4 id="rule-notIn">notIn=value1|value2|...</h4>
<p>The field under validation must not be included in the given list of values.</p>
<h4 id="rule-boolean">boolean</h4>
<p>The field under validation must be able to be cast as a boolean. Accepted input are true, false, 1, 0, "1", "0", "true", "false", "yes", "no", "on", "off".</p>
<h4 id="rule-accepted">accepted</h4>
<p>The field under validation must be "yes", "on", 1, or true.</p>
<h4 id="rule-declined">declined</h4>
<p>The field under validation must be "no", "off", 0, or false.</p>
<h4 id="rule-country">country, country.alpha2, country.alpha3, country.numeric</h4>
<p>The field under validation must be a valid ISO 3166-1 country code.</p>
<ul>
  <li><code>country</code> or <code>country.alpha2</code> - 2-letter code (e.g., US, GB, CN)</li>
  <li><code>country.alpha3</code> - 3-letter code (e.g., USA, GBR, CHN)</li>
  <li><code>country.numeric</code> - Numeric code (e.g., 840, 826, 156)</li>
</ul>
<pre>
type Address struct {
    Country string `valid:"required,country.alpha2"`
}
</pre>
<h4 id="rule-currency">currency, currency.fiat, currency.crypto</h4>
<p>The field under validation must be a valid currency code.</p>
<ul>
  <li><code>currency</code> or <code>currency.all</code> - Any currency (fiat + crypto)</li>
  <li><code>currency.fiat</code> - ISO 4217 3-letter code (e.g., USD, EUR, CNY)</li>
  <li><code>currency.crypto</code> - Cryptocurrency code (e.g., BTC, ETH, USDT, SOL)</li>
</ul>
<pre>
type Payment struct {
    Currency   string `valid:"required,currency"`        // fiat or crypto
    FiatOnly   string `valid:"required,currency.fiat"`   // fiat only
    CryptoOnly string `valid:"required,currency.crypto"` // crypto only
}
</pre>
<h4 id="rule-language">language, language.alpha2, language.alpha3</h4>
<p>The field under validation must be a valid ISO 639 language code.</p>
<ul>
  <li><code>language</code> or <code>language.alpha2</code> - 2-letter code (e.g., en, zh, ja)</li>
  <li><code>language.alpha3</code> - 3-letter code (e.g., eng, zho, jpn)</li>
</ul>
<pre>
type UserPreferences struct {
    Language string `valid:"required,language.alpha2"`
}
</pre>
<h4 id="rule-phone">phone, phone.e164</h4>
<p>The field under validation must be a valid phone number.</p>
<ul>
  <li><code>phone.e164</code> - E.164 format with country code (e.g., +14155551234)</li>
  <li><code>phone</code> - Any valid phone number format with country code</li>
</ul>
<pre>
type Contact struct {
    Phone string `valid:"required,phone.e164"`
}
</pre>
<h2>Type-Safe Date Validation (Alternative API)</h2>
<p>For those who prefer type-safe validation over struct tags, you can use the fluent API:</p>
<pre>
import "github.com/syssam/go-validator"

// Using standalone functions
if !validator.IsDateAfter(checkIn, validator.Today()) {
    return errors.New("check-in must be after today")
}

// Using fluent builder
rule := validator.Date("birth_date").Before(validator.T(validator.Today()).SubYears(18))
if err := rule.ValidateWithError(user.BirthDate); err != nil {
    return err
}

// Available helper functions
today := validator.Today()           // Today at midnight
tomorrow := validator.Tomorrow()     // Tomorrow at midnight
yesterday := validator.Yesterday()   // Yesterday at midnight

// Date arithmetic with TimeHelper
weekAgo := validator.T(validator.Today()).SubDays(7)
nextMonth := validator.T(validator.Today()).AddMonths(1)
years18Ago := validator.T(validator.Today()).SubYears(18)
</pre>
<h2>Custom Type Functions (Wrapper Type Support)</h2>
<p>Two mechanisms for unwrapping wrapper types before validation:</p>
<ul>
  <li><strong><code>RegisterAutoUnwrap</code></strong> — register a type matcher with method names for automatic unwrapping. The validator auto-builds optimized unwrappers with cached method indices. One registration handles all type parameter variants. Best for generic types like <code>graphql.Omittable[T]</code>.</li>
  <li><strong><code>RegisterCustomTypeFunc</code></strong> — explicit registration with direct type assertion. Fastest performance, best for types without <code>IsSet()/Value()</code> methods (e.g., <code>sql.NullString</code>).</li>
</ul>
<h4>gqlgen graphql.Omittable (RegisterAutoUnwrap — one registration for all variants)</h4>
<pre>
import (
    "reflect"
    "strings"
)

// One registration handles ALL Omittable[T] variants (300+ types).
// Specify the method names — the validator builds optimized unwrappers
// with cached method indices automatically.
validator.RegisterAutoUnwrap(
    func(t reflect.Type) bool {
        return strings.HasPrefix(t.Name(), "Omittable[") &&
            strings.Contains(t.PkgPath(), "gqlgen")
    },
    "IsSet", // Omittable.IsSet() bool — returns whether the value is set
    "Value", // Omittable.Value() T   — returns the inner value
)

type UpdateUserInput struct {
    Name  graphql.Omittable[*string] `json:"name" valid:"required"`
    Email graphql.Omittable[*string] `json:"email" valid:"required,email"`
    Age   graphql.Omittable[int]     `json:"age" valid:"min=0"`
}
</pre>
<h4>database/sql Null Types (RegisterCustomTypeFunc — fastest performance)</h4>
<pre>
import "database/sql"

// sql.Null* types use a .Valid field instead of IsSet()/Value() methods,
// so register each type explicitly for best performance.
validator.RegisterCustomTypeFunc(func(n sql.NullString) (any, bool) {
    if !n.Valid {
        return nil, false // NULL → skip validation
    }
    return n.String, true // valid → validate the string
})
validator.RegisterCustomTypeFunc(func(n sql.NullInt64) (any, bool) {
    if !n.Valid {
        return nil, false
    }
    return n.Int64, true
})
validator.RegisterCustomTypeFunc(func(n sql.NullFloat64) (any, bool) {
    if !n.Valid {
        return nil, false
    }
    return n.Float64, true
})

type User struct {
    Name  sql.NullString  `valid:"required"`
    Email sql.NullString  `valid:"required,email"`
    Age   sql.NullInt64   `valid:"min=18"`
    Score sql.NullFloat64 `valid:"between=0|100"`
}
</pre>
<p><strong>Behavior:</strong></p>
<ul>
  <li>When the wrapper reports "not set" (<code>false</code>): only <code>required</code> rules are checked — all other validation is skipped</li>
  <li>When the wrapper reports "set" (<code>true</code>): the inner value is extracted and validated normally with all rules</li>
  <li>Exact match (<code>RegisterCustomTypeFunc</code>) takes priority over <code>RegisterAutoUnwrap</code></li>
  <li>Auto-unwrap results are cached per type — matcher and method checks run only once per distinct type</li>
</ul>
<h2>Custom Validation Rules</h2>
<div class="highlight highlight-source-go">
  <pre>
  validator.CustomTypeTagMap.Set("customValidator", func CustomValidator(v reflect.Value, o reflect.Value, validTag *validator.ValidTag) bool {
    return false
  })
  </pre>
</div>
<h2>List of functions:</h2>
<h4>String Validation</h4>
<pre>
IsNumeric(str string) bool
IsInt(str string) bool
IsFloat(str string) bool
IsNull(str string) bool
IsEmail(str string) bool
IsAlpha(str string) bool
IsAlphaNum(str string) bool
IsAlphaDash(str string) bool
IsAlphaUnicode(str string) bool
IsAlphaNumUnicode(str string) bool
IsAlphaDashUnicode(str string) bool
IsIP(str string) bool
IsIPv4(str string) bool
IsIPv6(str string) bool
IsUUID3(str string) bool
IsUUID4(str string) bool
IsUUID5(str string) bool
IsUUID(str string) bool
IsURL(str string) bool
</pre>
<h4>Value Validation (Reflection-based)</h4>
<pre>
IsRequired(i interface{}) bool
IsBetween(i interface{}, params []string) (bool, error)
IsDigitsBetween(i interface{}, params []string) (bool, error)
IsMin(i interface{}, params []string) (bool, error)
IsMax(i interface{}, params []string) (bool, error)
IsSize(i interface{}, params []string) (bool, error)
IsDistinct(i interface{}) bool
IsGt(i interface{}, a interface{}) (bool, error)
IsGte(i interface{}, a interface{}) (bool, error)
IsLt(i interface{}, a interface{}) (bool, error)
IsLte(i interface{}, a interface{}) (bool, error)
</pre>
<h4>Date Validation</h4>
<pre>
Today() time.Time
Tomorrow() time.Time
Yesterday() time.Time
Now() time.Time
IsDateAfter(value, after time.Time) bool
IsDateAfterOrEqual(value, after time.Time) bool
IsDateBefore(value, before time.Time) bool
IsDateBeforeOrEqual(value, before time.Time) bool
IsDateBetween(value, start, end time.Time) bool
IsDateBetweenOrEqual(value, start, end time.Time) bool
IsToday(value time.Time) bool
IsFuture(value time.Time) bool
IsPast(value time.Time) bool
</pre>
<h4>Locale Validation (Country, Currency, Language, Phone)</h4>
<pre>
IsCountryCode(str string) bool
IsCountryAlpha2(str string) bool
IsCountryAlpha3(str string) bool
IsCountryNumeric(str string) bool
IsCurrencyCode(str string) bool
IsCurrencyNumeric(str string) bool
IsCurrency(str string) bool
IsLanguageCode(str string) bool
IsLanguageAlpha2(str string) bool
IsLanguageAlpha3(str string) bool
IsPhoneE164(str string) bool
IsPhone(str string, region string) bool
IsPhoneValid(str string) bool
IsPhoneMobile(str string, region string) bool
FormatPhoneE164(str string, region string) (string, error)
FormatPhoneNational(str string, region string) (string, error)
FormatPhoneInternational(str string, region string) (string, error)
</pre>
<h4>Generic Validation (Go 1.18+, Type-safe)</h4>
<pre>
IsGenericRequired[T comparable](value T) bool
IsGenericRequiredSlice[T any](value []T) bool
IsGenericRequiredMap[K comparable, V any](value map[K]V) bool
IsGenericMin[T NumericType](value, min T) bool
IsGenericMax[T NumericType](value, max T) bool
IsGenericBetween[T NumericType](value, min, max T) bool
IsGenericGt[T OrderedType](value, threshold T) bool
IsGenericGte[T OrderedType](value, threshold T) bool
IsGenericLt[T OrderedType](value, threshold T) bool
IsGenericLte[T OrderedType](value, threshold T) bool
IsGenericDistinct[T comparable](value []T) bool
IsGenericIn[T comparable](value T, allowed []T) bool
IsGenericNotIn[T comparable](value T, disallowed []T) bool
</pre>
