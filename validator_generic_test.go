package validator

import (
	"testing"
)

func TestIsGenericRequired(t *testing.T) {
	// string
	if IsGenericRequired("") {
		t.Error("empty string should be false")
	}
	if !IsGenericRequired("hello") {
		t.Error("non-empty string should be true")
	}

	// int
	if IsGenericRequired(0) {
		t.Error("zero int should be false")
	}
	if !IsGenericRequired(42) {
		t.Error("non-zero int should be true")
	}

	// pointer
	var p *int
	if IsGenericRequired(p) {
		t.Error("nil pointer should be false")
	}
	val := 1
	if !IsGenericRequired(&val) {
		t.Error("non-nil pointer should be true")
	}
}

func TestIsGenericRequiredSlice(t *testing.T) {
	if IsGenericRequiredSlice([]string(nil)) {
		t.Error("nil slice should be false")
	}
	if IsGenericRequiredSlice([]string{}) {
		t.Error("empty slice should be false")
	}
	if !IsGenericRequiredSlice([]string{"a"}) {
		t.Error("non-empty slice should be true")
	}
}

func TestIsGenericRequiredMap(t *testing.T) {
	if IsGenericRequiredMap(map[string]int(nil)) {
		t.Error("nil map should be false")
	}
	if IsGenericRequiredMap(map[string]int{}) {
		t.Error("empty map should be false")
	}
	if !IsGenericRequiredMap(map[string]int{"a": 1}) {
		t.Error("non-empty map should be true")
	}
}

func TestIsGenericMin(t *testing.T) {
	// int
	if !IsGenericMin(10, 5) {
		t.Error("10 >= 5 should be true")
	}
	if !IsGenericMin(5, 5) {
		t.Error("5 >= 5 should be true")
	}
	if IsGenericMin(3, 5) {
		t.Error("3 >= 5 should be false")
	}

	// float64
	if !IsGenericMin(3.14, 3.0) {
		t.Error("3.14 >= 3.0 should be true")
	}

	// int8
	if !IsGenericMin(int8(10), int8(5)) {
		t.Error("int8: 10 >= 5 should be true")
	}
}

func TestIsGenericMax(t *testing.T) {
	if !IsGenericMax(3, 5) {
		t.Error("3 <= 5 should be true")
	}
	if !IsGenericMax(5, 5) {
		t.Error("5 <= 5 should be true")
	}
	if IsGenericMax(10, 5) {
		t.Error("10 <= 5 should be false")
	}

	// float32
	if !IsGenericMax(float32(3.0), float32(5.0)) {
		t.Error("float32: 3.0 <= 5.0 should be true")
	}
}

func TestIsGenericBetween(t *testing.T) {
	if !IsGenericBetween(5, 1, 10) {
		t.Error("5 between 1-10 should be true")
	}
	if !IsGenericBetween(1, 1, 10) {
		t.Error("1 between 1-10 should be true")
	}
	if !IsGenericBetween(10, 1, 10) {
		t.Error("10 between 1-10 should be true")
	}
	if IsGenericBetween(0, 1, 10) {
		t.Error("0 between 1-10 should be false")
	}
	if IsGenericBetween(11, 1, 10) {
		t.Error("11 between 1-10 should be false")
	}

	// float64
	if !IsGenericBetween(5.5, 1.0, 10.0) {
		t.Error("5.5 between 1.0-10.0 should be true")
	}
}

func TestIsGenericGtGteLtLte(t *testing.T) {
	// Gt - int
	if !IsGenericGt(6, 5) {
		t.Error("6 > 5 should be true")
	}
	if IsGenericGt(5, 5) {
		t.Error("5 > 5 should be false")
	}

	// Gte - int
	if !IsGenericGte(6, 5) {
		t.Error("6 >= 5 should be true")
	}
	if !IsGenericGte(5, 5) {
		t.Error("5 >= 5 should be true")
	}
	if IsGenericGte(4, 5) {
		t.Error("4 >= 5 should be false")
	}

	// Lt - int
	if !IsGenericLt(4, 5) {
		t.Error("4 < 5 should be true")
	}
	if IsGenericLt(5, 5) {
		t.Error("5 < 5 should be false")
	}

	// Lte - int
	if !IsGenericLte(4, 5) {
		t.Error("4 <= 5 should be true")
	}
	if !IsGenericLte(5, 5) {
		t.Error("5 <= 5 should be true")
	}
	if IsGenericLte(6, 5) {
		t.Error("6 <= 5 should be false")
	}

	// String comparison (OrderedType includes string)
	if !IsGenericGt("b", "a") {
		t.Error("b > a should be true")
	}
	if !IsGenericLt("a", "b") {
		t.Error("a < b should be true")
	}
	if !IsGenericGte("abc", "abc") {
		t.Error("abc >= abc should be true")
	}

	// float64
	if !IsGenericGt(3.14, 3.0) {
		t.Error("3.14 > 3.0 should be true")
	}
}

func TestIsGenericDistinct(t *testing.T) {
	// int
	if !IsGenericDistinct([]int{1, 2, 3, 4, 5}) {
		t.Error("unique int slice should be true")
	}
	if !IsGenericDistinct([]int{}) {
		t.Error("empty slice should be true")
	}
	if IsGenericDistinct([]int{1, 2, 3, 2, 5}) {
		t.Error("duplicate int slice should be false")
	}

	// string
	if !IsGenericDistinct([]string{"a", "b", "c"}) {
		t.Error("unique string slice should be true")
	}
	if IsGenericDistinct([]string{"a", "b", "a"}) {
		t.Error("duplicate string slice should be false")
	}

	// float64
	if !IsGenericDistinct([]float64{1.1, 2.2, 3.3}) {
		t.Error("unique float64 slice should be true")
	}
}

func TestIsGenericIn(t *testing.T) {
	// string
	allowed := []string{"pending", "active", "completed"}
	if !IsGenericIn("active", allowed) {
		t.Error("active in list should be true")
	}
	if IsGenericIn("deleted", allowed) {
		t.Error("deleted not in list should be false")
	}

	// int
	if !IsGenericIn(2, []int{1, 2, 3}) {
		t.Error("2 in [1,2,3] should be true")
	}
	if IsGenericIn(5, []int{1, 2, 3}) {
		t.Error("5 in [1,2,3] should be false")
	}
}

func TestIsGenericNotIn(t *testing.T) {
	// string
	disallowed := []string{"banned", "deleted"}
	if !IsGenericNotIn("active", disallowed) {
		t.Error("active not in disallowed should be true")
	}
	if IsGenericNotIn("banned", disallowed) {
		t.Error("banned in disallowed should be false")
	}

	// int
	if !IsGenericNotIn(5, []int{1, 2, 3}) {
		t.Error("5 not in [1,2,3] should be true")
	}
	if IsGenericNotIn(2, []int{1, 2, 3}) {
		t.Error("2 not in [1,2,3] should be false")
	}
}

// Benchmarks
func BenchmarkIsGenericRequired(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsGenericRequired("hello")
	}
}

func BenchmarkIsGenericMin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		IsGenericMin(100, 18)
	}
}

func BenchmarkIsGenericDistinct(b *testing.B) {
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IsGenericDistinct(slice)
	}
}

func BenchmarkIsGenericIn(b *testing.B) {
	allowed := []string{"pending", "active", "completed", "archived"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		IsGenericIn("active", allowed)
	}
}
