package validator

import (
	"testing"
)

// Test complex nested struct validation
type ComplexNestedStruct struct {
	User   UserInfo    `valid:"required"`
	Config ConfigInfo  `valid:"required"`
	Items  []ItemInfo  `valid:"required,distinct"`
	Meta   interface{} // Interface field
}

type UserInfo struct {
	Name     string `valid:"required,min=2"`
	Email    string `valid:"required,email"`
	Age      int    `valid:"min=18,max=120"`
	Settings *UserSettings
}

type UserSettings struct {
	Theme    string `valid:"required"`
	Language string `valid:"required,size=2"`
}

type ConfigInfo struct {
	Version string `valid:"required"`
	Debug   bool
	Timeout int `valid:"gt=0,lt=3600"`
}

type ItemInfo struct {
	ID    string `valid:"required"`
	Value string `valid:"min=1"`
}

// Test complex validation scenarios
func TestComplexStructValidation(t *testing.T) {
	tests := []struct {
		name     string
		data     ComplexNestedStruct
		expected bool
	}{
		{
			"Valid complex struct",
			ComplexNestedStruct{
				User: UserInfo{
					Name:  "John Doe",
					Email: "john@example.com",
					Age:   25,
					Settings: &UserSettings{
						Theme:    "dark",
						Language: "en",
					},
				},
				Config: ConfigInfo{
					Version: "1.0.0",
					Debug:   true,
					Timeout: 30,
				},
				Items: []ItemInfo{
					{ID: "item1", Value: "value1"},
					{ID: "item2", Value: "value2"},
				},
				Meta: map[string]interface{}{"key": "value"},
			},
			true,
		},
		{
			"Invalid email",
			ComplexNestedStruct{
				User: UserInfo{
					Name:  "John Doe",
					Email: "invalid-email",
					Age:   25,
				},
				Config: ConfigInfo{
					Version: "1.0.0",
					Timeout: 30,
				},
				Items: []ItemInfo{
					{ID: "item1", Value: "value1"},
				},
			},
			false,
		},
		{
			"Age below minimum",
			ComplexNestedStruct{
				User: UserInfo{
					Name:  "John Doe",
					Email: "john@example.com",
					Age:   16, // Below minimum
				},
				Config: ConfigInfo{
					Version: "1.0.0",
					Timeout: 30,
				},
				Items: []ItemInfo{
					{ID: "item1", Value: "value1"},
				},
			},
			false,
		},
		{
			"Missing required fields",
			ComplexNestedStruct{
				User: UserInfo{
					Name: "", // Required field missing
					Age:  25,
				},
				Config: ConfigInfo{
					Timeout: 30,
				},
				Items: []ItemInfo{},
			},
			false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateStruct(test.data)
			actual := err == nil
			if actual != test.expected {
				t.Errorf("Expected %t for %s, got %t. Error: %v", test.expected, test.name, actual, err)
			}
		})
	}
}

// Test pointer field validation
type PointerStruct struct {
	Required *string `valid:"required"`
	Optional *string
	Number   *int `valid:"min=0"`
}

func TestPointerValidation(t *testing.T) {
	validString := "valid"
	invalidNumber := -5
	validNumber := 10

	tests := []struct {
		name     string
		data     PointerStruct
		expected bool
	}{
		{
			"Valid pointers",
			PointerStruct{
				Required: &validString,
				Number:   &validNumber,
			},
			true,
		},
		{
			"Missing required pointer",
			PointerStruct{
				Required: nil,
				Number:   &validNumber,
			},
			false,
		},
		{
			"Invalid number value",
			PointerStruct{
				Required: &validString,
				Number:   &invalidNumber,
			},
			false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateStruct(test.data)
			actual := err == nil
			if actual != test.expected {
				t.Errorf("Expected %t for %s, got %t. Error: %v", test.expected, test.name, actual, err)
			}
		})
	}
}

// Test slice and map validation
type CollectionStruct struct {
	StringSlice []string          `valid:"required,min=1"`
	IntSlice    []int             `valid:"max=5"`
	StringMap   map[string]string `valid:"required"`
	IntMap      map[string]int    `valid:"min=1,max=10"`
}

func TestCollectionValidation(t *testing.T) {
	tests := []struct {
		name     string
		data     CollectionStruct
		expected bool
	}{
		{
			"Valid collections",
			CollectionStruct{
				StringSlice: []string{"a", "b", "c"},
				IntSlice:    []int{1, 2, 3},
				StringMap:   map[string]string{"key": "value"},
				IntMap:      map[string]int{"count": 5},
			},
			true,
		},
		{
			"Empty required slice",
			CollectionStruct{
				StringSlice: []string{}, // Required but empty
				StringMap:   map[string]string{"key": "value"},
				IntMap:      map[string]int{"count": 5},
			},
			false,
		},
		{
			"Slice too large",
			CollectionStruct{
				StringSlice: []string{"a", "b", "c"},
				IntSlice:    []int{1, 2, 3, 4, 5, 6}, // Max 5
				StringMap:   map[string]string{"key": "value"},
				IntMap:      map[string]int{"count": 5},
			},
			false,
		},
		{
			"Map too large",
			CollectionStruct{
				StringSlice: []string{"a", "b", "c"},
				IntSlice:    []int{1, 2, 3},
				StringMap:   map[string]string{"key": "value"},
				IntMap: map[string]int{
					"1": 1, "2": 2, "3": 3, "4": 4, "5": 5,
					"6": 6, "7": 7, "8": 8, "9": 9, "10": 10, "11": 11, // Max 10
				},
			},
			false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateStruct(test.data)
			actual := err == nil
			if actual != test.expected {
				t.Errorf("Expected %t for %s, got %t. Error: %v", test.expected, test.name, actual, err)
			}
		})
	}
}

// Test omitempty functionality
type OmitEmptyStruct struct {
	Required string `valid:"required"`
	Optional string `valid:"omitempty,min=3"`
	Number   int    `valid:"omitempty,gt=0"`
}

func TestOmitEmptyValidation(t *testing.T) {
	tests := []struct {
		name     string
		data     OmitEmptyStruct
		expected bool
	}{
		{
			"Valid with optional fields",
			OmitEmptyStruct{
				Required: "value",
				Optional: "test",
				Number:   5,
			},
			true,
		},
		{
			"Valid with empty optional fields",
			OmitEmptyStruct{
				Required: "value",
				Optional: "", // Empty but omitempty
				Number:   0,  // Zero but omitempty
			},
			true,
		},
		{
			"Invalid optional field value",
			OmitEmptyStruct{
				Required: "value",
				Optional: "ab", // Too short when provided
				Number:   0,
			},
			false,
		},
		{
			"Missing required field",
			OmitEmptyStruct{
				Required: "", // Required field is empty
				Optional: "test",
				Number:   5,
			},
			false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := ValidateStruct(test.data)
			actual := err == nil
			if actual != test.expected {
				t.Errorf("Expected %t for %s, got %t. Error: %v", test.expected, test.name, actual, err)
			}
		})
	}
}
