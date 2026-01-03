package validator

import (
	"testing"
)

func TestIsRequired(t *testing.T) {
	// string
	if IsRequired("") {
		t.Error("empty string should be false")
	}
	if !IsRequired("hello") {
		t.Error("non-empty string should be true")
	}

	// int
	if IsRequired(0) {
		t.Error("zero int should be false")
	}
	if !IsRequired(42) {
		t.Error("non-zero int should be true")
	}

	// pointer
	var p *int
	if IsRequired(p) {
		t.Error("nil pointer should be false")
	}
	val := 1
	if !IsRequired(&val) {
		t.Error("non-nil pointer should be true")
	}
}

func TestIsRequiredSlice(t *testing.T) {
	if IsRequiredSlice([]string(nil)) {
		t.Error("nil slice should be false")
	}
	if IsRequiredSlice([]string{}) {
		t.Error("empty slice should be false")
	}
	if !IsRequiredSlice([]string{"a"}) {
		t.Error("non-empty slice should be true")
	}
}

func TestIsRequiredMap(t *testing.T) {
	if IsRequiredMap(map[string]int(nil)) {
		t.Error("nil map should be false")
	}
	if IsRequiredMap(map[string]int{}) {
		t.Error("empty map should be false")
	}
	if !IsRequiredMap(map[string]int{"a": 1}) {
		t.Error("non-empty map should be true")
	}
}

func TestIsMin(t *testing.T) {
	// int
	if !IsMin(10, 5) {
		t.Error("10 >= 5 should be true")
	}
	if !IsMin(5, 5) {
		t.Error("5 >= 5 should be true")
	}
	if IsMin(3, 5) {
		t.Error("3 >= 5 should be false")
	}

	// float64
	if !IsMin(3.14, 3.0) {
		t.Error("3.14 >= 3.0 should be true")
	}

	// int8
	if !IsMin(int8(10), int8(5)) {
		t.Error("int8: 10 >= 5 should be true")
	}
}

func TestIsMax(t *testing.T) {
	if !IsMax(3, 5) {
		t.Error("3 <= 5 should be true")
	}
	if !IsMax(5, 5) {
		t.Error("5 <= 5 should be true")
	}
	if IsMax(10, 5) {
		t.Error("10 <= 5 should be false")
	}

	// float32
	if !IsMax(float32(3.0), float32(5.0)) {
		t.Error("float32: 3.0 <= 5.0 should be true")
	}
}

func TestIsBetween(t *testing.T) {
	if !IsBetween(5, 1, 10) {
		t.Error("5 between 1-10 should be true")
	}
	if !IsBetween(1, 1, 10) {
		t.Error("1 between 1-10 should be true")
	}
	if !IsBetween(10, 1, 10) {
		t.Error("10 between 1-10 should be true")
	}
	if IsBetween(0, 1, 10) {
		t.Error("0 between 1-10 should be false")
	}
	if IsBetween(11, 1, 10) {
		t.Error("11 between 1-10 should be false")
	}

	// float64
	if !IsBetween(5.5, 1.0, 10.0) {
		t.Error("5.5 between 1.0-10.0 should be true")
	}
}

func TestIsGtGteLtLte(t *testing.T) {
	// Gt - int
	if !IsGt(6, 5) {
		t.Error("6 > 5 should be true")
	}
	if IsGt(5, 5) {
		t.Error("5 > 5 should be false")
	}

	// Gte - int
	if !IsGte(6, 5) {
		t.Error("6 >= 5 should be true")
	}
	if !IsGte(5, 5) {
		t.Error("5 >= 5 should be true")
	}
	if IsGte(4, 5) {
		t.Error("4 >= 5 should be false")
	}

	// Lt - int
	if !IsLt(4, 5) {
		t.Error("4 < 5 should be true")
	}
	if IsLt(5, 5) {
		t.Error("5 < 5 should be false")
	}

	// Lte - int
	if !IsLte(4, 5) {
		t.Error("4 <= 5 should be true")
	}
	if !IsLte(5, 5) {
		t.Error("5 <= 5 should be true")
	}
	if IsLte(6, 5) {
		t.Error("6 <= 5 should be false")
	}

	// String comparison (OrderedType includes string)
	if !IsGt("b", "a") {
		t.Error("b > a should be true")
	}
	if !IsLt("a", "b") {
		t.Error("a < b should be true")
	}
	if !IsGte("abc", "abc") {
		t.Error("abc >= abc should be true")
	}

	// float64
	if !IsGt(3.14, 3.0) {
		t.Error("3.14 > 3.0 should be true")
	}
}

func TestIsDistinct(t *testing.T) {
	// int
	if !IsDistinct([]int{1, 2, 3, 4, 5}) {
		t.Error("unique int slice should be true")
	}
	if !IsDistinct([]int{}) {
		t.Error("empty slice should be true")
	}
	if IsDistinct([]int{1, 2, 3, 2, 5}) {
		t.Error("duplicate int slice should be false")
	}

	// string
	if !IsDistinct([]string{"a", "b", "c"}) {
		t.Error("unique string slice should be true")
	}
	if IsDistinct([]string{"a", "b", "a"}) {
		t.Error("duplicate string slice should be false")
	}

	// float64
	if !IsDistinct([]float64{1.1, 2.2, 3.3}) {
		t.Error("unique float64 slice should be true")
	}
}

func TestIsIn(t *testing.T) {
	// string
	allowed := []string{"pending", "active", "completed"}
	if !IsIn("active", allowed) {
		t.Error("active in list should be true")
	}
	if IsIn("deleted", allowed) {
		t.Error("deleted not in list should be false")
	}

	// int
	if !IsIn(2, []int{1, 2, 3}) {
		t.Error("2 in [1,2,3] should be true")
	}
	if IsIn(5, []int{1, 2, 3}) {
		t.Error("5 in [1,2,3] should be false")
	}
}

func TestIsNotIn(t *testing.T) {
	// string
	disallowed := []string{"banned", "deleted"}
	if !IsNotIn("active", disallowed) {
		t.Error("active not in disallowed should be true")
	}
	if IsNotIn("banned", disallowed) {
		t.Error("banned in disallowed should be false")
	}

	// int
	if !IsNotIn(5, []int{1, 2, 3}) {
		t.Error("5 not in [1,2,3] should be true")
	}
	if IsNotIn(2, []int{1, 2, 3}) {
		t.Error("2 not in [1,2,3] should be false")
	}
}

// Benchmarks
func BenchmarkIsRequired(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsRequired("hello")
	}
}

func BenchmarkIsMin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsMin(100, 18)
	}
}

func BenchmarkIsDistinct(b *testing.B) {
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IsDistinct(slice)
	}
}

func BenchmarkIsIn(b *testing.B) {
	allowed := []string{"pending", "active", "completed", "archived"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IsIn("active", allowed)
	}
}
