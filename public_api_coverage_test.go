package validator

import (
	"testing"
	"time"
)

// okBool asserts a plain bool result.
func okBool(t *testing.T, name string, got, want bool) {
	t.Helper()
	if got != want {
		t.Errorf("%s = %v, want %v", name, got, want)
	}
}

// okErr asserts a (bool, error) result: no error and the expected bool.
func okErr(t *testing.T, name string, got bool, err error, want bool) {
	t.Helper()
	if err != nil {
		t.Errorf("%s returned unexpected error: %v", name, err)
		return
	}
	if got != want {
		t.Errorf("%s = %v, want %v", name, got, want)
	}
}

func TestParamComparisonFuncs(t *testing.T) {
	b, err := IsGtParam(5, []string{"3"})
	okErr(t, "IsGtParam(5,3)", b, err, true)
	b, err = IsGtParam(2, []string{"3"})
	okErr(t, "IsGtParam(2,3)", b, err, false)

	b, err = IsGteParam(3, []string{"3"})
	okErr(t, "IsGteParam(3,3)", b, err, true)

	b, err = IsLtParam(2, []string{"3"})
	okErr(t, "IsLtParam(2,3)", b, err, true)

	b, err = IsLteParam(3, []string{"3"})
	okErr(t, "IsLteParam(3,3)", b, err, true)
}

func TestTwoValueComparisonFuncs(t *testing.T) {
	b, err := IsGt(5, 3)
	okErr(t, "IsGt(5,3)", b, err, true)
	b, err = IsGte(3, 3)
	okErr(t, "IsGte(3,3)", b, err, true)
	b, err = IsLt(2, 3)
	okErr(t, "IsLt(2,3)", b, err, true)
	b, err = IsLte(3, 3)
	okErr(t, "IsLte(3,3)", b, err, true)
	b, err = IsSame(3, 3)
	okErr(t, "IsSame(3,3)", b, err, true)

	b, err = IsDifferent(3, 4)
	okErr(t, "IsDifferent(3,4)", b, err, true)
	b, err = IsDifferent(3, 3)
	okErr(t, "IsDifferent(3,3)", b, err, false)
}

func TestDigitAndMultipleFuncs(t *testing.T) {
	b, err := IsMultipleOf(9, []string{"3"})
	okErr(t, "IsMultipleOf(9,3)", b, err, true)
	b, err = IsMultipleOf(10, []string{"3"})
	okErr(t, "IsMultipleOf(10,3)", b, err, false)

	b, err = IsMaxDigits(123, []string{"3"})
	okErr(t, "IsMaxDigits(123,3)", b, err, true)
	b, err = IsMaxDigits(1234, []string{"3"})
	okErr(t, "IsMaxDigits(1234,3)", b, err, false)

	b, err = IsMinDigits(123, []string{"3"})
	okErr(t, "IsMinDigits(123,3)", b, err, true)
	b, err = IsMinDigits(12, []string{"3"})
	okErr(t, "IsMinDigits(12,3)", b, err, false)

	b, err = IsDecimalPrecision(1.23, []string{"2"})
	okErr(t, "IsDecimalPrecision(1.23,2)", b, err, true)
	b, err = IsDecimalPrecision(1.2, []string{"2"})
	okErr(t, "IsDecimalPrecision(1.2,2)", b, err, false)
}

func TestContainsFuncs(t *testing.T) {
	b, err := IsContains("hello world", []string{"world"})
	okErr(t, "IsContains(str)", b, err, true)
	b, err = IsContains([]string{"a", "b", "c"}, []string{"b"})
	okErr(t, "IsContains(slice)", b, err, true)
	b, err = IsContains("hello", []string{"zzz"})
	okErr(t, "IsContains(miss)", b, err, false)

	b, err = IsDoesntContain("hello", []string{"zzz"})
	okErr(t, "IsDoesntContain(ok)", b, err, true)
	b, err = IsDoesntContain("hello", []string{"ell"})
	okErr(t, "IsDoesntContain(hit)", b, err, false)
}

