package validator

import (
	"reflect"
	"sync"
)

// customTypeExtractFunc extracts the inner value from a wrapper type.
// Returns the unwrapped reflect.Value and whether the field should be validated.
type customTypeExtractFunc func(v reflect.Value) (reflect.Value, bool)

// TypeMatchFunc determines if a reflect.Type should be auto-unwrapped.
// Called once per distinct type; the result is cached.
type TypeMatchFunc func(t reflect.Type) bool

type autoUnwrapEntry struct {
	matcher     TypeMatchFunc
	boolMethod  string // method name returning bool (e.g., "IsSet", "Valid")
	valueMethod string // method name returning inner value (e.g., "Value", "Get")
}

var (
	// Exact type match — fast O(1) lookup via sync.Map
	customTypeFuncs sync.Map // reflect.Type -> customTypeExtractFunc

	// Auto-unwrap matchers — for generic types like graphql.Omittable[T]
	autoUnwrapMatchers   []autoUnwrapEntry
	autoUnwrapMatchersMu sync.RWMutex

	// Cache for auto-unwrap resolution
	autoUnwrapCache sync.Map // reflect.Type -> customTypeExtractFunc (nil = no match)
)

// RegisterCustomTypeFunc registers a value extractor for exact wrapper type W.
// Uses reflect.TypeFor[W]() for zero-alloc type registration and direct type
// assertion in the hot path — no reflection method calls needed.
//
// Best for types with a small number of variants (e.g., sql.NullString, sql.NullInt64).
//
// Example:
//
//	validator.RegisterCustomTypeFunc(func(n sql.NullString) (any, bool) {
//	    if !n.Valid { return nil, false }
//	    return n.String, true
//	})
func RegisterCustomTypeFunc[W any](extract func(W) (any, bool)) {
	t := reflect.TypeFor[W]()
	fn := func(v reflect.Value) (reflect.Value, bool) {
		w := v.Interface().(W)
		result, ok := extract(w)
		if !ok {
			return reflect.Value{}, false
		}
		return reflect.ValueOf(result), true
	}
	customTypeFuncs.Store(t, customTypeExtractFunc(fn))
}

// RegisterAutoUnwrap registers a type matcher for automatic unwrapping.
// When a struct type matches, the validator auto-builds an optimized unwrapper
// using the specified methods with cached method indices.
//
// Parameters:
//   - matcher: called once per distinct type to check if it should be unwrapped
//   - boolMethod: method name that returns bool (true = value is set)
//   - valueMethod: method name that returns the inner value
//
// The matcher and method lookups are cached per type — they run only once.
//
// Example for gqlgen's graphql.Omittable[T]:
//
//	validator.RegisterAutoUnwrap(
//	    func(t reflect.Type) bool {
//	        return strings.HasPrefix(t.Name(), "Omittable[") &&
//	            strings.Contains(t.PkgPath(), "gqlgen")
//	    },
//	    "IsSet",  // Omittable.IsSet() bool
//	    "Value",  // Omittable.Value() T
//	)
func RegisterAutoUnwrap(matcher TypeMatchFunc, boolMethod, valueMethod string) {
	autoUnwrapMatchersMu.Lock()
	autoUnwrapMatchers = append(autoUnwrapMatchers, autoUnwrapEntry{
		matcher:     matcher,
		boolMethod:  boolMethod,
		valueMethod: valueMethod,
	})
	autoUnwrapMatchersMu.Unlock()
	// Clear cache since new matcher may match previously cached types
	autoUnwrapCache = sync.Map{}
}

// ResetCustomTypeFuncs removes all registered custom type functions and auto-unwrap matchers.
// Useful for testing or reconfiguration.
func ResetCustomTypeFuncs() {
	customTypeFuncs.Range(func(key, _ any) bool {
		customTypeFuncs.Delete(key)
		return true
	})
	autoUnwrapMatchersMu.Lock()
	autoUnwrapMatchers = nil
	autoUnwrapMatchersMu.Unlock()
	autoUnwrapCache = sync.Map{}
}

// resolveCustomTypeFunc returns the extract function for the given type.
// Priority: exact match (fastest) > auto-unwrap matcher with cached method indices.
func resolveCustomTypeFunc(t reflect.Type) customTypeExtractFunc {
	// Fast path: exact type match
	if fn, ok := customTypeFuncs.Load(t); ok {
		return fn.(customTypeExtractFunc)
	}

	// Check auto-unwrap cache
	if cached, ok := autoUnwrapCache.Load(t); ok {
		if cached == nil {
			return nil
		}
		return cached.(customTypeExtractFunc)
	}

	// Slow path: check matchers and build unwrapper
	fn := tryAutoUnwrap(t)
	// Cache result (nil means "no match, don't check again")
	autoUnwrapCache.Store(t, fn)
	return fn
}

// tryAutoUnwrap checks if a struct type matches any registered auto-unwrap matcher,
// then verifies it has the specified methods. Builds a cached unwrapper
// using pre-resolved method indices for O(1) method dispatch.
func tryAutoUnwrap(t reflect.Type) customTypeExtractFunc {
	if t.Kind() != reflect.Struct {
		return nil
	}

	autoUnwrapMatchersMu.RLock()
	defer autoUnwrapMatchersMu.RUnlock()

	for _, entry := range autoUnwrapMatchers {
		if !entry.matcher(t) {
			continue
		}

		// Verify bool method: func() bool
		boolMethod, ok := t.MethodByName(entry.boolMethod)
		if !ok {
			continue
		}
		bmt := boolMethod.Type
		if bmt.NumIn() != 1 || bmt.NumOut() != 1 || bmt.Out(0).Kind() != reflect.Bool {
			continue
		}

		// Verify value method: func() T
		valMethod, ok := t.MethodByName(entry.valueMethod)
		if !ok {
			continue
		}
		vmt := valMethod.Type
		if vmt.NumIn() != 1 || vmt.NumOut() != 1 {
			continue
		}

		// Cache method indices — Value.Method(index) is O(1) vs MethodByName is O(n)
		boolIdx := boolMethod.Index
		valIdx := valMethod.Index

		return customTypeExtractFunc(func(v reflect.Value) (reflect.Value, bool) {
			if !v.Method(boolIdx).Call(nil)[0].Bool() {
				return reflect.Value{}, false
			}
			return v.Method(valIdx).Call(nil)[0], true
		})
	}

	return nil
}
