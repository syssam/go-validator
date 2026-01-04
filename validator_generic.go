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
// Generic validation functions - use generics for type-safe validation
// These use the "IsGeneric" prefix to distinguish from reflection-based versions
// =============================================================================

// IsGenericRequired checks if a comparable value is not zero.
func IsGenericRequired[T comparable](value T) bool {
	var zero T
	return value != zero
}

// IsGenericRequiredSlice checks if a slice is not empty.
func IsGenericRequiredSlice[T any](value []T) bool {
	return len(value) > 0
}

// IsGenericRequiredMap checks if a map is not empty.
func IsGenericRequiredMap[K comparable, V any](value map[K]V) bool {
	return len(value) > 0
}

// IsGenericMin checks if value >= min.
func IsGenericMin[T NumericType](value, min T) bool {
	return value >= min
}

// IsGenericMax checks if value <= max.
func IsGenericMax[T NumericType](value, max T) bool {
	return value <= max
}

// IsGenericBetween checks if min <= value <= max.
func IsGenericBetween[T NumericType](value, min, max T) bool {
	return value >= min && value <= max
}

// IsGenericGt checks if value > threshold.
func IsGenericGt[T OrderedType](value, threshold T) bool {
	return value > threshold
}

// IsGenericGte checks if value >= threshold.
func IsGenericGte[T OrderedType](value, threshold T) bool {
	return value >= threshold
}

// IsGenericLt checks if value < threshold.
func IsGenericLt[T OrderedType](value, threshold T) bool {
	return value < threshold
}

// IsGenericLte checks if value <= threshold.
func IsGenericLte[T OrderedType](value, threshold T) bool {
	return value <= threshold
}

// IsGenericDistinct checks if all slice elements are unique.
func IsGenericDistinct[T comparable](value []T) bool {
	seen := make(map[T]struct{}, len(value))
	for _, v := range value {
		if _, exists := seen[v]; exists {
			return false
		}
		seen[v] = struct{}{}
	}
	return true
}

// IsGenericIn checks if value is in allowed list.
func IsGenericIn[T comparable](value T, allowed []T) bool {
	for _, a := range allowed {
		if value == a {
			return true
		}
	}
	return false
}

// IsGenericNotIn checks if value is not in disallowed list.
func IsGenericNotIn[T comparable](value T, disallowed []T) bool {
	for _, d := range disallowed {
		if value == d {
			return false
		}
	}
	return true
}
