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
<h2>Type-Safe Date Validation (Alternative API)</h2>
<p>For those who prefer type-safe validation over struct tags, you can use the fluent API:</p>
<pre>
import "github.com/syssam/go-validator"

// Using standalone functions
if !validator.DateAfter(checkIn, validator.Today()) {
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
<h2>Custom Validation Rules</h2>
<div class="highlight highlight-source-go">
  <pre>
  validator.CustomTypeTagMap.Set("customValidator", func CustomValidator(v reflect.Value, o reflect.Value, validTag *validator.ValidTag) bool {
    return false
  })
  </pre>
</div>
<h2>List of functions:</h2>
<div class="highlight highlight-source-go">
  <pre>
    IsNumeric(str string) bool
    IsInt(str string) bool
    IsFloat(str string) bool
    IsNull(str string) bool
    ValidateBetween(i interface{}, params []string) (bool, error)
    ValidateDigitsBetween(i interface{}, params []string) (bool, error)
    ValidateDigitsBetweenInt64(value, left, right int64) bool
    ValidateDigitsBetweenFloat64(value, left, right float64) bool
    ValidateGt(i interface{}, a interface{}) (bool, error)
    ValidateGtFloat64(v, param float64) bool
    ValidateGte(i interface{}, a interface{}) (bool, error)
    ValidateGteFloat64(v, param float64) bool
    ValidateLt(i interface{}, a interface{}) (bool, error)
    ValidateLtFloat64(v, param float64) bool
    ValidateLte(i interface{}, a interface{}) (bool, error)
    ValidateLteFloat64(v, param float64) bool
    ValidateRequired(i interface{}) bool
    ValidateMin(i interface{}, params []string) (bool, error)
    ValidateMinFloat64(v, param float64) bool
    ValidateMax(i interface{}, params []string) (bool, error)
    ValidateMaxFloat64(v, param float64) bool
    ValidateSize(i interface{}, params []string) (bool, error)
    ValidateDistinct(i interface{}) bool
    ValidateEmail(str string) bool
    ValidateAlpha(str string) bool
    ValidateAlphaNum(str string) bool
    ValidateAlphaDash(str string) bool
    ValidateAlphaUnicode(str string) bool
    ValidateAlphaNumUnicode(str string) bool
    ValidateAlphaDashUnicode(str string) bool
    ValidateIP(str string) bool
    ValidateIPv4(str string) bool
    ValidateIPv6(str string) bool
    ValidateUUID3(str string) bool
    ValidateUUID4(str string) bool
    ValidateUUID5(str string) bool
    ValidateUUID(str string) bool
    ValidateURL(str string) bool

    // Date validation functions
    Today() time.Time
    Tomorrow() time.Time
    Yesterday() time.Time
    Now() time.Time
    DateAfter(value, after time.Time) bool
    DateAfterOrEqual(value, after time.Time) bool
    DateBefore(value, before time.Time) bool
    DateBeforeOrEqual(value, before time.Time) bool
    DateBetween(value, start, end time.Time) bool
    DateBetweenOrEqual(value, start, end time.Time) bool
    IsToday(value time.Time) bool
    IsFuture(value time.Time) bool
    IsPast(value time.Time) bool

    // Generic validation functions (Go 1.18+)
    IsMin[T NumericType](value, min T) bool
    IsMax[T NumericType](value, max T) bool
    IsBetween[T NumericType](value, min, max T) bool
    IsGt[T OrderedType](value, threshold T) bool
    IsGte[T OrderedType](value, threshold T) bool
    IsLt[T OrderedType](value, threshold T) bool
    IsLte[T OrderedType](value, threshold T) bool
    IsDistinct[T comparable](value []T) bool
    IsIn[T comparable](value T, allowed []T) bool
    IsNotIn[T comparable](value T, disallowed []T) bool
  </pre>
</div>
