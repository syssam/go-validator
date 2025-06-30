package validator

import (
	"math"
	"testing"
)

func TestToString(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected string
	}{
		{"hello", "hello"},
		{123, "123"},
		{int64(456), "456"},
		{uint64(789), "789"},
		{float64(3.14), "3.14"},
		{true, "true"},
		{false, "false"},
		{[]int{1, 2, 3}, "[1 2 3]"},
		{nil, "<nil>"},
		{struct{ Name string }{"test"}, "{test}"},
	}

	for _, test := range tests {
		result := ToString(test.input)
		if result != test.expected {
			t.Errorf("ToString(%v) = %s; expected %s", test.input, result, test.expected)
		}
	}
}

func TestToFloat(t *testing.T) {
	tests := []struct {
		input       string
		expected    float64
		expectError bool
	}{
		{"3.14", 3.14, false},
		{"0", 0.0, false},
		{"-2.5", -2.5, false},
		{"invalid", 0.0, true},
		{"", 0.0, true},
		{"1.23e10", 1.23e10, false},
	}

	for _, test := range tests {
		result, err := ToFloat(test.input)
		if result != test.expected {
			t.Errorf("ToFloat(%s) = %f; expected %f", test.input, result, test.expected)
		}
		if test.expectError && err == nil {
			t.Errorf("ToFloat(%s) expected error but got nil", test.input)
		}
		if !test.expectError && err != nil {
			t.Errorf("ToFloat(%s) unexpected error: %v", test.input, err)
		}
	}
}

func TestToBool(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"1", true},
		{"false", false},
		{"0", false},
		{"", false},
		{"yes", false},
		{"TRUE", false},
	}

	for _, test := range tests {
		result := ToBool(test.input)
		if result != test.expected {
			t.Errorf("ToBool(%s) = %t; expected %t", test.input, result, test.expected)
		}
	}
}

func TestToInt(t *testing.T) {
	tests := []struct {
		input       interface{}
		expected    int64
		expectError bool
	}{
		{int(123), 123, false},
		{int8(45), 45, false},
		{int16(1234), 1234, false},
		{int32(56789), 56789, false},
		{int64(9876543210), 9876543210, false},
		{uint(123), 123, false},
		{uint8(255), 255, false},
		{uint16(65535), 65535, false},
		{uint32(4294967295), 4294967295, false},
		{uint64(9223372036854775807), 9223372036854775807, false},
		{uint64(math.MaxUint64), 0, true}, // Exceeds int64 max
		{uint(math.MaxUint64), 0, true},   // Exceeds int64 max on 64-bit systems
		{"123", 123, false},
		{"-456", -456, false},
		{"invalid", 0, true},
		{3.14, 0, true}, // Unsupported type
	}

	for _, test := range tests {
		result, err := ToInt(test.input)
		if result != test.expected {
			t.Errorf("ToInt(%v) = %d; expected %d", test.input, result, test.expected)
		}
		if test.expectError && err == nil {
			t.Errorf("ToInt(%v) expected error but got nil", test.input)
		}
		if !test.expectError && err != nil {
			t.Errorf("ToInt(%v) unexpected error: %v", test.input, err)
		}
	}
}

func TestToUint(t *testing.T) {
	tests := []struct {
		input       string
		expected    uint64
		expectError bool
	}{
		{"123", 123, false},
		{"0", 0, false},
		{"18446744073709551615", math.MaxUint64, false},
		{"-1", 0, true},
		{"invalid", 0, true},
		{"", 0, true},
		{"0x10", 16, false}, // Hex
		{"010", 8, false},   // Octal
	}

	for _, test := range tests {
		result, err := ToUint(test.input)
		if result != test.expected {
			t.Errorf("ToUint(%s) = %d; expected %d", test.input, result, test.expected)
		}
		if test.expectError && err == nil {
			t.Errorf("ToUint(%s) expected error but got nil", test.input)
		}
		if !test.expectError && err != nil {
			t.Errorf("ToUint(%s) unexpected error: %v", test.input, err)
		}
	}
}