func TestPresenceFuncs(t *testing.T) {
	okBool(t, "IsAccepted(yes)", IsAccepted("yes"), true)
	okBool(t, "IsAccepted(bool)", IsAccepted(true), true)
	okBool(t, "IsAccepted(no)", IsAccepted("no"), false)

	okBool(t, "IsDeclined(no)", IsDeclined("no"), true)
	okBool(t, "IsDeclined(bool)", IsDeclined(false), true)
	okBool(t, "IsDeclined(yes)", IsDeclined("yes"), false)

	okBool(t, "IsProhibited(empty)", IsProhibited(""), true)
	okBool(t, "IsProhibited(nonempty)", IsProhibited("x"), false)

	okBool(t, "IsMissing(empty)", IsMissing(""), true)
	okBool(t, "IsMissing(nonempty)", IsMissing("x"), false)

	okBool(t, "IsPresent(x)", IsPresent("x"), true)

	okBool(t, "IsRequired(x)", IsRequired("x"), true)
	okBool(t, "IsRequired(empty)", IsRequired(""), false)
}

func TestTypeAndCollectionFuncs(t *testing.T) {
	b, err := IsList([]int{1, 2, 3})
	okErr(t, "IsList(slice)", b, err, true)
	b, err = IsList("not a list")
	okErr(t, "IsList(str)", b, err, false)

	b, err = IsArray([]int{1})
	okErr(t, "IsArray(slice)", b, err, true)
	b, err = IsArray(5)
	okErr(t, "IsArray(int)", b, err, false)

	b, err = IsString("hello")
	okErr(t, "IsString(str)", b, err, true)
	b, err = IsString(5)
	okErr(t, "IsString(int)", b, err, false)

	b, err = IsFilled("x")
	okErr(t, "IsFilled(x)", b, err, true)
	b, err = IsFilled("")
	okErr(t, "IsFilled(empty)", b, err, false)

	b, err = IsRequiredArrayKeys(map[string]int{"a": 1, "b": 2}, []string{"a", "b"})
	okErr(t, "IsRequiredArrayKeys(ok)", b, err, true)
	b, err = IsRequiredArrayKeys(map[string]int{"a": 1}, []string{"a", "b"})
	okErr(t, "IsRequiredArrayKeys(miss)", b, err, false)
}

func TestJSONAndRegexFuncs(t *testing.T) {
	b, err := IsJSON(`{"a":1}`)
	okErr(t, "IsJSON(obj)", b, err, true)
	b, err = IsJSON(`{bad json`)
	okErr(t, "IsJSON(bad)", b, err, false)

	b, err = IsBoolean(true)
	okErr(t, "IsBoolean(bool)", b, err, true)

	b, err = IsRegex("abc123", []string{"^[a-z]+[0-9]+$"})
	okErr(t, "IsRegex(match)", b, err, true)
	b, err = IsRegex("ABC", []string{"^[a-z]+$"})
	okErr(t, "IsRegex(nomatch)", b, err, false)

	b, err = IsNotRegex("ABC", []string{"^[a-z]+$"})
	okErr(t, "IsNotRegex(ok)", b, err, true)
	b, err = IsNotRegex("abc", []string{"^[a-z]+$"})
	okErr(t, "IsNotRegex(hit)", b, err, false)
}

