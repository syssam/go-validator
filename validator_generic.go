package validator

// Type constraints for generics
type (
	IntegerType interface {
		~int | ~int8 | ~int16 | ~int32 | ~int64 |
			~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
	}
	FloatType interface {
		~float32 | ~float64
	}
	NumericType interface {
		IntegerType | FloatType
	}
	OrderedType interface {
		IntegerType | FloatType | ~string
	}
)

// =============================================================================
// Generic validation functions - use generics for multi-type support
// =============================================================================

// IsRequired checks if a comparable value is not zero.
func IsRequired[T comparable](value T) bool {
	var zero T
	return value != zero
}

// IsRequiredSlice checks if a slice is not empty.
func IsRequiredSlice[T any](value []T) bool {
	return len(value) > 0
}

// IsRequiredMap checks if a map is not empty.
func IsRequiredMap[K comparable, V any](value map[K]V) bool {
	return len(value) > 0
}

// IsMin checks if value >= min.
func IsMin[T NumericType](value, min T) bool {
	return value >= min
}

// IsMax checks if value <= max.
func IsMax[T NumericType](value, max T) bool {
	return value <= max
}

// IsBetween checks if min <= value <= max.
func IsBetween[T NumericType](value, min, max T) bool {
	return value >= min && value <= max
}

// IsGt checks if value > threshold.
func IsGt[T OrderedType](value, threshold T) bool {
	return value > threshold
}

// IsGte checks if value >= threshold.
func IsGte[T OrderedType](value, threshold T) bool {
	return value >= threshold
}

// IsLt checks if value < threshold.
func IsLt[T OrderedType](value, threshold T) bool {
	return value < threshold
}

// IsLte checks if value <= threshold.
func IsLte[T OrderedType](value, threshold T) bool {
	return value <= threshold
}

// IsDistinct checks if all slice elements are unique.
func IsDistinct[T comparable](value []T) bool {
	seen := make(map[T]struct{}, len(value))
	for _, v := range value {
		if _, exists := seen[v]; exists {
			return false
		}
		seen[v] = struct{}{}
	}
	return true
}

// IsIn checks if value is in allowed list.
func IsIn[T comparable](value T, allowed []T) bool {
	for _, a := range allowed {
		if value == a {
			return true
		}
	}
	return false
}

// IsNotIn checks if value is not in disallowed list.
func IsNotIn[T comparable](value T, disallowed []T) bool {
	for _, d := range disallowed {
		if value == d {
			return false
		}
	}
	return true
}
