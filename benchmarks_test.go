package validator

import "testing"

func BenchmarkFieldsRequired(t *testing.B) {
	model := FieldsRequired{Name: "TEST", Email: "test@example.com"}
	expected := true
	for i := 0; i < t.N; i++ {
		actual, err := ValidateStruct(&model)
		if actual != expected {
			t.Errorf("Expected validateateStruct(%q) to be %v, got %v", model, expected, actual)
			if err != nil {
				t.Errorf("Got Error on validateateStruct(%q): %s", model, err)
			}
		}
	}
}