func TestDateValidationFuncs(t *testing.T) {
	b, err := IsDate("2020-01-02")
	okErr(t, "IsDate(valid)", b, err, true)
	b, err = IsDate("not-a-date")
	okErr(t, "IsDate(invalid)", b, err, false)

	b, err = IsDateFormat("2020-01-02", []string{"2006-01-02"})
	okErr(t, "IsDateFormat(match)", b, err, true)

	b, err = IsAfter("2020-06-01", []string{"2020-01-01"})
	okErr(t, "IsAfter(true)", b, err, true)
	b, err = IsAfter("2019-06-01", []string{"2020-01-01"})
	okErr(t, "IsAfter(false)", b, err, false)

	b, err = IsAfterOrEqual("2020-01-01", []string{"2020-01-01"})
	okErr(t, "IsAfterOrEqual(eq)", b, err, true)

	b, err = IsBefore("2019-06-01", []string{"2020-01-01"})
	okErr(t, "IsBefore(true)", b, err, true)

	b, err = IsBeforeOrEqual("2020-01-01", []string{"2020-01-01"})
	okErr(t, "IsBeforeOrEqual(eq)", b, err, true)
}

func TestInt64Comparators(t *testing.T) {
	okBool(t, "IsInt64Gt", IsInt64Gt(5, 3), true)
	okBool(t, "IsInt64Gte", IsInt64Gte(3, 3), true)
	okBool(t, "IsInt64Lt", IsInt64Lt(2, 3), true)
	okBool(t, "IsInt64Lte", IsInt64Lte(3, 3), true)
}

func TestUint64Comparators(t *testing.T) {
	okBool(t, "IsUint64Gt", IsUint64Gt(5, 3), true)
	okBool(t, "IsUint64Gte", IsUint64Gte(3, 3), true)
	okBool(t, "IsUint64Lt", IsUint64Lt(2, 3), true)
	okBool(t, "IsUint64Lte", IsUint64Lte(3, 3), true)
}

func TestStandaloneDateHelpers(t *testing.T) {
	base := time.Date(2020, 6, 15, 0, 0, 0, 0, time.UTC)
	earlier := base.AddDate(0, 0, -1)
	later := base.AddDate(0, 0, 1)

	okBool(t, "IsDateAfterOrEqual(after)", IsDateAfterOrEqual(later, base), true)
	okBool(t, "IsDateAfterOrEqual(equal)", IsDateAfterOrEqual(base, base), true)
	okBool(t, "IsDateBeforeOrEqual(before)", IsDateBeforeOrEqual(earlier, base), true)
	okBool(t, "IsDateBeforeOrEqual(equal)", IsDateBeforeOrEqual(base, base), true)
	okBool(t, "IsDateBetweenOrEqual(mid)", IsDateBetweenOrEqual(base, earlier, later), true)
	okBool(t, "IsDateBetweenOrEqual(edge)", IsDateBetweenOrEqual(earlier, earlier, later), true)

	// TodayUTC and Now exercise the clock helpers.
	if TodayUTC().Location() != time.UTC {
		t.Error("TodayUTC should be in UTC")
	}
	if Now().IsZero() {
		t.Error("Now should not be zero")
	}
}

func TestDateRuleFluentBuilder(t *testing.T) {
	base := time.Date(2020, 6, 15, 12, 0, 0, 0, time.UTC)
	lo := base.AddDate(0, 0, -1)
	hi := base.AddDate(0, 0, 1)

	okBool(t, "DateRule.BeforeOrEqual", Date().BeforeOrEqual(base).Validate(base), true)
	okBool(t, "DateRule.AfterOrEqual", Date().AfterOrEqual(base).Validate(base), true)
	okBool(t, "DateRule.BetweenOrEqual", Date().BetweenOrEqual(lo, hi).Validate(base), true)
	okBool(t, "DateRule.BetweenOrEqual(out)", Date().BetweenOrEqual(lo, hi).Validate(hi.AddDate(0, 0, 1)), false)
}

func TestTimeHelperArithmetic(t *testing.T) {
	base := time.Date(2020, 1, 15, 0, 0, 0, 0, time.UTC)

	if got := T(base).AddMonths(2); got.Month() != time.March {
		t.Errorf("AddMonths(2) month = %v, want March", got.Month())
	}
	if got := T(base).AddYears(3); got.Year() != 2023 {
		t.Errorf("AddYears(3) year = %d, want 2023", got.Year())
	}
}
